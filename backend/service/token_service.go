package service

import (
	"context"
	"crypto/rsa"
	"di/model"
	"di/repository"
	"di/util/errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

// tokenService used for injecting an implementation of TokenRepository
// for use in service methods along with keys and secrets for
// signing JWTs
type tokenService struct {
	TokenRepository       model.TokenRepository
	PrivKey               *rsa.PrivateKey
	PubKey                *rsa.PublicKey
	RefreshSecret         string
	IDExpirationSecs      int64
	RefreshExpirationSecs int64
}

// TokenServiceConfig will hold repositories that will eventually be injected into this
// this service layer
type TokenServiceConfig struct {
	RedisClient           *redis.Client
	PrivKey               *rsa.PrivateKey
	PubKey                *rsa.PublicKey
	RefreshSecret         string
	IDExpirationSecs      int64
	RefreshExpirationSecs int64
}

// idTokenCustomClaims holds structure of jwt claims of idToken
type idTokenCustomClaims struct {
	User *model.User `json:"user"`
	jwt.StandardClaims
}

func GetTokenServiceConfig(redisClient *redis.Client) (*TokenServiceConfig, error) {

	// privKeyFile := os.Getenv("./jwt/rsa_private.pem")
	priv, err := ioutil.ReadFile("./jwt/rsa_private.pem")

	if err != nil {
		return nil, fmt.Errorf("Could not read private key pem file: %w", err)
	}

	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(priv)

	if err != nil {
		return nil, fmt.Errorf("Could not parse private key: %w", err)
	}

	// pubKeyFile := os.Getenv("./jwt/rsa_public.pem")
	pub, err := ioutil.ReadFile("./jwt/rsa_public.pem")

	if err != nil {
		return nil, fmt.Errorf("Could not read public key pem file: %w", err)
	}

	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pub)

	if err != nil {
		return nil, fmt.Errorf("Could not parse public key: %w", err)
	}

	// load refresh token secret from env variable
	tokenRefreshSecret := os.Getenv("TOKEN_REFRESH_SECRET")

	// load expiration lengths from env variables and parse as int
	idTokenDuration := os.Getenv("ID_TOKEN_DURATION")
	refreshTokenDuration := os.Getenv("REFRESH_TOKEN_DURATION")

	idExp, err := strconv.ParseInt(idTokenDuration, 0, 64)
	if err != nil {
		return nil, fmt.Errorf("Could not parse ID_TOKEN_EXP as int: %w", err)
	}

	refreshExp, err := strconv.ParseInt(refreshTokenDuration, 0, 64)
	if err != nil {
		return nil, fmt.Errorf("Could not parse REFRESH_TOKEN_EXP as int: %w", err)
	}

	return &TokenServiceConfig{
		RedisClient:           redisClient,
		PrivKey:               privKey,
		PubKey:                pubKey,
		RefreshSecret:         tokenRefreshSecret,
		IDExpirationSecs:      idExp,
		RefreshExpirationSecs: refreshExp,
	}, nil
}

// NewTokenService is a factory function for
// initializing a UserService with its repository layer dependencies
func NewTokenService(config *TokenServiceConfig) TokenService {
	return &tokenService{
		TokenRepository:       repository.NewTokenRepository(config.RedisClient),
		PrivKey:               config.PrivKey,
		PubKey:                config.PubKey,
		RefreshSecret:         config.RefreshSecret,
		IDExpirationSecs:      config.IDExpirationSecs,
		RefreshExpirationSecs: config.RefreshExpirationSecs,
	}
}

// NewPairFromUser creates fresh id and refresh tokens for the current user
// If a previous token is included, the previous token is removed from
// the tokens repository
func (service *tokenService) NewFirstPairFromUser(ctx context.Context, u *model.User) (*model.TokenPair, error) {
	// No need to use a repository for idToken as it is unrelated to any data source
	idToken, err := generateIDToken(u, service.PrivKey, service.IDExpirationSecs)

	if err != nil {
		log.Printf("Error generating idToken for uid: %v. Error: %v\n", u.ID, err.Error())
		return nil, errors.NewInternal()
	}

	refreshToken, err := generateRefreshToken(u.ID, service.RefreshSecret, service.RefreshExpirationSecs)

	if err != nil {
		log.Printf("Error generating refreshToken for uid: %v. Error: %v\n", u.ID, err.Error())
		return nil, errors.NewInternal()
	}

	// set freshly minted refresh token to valid list
	if err := service.TokenRepository.SetRefreshToken(ctx, u.ID, refreshToken.ID, refreshToken.ExpiresIn); err != nil {
		log.Printf("Error storing tokenID for uid: %v. Error: %v\n", u.ID, err.Error())
		return nil, errors.NewInternal()
	}

	return &model.TokenPair{
		IDToken:      model.IDToken{SignedString: idToken},
		RefreshToken: model.RefreshToken{SignedString: refreshToken.SignedTokenString, ID: refreshToken.ID, UID: u.ID},
	}, nil
}

// NewPairFromUser creates fresh id and refresh tokens for the current user
// If a previous token is included, the previous token is removed from
// the tokens repository
func (service *tokenService) NewPairFromUser(ctx context.Context, u *model.User, refreshToken model.RefreshToken) (*model.TokenPair, error) {
	idToken, err := generateIDToken(u, service.PrivKey, service.IDExpirationSecs)

	if err != nil {
		log.Printf("Error generating idToken for uid: %v. Error: %v\n", u.ID, err.Error())
		return nil, errors.NewInternal()
	}

	return &model.TokenPair{
		IDToken:      model.IDToken{SignedString: idToken},
		RefreshToken: model.RefreshToken{SignedString: refreshToken.SignedString, ID: refreshToken.ID, UID: refreshToken.ID},
	}, nil
}

// Signout reaches out to the repository layer to delete all valid tokens for a user
func (service *tokenService) Signout(ctx context.Context, id uint) error {
	return service.TokenRepository.DeleteUserRefreshTokens(ctx, id)
}

// ValidateIDToken validates the id token jwt string
// It returns the user extract from the IDTokenCustomClaims
func (service *tokenService) ValidateIDToken(tokenString string) (*model.User, error) {
	claims, err := validateIDToken(tokenString, service.PubKey) // uses public RSA key

	// We'll just return unauthorized error in all instances of failing to verify user
	if err != nil {
		log.Printf("Unable to validate or parse idToken - Error: %v\n", err)
		return nil, errors.NewAuthorization("Unable to verify user from idToken")
	}

	return claims.User, nil
}

// ValidateRefreshToken checks to make sure the JWT provided by a string is valid
// and returns a RefreshToken if valid
func (service *tokenService) ValidateRefreshToken(tokenString string) (*model.RefreshToken, error) {
	// validate actual JWT with string a secret
	claims, err := validateRefreshToken(tokenString, service.RefreshSecret)

	// We'll just return unauthorized error in all instances of failing to verify user
	if err != nil {
		log.Printf("Unable to validate or parse refreshToken for token string: %s\n%v\n", tokenString, err)
		return nil, errors.NewAuthorization("Unable to verify user from refresh token")
	}

	// Standard claims store ID as a string. I want "model" to be clear our string
	// is a UUID. So we parse claims.Id as UUID
	tokenUUID, err := uuid.Parse(claims.Id)

	if err != nil {
		log.Printf("Claims ID could not be parsed as UUID: %s\n%v\n", claims.Id, err)
		return nil, errors.NewAuthorization("Unable to verify user from refresh token")
	}

	return &model.RefreshToken{
		SignedString: tokenString,
		ID:           uint(tokenUUID.ID()),
		UID:          claims.UID,
	}, nil
}

// generateIDToken generates an IDToken which is a jwt with myCustomClaims
// Could call this GenerateIDTokenString, but the signature makes this fairly clear
func generateIDToken(u *model.User, key *rsa.PrivateKey, duration int64) (string, error) {
	timestamp := time.Now().Unix()
	tokenExp := timestamp + duration

	claims := idTokenCustomClaims{
		User: u,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  timestamp,
			ExpiresAt: tokenExp,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedTokenString, err := token.SignedString(key)

	if err != nil {
		log.Println("Failed to sign id token string")
		return "", err
	}

	return signedTokenString, nil
}

// refreshTokenData holds the actual signed jwt string along with the ID
// We return the id so it can be used without re-parsing the JWT from signed string
type refreshTokenData struct {
	SignedTokenString string
	ID                uint
	ExpiresIn         time.Duration
}

// refreshTokenCustomClaims holds the payload of a refresh token
// This can be used to extract user id for subsequent
// application operations (IE, fetch user in Redis)
type refreshTokenCustomClaims struct {
	UID uint `json:"uid"`
	jwt.StandardClaims
}

// generateRefreshToken creates a refresh token
// The refresh token stores only the user's ID, a string
func generateRefreshToken(uid uint, key string, duration int64) (*refreshTokenData, error) {
	timestamp := time.Now()
	tokenExp := timestamp.Add(time.Duration(duration) * time.Second)
	tokenID, err := uuid.NewRandom() // v4 uuid in the google uuid lib

	if err != nil {
		log.Println("Failed to generate refresh token ID")
		return nil, err
	}

	claims := refreshTokenCustomClaims{
		UID: uid,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  timestamp.Unix(),
			ExpiresAt: tokenExp.Unix(),
			Id:        tokenID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedTokenString, err := token.SignedString([]byte(key))

	if err != nil {
		log.Println("Failed to sign refresh token string")
		return nil, err
	}

	return &refreshTokenData{
		SignedTokenString: signedTokenString,
		ID:                uint(tokenID.ID()),
		ExpiresIn:         tokenExp.Sub(timestamp),
	}, nil
}

func validateIDToken(tokenString string, key *rsa.PublicKey) (*idTokenCustomClaims, error) {
	claims := &idTokenCustomClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("ID token is invalid")
	}

	claims, ok := token.Claims.(*idTokenCustomClaims)

	if !ok {
		return nil, fmt.Errorf("ID token valid but couldn't parse claims")
	}

	return claims, nil
}

func validateRefreshToken(tokenString string, key string) (*refreshTokenCustomClaims, error) {
	claims := &refreshTokenCustomClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	// For now we'll just return the error and handle logging in service level
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("Refresh token is invalid")
	}

	claims, ok := token.Claims.(*refreshTokenCustomClaims)

	if !ok {
		return nil, fmt.Errorf("Refresh token valid but couldn't parse claims")
	}

	return claims, nil
}

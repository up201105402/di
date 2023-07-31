package service

import (
	"context"
	"crypto/rsa"
	"di/model"
	"di/repository"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type tokenService struct {
	I18n                  *i18n.Localizer
	TokenRepository       repository.TokenRepository
	PrivKey               *rsa.PrivateKey
	PubKey                *rsa.PublicKey
	RefreshSecret         string
	IDExpirationSecs      int64
	RefreshExpirationSecs int64
}

type TokenServiceConfig struct {
	RedisClient           *redis.Client
	PrivKey               *rsa.PrivateKey
	PubKey                *rsa.PublicKey
	RefreshSecret         string
	IDExpirationSecs      int64
	RefreshExpirationSecs int64
}

type idTokenCustomClaims struct {
	User *model.User `json:"user"`
	jwt.StandardClaims
}

func GetTokenServiceConfig(redisClient *redis.Client) (*TokenServiceConfig, error) {

	priv, err := ioutil.ReadFile("./jwt/rsa_private.pem")

	if err != nil {
		return nil, fmt.Errorf("Could not read private key pem file: %w", err)
	}

	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(priv)

	if err != nil {
		return nil, fmt.Errorf("Could not parse private key: %w", err)
	}

	pub, err := ioutil.ReadFile("./jwt/rsa_public.pem")

	if err != nil {
		return nil, fmt.Errorf("Could not read public key pem file: %w", err)
	}

	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pub)

	if err != nil {
		return nil, fmt.Errorf("Could not parse public key: %w", err)
	}

	tokenRefreshSecret, exists := os.LookupEnv("TOKEN_REFRESH_SECRET")

	if !exists {
		panic("TOKEN_REFRESH_SECRET is not defined!")
	}

	idTokenDuration, exists := os.LookupEnv("ID_TOKEN_DURATION")

	if !exists {
		panic("ID_TOKEN_DURATION is not defined!")
	}

	refreshTokenDuration, exists := os.LookupEnv("REFRESH_TOKEN_DURATION")

	if !exists {
		panic("REFRESH_TOKEN_DURATION is not defined!")
	}

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

func NewTokenService(config *TokenServiceConfig, i18n *i18n.Localizer) TokenService {
	return &tokenService{
		I18n:                  i18n,
		TokenRepository:       repository.NewTokenRepository(config.RedisClient),
		PrivKey:               config.PrivKey,
		PubKey:                config.PubKey,
		RefreshSecret:         config.RefreshSecret,
		IDExpirationSecs:      config.IDExpirationSecs,
		RefreshExpirationSecs: config.RefreshExpirationSecs,
	}
}

func (service *tokenService) NewFirstPairFromUser(ctx context.Context, u *model.User) (*model.TokenPair, error) {
	idToken, err := generateIDToken(u, service.PrivKey, service.IDExpirationSecs)

	if err != nil {
		errorMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "token.service.id-token.generate.failed",
			TemplateData: map[string]interface{}{
				"UID":    u.ID,
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})
		log.Printf(errorMessage)
		return nil, errors.New(errorMessage)
	}

	refreshToken, err := generateRefreshToken(u.ID, service.PrivKey, service.RefreshExpirationSecs)

	if err != nil {
		errorMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "token.service.refresh-token.generate.failed",
			TemplateData: map[string]interface{}{
				"UID":    u.ID,
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})
		log.Printf(errorMessage)
		return nil, errors.New(errorMessage)
	}

	if err := service.TokenRepository.SetRefreshToken(ctx, u.ID, refreshToken.ID, refreshToken.ExpiresIn); err != nil {
		errorMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "token.service.id-token.store.failed",
			TemplateData: map[string]interface{}{
				"UID":    u.ID,
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})
		log.Printf(errorMessage)
		return nil, errors.New(errorMessage)
	}

	return &model.TokenPair{
		IDToken:      model.IDToken{SignedString: idToken},
		RefreshToken: model.RefreshToken{SignedString: refreshToken.SignedTokenString, ID: refreshToken.ID, UID: u.ID},
	}, nil
}

func (service *tokenService) NewPairFromUser(ctx context.Context, u *model.User, refreshToken model.RefreshToken) (*model.TokenPair, error) {
	idToken, err := generateIDToken(u, service.PrivKey, service.IDExpirationSecs)

	if err != nil {
		errorMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "token.service.id-token.generate.failed",
			TemplateData: map[string]interface{}{
				"UID":    u.ID,
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})
		log.Printf(errorMessage)
		return nil, errors.New(errorMessage)
	}

	return &model.TokenPair{
		IDToken:      model.IDToken{SignedString: idToken},
		RefreshToken: model.RefreshToken{SignedString: refreshToken.SignedString, ID: refreshToken.ID, UID: refreshToken.ID},
	}, nil
}

func (service *tokenService) Signout(ctx context.Context, id uint) error {
	return service.TokenRepository.DeleteUserRefreshTokens(ctx, id)
}

func (service *tokenService) ValidateIDToken(tokenString string) (*model.User, error) {
	claims, err := validateIDToken(tokenString, service.PubKey)

	if err != nil {
		errorMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "token.service.id-token.validate.failed",
			TemplateData: map[string]interface{}{
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})
		log.Printf(errorMessage)
		return nil, errors.New(errorMessage)
	}

	return claims.User, nil
}

func (service *tokenService) ValidateRefreshToken(tokenString string) (*model.RefreshToken, error) {
	claims, err := validateRefreshToken(tokenString, service.PubKey)

	if err != nil {
		errorMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "token.service.refresh-token.validate.failed",
			TemplateData: map[string]interface{}{
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})
		log.Printf(errorMessage)
		return nil, errors.New(errorMessage)
	}

	tokenUUID, err := uuid.Parse(claims.Id)

	if err != nil {
		errorMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "token.service.claims.parse.failed",
			TemplateData: map[string]interface{}{
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})
		log.Printf(errorMessage)
		return nil, errors.New(errorMessage)
	}

	return &model.RefreshToken{
		SignedString: tokenString,
		ID:           uint(tokenUUID.ID()),
		UID:          claims.UID,
	}, nil
}

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

type refreshTokenData struct {
	SignedTokenString string
	ID                uint
	ExpiresIn         time.Duration
}

type refreshTokenCustomClaims struct {
	UID uint `json:"uid"`
	jwt.StandardClaims
}

func generateRefreshToken(uid uint, key *rsa.PrivateKey, duration int64) (*refreshTokenData, error) {
	timestamp := time.Now()
	tokenExp := timestamp.Add(time.Duration(duration) * time.Second)
	tokenID, err := uuid.NewRandom()

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

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedTokenString, err := token.SignedString(key)

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

func validateRefreshToken(tokenString string, key *rsa.PublicKey) (*refreshTokenCustomClaims, error) {
	claims := &refreshTokenCustomClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	// For now we'll just return the error and handle logging in service level
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("")
	}

	claims, ok := token.Claims.(*refreshTokenCustomClaims)

	if !ok {
		return nil, errors.New("")
	}

	return claims, nil
}

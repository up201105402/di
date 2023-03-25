package repository

import (
	"context"
	"log"
	"fmt"

	"gorm.io/gorm"

	"github.com/google/uuid"
	"github.com/jacobsngoodwin/memrizr/account/model"
	"github.com/jacobsngoodwin/memrizr/account/model/apperrors"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// PGUserRepository is data/repository implementation
// of service layer UserRepository
type pGUserRepository struct {
	DB *sqlx.DB
}

// NewUserRepository is a factory for initializing User Repositories
func NewUserRepository(db *sqlx.DB) model.UserRepository {
	// load env variables
	pgHost := os.Getenv("POSTGRES_HOST")
	pgPort := os.Getenv("POSTGRES_PORT")
	pgUser := os.Getenv("POSTGRES_USER")
	pgPassword := os.Getenv("POSTGRES_PASSWORD")
	pgDB := os.Getenv("POSTGRES_DB")
	pgSSL := os.Getenv("POSTGRES_SSL")

	pgConnString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", pgHost, pgPort, pgUser, pgPassword, pgDB, pgSSL)
	dsn := "host=postgres user=admin password=admin dbname=di port=5432 sslmode=disable TimeZone=Europe/Lisbon"
	db, err := gorm.Open()
	if err != nil {
		panic("failed to connect database")
	}

	return &pGUserRepository{
		DB: db,
	}
}

// Create reaches out to database SQLX api
func (r *pGUserRepository) Create(ctx context.Context, u *model.User) error {
	query := "INSERT INTO users (email, password) VALUES ($1, $2) RETURNING *"

	if err := r.DB.GetContext(ctx, u, query, u.Email, u.Password); err != nil {
		// check unique constraint
		if err, ok := err.(*pq.Error); ok && err.Code.Name() == "unique_violation" {
			log.Printf("Could not create a user with email: %v. Reason: %v\n", u.Email, err.Code.Name())
			return apperrors.NewConflict("email", u.Email)
		}

		log.Printf("Could not create a user with email: %v. Reason: %v\n", u.Email, err)
		return apperrors.NewInternal()
	}
	return nil
}

// FindByID fetches user by id
func (r *pGUserRepository) FindByID(ctx context.Context, uid uuid.UUID) (*model.User, error) {
	user := &model.User{}

	query := "SELECT * FROM users WHERE uid=$1"

	// we need to actually check errors as it could be something other than not found
	if err := r.DB.GetContext(ctx, user, query, uid); err != nil {
		return user, apperrors.NewNotFound("uid", uid.String())
	}

	return user, nil
}

// FindByEmail retrieves user row by email address
func (r *pGUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	user := &model.User{}

	query := "SELECT * FROM users WHERE email=$1"

	if err := r.DB.GetContext(ctx, user, query, email); err != nil {
		log.Printf("Unable to get user with email address: %v. Err: %v\n", email, err)
		return user, apperrors.NewNotFound("email", email)
	}

	return user, nil
}

// Update updates a user's properties
func (r *pGUserRepository) Update(ctx context.Context, u *model.User) error {
	query := `
		UPDATE users 
		SET name=:name, email=:email, website=:website
		WHERE uid=:uid
		RETURNING *;
	`

	nstmt, err := r.DB.PrepareNamedContext(ctx, query)

	if err != nil {
		log.Printf("Unable to prepare user update query: %v\n", err)
		return apperrors.NewInternal()
	}

	if err := nstmt.GetContext(ctx, u, u); err != nil {
		log.Printf("Failed to update details for user: %v\n", u)
		return apperrors.NewInternal()
	}

	return nil
}

// UpdateImage is used to separately update a user's image separate from
// other account details
func (r *pGUserRepository) UpdateImage(ctx context.Context, uid uuid.UUID, imageURL string) (*model.User, error) {
	query := `
		UPDATE users 
		SET image_url=$2
		WHERE uid=$1
		RETURNING *;
	`

	// must be instantiated to scan into ref using `GetContext`
	u := &model.User{}

	err := r.DB.GetContext(ctx, u, query, uid, imageURL)

	if err != nil {
		log.Printf("Error updating image_url in database: %v\n", err)
		return nil, apperrors.NewInternal()
	}

	return u, nil
}

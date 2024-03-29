package database

import (
	"context"
	"errors"
	"time"
	"fmt"
	"database/sql"

	"github.com/Ola-Daniel/qrcodebakery/assets"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/lib/pq"
)

const defaultTimeout = 3 * time.Second

type DB struct {
	*sqlx.DB
}

func New(dsn string, automigrate bool) (*DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	db, err := sqlx.ConnectContext(ctx, "postgres", "postgres://"+dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(2 * time.Hour)

	if automigrate {
		iofsDriver, err := iofs.New(assets.EmbeddedFiles, "migrations")
		if err != nil {
			return nil, err
		}

		migrator, err := migrate.NewWithSourceInstance("iofs", iofsDriver, "postgres://"+dsn)
		if err != nil {
			return nil, err
		}

		err = migrator.Up()
		switch {
		case errors.Is(err, migrate.ErrNoChange):
			break
		case err != nil:
			return nil, err
		}
	}

	return &DB{db}, nil
}






// UserNotFoundError represents an error when a user is not found.
type UserNotFoundError struct {
    UsernameOrEmail string
}



// Error returns the error message.
func (e *UserNotFoundError) Error() string {
    return fmt.Sprintf("user not found: %s", e.UsernameOrEmail)
}

var ErrUserNotFound = &UserNotFoundError{}


type User struct {
	ID            int     `db:"id"`
	Username      string  `db:"username"`
	Password_hash string  `db:"password_hash"`
	Email         *string `db:"email"`
}



func (db *DB) NewUser(username string, password_hash string, email string) error {
	_, err := db.Exec("INSERT INTO users (username, password_hash, email) VALUES ($1, $2, $3)", username, password_hash, email)
	return err
}


func (db *DB) GetUser(usernameOrEmail string) (*User, error) {
    var user User
    query := "SELECT id, username, password_hash, email FROM users WHERE username = $1 OR email = $2 LIMIT 1;"
    err := db.QueryRow(query, usernameOrEmail, usernameOrEmail).Scan(&user.ID, &user.Username, &user.Password_hash, &user.Email)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, ErrUserNotFound 
        }
        return nil, err
    }
    return &user, nil
}

// represents an error when a QR code is not found
type QRCodeNotFoundError struct {
	ID int
}


// Error returns the error message

func (e *QRCodeNotFoundError) Error() string {
	return "QR code not found with ID: " + fmt.Sprint(e.ID)
}


var ErrQRCodeNotFound = &QRCodeNotFoundError{}

type QRCode struct {
	ID         int        `db:"id"`
	UserID     int        `db:"user_id"`
	Data       int        `db:"data"`
	CreatedAt  time.Time  `db:"created_at"`
}


func (db *DB) CreateQRCode(userID int, data string) error {
	_, err := db.Exec("INSERT INTO qr_codes (user_id, data) VALUES ($1, $2)", userID, data)
	return err
}

func (db *DB) GetQRCodeByID(id int) (*QRCode, error) {
	var qrCode QRCode
	query := "SELECT id, user_id, data, created_at FROM qr_codes WHERE id = $1 LIMIT 1;"
	err := db.QueryRow(query, id).Scan(&qrCode.ID, &qrCode.UserID, &qrCode.Data, &qrCode.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrQRCodeNotFound
		}
		return nil, err
	}
	return &qrCode, nil
}


func (db *DB) UpdateQRCode(id int, data string) error {
	_, err := db.Exec("UPDATE qr_codes SET data = $1 WHERE id = $2", data, id)
	return err
}

func (db *DB) DeleteQRCode(id int) error {
	_, err := db.Exec("DELETE FROM qr_codes WHERE id = $1", id)
	return err
}




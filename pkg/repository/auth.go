package repository

import (
	"database/sql"
	"tracker/pkg/model"
)

type AuthRepository interface {
	GetToken(token string) (model.Token, error)
	CreateToken(token string, userEmail string) (int64, error)
}

type authRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) AuthRepository {
	return &authRepository{db}
}

func (r *authRepository) GetToken(tokenStr string) (model.Token, error) {
	query := `SELECT ID, Uuid, Token, UserEmail, Created FROM token WHERE Token = ?`
	var token model.Token

	row := r.db.QueryRow(query, tokenStr)
	err := row.Scan(
		&token.ID,
		&token.Uuid,
		&token.Token,
		&token.UserEmail,
		&token.Created,
	)

	if err != nil {
		return token, err
	}

	return token, nil
}

func (r *authRepository) CreateToken(token, userEmail string) (int64, error) {
	query := `INSERT INTO token (Token, UserEmail) VALUES (?, ?)`
	result, err := r.db.Exec(query, token, userEmail)

	if err != nil {
		return -1, err
	}

	lastInsertId, err := result.LastInsertId()

	if err != nil {
		return -1, err
	}

	return lastInsertId, nil
}

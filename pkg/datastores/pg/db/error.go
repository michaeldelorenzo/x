package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"
)

// https://www.postgresql.org/docs/current/errcodes-appendix.html
const (
	pgAuthenticationErrorCode                     = "28P01"
	pgUniqueViolationErrorCode                    = "23505"
	pgForeignKeyViolationErrorCode                = "23503"
	pgDataExceptionErrorCodePrefix                = "22"
	pgIntegrityConstraintViolationErrorCodePrefix = "23"
	pgCardinalityViolationErrorCode               = "21000"
	pgCheckViolationErrorCode                     = "23514"
	pgCanceledByUserRequestErrorCode              = "57014"
)

func NewDBError(err error, message string) ErrDB {
	var code, constraint string

	var pgErr *pq.Error
	if errors.As(err, &pgErr) {
		code = string(pgErr.Code)
		constraint = pgErr.Constraint
	}

	return ErrDB{
		err:        err,
		message:    message,
		code:       code,
		constraint: constraint,
	}

}

// ErrDB represents an error returned from postgres
type ErrDB struct {
	err        error
	code       string
	message    string
	constraint string
}

// Error satisfies the error interface
func (e ErrDB) Error() string {
	if e.err == nil {
		return e.message
	} else {
		return fmt.Errorf("%s: %w", e.message, e.err).Error()
	}
}

// Unwrap makes the error unwappable by the error lib
func (e ErrDB) Unwrap() error {
	return e.err
}

// IsAuthError checks if the error is an authentication error
func (e ErrDB) IsAuthError() bool {
	return e.code == pgAuthenticationErrorCode
}

// IsUniqueViolationError checks if the error is a unique violation error
func (e ErrDB) IsUniqueViolationError() bool {
	return e.code == pgUniqueViolationErrorCode
}

// IsValidationError checks if the error is not a unique violation error but instead a data or constraint error
func (e ErrDB) IsValidationError() bool {
	return !e.IsUniqueViolationError() && (strings.HasPrefix(e.code, pgDataExceptionErrorCodePrefix) || strings.HasPrefix(e.code, pgIntegrityConstraintViolationErrorCodePrefix))
}

// IsForeignKeyError checks if the error is a foreign key error, such as a linked row not existing
func (e ErrDB) IsForeignKeyError() bool {
	return e.code == pgForeignKeyViolationErrorCode
}

// IsCardinalityError checks if the error was caused by trying to update/modify the same row twice
func (e ErrDB) IsCardinalityError() bool {
	return e.code == pgCardinalityViolationErrorCode
}

// IsCheckError checks if the error was caused by failing a check constraint
func (e ErrDB) IsCheckError() bool {
	return e.code == pgCheckViolationErrorCode
}

// IsCanceledByUserRequestError checks if the error was caused by by a cancelled request
func (e ErrDB) IsCanceledByUserRequestError() bool {
	return e.code == pgCanceledByUserRequestErrorCode
}

func (e ErrDB) IsNoRowsError() bool {
	return errors.Is(e.err, sql.ErrNoRows)
}

// GetAffected returns a string representing what was affected (column or entity) based on the type
func (e ErrDB) GetAffected() string {
	switch {
	case e.IsForeignKeyError():
		// assumes foreign key constraints are named `entity_fkey` or `entity_fk“
		return strings.TrimSuffix(strings.TrimSuffix(e.constraint, "_fkey"), "_fk")
	case e.IsCheckError():
		return e.constraint
	case e.IsUniqueViolationError():
		return e.constraint
	default:
		return ""
	}
}

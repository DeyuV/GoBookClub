package users

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

func NewRepository(db *pgxpool.Pool) Repository {
	return &repositoryImpl{db: db}
}

type repositoryImpl struct {
	db *pgxpool.Pool
}

// VerifyUser verify if user exists
func (r *repositoryImpl) VerifyUser(ctx context.Context, userId uint) error {

	err := r.db.QueryRow(ctx, "SELECT * FROM users WHERE user_id = $1", userId)

	if err != nil {
		return errors.New("user doesn't exist")
	}
	return nil
}

func (r *repositoryImpl) GetUserById(ctx context.Context, userId uint) (*UserResponse, error) {
	var user UserResponse

	err := r.db.QueryRow(ctx, "SELECT * FROM users WHERE user_id = $1", userId).Scan(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *repositoryImpl) Create(ctx context.Context, FirstName, LastName, Email, PasswordHash string, Enabled bool, Role string) (*UserResponse, error) {
	/**
	TODO : use transactions

	BEGIN;
	INSERT INTO "users" ("first_name","last_name","email","password","enabled","role") VALUES ('Bogdan','Dinu','badu@badu.ro','123456',true,'ADMIN') RETURNING "user_id";
	COMMIT;

	-- on error
	ROLLBACK;

	Badu's way > (we can approach this later)

	// Transaction helper function, so developer won't forget to commit the transaction, as it usually happens
	func (m *TxDB) Transaction(ctx context.Context, opts pgx.TxOptions, fn func(sqlTx pgx.Tx) error) error {
		tx, err := m.BeginTx(ctx, opts)
		if err != nil {
			return err
		}

		// calling the dev function to perform operations in this transaction
		if err := fn(tx); err != nil {
			rollbackErr := tx.Rollback(ctx)
			if rollbackErr != nil {
				// we've failed to rollback - more important than failing queries from transaction
				return rollbackErr
			}

			return err
		}

		// finally, we commit the transaction
		return tx.Commit(ctx)
	}

	*/
	newUser := UserResponse{
		ID:        0,
		FirstName: FirstName,
		LastName:  LastName,
		Email:     Email,
		Password:  PasswordHash,
		Enabled:   Enabled,
		Role:      Role,
	}

	err := r.db.QueryRow(ctx, "INSERT INTO users(first_name, last_name, email, passw, en, rol) VALUES ($1,$2,$3,$4,$5,$6) RETURNING user_id",
		FirstName, LastName, Email, PasswordHash, Enabled, Role).Scan(&newUser.ID)
	if err != nil {
		log.Printf(err.Error())
		return nil, err
	}

	log.Printf("Created user with id : %d", newUser.ID)
	return &newUser, nil
}

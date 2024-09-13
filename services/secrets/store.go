package secrets

import (
	"database/sql"
	"fmt"
	"pm4devs-backend/types"
	"time"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (pg *Store) CreateSecret(secret *types.Secret) (int, error) {
	query := `
  INSERT INTO Secret (user_id, secret_type, encrypted_data, description)
  VALUES ($1, $2, $3, $4) RETURNING secret_id;
  `

	var secretId int
	err := pg.db.QueryRow(query, secret.UserID, secret.SecretType, secret.EncryptedData, secret.Description).Scan(&secretId)

	if err != nil {
		return secretId, err
	}
	return secretId, nil
}

func (pg *Store) GetAllSecret() ([]*types.Secret, error) {
	secrets := []*types.Secret{}

	query := `
  SELECT *
  FROM Secret;
  `
	rows, err := pg.db.Query(query)
	if err != nil {
		return []*types.Secret{}, err
	}
	for rows.Next() {
		secret := new(types.Secret)
		err := rows.Scan(
			&secret.SecretID,
			&secret.UserID,
			&secret.SecretType,
			&secret.EncryptedData,
			&secret.Description,
			&secret.CreatedAt,
			&secret.UpdatedAt,
		)
		if err != nil {
			return []*types.Secret{}, err
		}
		secrets = append(secrets, secret)
	}

	return secrets, nil
}

func (pg *Store) GetSecretById(secretId int) (*types.Secret, error) {
	query := `
  SELECT *
  FROM Secret
  WHERE secret_id = $1;
  `
	secret := new(types.Secret)
	err := pg.db.QueryRow(query, secretId).Scan(
		&secret.SecretID,
		&secret.UserID,
		&secret.SecretType,
		&secret.EncryptedData,
		&secret.Description,
		&secret.CreatedAt,
		&secret.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return secret, nil
}

func (pg *Store) GetAllSecretsByUserID(userId int) ([]*types.Secret, error) {
	secrets := []*types.Secret{}
	query := `
  SELECT *
  FROM Secret
  WHERE user_id = $1;
  `
	rows, err := pg.db.Query(query, userId)
	if err != nil {
		return []*types.Secret{}, err
	}
	for rows.Next() {
		secret := new(types.Secret)
		err := rows.Scan(
			&secret.SecretID,
			&secret.UserID,
			&secret.SecretType,
			&secret.EncryptedData,
			&secret.Description,
			&secret.CreatedAt,
			&secret.UpdatedAt,
		)
		if err != nil {
			return []*types.Secret{}, err
		}
		secrets = append(secrets, secret)
	}

	return secrets, nil
}

func (pg *Store) DeleteSecretById(secretId int) error {

	query := `
  DELETE FROM Secret 
  WHERE secret_id = $1; 
  `

	result, err := pg.db.Exec(query, secretId)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("secret with ID %d not found", secretId)
	}
	return nil
}

func (pg *Store) UpdateSecretById(secretId int, data types.UpdateSecretPayload) error {
	query := `
    UPDATE Secret
    SET encrypted_data = $1, description = $2, updated_at = $3
    WHERE secret_id = $4;
  `
	result, err := pg.db.Exec(query, data.EncryptedData, data.Description, time.Now(), secretId)
	if err != nil {
		return err
	}

	// Check how many rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("secret with ID %d not found", secretId)
	}

	return nil
}

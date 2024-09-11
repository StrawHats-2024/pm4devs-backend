package db

import (
	"pm4devs-backend/pkg/models"
)

func (pg *PostgresStore) CreateSecret(secret *models.Secret) (int, error) {
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

func (pg *PostgresStore) GetAllSecret() ([]*models.Secret, error) {
	secrets := []*models.Secret{}

	query := `
  SELECT *
  FROM Secret;
  `
	rows, err := pg.db.Query(query)
	if err != nil {
		return []*models.Secret{}, err
	}
	for rows.Next() {
		secret := new(models.Secret)
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
			return []*models.Secret{}, err
		}
		secrets = append(secrets, secret)
	}

	return secrets, nil
}

func (pg *PostgresStore) GetSecretById(secretId int) (*models.Secret, error) {
	query := `
  SELECT *
  FROM Secret
  WHERE secret_id = $1;
  `
	secret := new(models.Secret)
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

func (pg *PostgresStore) GetAllSecretsByUserID(userId int) ([]*models.Secret, error) {
	secrets := []*models.Secret{}
	query := `
  SELECT *
  FROM Secret
  WHERE user_id = $1;
  `
	rows, err := pg.db.Query(query, userId)
	if err != nil {
		return []*models.Secret{}, err
	}
	for rows.Next() {
		secret := new(models.Secret)
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
			return []*models.Secret{}, err
		}
		secrets = append(secrets, secret)
	}

	return secrets, nil
}

func (pg *PostgresStore) DeleteSecretById(secretId int) error {

	query := `
  DELETE FROM Secret 
  WHERE secret_id = $1; 
  `

	_, err := pg.db.Exec(query, secretId)
	if err != nil {
		return err
	}
	return nil
}

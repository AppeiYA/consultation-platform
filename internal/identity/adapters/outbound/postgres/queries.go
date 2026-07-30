package postgres

const (
	CREATE_USER = `
		INSERT INTO users (id, first_name, last_name, email, password_hash, role, is_verified, is_deleted, created_at, updated_at) 
		VALUES (:id, :first_name, :last_name, :email, :password_hash, :role, :is_verified, :is_deleted, :created_at, :updated_at)
	`
	FIND_USER_BY_ID = `
		SELECT id, first_name, last_name, email, password_hash, role, is_verified, is_deleted, created_at, updated_at 
		FROM users WHERE id = $1
	`
	FIND_USER_BY_EMAIL = `
		SELECT id, first_name, last_name, email, password_hash, role, is_verified, is_deleted, created_at, updated_at 
		FROM users WHERE email = $1
	`
	UPDATE_USER = `
		UPDATE users
		SET
		first_name = :first_name,
		last_name = :last_name,
		updated_at = :updated_at
		WHERE id = :id
	`
	DELETE_USER = `
		UPDATE users SET is_deleted = true, updated_at = :updated_at WHERE id = :id
	`
	VERIFY_USER = `
		UPDATE users SET is_verified = true, updated_at = :updated_at WHERE id = :id
	`
	RESTORE_USER = `
		UPDATE users SET is_deleted = false, updated_at = :updated_at WHERE id = :id
	`
)
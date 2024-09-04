-- Drop triggers
DROP TRIGGER update_updated_at_trigger_refresh_tokens ON refresh_tokens;
DROP TRIGGER update_updated_at_trigger_user_passwords ON user_passwords;
DROP TRIGGER update_updated_at_trigger ON users;

-- Drop functions
DROP FUNCTION update_updated_at();

-- Drop indexes
DROP INDEX idx_refresh_tokens_expires_at;
DROP INDEX idx_refresh_tokens_token;
DROP INDEX idx_refresh_tokens_user_id;
DROP INDEX idx_user_passwords_user_id;
DROP INDEX idx_users_user_id;
DROP INDEX idx_users_email;

-- Drop tables
DROP TABLE refresh_tokens;
DROP TABLE user_passwords;
DROP TABLE users;
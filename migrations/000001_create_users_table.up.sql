-- Create the users table
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  uuid VARCHAR(255) NOT NULL UNIQUE,
  email VARCHAR(255) NOT NULL UNIQUE,
  first_name VARCHAR(255) DEFAULT NULL,
  last_name VARCHAR(255) DEFAULT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create the user_passwords table
CREATE TABLE user_passwords (
  id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  password VARCHAR(255) NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  UNIQUE (user_id)
);

CREATE TABLE refresh_tokens (
  id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  token VARCHAR(255) NOT NULL,
  expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  
  -- Ensure the token is unique for each user with only one active token (deleted_at is NULL)
  CONSTRAINT unique_active_token_per_user UNIQUE (user_id, token)
);

-- Add partial index to ensure only one active token exists per user (deleted_at is NULL)
CREATE UNIQUE INDEX idx_unique_active_refresh_token 
ON refresh_tokens (user_id) 
WHERE deleted_at IS NULL;

-- Create indexes
CREATE INDEX idx_users_email ON users (email);
CREATE INDEX idx_users_user_id ON users (uuid);
CREATE INDEX idx_user_passwords_user_id ON user_passwords (user_id);
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens (user_id);
CREATE INDEX idx_refresh_tokens_token ON refresh_tokens (token);
CREATE INDEX idx_refresh_tokens_expires_at ON refresh_tokens (expires_at);

-- Create a function to update the updated_at column
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = CURRENT_TIMESTAMP AT TIME ZONE 'UTC';
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create a trigger to update the updated_at column on update
CREATE TRIGGER update_updated_at_trigger
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at();

-- Create a trigger to update the updated_at column on update for user_passwords
CREATE TRIGGER update_updated_at_trigger_user_passwords
BEFORE UPDATE ON user_passwords
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at();

-- Create a trigger to update the updated_at column on update for refresh_tokens
CREATE TRIGGER update_updated_at_trigger_refresh_tokens
BEFORE UPDATE ON refresh_tokens
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at();
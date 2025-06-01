-- Initialize the database with proper charset and collation
CREATE DATABASE IF NOT EXISTS notes_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- Create a dedicated user for the application (if not exists)
CREATE USER IF NOT EXISTS 'notes_user'@'%' IDENTIFIED BY 'notes_password';

-- Grant necessary privileges
GRANT SELECT, INSERT, UPDATE, DELETE, CREATE, ALTER, INDEX, DROP ON notes_db.* TO 'notes_user'@'%';

-- Flush privileges to ensure they take effect
FLUSH PRIVILEGES;

-- Use the notes database
USE notes_db;

-- Create indexes for better performance (GORM will create tables)
-- These will be applied after GORM creates the tables

-- Note: GORM will handle table creation, but we can add additional indexes here if needed
-- The following commands will be executed after the tables are created by GORM

-- Additional indexes for better query performance
-- ALTER TABLE users ADD INDEX idx_users_email (email);
-- ALTER TABLE users ADD INDEX idx_users_created_at (created_at);
-- ALTER TABLE notes ADD INDEX idx_notes_user_id_created_at (user_id, created_at);
-- ALTER TABLE notes ADD INDEX idx_notes_title (title);
-- ALTER TABLE notes ADD FULLTEXT INDEX idx_notes_search (title, content);

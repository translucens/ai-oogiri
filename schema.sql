CREATE DATABASE IF NOT EXISTS handson CHARACTER SET = 'utf8mb4';

USE handson

CREATE TABLE IF NOT EXISTS riddles (
    id INT AUTO_INCREMENT PRIMARY KEY,
    theme VARCHAR(255) NOT NULL,
    primary_answer VARCHAR(255) NOT NULL,
    secondary_answer VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX created_at_idx (created_at DESC)
);

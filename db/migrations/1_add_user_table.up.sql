DROP TABLE IF EXISTS user;
CREATE TABLE user (
    ID              INT AUTO_INCREMENT NOT NULL,
    Uuid            VARCHAR(36) DEFAULT (uuid()),
    Email           VARCHAR(255) NOT NULL UNIQUE,
    Username        VARCHAR(255) NOT NULL UNIQUE,
    Password        VARCHAR(255) NOT NULL,
    Created         DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY     (ID)
);
-- up
CREATE TABLE users (
    id uuid PRIMARY KEY NOT NULL, 
    username VARCHAR(255) UNIQUE NOT NULL,
    moniker VARCHAR(255) DEFAULT "username", -- for fancy customization purposes. 
    password VARCHAR(255) NOT NULL,
    recovery_hash TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL
);


-- down
DROP TABLE users
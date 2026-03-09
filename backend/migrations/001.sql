CREATE TABLE users (
    id uuid PRIMARY KEY NOT NULL, 
    username VARCHAR(255) UNIQUE NOT NULL,
    moniker VARCHAR(255) NOT NULL, -- for fancy customization purposes. 
    password VARCHAR(255) NOT NULL,
    recovery_hash TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL
);


CREATE TABLE user_stats (
  user_id uuid PRIMARY KEY NOT NULL,
  games_played INTEGER DEFAULT 0,
  wins_count INTEGER DEFAULT 0,
  best_score INTEGER DEFAULT 0,
  total_score INTEGER DEFAULT 0,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL
);
package database

import (
	"database/sql"
	"fmt"
)

const schema = `
CREATE TABLE IF NOT EXISTS songs (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	title TEXT NOT NULL,
	artist TEXT NOT NULL,
	duration INTEGER NOT NULL,
	location TEXT NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS queue (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	song_id INTEGER NOT NULL,
	position INTEGER NOT NULL,
	added_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (song_id) REFERENCES songs(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS history (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	song_id INTEGER NOT NULL,
	played_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	duration_played INTEGER NOT NULL,
	FOREIGN KEY (song_id) REFERENCES songs(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS state (
	key TEXT PRIMARY KEY,
	value TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_queue_position ON queue(position);
CREATE INDEX IF NOT EXISTS idx_history_played_at ON history(played_at);
`

func RunMigrations(db *sql.DB) error {
	_, err := db.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	return nil
}

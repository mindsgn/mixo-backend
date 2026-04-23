package queue

import (
	"database/sql"
	"fmt"
	"sync"
)

type Song struct {
	ID       int
	Title    string
	Artist   string
	Duration int
	Location string
}

type QueueItem struct {
	ID     int
	Song   Song
	Position int
}

type Manager struct {
	db     *sql.DB
	mu     sync.RWMutex
}

func New(db *sql.DB) *Manager {
	return &Manager{db: db}
}

func (m *Manager) Add(songID int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	tx, err := m.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Get current max position
	var maxPos int
	err = tx.QueryRow("SELECT COALESCE(MAX(position), 0) FROM queue").Scan(&maxPos)
	if err != nil {
		return fmt.Errorf("failed to get max position: %w", err)
	}

	// Insert at end of queue
	_, err = tx.Exec("INSERT INTO queue (song_id, position) VALUES (?, ?)", songID, maxPos+1)
	if err != nil {
		return fmt.Errorf("failed to add to queue: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (m *Manager) Remove(id int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	tx, err := m.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Get position of item to remove
	var position int
	err = tx.QueryRow("SELECT position FROM queue WHERE id = ?", id).Scan(&position)
	if err == sql.ErrNoRows {
		return fmt.Errorf("queue item not found")
	}
	if err != nil {
		return fmt.Errorf("failed to get queue item position: %w", err)
	}

	// Delete the item
	_, err = tx.Exec("DELETE FROM queue WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete from queue: %w", err)
	}

	// Update positions of remaining items
	_, err = tx.Exec("UPDATE queue SET position = position - 1 WHERE position > ?", position)
	if err != nil {
		return fmt.Errorf("failed to update queue positions: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (m *Manager) GetNext() (*Song, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	tx, err := m.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Get first item in queue
	var queueID, songID int
	err = tx.QueryRow("SELECT id, song_id FROM queue ORDER BY position ASC LIMIT 1").Scan(&queueID, &songID)
	if err == sql.ErrNoRows {
		return nil, nil // Queue is empty
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get next song: %w", err)
	}

	// Get song details
	var song Song
	err = tx.QueryRow("SELECT id, title, artist, duration, location FROM songs WHERE id = ?", songID).
		Scan(&song.ID, &song.Title, &song.Artist, &song.Duration, &song.Location)
	if err != nil {
		return nil, fmt.Errorf("failed to get song details: %w", err)
	}

	// Remove from queue
	_, err = tx.Exec("DELETE FROM queue WHERE id = ?", queueID)
	if err != nil {
		return nil, fmt.Errorf("failed to remove from queue: %w", err)
	}

	// Update positions
	_, err = tx.Exec("UPDATE queue SET position = position - 1")
	if err != nil {
		return nil, fmt.Errorf("failed to update queue positions: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &song, nil
}

func (m *Manager) GetAll() ([]QueueItem, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	rows, err := m.db.Query(`
		SELECT q.id, q.position, s.id, s.title, s.artist, s.duration, s.location
		FROM queue q
		JOIN songs s ON q.song_id = s.id
		ORDER BY q.position ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to get queue: %w", err)
	}
	defer rows.Close()

	var items []QueueItem
	for rows.Next() {
		var item QueueItem
		err := rows.Scan(&item.ID, &item.Position, &item.Song.ID, &item.Song.Title, &item.Song.Artist, &item.Song.Duration, &item.Song.Location)
		if err != nil {
			return nil, fmt.Errorf("failed to scan queue item: %w", err)
		}
		items = append(items, item)
	}

	return items, nil
}

func (m *Manager) Length() (int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var count int
	err := m.db.QueryRow("SELECT COUNT(*) FROM queue").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get queue length: %w", err)
	}
	return count, nil
}

package queue

import (
	"database/sql"
	"os"
	"testing"
	"github.com/mindsgn-studio/mixo-backend/internal/database"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	tmpDB := "/tmp/test_queue.db"

	db, err := database.New(tmpDB)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	
	t.Cleanup(func() {
		db.Close()
		os.Remove(tmpDB)
	})
	
	return db.DB
}

func addTestSong(t *testing.T, db *sql.DB, title, artist, location string, duration int) int {
	result, err := db.Exec("INSERT INTO songs (title, artist, duration, location) VALUES (?, ?, ?, ?)",
		title, artist, duration, location)
	if err != nil {
		t.Fatalf("Failed to add test song: %v", err)
	}
	id, _ := result.LastInsertId()
	return int(id)
}

func TestManager_Add(t *testing.T) {
	db := setupTestDB(t)

	songID := addTestSong(t, db, "Test Song", "Test Artist", "/test/path.mp3", 180)

	qm := New(db)
	err := qm.Add(songID)
	if err != nil {
		t.Fatalf("Failed to add to queue: %v", err)
	}

	items, err := qm.GetAll()
	if err != nil {
		t.Fatalf("Failed to get queue: %v", err)
	}

	if len(items) != 1 {
		t.Errorf("Expected 1 item in queue, got %d", len(items))
	}

	if items[0].Song.ID != songID {
		t.Errorf("Expected song ID %d, got %d", songID, items[0].Song.ID)
	}
}

func TestManager_AddMultiple(t *testing.T) {
	db := setupTestDB(t)

	songID1 := addTestSong(t, db, "Song 1", "Artist 1", "/path1.mp3", 180)
	songID2 := addTestSong(t, db, "Song 2", "Artist 2", "/path2.mp3", 200)
	songID3 := addTestSong(t, db, "Song 3", "Artist 3", "/path3.mp3", 150)

	qm := New(db)
	
	qm.Add(songID1)
	qm.Add(songID2)
	qm.Add(songID3)

	items, err := qm.GetAll()
	if err != nil {
		t.Fatalf("Failed to get queue: %v", err)
	}

	if len(items) != 3 {
		t.Errorf("Expected 3 items in queue, got %d", len(items))
	}

	// Check FIFO order
	if items[0].Song.ID != songID1 {
		t.Errorf("Expected first song to be %d, got %d", songID1, items[0].Song.ID)
	}
	if items[1].Song.ID != songID2 {
		t.Errorf("Expected second song to be %d, got %d", songID2, items[1].Song.ID)
	}
	if items[2].Song.ID != songID3 {
		t.Errorf("Expected third song to be %d, got %d", songID3, items[2].Song.ID)
	}
}

func TestManager_Remove(t *testing.T) {
	db := setupTestDB(t)

	songID1 := addTestSong(t, db, "Song 1", "Artist 1", "/path1.mp3", 180)
	songID2 := addTestSong(t, db, "Song 2", "Artist 2", "/path2.mp3", 200)
	songID3 := addTestSong(t, db, "Song 3", "Artist 3", "/path3.mp3", 150)

	qm := New(db)
	qm.Add(songID1)
	qm.Add(songID2)
	qm.Add(songID3)

	items, _ := qm.GetAll()
	queueID := items[1].ID // Remove middle item

	err := qm.Remove(queueID)
	if err != nil {
		t.Fatalf("Failed to remove from queue: %v", err)
	}

	items, err = qm.GetAll()
	if err != nil {
		t.Fatalf("Failed to get queue: %v", err)
	}

	if len(items) != 2 {
		t.Errorf("Expected 2 items in queue, got %d", len(items))
	}

	// Check positions are updated
	if items[0].Position != 1 {
		t.Errorf("Expected first item position to be 1, got %d", items[0].Position)
	}
	if items[1].Position != 2 {
		t.Errorf("Expected second item position to be 2, got %d", items[1].Position)
	}
}

func TestManager_GetNext(t *testing.T) {
	db := setupTestDB(t)

	songID1 := addTestSong(t, db, "Song 1", "Artist 1", "/path1.mp3", 180)
	songID2 := addTestSong(t, db, "Song 2", "Artist 2", "/path2.mp3", 200)

	qm := New(db)
	qm.Add(songID1)
	qm.Add(songID2)

	// Get first song
	song, err := qm.GetNext()
	if err != nil {
		t.Fatalf("Failed to get next song: %v", err)
	}

	if song == nil {
		t.Fatal("Expected song, got nil")
	}

	if song.ID != songID1 {
		t.Errorf("Expected song ID %d, got %d", songID1, song.ID)
	}

	// Check queue has one item left
	items, _ := qm.GetAll()
	if len(items) != 1 {
		t.Errorf("Expected 1 item in queue, got %d", len(items))
	}

	// Get second song
	song, err = qm.GetNext()
	if err != nil {
		t.Fatalf("Failed to get next song: %v", err)
	}

	if song.ID != songID2 {
		t.Errorf("Expected song ID %d, got %d", songID2, song.ID)
	}

	// Queue should be empty
	items, _ = qm.GetAll()
	if len(items) != 0 {
		t.Errorf("Expected 0 items in queue, got %d", len(items))
	}
}

func TestManager_GetNext_Empty(t *testing.T) {
	db := setupTestDB(t)

	qm := New(db)
	song, err := qm.GetNext()
	if err != nil {
		t.Fatalf("Failed to get next song: %v", err)
	}

	if song != nil {
		t.Error("Expected nil song for empty queue, got song")
	}
}

func TestManager_Length(t *testing.T) {
	db := setupTestDB(t)

	qm := New(db)

	length, err := qm.Length()
	if err != nil {
		t.Fatalf("Failed to get queue length: %v", err)
	}

	if length != 0 {
		t.Errorf("Expected length 0, got %d", length)
	}

	songID := addTestSong(t, db, "Test Song", "Test Artist", "/test.mp3", 180)
	qm.Add(songID)

	length, err = qm.Length()
	if err != nil {
		t.Fatalf("Failed to get queue length: %v", err)
	}

	if length != 1 {
		t.Errorf("Expected length 1, got %d", length)
	}
}

package playback

import (
	"database/sql"
	"io"
	"log"
	"github.com/mindsgn-studio/mixo-backend/internal/queue"
	"sync"
	"time"
)

type Engine struct {
	db         *sql.DB
	queue      *queue.Manager
	chunkChan  chan []byte
	currentSong *queue.Song
	mu         sync.RWMutex
	running    bool
	stopChan   chan struct{}
}

func New(db *sql.DB, q *queue.Manager) *Engine {
	return &Engine{
		db:        db,
		queue:     q,
		chunkChan: make(chan []byte, 100),
		stopChan:  make(chan struct{}),
	}
}

func (e *Engine) Start() {
	e.mu.Lock()
	if e.running {
		e.mu.Unlock()
		return
	}
	e.running = true
	e.mu.Unlock()

	go e.playbackLoop()
}

func (e *Engine) Stop() {
	e.mu.Lock()
	defer e.mu.Unlock()
	
	if e.running {
		close(e.stopChan)
		e.running = false
	}
}

func (e *Engine) playbackLoop() {
	for {
		select {
		case <-e.stopChan:
			return
		default:
			song, err := e.queue.GetNext()
			if err != nil {
				log.Printf("Error getting next song: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}

			if song == nil {
				// Queue is empty, wait
				time.Sleep(1 * time.Second)
				continue
			}

			e.setCurrentSong(song)
			e.playSong(song)
			e.addToHistory(song, song.Duration)
		}
	}
}

func (e *Engine) playSong(song *queue.Song) {
	streamer, err := NewFFmpegStreamer(song.Location)
	if err != nil {
		log.Printf("Error creating FFmpeg streamer: %v", err)
		return
	}
	defer streamer.Close()

	buffer := make([]byte, 4096)
	startTime := time.Now()

	for {
		select {
		case <-e.stopChan:
			return
		default:
			n, err := streamer.Read(buffer)
			if err != nil {
				if err == io.EOF {
					return
				}
				log.Printf("Error reading from stream: %v", err)
				return
			}

			// Real-time throttling
			elapsed := time.Since(startTime)
			expectedDuration := time.Duration(song.Duration) * time.Second
			if elapsed < expectedDuration {
				time.Sleep(10 * time.Millisecond)
			}

			chunk := make([]byte, n)
			copy(chunk, buffer[:n])
			
			select {
			case e.chunkChan <- chunk:
			case <-e.stopChan:
				return
			}
		}
	}
}

func (e *Engine) setCurrentSong(song *queue.Song) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.currentSong = song
	
	// Update state in database
	_, err := e.db.Exec("INSERT OR REPLACE INTO state (key, value) VALUES ('current_song', ?)", song.ID)
	if err != nil {
		log.Printf("Error updating current song state: %v", err)
	}
}

func (e *Engine) GetCurrentSong() *queue.Song {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.currentSong
}

func (e *Engine) GetChunkChan() <-chan []byte {
	return e.chunkChan
}

func (e *Engine) addToHistory(song *queue.Song, durationPlayed int) {
	_, err := e.db.Exec("INSERT INTO history (song_id, duration_played) VALUES (?, ?)", song.ID, durationPlayed)
	if err != nil {
		log.Printf("Error adding to history: %v", err)
	}
}

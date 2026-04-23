# Go Internet Radio Streaming Server

A real-time internet radio streaming server built with Go that broadcasts a single continuous audio stream to all connected listeners.

## Features

- **Single continuous audio stream** - All users hear the same playback in real-time
- **No seeking or per-user control** - Traditional radio-style broadcast
- **MP3 streaming over HTTP** - Chunked transfer encoding for real-time delivery
- **Fan-out broadcaster pattern** - Efficiently serves multiple clients
- **FIFO queue management** - Songs play in order they were added
- **Admin API** - RESTful endpoints for managing songs and queue
- **SQLite persistence** - Stores songs, queue, history, and state
- **FFmpeg integration** - Ensures consistent audio format
- **Slow client handling** - Automatically drops clients that can't keep up

## Requirements

- Go 1.21 or higher
- FFmpeg installed and available in PATH
- SQLite (included via Go driver)

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd radio/backend
```

2. Install dependencies:
```bash
go mod download
```

3. Ensure FFmpeg is installed:
```bash
ffmpeg -version
```

## Configuration

Create a `.env` file in the `backend` directory:

```env
PORT=8080
SONG_DIR=../songs
DB_PATH=./radio.db
STREAM_TIMEOUT=5
```

- `PORT`: HTTP server port (default: 8080)
- `SONG_DIR`: Directory containing MP3 files
- `DB_PATH`: SQLite database file path
- `STREAM_TIMEOUT`: Timeout for slow clients in seconds

## Usage

### Start the server

```bash
cd backend
go run cmd/server/main.go
```

The server will start on the configured port with the following endpoints:
- Stream: `http://localhost:8080/stream`
- Admin API: `http://localhost:8080/api/*`
- Health check: `http://localhost:8080/health`

### Add songs

```bash
curl -X POST http://localhost:8080/api/songs \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Song Title",
    "artist": "Artist Name",
    "duration": 180,
    "location": "/path/to/song.mp3"
  }'
```

### List songs

```bash
curl http://localhost:8080/api/songs
```

### Add song to queue

```bash
curl -X POST http://localhost:8080/api/queue/{song_id}
```

### Get current queue

```bash
curl http://localhost:8080/api/queue
```

### Get now playing

```bash
curl http://localhost:8080/api/now-playing
```

### Get playback history

```bash
curl http://localhost:8080/api/history?limit=50
```

### Listen to the stream

Use any media player that supports HTTP streaming:
```bash
ffplay http://localhost:8080/stream
```

Or open in a browser:
```
http://localhost:8080/stream
```

## Project Structure

```
radio/
├── backend/
│   ├── cmd/server/
│   │   └── main.go              # Application entry point
│   ├── internal/
│   │   ├── admin/
│   │   │   ├── handler.go      # Admin API handlers
│   │   │   └── routes.go       # Route registration
│   │   ├── config/
│   │   │   └── config.go       # Configuration loading
│   │   ├── database/
│   │   │   ├── migrations.go   # Database schema
│   │   │   └── sqlite.go       # SQLite connection
│   │   ├── playback/
│   │   │   ├── engine.go       # Playback engine
│   │   │   └── ffmpeg.go       # FFmpeg integration
│   │   ├── queue/
│   │   │   └── manager.go      # Queue management
│   │   └── stream/
│   │       ├── broadcaster.go   # Fan-out broadcaster
│   │       └── handler.go       # HTTP stream handler
│   ├── .env                    # Environment variables
│   ├── go.mod                  # Go module definition
│   ├── go.sum                  # Go dependencies
│   └── article.md              # Core concepts documentation
└── songs/                      # MP3 files directory
```

## API Endpoints

### Songs
- `POST /api/songs` - Add a new song
- `GET /api/songs` - List all songs
- `DELETE /api/songs/:id` - Delete a song

### Queue
- `POST /api/queue/:songId` - Add song to queue
- `GET /api/queue` - Get current queue
- `DELETE /api/queue/:id` - Remove from queue

### Status
- `GET /api/now-playing` - Get currently playing song
- `GET /api/history` - Get playback history

### Health
- `GET /health` - Health check

## Testing

Run unit tests:
```bash
cd backend
go test ./...
```

Run tests with coverage:
```bash
go test -cover ./...
```

## Core Concepts

For detailed explanations of the core concepts and patterns used in this project, see [article.md](backend/article.md).

Topics covered:
- Fan-out broadcaster pattern
- SQLite for persistence
- FFmpeg for audio streaming
- Chunked HTTP transfer encoding
- FIFO queue management
- Slow client detection
- Real-time throttling

## Version

Current version: v0.1.0

## License

MIT License

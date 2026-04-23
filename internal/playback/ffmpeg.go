package playback

import (
	"fmt"
	"io"
	"os/exec"
	"sync"
)

type FFmpegStreamer struct {
	cmd    *exec.Cmd
	stdout io.ReadCloser
	mu     sync.Mutex
}

func NewFFmpegStreamer(filePath string) (*FFmpegStreamer, error) {
	cmd := exec.Command("ffmpeg", "-i", filePath, "-f", "mp3", "-acodec", "libmp3lame", "-ar", "44100", "-ac", "2", "-b:a", "128k", "-")
	
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start ffmpeg: %w", err)
	}

	return &FFmpegStreamer{
		cmd:    cmd,
		stdout: stdout,
	}, nil
}

func (f *FFmpegStreamer) Read(p []byte) (int, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.stdout.Read(p)
}

func (f *FFmpegStreamer) Close() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	
	if f.cmd != nil {
		if err := f.cmd.Process.Kill(); err != nil {
			return fmt.Errorf("failed to kill ffmpeg: %w", err)
		}
		f.cmd.Wait()
	}
	if f.stdout != nil {
		f.stdout.Close()
	}
	return nil
}

func (f *FFmpegStreamer) IsRunning() bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.cmd != nil && f.cmd.Process != nil
}

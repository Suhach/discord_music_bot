package player

import (
	"fmt"
	"io"
	"os/exec"
	"sync"

	"github.com/bwmarrin/discordgo"
)

type Track struct {
	Title    string
	URL      string
	Platform string
	Duration string
}

type Player struct {
	queue        []*Track
	currentTrack *Track
	voiceConn    *discordgo.VoiceConnection
	mu           sync.Mutex
	isPlaying    bool
	isPaused     bool
	stopChan     chan bool
	cmd          *exec.Cmd
}

func New() *Player {
	return &Player{
		queue:    make([]*Track, 0),
		stopChan: make(chan bool),
	}
}

// AddToQueue добавляет трек в очередь
func (p *Player) AddToQueue(track *Track) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.queue = append(p.queue, track)
}

// GetQueue возвращает текущую очередь
func (p *Player) GetQueue() []*Track {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.queue
}

// Play начинает воспроизведение
func (p *Player) Play(vc *discordgo.VoiceConnection) error {
	p.voiceConn = vc

	if p.isPaused {
		p.isPaused = false
		return nil
	}

	if p.isPlaying {
		return nil
	}

	go p.playLoop()
	return nil
}

func (p *Player) playLoop() {
	for {
		p.mu.Lock()
		if len(p.queue) == 0 {
			p.isPlaying = false
			p.currentTrack = nil
			p.mu.Unlock()
			return
		}

		track := p.queue[0]
		p.queue = p.queue[1:]
		p.currentTrack = track
		p.isPlaying = true
		p.mu.Unlock()

		if err := p.playTrack(track); err != nil {
			fmt.Printf("Ошибка воспроизведения: %v\n", err)
		}

		select {
		case <-p.stopChan:
			p.mu.Lock()
			p.isPlaying = false
			p.currentTrack = nil
			p.mu.Unlock()
			return
		default:
		}
	}
}

func (p *Player) playTrack(track *Track) error {
	if p.voiceConn == nil {
		return fmt.Errorf("нет голосового подключения")
	}

	// Используем yt-dlp для получения прямой ссылки на аудио
	var audioURL string
	var err error

	if track.Platform == "youtube" || track.Platform == "soundcloud" {
		audioURL, err = p.getDirectURL(track.URL)
		if err != nil {
			return err
		}
	} else {
		audioURL = track.URL
	}

	// Запускаем ffmpeg для потоковой передачи
	ffmpeg := exec.Command("ffmpeg", "-i", audioURL, "-f", "s16le", "-ar", "48000", "-ac", "2", "pipe:1")
	stdout, err := ffmpeg.StdoutPipe()
	if err != nil {
		return err
	}

	if err := ffmpeg.Start(); err != nil {
		return err
	}

	p.cmd = ffmpeg

	// Отправляем аудио в Discord
	p.voiceConn.Speaking(true)
	defer p.voiceConn.Speaking(false)

	buffer := make([]byte, 3840)
	for {
		if p.isPaused {
			// Простая реализация паузы - ждем
			continue
		}

		n, err := stdout.Read(buffer)
		if err == io.EOF || err == io.ErrClosedPipe {
			break
		}
		if err != nil {
			return err
		}

		p.voiceConn.OpusSend <- buffer[:n]
	}

	ffmpeg.Wait()
	return nil
}

func (p *Player) getDirectURL(url string) (string, error) {
	cmd := exec.Command("yt-dlp", "-f", "bestaudio", "-g", url)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// Skip пропускает текущий трек
func (p *Player) Skip() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.isPlaying {
		return fmt.Errorf("ничего не воспроизводится")
	}

	if p.cmd != nil {
		p.cmd.Process.Kill()
	}

	return nil
}

// Pause ставит на паузу
func (p *Player) Pause() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.isPlaying {
		return fmt.Errorf("ничего не воспроизводится")
	}

	p.isPaused = !p.isPaused
	return nil
}

// Stop останавливает воспроизведение
func (p *Player) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.cmd != nil {
		p.cmd.Process.Kill()
	}

	p.stopChan <- true
	p.queue = make([]*Track, 0)
	p.isPlaying = false
	p.isPaused = false
	p.currentTrack = nil
}

// GetCurrentTrack возвращает текущий трек
func (p *Player) GetCurrentTrack() *Track {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.currentTrack
}

// IsPlaying проверяет, идет ли воспроизведение
func (p *Player) IsPlaying() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.isPlaying
}

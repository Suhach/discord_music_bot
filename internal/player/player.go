package player

import (
	"fmt"
	"log"
	"os/exec"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
)

type Track struct {
	Title string
	URL   string
}

type Player struct {
	mu         sync.Mutex
	Queue      []Track
	Voice      *discordgo.VoiceConnection
	isPlaying  bool
	currentIdx int
}

func New() *Player {
	return &Player{}
}

func (p *Player) PlayYouTube(s *discordgo.Session, i *discordgo.InteractionCreate, query string) {
	userVoice, err := findUserVoiceChannel(s, i.GuildID, i.Member.User.ID)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ You must be in a voice channel!",
			},
		})
		return
	}

	if p.Voice == nil {
		p.Voice, err = s.ChannelVoiceJoin(i.GuildID, userVoice, false, true)
		if err != nil {
			log.Println("Join error:", err)
			return
		}
	}

	out, err := exec.Command("yt-dlp", "-f", "bestaudio", "--get-title", "--get-url", query).Output()
	if err != nil {
		log.Println("yt-dlp error:", err)
		return
	}

	lines := splitLines(string(out))
	if len(lines) < 2 {
		return
	}

	track := Track{
		Title: lines[0],
		URL:   lines[1],
	}

	p.mu.Lock()
	p.Queue = append(p.Queue, track)
	p.mu.Unlock()

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("🎶 Added to queue: **%s**", track.Title),
		},
	})

	if !p.isPlaying {
		go p.playLoop()
	}
}

func (p *Player) playLoop() {
	p.isPlaying = true
	defer func() { p.isPlaying = false }()

	for len(p.Queue) > 0 {
		track := p.Queue[0]
		p.Queue = p.Queue[1:]

		log.Printf("▶️ Playing: %s", track.Title)

		// Настройка ffmpeg для вывода в stdout (raw PCM)
		cmd := exec.Command("ffmpeg", "-i", track.URL, "-f", "s16le", "-ar", "48000", "-ac", "2", "pipe:1")
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Println("ffmpeg pipe:", err)
			continue
		}

		err = cmd.Start()
		if err != nil {
			log.Println("ffmpeg start:", err)
			continue
		}

		// Создаём поток через dca
		done := make(chan error)
		dca.NewStream(stdout, p.Voice, done) // !!!

		// Ждём окончания проигрывания или ошибки
		if err := <-done; err != nil && err != dca.ErrVoiceNotPlaying {
			log.Println("Stream error:", err)
		}

		cmd.Wait()
	}

	p.Stop()
}

func (p *Player) Skip() {
	if p.Voice != nil {
		p.Voice.Speaking(false)
	}
}

func (p *Player) Stop() {
	if p.Voice != nil {
		p.Voice.Disconnect()
		p.Voice = nil
	}
	p.Queue = nil
}

func findUserVoiceChannel(s *discordgo.Session, guildID, userID string) (string, error) {
	g, err := s.State.Guild(guildID)
	if err != nil {
		return "", err
	}

	for _, vs := range g.VoiceStates {
		if vs.UserID == userID {
			return vs.ChannelID, nil
		}
	}
	return "", fmt.Errorf("user not in voice channel")
}

func splitLines(s string) []string {
	var res []string
	start := 0
	for i, ch := range s {
		if ch == '\n' || ch == '\r' {
			if start < i {
				res = append(res, s[start:i])
			}
			start = i + 1
		}
	}
	if start < len(s) {
		res = append(res, s[start:])
	}
	return res
}

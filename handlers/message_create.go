package handlers

import (
	"kyBot/status"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	if strings.HasPrefix(m.Content, "Wordle") {
		status.AddWordleStats(s, m)
	}

	if !strings.HasPrefix(m.Content, "k!") {
		return
	}

	trim := strings.TrimPrefix(m.Content, "k!")
	split_content := strings.SplitN(trim, " ", 2)
	if len(split_content) < 1 {
		return
	}
	command := strings.ToLower(split_content[0])

	switch command {
	case "test":
		status.SendWordleReminders(s)
	}
}
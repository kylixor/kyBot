package commands

import (
	"kyBot/config"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

var commands = make(map[string]*discordgo.ApplicationCommand)
var currentCommands = make(map[string]*discordgo.ApplicationCommand)

func AddCommand(cmd *discordgo.ApplicationCommand) {
	commands[cmd.Name] = cmd
}

func RegisterCommands(appid string, s *discordgo.Session) {
	// Get current commands to check if new ones need to be added
	currentCommandArray, err := s.ApplicationCommands(appid, "")
	if err != nil {
		log.Errorln("Error fetching current apppliation commands for guild", err)
	}

	if len(currentCommandArray) != 0 {
		log.Debug("Existing registered commands:")
		for _, command := range currentCommandArray {
			currentCommands[command.Name] = command
			log.Debugf("ID: %s | Name: %s", command.ID, command.Name)
		}
	}

	// Blank guildID means register commands globally across Discord
	guildID := ""
	if config.DEBUG {
		guildID = config.DEBUG_GUILD_ID
	}
	for _, command := range commands {
		if _, exists := currentCommands[command.Name]; !exists {
			command, err := s.ApplicationCommandCreate(appid, guildID, command)
			if err != nil {
				log.Errorln("Error creating application command", err, command)
			} else {
				log.Debugf("Registered command in [%s]: %s", guildID, command.Name)
			}
		}
	}
}
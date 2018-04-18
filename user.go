package main

import (
	"encoding/json"
	// "fmt"
	"os"

	"github.com/bwmarrin/discordgo"
)

type UserStateArray struct {
	GID   string      `json:"gID"`
	Users []UserState `json:"users"`
}

type UserState struct {
	Name         string   `json:"name"`
	UserID       string   `json:"userID"`
	CurrentCID   string   `json:"currentCID"`
	LastSeenCID  string   `json:"lastSeenCID"`
	PlayAnthem   bool     `json:"playAnthem"`
	Anthem       string   `json:"anthem"`
	NoiseCredits int      `json:"noiseCredits"`
	Dailies      bool     `json:"dailies"`
	Reminders    []string `json:"reminders"`
}

var USArray UserStateArray

func InitUserFile() {
	//Indent so its readable
	userData, err := json.MarshalIndent(USArray, "", "    ")
	if err != nil {
		panic(err)
	}
	//Open file
	jsonFile, err := os.Create("users.json")
	if err != nil {
		panic(err)
	}
	//Write to file
	_, err = jsonFile.Write(userData)
	if err != nil {
		panic(err)
	}
	//Cleanup
	jsonFile.Close()
}

func (u *UserStateArray) ReadUserFile() {
	file, err := os.Open("users.json")
	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&u)
	if err != nil {
		panic(err)
	}

	file.Close()
	u.WriteUserFile()
}

func (u *UserStateArray) WriteUserFile() {
	//Marshal global variable data
	jsonData, err := json.MarshalIndent(u, "", "    ")
	if err != nil {
		panic(err)
	}
	//Open file
	jsonFile, err := os.Create("users.json")
	if err != nil {
		panic(err)
	}
	//Write to file
	_, err = jsonFile.Write(jsonData)
	if err != nil {
		panic(err)
	}
	//Cleanup
	jsonFile.Close()
}

func (u *UserStateArray) ReadUser(s *discordgo.Session, i interface{}, code string) (UVS UserState, j int) {
	//Search through user array for specific user and return them
	switch code {

	case "VOICE":
		v := i.(*discordgo.VoiceStateUpdate)
		//Return user if they are inside array
		for j := range u.Users {
			if u.Users[j].UserID == v.UserID {
				return u.Users[j], j
			}
		}
		//Or create a new one if they cannot be found
		s.ChannelMessageSend(config.LogID, "```\nCannot find user...Creating new...\n```")
		user := u.CreateUser(s, v, "VOICE")
		return user, len(u.Users) - 1
	case "MSG":
		m := i.(*discordgo.MessageCreate)
		//Return user if they are inside array
		for j := range u.Users {
			if u.Users[j].UserID == m.Author.ID {
				return u.Users[j], j
			}
		}
		//Or create a new one if they cannot be found
		s.ChannelMessageSend(config.LogID, "```\nCannot find user...Creating new...\n```")
		user := u.CreateUser(s, m, "MSG")
		return user, len(u.Users) - 1
	default:
		panic("Incorrect code in ReadUser")
	}

}

func (u *UserStateArray) CreateUser(s *discordgo.Session, i interface{}, code string) (UVS UserState) {
	var user UserState

	switch code {

	case "VOICE":
		v := i.(*discordgo.VoiceStateUpdate)
		//Create user
		usr, _ := s.User(v.UserID)
		member, err := s.GuildMember(v.GuildID, v.UserID)
		user.Name = FormatAuthor(usr, member, err)
		user.UserID = v.UserID
		user.CurrentCID = v.ChannelID
		user.LastSeenCID = v.ChannelID
		user.NoiseCredits = 0
		user.PlayAnthem = true
		break
	case "MSG":
		m := i.(*discordgo.MessageCreate)
		usr, _ := s.User(m.Author.ID)
		channel, _ := s.Channel(m.ChannelID)
		member, err := s.GuildMember(channel.GuildID, m.Author.ID)
		user.Name = FormatAuthor(usr, member, err)
		user.UserID = m.Author.ID
		user.NoiseCredits = 0
		user.PlayAnthem = true
		break

	default:

	}

	u.Users = append(u.Users, user)
	u.WriteUserFile()
	return user
}

func (u *UserStateArray) UpdateUser(s *discordgo.Session, i interface{}, code string) bool {
	switch code {

	case "VOICE":
		v := i.(*discordgo.VoiceStateUpdate)
		//Get user object
		user, j := u.ReadUser(s, v, "VOICE")
		//Update user object
		//If the update is only a change in voice (mute, deafen, etc)
		if user.CurrentCID == v.ChannelID && user.LastSeenCID == v.ChannelID {
			return false
		}
		u.Users[j].CurrentCID = v.ChannelID
		if v.ChannelID != "" {
			u.Users[j].LastSeenCID = v.ChannelID
		}
		break
	case "MSG":
		//m := i.(*discordgo.MessageCreate)
		//user := ReadUser(s, m)
		break
	default:
		panic("Incorrect code sent to WriteUser")
	}

	u.WriteUserFile()
	return true
}

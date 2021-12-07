package kyDB

import (
	"os"
	"path"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB *gorm.DB
)

type User struct {
	gorm.Model
	ID            string `gorm:"primaryKey"` // User ID
	Name          string // Username
	Discriminator string // Unique identifier (#4712)
	CurrentVCID   string // Current voice channel ID
	PreviousVCID  string // Previous voice channel ID
	Credits       int    // Credits in account
	GotDailies    bool   // True if dailies have been claimed today
}

type Guild struct {
	gorm.Model
	ID           string `gorm:"primaryKey"` // discord guild ID
	Name         string // Name of guild
	DefaultCID   string // Default text channel
	MemberRoleID string // Role approved members will be assigned to
	Currency     string // Name of currency that bot uses (i.e. <gold> coins)
	Karma        int    // Bot's karma - tracked per guild
}

type Hangman struct {
	gorm.Model
	GuildID   string `gorm:"primaryKey"` // Guild game is attached to
	ChannelID string // ChannelID where game is played
	MessageID string // MessageID of current hangman game
	Word      string // Word/phrase for the game
	GameState int    // State of game, 1-7 until you lose
	WordState string // State of game's word
	Guessed   string // Characters/words that have been guessed
}

// Creates database
func CreateDBFile(path string) {

	file, err := os.Create(path)
	if err != nil {
		log.Fatalln("Error creating database", err.Error())
	}
	file.Close()
}

// Connects to database and returns the connection
func Connect() *gorm.DB {
	cwd, _ := os.Getwd()
	db_path := path.Join(cwd, "data", "db.sqlite")

	if _, err := os.Stat(db_path); os.IsNotExist(err) {
		CreateDBFile(db_path)
	}

	var err error
	DB, err = gorm.Open(sqlite.Open(db_path), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		log.Panicln("Failed to connect to database", err.Error())
	}

	return DB
}

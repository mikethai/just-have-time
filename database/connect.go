package database

import (
	"fmt"
	"log"
	"strconv"

	"github.com/mikethai/just-have-time/config"
	"github.com/mikethai/just-have-time/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Declare the variable for the database
var DB *gorm.DB

// ConnectDB connect to db
func ConnectDB() {
	var err error
	p := config.Config("POSTGRES_DB_PORT")
	port, err := strconv.ParseUint(p, 10, 32)

	if err != nil {
		log.Println("Port Error: " + err.Error())
	}

	// Connection URL to connect to Postgres Database
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.Config("POSTGRES_DB_HOST"), port, config.Config("POSTGRES_DB_USER"), config.Config("POSTGRES_DB_PASSWORD"), config.Config("POSTGRES_DB_NAME"))
	// Connect to the DB and initialize the DB variable
	DB, err = gorm.Open(postgres.Open(dsn))

	if err != nil {
		panic("failed to connect database")
	}

	fmt.Println("Connection Opened to Database")

	// Migrate the database
	DB.AutoMigrate(&model.User{}, &model.Song{}, &model.Hashtag{}, &model.Follow{}, &model.StorySong{})

	// DB.Migrator().DropTable(&model.User{}, &model.Song{}, &model.Hashtag{}, &model.Follow{}, &model.StorySong{})

	fmt.Println("Database Migrated")
}

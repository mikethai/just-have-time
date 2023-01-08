package database

import (
	"fmt"
	"log"

	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
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
	// p := config.Config("POSTGRES_DB_PORT")
	// port, err := strconv.ParseUint(p, 10, 32)

	if err != nil {
		log.Println("Port Error: " + err.Error())
	}

	var (
		dbHost = config.Config("POSTGRES_DB_HOST")
		dbUser = config.Config("POSTGRES_DB_USER")
		dbPwd  = config.Config("POSTGRES_DB_PASSWORD")
		dbName = config.Config("POSTGRES_DB_NAME")
	)

	// Connection URL to connect to Postgres Database
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, dbUser, dbName, dbPwd)

	// Connect to the DB and initialize the DB variable
	DB, err = gorm.Open(postgres.New(postgres.Config{
		DriverName: "cloudsqlpostgres",
		DSN:        dsn,
	}))

	if err != nil {
		panic("failed to connect database")
	}

	fmt.Println("Connection Opened to Database")

	// Migrate the database
	DB.AutoMigrate(&model.User{}, &model.Follow{}, &model.Song{})
	DB.AutoMigrate(&model.StorySong{}, &model.Hashtag{})

	// DB.Migrator().DropTable(&model.User{}, &model.Song{}, &model.Hashtag{}, &model.Follow{}, &model.StorySong{})

	fmt.Println("Database Migrated")
}

package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Configuration struct {
	DatabaseName        string
	DatabaseHost        string
	DatabaseUser        string
	DatabasePassword    string
	MigrateToVersion    string
	MigrationLocation   string
	FileStorageLocation string
	RabbitURL           string
	LogLevel            string
	Port                string
}

func GetConfiguration() Configuration {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println(err)
	}
	migrationLocation, set := os.LookupEnv("MIGRATION_LOCATION")
	if !set {
		migrationLocation = "migrationsx"
	}
	migrateToVersion, set := os.LookupEnv("MIGRATE")
	if !set {
		migrateToVersion = "latest"
	}
	staticFilesLocation, set := os.LookupEnv("FILES_LOCATION")
	if !set {
		staticFilesLocation = "file_storage"
	}

	return Configuration{
		DatabaseName:        os.Getenv("DB_NAME"),
		DatabaseHost:        os.Getenv("DB_HOST"),
		DatabaseUser:        os.Getenv("DB_USER"),
		DatabasePassword:    os.Getenv("DB_PASSWORD"),
		MigrateToVersion:    migrateToVersion,
		MigrationLocation:   migrationLocation,
		FileStorageLocation: staticFilesLocation,
		RabbitURL:           os.Getenv("RABBIT_URL"),
		Port:                os.Getenv("PORT_SERVER"),
	}
}

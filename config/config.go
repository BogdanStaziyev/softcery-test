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
	Port                string
}

// GetConfiguration returns configuration values from environment variables
func GetConfiguration() Configuration {
	//Load .env variables if exists
	err := godotenv.Load(".env")
	if err != nil {
		log.Println(err)
	}

	//Returns configuration of migration location path from environment variables or use default variable
	migrationLocation, set := os.LookupEnv("MIGRATION_LOCATION")
	if !set {
		migrationLocation = "migrations"
	}

	//Returns configuration of migration version values from environment variables or use default variable
	migrateToVersion, set := os.LookupEnv("MIGRATE")
	if !set {
		migrateToVersion = "latest"
	}

	//Returns configuration of image store path from environment variables or use default variable
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

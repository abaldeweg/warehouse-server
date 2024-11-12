package database

import (
	"fmt"
	"log"

	"github.com/abaldeweg/warehouse-server/gateway/core/models"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	viper.SetDefault("MYSQL_URL", "adm:pass@tcp(localhost:3306)/warehouse?charset=utf8mb4&parseTime=True&loc=Local")
	viper.SetDefault("SQLITE_NAME", "warehouse")
	viper.SetDefault("DATABASE", "sqlite")

	databaseType := viper.GetString("DATABASE")

	var driver gorm.Dialector

	switch databaseType {
	case "mysql":
		driver = mysql.Open(viper.GetString("MYSQL_URL"))
	case "sqlite":
		fallthrough
	default:
		driver = sqlite.Open(viper.GetString("SQLITE_NAME") + ".db")
	}

	db, err := gorm.Open(driver, &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("Connected to database!")

	if databaseType != "mysql" {
		runMigrations(db)
	}

	return db
}

func runMigrations(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.Author{},
		&models.Branch{},
		&models.Condition{},
		&models.Tag{},
	)

	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	fmt.Println("Migrations run successfully!")
}

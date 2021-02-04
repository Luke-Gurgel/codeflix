package db

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/Luke-Gurgel/codeflix/domain/model"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	_ "gorm.io/driver/sqlite"
)

func init() {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	err := godotenv.Load(basepath + "/../../.env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func ConnectDB(env string) *gorm.DB {
	dsn := os.Getenv("dsn")
	dbType := os.Getenv("dbType")

	if env == "test" {
		dsn = os.Getenv("dsnTest")
	}

	db, err := gorm.Open(dbType, dsn)

	if err != nil {
		log.Fatalf("Error connecting to Database: %v", err)
		panic(err)
	}

	if os.Getenv("debug") == "true" {
		db.LogMode(true)
	}

	if os.Getenv("AutoMigrateDB") == "true" {
		db.AutoMigrate(&model.Bank{}, &model.Account{}, &model.PixKey{}, &model.Transaction{})
	}

	return db
}

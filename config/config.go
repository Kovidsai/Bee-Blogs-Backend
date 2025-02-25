package config

import (
	"MINIPROJECT/models"
	"fmt"
	"log"
	"os"

	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func ConnectDB() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("error Loading env file")
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	// MySQL connection string
	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbPort, dbName)
	orm.RegisterDriver("mysql", orm.DRMySQL)
	err = orm.RegisterDataBase("default", "mysql", dataSource)
	if err != nil {
		log.Fatal("DataBase Connection Failed", err)
	}

	orm.RegisterModel(new(models.User), new(models.Post), new(models.Comment))
	orm.RunSyncdb("default", false, true)
	fmt.Printf("DataBase Connection Succeded")
}

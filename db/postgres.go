package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"password-lock/config"
)

func ConnectToDatabase(cfg *config.Config) *gorm.DB {
	db, err := gorm.Open(postgres.Open(cfg.PgUrl), nil)
	if err != nil {
		log.Println("Could not connect to database: ", err.Error())
	} else {
		log.Println(" Successfully connected to database")
	}

	return db
}

const (
	USERS_TABLE                   = "users"
	ENTITIES_TABLE                = "entities"
	TOKENS_TABLE                  = "tokens"
	PERSONAL_QUESTIONS_TABLE      = "personal_questions"
	USER_PERSONAL_QUESTIONS_TABLE = "user_personal_questions"
)

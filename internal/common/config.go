package common

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var Conf *QuizConfig

type Config interface {
	Parse()
}

type QuizConfig struct {
	PostgresDsn string
	Addr        string
}

func NewQuizConfig() *QuizConfig {
	config := new(QuizConfig)
	config.parse()
	return config
}

func (config *QuizConfig) parse() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла:", err)
	}
	config.PostgresDsn = parseVar("GOOSE_DBSTRING")
	config.Addr = parseVar("ADDR")
}

func parseVar(varName string) string {
	variable := os.Getenv(varName)
	if variable == "" {
		panic(varName + " not provided")
	}
	return variable
}

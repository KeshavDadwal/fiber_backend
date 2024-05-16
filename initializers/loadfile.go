package initializers

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadVariable() {
	err := godotenv.Load()

	if err!=nil{
		log.Fatal("err to load the data there")
	}

}
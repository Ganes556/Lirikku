package main

import (
	"github.com/Lirikku/configs"
	"github.com/joho/godotenv"
)

func init(){
	godotenv.Load(".env.example")
	// godotenv.Load(".env.production")
	configs.InitDB()
}

func main() {	
}
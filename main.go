package main

import (
	"github.com/Lirikku/configs"
	"github.com/joho/godotenv"
)

func init(){
	godotenv.Load(".env.dev")
	configs.InitDB()
}

func main() {
}
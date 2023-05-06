package main

import (
	"github.com/Lirikku/configs"
	"github.com/Lirikku/routes"
	"github.com/joho/godotenv"
)

func init(){
	godotenv.Load(".env.production")
	configs.InitDB()
}

func main() {	
	e := routes.NewRoute()
	e.Logger.Fatal(e.Start(":8000"))
}
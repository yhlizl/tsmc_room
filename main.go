package main

import (
	_ "go-gofram-chat/boot"
	_ "go-gofram-chat/router"

	"go-gofram-chat/app/models"

	"github.com/gogf/gf/frame/g"
)

func init() {
	// viper.SetConfigType("json")
	// if err := viper.ReadConfig(bytes.NewBuffer(conf.AppJsonConfig)); err != nil {
	// 	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
	// 		// Config file not found; ignore error if desired
	// 		log.Println("no such config file")
	// 	} else {
	// 		// Config file was found but another error was produced
	// 		log.Println("read config error")
	// 	}
	// 	log.Fatal(err) //read config fatal fail

	// }
	models.InitDB()
}
func main() {
	g.Server().Run()
}

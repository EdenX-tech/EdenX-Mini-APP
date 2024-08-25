package main

import (
	"ginDemo/config"
	"ginDemo/database"
	"ginDemo/route"
	"github.com/spf13/viper"
	"log"
)

func main() {

	config.Setup()

	database.Init()

	r := route.InitRouter()

	// 设置受信任的代理
	err := r.SetTrustedProxies([]string{viper.GetString("service_ip.ip")})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Starting server on :3000")
	if err := r.Run(":3000"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

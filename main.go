package main

import (
	"ginDemo/config"
	"ginDemo/database"
	"ginDemo/route"
	"ginDemo/services"
	"github.com/robfig/cron/v3"
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

	c := cron.New()
	_, err = c.AddFunc("* * * * *", func() {
		services.RestoreEnergy()
	})
	if err != nil {
		log.Printf("Failed to add cron job: %v", err)
		return
	}
	c.Start()

	defer c.Stop()

	if err := r.Run(":3000"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

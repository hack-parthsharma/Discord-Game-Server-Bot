package main

import (
	"os"
	"os/signal"
	"time"

	"github.com/Asthetic/DiscordGameServerBot/config"
	"github.com/Asthetic/DiscordGameServerBot/discord"
	"github.com/Asthetic/DiscordGameServerBot/network"
	"github.com/Asthetic/DiscordGameServerBot/storage"
	log "github.com/sirupsen/logrus"
)

func main() {
	config, err := config.NewConfig()
	if err != nil {
		log.Error(err, "unable to read configuration file")
		return
	}

	discordBot, err := discord.New(config.DiscordCfg)
	if err != nil {
		log.WithError(err).Error("Failed to initialize discord session")
		return
	}

	defer discordBot.Close()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	done := make(chan bool)

	for {
		getIP(discordBot)

		select {
		case <-done:
			return
		case s := <-sig:
			log.Infof("Got signal to kill program: %v", s)
			return
		case <-ticker.C:
			getIP(discordBot)
		}
	}
}

func getIP(discord *discord.Discord) {
	currentIP, err := storage.GetIP()
	if err != nil {
		log.WithError(err).Errorf("Error fetching file: %v", err)
	}

	ip, err := network.GetPublicIP()
	if err != nil {
		log.WithError(err).Errorf("Error getting IP address: %v")
		return
	}

	log.Infof("Sucessfully got public IP address: %v", ip)

	if currentIP != ip {
		currentIP = ip
		err = storage.WriteIP(network.Network{IP: currentIP})
		if err != nil {
			log.WithError(err).Errorf("Error writing IP to local storage")
		}

		discord.SendUpdatedIP(ip)
	}
}

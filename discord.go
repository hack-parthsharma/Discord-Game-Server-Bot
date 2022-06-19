package discord

import (
	"fmt"

	"github.com/Asthetic/DiscordGameServerBot/config"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

// Discord contains the discord session and connection settings
type Discord struct {
	session  *discordgo.Session
	token    string
	channels []string
}

// New creates a new Discord session with config
func New(cfg config.Discord) (*Discord, error) {
	dg, err := discordgo.New(fmt.Sprintf("Bot %s", cfg.Token))
	if err != nil {
		return nil, err
	}

	if err = dg.Open(); err != nil {
		return nil, err
	}

	return &Discord{
		session:  dg,
		token:    cfg.Token,
		channels: cfg.Channels,
	}, nil
}

// Close cleanly closes the current Discord session
func (d *Discord) Close() {
	d.session.Close()
}

// SendUpdatedIP sends the updated IP address to the configured channels
func (d *Discord) SendUpdatedIP(ip string) {
	for _, channel := range d.channels {
		msg := formatMessage(ip)
		result, err := d.session.ChannelMessageSendComplex(channel, msg)
		if err != nil {
			log.WithError(err).Errorf("Unable to post message to channel: %v", channel)
		} else {
			log.Info("Sucessfully sent Discord message")
		}

		if err = d.pinMessage(result); err != nil {
			log.WithError(err).Error("failed to pin new message")
		} else {
			log.Info("Sucessfully pinned Discord message")
		}
	}
}

func formatMessage(ip string) *discordgo.MessageSend {
	return &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Type:  discordgo.EmbedTypeRich,
			Title: "Minecraft Server IP Updated",
			Color: 5439264,
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL:    "https://i.imgur.com/6rp13.png",
				Width:  128,
				Height: 128,
			},
			Fields: formatFields(ip),
		},
	}
}

func formatFields(ip string) []*discordgo.MessageEmbedField {
	fields := []*discordgo.MessageEmbedField{}
	field := &discordgo.MessageEmbedField{
		Name:   "Minecraft",
		Value:  ip,
		Inline: true,
	}

	fields = append(fields, field)
	return fields
}

func (d *Discord) pinMessage(msg *discordgo.Message) error {
	pinnedMsgs, err := d.session.ChannelMessagesPinned(msg.ChannelID)
	if err != nil {
		return err
	}

	for _, pinned := range pinnedMsgs {
		if pinned.Author.ID == d.session.State.User.ID {
			if err = d.session.ChannelMessageUnpin(pinned.ChannelID, pinned.ID); err != nil {
				log.WithError(err).Errorf("failed to unpin message before pinning new one for user id: %v, author: %v, message id: %v", pinned.Author.ID, pinned.Author.Username, pinned.ID)
			}

		}
	}

	return d.session.ChannelMessagePin(msg.ChannelID, msg.ID)
}

package main

import (
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
)

const allowOption = "-Xallow-kotlin"

func main() {
	token := flag.String("token", "", "Bot token")
	flag.Parse()

	if *token == "" {
		panic("Token must be provided.")
	}

	discord, err := discordgo.New("Bot " + *token)

	discord.Identify.Intents = discordgo.IntentsGuildMessages

	if err != nil {
		fmt.Println(err)
	}

	discord.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		content := m.Content

		if matches, err := regexp.Match("([КкKk]+)(.{0,5})([ОоOo0]+)(.{0,5})([ТтTt]+)(.{0,5})([Ll]+)(.{0,5})([Ii]+)(.{0,5})([Nn]+)(.{0,5})", []byte(content)); err == nil && matches {
			if !m.Author.Bot {
				if !strings.Contains(content, allowOption) {
					err := s.ChannelMessageDelete(m.ChannelID, m.Message.ID)
					if err != nil {
						fmt.Println(err)
					}

					_, err2 := s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
						Content: m.Author.Mention(),
						Embed: &discordgo.MessageEmbed{
							Color:       0xb061f1, // purple
							Description: "Kotlin is prohibited by the owner. Include `-Xallow-kotlin` in your message to avoid the rule.",
						},
					})
					if err2 != nil {
						fmt.Println(err2)
					}
				}
			}
		} else {
			if err != nil {
				fmt.Println(err)
			}
		}
	})

	if err := discord.Open(); err != nil {
		fmt.Println("Epic fail:", err)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	if err := discord.Close(); err != nil {
		panic(err)
	}
}

package embed

import "github.com/bwmarrin/discordgo"

func Wrap(title string, text string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       title,
		Description: text,
	}
}

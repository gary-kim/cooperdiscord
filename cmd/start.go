//    Copyright (C) 2020 Gary Kim <gary@garykim.dev>, All Rights Reserved
//
//    This program is free software: you can redistribute it and/or modify
//    it under the terms of the GNU Affero General Public License as published
//    by the Free Software Foundation, either version 3 of the License, or
//    (at your option) any later version.
//
//    This program is distributed in the hope that it will be useful,
//    but WITHOUT ANY WARRANTY; without even the implied warranty of
//    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//    GNU Affero General Public License for more details.
//
//    You should have received a copy of the GNU Affero General Public License
//    along with this program.  If not, see <https://www.gnu.org/licenses/>.

package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/mattn/go-shellwords"
	"github.com/spf13/cobra"

	"gomod.garykim.dev/cooperdiscord/cooper"
)

var (
	guilds        string
	commandPrefix string
	guildList     []string
	parser        *shellwords.Parser
	courses       []cooper.CourseInfo
)

var helpMessage = `
Usage:
  $PREFIX course-search "CH-110"
`

func init() {
	start := &cobra.Command{
		Use:   "start",
		Short: "Start the bot",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Initial setup
			for _, g := range strings.Split(guilds, ";") {
				guildList = append(guildList, g)
			}
			parser = shellwords.NewParser()
			helpMessage = strings.Replace(helpMessage, "$PREFIX", commandPrefix, -1)
			coursest, err := cooper.ScrapeInfo()
			courses = coursest
			if err != nil {
				return err
			}

			discord, err := discordgo.New(Token)
			if err != nil {
				return err
			}
			discord.Identify.Intents = discordgo.MakeIntent(
				discordgo.IntentsAllWithoutPrivileged,
			)
			discord.AddHandler(onMessageHandler)
			err = discord.Open()
			if err != nil {
				return err
			}
			log.Printf("Sucessfully started")
			c := make(chan bool)
			_ = <-c
			return nil
		},
	}

	start.PersistentFlags().StringVarP(&guilds, "guilds", "g", "", "Discord token (ENV: DISCORD_GUILS)")
	if os.Getenv("DISCORD_GUILDS") != "" {
		_ = start.PersistentFlags().Set("guilds", os.Getenv("DISCORD_GUILDS"))
	}

	start.PersistentFlags().StringVarP(&commandPrefix, "prefix", "p", "/cooper", "Command prefix")

	Root.AddCommand(start)
}

func onMessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m == nil || !isEnabledGuild(m.GuildID) || m.Message == nil || !strings.HasPrefix(m.Message.Content, commandPrefix) {
		return
	}

	log.Printf("Received new message: %s from %s\n", m.Message.Content, m.Author.String())

	command, err := parser.Parse(m.Message.Content)
	if err != nil {
		log.Printf("Could not parse command: %s\n", err)
		return
	}
	if len(command) < 2 || command[1] == "--help" || command[1] == "help" || command[1] == "-h" || command[1] == "?" {
		err = printHelpMessage(s, m)
		if err != nil {
			log.Printf("Could not send help message: %s\n", err)
		}
		return
	}

	if len(command) > 2 && command[1] == "course-search" {
		course := findCourseByID(command[2])
		if course == nil {
			_, err = s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
				Content:   "Could not find course with code " + command[2],
				Reference: m.Reference(),
			})
			if err != nil {
				log.Printf("Could not send not found response: %s\n", err)
			}
			return
		}
		message := courseToMessage(course)
		message.Reference = m.Reference()
		_, err = s.ChannelMessageSendComplex(m.ChannelID, message)
		if err != nil {
			log.Printf("Could not send response: %s\n", err)
		}
		return
	}
	if strings.ToLower(command[1]) == "care" {
		if err := printCooperCareInfo(s, m); err != nil {
			log.Printf("Could not send response: %s\n", err)
		}
		return
	}
	err = printHelpMessage(s, m)
	if err != nil {
		log.Printf("Could not send help message: %s\n", err)
	}
	return
}

func printCooperCareInfo(s *discordgo.Session, m *discordgo.MessageCreate) error {
	_, err := s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Reference: m.Reference(),
		Embed: &discordgo.MessageEmbed{
			URL:         "https://cooper.care",
			Title:       "Cooper Care",
			Description: "Do not worry! Cooper Care is available for Cooper Union students. Visit this link for the best care Cooper can give...",
		},
	})
	return err
}

func printHelpMessage(s *discordgo.Session, m *discordgo.MessageCreate) error {
	_, err := s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Reference: m.Reference(),
		Embed: &discordgo.MessageEmbed{
			URL:         "https://github.com/gary-kim/cooperdiscord",
			Title:       "Cooper Union Discord Bot",
			Description: helpMessage,
		},
	})
	return err
}

func courseToMessage(info *cooper.CourseInfo) *discordgo.MessageSend {
	dashCode := strings.Replace(info.Code, " ", "-", -1)
	message := "**Code:** " + dashCode
	if len(info.Codes) > 1 {
		message += " (AKA: " + strings.Join(info.Codes, ", ") + ")"
	}
	message += "\n" +
		"**Name:** " + info.Name + "\n" +
		"**Description:** " + info.Description + "\n" +
		"**Extra Info:** " + info.ExtraInfo

	messageEmbed := &discordgo.MessageEmbed{
		URL:         "https://dtss.cooper.edu/Student/Student/Courses/Search?keyword=" + dashCode,
		Title:       "DTSS: " + info.Code,
		Description: "DTSS page for " + info.Name,
	}

	return &discordgo.MessageSend{
		Content: message,
		Embed:   messageEmbed,
	}
}

func findCourseByID(id string) *cooper.CourseInfo {
	id = strings.TrimSpace(id)
	id = strings.Replace(id, "-", " ", -1)
	id = strings.ToUpper(id)
	for _, course := range courses {
		if course.Code == id {
			return &course
		}
	}
	return nil
}

func isEnabledGuild(guildID string) bool {
	if guilds == "" {
		return true
	}
	for _, g := range guildList {
		if g == guildID {
			return true
		}
	}
	return false
}

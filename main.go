package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	token := os.Getenv("DISCORD_TOKEN")
	guildID := os.Getenv("DISCORD_GUILD")
	voiceCategory := os.Getenv("VOICE_CATEGORY")
	voiceChannel := os.Getenv("VOICE_CHANNEL")
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Error creating discord session: %s", err)
	}
	ordre_img, err := os.ReadFile("resources/ordre.png")
	if err != nil {
		log.Fatalf("Error reading ordre.png: %s", err)
	}
	ordre_img_bytes := bytes.NewReader(ordre_img)
	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "v5",
			Description: "Donne le lien de la v5 du forum",
		},
		{
			Name:        "v4",
			Description: "Donne le lien de la v4 du forum",
		},
		{
			Name:        "v3",
			Description: "Donne le lien de la v3 du forum",
		},
		{
			Name:        "v2",
			Description: "Donne le lien de la v2 du forum",
		},
		{
			Name:        "ordre",
			Description: "Donne un ordre d'apprentissage des tricks pour les débutants",
		},
	}
	commandHandlers := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"v5": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Voici le lien de la v5 du forum : https://forum.penspinning-france.fr/",
				},
			})
		},
		"v4": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Voici le lien de la v4 du forum : https://thefpsb.1fr1.net/",
				},
			})
		},
		"v3": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Voici le lien de la v3 du forum : https://thefpsb.penspinning.fr/index.php",
				},
			})
		},
		"v2": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Voici le lien de la v2 du forum : https://thefpsbv2.penspinning.fr/",
				},
			})
		},
		"ordre": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{

					Files: []*discordgo.File{
						{
							Name:        "ordre.png",
							ContentType: "image/png",
							Reader:      ordre_img_bytes,
						},
					},
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       "Ordre d'apprentissage",
							Description: "Voici l'ordre d'apprentissage privilégié pour les débutants !",
							Image: &discordgo.MessageEmbedImage{
								URL: "attachment://ordre.png",
							},
						},
					},
				},
			})
		}}

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	tempVoiceChannels := make(map[string][]*discordgo.User)

	s.AddHandler(func(s *discordgo.Session, v *discordgo.VoiceStateUpdate) {
		before := v.BeforeUpdate
		channelID := v.ChannelID
		user, err := s.User(v.UserID)
		if err != nil {
			log.Fatalf("Error getting user: %s", err)
		}
		// Member joins voice creation channel
		if channelID == voiceChannel && (before == nil || before.ChannelID != voiceChannel) {
			username := user.Username
			log.Printf("%s joined temp chan creation channel", username)
			ch, err := s.GuildChannelCreateComplex(guildID, discordgo.GuildChannelCreateData{
				Name:     fmt.Sprintf("Salon de %s", username),
				Type:     discordgo.ChannelTypeGuildVoice,
				ParentID: voiceCategory,
			})
			if err != nil {
				log.Fatalf("Error creating voice channel: %s", err)
			}
			log.Printf("Created temp channel with ID %s and name %s", ch.ID, fmt.Sprintf("Salon de %s", username))
			tempVoiceChannels[ch.ID] = []*discordgo.User{} // Adds new channel to internal map of temporary channels
			s.GuildMemberMove(guildID, v.UserID, &ch.ID)   // Moves member to new channel
		}
		_, isTemp := tempVoiceChannels[channelID]
		// If user joins temporary channel
		if isTemp {
			tempVoiceChannels[channelID] = append(tempVoiceChannels[channelID], user) // Adds user to temporary channel
		}
		// If user was in a channel
		if before != nil {
			beforeChannel, err := s.Channel(before.ChannelID)
			if err != nil {
				log.Fatalf("Error getting channel: %s", err)
			}
			members, isTemp := tempVoiceChannels[before.ChannelID]
			// If user was in a temporary channel
			if isTemp {
				for _, u := range tempVoiceChannels[before.ChannelID] {
					if u.ID == v.UserID {
						tempVoiceChannels[before.ChannelID] = append(tempVoiceChannels[before.ChannelID][:0], tempVoiceChannels[before.ChannelID][1:]...) // Removes user from temporary channel in internal map
						members = append(members[:0], members[1:]...)                                                                                     // Removes user from members list
						break
					}
				}
				// If user has left channel
				if v.ChannelID != before.ChannelID {
					log.Printf("%s left temp chan %s", user.Username, beforeChannel.Name)
					// If channel is empty
					if len(members) == 0 {
						s.ChannelDelete(before.ChannelID) // Deletes channel
						log.Printf("Deleted temp channel with ID %s and name %s", before.ChannelID, beforeChannel.Name)
						delete(tempVoiceChannels, before.ChannelID) // Removes channel from internal map
					}
				}
			}
		}
	})

	err = s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, guildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	log.Println("Gracefully shutting down.")
}

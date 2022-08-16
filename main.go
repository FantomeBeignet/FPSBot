package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

type Spinner struct {
	Key     string `json:"key"`
	Name    string `json:"name"`
	Twitter string `json:"twitter,omitempty"`
	Youtube string `json:"youtube,omitempty"`
}

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
	spinnerdex_api := os.Getenv("SPINNERDEX_API")
	if err != nil {
		log.Fatalf("Error creating discord session: %s", err)
	}
	ordre_img, err := os.ReadFile("resources/ordre.png")
	if err != nil {
		log.Fatalf("Error reading ordre.png: %s", err)
	}
	ordre_img_bytes := bytes.NewReader(ordre_img)
	sdex_img, err := os.ReadFile("resources/Spinnerdex.png")
	if err != nil {
		log.Fatalf("Error reading FPSBot.png: %s", err)
	}
	sdex_bytes := bytes.NewReader(sdex_img)
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
		{
			Name:        "spinner",
			Description: "Renvoie les liens vers le Twitter et le Youtube du spinner demandé",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "nom",
					Description: "Le nom du spinner",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
				},
			},
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
		},
		"spinner": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			name := i.ApplicationCommandData().Options[0].Value
			query := fmt.Sprintf("%s/spinner/%s", spinnerdex_api, name)
			info, err := http.Get(query)
			if err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Erreur lors de la récupération des informations du spinner",
					},
				})
				return
			}
			defer info.Body.Close()
			var spinner Spinner
			err = json.NewDecoder(info.Body).Decode(&spinner)
			if err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Erreur lors de la récupération des informations du spinner",
					},
				})
				return
			}
			log.Printf("Spinner %s was requested in SpinnerDex", name)
			if spinner.Name == "" {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{
							{
								Description: fmt.Sprintf("Le spinner ***%s*** est introuvable", name),
								Color:       0xFF0000,
							},
						},
					},
				})
				return
			}
			if spinner.Twitter == "" {
				spinner.Twitter = "*Aucun Twitter trouvé*"
			}
			if spinner.Youtube == "" {
				spinner.Youtube = "*Aucun YouTube trouvé*"
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Files: []*discordgo.File{
						{
							Name:        "spinnerdex.png",
							ContentType: "image/png",
							Reader:      sdex_bytes,
						},
					},
					Embeds: []*discordgo.MessageEmbed{
						{
							Title: spinner.Name,
							Fields: []*discordgo.MessageEmbedField{
								{
									Name:  "Twitter",
									Value: spinner.Twitter,
								},
								{
									Name:  "YouTube",
									Value: spinner.Youtube,
								},
							},
							Footer: &discordgo.MessageEmbedFooter{
								Text: "Powered by SpinnerDex",
							},
							Thumbnail: &discordgo.MessageEmbedThumbnail{
								URL: "attachment://spinnerdex.png",
							},
						},
					},
				},
			},
			)
		},
	}

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
				// If user has left channel
				if v.ChannelID != before.ChannelID {
					for i, u := range tempVoiceChannels[before.ChannelID] {
						if u.ID == v.UserID {
							tempVoiceChannels[before.ChannelID] = append(tempVoiceChannels[before.ChannelID][:i], tempVoiceChannels[before.ChannelID][i+1:]...) // Removes user from temporary channel in internal map
							members = append(members[:i], members[i+1:]...)                                                                                     // Removes user from members list
							break
						}
					}
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

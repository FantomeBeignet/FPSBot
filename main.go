package main

import (
	"bytes"
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
	GuildID := os.Getenv("DISCORD_GUILD")
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

	err = s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, GuildID, v)
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

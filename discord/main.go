package discord

// code for a discord bot that has a slash command just to test the discord bot
import (
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "ping",
			Description: "Responds with Pong!",
		},
	}
	GuildID string
)

// StartBot starts the Discord bot.
func StartBot() error {
	// Create a new Discord session using the provided bot token
	dg, err := discordgo.New("Bot " + os.Getenv("DISCORD_BOT_TOKEN"))
	if err != nil {
		return err
	}

	// Add a handler for the "ready" event
	dg.AddHandler(ready)

	// Add a handler for the "interactionCreate" event
	dg.AddHandler(interactionCreate)

	// Open a websocket connection to Discord and begin listening
	err = dg.Open()
	if err != nil {
		return err
	}
	defer dg.Close()

	// Wait for the bot to be ready before attempting to register commands
	log.Println("Bot is connecting to Discord...")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	// Clean up commands on exit - this part is already good but we'll keep it
	if GuildID != "" {
		log.Println("Cleaning up commands...")
		// We need to get the current commands to delete them
		cmds, err := dg.ApplicationCommands(dg.State.User.ID, GuildID)
		if err != nil {
			log.Printf("Error fetching commands: %v", err)
		}

		for _, cmd := range cmds {
			err := dg.ApplicationCommandDelete(dg.State.User.ID, GuildID, cmd.ID)
			if err != nil {
				log.Printf("Error deleting command %v: %v", cmd.Name, err)
			}
		}
	}

	return nil
}

// ready is called when the bot receives the "ready" event
func ready(s *discordgo.Session, event *discordgo.Ready) {
	// Register commands after the bot is ready
	log.Printf("Bot is ready as %s#%s\n", event.User.Username, event.User.Discriminator)

	// Register commands using the bot's actual ID from the session
	log.Println("Registering commands...")
	for _, cmd := range commands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, GuildID, cmd)
		if err != nil {
			log.Printf("Error registering command %v: %v", cmd.Name, err)
			// Don't panic here, just log the error
		} else {
			log.Printf("Registered command: %s", cmd.Name)
		}
	}
}

// interactionCreate is called when the bot receives the "interactionCreate" event
func interactionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand {
		switch i.ApplicationCommandData().Name {
		case "ping":
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Pong!",
				},
			})
		}
	}
}

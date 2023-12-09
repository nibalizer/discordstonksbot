package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	finnhub "github.com/Finnhub-Stock-API/finnhub-go/v2"
	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	Token string
)

func init() {

	Token = os.Getenv("DISCORD_BOT_TOKEN")
}

func Quote(symbol string) string {
	cfg := finnhub.NewConfiguration()
	apiKey := os.Getenv("FINNHUB_API_KEY")
	if len(apiKey) == 0 {
		log.Fatal("FINNHUB_API_KEY not set!")
	}
	cfg.AddDefaultHeader("X-Finnhub-Token", apiKey)
	finnhubClient := finnhub.NewAPIClient(cfg).DefaultApi

	// Quote
	quote, _, err := finnhubClient.Quote(context.Background()).Symbol(symbol).Execute()
	if err != nil {
		log.Print(err)
		return ""
	}
	if *quote.C == 0 {
		result := "Symbol not found"
		log.Print(result)
		return result
	}

	result := fmt.Sprintf("%s: $%.2f(%.2f%%)", symbol, *quote.C, *quote.Dp)
	fmt.Println(result)
	return result

}

func main() {
	apiKey := os.Getenv("FINNHUB_API_KEY")
	if len(apiKey) == 0 {
		log.Fatal("FINNHUB_API_KEY not set!")
	}

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(genMessageCreate())

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	scc := make(chan os.Signal, 1)
	signal.Notify(scc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-scc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func genMessageCreate() func(s *discordgo.Session, m *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {

		// Ignore all messages created by the bot itself
		// This isn't required in this specific example but it's a good practice.
		if m.Author.ID == s.State.User.ID {
			return
		}
		// If the message is "ping" reply with "Pong!"
		if m.Content == "ping" {
			s.ChannelMessageSend(m.ChannelID, "Pong!")
		}

		// If the message is "pong" reply with "Ping!"
		if m.Content == "pong" {
			s.ChannelMessageSend(m.ChannelID, "Ping!")
		}

		if strings.HasPrefix(m.Content, "!q") {
			symbols := strings.Split(strings.ToUpper(strings.Split(m.Content, " ")[1]), ",")
			for _, symbol := range symbols {
				fmt.Printf("symbol = %+v\n", symbol)
				resp := Quote(symbol)
				s.ChannelMessageSend(m.ChannelID, resp)
			}
		}
	}
}

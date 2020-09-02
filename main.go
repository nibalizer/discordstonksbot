package main

import (
	"context"
	"fmt"
	finnhub "github.com/Finnhub-Stock-API/finnhub-go"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	stonksV1 "github.com/nibalizer/stonksapi/v1"
)

// Variables used for command line parameters
var (
	Token string
)

func init() {

	Token = os.Getenv("DISCORD_BOT_TOKEN")
}

func main() {
	// initialize data and auth
	// initialize symbol -> description data
	stonksDataPath := os.Getenv("STONKS_DATA_PATH")
	if stonksDataPath == "" {
		log.Fatal("Please set STONKS_DATA_PATH")
	}

	records, err := stonksV1.GetStonksDataFromCSV(stonksDataPath)
	if err != nil {
		log.Fatal(err)
	}
	// initialize finnhub api auth
	finnhubClient := finnhub.NewAPIClient(finnhub.NewConfiguration()).DefaultApi
	finnhubAuth := context.WithValue(context.Background(), finnhub.ContextAPIKey, finnhub.APIKey{
		Key: os.Getenv("FINNHUB_API_KEY"),
	})

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(genMessageCreate(records, finnhubClient, finnhubAuth))

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
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func genMessageCreate(records [][]string, finnhubClient *finnhub.DefaultApiService, finnhubAuth context.Context) func(s *discordgo.Session, m *discordgo.MessageCreate) {
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

		if strings.HasPrefix(m.Content, "!quote") {
			resp, err := quote(m.Content, records, finnhubClient, finnhubAuth)
			if err != nil {
				log.Printf("Error: %s\n", err)
			}
			s.ChannelMessageSend(m.ChannelID, resp)
		}
		if strings.HasPrefix(m.Content, "!q") {
			resp, err := quote(m.Content, records, finnhubClient, finnhubAuth)
			if err != nil {
				log.Printf("Error: %s\n", err)
			}
			s.ChannelMessageSend(m.ChannelID, resp)
		}
	}
}

func quote(args string, records [][]string, finnhubClient *finnhub.DefaultApiService, finnhubAuth context.Context) (msg string, err error) {
	symbol := strings.ToUpper(strings.Split(args, " ")[1])

	log.Printf("Looking up stock quote: %s\n", symbol)
	detail, err := stonksV1.Quote(symbol, true, records, finnhubClient, finnhubAuth)
	if err != nil {
		log.Printf("Error getting stock quote %s", err)
		return "", err
	}
	log.Printf("%+v\n", detail)

	return detail.FormattedDetail, nil
}

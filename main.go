package main

import (
	"fmt"
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
	// initialize stonksClient
	stonksDataPath := os.Getenv("STONKS_DATA_PATH")
	key := os.Getenv("FINNHUB_API_KEY")
	sc := stonksV1.NewStonksClient(key, stonksDataPath)

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(genMessageCreate(sc))

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
func genMessageCreate(sc *stonksV1.StonksClient) func(s *discordgo.Session, m *discordgo.MessageCreate) {
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
			resp, err := quote(m.Content, sc)
			if err != nil {
				log.Printf("Error: %s\n", err)
			}
			s.ChannelMessageSend(m.ChannelID, resp)
		}
		if strings.HasPrefix(m.Content, "!short") {
			symbol := strings.Split(m.Content, " ")[1]
			resp, err := short(symbol, sc)
			if err != nil {
				log.Printf("Error: %s\n", err)
			}
			s.ChannelMessageSend(m.ChannelID, resp)
		}
		if strings.HasPrefix(m.Content, "!q") {
			symbols := strings.Split(strings.ToUpper(strings.Split(m.Content, " ")[1]), ",")
			for _, symbol := range symbols {
				resp, err := quote(symbol, sc)
				if err != nil {
					log.Printf("Error: %s\n", err)
				}
				s.ChannelMessageSend(m.ChannelID, resp)
			}
		}
	}
}

func short(symbol string, sc *stonksV1.StonksClient) (msg string, err error) {

	log.Printf("Looking up short interest on quote: %s\n", symbol)
	detail, err := sc.GetShortInterestBeta(symbol)
	if err != nil {
		log.Printf("Error getting short interest %s", err)
		return "", err
	}
	log.Printf("%+v\n", detail)
	res := "```"
	for _, item := range detail.Data {
		res += fmt.Sprintf("%s: %d\n", item.Date, item.ShortInterest)
	}
	res += "```"
	fmt.Printf("res: %s", res)

	return res, nil
}

func quote(symbol string, sc *stonksV1.StonksClient) (msg string, err error) {

	log.Printf("Looking up stock quote: %s\n", symbol)
	detail, err := sc.Quote(symbol)
	if err != nil {
		log.Printf("Error getting stock quote %s", err)
		return "", err
	}
	log.Printf("%+v\n", detail)

	return detail.FormattedDetail, nil
}

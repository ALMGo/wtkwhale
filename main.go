package main

import (
	"github.com/almaclaine/ethplorer"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/robfig/cron"
	"go.uber.org/zap"
	"math/rand"
	"os"
	"strconv"
	"time"
)

var logger *zap.Logger

var WTK_ADDRESS = "0xdf9d4674a430bdcc096a3a403128357ab36844ba"
var ETHPLORER_KEY = os.Getenv("ETHPLORER_KEY")

var topHolders *ethplorer.TopTokenHolders
var config Config

type Credentials struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
}

func getClient(creds *Credentials) (*twitter.Client, error) {
	// Pass in your consumer key (API Key) and your Consumer Secret (API Secret)
	config := oauth1.NewConfig(creds.ConsumerKey, creds.ConsumerSecret)
	// Pass in your Access Token and your Access Token Secret
	token := oauth1.NewToken(creds.AccessToken, creds.AccessTokenSecret)

	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	// Verify Credentials
	verifyParams := &twitter.AccountVerifyParams{
		SkipStatus:   twitter.Bool(true),
		IncludeEmail: twitter.Bool(true),
	}

	// we can retrieve the user and verify if the credentials
	// we have used successfully allow us to log in!
	_, _, err := client.Accounts.VerifyCredentials(verifyParams)
	if err != nil {
		return nil, err
	}

	return client, nil
}

var credentials = Credentials{
	ConsumerKey:       os.Getenv("TWITTER_API_KEY"),
	ConsumerSecret:    os.Getenv("TWITTER_API_SECRET"),
	AccessToken:       os.Getenv("TWITTER_ACCESS_TOKEN"),
	AccessTokenSecret: os.Getenv("TWITTER_ACCESS_SECRET"),
}

var tweets = make([]string, 0)

func main() {
	rand.Seed(time.Now().UnixNano())
	db, err := sqlx.Connect("sqlite3", "wtkwhales.db")
	topHolders, _ = updateTopHolders(db)
	logger, _ = zap.NewProduction()

	if err != nil {
		logger.Error("error Connecting To Database")
		os.Exit(1)
	}

	config, err = setupConfig(db)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	client, err := getClient(&credentials)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	c := cron.New()

	// Get Top holders every hour
	c.AddFunc("0 0 * * * *", func() {
		logger.Info("Getting Top Holders")
		topHolders, _ = updateTopHolders(db)
	})

	// watchWhales every minute
	c.AddFunc("0 * * * * *", func() {
		logger.Info("Watching Whales")
		watchWhales(db)
		if len(tweets) > 0 {
			tweet := tweets[len(tweets) - 1]
			tweets = tweets[:len(tweets)-1]
			err = sendTweet(client, tweet)
			if err != nil {
				logger.Error(err.Error())
			}
			logger.Info("Just tweeted, length of remaining tweets: " + strconv.Itoa(len(tweets)))
		}
	})

	//watchWhales(db)

	c.Start()
	select {}
}

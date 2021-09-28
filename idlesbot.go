package main

import (
	// other imports

	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

// TweetOdds represents how often the bot will run and Tweet a lyric. Setting this to 7, and running the bot every hour
// means that the bot has a 1 in 8 chance (0 thru 7 are eight numbers) of running that hour.
var TweetOdds = 10

// Credentials stores all of our access/consumer tokens
// and secret keys needed for authentication against
// the twitter REST API.
type Credentials struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
}

// getClient is a helper function that will return a twitter client
// that we can subsequently use to send tweets, or to stream new tweets
// this will take in a pointer to a Credential struct which will contain
// everything needed to authenticate and return a pointer to a twitter Client
// or an error
func getClient(creds *Credentials) (*twitter.Client, error) {
	config := oauth1.NewConfig(creds.ConsumerKey, creds.ConsumerSecret)
	token := oauth1.NewToken(creds.AccessToken, creds.AccessTokenSecret)

	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	// Verify Credentials
	verifyParams := &twitter.AccountVerifyParams{
		SkipStatus:   twitter.Bool(true),
		IncludeEmail: twitter.Bool(true),
	}

	user, _, err := client.Accounts.VerifyCredentials(verifyParams)
	if err != nil {
		return nil, err
	}

	log.Printf("User's ACCOUNT:\n%+v\n", user)
	return client, nil
}

// getLyric picks a line from the given file at random and returns it
func getLyric() string {
	rand.Seed(time.Now().UnixNano())

	lyricfile, err := ioutil.ReadFile("idles.txt")
	if err != nil {
		panic(err)
	}
	lyrics := strings.Split(string(lyricfile), "\n")

	length := len(lyrics)
	r := rand.Intn(length)
	line := lyrics[r]
	lyric := strings.Replace(string(line), "/ ", "\n", -1)
	return lyric
}

func main() {
	creds := Credentials{
		AccessToken:       os.Getenv("ACCESS_TOKEN"),
		AccessTokenSecret: os.Getenv("ACCESS_TOKEN_SECRET"),
		ConsumerKey:       os.Getenv("API_KEY"),
		ConsumerSecret:    os.Getenv("API_SECRET_KEY"),
	}

	// We roll a dice here and if the result is a 0, we proceed. For any other number, we exit.
	rand.Seed(time.Now().UnixNano())
	dice := rand.Intn(TweetOdds)
	if dice != 0 {
		log.Println("Not tweeting this time, dice roll was", dice)
		os.Exit(0)
	}

	// We rolled a 0, so let's tweet.
	client, err := getClient(&creds)
	if err != nil {
		log.Println("Error getting Twitter Client")
		log.Println(err)
	}

	lyric := getLyric()
	log.Println(lyric)

	tweet, resp, err := client.Statuses.Update(lyric, nil)
	if err != nil {
		log.Println(err)
	}
	log.Printf("%+v\n", resp)
	log.Printf("%+v\n", tweet)
	os.Exit(0)
}

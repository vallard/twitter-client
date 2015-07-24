package twitstream

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/garyburd/go-oauth/oauth"
	"github.com/joeshaw/envdecode"
)

type TwitStream struct {
	authClient *oauth.Client
	creds      *oauth.Credentials
	conn       net.Conn
	httpClient *http.Client
	reader     io.ReadCloser
}

// this can be augmented with more fields.
// see: https://dev.twitter.com/rest/reference/post/statuses/update
type user struct {
	Screen_Name string
}

type tweet struct {
	User user
	Text string
}

func (t *TwitStream) CloseConn() {
	if t.conn != nil {
		t.conn.Close()
	}
	fmt.Println("Closing t.reader")
	if t.reader != nil {
		t.reader.Close()
	}
}

func (t *TwitStream) Get(stuff string) {
	form := url.Values{"track": strings.Split(stuff, ",")}
	formEnc := form.Encode()
	u, _ := url.Parse("https://stream.twitter.com/1.1/statuses/filter.json")
	req, err := http.NewRequest("POST", u.String(), strings.NewReader(formEnc))
	if err != nil {
		log.Println("creating filter request failed:", err)
	}
	req.Header.Set("Authorization", t.authClient.AuthorizationHeader(t.creds, "POST", u, form))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(formEnc)))
	fmt.Println("Doing request")
	resp, err := t.httpClient.Do(req)
	if err != nil {
		log.Println("StatusCode =", resp.StatusCode)
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Println("StatusCode =", resp.StatusCode)
		return
	}
	log.Println("StatusCode: ", resp.StatusCode)
	t.reader = resp.Body
	decoder := json.NewDecoder(t.reader)
	for {
		var tweet tweet

		if err := decoder.Decode(&tweet); err == nil {
			fmt.Println(tweet.User.Screen_Name, ":", tweet.Text)
		} else {
			break
		}
	}
}

func New() *TwitStream {
	var ts struct {
		ConsumerKey    string `env:"STRIPSTOCK_TWITTER_CONSUMER_KEY,required"`
		ConsumerSecret string `env:"STRIPSTOCK_TWITTER_CONSUMER_SECRET,required"`
		AccessToken    string `env:"STRIPSTOCK_TWITTER_ACCESS_TOKEN,required"`
		AccessSecret   string `env:"STRIPSTOCK_TWITTER_ACCESS_SECRET,required"`
	}
	var conn net.Conn

	if err := envdecode.Decode(&ts); err != nil {
		log.Fatalln(err)
	}
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				if conn != nil {
					conn.Close()
					conn = nil
				}
				netc, err := net.DialTimeout(netw, addr, 5*time.Second)
				if err != nil {
					return nil, err
				}
				conn = netc
				return netc, nil
			},
		},
	}

	var r io.ReadCloser

	return &TwitStream{
		authClient: &oauth.Client{
			Credentials: oauth.Credentials{
				Token:  ts.ConsumerKey,
				Secret: ts.ConsumerSecret,
			},
		},
		creds: &oauth.Credentials{
			Token:  ts.AccessToken,
			Secret: ts.AccessSecret,
		},
		conn:       conn,
		httpClient: client,
		reader:     r,
	}
}

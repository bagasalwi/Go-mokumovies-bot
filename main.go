package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"unicode/utf8"

	"github.com/bwmarrin/discordgo"
)

// Varibel yang digunakan
var (
	Token    string
	APIToken string
	config   *Config // ambil data dari Config struct
)

type Config struct {
	Token    string `json:"token"`
	APIToken string `json:"api_token"`
}

func getToken() error {
	fmt.Println("get Token file...")
	file, err := ioutil.ReadFile("./config.json") // ioutil buat baca json file.

	//Buat handle err dari ioutil
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// Unmarshal hasil filenya ke var config
	err = json.Unmarshal(file, &config)

	//HHandle error dari unmarshal
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// taro tokennya ke var Token & APIToken
	Token = config.Token
	APIToken = config.APIToken
	fmt.Println("Token : " + Token)
	fmt.Println("APIToken : " + APIToken)

	// kalo ndak ada error nil
	return nil
}

const goAPIURL = "https://www.omdbapi.com/?apikey="

func main() {

	err := getToken()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

type MovieDetail struct {
	Title    string `json:"Title"`
	Year     string `json:"Year"`
	Rated    string `json:"Rated"`
	Released string `json:"Released"`
	Runtime  string `json:"Runtime"`
	Genre    string `json:"Genre"`
	Director string `json:"Director"`
	Writer   string `json:"Writer"`
	Actors   string `json:"Actors"`
	Plot     string `json:"Plot"`
	Language string `json:"Language"`
	Country  string `json:"Country"`
	Awards   string `json:"Awards"`
	Poster   string `json:"Poster"`
	Ratings  []struct {
		Source string `json:"Source"`
		Value  string `json:"Value"`
	} `json:"Ratings"`
	Metascore  string `json:"Metascore"`
	ImdbRating string `json:"imdbRating"`
	ImdbVotes  string `json:"imdbVotes"`
	ImdbID     string `json:"imdbID"`
	Type       string `json:"Type"`
	Dvd        string `json:"DVD"`
	BoxOffice  string `json:"BoxOffice"`
	Production string `json:"Production"`
	Website    string `json:"Website"`
	Response   string `json:"Response"`
}
type OmdbAPI struct {
	Search []struct {
		Title  string `json:Title`
		Year   string `json:Year`
		ImdbID string `json:imdbID`
		Type   string `json:Type`
		Poster string `json:Poster`
	}
	Result   string `json:totalResults`
	Response string `json:Response`
}

// fungsi utama messageCreate
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}

	var message strings.Builder
	var resmovies strings.Builder

	// !find The Batman
	findfilm := strings.Contains(m.Content, "!find")
	// !detail tt0372784
	findid := strings.Contains(m.Content, "!detail")

	if findid {
		movieid := strings.Replace(m.Content, "!detail", "", -1)
		// di trim dlu karena ga pake space
		countChecker := utf8.RuneCountInString(strings.TrimSpace(movieid))

		fmt.Println(countChecker)
		if countChecker <= 9 {
			movieDetailResponse, err := http.Get(goAPIURL + APIToken + "&i=" + strings.TrimSpace(movieid))
			if err != nil {
				fmt.Println(err)
			}

			defer movieDetailResponse.Body.Close()

			// rubah movieResponse jadi Byte bakal di proses unmarshall
			bodyMovie, err := ioutil.ReadAll(movieDetailResponse.Body)
			if err != nil {
				fmt.Println(err)
			}

			if movieDetailResponse.StatusCode == 200 {
				var res MovieDetail

				err := json.Unmarshal(bodyMovie, &res)
				if err != nil {
					panic(err)
				}

				if res.Response == "True" {
					// Message untuk info detail film
					message.WriteString("Title : " + res.Title + "\n")
					message.WriteString("Released : " + res.Released + "\n")
					message.WriteString("Plot : " + res.Plot + "\n")
					message.WriteString("Genre : " + res.Genre + "\n")
					message.WriteString("Rating IMDB : " + res.ImdbRating + "\n")

					resp, err := http.Get(res.Poster)
					if err != nil {
						fmt.Println(err)
					}
					defer resp.Body.Close()

					_, err = s.ChannelFileSend(m.ChannelID, res.ImdbID+".png", resp.Body)
					if err != nil {
						fmt.Println(err)
					}
				} else {
					message.WriteString("Maap mint salah itu idnya :')")
				}
				s.ChannelMessageSend(m.ChannelID, message.String())
			} else {
				s.ChannelMessageSend(m.ChannelID, "Maap mint salah itu idnya :')")
			}
		} else {
			s.ChannelMessageSend(m.ChannelID, "Maap mint salah itu idnya :')")
		}

	}

	if findfilm {
		movieName := strings.Replace(m.Content, "!find", "", -1)

		movieResponse, err := http.Get(goAPIURL + APIToken + "&s=" + strings.TrimSpace(movieName))
		if err != nil {
			fmt.Println(err)
		}
		defer movieResponse.Body.Close()

		// rubah movieResponse jadi Byte bakal di proses unmarshall
		body, err := ioutil.ReadAll(movieResponse.Body)
		if err != nil {
			fmt.Println(err)
		}

		if movieResponse.StatusCode == 200 {
			var res OmdbAPI

			err := json.Unmarshal(body, &res)
			if err != nil {
				panic(err)
			}

			// fast response before looping
			message.WriteString("Nih hasil pencarian Moku untuk " + strings.TrimSpace(movieName) + " | Result : " + res.Result + "\n")
			s.ChannelMessageSend(m.ChannelID, message.String())
			message.Reset()

			// fmt.Printf("TotalResults : %s", res.Result+"\n")
			// fmt.Printf("Response : %s", res.Response+"\n")

			if res.Response == "True" {
				// Loop data dari arrayJson
				for _, searchData := range res.Search {
					resmovies.WriteString("• " + searchData.Title + " | imdbID : " + searchData.ImdbID + "\n")
				}

				_, err = s.ChannelMessageSend(m.ChannelID, resmovies.String())
				if err != nil {
					fmt.Println(err)
				}

				// message.WriteString("Total Result :  " + res.Result)
			} else {
				message.WriteString("Maap mint judulnya ndak ketemu :P")
			}

			s.ChannelMessageSend(m.ChannelID, message.String())
		} else {
			message.WriteString("Maap mint judulnya ndak ketemu :P")
			s.ChannelMessageSend(m.ChannelID, message.String())
		}
	}

	if m.Content == "!mokuhelp" {
		message.WriteString("Perkenalkan ini adalah bot buatan moku yang menggunakan Go & OMDBAPI (IMDB unofficial API)!! \n!mokuhelp - main info \n!moku - for surprise!! \n!find <nama film>  - cari film berbentuk list \n!detail <imdbID> - cari film berbentuk detail")

		s.ChannelMessageSend(m.ChannelID, message.String())
	}

	if m.Content == "!moku" {
		message.WriteString("Pembuatnya ganteng bukan main! follow @mokultur / https://mokumedia.online :D")

		s.ChannelMessageSend(m.ChannelID, message.String())
	}

}

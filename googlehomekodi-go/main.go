package main

import (
	"errors"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"github.com/sevren/googlehome/googlehomekodi-go/api"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type reqBody struct {
	Token string
}

type kodi struct {
	Id       string
	Protocol string
	Ip       string
	Port     int
	User     string
	Pass     string
}

type config struct {
	Kodi          kodi
	AuthToken     string
	ListenerPort  int
	YoutubeApiKey string
}

var cfg *config

// Only using this function to get the outbound ip addr -- https://stackoverflow.com/a/37382208
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

/**
* Define a set of matchers that map the requests to functions
* A webrequest can have 0, 1 or 2 parameters we need to match for the appropriate api call
**/

//SimpleMatchers - No additional parameters needed to call the kodi api
func matcherSimple(s string) (func(), error) {

	switch s {
	case "/syncLibrary":
		return api.SyncLibrary, nil
	case "/navselect":
		return api.NavSelect, nil
	}
	return nil, errors.New("No Match!")

}

//TextParameter Matchers work on api calls that take 1 string parameter as input
func matcherTextParameter(s string) (func(s string), error) {

	switch s {
	case "/playtvshow":
		return api.HandleTvShow, nil
	}
	return nil, errors.New("No Match!")
}

//TextAndNumber Matcher work on api calls that 1 string and 1 number as input
func matcherTextAndNumberParameter(s string) (func(s string, d int), error) {

	switch s {
	case "/playepisode":
		return api.PlayEpisode, nil
	}

	return nil, errors.New("No Match!")

}

func extractTextParameter(s url.Values) (string, error) {

	queryParam := s.Get("q")

	if len(queryParam) > 0 {
		return queryParam, nil
	}
	return "", errors.New("Query Param 'q' expected!")
}

func extractNumberParameter(s url.Values) (int, error) {

	q := s.Get("e")

	if len(q) > 0 {

		n, err := strconv.Atoi(q)

		if err != nil {
			return 0, errors.New("Could not convert Query Parameter 'e' to integer")
		}
		return n, nil
	}
	return 0, errors.New("Query Param 'e' expected!")
}

//validates the body of the request for the token //TODO fixme
func validateBody(b string) error {
	return nil
}

func ExecuteReq(s string) {

	log.Infof("Executing Request to kodi instance", s)
	request := gorequest.New()
	resp, body, errs := request.Post(fmt.Sprintf("%s://%s:%d", cfg.Kodi.Protocol, cfg.Kodi.Ip, cfg.Kodi.Port)).
		Set("Notes", "gorequst is coming!").
		Send(`{"name":"backy", "species":"dog"}`).
		End()

	if errs != nil {
		log.Errorf("Posting to Kodi went wrong", errs)
		return
	}

	log.Info(resp, body)

}

//TODO Find a golang way to make this nicer..
func handle(w http.ResponseWriter, req *http.Request) {

	p := req.URL.Path
	log.Infof("URL Path: %s", p)

	//validation must occur first otherwise exit early
	if err := validateBody(""); err != nil {
		log.Errorf("Authentication Error %s", err)
		return
	}
	//check to see if there is a parameter
	query := req.URL.Query()

	switch len(query) {
	case 0:
		f, err := matcherSimple(p)
		if err != nil {
			log.Errorf("Simple Matcher error: \n", err)
			return
		}
		f()
	case 1:
		f, err := matcherTextParameter(p)
		if err != nil {
			log.Errorf("Text Matcher error: \n", err)
			return
		}

		q, err := extractTextParameter(query)
		if err != nil {
			log.Error(err)
			return
		}
		f(q)
	case 2:
		f, err := matcherTextAndNumberParameter(p)
		if err != nil {
			log.Errorf("Text And Number matcher error", err)
			return
		}

		q, err := extractTextParameter(query)
		if err != nil {
			log.Error(err)
			return
		}
		n, err := extractNumberParameter(query)
		f(q, n)
	}

}

func setupViper() {

	viper.BindPFlags(flag.CommandLine)
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()
}

func main() {

	_ = flag.String("configPath", "./", "Sets the configuration path for the application configuration")
	flag.Parse()

	setupViper()

	viper.SetConfigName("cfg")
	viper.AddConfigPath(viper.GetString("configPath"))

	err := viper.ReadInConfig()

	if err != nil {
		log.Fatalf("Unable to read configuration file %s", err)
	}

	err = viper.Unmarshal(&cfg)

	if err != nil {
		log.Fatalf("Unable to unmarshal the configuration file: %v, exiting!", err)
	}

	log.Infof("Application started with the following cfg: %+v", cfg)

	log.Infof("Booting up Google Home Kodi Go Client on %s:%d", GetOutboundIP(), cfg.ListenerPort)

	http.HandleFunc("/", handle)
	err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.ListenerPort), nil)
	if err != nil {
		panic(err)
	}
}

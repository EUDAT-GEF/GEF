package server

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/EUDAT-GEF/GEF/gefserver/def"
	"golang.org/x/oauth2"
)

func init() {
	if os.Getenv("GEF_B2ACCESS_CONSUMER_KEY") == "" {
		log.Println("ERROR: GEF_B2ACCESS_CONSUMER_KEY environment variable not found")
	}
	if os.Getenv("GEF_B2ACCESS_SECRET_KEY") == "" {
		log.Println("ERROR: GEF_B2ACCESS_SECRET_KEY environment variable not found")
	}
}

var oauthCfg oauth2.Config
var userInfoURL string

func initB2Access(cfg def.B2AccessConfig) {
	if !strings.HasSuffix(cfg.BaseURL, "/") {
		cfg.BaseURL += "/"
	}
	oauthCfg = oauth2.Config{
		ClientID:     os.Getenv("GEF_B2ACCESS_CONSUMER_KEY"),
		ClientSecret: os.Getenv("GEF_B2ACCESS_SECRET_KEY"),
		Endpoint: oauth2.Endpoint{
			AuthURL:  cfg.BaseURL + "oauth2-as/oauth2-authz",
			TokenURL: cfg.BaseURL + "oauth2/token",
		},
		RedirectURL: cfg.RedirectURL,
		Scopes:      []string{"USER_PROFILE", "GENERATE_USER_CERTIFICATE", "email", "profile"},
	}
	userInfoURL = cfg.BaseURL + "oauth2/userinfo"
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	session, err := cookieStore.Get(r, sessionName)
	if err != nil {
		str := fmt.Sprintf("Cookie store error: %s", err.Error())
		log.Println(str)
		http.Error(w, str, 500)
		return
	}
	Response{w}.Ok(jmap("name", session.Values["userName"], "email", session.Values["userEmail"]))
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	statebuf := make([]byte, 16)
	rand.Read(statebuf)
	state := base64.URLEncoding.EncodeToString(statebuf)
	// log.Println("created random state: ", state)

	session, err := cookieStore.Get(r, sessionName)
	if err != nil {
		str := fmt.Sprintf("Cookie store error: %s", err.Error())
		log.Println(str)
		http.Error(w, str, 500)
		return
	}
	session.Values["state"] = state
	session.Save(r, w)

	url := oauthCfg.AuthCodeURL(state)
	// log.Println("redirect to url: ", url)
	http.Redirect(w, r, url, 302)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	session, err := cookieStore.Get(r, sessionName)
	if err != nil {
		str := fmt.Sprintf("Cookie store error: %s", err.Error())
		log.Println(str)
		http.Error(w, str, 500)
		return
	}

	// log.Println("received state: ", r.URL.Query().Get("state"))
	if r.URL.Query().Get("state") != session.Values["state"] {
		str := fmt.Sprintf("State mismatch, CSRF or cookies not enabled?")
		log.Println(str)
		http.Error(w, str, 500)
		return
	}

	// log.Println("received code: ", r.URL.Query().Get("code"))
	tkn, err := oauthCfg.Exchange(oauth2.NoContext, r.URL.Query().Get("code"))
	if err != nil {
		str := fmt.Sprintf("Error while getting access token: %s", err.Error())
		log.Println(str)
		http.Error(w, str, 500)
		return
	}

	// log.Println("exchanged for token: ", tkn)
	if !tkn.Valid() {
		str := fmt.Sprintf("Received invalid access token: %s", tkn)
		log.Println(str)
		http.Error(w, str, 500)
		return
	}

	resp, err := oauthCfg.Client(oauth2.NoContext, tkn).Get(userInfoURL)
	if err != nil {
		str := fmt.Sprintf("Error while getting user info: %s", err.Error())
		log.Println(str)
		http.Error(w, str, 500)
		return
	}

	decoder := json.NewDecoder(resp.Body)
	type userType struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	var user userType
	err = decoder.Decode(&user)
	if err != nil {
		str := fmt.Sprintf("Error while deserialising user info: %s", err.Error())
		log.Println(str)
		http.Error(w, str, 500)
		return
	}

	session.Values["userName"] = user.Name
	session.Values["userEmail"] = user.Email
	session.Save(r, w)

	http.Redirect(w, r, "/", 302)
}

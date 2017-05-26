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
		Response{w}.ServerError("Cookie store error", err)
		return
	}

	userObject := jmap(
		"Name", session.Values["userName"],
		"Email", session.Values["userEmail"],
		"Error", session.Values["userError"])
	Response{w}.Ok(jmap("User", userObject))
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	statebuf := make([]byte, 16)
	rand.Read(statebuf)
	state := base64.URLEncoding.EncodeToString(statebuf)

	session, err := cookieStore.Get(r, sessionName)
	if err != nil {
		Response{w}.ServerError("Cookie store error", err)
		return
	}
	session.Values["state"] = state
	session.Save(r, w)

	url := oauthCfg.AuthCodeURL(state)
	http.Redirect(w, r, url, 302)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	session, err := cookieStore.Get(r, sessionName)
	if err != nil {
		Response{w}.ServerError("Cookie store error", err)
		return
	}

	userError := func(str string) {
		session.Values["userError"] = str
		delete(session.Values, "userName")
		delete(session.Values, "userEmail")
		session.Save(r, w)
	}

	// log.Println("received state: ", r.URL.Query().Get("state"))
	if r.URL.Query().Get("state") != session.Values["state"] {
		userError("State mismatch, CSRF or cookies not enabled?")
		return
	}

	// log.Println("received code: ", r.URL.Query().Get("code"))
	tkn, err := oauthCfg.Exchange(oauth2.NoContext, r.URL.Query().Get("code"))
	if err != nil {
		userError(fmt.Sprintf("Error while getting access token: %s", err))
		return
	}

	// log.Println("exchanged for token: ", tkn)
	if !tkn.Valid() {
		userError(fmt.Sprintf("Received invalid access token: %s", tkn))
		return
	}

	resp, err := oauthCfg.Client(oauth2.NoContext, tkn).Get(userInfoURL)
	if err != nil {
		userError(fmt.Sprintf("Error while getting user info: %s", err.Error()))
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
		userError(fmt.Sprintf("Error while deserialising user info: %s", err.Error()))
		return
	}

	delete(session.Values, "userError")
	session.Values["userName"] = user.Name
	session.Values["userEmail"] = user.Email
	session.Save(r, w)

	http.Redirect(w, r, "/", 302)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := cookieStore.Get(r, sessionName)
	if err != nil {
		Response{w}.ServerError("Cookie store error", err)
		return
	}
	delete(session.Values, "userError")
	delete(session.Values, "userName")
	delete(session.Values, "userEmail")
	session.Save(r, w)
	http.Redirect(w, r, "/", 302)
}

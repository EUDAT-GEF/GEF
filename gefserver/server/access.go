package server

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/EUDAT-GEF/GEF/gefserver/db"
	"github.com/EUDAT-GEF/GEF/gefserver/def"
	"github.com/gorilla/mux"
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

func (s *Server) userHandler(w http.ResponseWriter, r *http.Request) {
	user, err := s.getLoggedInUser(r)
	if err != nil {
		Response{w}.Ok(jmap("Error", err.Error()))
	} else if user != nil {
		Response{w}.Ok(jmap("User", user))
	} else {
		Response{w}.Ok("{}")
	}
}

func (s *Server) logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := cookieStore.Get(r, sessionName)
	if err != nil {
		Response{w}.ServerError("Cookie store error", err)
		return
	}
	delete(session.Values, "userID")
	session.Save(r, w)
	http.Redirect(w, r, "/", 302)
}

func (s *Server) newTokenHandler(w http.ResponseWriter, r *http.Request) {
	tokenName := r.FormValue("tokenName")
	if tokenName == "" {
		vars := mux.Vars(r)
		tokenName = vars["tokenName"]
	}
	logParam("tokenName", tokenName)

	user, err := s.getLoggedInUser(r)
	if err != nil {
		Response{w}.ServerError("User error", err)
		return
	} else if user == nil {
		Response{w}.Unauthorized()
		return
	}

	expire := time.Now().AddDate(10, 0, 0) // 10 years from now on
	token, err := s.db.NewUserToken(user.ID, tokenName, expire)
	if err != nil {
		Response{w}.ServerError("cannot create new token", err)
		return
	}
	Response{w}.Created(jmap("Token", token))
}

func (s *Server) listTokenHandler(w http.ResponseWriter, r *http.Request) {
	user, err := s.getLoggedInUser(r)
	if err != nil {
		Response{w}.ServerError("User error", err)
		return
	} else if user == nil {
		Response{w}.Unauthorized()
		return
	}

	tokens, err := s.db.GetUserTokens(user.ID)
	if err != nil {
		Response{w}.ServerError("cannot list tokens", err)
		return
	}
	for i := range tokens {
		tokens[i].Secret = tokens[i].Secret[:8] + "..."
	}
	Response{w}.Ok(jmap("Tokens", tokens))
}

func (s *Server) removeTokenHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenIDStr := vars["tokenID"]
	tokenID, err := strconv.ParseInt(tokenIDStr, 10, 64)
	if err != nil {
		Response{w}.ClientError("token id must be an int", err)
		return
	}
	token, err := s.db.GetTokenByID(tokenID)
	if err != nil {
		Response{w}.ClientError("token error", err)
		return
	}

	user, err := s.getLoggedInUser(r)
	if err != nil {
		Response{w}.ServerError("User error", err)
		return
	} else if user == nil {
		Response{w}.Unauthorized()
		return
	}

	if token.UserID != user.ID {
		Response{w}.Forbidden()
		return
	}
	err = s.db.DeleteUserToken(user.ID, token.ID)
	if err != nil {
		Response{w}.ServerError("delete token error", err)
		return
	}
	Response{w}.Ok(jmap("Token", token))
}

func (s *Server) getLoggedInUser(r *http.Request) (*db.User, error) {
	session, err := cookieStore.Get(r, sessionName)
	if err != nil {
		return nil, def.Err(err, "Cookie store error")
	}

	userIDValue := session.Values["userID"]
	if userIDValue == nil {
		return nil, nil // no logged in user, but also no error
	}

	userID, ok := userIDValue.(int64)
	if !ok {
		return nil, def.Err(err, "Session userID is not an integer")
	}

	user, err := s.db.GetUserByID(userID)
	if err != nil {
		return nil, def.Err(err, "GetUserByID error")
	}
	return user, nil
}

///////////////////////////////////////////////////////////////////////////////

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

func (s *Server) oauthLoginHandler(w http.ResponseWriter, r *http.Request) {
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

func (s *Server) oauthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	session, err := cookieStore.Get(r, sessionName)
	if err != nil {
		Response{w}.ServerError("Cookie store error", err)
		return
	}

	// log.Println("received state: ", r.URL.Query().Get("state"))
	if r.URL.Query().Get("state") != session.Values["state"] {
		Response{w}.ServerError("State mismatch, CSRF or cookies not enabled?", nil)
		return
	}

	// log.Println("received code: ", r.URL.Query().Get("code"))
	tkn, err := oauthCfg.Exchange(oauth2.NoContext, r.URL.Query().Get("code"))
	if err != nil {
		Response{w}.ServerError("Error while getting access token", err)
		return
	}

	// log.Println("exchanged for token: ", tkn)
	if !tkn.Valid() {
		Response{w}.ServerError(fmt.Sprintf("Received invalid access token: %s", tkn), nil)
		return
	}

	resp, err := oauthCfg.Client(oauth2.NoContext, tkn).Get(userInfoURL)
	if err != nil {
		Response{w}.ServerError("Error while getting user info", err)
		return
	}

	type userInfoType struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	var userInfo userInfoType
	err = json.NewDecoder(resp.Body).Decode(&userInfo)
	if err != nil {
		Response{w}.ServerError("Error while deserialising user info", err)
		return
	}

	user, err := s.db.GetUserByEmail(userInfo.Email)
	if err != nil {
		Response{w}.ServerError("Database error while retrieving user info", err)
		return
	}
	if user == nil {
		// user not found
		log.Println("new user: ", userInfo)
		now := time.Now()
		user = &db.User{
			Name:    userInfo.Name,
			Email:   userInfo.Email,
			Created: now,
			Updated: now,
		}
		_user, err := s.db.AddUser(*user)
		if err != nil {
			Response{w}.ServerError("Error while adding user to db", err)
			return
		}
		user = &_user
	}

	if user.Name != userInfo.Name {
		user.Name = userInfo.Name
		user.Updated = time.Now()
		s.db.UpdateUser(*user)
	}

	session.Values["userID"] = user.ID
	session.Save(r, w)
	http.Redirect(w, r, "/", 302)
}

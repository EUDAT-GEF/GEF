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

var (
	AccessTokenQueryParam = "access_token"
	AccessTokenCookieKey  = "UIAccessToken"
)

func init() {
	if os.Getenv("GEF_B2ACCESS_CONSUMER_KEY") == "" {
		log.Println("ERROR: GEF_B2ACCESS_CONSUMER_KEY environment variable not found")
	}
	if os.Getenv("GEF_B2ACCESS_SECRET_KEY") == "" {
		log.Println("ERROR: GEF_B2ACCESS_SECRET_KEY environment variable not found")
	}
}

// user  related handlers

func (s *Server) userHandler(w http.ResponseWriter, r *http.Request) {
	user, err := s.getCurrentUser(r)
	if err != nil {
		Response{w}.Ok(jmap("Error", err.Error()))
		return
	}
	if user == nil {
		Response{w}.Ok("{}")
		return
	}

	roles, err := s.db.GetUserRoles(user.ID)
	if err != nil {
		Response{w}.ServerError("Get user roles error", err)
		return
	}
	Response{w}.Ok(jmap(
		"User", user,
		"IsSuperAdmin", s.isSuperAdmin(user),
		"Roles", roles))
}

func (s *Server) logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := cookieStore.Get(r, sessionName)
	if err != nil {
		Response{w}.ServerError("Cookie store error", err)
		return
	}
	delete(session.Values, AccessTokenCookieKey)
	session.Save(r, w)
	http.Redirect(w, r, "/", 302)
}

func (s *Server) newTokenHandler(w http.ResponseWriter, r *http.Request) {
	allow, user := Authorization{s, w, r}.allowCreateToken()
	if user == nil || !allow {
		return
	}

	tokenName := r.FormValue("tokenName")
	if tokenName == "" {
		vars := mux.Vars(r)
		tokenName = vars["tokenName"]
	}
	logParam("tokenName", tokenName)

	expire := time.Now().AddDate(10, 0, 0) // 10 years from now on
	token, err := s.db.NewUserToken(user.ID, tokenName, expire)
	if err != nil {
		Response{w}.ServerError("cannot create new token", err)
		return
	}
	Response{w}.Created(jmap("Token", token))
}

func (s *Server) listTokenHandler(w http.ResponseWriter, r *http.Request) {
	allow, user := Authorization{s, w, r}.allowGetTokens()
	if user == nil || !allow {
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
		Response{w}.ClientError("tokenID must be an int", err)
		return
	}
	token, err := s.db.GetTokenByID(tokenID)
	if err != nil {
		Response{w}.ClientError("db get token error", err)
		return
	}

	allow, user := Authorization{s, w, r}.allowDeleteToken(token)
	if !allow {
		return
	}

	err = s.db.DeleteUserToken(user.ID, token.ID)
	if err != nil {
		Response{w}.ServerError("delete token error", err)
		return
	}
	Response{w}.Ok("")
}

///////////////////////////////////////////////////////////////////////////////

// role related handlers

func (s *Server) listRolesHandler(w http.ResponseWriter, r *http.Request) {
	allow, _ := Authorization{s, w, r}.allowListRoles()
	if !allow {
		return
	}

	roles, err := s.db.ListRoles()
	if err != nil {
		Response{w}.ServerError("db error", err)
		return
	}
	Response{w}.Ok(jmap("Roles", roles))
}

func (s *Server) listRoleUsersHandler(w http.ResponseWriter, r *http.Request) {
	allow, _ := Authorization{s, w, r}.allowListRoleUsers()
	if !allow {
		return
	}

	vars := mux.Vars(r)
	roleIDStr := vars["roleID"]
	roleID, err := strconv.ParseInt(roleIDStr, 10, 64)
	if err != nil {
		Response{w}.ClientError("roleID must be an int", err)
		return
	}

	roles, err := s.db.GetRoleUsers(roleID)
	if err != nil {
		Response{w}.ServerError("db error", err)
		return
	}
	Response{w}.Ok(jmap("Users", roles))
}

func (s *Server) newRoleUserHandler(w http.ResponseWriter, r *http.Request) {
	allow, _ := Authorization{s, w, r}.allowNewRoleUser()
	if !allow {
		return
	}

	vars := mux.Vars(r)
	roleIDStr := vars["roleID"]
	roleID, err := strconv.ParseInt(roleIDStr, 10, 64)
	if err != nil {
		Response{w}.ClientError("roleID must be an int", err)
		return
	}
	userEmail := r.FormValue("userEmail")
	logParam("userEmail", userEmail)
	user, err := s.db.GetUserByEmail(userEmail)
	if err != nil {
		Response{w}.ServerError("db error", err)
		return
	}
	if user == nil {
		Response{w}.ClientError("User not found", nil)
		return
	}

	err = s.db.AddRoleToUser(user.ID, roleID)
	if err != nil {
		Response{w}.ServerError("db error", err)
		return
	}

	Response{w}.Ok("")
}

func (s *Server) removeRoleUserHandler(w http.ResponseWriter, r *http.Request) {
	allow, _ := Authorization{s, w, r}.allowDeleteRoleUser()
	if !allow {
		return
	}

	vars := mux.Vars(r)
	roleIDStr := vars["roleID"]
	roleID, err := strconv.ParseInt(roleIDStr, 10, 64)
	if err != nil {
		Response{w}.ClientError("roleID must be an int", err)
		return
	}
	userIDStr := vars["userID"]
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		Response{w}.ClientError("userID must be an int", err)
		return
	}

	err = s.db.DeleteRoleFromUser(userID, roleID)
	if err != nil {
		Response{w}.ServerError("db error", err)
		return
	}

	Response{w}.Ok("")
}

///////////////////////////////////////////////////////////////////////////////

func (s *Server) isSuperAdmin(user *db.User) bool {
	if user.Email == s.administration.SuperAdminEmail {
		return true
	}
	return s.db.HasSuperAdminRole(user.ID)
}

func (s *Server) getCurrentUser(r *http.Request) (*db.User, error) {
	accessToken := ""

	session, err := cookieStore.Get(r, sessionName)
	if err != nil {
		return nil, def.Err(err, "Cookie store error")
	}

	accessTokenInterface := session.Values[AccessTokenCookieKey]
	if accessTokenInterface != nil {
		var ok bool
		accessToken, ok = accessTokenInterface.(string)
		if !ok {
			return nil, def.Err(err, "Bad cookie value type")
		}
		// log.Println("\tuser authenticated by cookie")
	}
	if accessToken == "" {
		accessTokenList, ok := r.URL.Query()[AccessTokenQueryParam]
		if !ok || len(accessTokenList) == 0 {
			// no access_token
			return nil, nil
		}
		if len(accessTokenList) > 1 {
			return nil, def.Err(err, "too many access token parameters")
		}
		accessToken = accessTokenList[0]
		if accessToken == "" {
			return nil, nil // no user, no error
		}
		// log.Println("\tuser authenticated by access_token")
	}

	token, err := s.db.GetTokenBySecret(accessToken)
	if err != nil || token.ID == 0 || token.Secret != accessToken {
		return nil, def.Err(err, "bad access token")
	}

	user, err := s.db.GetUserByID(token.UserID)
	if err != nil || user.ID == 0 {
		return nil, def.Err(err, "GetUserByID error")
	}

	return user, nil
}

// getUserOrWriteError returns the logged in user or writes errors into the http stream
func (s *Server) getUserOrWriteError(w http.ResponseWriter, r *http.Request) *db.User {
	user, err := s.getCurrentUser(r)
	if err != nil {
		Response{w}.ServerError("User error", err)
		return nil
	}
	if user == nil {
		Response{w}.Unauthorized()
		return nil
	}
	return user
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

	// extract user info provided by B2ACCESS
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

	// put user in the database if not already there
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

	// update user data if necessary
	if user.Name != userInfo.Name {
		user.Name = userInfo.Name
		user.Updated = time.Now()
		s.db.UpdateUser(*user)
	}

	// create access token for the UI if necessary
	tokenList, err := s.db.GetUserTokens(user.ID)
	if err != nil {
		Response{w}.ServerError("Error while retrieving user tokens", err)
		return
	}
	accessToken := ""
	for _, token := range tokenList {
		if token.Name == AccessTokenCookieKey {
			accessToken = token.Secret
			break
		}
	}
	if accessToken == "" {
		expire := time.Now().AddDate(100, 0, 0)
		token, err := s.db.NewUserToken(user.ID, AccessTokenCookieKey, expire)
		if err != nil {
			Response{w}.ServerError("Error while creating user token", err)
			return
		}
		accessToken = token.Secret
	}

	session.Values[AccessTokenCookieKey] = accessToken
	session.Save(r, w)
	http.Redirect(w, r, "/", 302)
}

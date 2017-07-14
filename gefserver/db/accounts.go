package db

import (
	"crypto/rand"
	"encoding/base64"
	"strings"
	"time"

	"github.com/EUDAT-GEF/GEF/gefserver/def"
)

// User struct, also used to serialize JSON
type User struct {
	ID      int64
	Name    string
	Email   string
	Created time.Time
	Updated time.Time
}

// Token struct, also used to serialize JSON
type Token struct {
	ID     int64
	Name   string
	Secret string
	UserID int64
	Expire time.Time
}

///////////////////////////////////////////////////////////////////////////////

// AddUser adds a user to the database
func (d *Db) AddUser(user User) (User, error) {
	u := d.user2UserTable(user)
	err := d.db.Insert(&u)
	user.ID = u.ID
	return user, err
}

// UpdateUser updates user information in the database
func (d *Db) UpdateUser(user User) error {
	u := d.user2UserTable(user)
	_, err := d.db.Update(&u)
	return err
}

// GetUserByID searches for a user in the database
func (d *Db) GetUserByID(id int64) (*User, error) {
	user, err := d.db.Get(UserTable{}, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil
	}
	return d.userTable2user(user.(*UserTable)), nil
}

// GetUserByEmail searches for a user in the database
func (d *Db) GetUserByEmail(email string) (*User, error) {
	var user UserTable
	err := d.db.SelectOne(&user, "SELECT * FROM users WHERE Email=?", email)
	if err != nil {
		if strings.HasSuffix(err.Error(), "no rows in result set") {
			return nil, nil
		}
		return nil, err
	}
	return d.userTable2user(&user), nil
}

// GetUserTokens returns the user's tokens
func (d *Db) GetUserTokens(userID int64) ([]Token, error) {
	var dbtokens []TokenTable
	_, err := d.db.Select(&dbtokens, "SELECT * FROM tokens WHERE UserID=?", userID)
	if err != nil {
		return nil, err
	}
	tokens := make([]Token, 0, 0)
	for _, t := range dbtokens {
		tokens = append(tokens, d.tokenTable2token(t))
	}
	return tokens, nil
}

// GetTokenByID returns a token by its ID
func (d *Db) GetTokenByID(tokenID int64) (Token, error) {
	var token TokenTable
	err := d.db.SelectOne(&token, "SELECT * FROM tokens WHERE ID=?", tokenID)
	if err != nil {
		return Token{}, err
	}
	return d.tokenTable2token(token), nil
}

// GetTokenBySecret returns a token by its Secret
func (d *Db) GetTokenBySecret(accessToken string) (Token, error) {
	var token TokenTable
	err := d.db.SelectOne(&token, "SELECT * FROM tokens WHERE Secret=?", accessToken)
	if err != nil {
		return Token{}, err
	}
	return d.tokenTable2token(token), nil
}

// NewUserToken creates a new user access token
func (d *Db) NewUserToken(userID int64, name string, expire time.Time) (Token, error) {
	// generate crypto random string
	buf := make([]byte, 36)
	_, err := rand.Read(buf)
	if err != nil {
		return Token{}, def.Err(err, "Cannot read crypto/rand")
	}

	token := TokenTable{
		Name:   name,
		Secret: base64.RawURLEncoding.EncodeToString(buf),
		UserID: userID,
		Expire: expire,
	}
	err = d.db.Insert(&token)
	if err != nil {
		return Token{}, err
	}
	return d.tokenTable2token(token), nil
}

// DeleteUserToken deletes the token from the database
func (d *Db) DeleteUserToken(userID, tokenID int64) error {
	_, err := d.db.Exec("DELETE FROM tokens WHERE UserID=? AND ID=?", userID, tokenID)
	return err
}

func (d *Db) user2UserTable(user User) UserTable {
	return UserTable{
		ID:      int64(user.ID),
		Name:    user.Name,
		Email:   user.Email,
		Created: user.Created,
		Updated: user.Updated,
	}
}

func (d *Db) userTable2user(user *UserTable) *User {
	if user == nil {
		return nil
	}
	return &User{
		ID:      user.ID,
		Name:    user.Name,
		Email:   user.Email,
		Created: user.Created,
		Updated: user.Updated,
	}
}

func (d *Db) tokenTable2token(token TokenTable) Token {
	return Token{
		ID:     token.ID,
		Name:   token.Name,
		Secret: token.Secret,
		UserID: token.UserID,
		Expire: token.Expire,
	}
}

///////////////////////////////////////////////////////////////////////////////

package db

import (
	"crypto/rand"
	"encoding/base64"
	"log"
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

type Community struct {
	ID          int64
	Name        string
	Description string
}

var (
	SuperAdminRoleName      = "SuperAdministrator"
	CommunityAdminRoleName  = "Administrator"
	CommunityMemberRoleName = "Member"
)

type Role struct {
	ID          int64
	Name        string
	CommunityID int64
	Description string
}

// GetCommunityByID gets a community from the database
func (d *Db) GetCommunityByID(communityID int64) (Community, error) {
	var c CommunityTable
	err := d.db.SelectOne(&c, "SELECT * FROM Communities WHERE ID=?", communityID)
	if err != nil {
		return Community{}, err
	}
	return commTable2comm(c), nil
}

// GetCommunityByName gets a community from the database
func (d *Db) GetCommunityByName(name string) (Community, error) {
	var c CommunityTable
	err := d.db.SelectOne(&c, "SELECT * FROM Communities WHERE Name=?", name)
	if err != nil {
		return Community{}, err
	}
	return commTable2comm(c), nil
}

// ListCommunities gets the list of all the communities
func (d *Db) ListCommunities() ([]Community, error) {
	var cs []CommunityTable
	_, err := d.db.Select(&cs, "SELECT * FROM Communities")
	if err != nil {
		return nil, err
	}

	var comms []Community
	for _, c := range cs {
		comms = append(comms, commTable2comm(c))
	}
	return comms, nil
}

func commTable2comm(c CommunityTable) Community {
	return Community{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
	}
}

// AddCommunity adds a community to the database, with checks
func (d *Db) AddCommunity(name, description string, addDefaultRoles bool) (Community, error) {
	comm, err := d.GetCommunityByName(name)
	if err == nil {
		return comm, nil
	}
	if !isNoResultsError(err) {
		return comm, err
	}

	c := CommunityTable{
		Name:        name,
		Description: description,
	}

	err = d.db.Insert(&c)
	if err != nil {
		return Community{}, err
	}
	if c.ID == 0 {
		panic("Community with ID 0 just inserted, which is illegal")
	}
	comm = Community{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
	}

	if addDefaultRoles {
		_, err := d.AddRole(CommunityAdminRoleName, comm.ID, "Community Administrator")
		if err != nil {
			return comm, err
		}
		_, err = d.AddRole(CommunityMemberRoleName, comm.ID, "Community Member")
		if err != nil {
			return comm, err
		}
	}

	return comm, nil
}

// GetRoleByName gets a named role from the database
func (d *Db) GetRoleByName(name string, communityID int64) (Role, error) {
	var r RoleTable
	err := d.db.SelectOne(&r, "SELECT * FROM roles WHERE Name=? AND CommunityID=?", name, communityID)
	if err != nil {
		return Role{}, err
	}
	return Role{
		ID:          r.ID,
		Name:        r.Name,
		CommunityID: r.CommunityID,
		Description: r.Description,
	}, nil
}

// AddRole adds a role to the database, but only if it's not already present
// For roles that have no associated community use communityID: 0
func (d *Db) AddRole(name string, communityID int64, description string) (Role, error) {
	if name == SuperAdminRoleName && communityID != 0 {
		return Role{}, def.Err(nil, "This role name is reserved")
	}
	role, err := d.GetRoleByName(name, communityID)
	if err == nil {
		return role, nil
	}
	if !isNoResultsError(err) {
		return role, err
	}

	r := RoleTable{
		Name:        name,
		CommunityID: communityID,
		Description: description,
	}

	err = d.db.Insert(&r)
	if err != nil {
		return Role{}, err
	}
	return Role{
		ID:          r.ID,
		Name:        r.Name,
		CommunityID: r.CommunityID,
		Description: r.Description,
	}, nil
}

// AddRoleToUser
func (d *Db) AddRoleToUser(userID, roleID int64) error {
	ur := UserRoleTable{
		UserID: userID,
		RoleID: roleID,
	}
	err := d.db.Insert(&ur)
	return err
}

// GetUserRoles
func (d *Db) GetUserRoles(userID int64) ([]Role, error) {
	var urs []UserRoleTable
	_, err := d.db.Select(&urs, "SELECT * FROM userroles WHERE UserID=? ORDER BY RoleID", userID)
	if err != nil {
		return nil, err
	}

	var roles []Role
	for _, ur := range urs {
		var r RoleTable
		err := d.db.SelectOne(&r, "SELECT * FROM roles WHERE ID=?", ur.RoleID)
		if err != nil {
			return roles, err
		}
		roles = append(roles, Role{
			ID:          r.ID,
			Name:        r.Name,
			CommunityID: r.CommunityID,
			Description: r.Description,
		})
	}
	return roles, nil
}

func (d *Db) HasSuperAdminRole(userID int64) bool {
	roles, err := d.GetUserRoles(userID)
	if err != nil {
		log.Println(def.Err(err, "HasSuperAdminRole/GetUserRoles failed"))
		return false
	}
	for _, r := range roles {
		if r.Name == SuperAdminRoleName && r.CommunityID == 0 {
			return true
		}
	}
	return false
}

func (d *Db) IsServiceOwner(userID int64, serviceID ServiceID) bool {
	var x OwnerTable
	err := d.db.SelectOne(&x, "SELECT * FROM owners WHERE UserID=? AND ObjectType=? AND ObjectID=?",
		userID, "Service", serviceID)
	return err == nil
}

func (d *Db) IsJobOwner(userID int64, jobID JobID) bool {
	var x OwnerTable
	err := d.db.SelectOne(&x, "SELECT * FROM owners WHERE UserID=? AND ObjectType=? AND ObjectID=?",
		userID, "Job", jobID)
	return err == nil
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

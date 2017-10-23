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
	ID            int64
	Name          string
	Description   string
	CommunityID   int64
	CommunityName string
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
		return comm, def.Err(err, "error in GetCommunityByName:"+name)
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
			return comm, def.Err(err, "error in AddRole/admin")
		}
		_, err = d.AddRole(CommunityMemberRoleName, comm.ID, "Community Member")
		if err != nil {
			return comm, def.Err(err, "error in AddRole/member")
		}
	}

	return comm, nil
}

//  ListRoles returns the list of all roles
func (d *Db) ListRoles() ([]Role, error) {
	var roles []Role
	_, err := d.db.Select(&roles,
		`SELECT r.ID ID, r.Name Name, r.Description Description
		FROM roles r
		WHERE r.CommunityID = ?`,
		0)
	if err != nil {
		return nil, err
	}

	var communityRoles []Role
	_, err = d.db.Select(&communityRoles,
		`SELECT r.ID ID, r.Name Name, r.Description Description, c.ID CommunityID, c.Name CommunityName
		FROM roles r, communities c
		WHERE r.CommunityID = c.ID`)
	if err != nil {
		return roles, err
	}
	roles = append(roles, communityRoles...)
	return roles, nil
}

// GetRoleByID gets a role from the database
func (d *Db) GetRoleUsers(roleID int64) ([]User, error) {
	var users []User
	_, err := d.db.Select(
		&users,
		`SELECT u.ID ID, u.Name Name, u.Email Email, u.Created Created, u.Updated Updated
		FROM roles r, userroles ur, users u
		WHERE r.ID == ? AND r.ID = ur.RoleID AND ur.UserID = u.ID`,
		roleID)
	return users, err
}

// GetRoleByName gets a named role from the database
func (d *Db) GetRoleByName(name string, communityID int64) (Role, error) {
	var r Role
	var err error
	if communityID == 0 {
		err = d.db.SelectOne(&r,
			"SELECT ID, Name, Description FROM roles WHERE Name=? AND CommunityID=?",
			name, communityID)
	} else {
		err = d.db.SelectOne(&r,
			`SELECT r.ID ID, r.Name Name, r.Description Description, c.ID CommunityID, c.Name CommunityName
			FROM roles r, communities c
			WHERE r.CommunityID = c.ID AND r.Name=? AND r.CommunityID=?`,
			name, communityID)
	}
	return r, err
}

// AddRole adds a role to the database, but only if it's not already present
// For roles that have no associated community use communityID: 0
func (d *Db) AddRole(name string, communityID int64, description string) (Role, error) {
	if name == SuperAdminRoleName && communityID != 0 {
		return Role{}, def.Err(nil, "The SuperAdmin role name can only be used with communityID 0")
	}
	role, err := d.GetRoleByName(name, communityID)
	if err == nil {
		return role, nil
	}
	if !isNoResultsError(err) {
		return role, def.Err(err, "error in GetRoleByName:"+name)
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

	role, err = d.GetRoleByName(name, communityID)
	if err != nil {
		return Role{}, def.Err(err, "error in GetRoleByName:"+name)
	}
	return role, nil
}

// AddRoleToUser
func (d *Db) AddRoleToUser(userID, roleID int64) error {
	var urs []UserRoleTable
	_, err := d.db.Select(&urs, `SELECT * FROM userroles WHERE UserID=? AND RoleID=?`, userID, roleID)
	if len(urs) > 0 {
		return nil
	}
	ur := UserRoleTable{
		UserID: userID,
		RoleID: roleID,
	}
	err = d.db.Insert(&ur)
	return err
}

// DeleteRoleFromUser
func (d *Db) DeleteRoleFromUser(userID, roleID int64) error {
	_, err := d.db.Exec("DELETE FROM userroles WHERE UserID=? AND RoleID=?", userID, roleID)
	return err
}

// GetUserRoles
func (d *Db) GetUserRoles(userID int64) ([]Role, error) {
	var roles []Role
	_, err := d.db.Select(&roles,
		`SELECT r.ID ID, r.Name Name, r.Description Description
		FROM userroles ur, roles r
		WHERE ur.UserID = ? AND ur.RoleID = r.ID AND r.CommunityID = ?`,
		userID, 0)
	if err != nil {
		return nil, err
	}

	var communityRoles []Role
	_, err = d.db.Select(
		&communityRoles,
		`SELECT r.ID ID, r.Name Name, r.Description Description, c.ID CommunityID, c.Name CommunityName
		FROM userroles ur, roles r, communities c
		WHERE ur.UserID = ? AND ur.RoleID = r.ID AND r.CommunityID = c.ID`,
		userID)
	if err != nil {
		return roles, err
	}
	roles = append(roles, communityRoles...)
	return roles, nil
}

// HasSuperAdminRole checks if a a user has super admin privileges
func (d *Db) HasSuperAdminRole(userID int64) bool {
	roles, err := d.GetUserRoles(userID)
	if err != nil {
		log.Println(def.Err(err, "GetUserRoles failed"))
		return false
	}
	for _, r := range roles {
		if r.Name == SuperAdminRoleName && r.CommunityID == 0 {
			return true
		}
	}
	return false
}

// IsServiceOwner checks if a certain user owns a certain service
func (d *Db) IsServiceOwner(userID int64, serviceID ServiceID) bool {
	var x OwnerTable
	err := d.db.SelectOne(&x,
		"SELECT * FROM owners WHERE UserID=? AND ObjectType=? AND ObjectID=?",
		userID, "Service", string(serviceID))
	if err != nil && !isNoResultsError(err) {
		log.Printf("ERROR in IsServiceOwner: %#v", err)
	}
	return err == nil
}

// IsJobOwner checks if a certain user owns a certain job
func (d *Db) IsJobOwner(userID int64, jobID JobID) bool {
	var x OwnerTable
	err := d.db.SelectOne(&x,
		"SELECT * FROM owners WHERE UserID=? AND ObjectType=? AND ObjectID=?",
		userID, "Job", string(jobID))
	if err != nil && !isNoResultsError(err) {
		log.Printf("ERROR in IsJobOwner: %#v", err)
	}
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

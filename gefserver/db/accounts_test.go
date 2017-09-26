package db

import (
	"log"
	"os"
	"testing"
	"time"

	. "github.com/EUDAT-GEF/GEF/gefserver/tests"
)

var (
	name1  = "user1"
	email1 = "user1@example.com"
	name2  = "user2"
	email2 = "user2@example.com"
)

func TestRolesUser(t *testing.T) {
	db, file, err := InitDbForTesting()
	CheckErr(t, err)
	defer db.Close()
	defer os.Remove(file)

	user1 := AddTestUser(t, db, name1, email1)

	superAdminRole, err := db.GetRoleByName(SuperAdminRoleName, 0)
	CheckErr(t, err)

	err = db.AddRoleToUser(user1.ID, superAdminRole.ID)
	CheckErr(t, err)

	users, err := db.GetRoleUsers(superAdminRole.ID)
	CheckErr(t, err)
	log.Println(users)
	ExpectEquals(t, len(users), 1)
	ExpectEquals(t, &users[0], user1)

	roles, err := db.GetUserRoles(user1.ID)
	CheckErr(t, err)
	log.Println(roles)
	ExpectEquals(t, roles[0], superAdminRole)
}

func TestServiceAndOwnership(t *testing.T) {
	db, file, err := InitDbForTesting()
	CheckErr(t, err)
	defer db.Close()
	defer os.Remove(file)

	user1 := AddTestUser(t, db, name1, email1)

	service := Service{
		ID:           ServiceID("service_test_id"),
		ConnectionID: ConnectionID(1),
		Name:         "service name",
		RepoTag:      "service repo tag",
		Description:  "service description",
	}
	err = db.AddService(user1.ID, service)
	CheckErr(t, err)
	Expect(t, db.IsServiceOwner(user1.ID, service.ID))

	err = db.RemoveService(service.ID)
	CheckErr(t, err)
	Expect(t, !db.IsServiceOwner(user1.ID, service.ID))

	_, err = db.GetService(service.ID)
	Expect(t, isNoResultsError(err))
}

func TestUsersAndTokens(t *testing.T) {
	db, file, err := InitDbForTesting()
	CheckErr(t, err)
	defer db.Close()
	defer os.Remove(file)

	user1 := AddTestUser(t, db, name1, email1)
	user2 := AddTestUser(t, db, name2, email2)
	testUserTokens(t, db, user1, user2)
}

func TestCommunityAndUserRoles(t *testing.T) {
	db, file, err := InitDbForTesting()
	CheckErr(t, err)
	defer db.Close()
	defer os.Remove(file)

	c1 := AddTestCommunity(t, db, "community1")
	c2 := AddTestCommunity(t, db, "community2")
	c0, err := db.GetCommunityByName("EUDAT")
	CheckErr(t, err)

	cs, err := db.ListCommunities()
	CheckErr(t, err)
	ExpectEquals(t, cs, []Community{c0, c1, c2})

	r1 := AddTestRole(t, db, c1, "testrole1")
	r2 := AddTestRole(t, db, c1, "testrole2")
	r3 := AddTestRole(t, db, c1, "testrole3")

	r1c2 := AddTestRole(t, db, c2, "testrole1c2")
	r2c2 := AddTestRole(t, db, c2, "testrole2c2")

	var roleset = map[int64]int{
		r1.ID: 0, r2.ID: 0, r3.ID: 0, r1c2.ID: 0, r2c2.ID: 0,
	}
	ExpectEquals(t, len(roleset), 5) // all different IDs

	roles, err := db.ListRoles()
	CheckErr(t, err)
	ExpectEquals(t, roles[len(roles)-5:], []Role{r1, r2, r3, r1c2, r2c2})

	user1 := AddTestUser(t, db, name1, email1)
	user2 := AddTestUser(t, db, name2, email2)

	err = db.AddRoleToUser(user1.ID, r1.ID)
	CheckErr(t, err)
	db.AddRoleToUser(user1.ID, r2.ID)
	CheckErr(t, err)
	db.AddRoleToUser(user1.ID, r2c2.ID)
	CheckErr(t, err)

	db.AddRoleToUser(user2.ID, r2.ID)
	CheckErr(t, err)
	db.AddRoleToUser(user2.ID, r3.ID)
	CheckErr(t, err)
	db.AddRoleToUser(user2.ID, r1c2.ID)
	CheckErr(t, err)

	roles, err = db.GetUserRoles(user1.ID)
	CheckErr(t, err)
	ExpectEquals(t, roles, []Role{r1, r2, r2c2})

	roles, err = db.GetUserRoles(user2.ID)
	CheckErr(t, err)
	ExpectEquals(t, roles, []Role{r2, r3, r1c2})
}

func AddTestRole(t *testing.T, db Db, community Community, name string) Role {
	r, err := db.AddRole(name, community.ID, "description of testrole "+name)
	CheckErr(t, err)
	r2, err := db.AddRole(name, community.ID, "description of testrole "+name)
	CheckErr(t, err)
	ExpectEquals(t, r, r2)

	r2, err = db.GetRoleByName(name, community.ID)
	CheckErr(t, err)
	ExpectEquals(t, r, r2)
	return r
}

func AddTestCommunity(t *testing.T, db Db, name string) Community {
	c, err := db.AddCommunity(name, "Description of test community: "+name, true)
	CheckErr(t, err)
	Expect(t, c.ID != 0)

	com, err := db.GetCommunityByID(c.ID)
	CheckErr(t, err)
	ExpectEquals(t, c, com)

	com, err = db.GetCommunityByName(c.Name)
	CheckErr(t, err)
	ExpectEquals(t, c, com)

	// expect automatically created roles
	_, err = db.GetRoleByName(CommunityAdminRoleName, c.ID)
	CheckErr(t, err)
	_, err = db.GetRoleByName(CommunityMemberRoleName, c.ID)
	CheckErr(t, err)

	return com
}

func AddTestUser(t *testing.T, db Db, name string, email string) *User {
	user, err := db.GetUserByEmail(email)
	CheckErr(t, err)
	ExpectNotNil(t, user)

	now := time.Now()
	user = &User{
		Name:    name,
		Email:   email,
		Created: now,
		Updated: now,
	}

	_user, err := db.AddUser(*user)
	CheckErr(t, err)
	ExpectEquals(t, user.Name, name)
	ExpectEquals(t, user.Email, email)
	user = &_user

	Expect(t, user.ID != 0)
	user, err = db.GetUserByID(user.ID)
	CheckErr(t, err)
	ExpectNotNil(t, user)
	ExpectEquals(t, user.Name, name)
	ExpectEquals(t, user.Email, email)

	user.Name = "xx"
	db.UpdateUser(*user)
	user, err = db.GetUserByEmail(email)
	CheckErr(t, err)
	ExpectEquals(t, user.Name, "xx")

	user.Name = name
	db.UpdateUser(*user)
	user, err = db.GetUserByEmail(email)
	CheckErr(t, err)
	ExpectEquals(t, user.Name, name)

	return user
}

func testUserTokens(t *testing.T, db Db, user1, user2 *User) {
	expire := time.Now().AddDate(10, 0, 0)

	t1, err := db.NewUserToken(user1.ID, "token1", expire)
	CheckErr(t, err)
	t2, err := db.NewUserToken(user1.ID, "token2", expire)
	CheckErr(t, err)
	t3, err := db.NewUserToken(user1.ID, "token3", expire)
	CheckErr(t, err)

	tt1, err := db.NewUserToken(user2.ID, "token 1 user 2", expire)
	CheckErr(t, err)
	tt2, err := db.NewUserToken(user2.ID, "token 2 user 2", expire)
	CheckErr(t, err)

	tokenList, err := db.GetUserTokens(user1.ID)
	CheckErr(t, err)
	ExpectEqualTokens(t, tokenList[0], t1)
	ExpectEqualTokens(t, tokenList[1], t2)
	ExpectEqualTokens(t, tokenList[2], t3)

	token, err := db.GetTokenByID(t1.ID)
	CheckErr(t, err)
	ExpectEqualTokens(t, token, t1)

	token, err = db.GetTokenBySecret(t2.Secret)
	CheckErr(t, err)
	ExpectEqualTokens(t, token, t2)

	err = db.DeleteUserToken(user1.ID, t2.ID)
	CheckErr(t, err)

	tokenList, err = db.GetUserTokens(user1.ID)
	CheckErr(t, err)
	ExpectEqualTokens(t, tokenList[0], t1)
	ExpectEqualTokens(t, tokenList[1], t3)

	tokenList, err = db.GetUserTokens(user2.ID)
	CheckErr(t, err)
	ExpectEqualTokens(t, tokenList[0], tt1)
	ExpectEqualTokens(t, tokenList[1], tt2)
}

func ExpectEqualTokens(t *testing.T, t1, t2 Token) {
	ExpectEquals(t, t1.ID, t2.ID)
	ExpectEquals(t, t1.Name, t2.Name)
	ExpectEquals(t, t1.Secret, t2.Secret)
	ExpectEquals(t, t1.Expire.Unix(), t2.Expire.Unix())
}

package tests

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/EUDAT-GEF/GEF/gefserver/db"
	"github.com/EUDAT-GEF/GEF/gefserver/def"
	"github.com/EUDAT-GEF/GEF/gefserver/pier"
)

const testPID = "11304/a3d012ca-4e23-425e-9e2a-1e6a195b966f"

var configFilePath = "../config.json"
var internalServicesFolder = "../../services/_internal"

var (
	name1  = "user1"
	email1 = "user1@example.com"
	name2  = "user2"
	email2 = "user2@example.com"
)

func TestClient(t *testing.T) {

	config, err := def.ReadConfigFile(configFilePath)
	checkMsg(t, err, "reading config files")

	db, file, err := db.InitDbForTesting()
	checkMsg(t, err, "creating db")
	defer db.Close()
	defer os.Remove(file)

	user1 := addUser(t, db, name1, email1)
	user2 := addUser(t, db, name2, email2)
	testUserTokens(t, db, user1, user2)

	pier, err := pier.NewPier(&db, config.TmpDir)
	checkMsg(t, err, "creating new pier")

	err = pier.SetDockerConnection(config.Docker, config.Limits, config.Timeouts, internalServicesFolder)
	checkMsg(t, err, "setting docker connection")

	before, err := db.ListServices()
	checkMsg(t, err, "listing services failed")

	service, err := pier.BuildService("./clone_test")
	checkMsg(t, err, "build service failed")
	log.Println("test service built:", service)

	after, err := db.ListServices()
	checkMsg(t, err, "listing services failed")

	errstr := "Cannot find new service in list"
	for _, x := range after {
		if x.ID == service.ID {
			errstr = ""
			break
		}
	}

	if errstr != "" {
		t.Error("before is: ", len(before), before)
		t.Error("service is: ", service)
		t.Error("after is: ", len(after), after)
		t.Error("")
		t.Error(errstr)
		t.Fail()
		return
	}

	job, err := pier.RunService(service, testPID)
	for {
		runningJob, err := db.GetJob(job.ID)
		checkMsg(t, err, "running job failed")
		if runningJob.State.Code > -1 {
			expect(t, runningJob.State.Code == 0, "job failed: "+runningJob.State.Error)
			break
		}
	}

	checkMsg(t, err, "running service failed")
	log.Println("test job: ", job)

	jobList, err := db.ListJobs()
	if len(jobList) == 0 {
		t.Error("cannot find any job")
		t.FailNow()
	}

	found := false
	for _, j := range jobList {
		if j.ID == job.ID {
			found = true
		}
	}
	if !found {
		t.Error("cannot find the new job in the list of all jobs")
		t.FailNow()
	}

	j, err := db.GetJob(job.ID)
	checkMsg(t, err, "getting job failed")
	if j.ID != job.ID {
		t.Error("job retrieval based on id failure")
		t.FailNow()
	}
}

func TestExecution(t *testing.T) {
	log.Println("TestB running")
	config, err := def.ReadConfigFile(configFilePath)
	checkMsg(t, err, "reading config files")

	db, dbfile, err := db.InitDbForTesting()
	checkMsg(t, err, "creating db")
	defer db.Close()
	defer os.Remove(dbfile)

	addUser(t, db, name1, email1)
	addUser(t, db, name2, email2)

	pier, err := pier.NewPier(&db, config.TmpDir)
	checkMsg(t, err, "creating new pier")

	err = pier.SetDockerConnection(config.Docker, config.Limits, config.Timeouts, internalServicesFolder)
	checkMsg(t, err, "setting docker connection")

	service, err := pier.BuildService("./clone_test")
	checkMsg(t, err, "build service failed")
	log.Println("test service built:", service)

	job, err := pier.RunService(service, testPID)
	checkMsg(t, err, "running service failed")

	log.Println("test job: ", job)
	jobid := job.ID

	for job.State.Code == -1 {
		job, err = db.GetJob(jobid)
		checkMsg(t, err, "getting job failed")
	}

	if job.State.Error != "" {
		log.Println("test job error:")
		log.Println("state: ", job.State)
		for i, t := range job.Tasks {
			log.Println("task ", i, ":", t)
		}
	}
	expect(t, job.State.Error == "", "job error")

	files, err := pier.ListFiles(job.OutputVolume, "")
	checkMsg(t, err, "getting volume failed")

	expect(t, len(files) == 1, "bad returned files")

	_, err = pier.RemoveJob(jobid)
	checkMsg(t, err, "removing job failed")
}

///////////////////////////////////////////////////////////////////////////////

func addUser(t *testing.T, d db.Db, name string, email string) *db.User {
	user, err := d.GetUserByEmail(email)
	checkMsg(t, err, "getting user by email failed")
	expect(t, user == nil, "user is not nil")

	now := time.Now()
	user = &db.User{
		Name:    name,
		Email:   email,
		Created: now,
		Updated: now,
	}

	_user, err := d.AddUser(*user)
	checkMsg(t, err, "add user failed")
	expect(t, user.Name == name && user.Email == email,
		"add user returns garbage")
	user = &_user

	expect(t, user.ID != 0, "user id is 0")
	user, err = d.GetUserByID(user.ID)
	checkMsg(t, err, "getting user by ID failed")
	expect(t, user != nil, "getting user by ID returned nil but shouldn't")

	user.Name = "xx"
	d.UpdateUser(*user)
	user, err = d.GetUserByEmail(email)
	checkMsg(t, err, "getting user by email failed")
	expect(t, user.Name == "xx", "bad user name")

	user.Name = name
	d.UpdateUser(*user)
	user, err = d.GetUserByEmail(email)
	checkMsg(t, err, "getting user by email failed")
	expect(t, user.Name == name, "bad updated user name")

	return user
}

func testUserTokens(t *testing.T, d db.Db, user1, user2 *db.User) {
	expire := time.Now().AddDate(10, 0, 0)

	t1, err := d.NewUserToken(user1.ID, "token1", expire)
	checkMsg(t, err, "error in NewUserToken")
	t2, err := d.NewUserToken(user1.ID, "token2", expire)
	checkMsg(t, err, "error in NewUserToken")
	t3, err := d.NewUserToken(user1.ID, "token3", expire)
	checkMsg(t, err, "error in NewUserToken")

	tt1, err := d.NewUserToken(user2.ID, "token 1 user 2", expire)
	checkMsg(t, err, "error in NewUserToken")
	tt2, err := d.NewUserToken(user2.ID, "token 2 user 2", expire)
	checkMsg(t, err, "error in NewUserToken")

	tokenList, err := d.GetUserTokens(user1.ID)
	checkMsg(t, err, "error in GetUserTokens")
	expectEqualTokens(t, tokenList[0], t1, "token1 mismatch")
	expectEqualTokens(t, tokenList[1], t2, "token2 mismatch")
	expectEqualTokens(t, tokenList[2], t3, "token3 mismatch")

	token, err := d.GetTokenByID(t1.ID)
	checkMsg(t, err, "error in GetTokenByID")
	expectEqualTokens(t, token, t1, "t2 by id mismatch")

	token, err = d.GetTokenBySecret(t2.Secret)
	checkMsg(t, err, "error in GetTokenBySecret")
	expectEqualTokens(t, token, t2, "t2 by secret mismatch")

	err = d.DeleteUserToken(user1.ID, t2.ID)
	checkMsg(t, err, "error in DeleteUserToken")

	tokenList, err = d.GetUserTokens(user1.ID)
	checkMsg(t, err, "error in GetUserTokens")
	expectEqualTokens(t, tokenList[0], t1, "first token mismatch")
	expectEqualTokens(t, tokenList[1], t3, "second token mismatch")

	tokenList, err = d.GetUserTokens(user2.ID)
	checkMsg(t, err, "error in GetUserTokens")
	expectEqualTokens(t, tokenList[0], tt1, "u2 first token mismatch")
	expectEqualTokens(t, tokenList[1], tt2, "u2 second token mismatch")
}

func expectEqualTokens(t *testing.T, t1, t2 db.Token, msg string) {
	expect(t, t1.ID == t2.ID, msg)
	expect(t, t1.Name == t2.Name, msg)
	expect(t, t1.Secret == t2.Secret, msg)
	expect(t, t1.Expire.Unix() == t2.Expire.Unix(), msg)
}

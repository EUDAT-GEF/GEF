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
	name3  = "user3"
	email3 = "user3@example.com"
)

func TestClient(t *testing.T) {
	config, err := def.ReadConfigFile(configFilePath)
	CheckErr(t, err)

	// overwrite this because when testing we're in a different working directory
	config.Pier.InternalServicesFolder = internalServicesFolder

	db, file, err := db.InitDbForTesting()
	CheckErr(t, err)
	defer db.Close()
	defer os.Remove(file)
	user, _ := AddUserWithToken(t, db, name1, email1)

	pier, err := pier.NewPier(&db, config.Pier, config.TmpDir, config.Timeouts)
	CheckErr(t, err)

	connID, err := pier.AddDockerConnection(0, config.Docker)
	CheckErr(t, err)

	before, err := db.ListServices()
	CheckErr(t, err)

	service, err := pier.BuildService(connID, user.ID, "./clone_test")
	CheckErr(t, err)
	log.Print("test service built: ", service.ID, " ", service.ImageID)
	log.Printf("test service built: %#v", service)

	after, err := db.ListServices()
	CheckErr(t, err)

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

	job, err := pier.RunService(user.ID, service.ID, testPID, config.Limits, config.Timeouts)
	CheckErr(t, err)
	for job.State.Code == -1 {
		job, err = db.GetJob(job.ID)
		CheckErr(t, err)
	}

	log.Print("test job: ", job.ID)
	// log.Printf("test job: %#v", job)

	jobList, err := db.ListJobs()
	Expect(t, len(jobList) != 0)

	found := false
	for _, j := range jobList {
		if j.ID == job.ID {
			found = true
		}
	}
	Expect(t, found)

	j, err := db.GetJob(job.ID)
	CheckErr(t, err)
	ExpectEquals(t, j.ID, job.ID)
}

func TestExecution(t *testing.T) {
	config, err := def.ReadConfigFile(configFilePath)
	CheckErr(t, err)

	// overwrite this because when testing we're in a different working directory
	config.Pier.InternalServicesFolder = internalServicesFolder

	db, dbfile, err := db.InitDbForTesting()
	CheckErr(t, err)
	defer db.Close()
	defer os.Remove(dbfile)
	user, _ := AddUserWithToken(t, db, name1, email1)

	p, err := pier.NewPier(&db, config.Pier, config.TmpDir, config.Timeouts)
	CheckErr(t, err)

	connID, err := p.AddDockerConnection(0, config.Docker)
	CheckErr(t, err)

	service, err := p.BuildService(connID, user.ID, "./clone_test")
	CheckErr(t, err)
	log.Print("test service built: ", service.ID, " ", service.ImageID)

	job, err := p.RunService(user.ID, service.ID, testPID, config.Limits, config.Timeouts)
	CheckErr(t, err)

	log.Print("test job: ", job.ID)
	jobid := job.ID

	for job.State.Code == -1 {
		job, err = db.GetJob(jobid)
		CheckErr(t, err)
	}

	if job.State.Error != "" {
		for i, t := range job.Tasks {
			log.Println("task ", i, ":", t)
		}
	}
	ExpectEquals(t, job.State.Error, "")

	files, err := p.ListFiles(job.OutputVolume[0].VolumeID, "", config.Limits, config.Timeouts)
	CheckErr(t, err)
	ExpectEquals(t, len(files), 1)

	_, err = p.RemoveJob(user.ID, jobid)
	CheckErr(t, err)
}

func TestJobTimeOut(t *testing.T) {
	config, err := def.ReadConfigFile(configFilePath)
	CheckErr(t, err)
	config.Timeouts.JobExecution = 3  // forcing a small time out
	config.Timeouts.CheckInterval = 2 // forcing a small check interval

	// overwrite this because when testing we're in a different working directory
	config.Pier.InternalServicesFolder = internalServicesFolder

	db, dbfile, err := db.InitDbForTesting()
	CheckErr(t, err)
	defer db.Close()
	defer os.Remove(dbfile)
	user, _ := AddUserWithToken(t, db, name1, email1)

	p, err := pier.NewPier(&db, config.Pier, config.TmpDir, config.Timeouts)
	CheckErr(t, err)

	connID, err := p.AddDockerConnection(0, config.Docker)
	CheckErr(t, err)

	service, err := p.BuildService(connID, user.ID, "./timeout_test")
	CheckErr(t, err)
	log.Print("test service built: ", service.ID, " ", service.ImageID)
	// log.Printf("test service built: %#v", service)

	timedOutjob, err := p.RunService(user.ID, service.ID, testPID, config.Limits, config.Timeouts)
	CheckErr(t, err)

	log.Print("test timed out job: ", timedOutjob.ID)
	// log.Printf("test timed out job: %#v", timedOutjob)
	timedOutJobId := timedOutjob.ID

	for timedOutjob.State.Code == -1 {
		timedOutjob, err = db.GetJob(timedOutJobId)
		CheckErr(t, err)
	}

	ExpectEquals(t, timedOutjob.State.Error, pier.JobTimeOutError)
}

func AddUserWithToken(t *testing.T, d db.Db, name, email string) (db.User, db.Token) {
	now := time.Now()
	user, err := d.AddUser(db.User{
		Name:    name,
		Email:   email,
		Created: now,
		Updated: now,
	})
	CheckErr(t, err)

	expire := time.Now().AddDate(0, 0, 1)
	token, err := d.NewUserToken(user.ID, "TestToken", expire)
	CheckErr(t, err)
	Expect(t, !d.HasSuperAdminRole(user.ID))

	return user, token
}

func MakeMember(t *testing.T, d db.Db, communityName string, userID int64) {
	c, err := d.GetCommunityByName(communityName)
	CheckErr(t, err)
	member, err := d.GetRoleByName(db.CommunityMemberRoleName, c.ID)
	d.AddRoleToUser(userID, member.ID)
}

func MakeAdmin(t *testing.T, d db.Db, communityName string, userID int64) {
	c, err := d.GetCommunityByName(communityName)
	CheckErr(t, err)
	admin, err := d.GetRoleByName(db.CommunityAdminRoleName, c.ID)
	d.AddRoleToUser(userID, admin.ID)
}

func SetSuperAdmin(t *testing.T, d db.Db, userID int64) {
	superadminrole, err := d.GetRoleByName(db.SuperAdminRoleName, 0)
	CheckErr(t, err)
	err = d.AddRoleToUser(userID, superadminrole.ID)
	CheckErr(t, err)
	Expect(t, d.HasSuperAdminRole(userID))
}

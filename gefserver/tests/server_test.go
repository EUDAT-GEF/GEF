package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/EUDAT-GEF/GEF/gefserver/db"
	"github.com/EUDAT-GEF/GEF/gefserver/def"
	"github.com/EUDAT-GEF/GEF/gefserver/pier"
	"github.com/EUDAT-GEF/GEF/gefserver/server"

	"strings"
	"fmt"
)



// JobVolume points to volumes bound to a particular job
type JobVolume struct {
	VolumeID string
	Name string
}

func TestServer(t *testing.T) {
	config, err := def.ReadConfigFile(configFilePath)
	CheckErr(t, err)

	// overwrite this because when testing we're in a different working directory
	config.Pier.InternalServicesFolder = internalServicesFolder

	db, dbfile, err := db.InitDbForTesting()
	CheckErr(t, err)
	defer db.Close()
	defer os.Remove(dbfile)

	superuser, token := AddUserWithToken(t, db, "admin", "admin@example.com")
	SetSuperAdmin(t, db, superuser.ID)
	superToken := token.Secret

	_, token = AddUserWithToken(t, db, name1, email1)
	userToken := token.Secret

	member, token := AddUserWithToken(t, db, name2, email2)
	MakeMember(t, db, "EUDAT", member.ID)
	memberToken := token.Secret

	admin, token := AddUserWithToken(t, db, name3, email3)
	MakeAdmin(t, db, "EUDAT", admin.ID)
	adminToken := token.Secret

	pier, err := pier.NewPier(&db, config.Pier, config.TmpDir, config.Timeouts)
	CheckErr(t, err)

	// connID, err := pier.AddDockerConnection(0, config.Docker)
	_, err = pier.AddDockerConnection(0, config.Docker)
	CheckErr(t, err)

	var srv *httptest.Server
	{
		s, err := server.NewServer(config, pier, &db)
		CheckErr(t, err)
		srv = httptest.NewServer(s.Server.Handler)
	}
	defer srv.Close()

	baseURL := srv.URL + "/api/"

	// test get info
	info, code := getRes(t, gefurl(baseURL, ""))
	ExpectEquals(t, code, 200)
	ExpectEquals(t, info["Version"], server.Version)

	// test get anonymous user info
	userData, code := getRes(t, gefurl(baseURL+"user", ""))
	ExpectEquals(t, code, 200)
	ExpectEquals(t, userData, make(map[string]interface{}))

	// test get superadmin info
	userData, code = getRes(t, gefurl(baseURL+"user", superToken))
	ExpectEquals(t, code, 200)
	user := userData["User"].(map[string]interface{})
	ExpectEquals(t, user["Name"], "admin")
	ExpectEquals(t, userData["IsSuperAdmin"], true)

	// test get simple user info
	userData, code = getRes(t, gefurl(baseURL+"user", userToken))
	ExpectEquals(t, code, 200)
	user = userData["User"].(map[string]interface{})
	ExpectEquals(t, user["Name"], name1)
	ExpectEquals(t, user["Email"], email1)
	ExpectEquals(t, userData["IsSuperAdmin"], false)

	// test get member info
	userData, code = getRes(t, gefurl(baseURL+"user", memberToken))
	ExpectEquals(t, code, 200)
	user = userData["User"].(map[string]interface{})
	ExpectEquals(t, user["Name"], name2)
	ExpectEquals(t, user["Email"], email2)
	ExpectEquals(t, userData["IsSuperAdmin"], false)
	userRoleList := userData["Roles"].([]interface{})
	userRole := userRoleList[0].(map[string]interface{})
	ExpectEquals(t, userRole["Name"], "Member")
	ExpectEquals(t, userRole["CommunityID"], 1.0) // json represents numbers as floats

	// test get admin info
	userData, code = getRes(t, gefurl(baseURL+"user", adminToken))
	ExpectEquals(t, code, 200)
	user = userData["User"].(map[string]interface{})
	ExpectEquals(t, user["Name"], name3)
	ExpectEquals(t, user["Email"], email3)
	ExpectEquals(t, userData["IsSuperAdmin"], false)
	userRoleList = userData["Roles"].([]interface{})
	userRole = userRoleList[0].(map[string]interface{})
	ExpectEquals(t, userRole["Name"], "Administrator")
	ExpectEquals(t, userRole["CommunityID"], 1.0) // json represents numbers as floats

	// test create a build
	build, code := postRes(t, gefurl(baseURL+"builds", ""), nil)
	ExpectEquals(t, code, 401)
	build, code = postRes(t, gefurl(baseURL+"builds", userToken), nil)
	ExpectEquals(t, code, 403)
	build, code = postRes(t, gefurl(baseURL+"builds", memberToken), nil)
	ExpectEquals(t, code, 403)
	build, code = postRes(t, gefurl(baseURL+"builds", superToken), nil)
	ExpectEquals(t, code, 201)
	build, code = postRes(t, gefurl(baseURL+"builds", adminToken), nil)
	ExpectEquals(t, code, 201)

	// test make a service
	buildURL := srv.URL + build["Location"].(string)
	res, code := uploadDir(t, gefurl(buildURL, ""), "./clone_test")
	ExpectEquals(t, code, 401)
	res, code = uploadDir(t, gefurl(buildURL, userToken), "./clone_test")
	ExpectEquals(t, code, 403)
	res, code = uploadDir(t, gefurl(buildURL, memberToken), "./clone_test")
	ExpectEquals(t, code, 403)
	res, code = uploadDir(t, gefurl(buildURL, adminToken), "./clone_test")
	ExpectEquals(t, code, 200)
	service := res["Service"].(map[string]interface{})

	ExpectEquals(t, service["Name"], "Test Clone")
	serviceID := service["ID"].(string)

	res, code = getRes(t, gefurl(baseURL+"services", ""))
	ExpectEquals(t, code, 200)
	services := res["Services"].([]interface{})
	ExpectNotNil(t, services)
	Expect(t, len(services) > 0)

	// test create a job
	jobParams := map[string]string{
		"serviceID": serviceID,
		"pid":       testPID,
	}
	jobLink, code := postRes(t, gefurl(baseURL+"jobs", ""), jobParams)
	ExpectEquals(t, code, 401)
	jobLink, code = postRes(t, gefurl(baseURL+"jobs", userToken), jobParams)
	ExpectEquals(t, code, 403)
	jobLink, code = postRes(t, gefurl(baseURL+"jobs", adminToken), jobParams)
	ExpectEquals(t, code, 201)
	jobLink, code = postRes(t, gefurl(baseURL+"jobs", memberToken), jobParams)
	ExpectEquals(t, code, 201)

	jobURL := srv.URL + jobLink["Location"].(string)
	res, code = getRes(t, gefurl(jobURL, ""))
	ExpectEquals(t, code, 200)
	job := res["Job"].(map[string]interface{})

	ExpectEquals(t, job["ServiceID"], serviceID)
	jobID := job["ID"].(string)

	// test get the list of all jobs
	res, code = getRes(t, gefurl(baseURL+"jobs", ""))
	ExpectEquals(t, code, 200)
	jobs := res["Jobs"].([]interface{})
	ExpectNotNil(t, jobs)
	Expect(t, len(jobs) > 0)

	// loop until the job ends
	exitCode := -1
	for exitCode < 0 {
		time.Sleep(1 * time.Second)
		res, code = getRes(t, gefurl(baseURL+"jobs/"+jobID, ""))
		ExpectEquals(t, code, 200)
		job = res["Job"].(map[string]interface{})
		ExpectNotNil(t, job)
		exitCode = int(job["State"].(map[string]interface{})["Code"].(float64))
	}

	firstOutputVolume := job["OutputVolume"].([]interface{})[0].(map[string]interface{})["VolumeID"].(string)

	// test the job console output
	var console interface{}
	tasks := job["Tasks"].([]interface{})
	if len(tasks) > 1 {
		console = tasks[1].(map[string]interface{})["ConsoleOutput"]
	}

	Expect(t, console == "'/test.txt' -> '/mydata/output/test.txt'\n" ||
		strings.HasSuffix(console.(string),
			"/logs is an experimental feature introduced in Docker 1.13. Unfortunately, it is not yet supported by the Docker client we use"))

	// test get the job file system output
	volURL := baseURL + "volumes/" + firstOutputVolume + "/"

	fmt.Println(volURL)
	res, code = getRes(t, gefurl(volURL, userToken))
	ExpectEquals(t, code, 403)
	res, code = getRes(t, gefurl(volURL, memberToken))
	ExpectEquals(t, code, 200)
	files := res["volumeContent"].([]interface{})

	ExpectEquals(t, len(files), 1)
	file := files[0].(map[string]interface{})

	ExpectEquals(t, file["name"].(string), "test.txt")
	ExpectEquals(t, file["path"].(string), "")

	fileURL := volURL + "text.txt"
	res, code = getRes(t, gefurlFileContent(fileURL, userToken))
	ExpectEquals(t, code, 403)

	// TODO: the following test does not work (server error)
	// content, code := getResString(t, gefurlFileContent(fileURL, memberToken))
	// ExpectEquals(t, code, 200)
	// ExpectEquals(t, content, "Hi there")

	// test service removal
	servicesURL := baseURL + "services/"

	res, code = deleteRes(t, gefurl(servicesURL+serviceID, userToken))
	ExpectEquals(t, code, 403)

	res, code = deleteRes(t, gefurl(servicesURL+serviceID, superToken))
	ExpectEquals(t, code, 200)
}

///////////////////////////////////////////////////////////////////////////////

func deleteRes(t *testing.T, url string) (map[string]interface{}, int) {
	body, code := getBody(t, url, "DELETE")
	defer body.Close()
	if code >= 400 {
		return nil, code
	}
	return readJSON(t, body), code
}

func getRes(t *testing.T, url string) (map[string]interface{}, int) {
	body, code := getBody(t, url, "GET")
	defer body.Close()
	if code >= 400 {
		return nil, code
	}
	return readJSON(t, body), code
}

func getResString(t *testing.T, url string) (string, int) {
	body, code := getBody(t, url, "GET")
	defer body.Close()
	if code >= 400 {
		return "", code
	}
	return readString(t, body), code
}

func postRes(t *testing.T, url string, params map[string]string) (map[string]interface{}, int) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	if params != nil {
		for k, v := range params {
			err := writer.WriteField(k, v)
			CheckErr(t, err)
		}
	}
	err := writer.Close()
	CheckErr(t, err)

	req, err := http.NewRequest("POST", url, body)
	CheckErr(t, err)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := http.DefaultClient.Do(req)
	CheckErr(t, err)
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return nil, res.StatusCode
	}
	return readJSON(t, res.Body), res.StatusCode
}

func uploadDir(t *testing.T, uri string, path string) (map[string]interface{}, int) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	fileinfos, err := ioutil.ReadDir(path)
	CheckErr(t, err)
	ExpectEquals(t, len(fileinfos), 2)
	for _, fi := range fileinfos {
		path := filepath.Join(path, fi.Name())
		file, err := os.Open(path)
		CheckErr(t, err)
		defer file.Close()

		part, err := writer.CreateFormFile(fi.Name(), filepath.Base(path))
		CheckErr(t, err)
		_, err = io.Copy(part, file)
		CheckErr(t, err)
	}
	err = writer.Close()
	CheckErr(t, err)

	req, err := http.NewRequest("POST", uri, body)
	CheckErr(t, err)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := http.DefaultClient.Do(req)
	CheckErr(t, err)
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return nil, res.StatusCode
	}
	return readJSON(t, res.Body), res.StatusCode
}

func getBody(t *testing.T, url string, method string) (io.ReadCloser, int) {
	request, err := http.NewRequest(method, url, nil)
	CheckErr(t, err)
	res, err := http.DefaultClient.Do(request)
	CheckErr(t, err)
	return res.Body, res.StatusCode
}

func gefurl(url, accessToken string) string {
	if accessToken != "" {
		url += "?access_token=" + accessToken
	}
	return url
}

func gefurlFileContent(url, accessToken string) string {
	if accessToken != "" {
		url += "?access_token=" + accessToken + "&content=1"
	} else {
		url += "?content=1"
	}
	return url
}

func readJSON(t *testing.T, body io.Reader) map[string]interface{} {
	data, err := ioutil.ReadAll(body)
	CheckErr(t, err)

	var j map[string]interface{}
	err = json.Unmarshal([]byte(data), &j)
	CheckErr(t, err)
	return j
}

func readString(t *testing.T, body io.Reader) string {
	data, err := ioutil.ReadAll(body)
	CheckErr(t, err)
	return string(data)
}

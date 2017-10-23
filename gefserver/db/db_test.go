package db

import (
	"os"
	"testing"

	"github.com/EUDAT-GEF/GEF/gefserver/def"
	. "github.com/EUDAT-GEF/GEF/gefserver/tests"
)

func TestConnection(t *testing.T) {
	db, file, err := InitDbForTesting()
	CheckErr(t, err)
	defer db.Close()
	defer os.Remove(file)

	connection1 := def.DockerConfig{
		Endpoint:    "http://first.example.com",
		Description: "description",
		TLSVerify:   false,
		CertPath:    "",
		KeyPath:     "",
		CAPath:      "",
	}

	// first connection
	connID1, err := db.AddConnection(0, connection1)
	CheckErr(t, err)
	ExpectNotEquals(t, connID1, 0)

	id, err := db.GetFirstConnectionID()
	CheckErr(t, err)
	ExpectEquals(t, connID1, id)

	cmap, err := db.GetConnections()
	CheckErr(t, err)
	ExpectEquals(t, len(cmap), 1)
	ExpectNotNil(t, cmap[connID1])
	ExpectEquals(t, cmap[connID1], connection1)

	// second connection
	connection2 := connection1
	connection2.Endpoint = "http://second.example.com"

	connID2, err := db.AddConnection(0, connection2)
	CheckErr(t, err)
	ExpectNotEquals(t, connID2, 0)
	ExpectNotEquals(t, connID1, connID2)

	id, err = db.GetFirstConnectionID()
	CheckErr(t, err)
	ExpectEquals(t, connID1, id)

	cmap, err = db.GetConnections()
	CheckErr(t, err)
	ExpectEquals(t, len(cmap), 2)
	ExpectNotNil(t, cmap[connID1])
	ExpectEquals(t, cmap[connID1], connection1)
	ExpectNotNil(t, cmap[connID2])
	ExpectEquals(t, cmap[connID2], connection2)

	// change something in the first connection
	connection1.Description = "different description"
	connID1_1, err := db.AddConnection(0, connection1)
	CheckErr(t, err)
	ExpectEquals(t, connID1, connID1_1)

	id, err = db.GetFirstConnectionID()
	CheckErr(t, err)
	ExpectEquals(t, connID1, id)

	cmap, err = db.GetConnections()
	CheckErr(t, err)
	ExpectEquals(t, len(cmap), 2)
	ExpectNotNil(t, cmap[connID1])
	ExpectEquals(t, cmap[connID1], connection1)
	ExpectNotNil(t, cmap[connID2])
	ExpectEquals(t, cmap[connID2], connection2)
}

package main

import (
	// "fmt"
	"io"
	"log"
	"os"

	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	cpac "github.com/septianw/jas-contact"
	"github.com/septianw/jas/common"
	"github.com/septianw/jas/types"
	"github.com/stretchr/testify/assert"

	"net/http"
	"net/http/httptest"
)

type header map[string]string
type headers []header
type payload struct {
	Method string
	Url    string
	Body   io.Reader
}
type expectation struct {
	Code int
	Body string
}
type quest struct {
	pload  payload
	heads  headers
	expect expectation
}
type quests []quest

var LastPostID int64

func getArm() (*gin.Engine, *httptest.ResponseRecorder) {
	router := gin.New()
	gin.SetMode(gin.ReleaseMode)
	Router(router)

	recorder := httptest.NewRecorder()
	return router, recorder
}

func handleErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func doTheTest(load payload, heads headers) *httptest.ResponseRecorder {
	var router, recorder = getArm()

	req, err := http.NewRequest(load.Method, load.Url, load.Body)
	handleErr(err)

	if len(heads) != 0 {
		for _, head := range heads {
			for key, value := range head {
				req.Header.Set(key, value)
			}
		}
	}
	router.ServeHTTP(recorder, req)

	return recorder
}

func SetupRouter() *gin.Engine {
	return gin.New()
}

func SetEnvironment() {
	var rt types.Runtime
	var Dbconf types.Dbconf

	Dbconf.Database = "ipoint"
	Dbconf.Host = "localhost"
	Dbconf.Pass = "dummypass"
	Dbconf.Port = 3306
	Dbconf.Type = "mysql"
	Dbconf.User = "asep"

	rt.Dbconf = Dbconf
	rt.Libloc = "/home/asep/gocode/src/github.com/septianw/jas/libs"

	common.WriteRuntime(rt)
}

func UnsetEnvironment() {
	os.Remove("/tmp/shinyRuntimeFile")
}

func TestContactPostPositive(t *testing.T) {
	SetEnvironment()
	defer UnsetEnvironment()
	contactIn := cpac.ContactIn{
		"Pramitha",
		"Utami",
		"Mrs",
		"karyawan",
	}
	nc, err := json.Marshal(contactIn)
	common.ErrHandler(err)
	NewContact := strings.NewReader(string(nc))

	q := quest{
		payload{"POST", "/api/v1/contact/", NewContact},
		headers{},
		expectation{201, "contact post"},
	}

	rec := doTheTest(q.pload, q.heads)

	ci, err := cpac.FindContact(contactIn)
	if err != nil || len(ci) == 0 {
		t.Fail()
	}
	t.Logf("\n%+v\n", ci)
	LastPostID = ci[0].Id
	cjson, err := json.Marshal(ci[0])
	if err != nil {
		t.Fail()
	}

	assert.Equal(t, q.expect.Code, rec.Code)
	assert.Equal(t, string(cjson)+"\n", rec.Body.String())
}

func TestContactAllPositive(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	SetEnvironment()
	defer UnsetEnvironment()

	contacts := cpac.GetContact(-1, 2, 0)
	if len(contacts) < 2 {
		t.Fail()
	}
	contactsAllJSON, err := json.Marshal(contacts)
	common.ErrHandler(err)

	contacts = cpac.GetContact(LastPostID, 0, 0)
	if len(contacts) != 1 {
		t.Fail()
	}
	contact := contacts[0]
	contactIdJSON, err := json.Marshal(contact)
	common.ErrHandler(err)

	// log.Println(string(contactIdJSON))
	// log.Println(string(contactUpdatedJSON))

	qs := quests{
		quest{
			payload{"GET", "/api/v1/contact/all/2/0", nil},
			headers{},
			expectation{200, string(contactsAllJSON) + "\n"},
		},
		quest{
			payload{"GET", fmt.Sprintf("/api/v1/contact/%d", LastPostID), nil},
			headers{},
			expectation{200, string(contactIdJSON) + "\n"},
		},
	}

	for _, q := range qs {
		rec := doTheTest(q.pload, q.heads)
		assert.Equal(t, q.expect.Code, rec.Code)
		assert.Equal(t, q.expect.Body, rec.Body.String())
	}
}

func TestContactPutPositive(t *testing.T) {
	SetEnvironment()
	defer UnsetEnvironment()

	UpdateContact := strings.NewReader(`{
	"firstname": "Pramitha",
	"lastname": "Utami",
	"prefix": "Mr",
	"type": "konsumen"
	}`)
	contactUpdatedJSON, err := json.Marshal(cpac.ContactOut{
		LastPostID,
		"Pramitha",
		"Utami",
		"Mr",
		"konsumen",
	})
	common.ErrHandler(err)

	q := quest{
		payload{"PUT", fmt.Sprintf("/api/v1/contact/%d", LastPostID), UpdateContact},
		headers{},
		expectation{200, string(contactUpdatedJSON) + "\n"},
	}

	rec := doTheTest(q.pload, q.heads)
	assert.Equal(t, q.expect.Code, rec.Code)
	assert.Equal(t, q.expect.Body, rec.Body.String())
}

func TestContactDeletePositive(t *testing.T) {
	SetEnvironment()
	defer UnsetEnvironment()

	contactUpdatedJSON, err := json.Marshal(cpac.ContactOut{
		LastPostID,
		"Pramitha",
		"Utami",
		"Mr",
		"konsumen",
	})
	common.ErrHandler(err)

	q := quest{
		payload{"DELETE", fmt.Sprintf("/api/v1/contact/%d", LastPostID), nil},
		headers{},
		expectation{200, string(contactUpdatedJSON)},
	}

	rec := doTheTest(q.pload, q.heads)

	assert.Equal(t, q.expect.Code, rec.Code)
	assert.Equal(t, q.expect.Body+"\n", rec.Body.String())
}

func TestContactPostNegative(t *testing.T) {
	SetEnvironment()
	defer UnsetEnvironment()

	q := quest{
		payload{"POST", "/api/v1/contact/", strings.NewReader(`{ "firstname: "satu", "dua": tiga }`)},
		headers{},
		expectation{400, `{"code":2110,"message":"INPUT_VALIDATION_FAIL: invalid character 's' after object key"}`},
	}

	rec := doTheTest(q.pload, q.heads)

	assert.Equal(t, q.expect.Code, rec.Code)
	assert.Equal(t, q.expect.Body+"\n", rec.Body.String())
}

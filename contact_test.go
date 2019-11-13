package contact

import (
	"testing"

	"log"
	"os"
	"reflect"

	"github.com/septianw/jas/common"
	"github.com/septianw/jas/types"
)

var contactid int64

func SetEnvironment() {
	var rt types.Runtime
	var Dbconf types.Dbconf

	Dbconf.Database = "jasdev"
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

func TestInsertContact(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var ct ContactIn

	SetEnvironment()
	defer UnsetEnvironment()

	ct.Firstname = "firstname test"
	ct.Lastname = "lastname test"
	ct.Prefix = "Jr."
	ct.Type = "mitra"
	ct.Phone = "089455233444"
	ct.Email = "mama@mia.me"

	id, err := InsertContact(ct)

	contacts := GetContact(id, 0, 0)

	if len(contacts) == 0 {
		t.Fail()
	}

	if err != nil {
		t.Fail()
	}

	if id == 0 {
		t.Fail()
	}

	contactid = id

	t.Logf("Last insert id: %d", id)
	t.Logf("Error %+v", err)
	t.Logf("contacts: %+v", contacts)
}

func TestGetContact(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	SetEnvironment()
	defer UnsetEnvironment()

	contacts := GetContact(contactid, 0, 0)

	if len(contacts) != 1 {
		t.Fail()
	}

	contacts = GetContact(-1, 5, 0)

	if len(contacts) != 5 {
		t.Fail()
	}
}

func TestUpdateContact(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var ct, cto ContactIn

	SetEnvironment()
	defer UnsetEnvironment()

	ct.Firstname = "Lyda"
	ct.Lastname = "Media"
	ct.Prefix = "V"
	ct.Type = "mitra"
	rec, err := UpdateContact(contactid, ct)
	log.Println(err)
	log.Println(rec)

	log.Printf("%+v", ct)

	contacts := GetContact(contactid, 0, 0)
	log.Printf("%+v", contacts[0])
	if len(contacts) != 1 {
		t.Fail()
	}
	cto.Firstname = contacts[0].Firstname
	cto.Lastname = contacts[0].Lastname
	cto.Prefix = contacts[0].Prefix
	cto.Type = contacts[0].Type

	if !reflect.DeepEqual(ct, cto) {
		t.Fail()
	}
}

func TestDeleteContact(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	SetEnvironment()
	defer UnsetEnvironment()

	contact, err := DeleteContact(contactid)
	if err != nil {
		t.Fail()
	}
	log.Println(contact)
	log.Println(err)
}

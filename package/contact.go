package contact

import (
	"fmt"

	"database/sql"
	"errors"
	"log"
	"path/filepath"
	"strings"

	"github.com/septianw/jas/common"

	"github.com/gin-gonic/gin"
)

const Version = "0.1.1"

func getdbobj() (db *sql.DB, err error) {
	rt := common.ReadRuntime()
	dbs := common.LoadDatabase(filepath.Join(rt.Libloc, "database.so"), rt.Dbconf)
	db, err = dbs.OpenDb(rt.Dbconf)
	return
}

func Query(q string) (*sql.Rows, error) {
	db, err := getdbobj()
	common.ErrHandler(err)
	defer db.Close()

	return db.Query(q)
}

func Exec(q string) (sql.Result, error) {
	db, err := getdbobj()
	common.ErrHandler(err)
	defer db.Close()

	return db.Exec(q)
}

type Contact struct {
	Contactid int64
	Fname     string
	Lname     string
	Prefix    string
}

type Contacttype struct {
	Ctypeid int64
	Name    string
}

type Contactwtype struct {
	Contact_contactid   int64
	Contacttype_ctypeid int64
}

type ContactIn struct {
	Firstname string `json:"firstname" binding:"required"`
	Lastname  string `json:"lastname" binding:"required"`
	Prefix    string `json:"prefix" binding:"required"`
	Type      string `json:"type" binding:"required"`
}

type ContactOut struct {
	Id        int64  `json:"id" binding:"required"`
	Firstname string `json:"firstname" binding:"required"`
	Lastname  string `json:"lastname" binding:"required"`
	Prefix    string `json:"prefix" binding:"required"`
	Type      string `json:"type" binding:"required"`
}

/*
ERROR CODE LEGEND:
error containt 4 digits,
first digit represent error location either module or main app
1 for main app
2 for module

second digit represent error at level app or database
1 for app
2 for database

third digit represent error with input variable or variable manipulation
0 for skipping this error
1 for input validation error
2 for variable manipulation error

fourth digit represent error with logic, this type of error have
increasing error number based on which part of code that error.
0 for skipping this error
1 for unknown logical error
2 for whole operation fail, operation end unexpectedly
*/

const DATABASE_EXEC_FAIL = 2200
const MODULE_OPERATION_FAIL = 2102
const INPUT_VALIDATION_FAIL = 2110

const VERSION = "0.1.0"

var NOT_ACCEPTABLE = gin.H{"code": "NOT_ACCEPTABLE", "message": "You are trying to request something not acceptible here."}
var NOT_FOUND = gin.H{"code": "NOT_FOUND", "message": "You are find something we can't found it here."}

func GetContact(id, limit, offset int64) (records []ContactOut) {
	var record ContactOut
	q := `
		select distinct
			c.contactid id,
			c.fname firstname,
			c.lname lastname,
			c.prefix prefix,
			ct.name type
		from
			contact c
		join contactwtype cwt
			on c.contactid = cwt.contact_contactid
		join contacttype ct
			on ct.ctypeid = cwt.contacttype_ctypeid
		where deleted = 0
		group by
			id,
			firstname,
			lastname %s`
	constr := `limit %d offset %d`

	qid := `
		select distinct
			c.contactid id,
			c.fname firstname,
			c.lname lastname,
			c.prefix prefix,
			ct.name type
		from
			contact c
		join contactwtype cwt
			on c.contactid = cwt.contact_contactid
		join contacttype ct
			on ct.ctypeid = cwt.contacttype_ctypeid
		%s
		group by
			id,
			firstname,
			lastname`
	constrid := "where c.contactid = %d and deleted = 0"

	if id == -1 {
		if limit == 0 {
			constr = fmt.Sprintf(constr, 10, 0)
			q = fmt.Sprintf(q, constr)
		} else {
			constr = fmt.Sprintf(constr, limit, offset)
			q = fmt.Sprintf(q, constr)
		}
		log.Println(q)
		rows, err := Query(q)
		defer rows.Close()
		common.ErrHandler(err)

		for rows.Next() {
			err := rows.Scan(&record.Id, &record.Firstname, &record.Lastname, &record.Prefix, &record.Type)
			common.ErrHandler(err)

			records = append(records, record)
		}
	} else {
		constrid = fmt.Sprintf(constrid, id)
		qid = fmt.Sprintf(qid, constrid)

		log.Println(qid)
		rows, err := Query(qid)
		common.ErrHandler(err)

		for rows.Next() {
			err := rows.Scan(&record.Id, &record.Firstname, &record.Lastname, &record.Prefix, &record.Type)
			common.ErrHandler(err)

			records = append(records, record)
		}
	}

	return
}

func FindContact(contacin ContactIn) (records []ContactOut, err error) {
	var contact Contact
	var sbQContact strings.Builder
	// var contacts []Contact
	sbQContact.WriteString("SELECT contactid FROM contact WHERE ")

	if strings.Compare(contacin.Firstname, "") != 0 {
		_, err = sbQContact.WriteString(fmt.Sprintf("fname = '%s' ", contacin.Firstname))
		common.ErrHandler(err)
	}
	if strings.Compare(contacin.Lastname, "") != 0 {
		_, err = sbQContact.WriteString(fmt.Sprintf("and lname = '%s' ", contacin.Lastname))
		common.ErrHandler(err)
	}
	if strings.Compare(contacin.Prefix, "") != 0 {
		_, err = sbQContact.WriteString(fmt.Sprintf("and prefix = '%s'", contacin.Prefix))
		common.ErrHandler(err)
	}

	q := sbQContact.String()
	log.Printf("\n%+v\n", q)
	rows, err := Query(q)
	// log.Printf("\n%+v err: %+v\n", rows, err)
	if err != nil {
		return
	}
	for rows.Next() {
		rows.Scan(&contact.Contactid)
		// contacts = append(contacts, contact)
	}
	records = GetContact(contact.Contactid, 0, 0)
	return
}

func InsertContact(contactin ContactIn) (id int64, err error) {
	var contactType Contacttype
	var q string

	if strings.Compare(contactin.Type, "") != 0 {
		q = fmt.Sprintf("SELECT * FROM contacttype WHERE name = '%s'", contactin.Type)
		rows, err := Query(q)
		if err != nil {
			return 0, err
		}
		for rows.Next() {
			rows.Scan(&contactType.Ctypeid, &contactType.Name)
		}
		rows.Close()
	}

	// ambil contact type

	q = fmt.Sprintf(
		"INSERT INTO `contact` (`fname`, `lname`, `prefix`, `deleted`) VALUES ('%s','%s','%s','0')",
		contactin.Firstname,
		contactin.Lastname,
		contactin.Prefix,
	)

	result, err := Exec(q)
	if err != nil {
		return 0, err
	}
	contactID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	if contactType == (Contacttype{}) {
		// sambungkan contact dengan contact type
		q = fmt.Sprintf("INSERT INTO `contactwtype` VALUES ('%d','%d')", contactID, contactType.Ctypeid)
		_, err = Exec(q)
		if err != nil {
			return 0, err
		}
	}

	// ambil record tersimpan
	contacts := GetContact(contactID, 0, 0)
	if len(contacts) == 0 {
		return 0, errors.New("Insert contact fail.")
	}

	return contactID, nil
}

func UpdateContact(contactId int64, contactin ContactIn) (id int64, err error) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// FIXME: terlalu banyak query, cek yang lama dengan yang baru aja dulu, baru update yang berubah aja.
	// var contactDiff ContactOut
	var sbContact strings.Builder
	var dbContactwType strings.Builder
	var qCheckType, qSetContact string
	var set bool = false
	var setField []string
	var ctype Contacttype
	var oldCtypeId int64
	// var result sql.Result

	_, err = sbContact.WriteString("update contact set ")
	common.ErrHandler(err)

	_, err = dbContactwType.WriteString(`update contactwtype
			set contacttype_ctypeid = %d
		where
			contact_contactid = %d and
			contacttype_ctypeid = %d`)
	common.ErrHandler(err)

	contacts := GetContact(contactId, 0, 0)
	if len(contacts) == 0 {
		return 0, errors.New("Contact not found.")
	}
	// log.Println(contacts)
	// log.Println(contactin)

	// compare every value, and find the different
	contact := contacts[0]
	if (strings.Compare(contact.Firstname, contactin.Firstname) != 0) &&
		(strings.Compare(contactin.Firstname, "") != 0) {
		set = true
		setField = append(setField, fmt.Sprintf("fname = '%s'", contactin.Firstname))
	}
	if (strings.Compare(contact.Lastname, contactin.Lastname) != 0) &&
		(strings.Compare(contactin.Lastname, "") != 0) {
		set = true
		setField = append(setField, fmt.Sprintf("lname = '%s'", contactin.Lastname))
	}
	if (strings.Compare(contact.Prefix, contactin.Prefix) != 0) &&
		(strings.Compare(contactin.Prefix, "") != 0) {
		set = true
		setField = append(setField, fmt.Sprintf("prefix = '%s'", contactin.Prefix))
	}
	if (strings.Compare(contact.Type, contactin.Type) != 0) &&
		(strings.Compare(contactin.Type, "") != 0) {
		qCheckType = fmt.Sprintf("select * from contacttype where name = '%s'", contactin.Type)
	}

	if set {
		// sbContact.WriteString(" ")
		sbContact.WriteString(strings.Join(setField, ", "))
		sbContact.WriteString(fmt.Sprintf(" where contactid = %d", contactId))
		qSetContact = sbContact.String()
	}

	// kalau yang berubah hanya di table contact saja
	if strings.Compare(qSetContact, "") != 0 {
		log.Println(qSetContact)
		_, err = Exec(qSetContact)
		if err != nil {
			return 0, err
		}
	}

	// kalau yang berubah hanya type contact saja.
	if strings.Compare(qCheckType, "") != 0 {
		log.Println(qCheckType)
		rows, err := Query(qCheckType)
		if err != nil {
			return 0, err
		}

		for rows.Next() {
			rows.Scan(&ctype.Ctypeid, &ctype.Name)
		}

		if ctype == (Contacttype{}) {
			return 0, errors.New("Contact type not found.")
		}

		q := fmt.Sprintf(`select ctypeid from contacttype where name = '%s'`, contact.Type)
		rows.Close()
		rows, err = Query(q)
		for rows.Next() {
			rows.Scan(&oldCtypeId)
		}

		q = fmt.Sprintf(dbContactwType.String(), ctype.Ctypeid, contactId, oldCtypeId)
		_, err = Exec(q)
		if err != nil {
			return 0, err
		}
	}

	return
}

// ini jenisnya upsert
func SetContact(contactId int64, contactin ContactIn) (id int64, err error) {
	// var input cpac.ContactIn
	// var contactType Contacttype
	// var oldcontacttype, newcontacttype Contacttype

	ctk := GetContact(contactId, 0, 0)

	if len(ctk) == 0 { // kalau kontak tidak ditemukan
		return InsertContact(contactin)
	} else {
		// kalau kontak ditemukan berarti update.
		return UpdateContact(contactId, contactin)
	}
}

func DeleteContact(contactID int64) (contact ContactOut, err error) {
	contacts := GetContact(contactID, 0, 0)

	if len(contacts) != 1 {
		return ContactOut{}, errors.New("Contact not found.")
	}

	contact = contacts[0]

	q := fmt.Sprintf(`UPDATE contact SET deleted=1 WHERE contactid = %d`, contactID)
	result, err := Exec(q)
	if err != nil {
		return
	}

	aff, err := result.RowsAffected()
	if err != nil {
		return
	}
	log.Printf("Record %+v deleted, row affected %d", contacts[0], aff)

	return
}

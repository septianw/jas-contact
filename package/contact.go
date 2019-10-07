package contact

import (
	"fmt"

	"database/sql"
	"log"
	"path/filepath"

	"github.com/septianw/jas/common"

	"github.com/gin-gonic/gin"
)

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
	// var contacts []Contact
	q := fmt.Sprintf(`SELECT contactid
		FROM contact
		WHERE
			fname = '%s' and
			lname = '%s' and
			prefix = '%s'`,
		contacin.Firstname,
		contacin.Lastname,
		contacin.Prefix,
	)
	log.Printf("\n%+v\n", q)
	rows, err := Query(q)
	log.Printf("\n%+v err: %+v\n", rows, err)
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

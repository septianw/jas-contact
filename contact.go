package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"strings"

	"fmt"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/septianw/jas/common"
)

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

func Bootstrap() {
	log.Println("Contact module bootstrap.")
}

func Router(r *gin.Engine) {
	r.Any("/api/v1/contact/*path1", deflt)
}

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

// func GetContactType

func deflt(c *gin.Context) {
	segments := strings.Split(c.Param("path1"), "/")
	// log.Printf("\n%+v\n", c.Request.Method)
	// log.Printf("\n%+v\n", c.Param("path1"))
	// log.Printf("\n%+v\n", segments)
	// log.Printf("\n%+v\n", len(segments))
	switch c.Request.Method {
	case "POST":
		if strings.Compare(segments[1], "") == 0 {
			PostContactHandler(c)
		} else {
			c.AbortWithStatusJSON(http.StatusMethodNotAllowed, NOT_ACCEPTABLE)
		}
		break
	case "GET":
		if strings.Compare(segments[1], "all") == 0 {
			GetContactAllHandler(c)
		} else if i, e := strconv.Atoi(segments[1]); (e == nil) && (i > 0) {
			GetContactIdHandler(c)
		} else {
			c.AbortWithStatusJSON(http.StatusNotAcceptable, NOT_ACCEPTABLE)
		}
		break
	case "PUT":
		if i, e := strconv.Atoi(segments[1]); (e == nil) && (i > 0) {
			PutContactIdHandler(c)
		} else {
			c.AbortWithStatusJSON(http.StatusMethodNotAllowed, NOT_ACCEPTABLE)
		}
		break
	case "DELETE":
		if i, e := strconv.Atoi(segments[1]); (e == nil) && (i > 0) {
			DeleteContactIdHandler(c)
		} else {
			c.AbortWithStatusJSON(http.StatusMethodNotAllowed, NOT_ACCEPTABLE)
		}
		break
	default:
		c.AbortWithStatusJSON(http.StatusMethodNotAllowed, NOT_ACCEPTABLE)
		break
	}
	// c.String(http.StatusOK, "hai")
}

func PostContactHandler(c *gin.Context) {
	var input ContactIn
	var contactType Contacttype

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": INPUT_VALIDATION_FAIL,
			"message": fmt.Sprintf("INPUT_VALIDATION_FAIL: %s", err.Error())})
		return
	}

	// insert ke table contact
	q := fmt.Sprintf(
		"INSERT INTO `contact` (`fname`, `lname`, `prefix`, `deleted`) VALUES ('%s','%s','%s','0')",
		input.Firstname,
		input.Lastname,
		input.Prefix,
	)
	result, err := Exec(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": DATABASE_EXEC_FAIL,
			"message": fmt.Sprintf("DATABASE_EXEC_FAIL: %s", err.Error())})
		return
	}
	contactID, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": DATABASE_EXEC_FAIL,
			"message": fmt.Sprintf("DATABASE_EXEC_FAIL: %s", err.Error())})
		return
	}

	// ambil contact type
	q = fmt.Sprintf("SELECT * FROM contacttype WHERE name = '%s'", input.Type)
	rows, err := Query(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": DATABASE_EXEC_FAIL,
			"message": fmt.Sprintf("DATABASE_EXEC_FAIL: %s", err.Error())})
		return
	}
	for rows.Next() {
		rows.Scan(&contactType.Ctypeid, &contactType.Name)
	}

	// sambungkan contact dengan contact type
	q = fmt.Sprintf("INSERT INTO `contactwtype` VALUES ('%d','%d')", contactID, contactType.Ctypeid)
	_, err = Exec(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": DATABASE_EXEC_FAIL,
			"message": fmt.Sprintf("DATABASE_EXEC_FAIL: %s", err.Error())})
		return
	}

	// ambil record tersimpan
	contacts := GetContact(contactID, 0, 0)
	if len(contacts) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"code": MODULE_OPERATION_FAIL,
			"message": fmt.Sprintf("MODULE_OPERATION_FAIL: insert contact fail, inserted %d", len(contacts))})
		return
	}
	contact := contacts[0]
	c.JSON(http.StatusCreated, contact)
	return
}

func GetContactAllHandler(c *gin.Context) {
	var records []ContactOut
	var segments = strings.Split(c.Param("path1"), "/")
	var l, o int64
	var limit, offset int
	var err error

	if len(segments) == 3 {
		limit = 10
		offset, err = strconv.Atoi(segments[2])
	} else if len(segments) == 4 {
		limit, err = strconv.Atoi(segments[3])
		offset, err = strconv.Atoi(segments[2])
	} else {
		limit = 10
		offset = 0
	}

	if err == nil { // tidak ada error dari konversi
		l = int64(limit)
		o = int64(offset)
	}

	records = GetContact(-1, l, o)

	c.JSON(http.StatusOK, records)
	return
}

func GetContactIdHandler(c *gin.Context) {
	var records []ContactOut
	var record ContactOut
	var segments = strings.Split(c.Param("path1"), "/")
	var id int64 = 1

	i, e := strconv.Atoi(segments[1])

	if e == nil { // konversi berhasil
		id = int64(i)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"code": INPUT_VALIDATION_FAIL,
			"message": fmt.Sprintf("INPUT_VALIDATION_FAIL: %s", e.Error())})
		return
	}

	records = GetContact(id, 0, 0)
	if len(records) > 0 {
		record = records[0]
	} else {
		c.JSON(http.StatusNotFound, NOT_FOUND)
		return
	}

	c.JSON(http.StatusOK, record)
	return
}

func PutContactIdHandler(c *gin.Context) {
	var records []ContactOut
	var record ContactOut
	var oldcontacttype, newcontacttype Contacttype
	var contactwtype Contactwtype
	var segments = strings.Split(c.Param("path1"), "/")
	var id int64 = 1
	var input ContactIn

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": INPUT_VALIDATION_FAIL,
			"message": fmt.Sprintf("INPUT_VALIDATION_FAIL: %s", err.Error())})
		return
	}

	i, e := strconv.Atoi(segments[1])
	if e == nil { // konversi berhasil
		id = int64(i)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"code": INPUT_VALIDATION_FAIL,
			"message": fmt.Sprintf("INPUT_VALIDATION_FAIL: %s", e.Error())})
		return
	}

	records = GetContact(id, 0, 0)
	if len(records) == 0 {
		c.JSON(http.StatusNotFound, NOT_FOUND)
		return
	}
	record = records[0]

	q := fmt.Sprintf("select * from contacttype where name = '%s'", record.Type)
	rows, err := Query(q)
	common.ErrHandler(err)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": DATABASE_EXEC_FAIL,
			"message": fmt.Sprintf("DATABASE_EXEC_FAIL: %s", err.Error())})
		return
	}
	for rows.Next() {
		rows.Scan(&oldcontacttype.Ctypeid, &oldcontacttype.Name)
	}

	q = fmt.Sprintf("select * from contacttype where name = '%s'", input.Type)
	rows, err = Query(q)
	common.ErrHandler(err)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": DATABASE_EXEC_FAIL,
			"message": fmt.Sprintf("DATABASE_EXEC_FAIL: %s", err.Error())})
		return
	}
	for rows.Next() {
		rows.Scan(&newcontacttype.Ctypeid, &newcontacttype.Name)
	}

	q = fmt.Sprintf("select * from contactwtype where contact_contactid = %d and contacttype_ctypeid = %d", record.Id, oldcontacttype.Ctypeid)
	rows, err = Query(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": DATABASE_EXEC_FAIL,
			"message": fmt.Sprintf("DATABASE_EXEC_FAIL: %s", err.Error())})
		return
	}
	for rows.Next() {
		rows.Scan(&contactwtype.Contact_contactid, &contactwtype.Contacttype_ctypeid)
	}

	q = fmt.Sprintf(`update contact
		set fname = '%s',
			lname = '%s',
		    prefix = '%s'
		where contactid = %d`, input.Firstname, input.Lastname, input.Prefix, id)
	result, err := Exec(q)
	common.ErrHandler(err)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": DATABASE_EXEC_FAIL,
			"message": fmt.Sprintf("DATABASE_EXEC_FAIL: %s", err.Error())})
		return
	}
	log.Printf("\nresult: %+v\n", result)

	q = fmt.Sprintf(`update contactwtype
			set contacttype_ctypeid = %d
		where
			contact_contactid = %d and
			contacttype_ctypeid = %d`, newcontacttype.Ctypeid, id, oldcontacttype.Ctypeid)
	result, err = Exec(q)
	common.ErrHandler(err)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": DATABASE_EXEC_FAIL,
			"message": fmt.Sprintf("DATABASE_EXEC_FAIL: %s", err.Error())})
		return
	}
	log.Printf("\nresult: %+v\n", result)

	records = GetContact(id, 0, 0)
	record = records[0]

	c.JSON(http.StatusOK, record)
	return
}

func DeleteContactIdHandler(c *gin.Context) {
	var segments = strings.Split(c.Param("path1"), "/")
	var oldcontacttype Contacttype
	var id int64 = 1

	i, e := strconv.Atoi(segments[1])
	if e == nil { // konversi berhasil
		id = int64(i)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"code": INPUT_VALIDATION_FAIL,
			"message": fmt.Sprintf("INPUT_VALIDATION_FAIL: %s", e.Error())})
		return
	}

	contacts := GetContact(id, 0, 0)
	log.Printf("\ncontacts: %+v\n", len(contacts))
	if len(contacts) > 0 {
		contact := contacts[0]

		q := fmt.Sprintf("select * from contacttype where name = '%s'", contact.Type)
		rows, err := Query(q)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": DATABASE_EXEC_FAIL,
				"message": fmt.Sprintf("DATABASE_EXEC_FAIL: %s", err.Error())})
			return
		}
		for rows.Next() {
			rows.Scan(&oldcontacttype.Ctypeid, &oldcontacttype.Name)
		}

		q = fmt.Sprintf(`UPDATE contact SET deleted=1 WHERE contactid = %d`, id)
		_, err = Exec(q)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": DATABASE_EXEC_FAIL,
				"message": fmt.Sprintf("DATABASE_EXEC_FAIL: %s", err.Error())})
			return
		}

		// rec := GetContact(id, 0, 0)
		// log.Printf("\nrec: %+v\n", rec)

		c.JSON(http.StatusOK, contact)
	} else {
		c.JSON(http.StatusNotFound, NOT_FOUND)
	}
}

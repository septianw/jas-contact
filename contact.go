package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"fmt"

	cpac "github.com/septianw/jas-contact/package"

	"github.com/gin-gonic/gin"
	"github.com/septianw/jas/common"
)

func Bootstrap() {
	log.Println("Contact module bootstrap.")
}

func Router(r *gin.Engine) {
	r.Any("/api/v1/contact/*path1", deflt)
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
			c.AbortWithStatusJSON(http.StatusMethodNotAllowed, cpac.NOT_ACCEPTABLE)
		}
		break
	case "GET":
		if strings.Compare(segments[1], "all") == 0 {
			GetContactAllHandler(c)
		} else if i, e := strconv.Atoi(segments[1]); (e == nil) && (i > 0) {
			GetContactIdHandler(c)
		} else {
			c.AbortWithStatusJSON(http.StatusNotAcceptable, cpac.NOT_ACCEPTABLE)
		}
		break
	case "PUT":
		if i, e := strconv.Atoi(segments[1]); (e == nil) && (i > 0) {
			PutContactIdHandler(c)
		} else {
			c.AbortWithStatusJSON(http.StatusMethodNotAllowed, cpac.NOT_ACCEPTABLE)
		}
		break
	case "DELETE":
		if i, e := strconv.Atoi(segments[1]); (e == nil) && (i > 0) {
			DeleteContactIdHandler(c)
		} else {
			c.AbortWithStatusJSON(http.StatusMethodNotAllowed, cpac.NOT_ACCEPTABLE)
		}
		break
	default:
		c.AbortWithStatusJSON(http.StatusMethodNotAllowed, cpac.NOT_ACCEPTABLE)
		break
	}
	// c.String(http.StatusOK, "hai")
}

func PostContactHandler(c *gin.Context) {
	var input cpac.ContactIn
	var contactType cpac.Contacttype

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": cpac.INPUT_VALIDATION_FAIL,
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
	result, err := cpac.Exec(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": cpac.DATABASE_EXEC_FAIL,
			"message": fmt.Sprintf("DATABASE_EXEC_FAIL: %s", err.Error())})
		return
	}
	contactID, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": cpac.DATABASE_EXEC_FAIL,
			"message": fmt.Sprintf("DATABASE_EXEC_FAIL: %s", err.Error())})
		return
	}

	// ambil contact type
	q = fmt.Sprintf("SELECT * FROM contacttype WHERE name = '%s'", input.Type)
	rows, err := cpac.Query(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": cpac.DATABASE_EXEC_FAIL,
			"message": fmt.Sprintf("DATABASE_EXEC_FAIL: %s", err.Error())})
		return
	}
	for rows.Next() {
		rows.Scan(&contactType.Ctypeid, &contactType.Name)
	}

	// sambungkan contact dengan contact type
	q = fmt.Sprintf("INSERT INTO `contactwtype` VALUES ('%d','%d')", contactID, contactType.Ctypeid)
	_, err = cpac.Exec(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": cpac.DATABASE_EXEC_FAIL,
			"message": fmt.Sprintf("DATABASE_EXEC_FAIL: %s", err.Error())})
		return
	}

	// ambil record tersimpan
	contacts := cpac.GetContact(contactID, 0, 0)
	if len(contacts) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"code": cpac.MODULE_OPERATION_FAIL,
			"message": fmt.Sprintf("MODULE_OPERATION_FAIL: insert contact fail, inserted %d", len(contacts))})
		return
	}
	contact := contacts[0]
	c.JSON(http.StatusCreated, contact)
	return
}

func GetContactAllHandler(c *gin.Context) {
	var records []cpac.ContactOut
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

	records = cpac.GetContact(-1, l, o)

	c.JSON(http.StatusOK, records)
	return
}

func GetContactIdHandler(c *gin.Context) {
	var records []cpac.ContactOut
	var record cpac.ContactOut
	var segments = strings.Split(c.Param("path1"), "/")
	var id int64 = 1

	i, e := strconv.Atoi(segments[1])

	if e == nil { // konversi berhasil
		id = int64(i)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"code": cpac.INPUT_VALIDATION_FAIL,
			"message": fmt.Sprintf("INPUT_VALIDATION_FAIL: %s", e.Error())})
		return
	}

	records = cpac.GetContact(id, 0, 0)
	if len(records) > 0 {
		record = records[0]
	} else {
		c.JSON(http.StatusNotFound, cpac.NOT_FOUND)
		return
	}

	c.JSON(http.StatusOK, record)
	return
}

func PutContactIdHandler(c *gin.Context) {
	var records []cpac.ContactOut
	var record cpac.ContactOut
	var oldcontacttype, newcontacttype cpac.Contacttype
	var contactwtype cpac.Contactwtype
	var segments = strings.Split(c.Param("path1"), "/")
	var id int64 = 1
	var input cpac.ContactIn

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": cpac.INPUT_VALIDATION_FAIL,
			"message": fmt.Sprintf("INPUT_VALIDATION_FAIL: %s", err.Error())})
		return
	}

	i, e := strconv.Atoi(segments[1])
	if e == nil { // konversi berhasil
		id = int64(i)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"code": cpac.INPUT_VALIDATION_FAIL,
			"message": fmt.Sprintf("INPUT_VALIDATION_FAIL: %s", e.Error())})
		return
	}

	records = cpac.GetContact(id, 0, 0)
	if len(records) == 0 {
		c.JSON(http.StatusNotFound, cpac.NOT_FOUND)
		return
	}
	record = records[0]

	q := fmt.Sprintf("select * from contacttype where name = '%s'", record.Type)
	rows, err := cpac.Query(q)
	common.ErrHandler(err)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": cpac.DATABASE_EXEC_FAIL,
			"message": fmt.Sprintf("DATABASE_EXEC_FAIL: %s", err.Error())})
		return
	}
	for rows.Next() {
		rows.Scan(&oldcontacttype.Ctypeid, &oldcontacttype.Name)
	}

	q = fmt.Sprintf("select * from contacttype where name = '%s'", input.Type)
	rows, err = cpac.Query(q)
	common.ErrHandler(err)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": cpac.DATABASE_EXEC_FAIL,
			"message": fmt.Sprintf("DATABASE_EXEC_FAIL: %s", err.Error())})
		return
	}
	for rows.Next() {
		rows.Scan(&newcontacttype.Ctypeid, &newcontacttype.Name)
	}

	q = fmt.Sprintf("select * from contactwtype where contact_contactid = %d and contacttype_ctypeid = %d", record.Id, oldcontacttype.Ctypeid)
	rows, err = cpac.Query(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": cpac.DATABASE_EXEC_FAIL,
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
	result, err := cpac.Exec(q)
	common.ErrHandler(err)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": cpac.DATABASE_EXEC_FAIL,
			"message": fmt.Sprintf("DATABASE_EXEC_FAIL: %s", err.Error())})
		return
	}
	log.Printf("\nresult: %+v\n", result)

	q = fmt.Sprintf(`update contactwtype
			set contacttype_ctypeid = %d
		where
			contact_contactid = %d and
			contacttype_ctypeid = %d`, newcontacttype.Ctypeid, id, oldcontacttype.Ctypeid)
	result, err = cpac.Exec(q)
	common.ErrHandler(err)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": cpac.DATABASE_EXEC_FAIL,
			"message": fmt.Sprintf("DATABASE_EXEC_FAIL: %s", err.Error())})
		return
	}
	log.Printf("\nresult: %+v\n", result)

	records = cpac.GetContact(id, 0, 0)
	record = records[0]

	c.JSON(http.StatusOK, record)
	return
}

func DeleteContactIdHandler(c *gin.Context) {
	var segments = strings.Split(c.Param("path1"), "/")
	var oldcontacttype cpac.Contacttype
	var id int64 = 1

	i, e := strconv.Atoi(segments[1])
	if e == nil { // konversi berhasil
		id = int64(i)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"code": cpac.INPUT_VALIDATION_FAIL,
			"message": fmt.Sprintf("INPUT_VALIDATION_FAIL: %s", e.Error())})
		return
	}

	contacts := cpac.GetContact(id, 0, 0)
	log.Printf("\ncontacts: %+v\n", len(contacts))
	if len(contacts) > 0 {
		contact := contacts[0]

		q := fmt.Sprintf("select * from contacttype where name = '%s'", contact.Type)
		rows, err := cpac.Query(q)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": cpac.DATABASE_EXEC_FAIL,
				"message": fmt.Sprintf("DATABASE_EXEC_FAIL: %s", err.Error())})
			return
		}
		for rows.Next() {
			rows.Scan(&oldcontacttype.Ctypeid, &oldcontacttype.Name)
		}

		q = fmt.Sprintf(`UPDATE contact SET deleted=1 WHERE contactid = %d`, id)
		_, err = cpac.Exec(q)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": cpac.DATABASE_EXEC_FAIL,
				"message": fmt.Sprintf("DATABASE_EXEC_FAIL: %s", err.Error())})
			return
		}

		// rec := GetContact(id, 0, 0)
		// log.Printf("\nrec: %+v\n", rec)

		c.JSON(http.StatusOK, contact)
	} else {
		c.JSON(http.StatusNotFound, cpac.NOT_FOUND)
	}
}

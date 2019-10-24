package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"fmt"

	cpac "github.com/septianw/jas-contact"

	"github.com/gin-gonic/gin"
	"github.com/septianw/jas/common"
)

// const Version = cpac.Version

func main() {
	return
}

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
		common.SendHttpError(c, common.INPUT_VALIDATION_FAIL_CODE, err)
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
		common.SendHttpError(c, common.DATABASE_EXEC_FAIL_CODE, err)
		return
	}
	contactID, err := result.LastInsertId()
	if err != nil {
		common.SendHttpError(c, common.DATABASE_EXEC_FAIL_CODE, err)
		return
	}

	// ambil contact type
	q = fmt.Sprintf("SELECT * FROM contacttype WHERE name = '%s'", input.Type)
	rows, err := cpac.Query(q)
	if err != nil {
		common.SendHttpError(c, common.DATABASE_EXEC_FAIL_CODE, err)
		return
	}
	for rows.Next() {
		rows.Scan(&contactType.Ctypeid, &contactType.Name)
	}

	// sambungkan contact dengan contact type
	q = fmt.Sprintf("INSERT INTO `contactwtype` VALUES ('%d','%d')", contactID, contactType.Ctypeid)
	_, err = cpac.Exec(q)
	if err != nil {
		common.SendHttpError(c, common.DATABASE_EXEC_FAIL_CODE, err)
		return
	}

	// ambil record tersimpan
	contacts := cpac.GetContact(contactID, 0, 0)
	if len(contacts) == 0 {
		common.SendHttpError(c, common.MODULE_OPERATION_FAIL_CODE,
			errors.New(fmt.Sprintf("MODULE_OPERATION_FAIL: insert contact fail, inserted %d", len(contacts))))
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
		offset = 0
		limit, err = strconv.Atoi(segments[2])
		if err != nil {
			common.ErrHandler(err)
			common.SendHttpError(c, common.INPUT_VALIDATION_FAIL_CODE, errors.New(
				fmt.Sprintf("%+v should be numeric", segments[2])))
			return
		}
	} else if len(segments) == 4 {
		offset, err = strconv.Atoi(segments[3])
		if err != nil {
			log.Println(err.Error())
			common.SendHttpError(c, common.INPUT_VALIDATION_FAIL_CODE, errors.New(
				fmt.Sprintf("%+v should be numeric", segments[3])))
			return
		}
		limit, err = strconv.Atoi(segments[2])
		if err != nil {
			log.Println(err.Error())
			common.SendHttpError(c, common.INPUT_VALIDATION_FAIL_CODE, errors.New(
				fmt.Sprintf("%+v should be numeric", segments[2])))
			return
		}
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
		common.SendHttpError(c, common.INPUT_VALIDATION_FAIL_CODE, e)
		return
	}

	records = cpac.GetContact(id, 0, 0)
	if len(records) > 0 {
		record = records[0]
	} else {
		common.SendHttpError(c, common.RECORD_NOT_FOUND_CODE, errors.New("You are find something we can't found it here."))
		return
	}

	c.JSON(http.StatusOK, record)
	return
}

func PutContactIdHandler(c *gin.Context) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var records []cpac.ContactOut
	// var record cpac.ContactOut
	// var oldcontacttype, newcontacttype cpac.Contacttype
	// var contactwtype cpac.Contactwtype
	var segments = strings.Split(c.Param("path1"), "/")
	var id int64
	var input cpac.ContactIn

	if err := c.ShouldBindJSON(&input); err != nil {
		common.SendHttpError(c, common.INPUT_VALIDATION_FAIL_CODE, err)
		return
	}

	i, e := strconv.Atoi(segments[1])
	if e == nil { // konversi berhasil
		id = int64(i)
	} else {
		common.SendHttpError(c, common.INPUT_VALIDATION_FAIL_CODE, e)
		return
	}

	records = cpac.GetContact(id, 0, 0)
	log.Println(records)

	_, err := cpac.UpdateContact(id, input)

	if err != nil {
		if strings.Compare("Contact not found.", err.Error()) == 0 {
			common.SendHttpError(c, common.RECORD_NOT_FOUND_CODE, err)
			return
		} else {
			common.SendHttpError(c, common.DATABASE_EXEC_FAIL_CODE, err)
			return
		}
	}

	on := cpac.GetContact(id, 0, 0)
	log.Println(on)

	c.JSON(http.StatusOK, on[0])
	return
}

func DeleteContactIdHandler(c *gin.Context) {
	var segments = strings.Split(c.Param("path1"), "/")
	var id int64 = 1

	i, e := strconv.Atoi(segments[1])
	if e == nil { // konversi berhasil
		id = int64(i)
	} else {
		common.SendHttpError(c, common.INPUT_VALIDATION_FAIL_CODE, e)
		return
	}

	// contacts := cpac.GetContact(id, 0, 0)
	contact, err := cpac.DeleteContact(id)
	if err != nil {
		common.SendHttpError(c, common.DATABASE_EXEC_FAIL_CODE, err)
		return
	} else if (err != nil) && (strings.Compare("Contact not found.", err.Error()) == 0) {
		common.SendHttpError(c, common.RECORD_NOT_FOUND_CODE, err)
		return
	}

	c.JSON(http.StatusOK, contact)
}

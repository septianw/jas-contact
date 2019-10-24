package main

import (
	"fmt"

	"net/http"

	"github.com/gin-gonic/gin"
	pak "github.com/septianw/jas/common"
)

func main() {

	lib := pak.LoadSo("/home/asep/gocode/src/github.com/septianw/jas-contact/test/contact.so")
	bootsym, err := lib.Lookup("Bootstrap")
	pak.ErrHandler(err)

	routersym, err := lib.Lookup("Router")
	pak.ErrHandler(err)

	bootstrap := bootsym.(func())
	router := routersym.(func(*gin.Engine))

	bootstrap()

	e := gin.Default()

	router(e)

	srv := &http.Server{
		Addr:    "0.0.0.0:4519",
		Handler: e,
	}

	srv.ListenAndServe()

	fmt.Println("vim-go")
}

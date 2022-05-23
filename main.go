package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net"
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9-)?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func isValid(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
}

func getdomain(email string) (domain string, err error) {
	splitemail := strings.Split(email, "@")
	if len(splitemail) < 2 {
		err = errors.New("Not a valid email!, missing the @ symbol")
		return "", err
	}
	domain = splitemail[1]
	//fmt.Println(domain)
	return domain, err
}

func checkMx(domain string) (bool bool, err error) {
	bool = false
	_, err = net.LookupMX(domain)
	if err != nil {
		fmt.Println("Bad Domain: ", err)
		bool = false
		return false, err
	}
	//fmt.Println(mx)
	bool = true
	return true, err
}

func main() {
	router := gin.Default()
	router.GET("/validate/:Email", verifyHandler)
	err := router.Run("localhost:8081")
	if err != nil {
		log.Fatal(err)
	}
	//email := "nowayout@netnaija.com"
	//validator(email)
}

func verifyHandler(c *gin.Context) {
	email := c.Param("Email")
	validator(email)
	
}

func validator(email string) {
	if isValid(email) {
		domain, err := getdomain(email)
		_, err = checkMx(domain)
		if err != nil {
			return
		}

		fmt.Println("Email:", email, "is valid")
		return
	}
	if !isValid(email) {
		fmt.Println("Bad email")
		return
	}
}

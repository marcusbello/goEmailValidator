package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type Result struct {
	Email    string `json:"email"`
	Domain   string `json:"domain"`
	Validity string `json:"validity"`
	Reason   string `json:"reason"`
}

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9-)?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func isValid(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
}

func getDomain(email string) (domain string, err error) {
	split := strings.Split(email, "@")
	if len(split) < 2 {
		err = errors.New("not a valid email")
		return "", err
	}
	domain = split[1]
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

	port := os.Getenv("PORT")
	router := gin.Default()
	router.GET("/validate/:Email", verifyHandler)
	err := router.Run(":" + port)
	if err != nil {
		log.Fatal(err)
	}
}

func verifyHandler(c *gin.Context) {
	email := c.Param("Email")
	result := validator(email)
	c.JSON(http.StatusOK, result)
}

func validator(email string) (result Result) {
	result = Result{
		Email:    email,
		Domain:   "nil",
		Validity: "nil",
		Reason:   "nil",
	}
	domain, err := getDomain(email)
	result.Domain = domain
	if isValid(email) {
		_, err = checkMx(domain)
		if err != nil {
			fmt.Println("Error on domain")
			result.Reason = "No MX record set for domain"
			result.Validity = "Not a valid email"
			return
		}
		fmt.Println("Email:", email, "is valid")
		valid := "Email is valid"
		result.Validity = valid
		return
	}
	if !isValid(email) {
		fmt.Println("Bad email syntax")
		result.Reason = "Bad email syntax"
		result.Validity = "Not a valid email"
		return
	}
	return result
}

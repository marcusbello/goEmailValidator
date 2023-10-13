package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
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
		return bool, err
	}
	//fmt.Println(mx)
	bool = true
	return true, err
}

func main() {

	//port := os.Getenv("PORT")
	port := "3030"
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "Okay!")
	})
	router.POST("/validate-email", verifyHandler)
	err := router.Run(":" + port)
	if err != nil {
		log.Fatal(err)
	}
}

func verifyHandler(c *gin.Context) {
	body := struct {
		Email string
	}{}
	if err := c.BindJSON(&body); err != nil {
		errResponse := struct {
			Error string
		}{
			Error: err.Error(),
		}
		c.JSON(http.StatusBadRequest, errResponse)
	}
	result := validator(body.Email)
	c.JSON(http.StatusOK, result)
}

func validator(email string) (result *Result) {
	result = &Result{
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
			fmt.Println("error on domain")
			result.Reason = "no MX record set for domain"
			result.Validity = "not a valid email"
			return
		}
		fmt.Println("Email:", email, "is valid")
		result.Validity = "email is valid"
		result.Reason = "syntax and MX records checks out"
		return
	} else {
		fmt.Println("bad email syntax")
		result.Reason = "bad email syntax"
		result.Validity = err.Error()
	}
	return
}

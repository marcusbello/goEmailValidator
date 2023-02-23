# goEmailValidator
Email Validator Api written in Golang using Gin gonic as the webserver for routing api

## How it works: - 

Example if you send a GET or POST request to localhost e.g 

`http://localhost:3030/nowix@email.com`

This will be your response formatted in json:

```
{
"email": "nowix@email.com",
"domain": "email.com",
"validity": "email is valid",
"reason": "nil"
}
```
The validation is done by checking regex and then checking the domain mx if it exists.



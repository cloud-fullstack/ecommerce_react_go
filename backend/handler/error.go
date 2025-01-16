package handler

import (
	"fmt"
	"github.com/apex/log"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/lithammer/shortuuid/v3"
	"unicode"
)

// Logs fields with reference error code and a message, and responds to c's HTTP request with message and a reference error code.
// Logs can be searched by the reference error code
// userFields are fields sent to the user as JSON.
func logRespondError(c *gin.Context, fields log.Fields, code int, message string, userFields ...gin.H) {
	ref := shortuuid.New()
	fields["error_code"] = ref
	log.WithFields(fields).Errorf(message)
	retFields := gin.H{
		"error":           true,
		"message":         message,
		"reference_error": ref,
	}

	for _, v := range userFields {
		for key, value := range v {
			retFields[key] = value
		}
	}

	c.JSON(code, retFields)
}

func spaceCaps(s string) string {
	var s2 string
	for _, char := range s {
		if unicode.IsUpper(char) {
			s2 = s2 + " "
		}
		s2 = s2 + string(char)
	}
	return s2
}

func translate(err validator.FieldError) string {
	phrase := ""

	switch err.Tag() {
	case "gte":
		phrase = "must be greater than"
	case "required":
		phrase = "is required"
	case "ltfield", "ltefield":
		phrase = "must be less than"
	}


	return fmt.Sprintf(`%v %v %v`, spaceCaps(err.Field()), phrase, spaceCaps(err.Param()))
}

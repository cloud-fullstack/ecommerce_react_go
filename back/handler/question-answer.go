package handler

import (
	"github.com/apex/log"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"shopa/db"
	"strings"
)

func InsertQuestion(c *gin.Context) {
	dbConn := c.MustGet("dbConn").(*pgxpool.Pool)
	askingAvatar := c.MustGet("authAvatar").(string)
	in := struct {
		ProductID    string `json:"product_id"`
		QuestionText string `json:"question_text"`
	}{}
	err := c.ShouldBindJSON(&in)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "binding in:"+err.Error())
		return
	}
	in.QuestionText = formatQuestion(in.QuestionText)
	// make sure owner is not asking the question
	productOwner, err := db.ProductOwner(c.Request.Context(), dbConn, in.ProductID)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "getting product owner:"+err.Error())
		return
	}
	if productOwner == askingAvatar {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 400, "cannot create a question on own product")
		return
	}
	err = db.InsertQuestion(c.Request.Context(), dbConn, askingAvatar, in.ProductID, in.QuestionText)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "inserting question"+err.Error())
		return
	}

	c.JSON(200, gin.H{
		"message":    "inserted question",
		"product_id": in.ProductID,
		"question":   in.QuestionText,
	})
}

func formatQuestion(q string) string {
	firstLetter := strings.ToUpper(string(q[0]))
	q = firstLetter + q[1:]
	if q[len(q)-1] != '?' {
		q = q + "?"
	}
	return q
}

func DeleteQuestion(c *gin.Context) {
	dbConn := c.MustGet("dbConn").(*pgxpool.Pool)
	authAvatar := c.MustGet("authAvatar").(string)
	questionID := c.Param("questionID")
	questionOwner, err := db.QuestionOwner(c.Request.Context(), dbConn, questionID)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "grabbing question owner"+err.Error())
		return
	}
	if questionOwner != authAvatar {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 400, "cannot delete question you do not own"+err.Error())
		return
	}
	err = db.DeleteQuestion(c.Request.Context(), dbConn, questionID)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "deleting question"+err.Error())
		return
	}
	c.JSON(200, gin.H{
		"message": "deleted question",
		"id":      questionID,
	})
}

func InsertAnswer(c *gin.Context) {
	dbConn := c.MustGet("dbConn").(*pgxpool.Pool)
	authAvatar := c.MustGet("authAvatar").(string)
	in := struct {
		QuestionID string `json:"question_id"`
		AnswerText string `json:"answer_text"`
	}{}
	err := c.ShouldBindJSON(&in)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "binding in:"+err.Error())
		return
	}
	// do not let them answer their own question
	questionOwner, err := db.QuestionOwner(c.Request.Context(), dbConn, in.QuestionID)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "grabbing question owner"+err.Error())
		return
	}
	if questionOwner == authAvatar {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 400, "cannot answer own question")
		return
	}

	productID, err := db.QuestionProductID(c.Request.Context(), dbConn, in.QuestionID)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "grabbing question product id"+err.Error())
		return
	}

	boughtProduct, err := db.BoughtProduct(c.Request.Context(), dbConn, productID, authAvatar)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "grabbing boughtness of product"+err.Error())
		return
	}

	productOwner, err := db.ProductOwner(c.Request.Context(), dbConn, productID)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "grabbing product owner"+err.Error())
		return
	}

	if !boughtProduct && authAvatar != productOwner {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 400, "must be product owner or product buyer to answer")
		return
	}

	err = db.InsertAnswer(c.Request.Context(), dbConn, authAvatar, in.QuestionID, in.AnswerText)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "inserting answer"+err.Error())
		return
	}

	c.JSON(200, gin.H{
		"message":     "inserted answer",
		"question_id": in.QuestionID,
		"answer":      in.AnswerText,
	})
}

func DeleteAnswer(c *gin.Context) {
	dbConn := c.MustGet("dbConn").(*pgxpool.Pool)
	authAvatar := c.MustGet("authAvatar").(string)
	answerID := c.Param("answerID")
	// ensure answer deleter is the owner of the answer
	answerOwner, err := db.AnswerOwner(c.Request.Context(), dbConn, answerID)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "grabbing answer owner"+err.Error())
		return
	}
	if answerOwner != authAvatar {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 400, "cannot delete answer you do not own")
		return
	}
	err = db.DeleteAnswer(c.Request.Context(), dbConn, answerID)
	if err != nil {
		logRespondError(c, log.Fields{
			"IP": c.ClientIP(),
		}, 500, "deleting answer"+err.Error())
		return
	}
	c.JSON(200, gin.H{
		"message": "deleted answer",
		"id":      answerID,
	})
}

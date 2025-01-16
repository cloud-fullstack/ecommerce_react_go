package db

import (
	"context"
	"database/sql"
)

type Question struct {
	ID                          string   `json:"id"`
	Asker                       string   `json:"asker"`
	AskerLegacyName             string   `json:"asker_legacy_name"`
	ProductID                   string   `json:"product_id"`
	QuestionText                string   `json:"question_text"`
	AvatarViewingIsBuyer        bool     `json:"avatar_viewing_is_buyer"`
	AvatarViewingIsProductOwner bool     `json:"avatar_viewing_is_product_owner"`
	Answers                     []Answer `json:"answers"`
}

type Answer struct {
	ID                     string `json:"id"`
	Answerer               string `json:"answerer"`
	AnswererLegacyName     string `json:"answerer_legacy_name"`
	AnswerText             string `json:"answer_text"`
	AnswererIsProductOwner bool   `json:"answerer_is_product_owner"`
}

type AnswerNullable struct {
	ID                     sql.NullString `json:"id"`
	Answerer               sql.NullString `json:"answerer"`
	AnswererLegacyName     sql.NullString `json:"answerer_legacy_name"`
	AnswerText             sql.NullString `json:"answer_text"`
	AnswererIsProductOwner sql.NullBool   `json:"answerer_is_product_owner"`
}

func ProductsQuestionsAndAnswers(ctx context.Context, conn Conn, productID, avatarViewing string) ([]Question, error) {
	if avatarViewing == "" {
		avatarViewing = "00000000-0000-0000-0000-000000000000"
	}
	rows, err := conn.Query(ctx, `
	SELECT question.*,
    a.legacyname AS asker_legacy_name,
    answer.*,
    a2.legacyname as answerer_legacy_name,
    answer.owner IN
        (select product.owner from product where product.id = $1) as answerer_is_product_owner,
    $2 IN (SELECT customer_order.buyer from customer_order
           JOIN customer_order_product_line copl on copl.customer_order = customer_order.id
           WHERE product_id = $1 AND demo = false) as isABuyer,
    $2 IN (SELECT owner from product
        WHERE product.id = $1) as isProductOwner
	FROM question
		 LEFT JOIN avatar a ON question.owner = a.uuid
		 LEFT JOIN answer ON question.id = answer.question
		 LEFT JOIN avatar a2 ON answer.owner = a2.uuid
	WHERE question.product = $1`,
		productID, avatarViewing)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	questions := make(map[string]Question) // map[questionID]
	for rows.Next() {
		q := Question{}
		a := AnswerNullable{}
		err = rows.Scan(
			&q.ID,
			&q.Asker,
			&q.ProductID,
			&q.QuestionText,
			&q.AskerLegacyName,
			&a.ID,
			&a.Answerer,
			nil, // questionID
			&a.AnswerText,
			&a.AnswererLegacyName,
			&a.AnswererIsProductOwner,
			&q.AvatarViewingIsBuyer,
			&q.AvatarViewingIsProductOwner,
		)
		if err != nil {
			return nil, err
		}

		var currQuestion = questions[q.ID]
		if _, exists := questions[q.ID]; !exists {
			currQuestion = q
		}
		if a.ID.Valid {
			currQuestion.Answers = append(currQuestion.Answers, Answer{
				ID:                     a.ID.String,
				Answerer:               a.Answerer.String,
				AnswererLegacyName:     a.AnswererLegacyName.String,
				AnswerText:             a.AnswerText.String,
				AnswererIsProductOwner: a.AnswererIsProductOwner.Bool,
			})
		}
		questions[q.ID] = currQuestion
	}
	questionsSlice := []Question{}
	for _, q := range questions {
		questionsSlice = append(questionsSlice, q)
	}
	return questionsSlice, nil
}

func InsertQuestion(ctx context.Context, conn Conn, avatarID, productID, question string) error {
	_, err := conn.Exec(ctx, `
		INSERT INTO question (id, owner, product, text)
		VALUES (uuid_generate_v4(), $1, $2, $3)
	`, avatarID, productID, question)
	if err != nil {
		return err
	}
	return nil
}

func DeleteQuestion(ctx context.Context, conn Conn, questionID string) error {
	tx, err := conn.Begin(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `
		DELETE FROM answer
		WHERE question = $1
	`, questionID)
	_, err = tx.Exec(ctx, `
		DELETE FROM question
		WHERE id = $1
	`, questionID)
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	return nil
}

func QuestionOwner(ctx context.Context, conn Conn, questionID string) (string, error) {
	var owner string
	err := conn.QueryRow(ctx, `
		SELECT owner
		FROM question
		WHERE id = $1
	`, questionID).Scan(&owner)
	if err != nil {
		return "", err
	}
	return owner, nil
}

func QuestionProductID(ctx context.Context, conn Conn, questionID string) (string, error) {
	var productID string
	err := conn.QueryRow(ctx, `
		SELECT product
		FROM question
		WHERE id = $1
	`, questionID).Scan(&productID)
	if err != nil {
		return "", err
	}
	return productID, nil
}

func InsertAnswer(ctx context.Context, conn Conn, avatarID, questionID, answer string) error {
	_, err := conn.Exec(ctx, `
		INSERT INTO answer
		VALUES (uuid_generate_v4(), $1, $2, $3)
	`, avatarID, questionID, answer)
	if err != nil {
		return err
	}
	return nil
}

func DeleteAnswer(ctx context.Context, conn Conn, answerID string) error {
	_, err := conn.Exec(ctx, `
		DELETE FROM answer
		WHERE id = $1
	`, answerID)
	if err != nil {
		return err
	}
	return nil
}

func AnswerOwner(ctx context.Context, conn Conn, answerID string) (string, error) {
	var owner string
	err := conn.QueryRow(ctx, `
		SELECT owner
		FROM answer
		WHERE id = $1
	`, answerID).Scan(&owner)
	if err != nil {
		return "", err
	}
	return owner, nil
}

type FAQ struct {
	ID           string `json:"id"`
	ProductID    string `json:"product_id"`
	QuestionText string `json:"question_text"`
	AnswerText   string `json:"answer_text"`
}

func GetProductFAQs(ctx context.Context, conn Conn, productID string) ([]FAQ, error) {
	var faqs []FAQ
	rows, err := conn.Query(ctx, `
		SELECT id, product, question_text, answer_text
		FROM faq
		WHERE product = $1
	`, productID)
	defer rows.Close()
	for rows.Next() {
		var f FAQ
		err := rows.Scan(&f.ID, &f.ProductID, &f.QuestionText, &f.AnswerText)
		if err != nil {
			return nil, err
		}
		faqs = append(faqs, f)
	}
	if err != nil {
		return nil, err
	}
	return faqs, nil
}

func InsertFAQ(ctx context.Context, conn Conn, productID, question, answer string) error {
	_, err := conn.Exec(ctx, `
		INSERT INTO faq
		VALUES (uuid_generate_v4(), $1, $2, $3)
	`, productID, question, answer)
	if err != nil {
		return err
	}
	return nil
}

func DeleteAllFAQs(ctx context.Context, conn Conn, productID string) error {
	_, err := conn.Exec(ctx, `
		DELETE FROM faq
		WHERE product = $1
	`, productID)
	if err != nil {
		return err
	}
	return nil
}

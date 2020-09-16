package actions

import (
	"time"

	"github.com/YaleSpinup/tweaser/helpers"
	"github.com/YaleSpinup/tweaser/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
	"github.com/pkg/errors"
)

// QuestionsList returns the list of questions.
// /v1/tweaser/questions[?user_id=someguy]
func QuestionsList(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	questions := []models.Question{}
	if userid := c.Param("user_id"); userid == "" {
		// Paginate results. Params "page" and "per_page" control pagination.
		// Default values are "page=1" and "per_page=20".
		q := tx.PaginateFromParams(c.Params())

		// Retrieve all Questions from the DB
		if err := q.All(&questions); err != nil {
			return errors.WithStack(err)
		}
	} else {
		campaigns := []models.Campaign{}
		cq := tx.Select("id").Where("start_date <= ?", time.Now()).Where("end_date > ?", time.Now()).Where("enabled = true")
		err := cq.All(&campaigns)
		if err != nil {
			return c.Render(404, r.JSON([]string{}))
		}

		var campaignIDs []interface{}
		for _, c := range campaigns {
			campaignIDs = append(campaignIDs, c.ID.String())
		}

		// Paginate results. Params "page" and "per_page" control pagination.
		// Default values are "page=1" and "per_page=20".
		q := tx.PaginateFromParams(c.Params())
		q = q.Where("questions.enabled = true")
		q = q.Where("questions.campaign_id IN (?)", campaignIDs...)
		q = q.Where("id NOT in (select question_id FROM responses WHERE user_id = (?))", userid)
		err = q.All(&questions)
		if err != nil {
			return c.Render(404, r.JSON([]string{}))
		}

		// Generate a token for each question
		for i, q := range questions {
			mt := helpers.ModelToken{
				ID:     q.ID,
				Secret: CryptToken,
				UserID: userid,
			}

			token, err := mt.Generate()
			if err != nil {
				return c.Render(500, r.JSON("Internal server error."))
			}

			questions[i].Token = token

			// Get the enabled answers for the question and append them to the questions response
			answers := []models.Answer{}
			err = tx.Where("question_id = ?", q.ID).Where("enabled = true").All(&answers)
			if err != nil {
				return c.Render(500, r.JSON("Internal server error."))
			}
			questions[i].Answers = answers
		}
	}

	return c.Render(200, r.JSON(questions))
}

// QuestionsGet gets a question by ID.
// /v1/tweaser/questions/{question_id}
func QuestionsGet(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Question
	question := &models.Question{}

	// To find the Question the parameter question_id is used.
	if err := tx.Find(question, c.Param("question_id")); err != nil {
		return c.Error(404, err)
	}

	return c.Render(200, r.JSON(question))
}

// QuestionsGetAnswers gets the answers for a question by question ID.
// /v1/tweaser/questions/{question_id}/answers
func QuestionsGetAnswers(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Question
	question := &models.Question{}

	// To find the Question the parameter question_id is used.
	if err := tx.Eager("Answers").Find(question, c.Param("question_id")); err != nil {
		return c.Error(404, err)
	}

	return c.Render(200, r.JSON(question.Answers))
}

// QuestionsGetResponses gets the responses for a question by question ID.
// /v1/tweaser/questions/{question_id}/responses[?extended=true]
func QuestionsGetResponses(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	if e := c.Param("extended"); e != "" {
		// Allocate empty Responses model
		responses := models.Responses{}
		if err := tx.Eager().Where("question_id = (?)", c.Param("question_id")).All(&responses); err != nil {
			return c.Error(404, err)
		}
		return c.Render(200, r.JSON(responses))
	}

	question := &models.Question{}
	if err := tx.Eager("Answers").Find(question, c.Param("question_id")); err != nil {
		return c.Error(404, err)
	}

	counts := map[string]int{}
	answers := map[string]string{}
	for _, a := range question.Answers {
		if !a.Enabled {
			continue
		}

		id := a.ID.String()
		rq := []models.ResponseAnswer{}
		count, err := tx.Where("answer_id = (?)", id).Count(&rq)
		if err != nil {
			return c.Error(404, err)
		}
		counts[id] = count
		answers[id] = a.Text
	}

	resp := struct {
		Count   map[string]int    `json:"count"`
		Answers map[string]string `json:"answers"`
	}{
		Count:   counts,
		Answers: answers,
	}
	return c.Render(200, r.JSON(resp))
}

// QuestionsCreate creates an question.
func QuestionsCreate(c buffalo.Context) error {
	// Allocate an empty Question
	question := &models.Question{}

	// bind the request body to the new question
	if err := c.Bind(question); err != nil {
		return errors.WithStack(err)
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Validate the posted data and save it to the database
	verrs, err := tx.ValidateAndCreate(question)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		return c.Render(422, r.JSON(verrs))
	}

	return c.Render(201, r.JSON(question))
}

// QuestionsUpdate updates an question.
func QuestionsUpdate(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Question
	question := &models.Question{}

	if err := tx.Find(question, c.Param("question_id")); err != nil {
		return c.Error(404, err)
	}

	// bind the request body to the question
	if err := c.Bind(question); err != nil {
		return errors.WithStack(err)
	}

	verrs, err := tx.ValidateAndUpdate(question)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		return c.Render(422, r.JSON(verrs))
	}

	return c.Render(200, r.JSON(question))
}

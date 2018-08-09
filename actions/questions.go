package actions

import (
	"github.com/YaleUniversity/tweaser/helpers"
	"github.com/YaleUniversity/tweaser/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
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

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())

	if userid := c.Param("user_id"); userid == "" {
		// Retrieve all Questions from the DB
		if err := q.All(&questions); err != nil {
			return errors.WithStack(err)
		}
	} else {
		// select * from questions where id not in (select question_id from responses where user_id ="cf322");
		q = q.RawQuery("SELECT * FROM questions WHERE id NOT IN (select question_id FROM responses WHERE user_id = ?) and enabled = true", userid)
		err := q.Eager("Answers").All(&questions)
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
	if err := tx.Eager("Answers").Find(question, c.Param("question_id")); err != nil {
		return c.Error(404, err)
	}

	return c.Render(200, r.JSON(question))
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

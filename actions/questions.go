package actions

import (
	"git.yale.edu/spinup/tweaser/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
)

// QuestionsList returns the list of questions.
// /v1/tweaser/questions
func QuestionsList(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	questions := &models.Questions{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.Eager().PaginateFromParams(c.Params())

	// Retrieve all Questions from the DB
	if err := q.All(questions); err != nil {
		return errors.WithStack(err)
	}

	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)

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

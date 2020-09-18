package actions

import (
	"github.com/YaleSpinup/tweaser/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
	"github.com/pkg/errors"
)

// AnswersList gets a paginated list of answers.
func AnswersList(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	answers := &models.Answers{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())

	// Retrieve all Answers from the DB
	if err := q.All(answers); err != nil {
		return errors.WithStack(err)
	}

	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)

	return c.Render(200, r.JSON(answers))
}

// AnswersGet gets an answer by ID.
func AnswersGet(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Answer
	answer := &models.Answer{}

	// To find the Answer the parameter answer_id is used.
	if err := tx.Find(answer, c.Param("answer_id")); err != nil {
		return c.Error(404, err)
	}

	return c.Render(200, r.JSON(answer))
}

// AnswersCreate creates an answer.
func AnswersCreate(c buffalo.Context) error {
	// Allocate an empty Answer
	answer := &models.Answer{}

	// bind the request body to the new answer
	if err := c.Bind(answer); err != nil {
		return errors.WithStack(err)
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Validate the posted data and save it to the database
	verrs, err := tx.ValidateAndCreate(answer)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		return c.Render(422, r.JSON(verrs))
	}

	return c.Render(201, r.JSON(answer))
}

// AnswersUpdate updates an answer.
func AnswersUpdate(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Answer
	answer := &models.Answer{}

	if err := tx.Find(answer, c.Param("answer_id")); err != nil {
		return c.Error(404, err)
	}

	// bind the request body to the answer
	if err := c.Bind(answer); err != nil {
		return errors.WithStack(err)
	}

	verrs, err := tx.ValidateAndUpdate(answer)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		return c.Render(422, r.JSON(verrs))
	}

	return c.Render(200, r.JSON(answer))
}

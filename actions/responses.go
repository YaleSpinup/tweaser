package actions

import (
	"log"

	"github.com/YaleUniversity/tweaser/helpers"
	"github.com/YaleUniversity/tweaser/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
)

// ResponsesList default implementation.
func ResponsesList(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	responses := &models.Responses{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())

	// Retrieve all Responses from the DB
	if err := q.All(responses); err != nil {
		return errors.WithStack(err)
	}

	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)

	return c.Render(200, r.JSON(responses))
}

// ResponsesGet default implementation.
func ResponsesGet(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Question
	response := &models.Response{}

	// To find the Response the parameter response_id is used.
	if err := tx.Find(response, c.Param("response_id")); err != nil {
		return c.Error(404, err)
	}

	return c.Render(200, r.JSON(response))
}

// ResponsesCreate default implementation.
func ResponsesCreate(c buffalo.Context) error {
	token := c.Param("token")
	if token == "" {
		return c.Error(403, errors.New("Unauthorized."))
	}

	// Allocate an empty Response
	response := &models.Response{}

	// bind the request body to the new response
	if err := c.Bind(response); err != nil {
		return errors.WithStack(err)
	}

	mt := helpers.ModelToken{
		ID:     response.QuestionID,
		Secret: CryptToken,
		UserID: response.UserID,
	}
	err := mt.Validate(token)
	if err != nil {
		return c.Render(403, r.JSON("Unauthorized. Invalid Token."))
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	if err = tx.Find(&response.Question, response.QuestionID); err != nil {
		return c.Render(404, r.JSON("Question Not Found."))
	}

	// Validate the posted data and save it to the database
	verrs, err := tx.ValidateAndCreate(response)
	if err != nil {
		return errors.WithStack(err)
	}

	for _, a := range response.Answers {
		ra := &models.ResponseAnswer{
			ResponseID: response.ID,
			AnswerID:   a.ID,
			QuestionID: response.QuestionID,
		}

		// Validate the posted data and save it to the database
		ve, err := tx.ValidateAndCreate(ra)
		if err != nil {
			return errors.WithStack(err)
		}
		verrs.Append(ve)
	}

	if verrs.HasAny() {
		return c.Render(422, r.JSON(verrs))
	}

	log.Println("done with update")

	return c.Render(202, r.JSON("submitted"))
}

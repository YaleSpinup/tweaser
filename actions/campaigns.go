package actions

import (
	"git.yale.edu/spinup/tweaser/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
)

// CampaignsList gets a paginated list of campaigns.
// /v1/tweaser/campaigns
func CampaignsList(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	campaigns := &models.Campaigns{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())

	// Retrieve all Campaigns from the DB
	if err := q.All(campaigns); err != nil {
		return errors.WithStack(err)
	}

	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)

	return c.Render(200, r.JSON(campaigns))
}

// CampaignsGet gets a Campaign by ID.
// /v1/tweaser/campaigns/{campaign_id}
func CampaignsGet(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Campaign
	campaign := &models.Campaign{}

	// To find the Campaign the parameter campaign_id is used.
	if err := tx.Eager("Questions").Find(campaign, c.Param("campaign_id")); err != nil {
		return c.Error(404, err)
	}

	return c.Render(200, r.JSON(campaign))
}

// CampaignsGetQuestions gets the questions for a given Campaign.
// /v1/tweaser/campaigns/{campaign_id}/questions
func CampaignsGetQuestions(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Campaign
	campaign := &models.Campaign{}

	// To find the Campaign the parameter campaign_id is used.
	if err := tx.Eager("Questions").Find(campaign, c.Param("campaign_id")); err != nil {
		return c.Error(404, err)
	}

	return c.Render(200, r.JSON(campaign.Questions))
}

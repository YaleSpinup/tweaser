package actions

import (
	"strconv"
	"time"

	"github.com/YaleUniversity/tweaser/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
)

// CampaignsList gets a paginated list of campaigns.
// GET /v1/tweaser/campaigns[?active=true|false][&enabled=true|false]
func CampaignsList(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	campaigns := &models.Campaigns{}

	q := tx.Q()
	if active, err := strconv.ParseBool(c.Param("active")); err == nil {
		if active {
			q = q.Where("start_date <= ?", time.Now()).Where("end_date > ?", time.Now())
		} else {
			q = q.Where("(start_date > ? or end_date <= ?)", time.Now(), time.Now())
		}
	}

	if enabled, err := strconv.ParseBool(c.Param("enabled")); err == nil {
		q = q.Where("enabled = ?", enabled)
	}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q = q.PaginateFromParams(c.Params())

	// Retrieve all Campaigns from the DB
	if err := q.All(campaigns); err != nil {
		return errors.WithStack(err)
	}

	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)

	return c.Render(200, r.JSON(campaigns))
}

// CampaignsGet gets a Campaign by ID
// GET /v1/tweaser/campaigns/{campaign_id}
func CampaignsGet(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Campaign
	campaign := &models.Campaign{}

	// To find the Campaign the parameter campaign_id is used.
	if err := tx.Find(campaign, c.Param("campaign_id")); err != nil {
		return c.Error(404, err)
	}

	return c.Render(200, r.JSON(campaign))
}

// CampaignsGetQuestions gets a Campaign's questions by campaign ID.
// GET /v1/tweaser/campaigns/{campaign_id}/questions
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

// CampaignsCreate creates a new Campaign in the database
// POST /v1/tweaser/campaigns
func CampaignsCreate(c buffalo.Context) error {
	// Allocate an empty Campaign
	campaign := &models.Campaign{}

	// bind the request body to the new campaign
	if err := c.Bind(campaign); err != nil {
		return errors.WithStack(err)
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Validate the posted data and save it to the database
	verrs, err := tx.ValidateAndCreate(campaign)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		return c.Render(422, r.JSON(verrs))
	}

	return c.Render(201, r.JSON(campaign))
}

// CampaignsUpdate changes a Campaign in the DB.
// PUT /v1/tweaser/campaigns/{campaign_id}
func CampaignsUpdate(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Campaign
	campaign := &models.Campaign{}

	if err := tx.Find(campaign, c.Param("campaign_id")); err != nil {
		return c.Error(404, err)
	}

	// Bind Campaign to request body
	if err := c.Bind(campaign); err != nil {
		return errors.WithStack(err)
	}

	verrs, err := tx.ValidateAndUpdate(campaign)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		return c.Render(422, r.JSON(verrs))
	}

	return c.Render(200, r.JSON(campaign))
}

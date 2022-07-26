package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

type Question struct {
	ID         uuid.UUID `json:"id" db:"id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
	Text       string    `json:"text" db:"text"`
	Campaign   Campaign  `belongs_to:"campaign" json:"-"`
	CampaignID uuid.UUID `json:"campaign_id" db:"campaign_id"`
	Enabled    bool      `json:"enabled" db:"enabled"`
	Answers    Answers   `has_many:"answers" json:"answers,omitempty"`
	Type       string    `json:"type" db:"type"`
	Token      string    `json:"token,omitempty" db:"-"`
}

// String is not required by pop and may be deleted
func (q Question) String() string {
	jq, _ := json.Marshal(q)
	return string(jq)
}

// Questions is not required by pop and may be deleted
type Questions []Question

// String is not required by pop and may be deleted
func (q Questions) String() string {
	jq, _ := json.Marshal(q)
	return string(jq)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (q *Question) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: q.Text, Name: "Text"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (q *Question) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (q *Question) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

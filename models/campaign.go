package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

type Campaign struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Name      string    `json:"name" db:"name"`
	StartDate time.Time `json:"start_date" db:"start_date"`
	EndDate   time.Time `json:"end_date" db:"end_date"`
	Enabled   bool      `json:"enabled" db:"enabled"`
	Questions Questions `has_many:"questions" json:"questions,omitempty"`
}

// String is not required by pop and may be deleted
func (c Campaign) String() string {
	jc, _ := json.Marshal(c)
	return string(jc)
}

// Campaigns is not required by pop and may be deleted
type Campaigns []Campaign

// String is not required by pop and may be deleted
func (c Campaigns) String() string {
	jc, _ := json.Marshal(c)
	return string(jc)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (c *Campaign) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: c.Name, Name: "Name"},
		&validators.TimeIsPresent{Field: c.StartDate, Name: "StartDate"},
		&validators.TimeIsPresent{Field: c.EndDate, Name: "EndDate"},
		&validators.TimeIsBeforeTime{FirstName: "StartDate", FirstTime: c.StartDate, SecondName: "EndTime", SecondTime: c.EndDate},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (c *Campaign) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (c *Campaign) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

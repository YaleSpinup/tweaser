package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

type Response struct {
	ID         uuid.UUID `json:"id" db:"id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
	UserID     string    `json:"user_id" db:"user_id"`
	Text       string    `json:"text" db:"text"`
	Question   Question  `belongs_to:"question" json:"-"`
	QuestionID uuid.UUID `json:"question_id" db:"question_id"`
	Answers    Answers   `many_to_many:"response_answers"`
}

// String is not required by pop and may be deleted
func (r Response) String() string {
	jr, _ := json.Marshal(r)
	return string(jr)
}

// Responses is not required by pop and may be deleted
type Responses []Response

// String is not required by pop and may be deleted
func (r Responses) String() string {
	jr, _ := json.Marshal(r)
	return string(jr)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (r *Response) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: r.QuestionID, Name: "QuestionID"},
		&validators.StringIsPresent{Field: r.UserID, Name: "UserID"},
		&UserAlreadyResponded{UserID: r.UserID, QuestionID: r.QuestionID, tx: tx, Name: "UserAlreadyResponded"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (r *Response) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (r *Response) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

type UserAlreadyResponded struct {
	Name       string
	UserID     string
	QuestionID uuid.UUID
	tx         *pop.Connection
}

func (u *UserAlreadyResponded) IsValid(errors *validate.Errors) {
	response := Response{}
	query := u.tx.Where("question_id = ?", u.QuestionID).Where("user_id = ?", u.UserID)
	err := query.First(&response)
	if err == nil {
		// found response from this user
		errors.Add(validators.GenerateKey(u.Name), fmt.Sprintf("User %s has already responded to question %s.", u.UserID, u.QuestionID))
	}
}

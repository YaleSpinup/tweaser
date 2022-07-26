package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

type ResponseAnswer struct {
	ID         uuid.UUID `json:"id" db:"id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
	Answer     Answer    `belongs_to:"answer" db:"-"`
	Response   Response  `belongs_to:"response" db:"-"`
	AnswerID   uuid.UUID `json:"answer_id" db:"answer_id"`
	ResponseID uuid.UUID `json:"response_id" db:"response_id"`
	QuestionID uuid.UUID `json:"-" db:"-"`
}

// String is not required by pop and may be deleted
func (r ResponseAnswer) String() string {
	jr, _ := json.Marshal(r)
	return string(jr)
}

// ResponseAnswers is not required by pop and may be deleted
type ResponseAnswers []ResponseAnswer

// String is not required by pop and may be deleted
func (r ResponseAnswers) String() string {
	jr, _ := json.Marshal(r)
	return string(jr)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (r *ResponseAnswer) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: r.ResponseID, Name: "ResponseID"},
		&validators.UUIDIsPresent{Field: r.AnswerID, Name: "AnswerID"},
		&AnswerBelongsToQuestion{Name: "AnswerBelongsToQuestion", AnswerID: r.AnswerID, QuestionID: r.QuestionID, tx: tx},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (r *ResponseAnswer) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (r *ResponseAnswer) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// AnswerBelongsToQuestion is a custom validator for the answer/response join table
type AnswerBelongsToQuestion struct {
	Name       string
	AnswerID   uuid.UUID
	QuestionID uuid.UUID
	tx         *pop.Connection
}

// IsValid validates that the posted answer belongs to the correct question
func (v *AnswerBelongsToQuestion) IsValid(errors *validate.Errors) {
	answer := &Answer{}
	err := v.tx.Find(answer, v.AnswerID)
	if err != nil {
		errors.Add(validators.GenerateKey(v.Name), "Answer ID not found in db")
	}

	if answer.QuestionID != v.QuestionID {
		errors.Add(validators.GenerateKey(v.Name), "Given answer does not belong to the given question in response")
	}
}

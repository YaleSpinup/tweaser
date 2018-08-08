package helpers

import (
	"encoding/base64"
	"encoding/json"

	log "github.com/sirupsen/logrus"

	"github.com/gobuffalo/uuid"
	"golang.org/x/crypto/bcrypt"
)

// ModelToken is the object used to generate a token loosesly associated with a model
type ModelToken struct {
	UserID string    `json:"user_id"`
	ID     uuid.UUID `json:"id"`
	Secret string    `json:"secret"`
}

// Generate creates a base64 encoded renewal token
func (r *ModelToken) Generate() (string, error) {
	str, err := json.Marshal(r)
	if err != nil {
		return "", err
	}

	log.Debugf("Marshalled secret JSON string %s", str)

	token, err := bcrypt.GenerateFromPassword(str, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	log.Debugln("Secret hash:", string(token))

	return base64.StdEncoding.EncodeToString(token), nil
}

// Validate validates a base64 encoded renewal token
func (r *ModelToken) Validate(token string) error {
	decodedToken, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		log.Error("Failed to decode base64", err)
		return err
	}

	str, err := json.Marshal(r)
	if err != nil {
		log.Error("Failed to marshall JSON", err)
		return err
	}

	log.Debugf("Comparing decodedToken: %s with renewalSecret %s", decodedToken, string(str))

	return bcrypt.CompareHashAndPassword(decodedToken, str)
}

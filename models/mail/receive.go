package models

import (
	"encoding/json"
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
	"time"
)

// MailReceive is used by pop to map your mail_receives database table to your go code.
type MailReceive struct {
	ID         int       `json:"id" db:"id"`
	UID        uuid.UUID `json:"uid" db:"uid"`
	From       string    `json:"from" db:"from"`
	To         string    `json:"to" db:"to"`
	Subject    string    `json:"subject" db:"subject"`
	Body       string    `json:"body" db:"body"`
	Attachment string    `json:"attachment" db:"attachment"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// String is not required by pop and may be deleted
func (m MailReceive) String() string {
	jm, _ := json.Marshal(m)
	return string(jm)
}

// MailReceives is not required by pop and may be deleted
type MailReceives []MailReceive

// String is not required by pop and may be deleted
func (m MailReceives) String() string {
	jm, _ := json.Marshal(m)
	return string(jm)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (m *MailReceive) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.IntIsPresent{Field: m.ID, Name: "ID"},
		&validators.StringIsPresent{Field: m.From, Name: "From"},
		&validators.StringIsPresent{Field: m.To, Name: "To"},
		&validators.StringIsPresent{Field: m.Subject, Name: "Subject"},
		&validators.StringIsPresent{Field: m.Body, Name: "Body"},
		&validators.StringIsPresent{Field: m.Attachment, Name: "Attachment"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (m *MailReceive) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (m *MailReceive) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

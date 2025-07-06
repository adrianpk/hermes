package ssg

import (
	"encoding/json"

	"github.com/adrianpk/hermes/internal/am"
	"github.com/google/uuid"
)

const (
	contentType = "content"
)

type Content struct {
	*am.BaseModel
	UserID  uuid.UUID
	Heading string `json:"heading"`
	Body    string `json:"body"`
	Status  string
}

func NewContent(heading, body string) Content {
	return Content{
		BaseModel: am.NewModel(am.WithType(contentType)),
		Heading:   heading,
		Body:      body,
	}
}

func (r *Content) Slug() string {
	return am.Normalize(r.Heading) + "-" + r.ShortID()
}

func (r *Content) OptLabel() string {
	return r.Heading
}

// UnmarshalJSON ensures Model is always initialized after unmarshal.
func (r *Content) UnmarshalJSON(data []byte) error {
	type Alias Content
	temp := &Alias{}
	if err := json.Unmarshal(data, temp); err != nil {
		return err
	}
	*r = Content(*temp)
	if r.BaseModel == nil {
		r.BaseModel = am.NewModel(am.WithType(contentType))
	}
	return nil
}

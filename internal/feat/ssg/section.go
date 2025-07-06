package ssg

import (
	"encoding/json"

	"github.com/adrianpk/hermes/internal/am"
	"github.com/google/uuid"
)

const (
	sectionType = "section"
)

type Section struct {
	*am.BaseModel
	UserID      uuid.UUID
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Path        string    `json:"path"`
	LayoutID    uuid.UUID `json:"layout_id"`
	Image       string    `json:"image"`
	Header      string    `json:"header"`
}

func NewSection(name, description, path string, layoutID uuid.UUID) Section {
	return Section{
		BaseModel:   am.NewModel(am.WithType(sectionType)),
		Name:        name,
		Description: description,
		Path:        path,
		LayoutID:    layoutID,
	}
}

func (s *Section) Slug() string {
	return am.Normalize(s.Name) + "-" + s.ShortID()
}

// UnmarshalJSON ensures Model is always initialized after unmarshal.
func (s *Section) UnmarshalJSON(data []byte) error {
	type Alias Section
	temp := &Alias{}
	if err := json.Unmarshal(data, temp); err != nil {
		return err
	}
	*s = Section(*temp)
	if s.BaseModel == nil {
		s.BaseModel = am.NewModel(am.WithType(sectionType))
	}
	return nil
}

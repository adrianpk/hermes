package ssg

import (
	"encoding/json"

	"github.com/adrianpk/hermes/internal/am"
	"github.com/google/uuid"
)

const (
	layoutType = "layout"
)

type Layout struct {
	*am.BaseModel
	Name        string `json:"name"`
	Description string `json:"description"`
	Code        string `json:"header"`
}

func Newlayout(name, description, path string, layoutID uuid.UUID) Layout {
	return Layout{
		BaseModel:   am.NewModel(am.WithType(layoutType)),
		Name:        name,
		Description: description,
		Code:        path,
	}
}

func (s *Layout) Slug() string {
	return am.Normalize(s.Name) + "-" + s.ShortID()
}

func (s Layout) OptValue() string {
	return s.ID().String()
}

func (s Layout) OptLabel() string {
	return s.Name
}

// UnmarshalJSON ensures Model is always initialized after unmarshal.
func (s *Layout) UnmarshalJSON(data []byte) error {
	type Alias Layout
	temp := &Alias{}
	if err := json.Unmarshal(data, temp); err != nil {
		return err
	}
	*s = Layout(*temp)
	if s.BaseModel == nil {
		s.BaseModel = am.NewModel(am.WithType(layoutType))
	}
	return nil
}

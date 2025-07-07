package ssg

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"

	"github.com/adrianpk/hermes/internal/am"
	"github.com/google/uuid"
)

type Seeder struct {
	*am.JSONSeeder
	repo Repo
}

type SeedData struct {
	Layouts  []Layout  `json:"layouts"`
	Sections []Section `json:"sections"`
}

func NewSeeder(assetsFS embed.FS, engine string, repo Repo) *Seeder {
	return &Seeder{
		JSONSeeder: am.NewJSONSeeder(ssgFeat, assetsFS, engine),
		repo:       repo,
	}
}

func (s *Seeder) Setup(ctx context.Context) error {
	if err := s.JSONSeeder.Setup(ctx); err != nil {
		return err
	}
	return s.SeedAll(ctx)
}

func (s *Seeder) SeedAll(ctx context.Context) error {
	s.Log().Info("Seeding SSG data...")
	byFeature, err := s.JSONSeeder.LoadJSONSeeds()
	if err != nil {
		return fmt.Errorf("failed to load JSON seeds: %w", err)
	}
	const ssgFeat = "ssg"
	for feature, seeds := range byFeature {
		if feature != ssgFeat {
			continue
		}
		fmt.Printf("Seeding feature: %s\n", feature)
		for _, seed := range seeds {
			applied, err := s.JSONSeeder.SeedApplied(seed.Datetime, seed.Name, feature)
			if err != nil {
				return fmt.Errorf("failed to check if seed was applied: %w", err)
			}
			if applied {
				s.Log().Debugf("Seed already applied: %s-%s [%s]", seed.Datetime, seed.Name, feature)
				continue
			}

			var data SeedData
			err = json.Unmarshal([]byte(seed.Content), &data)
			if err != nil {
				return fmt.Errorf("failed to unmarshal %s seed: %w", feature, err)
			}

			err = s.seedData(ctx, &data)
			if err != nil {
				return err
			}

			err = s.JSONSeeder.ApplyJSONSeed(seed.Datetime, seed.Name, feature, seed.Content)
			if err != nil {
				s.Log().Errorf("error recording JSON seed: %v", err)
			}
		}
	}
	return nil
}

func (s *Seeder) seedData(ctx context.Context, data *SeedData) error {
	ctx, tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("error at beginning tx for seedLayouts: %w", err)
	}
	defer tx.Rollback()

	// Seed layouts and build ref->UUID map
	layoutRefToID := make(map[string]uuid.UUID)
	for i := range data.Layouts {
		l := &data.Layouts[i]
		l.GenCreateValues()
		err := s.repo.CreateLayout(ctx, *l)
		if err != nil {
			return fmt.Errorf("error inserting layout: %w", err)
		}
		if l.Name != "" {
			layoutRefToID[l.Name] = l.ID()
		}
		if l.RefValue != "" {
			layoutRefToID[l.RefValue] = l.ID()
		}
	}

	// Seed sections, resolving layout_ref
	for i := range data.Sections {
		sec := &data.Sections[i]
		sec.GenCreateValues()
		if sec.LayoutID == uuid.Nil && sec.Path == "/" && sec.Name == "root" {
			// Try to resolve layout_ref if present
			if ref, ok := any(sec).(interface{ LayoutRef() string }); ok {
				layoutRef := ref.LayoutRef()
				if layoutRef != "" {
					if id, found := layoutRefToID[layoutRef]; found {
						sec.LayoutID = id
					}
				}
			}
			// Fallback: try to resolve from a custom field if present
			if sec.LayoutID == uuid.Nil {
				if id, found := layoutRefToID["alt"]; found {
					sec.LayoutID = id
				}
			}
		}
		err := s.repo.CreateSection(ctx, *sec)
		if err != nil {
			return fmt.Errorf("error inserting section: %w", err)
		}
	}

	return tx.Commit()
}

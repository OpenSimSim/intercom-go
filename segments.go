package intercom

import (
	"context"
	"fmt"
)

// SegmentService handles interactions with the API through a SegmentRepository.
type SegmentService struct {
	Repository SegmentRepository
}

// Segment represents an Segment in Intercom.
type Segment struct {
	ID         string `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	CreatedAt  int64  `json:"created_at,omitempty"`
	UpdatedAt  int64  `json:"updated_at,omitempty"`
	PersonType string `json:"person_type,omitempty"`
}

// SegmentList, an object holding a list of Segments
type SegmentList struct {
	Segments []Segment `json:"segments,omitempty"`
}

// List all Segments for the App
func (t *SegmentService) List(ctx context.Context) (SegmentList, error) {
	return t.Repository.list(ctx)
}

// Find a particular Segment in the App
func (t *SegmentService) Find(ctx context.Context, id string) (Segment, error) {
	return t.Repository.find(ctx, id)
}

func (s Segment) String() string {
	return fmt.Sprintf("[intercom] segment { id: %s, type: %s }", s.ID, s.PersonType)
}

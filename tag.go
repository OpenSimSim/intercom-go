package intercom

import (
	"context"
	"fmt"
)

// TagService handles interactions with the API through a TagRepository.
type TagService struct {
	Repository TagRepository
}

// Tag represents an Tag in Intercom.
type Tag struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// TagList, an object holding a list of Tags
type TagList struct {
	Tags []Tag `json:"tags,omitempty"`
}

// List all Tags for the App
func (t *TagService) List(ctx context.Context) (TagList, error) {
	return t.Repository.list(ctx)
}

// Save a new Tag for the App.
func (t *TagService) Save(ctx context.Context, tag *Tag) (Tag, error) {
	return t.Repository.save(ctx, tag)
}

// Delete a Tag
func (t *TagService) Delete(ctx context.Context, id string) error {
	return t.Repository.delete(ctx, id)
}

// Tag Users or Companies using a TaggingList.
func (t *TagService) Tag(ctx context.Context, taggingList *TaggingList) (Tag, error) {
	return t.Repository.tag(ctx, taggingList)
}

func (t Tag) String() string {
	return fmt.Sprintf("[intercom] tag { id: %s name: %s }", t.ID, t.Name)
}

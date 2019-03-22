package intercom

import (
	"context"
	"fmt"
)

// ContactService handles interactions with the API through a ContactRepository.
type ContactService struct {
	Repository ContactRepository
}

// ContactList holds a list of Contacts and paging information
type ContactList struct {
	Pages       PageParams
	Contacts    []Contact
	ScrollParam string `json:"scroll_param,omitempty"`
}

// Contact represents a Contact within Intercom.
// Not all of the fields are writeable to the API, non-writeable fields are
// stripped out from the request. Please see the API documentation for details.
type Contact struct {
	ID                     string                 `json:"id,omitempty"`
	Email                  string                 `json:"email,omitempty"`
	Phone                  string                 `json:"phone,omitempty"`
	UserID                 string                 `json:"user_id,omitempty"`
	Name                   string                 `json:"name,omitempty"`
	Avatar                 *UserAvatar            `json:"avatar,omitempty"`
	LocationData           *LocationData          `json:"location_data,omitempty"`
	LastRequestAt          int64                  `json:"last_request_at,omitempty"`
	CreatedAt              int64                  `json:"created_at,omitempty"`
	UpdatedAt              int64                  `json:"updated_at,omitempty"`
	SessionCount           int64                  `json:"session_count,omitempty"`
	LastSeenIP             string                 `json:"last_seen_ip,omitempty"`
	SocialProfiles         *SocialProfileList     `json:"social_profiles,omitempty"`
	UnsubscribedFromEmails *bool                  `json:"unsubscribed_from_emails,omitempty"`
	UserAgentData          string                 `json:"user_agent_data,omitempty"`
	Tags                   *TagList               `json:"tags,omitempty"`
	Segments               *SegmentList           `json:"segments,omitempty"`
	Companies              *CompanyList           `json:"companies,omitempty"`
	CustomAttributes       map[string]interface{} `json:"custom_attributes,omitempty"`
	UpdateLastRequestAt    *bool                  `json:"update_last_request_at,omitempty"`
	NewSession             *bool                  `json:"new_session,omitempty"`
}

type contactListParams struct {
	PageParams
	SegmentID string `url:"segment_id,omitempty"`
	TagID     string `url:"tag_id,omitempty"`
	Email     string `url:"email,omitempty"`
}

// FindByID looks up a Contact by their Intercom ID.
func (c *ContactService) FindByID(ctx context.Context, id string) (Contact, error) {
	return c.findWithIdentifiers(ctx, UserIdentifiers{ID: id})
}

// FindByUserID looks up a Contact by their UserID (automatically generated server side).
func (c *ContactService) FindByUserID(ctx context.Context, userID string) (Contact, error) {
	return c.findWithIdentifiers(ctx, UserIdentifiers{UserID: userID})
}

func (c *ContactService) findWithIdentifiers(ctx context.Context, identifiers UserIdentifiers) (Contact, error) {
	return c.Repository.find(ctx, identifiers)
}

// List all Contacts for App.
func (c *ContactService) List(ctx context.Context, params PageParams) (ContactList, error) {
	return c.Repository.list(ctx, contactListParams{PageParams: params})
}

// List all Contacts for App via Scroll API
func (c *ContactService) Scroll(ctx context.Context, scrollParam string) (ContactList, error) {
	return c.Repository.scroll(ctx, scrollParam)
}

// ListByEmail looks up a list of Contacts by their Email.
func (c *ContactService) ListByEmail(ctx context.Context, email string, params PageParams) (ContactList, error) {
	return c.Repository.list(ctx, contactListParams{PageParams: params, Email: email})
}

// List Contacts by Segment.
func (c *ContactService) ListBySegment(ctx context.Context, segmentID string, params PageParams) (ContactList, error) {
	return c.Repository.list(ctx, contactListParams{PageParams: params, SegmentID: segmentID})
}

// List Contacts By Tag.
func (c *ContactService) ListByTag(ctx context.Context, tagID string, params PageParams) (ContactList, error) {
	return c.Repository.list(ctx, contactListParams{PageParams: params, TagID: tagID})
}

// Create Contact
func (c *ContactService) Create(ctx context.Context, contact *Contact) (Contact, error) {
	return c.Repository.create(ctx, contact)
}

// Update Contact
func (c *ContactService) Update(ctx context.Context, contact *Contact) (Contact, error) {
	return c.Repository.update(ctx, contact)
}

// Convert Contact to User
func (c *ContactService) Convert(ctx context.Context, contact *Contact, user *User) (User, error) {
	return c.Repository.convert(ctx, contact, user)
}

// Delete Contact
func (c *ContactService) Delete(ctx context.Context, contact *Contact) (Contact, error) {
	return c.Repository.delete(ctx, contact.ID)
}

// MessageAddress gets the address for a Contact in order to message them
func (c Contact) MessageAddress() MessageAddress {
	return MessageAddress{
		Type:   "contact",
		ID:     c.ID,
		Email:  c.Email,
		UserID: c.UserID,
	}
}

func (c Contact) String() string {
	return fmt.Sprintf("[intercom] contact { id: %s name: %s, user_id: %s, email: %s }", c.ID, c.Name, c.UserID, c.Email)
}

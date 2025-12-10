package contacts

// PhoneDTO represents a phone number data transfer object. Used in both requests and responses.
// It matches child table contact_phone_numbers (numbers+label).

type PhoneDTO struct {
	UserID           string `json:"user_id,omitempty"`           // contact_phone_numbers.user_id
	ContactID        string `json:"contact_id,omitempty"`        // contact_phone_numbers.contact_id
	Label            string `json:"label"`                       // contact_phone_numbers.label
	Number           string `json:"number"`                      // contact_phone_numbers.number
	NormalizedNumber string `json:"normalized_number,omitempty"` // contact_phone_numbers.normalized_number
	IsPrimary        bool   `json:"is_primary"`                  // contact_phone_numbers.is_primary
}

// EmailDTO matches contact_emails (emails+label).
type EmailDTO struct {
	UserID          string `json:"user_id"`          // contact_emails.user_id
	ContactID       string `json:"contact_id"`       // contact_emails.contact_id
	Label           string `json:"label"`            // contact_emails.label
	Email           string `json:"email"`            // contact_emails.email
	NormalizedEmail string `json:"normalized_email"` // contact_emails.normalized_email
	IsPrimary       bool   `json:"is_primary"`       // contact_emails.is_primary
}

// CreateContactRequest is what the frontend sends to POST /contacts
// Fields map 1:1 to columns in the contacts table + child table collections
type CreateContactRequest struct {
	DisplayName string     `json:"display_name"`     // contacts.display_name
	FirstName   string     `json:"first_name"`       // contacts.first_name
	Surname     string     `json:"surname"`          // contacts.surname
	Note        string     `json:"note,omitempty"`   // contacts.note (optional)
	Source      string     `json:"source"`           // contacts.source
	Phones      []PhoneDTO `json:"phones,omitempty"` // contact_phone_numbers (optional)
	Emails      []EmailDTO `json:"emails,omitempty"` // contact_emails (optional)
}

// ContactResponse is what we return to the client after creation/fetch.
// It wraps the main contacts row + its phones/emails.
type ContactResponse struct {
	ID          string     `json:"id"` // contacts.id (uuid)
	DisplayName string     `json:"displayName"`
	FirstName   string     `json:"firstName"`
	Surname     string     `json:"surname"`
	Note        string     `json:"note,omitempty"`
	Source      string     `json:"source"`
	Phones      []PhoneDTO `json:"phones,omitempty"`
	Emails      []EmailDTO `json:"emails,omitempty"`
}

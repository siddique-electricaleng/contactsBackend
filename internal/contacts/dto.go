package contacts

// PhoneDTO represents a phone number data transfer object. Used in both requests and responses.
// It matches child table contact_phone_numbers (numbers+label).

type PhoneDTO struct {
	UserID           string `json:"user_id,omitempty" example:"0a0915ca-f33b-4140-b11b-c7a6cf9e2acf"`    // contact_phone_numbers.user_id
	ContactID        string `json:"contact_id,omitempty" example:"95a6d5ac-5e28-43e9-aa30-81021e5f288e"` // contact_phone_numbers.contact_id
	Label            string `json:"label" example:"mobile"`                                              // contact_phone_numbers.label
	Number           string `json:"number" example:"+8801681602515"`                                     // contact_phone_numbers.number
	NormalizedNumber string `json:"normalized_number,omitempty" example:"8801681602515"`                 // contact_phone_numbers.normalized_number
	IsPrimary        bool   `json:"is_primary" example:"true"`                                           // contact_phone_numbers.is_primary
}

// EmailDTO matches contact_emails (emails+label).
type EmailDTO struct {
	UserID          string `json:"user_id,omitempty" example:"0a0915ca-f33b-4140-b11b-c7a6cf9e2acf"`    // contact_emails.user_id
	ContactID       string `json:"contact_id,omitempty" example:"95a6d5ac-5e28-43e9-aa30-81021e5f288e"` // contact_emails.contact_id
	Label           string `json:"label" example:"work"`                                                // contact_emails.label
	Email           string `json:"email" example:"abu.bakr@fiberathome.net"`                            // contact_emails.email
	NormalizedEmail string `json:"normalized_email,omitempty" example:"abu.bakr@fiberathome.net"`       // contact_emails.normalized_email
	IsPrimary       bool   `json:"is_primary" example:"false"`                                          // contact_emails.is_primary
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
	ID          string     `json:"id" example:"95a6d5ac-5e28-43e9-aa30-81021e5f288e"` // contacts.id (uuid)
	DisplayName string     `json:"displayName" example:"Abu Bakr Siddique"`
	FirstName   string     `json:"firstName" example:"Abu Bakr"`
	Surname     string     `json:"surname" example:"Siddique"`
	Note        string     `json:"note,omitempty" example:"Programmer from FGL"`
	Source      string     `json:"source" example:"DB_STORED"`
	Phones      []PhoneDTO `json:"phones,omitempty"`
	Emails      []EmailDTO `json:"emails,omitempty"`
}

type ContactCreateData struct {
	Code        int                   `json:"code" example:"200"`
	Count       int                   `json:"count" example:"1"`
	ContactInfo ContactIDWithDispName `json:"contactInfo"`
}
type CreateContactCountResponse struct {
	Message string `json:"message" example:"1 contact created and stored"`
	Data    ContactCreateData
}

// For showing relevant display_name and contact_id per contact

type ContactIDWithDispName struct {
	ContactID   string `json:"contactID" example:"c16dfb9c-6e26-48c6-acbd-44dd3a7827d1"`
	DisplayName string `json:"displayName" example:"Abu Bakr Siddique"`
}

// For Bulk Upload
type BulkCreateContactRequest struct {
	Contacts []CreateContactRequest `json:"contacts"`
}

type BulkContactCreateData struct {
	Code        int                     `json:"code" example:"200"`
	Count       int                     `json:"count" example:"3"`
	ContactInfo []ContactIDWithDispName `json:"contactInfo"`
}

type CreateBulkContactCountResponse struct {
	Message string `json:"message" example:"3 contacts created and stored"`
	Data    BulkContactCreateData
}

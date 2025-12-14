package contacts

// These structs are to be used for Swagger Docs only and not necessarily inside the code

type CodeOnly struct {
	Code int `json:"code" example:"200"`
}

// --------------------- READ / GET contacts ------------------------------
type ContactsData struct {
	CodeOnly
	Contacts []ContactResponse `json:"contacts"`
}

type ListContactsResponse struct {
	Message string       `json:"message" example:"contacts fetched successfully"`
	Data    ContactsData `json:"data"`
}

// --------------------- CREATE / POST contacts ------------------------------
// type ContactData struct {
// 	CodeOnly
// 	Contact ContactResponse `json:"contact"`
// }

// type CreateContactResponse struct {
// 	Message string      `json:"message" example:"contact created and stored"`
// 	Data    ContactData `json:"data"`
// }

// ------------------ HTTP Error Envelope ------------------------------------
type ErrorEnvelope struct {
	Message string `example:"unauthorized"`
}

// ------------------- Contact Payload for POST --------------------------------
type PhoneCreateRequest struct {
	Label     string `json:"label" example:"mobile"`
	Number    string `json:"number" example:"+8801681602515"`
	IsPrimary bool   `json:"is_primary" example:"true"`
}

type EmailCreateRequest struct {
	Label     string `json:"label" example:"work"`
	Email     string `json:"email" example:"abu.bakr@fiberathome.net"`
	IsPrimary bool   `json:"is_primary" example:"true"`
}

type CreateContactRequestDocs struct {
	DisplayName string               `json:"display_name" example:"Abu Bakr Siddique"`
	FirstName   string               `json:"first_name" example:"Abu Bakr"`
	Surname     string               `json:"surname" example:"Siddique"`
	Note        string               `json:"note" example:"Programmer from FGL"`
	Source      string               `json:"source" example:"DB_STORED"`
	Phones      []PhoneCreateRequest `json:"phones"`
	Emails      []EmailCreateRequest `json:"emails"`
}

// Bulk upload docs

// BulkCreateContactRequestDocs is used only in Swagger annotations.
type BulkCreateContactRequestDocs struct {
	Contacts []CreateContactRequestDocs `json:"contacts"`
}

// Common "code + contacts" wrapper for responses.
// type BulkCreateContactsData struct {
// 	Code     int               `json:"code"     example:"201"`
// 	Contacts []ContactResponse `json:"contacts"`
// }

// // Full response body for /contacts/sync
// type BulkCreateContactsResponse struct {
// 	Message string                 `json:"message" example:"contacts created successfully"`
// 	Data    BulkCreateContactsData `json:"data"`
// }

// New full body response for /contacts/sync
type BulkCountResponse struct {
	Message string `json:"message" example:"3 contacts created successfully"`
	Data    struct {
		Code  int `json:"code" example:"200"`
		Count int `json:"count" example:"3"`
	} `json:"data"`
}

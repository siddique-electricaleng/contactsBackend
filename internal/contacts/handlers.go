package contacts

import (
	"contacts/internal/json"
	appMiddleWare "contacts/internal/middleware"
	"contacts/internal/transport"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

var (
	ErrInvalidLength = errors.New("invalid userID obtained")
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

// // CreateContacts godoc
// @Summary      Create/Update - Upsert a Contact with required phone and email fields
// @Description  After verifying if user is authorized using bearer token, a contact can be created or existing contact can be updated and stored in Pico Database through this API, additionally a single contact can have multiple phone numbers and emails
// @Tags         contacts
// @Accept       json
// @Produce      json
// @Security	BearerAuth
// @Param        body  body      CreateContactRequestDocs  true  "Contact payload"
// @Success      201    {object}  CreateContactCountResponse
// @Failure      400    {string}  string  "invalid body"
// @Failure      401    {string}  string  "unauthorized"
// @Failure      409    {string}  string  "email or username already exists"
// @Failure      500    {string}  string  "err.Error() messages"
// @Router       /contacts [post]
func (h *Handler) CreateContacts(w http.ResponseWriter, r *http.Request) {
	// Decode json body
	var reqBody CreateContactRequest

	if err := json.Read(r, &reqBody); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	// Get userID string from auth middleware context
	userIDString, ok := appMiddleWare.UserIDFromContext(r.Context())

	if !ok || userIDString == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	userIDUUID, err := uuid.Parse(userIDString)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	}
	// Calling the service
	contactResp, contactInfo, err := h.service.CreateContacts(r.Context(), userIDUUID, reqBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Wrap response in transport envelope - helper
	/*
		resp := transport.Envelope{
			Message: fmt.Sprintf("%d contact created and stored", contactResp),
			Data: map[string]any{
				"code":    http.StatusOK,
				"contact": contactResp,
			},
		}
	*/
	resp := CreateContactCountResponse{
		Message: fmt.Sprintf("%d contact created and stored", contactResp),
		Data: ContactCreateData{
			Code:        http.StatusOK,
			Count:       contactResp,
			ContactInfo: contactInfo,
		},
	}

	json.Write(w, http.StatusCreated, resp)
}

// BulkCreateContacts godoc
// @Summary     Bulk create contacts
// @Description Create multiple contacts for the authenticated user in a single request.
// @Tags        contacts
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       body body BulkCreateContactRequestDocs true "Bulk contact payload"
// @Success     201 {object} CreateBulkContactCountResponse
// @Failure     400 {string} string "invalid body"
// @Failure     401 {string} string "unauthorized"
// @Failure     500 {string} string "err.Error() messages"
// @Router      /contacts/sync [post]
func (h *Handler) BulkCreateContacts(w http.ResponseWriter, r *http.Request) {
	userIDString, ok := appMiddleWare.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := uuid.Parse(userIDString)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Decode JSON body
	var req BulkCreateContactRequest
	if err := json.Read(r, &req); err != nil {
		http.Error(w, "invalid body "+err.Error(), http.StatusBadRequest)
		return
	}
	if len(req.Contacts) == 0 {
		http.Error(w, "contacts array cannot be empty", http.StatusBadRequest)
		return
	}

	total, infos, err := h.service.BulkCreateContacts(r.Context(), userID, req.Contacts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	/*
		resp := BulkCreateContactsResponse{
			Message: "contacts created successfully",
			Data: BulkCreateContactsData{
				Code:     http.StatusCreated,
				Contacts: created,
			},
		}
	*/
	// resp := BulkCountResponse{
	// 	Message: fmt.Sprintf("%d contacts created successfully", created),
	// }
	// resp.Data.Code = http.StatusOK
	// resp.Data.Count = created

	resp := CreateBulkContactCountResponse{
		Message: fmt.Sprintf("%d contacts created successfully", total),
		Data: BulkContactCreateData{
			Code:        http.StatusOK,
			Count:       total,
			ContactInfo: infos,
		},
	}

	json.Write(w, http.StatusCreated, resp)
}

// Get/Fetch Contacts godoc
// @Summary     Get all the contacts stored in Pico for the specific user
// @Description After verifying if user is authorized using bearer token, fetch/get all the contacts stored for the specific user.
// @Tags        contacts
// @Produce     json
// @Security	BearerAuth
// @Success     200 {object} ListContactsResponse
// @Failure     401 {string} string "unauthorized"
// @Failure     500 {string} string "err.Error() messages"
// @Router      /contacts [get]
func (h *Handler) ListContactsForUserWithDetails(w http.ResponseWriter, r *http.Request) {

	userIDString, ok := appMiddleWare.UserIDFromContext(r.Context())

	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := uuid.Parse(userIDString)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	contacts, err := h.service.ListContactsForUserWithDetails(r.Context(), userID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := transport.Envelope{
		Message: "contacts fetched successfully",
		Data: map[string]any{
			"code":     http.StatusOK,
			"contacts": contacts,
		},
	}

	json.Write(w, http.StatusOK, response)
}

package contacts

import (
	"contacts/internal/json"
	appMiddleWare "contacts/internal/middleware"
	"contacts/internal/transport"
	"errors"
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

func (h *Handler) CreateContacts(w http.ResponseWriter, r *http.Request) {
	// Decode json body
	var reqBody CreateContactRequest

	if err := json.Read(r, &reqBody); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	// Get userID from auth middleware context
	userIDString, ok := appMiddleWare.UserIDFromContext(r.Context())

	if !ok || userIDString == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	userIDUUID, err := uuid.Parse(userIDString)
	if err != nil {
		http.Error(w, "error converting string userID to UUID", http.StatusInternalServerError)
	}
	// Calling the service
	contactResp, err := h.service.CreateContacts(r.Context(), userIDUUID, reqBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Wrap response in transport envelope - helper
	resp := transport.Envelope{
		Message: "contact created and stored",
		Data: map[string]any{
			"code":    http.StatusCreated,
			"contact": contactResp,
		},
	}

	json.Write(w, http.StatusCreated, resp)
}

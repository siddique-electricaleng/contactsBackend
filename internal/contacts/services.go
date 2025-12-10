package contacts

import (
	repo "contacts/internal/adapters/postgresql/sqlc"
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service interface {
	CreateContacts(ctx context.Context, userID uuid.UUID, req CreateContactRequest) (ContactResponse, error)
}

type svc struct {
	query repo.Querier
}

func NewService(r *repo.Queries) Service {
	return &svc{query: r}
}

/* Internal Helpers
1. Normalize  Email
2. Normalize Phone Numbers
3. Compile first_name + surname to make display_name

*/

func normalizeEmail(email string) string {
	return strings.TrimSpace(strings.ToLower(email))
}

func normalizePhone(number string) (normalizedNum string) {

	// Remove +, spaces and dashes
	normalizedNum = strings.TrimSpace(number)
	normalizedNum = strings.ReplaceAll(normalizedNum, " ", "")
	normalizedNum = strings.ReplaceAll(normalizedNum, "-", "")
	normalizedNum = strings.ReplaceAll(normalizedNum, "+", "")

	return
}

func joinName(firstName string, surname string) string {
	return fmt.Sprintf("%s %s", firstName, surname)
}

func (s *svc) CreateContacts(ctx context.Context, userID uuid.UUID, req CreateContactRequest) (ContactResponse, error) {

	// 1. Join first_name and surname to form display_name

	displayName := joinName(req.FirstName, req.Surname)

	// 2a. Insert data into the contacts table
	insertContact, err := s.query.CreateContact(ctx, repo.CreateContactParams{
		UserID: pgtype.UUID{
			Bytes: userID,
			Valid: true,
		},
		DisplayName: displayName,
		FirstName: pgtype.Text{
			String: req.FirstName,
			Valid:  true,
		},
		Surname: pgtype.Text{
			String: req.Surname,
			Valid:  true,
		},
		Note: pgtype.Text{
			String: req.Note,
			Valid:  true,
		},
		Source: req.Source,
	})

	// 2b. check for errors upon data insertion into contacts table
	if err != nil {
		return ContactResponse{}, fmt.Errorf("error Inserting/Updating contact information")
	}

	// 3. Upsert Phones
	phoneDTOs := []PhoneDTO{}

	for _, phone := range req.Phones {
		normalized := normalizePhone(phone.Number)

		insertContactNumber, err := s.query.UpsertContactPhone(ctx, repo.UpsertContactPhoneParams{
			UserID: pgtype.UUID{
				Bytes: userID,
				Valid: true,
			},
			ContactID: insertContact.ID,
			Number:    phone.Number,
			Label: pgtype.Text{
				String: phone.Label,
				Valid:  true,
			},
			NormalizedNumber: normalized,
			IsPrimary:        phone.IsPrimary,
		})

		if err != nil {
			return ContactResponse{}, fmt.Errorf("error Inserting/Updating number")
		}

		phoneDTOs = append(phoneDTOs, PhoneDTO{
			UserID:           insertContact.UserID.String(),
			ContactID:        insertContact.ID.String(),
			Label:            insertContactNumber.Label.String,
			Number:           insertContactNumber.Number,
			NormalizedNumber: insertContactNumber.NormalizedNumber,
			IsPrimary:        insertContactNumber.IsPrimary,
		})
	}

	// 4. Upsert Email
	emailDTOs := []EmailDTO{}

	for _, email := range req.Emails {
		normalized := normalizeEmail(email.Email)

		insertContactEmail, err := s.query.UpsertContactEmail(ctx, repo.UpsertContactEmailParams{
			UserID: pgtype.UUID{
				Bytes: userID,
				Valid: true,
			},
			ContactID: insertContact.ID,
			Email:     email.Email,
			Label: pgtype.Text{
				String: email.Label,
				Valid:  true,
			},
			NormalizedEmail: normalized,
			IsPrimary:       email.IsPrimary,
		})

		if err != nil {
			return ContactResponse{}, fmt.Errorf("error Inserting/Updating email")
		}

		emailDTOs = append(emailDTOs, EmailDTO{
			UserID:          insertContact.UserID.String(),
			ContactID:       insertContact.ID.String(),
			Label:           insertContactEmail.Label.String,
			Email:           insertContactEmail.Email,
			NormalizedEmail: insertContactEmail.NormalizedEmail,
			IsPrimary:       insertContactEmail.IsPrimary,
		})
	}

	// 5. Build ContactResponseDTO
	resp := ContactResponse{
		ID:          insertContact.ID.String(),
		DisplayName: insertContact.DisplayName,
		FirstName:   insertContact.FirstName.String,
		Surname:     insertContact.Surname.String,
		Note:        insertContact.Note.String,
		Source:      insertContact.Source,
		Phones:      phoneDTOs,
		Emails:      emailDTOs,
	}

	return resp, nil
}

// func (s *svc) ListContactsForUser(ctx context.Context, userID string, limit, offset int32)([]ContactDTO)

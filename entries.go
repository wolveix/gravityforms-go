package gravityforms

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

type Entry struct {
	ID              int               `json:"-"`
	FormID          string            `json:"form_id,omitempty"`
	PostID          string            `json:"post_id,omitempty"`
	DateCreated     string            `json:"date_created,omitempty"`
	DateUpdated     string            `json:"date_updated,omitempty"`
	IsFulfilled     string            `json:"is_fulfilled,omitempty"`
	IsStarred       string            `json:"is_starred,omitempty"`
	IsRead          string            `json:"is_read,omitempty"`
	IP              string            `json:"ip,omitempty"`
	SourceURL       string            `json:"source_url,omitempty"`
	UserAgent       string            `json:"user_agent,omitempty"`
	Currency        string            `json:"currency,omitempty"`
	CreatedBy       string            `json:"created_by,omitempty"`
	Status          string            `json:"status,omitempty"`
	PaymentAmount   string            `json:"payment_amount,omitempty"`
	PaymentDate     string            `json:"payment_date,omitempty"`
	PaymentStatus   string            `json:"payment_status,omitempty"`
	TransactionID   string            `json:"transaction_id,omitempty"`
	TransactionType string            `json:"transaction_type,omitempty"`
	Fields          map[string]string `json:"-"`
}

func (e *Entry) GetField(id string) string {
	if _, ok := e.Fields[id]; ok {
		return e.Fields[id]
	}

	return ""
}

func (e *Entry) SetField(id string, value string) {
	if e.Fields == nil {
		e.Fields = make(map[string]string)
	}

	e.Fields[id] = value
}

// MarshalJSON custom function to handle predefined fields and dynamic fields
func (e *Entry) MarshalJSON() ([]byte, error) {
	jsonData, err := json.Marshal(*e)
	if err != nil {
		return nil, err
	}

	preparedData := make(map[string]string)
	preparedData["id"] = strconv.Itoa(e.ID)

	if err = json.Unmarshal(jsonData, &preparedData); err != nil {
		return nil, err
	}

	for key, value := range e.Fields {
		preparedData[key] = value
	}

	return json.Marshal(preparedData)
}

// UnmarshalJSON handles unmarshalling the known field IDs into strict fields, and then puts the rest into Fields.
func (e *Entry) UnmarshalJSON(data []byte) error {
	type Alias Entry

	alias := Alias{
		Fields: make(map[string]string),
	}

	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}

	var fields map[string]interface{}
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}

	for key, value := range fields {
		switch key {
		case "id":
			// Handle ID field specifically.
			if str, ok := value.(string); ok {
				id, err := strconv.ParseInt(str, 10, 64)
				if err == nil {
					alias.ID = int(id)
				}
			} else if num, ok := value.(float64); ok {
				alias.ID = int(num)
			}
		case "form_id", "post_id", "date_created", "date_updated", "is_fulfilled", "is_starred",
			"is_read", "ip", "source_url", "user_agent", "currency", "created_by", "status",
			"payment_amount", "payment_date", "payment_status", "transaction_id", "transaction_type":
			// Skip known fields.
			continue
		default:
			// Convert value to string based on type
			switch v := value.(type) {
			case string:
				alias.Fields[key] = v
			case bool:
				alias.Fields[key] = strconv.FormatBool(v)
			case float64:
				alias.Fields[key] = strconv.FormatFloat(v, 'f', -1, 64)
			case int:
				alias.Fields[key] = strconv.Itoa(v)
			case nil:
				alias.Fields[key] = ""
			default:
				// For any other type, try to JSON marshal it to string.
				if b, err := json.Marshal(v); err == nil {
					alias.Fields[key] = string(b)
				} else {
					alias.Fields[key] = fmt.Sprintf("%v", v)
				}
			}
		}
	}

	*e = Entry(alias)
	return nil
}

// CreateEntry creates the given entry within Gravity Forms.
func (s *Service) CreateEntry(formID int, entry *Entry) error {
	entry.FormID = strconv.Itoa(formID)

	response := struct {
		ID int `json:"id"`
	}{}

	if _, err := s.makeRequest(http.MethodPost, "entries", &entry, &response); err != nil {
		return err
	}

	entry.ID = response.ID

	return nil
}

// DeleteEntry deletes the given entry within Gravity Forms
func (s *Service) DeleteEntry(id int) error {
	if id == 0 {
		return errors.New("missing entry id")
	}

	if _, err := s.makeRequest(http.MethodDelete, "entries/"+strconv.Itoa(id), nil, nil); err != nil {
		return err
	}

	return nil
}

// GetEntries called GetEntriesByFormID with 0 as the given form ID.
func (s *Service) GetEntries() ([]*Entry, error) {
	return s.GetEntriesByFormID(0)
}

// GetEntriesByFormID returns all entries for the given form ID, utilizing Gravity Forms' pagination.
func (s *Service) GetEntriesByFormID(formID int) ([]*Entry, error) {
	var entries []*Entry
	currentPage := 1
	pageSize := 100

	for {
		obj := struct {
			Entries    []*Entry `json:"entries"`
			TotalCount int      `json:"total_count"`
		}{}

		if _, err := s.makeRequest(http.MethodGet, fmt.Sprintf("forms/%d/entries?paging[page_size]=%d&paging[current_page]=%d", formID, pageSize, currentPage), nil, &obj); err != nil {
			return nil, err
		}

		entries = append(entries, obj.Entries...)

		if obj.TotalCount <= (currentPage * pageSize) {
			break
		}

		currentPage++
	}

	if len(entries) == 0 {
		return nil, errors.New("no entries found")
	}

	return entries, nil
}

func (s *Service) GetEntryByID(id int) (*Entry, error) {
	if id == 0 {
		return nil, errors.New("missing entry id")
	}

	var entry *Entry

	if _, err := s.makeRequest(http.MethodGet, "entries/"+strconv.Itoa(id), nil, &entry); err != nil {
		return nil, err
	}

	return entry, nil
}

// UpdateEntry updates the given entry within Gravity Forms.
func (s *Service) UpdateEntry(id int, entry *Entry) error {
	if id == 0 {
		return errors.New("missing entry id")
	}

	if _, err := s.makeRequest(http.MethodPut, "entries/"+strconv.Itoa(id), entry, nil); err != nil {
		return err
	}

	return nil
}

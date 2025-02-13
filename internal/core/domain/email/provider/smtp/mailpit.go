package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Mailpit struct {
	instanceURL string
}

type EmailRequest struct {
	Subject string `json:"Subject"`
	Text    string `json:"Text"`
	To      []To   `json:"To"`
}

type From struct {
	Email string `json:"Email"`
	Name  string `json:"Name"`
}

type To struct {
	Email string `json:"Email"`
	Name  string `json:"Name"`
}

func (m *Mailpit) Send(r Request) error {

	requestBody := EmailRequest{
		Subject: r.Subject,
		Text:    r.Text,
		To: []To{
			{Email: r.To, Name: ""},
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %v", err)
	}

	req, err := http.NewRequest("POST", m.instanceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// return fmt.Errorf("error: received status code %d", resp.StatusCode)
		return ErrSysUnknown
	}

	return nil
}

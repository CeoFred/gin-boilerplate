package sendgrid

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Email struct {
	Personalizations []Personalization `json:"personalizations"`
	From             EmailAddress      `json:"from"`
	Content          []Content         `json:"content"`
}

type Personalization struct {
	To      []EmailAddress `json:"to"`
	Subject string         `json:"subject"`
}

type EmailAddress struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

type Content struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Client struct {
	apiKey string
}

func NewClient(apiKey string) *Client {

	return &Client{
		apiKey: apiKey,
	}
}

func (c *Client) Send(person *EmailAddress, from, fromName, subject, content string) error {

	body := &Email{
		Personalizations: []Personalization{
			{
				To: []EmailAddress{
					{
						Email: person.Email,
						Name:  person.Name,
					},
				},
				Subject: subject,
			},
		},
		From: EmailAddress{
			Email: from,
			Name:  fromName,
		},
		Content: []Content{
			{
				Type:  "text/html",
				Value: content,
			},
		},
	}

	reqBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://api.sendgrid.com/v3/mail/send", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Convert the response body to a string for logging
	respBodyStr := string(respBody)

	if resp.StatusCode >= 203 {
		return fmt.Errorf("sendgrid: error sending email, status code: %d, message: %s ", resp.StatusCode, respBodyStr)
	}

	return nil
}

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
	apiKey  string
	From    string
	Subject string
	Content string
	Name    string
}

func NewClient(apiKey, from, fromName, subject, content string) *Client {

	return &Client{
		apiKey:  apiKey,
		From:    from,
		Subject: subject,
		Content: content,
		Name:    fromName,
	}
}

func (c *Client) Send(person *EmailAddress) error {

	body := &Email{
		Personalizations: []Personalization{
			{
				To: []EmailAddress{
					{
						Email: person.Email,
						Name:  person.Name,
					},
				},
				Subject: c.Subject,
			},
		},
		From: EmailAddress{
			Email: c.From,
			Name:  c.Name,
		},
		Content: []Content{
			{
				Type:  "text/html",
				Value: c.Content,
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

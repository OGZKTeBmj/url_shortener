package dto

import (
	"fmt"
	"net/url"

	"github.com/OGZKTeBmj/url_shortener/internal/domain"
)

type ShortInput struct {
	URL     string `json:"url"`
	UserID  domain.UUID
	IP      string
	IsGuest bool
}

func (s *ShortInput) Validate() error {
	u, err := url.ParseRequestURI(s.URL)
	if err != nil {
		return fmt.Errorf("invalid url")
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("support only http(s)")
	}

	return nil
}

type UserInput struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (s *UserInput) Validate() error {
	if len(s.Name) <= 3 {
		return fmt.Errorf("name length must be more than 3")
	}
	if len(s.Password) <= 8 {
		return fmt.Errorf("name length must be more than 8")
	}
	return nil
}

type RefreshInput struct {
	Token string `json:"refresh_token"`
}

type UserURLOutput struct {
	Short  string `json:"short"`
	URL    string `json:"url"`
	Visits int64  `json:"visits"`
}

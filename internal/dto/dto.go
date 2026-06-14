package dto

import (
	"fmt"
)

type ShortInput struct {
	URL     string `json:"url"`
	IsGuest bool
}

func (s *ShortInput) Validate() error {
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

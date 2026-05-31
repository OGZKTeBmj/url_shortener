package dto

type ShortInput struct {
	URL     string `json:"url"`
	IsGuest bool
}

func (s *ShortInput) Validate() error {
	return nil
}

package dto

type ShortInput struct {
	URL string `json:"url"`
}

func (s *ShortInput) Validate() error {
	return nil
}

package domain

import "github.com/google/uuid"

type UUID [16]byte

func (u UUID) String() string {
	return uuid.UUID(u).String()
}

type User struct {
	UUID     UUID
	Name     string
	PassHash []byte
}

type RefreshSession struct {
	UserID UUID
}

type UserURL struct {
	URL    string
	Short  string
	Visits int64
}

package model

import "time"

type Customer struct {
	ID             int        `json:"id"`
	PhoneNumber    string     `json:"phoneNumber"`
	Email          string     `json:"email"`
	LinkedID       int        `json:"linkedId"`
	LinkPrecedence string     `json:"linkPrecedence"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
	DeletedAt      *time.Time `json:"deletedAt"`
}

type Result struct {
	PrimaryContactID    int      `json:"primaryContactId"`
	Emails              []string `json:"emails"`
	PhoneNumbers        []string `json:"phoneNumbers"`
	SecondaryContactIDs []int    `json:"secondaryContactIds"`
}

type Payload struct {
	PhoneNumber string `json:"phoneNumber"`
	Email       string `json:"email"`
}

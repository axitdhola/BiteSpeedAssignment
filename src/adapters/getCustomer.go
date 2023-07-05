package adapters

import (
	"time"

	"github.com/axitdhola/BiteSpeedAssignment/src/model"
)

func GetCustomerRequest() *model.Customer {
	return &model.Customer{
		ID:             -1,
		PhoneNumber:    "",
		Email:          "",
		LinkedID:       -1,
		LinkPrecedence: "",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		DeletedAt:      nil,
	}
}

func GetReults() *model.Result {
	return &model.Result{
		PrimaryContactID:    -1,
		Emails:              []string{},
		PhoneNumbers:        []string{},
		SecondaryContactIDs: []int{},
	}
}

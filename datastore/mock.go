package datastore

import (
	"errors"

	"github.com/customerio/homework/serve"
)

type Mock struct{}

var mockCustomer1 = &serve.Customer{
	ID: 1,
	Attributes: map[string]string{
		"email":  "customer1@example.com",
		"tier":   "S",
		"type":   "temporary",
		"animal": "tiger",
	},
	Events: map[string]int{
		"played_song": 5,
	},
	LastUpdated: 1625181700,
}

var mockCustomer2 = &serve.Customer{
	ID: 2,
	Attributes: map[string]string{
		"email":  "customer2@example.com",
		"tier":   "A",
		"type":   "permanent",
		"animal": "none",
	},
	Events: map[string]int{
		"played_song": 1,
	},
	LastUpdated: 1625180000,
}

func (d Mock) Get(id int) (*serve.Customer, error) {
	switch id {
	case 1:
		return mockCustomer1, nil
	case 2:
		return mockCustomer2, nil
	default:
		return nil, serve.ErrNotFound
	}
}

func (d Mock) List(page, count int) ([]*serve.Customer, error) {
	return []*serve.Customer{mockCustomer1, mockCustomer2}, nil
}

func (m Mock) Create(id int, attributes map[string]string) (*serve.Customer, error) {
	return nil, errors.New("unimplemented")
}

// Update is intentionally naive
func (m Mock) Update(id int, attributes map[string]string) (*serve.Customer, error) {
	switch id {
	case 1:
		mockCustomer1.Attributes = attributes
		return mockCustomer1, nil
	case 2:
		mockCustomer2.Attributes = attributes
		return mockCustomer2, nil
	default:
		return nil, serve.ErrNotFound
	}
}

func (m Mock) Delete(id int) error {
	return nil
}

func (m Mock) TotalCustomers() (int, error) {
	return 2, nil
}

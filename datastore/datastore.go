package datastore

import (
	"errors"

	"github.com/customerio/homework/serve"
)

type Datastore struct{}

func (d Datastore) Get(id int) (*serve.Customer, error) {
	return nil, errors.New("unimplemented")
}

func (d Datastore) List(page, count int) ([]*serve.Customer, error) {
	return nil, errors.New("unimplemented")
}

func (m Datastore) Create(id int, attributes map[string]string) (*serve.Customer, error) {
	return nil, errors.New("unimplemented")
}

func (m Datastore) Update(id int, attributes map[string]string) (*serve.Customer, error) {
	return nil, errors.New("unimplemented")
}

func (m Datastore) Delete(id int) error {
	return errors.New("unimplemented")
}

func (m Datastore) TotalCustomers() (int, error) {
	return 0, errors.New("unimplemented")
}

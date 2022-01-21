package serve

import "errors"

var ErrNotFound = errors.New("not found")

func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

type Event struct {
	Name  string
	Count int
}

type Customer struct {
	ID          int               `json:"id"`
	Attributes  map[string]string `json:"attributes"`
	Events      map[string]int    `json:"events"`
	LastUpdated int               `json:"last_updated"`
}

type Datastore interface {
	List(page, count int) ([]*Customer, error)
	Get(id int) (*Customer, error)
	Create(id int, attributes map[string]string) (*Customer, error)
	Update(id int, attributes map[string]string) (*Customer, error)
	Delete(id int) error
	TotalCustomers() (int, error)
}

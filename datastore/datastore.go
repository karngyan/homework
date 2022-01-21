package datastore

import (
	"strconv"
	"time"

	"github.com/customerio/homework/serve"
	"github.com/customerio/homework/stream"
	"github.com/customerio/homework/utils"
	"github.com/hashicorp/go-memdb"
	"github.com/labstack/gommon/log"
)

const (
	customerTableName = "customer"
)

// Datastore - in memory concurrent map based data store
type Datastore struct {
	customers *memdb.MemDB
}

// CreateDatastore - creates data store by summarizedEvents and summarizedAttributes
func CreateDatastore(summarizedAttributes map[string]stream.Record, summarizedEvents map[string]map[string]int) (Datastore, error) {

	// Create the DB schema
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			customerTableName: &memdb.TableSchema{
				Name: customerTableName,
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.IntFieldIndex{Field: "ID"},
					},
				},
			},
		},
	}

	// Create a new database
	csm, err := memdb.NewMemDB(schema)
	if err != nil {
		return Datastore{}, err
	}

	txn := csm.Txn(true)
	for k, rec := range summarizedAttributes {
		var customerId int
		var err error

		if customerId, err = strconv.Atoi(k); err != nil {
			return Datastore{}, err
		}

		events := summarizedEvents[k]
		if events == nil {
			// expected by the verify-script
			events = make(map[string]int)
		}

		cs := &serve.Customer{
			ID:          customerId,
			Attributes:  rec.Data,
			Events:      events,
			LastUpdated: int(rec.Timestamp),
		}

		if err := txn.Insert(customerTableName, cs); err != nil {
			log.Error(err)
			return Datastore{}, err
		}
	}
	// commit all writes
	txn.Commit()

	return Datastore{
		customers: csm,
	}, nil
}

func (d Datastore) Get(id int) (*serve.Customer, error) {

	txn := d.customers.Txn(false)
	defer txn.Abort()

	if raw, err := txn.First(customerTableName, "id", id); err == nil && raw != nil {
		customer := raw.(*serve.Customer)
		return customer, nil
	}

	return nil, serve.ErrNotFound
}

func (d Datastore) List(page, count int) ([]*serve.Customer, error) {
	// since the data is stored as a concurrent map, skipping optimizing on time complexity for List
	var start = (page-1)*count + 1
	var end = start + count - 1
	var iter = 1

	cs := make([]*serve.Customer, 0, count)

	txn := d.customers.Txn(false)
	defer txn.Abort()

	// List all the customers
	it, err := txn.Get(customerTableName, "id")
	if err != nil {
		return nil, err
	}

	for obj := it.Next(); obj != nil; obj = it.Next() {
		if iter > end {
			break
		}

		if iter >= start && iter <= end {
			cs = append(cs, obj.(*serve.Customer))
		}
		iter++
	}

	return cs, nil
}

func (d Datastore) Create(id int, attributes map[string]string) (*serve.Customer, error) {

	customer := &serve.Customer{
		ID:          id,
		Attributes:  attributes,
		Events:      nil,
		LastUpdated: int(time.Now().Unix()),
	}

	txn := d.customers.Txn(true)
	if err := txn.Insert(customerTableName, customer); err != nil {
		return nil, err
	}

	txn.Commit()
	return customer, nil
}

func (d Datastore) Update(id int, attributes map[string]string) (*serve.Customer, error) {

	var customer *serve.Customer
	var err error

	if customer, err = d.Get(id); err != nil {
		return nil, err
	}

	customer.Attributes = utils.MergeMaps(attributes, customer.Attributes, true)
	customer.LastUpdated = int(time.Now().Unix())
	txn := d.customers.Txn(true)
	if err := txn.Insert(customerTableName, customer); err != nil {
		return nil, err
	}
	txn.Commit()
	return customer, nil
}

func (d Datastore) Delete(id int) error {

	var customer *serve.Customer
	var err error

	if customer, err = d.Get(id); err != nil {
		return err
	}

	txn := d.customers.Txn(true)
	if err := txn.Delete(customerTableName, customer); err != nil {
		return err
	}
	txn.Commit()
	return nil
}

// TotalCustomers - it iterates over all entries and returns a total count
// Since we're using an in memory db, this is slow; could be optimized when using mysql
func (d Datastore) TotalCustomers() (int, error) {
	var count = 0

	txn := d.customers.Txn(false)
	defer txn.Abort()

	// Get all iterator
	it, err := txn.Get(customerTableName, "id")
	if err != nil {
		return count, err
	}

	for obj := it.Next(); obj != nil; obj = it.Next() {
		count++
	}
	return count, nil
}

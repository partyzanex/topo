package topo

import (
	"fmt"
	"github.com/pkg/errors"
)

// interface for sortable entities
type SortableEntity interface {
	// returns entity id
	Self() interface{}

	// returns id of parent entity
	Parent() interface{}
}

// map for entities
type vertexStore map[interface{}]interface{}

// service for sorting and searching entities
type TopologicalSorter struct {
	// private storage
	storage vertexStore
}

// constructor
func New() *TopologicalSorter {
	return &TopologicalSorter{
		storage: make(vertexStore),
	}
}

// checking and extracting self and parent values from entity
func (TopologicalSorter) extract(entity SortableEntity) (parent, self interface{}, ok bool) {
	parent = entity.Parent()
	self = entity.Self()

	return parent, self, parent != self
}

// check for exists by parent and self values
func (ts *TopologicalSorter) Exists(parent, self int) bool {
	if storage, ok := ts.storage[parent]; ok {
		if _, ok := storage.(vertexStore)[self]; ok {
			return true
		}
	}

	return false
}

// adding entities for sorting
func (ts *TopologicalSorter) Push(entity SortableEntity) error {
	parent, self, ok := ts.extract(entity)
	if !ok {
		return errors.New("entity is not valid")
	}

	storage, ok := ts.storage[parent].(vertexStore)
	if ok {
		if _, ok = storage[self]; ok {
			return errors.New("vertex is defined")
		}
	} else {
		ts.storage[parent] = make(vertexStore)
		storage = make(vertexStore)
	}

	storage[self] = entity
	ts.storage[parent] = storage

	return nil
}

// adding from entities slice for sorting
func (ts *TopologicalSorter) PushAll(entities ...SortableEntity) (err error) {
	for _, entity := range entities {
		if err = ts.Push(entity); err != nil {
			return err
		}
	}

	return nil
}

// getting children entities by parent id
func (ts TopologicalSorter) Child(parent interface{}) ([]interface{}, error) {
	if storage, ok := ts.storage[parent].(vertexStore); ok {
		if n := len(storage); n > 0 {
			results := make([]interface{}, n)

			i := 0
			for _, item := range storage {
				results[i] = item
				i++
			}

			return results, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("key %d is not found", parent))
}

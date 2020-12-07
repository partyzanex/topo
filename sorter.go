package topo

import (
	"fmt"
	"github.com/pkg/errors"
)

var ErrVertexDefined = errors.New("vertex is defined")

// SortableEntity interface for sortable entities
type SortableEntity interface {
	// returns entity id
	Self() interface{}
	// returns id of parent entity
	Parent() interface{}
	// set children elements
	SetChildren(interface{})
}

// vertexStore is a map for entities
type vertexStore map[interface{}]interface{}

// TopologicalSorter represents service for sorting and searching entities
type TopologicalSorter struct {
	// private storage
	storage vertexStore
}

// New is a constructor for TopologicalSorter
func New() *TopologicalSorter {
	return &TopologicalSorter{
		storage: make(vertexStore),
	}
}

// extract a checking and extracting self and parent values from entity
func extract(entity SortableEntity) (parent, self interface{}, ok bool) {
	parent = entity.Parent()
	self = entity.Self()

	return parent, self, parent != self
}

// Exists checks for exists by parent and self values
func (ts *TopologicalSorter) Exists(parent, self int) bool {
	if storage, ok := ts.storage[parent]; ok {
		if _, ok := storage.(vertexStore)[self]; ok {
			return true
		}
	}

	return false
}

// Push adds entities for sorting
func (ts *TopologicalSorter) Push(entity SortableEntity) error {
	parent, self, ok := extract(entity)
	if !ok {
		return errors.New("entity is not valid")
	}

	storage, ok := ts.storage[parent].(vertexStore)
	if ok {
		if _, ok = storage[self]; ok {
			return ErrVertexDefined
		}
	} else {
		ts.storage[parent] = make(vertexStore)
		storage = make(vertexStore)
	}

	storage[self] = entity
	ts.storage[parent] = storage

	return nil
}

// PushAll adds from entities slice for sorting
func (ts *TopologicalSorter) PushAll(entities ...SortableEntity) (err error) {
	for i := range entities {
		if err = ts.Push(entities[i]); err != nil {
			return err
		}
	}

	return nil
}

// Child returns children entities by parent id
func (ts TopologicalSorter) Child(parent interface{}) ([]interface{}, error) {
	if storage, ok := ts.storage[parent].(vertexStore); ok {
		if n := len(storage); n > 0 {
			results := make([]interface{}, n)

			i := 0
			for _, item := range storage {
				result := item.(SortableEntity)
				children, _ := ts.Child(result.Self())
				result.SetChildren(children)
				results[i] = result
				i++
			}

			return results, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("key %v is not found", parent))
}

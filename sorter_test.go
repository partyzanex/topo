package topo

import (
	"fmt"
	"github.com/pkg/errors"
	"testing"
)

type category struct {
	ID, ParentID int
	Name         string
}

func (c category) Self() interface{} {
	return c.ID
}

func (c category) Parent() interface{} {
	return c.ParentID
}

var cats = []*category{
	{
		ID:   1,
		Name: "Category 1",
	},
	{
		ID:       2,
		ParentID: 1,
		Name:     "Category 2",
	},
	{
		ID:   3,
		Name: "Category 3",
	},
	{
		ID:       5,
		ParentID: 3,
		Name:     "Category 5",
	},
	{
		ID:       4,
		ParentID: 1,
		Name:     "Category 4",
	},
	{
		ID:       7,
		ParentID: 5,
		Name:     "Category 7",
	},
	{
		ID:       6,
		ParentID: 2,
		Name:     "Category 6",
	},
	{
		ID:       8,
		ParentID: 2,
		Name:     "Category 8",
	},
	{
		ID:       9,
		ParentID: 2,
		Name:     "Category 9",
	},
	{
		ID:       10,
		ParentID: 8,
		Name:     "Category 10",
	},
	{
		ID:       11,
		ParentID: 8,
		Name:     "Category 11",
	},
	{
		ID:       12,
		ParentID: 2,
		Name:     "Category 12",
	},
}

func TestTopologicalSorter_Push(t *testing.T) {
	sorter := New()

	for _, cat := range cats {
		if err := sorter.Push(cat); err != nil {
			t.Fatal(err)
		}
	}
}

func TestTopologicalSorter_Exists(t *testing.T) {
	sorter := New()

	for _, cat := range cats {
		if err := sorter.Push(cat); err != nil {
			t.Fatal(err)
		}

		if !sorter.Exists(cat.ParentID, cat.ID) {
			t.Fatal(fmt.Sprintf("exists error on %d:%d", cat.ParentID, cat.ID))
		}
	}

	if sorter.Exists(3, 0) || sorter.Exists(2, 4) {
		t.Fatal("not exists error")
	}
}

func TestTopologicalSorter_Child(t *testing.T) {
	sorter := New()

	for _, cat := range cats {
		if err := sorter.Push(cat); err != nil {
			t.Fatal(err)
		}
	}

	child, err := sorter.Child(0)
	if err != nil {
		t.Fatal(errors.Wrap(err, "0"))
	}

	if len(child) != 2 {
		t.Fatal("error 2")
	}

	cat, ok := child[0].(*category)
	if !ok {
		t.Fatal("error category pointer")
	}

	if cat.ID != 1 && cat.ID != 3 {
		t.Fatal("error 3")
	}

	child, err = sorter.Child(2)
	if err != nil {
		t.Fatal(errors.Wrap(err, "0"))
	}

	if len(child) != 4 {
		t.Fatal("error 4")
	}

	cat, ok = child[0].(*category)
	if !ok {
		t.Fatal("error category pointer 2")
	}
}

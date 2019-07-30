package topo_test

import (
	"github.com/partyzanex/topo"
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
	{
		ID:   21,
		Name: "Category 1",
	},
	{
		ID:       22,
		ParentID: 1,
		Name:     "Category 2",
	},
	{
		ID:   23,
		Name: "Category 3",
	},
	{
		ID:       25,
		ParentID: 3,
		Name:     "Category 5",
	},
	{
		ID:       24,
		ParentID: 1,
		Name:     "Category 4",
	},
	{
		ID:       27,
		ParentID: 25,
		Name:     "Category 7",
	},
	{
		ID:       26,
		ParentID: 22,
		Name:     "Category 6",
	},
	{
		ID:       28,
		ParentID: 22,
		Name:     "Category 8",
	},
	{
		ID:       29,
		ParentID: 2,
		Name:     "Category 9",
	},
	{
		ID:       210,
		ParentID: 8,
		Name:     "Category 10",
	},
	{
		ID:       211,
		ParentID: 8,
		Name:     "Category 11",
	},
	{
		ID:       212,
		ParentID: 22,
		Name:     "Category 12",
	},
}

func TestTopologicalSorter_Push(t *testing.T) {
	sorter := topo.New()

	for _, cat := range cats {
		if err := sorter.Push(cat); err != nil {
			t.Fatalf("pushing failed: %s", err)
		}
	}
}

func TestTopologicalSorter_Exists(t *testing.T) {
	sorter := topo.New()

	for _, cat := range cats {
		if err := sorter.Push(cat); err != nil {
			t.Fatalf("pushing failed: %s", err)
		}

		if !sorter.Exists(cat.ParentID, cat.ID) {
			t.Fatalf("exists error on %d:%d", cat.ParentID, cat.ID)
		}
	}

	if sorter.Exists(3, 0) || sorter.Exists(2, 4) {
		t.Fatal("not exists error")
	}
}

func TestTopologicalSorter_Child(t *testing.T) {
	sorter := topo.New()

	for _, cat := range cats {
		if err := sorter.Push(cat); err != nil {
			t.Fatalf("pushing failed: %s", err)
		}
	}

	child, err := sorter.Child(nil)
	if err != nil {
		t.Fatalf("getting childs failed: %s", err)
	}

	if len(child) != 4 {
		t.Fatalf("wrong count of childs: except %d, got %d", 4, len(child))
	}

	cat, ok := child[0].(*category)
	if !ok {
		t.Fatal("error category pointer")
	}

	if cat.ID != 1 && cat.ID != 3 && cat.ID != 21 && cat.ID != 23 {
		t.Fatal("invalid elements")
	}

	child, err = sorter.Child(2)
	if err != nil {
		t.Fatalf("getting childs failed: %s", err)
	}

	if len(child) != 5 {
		t.Fatalf("wrong count of childs: except %d, got %d", 5, len(child))
	}
}

func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		topo.New()
	}
}

func BenchmarkTopologicalSorter_Push(b *testing.B) {
	id := 1
	var lastID int

	sorter := topo.New()

	for i := 0; i < b.N; i++ {
		sorter.Push(&category{
			ID:       id,
			ParentID: lastID,
			Name:     "Test",
		})

		if id%5 == 0 {
			lastID = id
		}
		id++
	}
}

func BenchmarkTopologicalSorter_PushAll(b *testing.B) {
	//sorter := New()
	id := 1
	var lastID int

	var items []*category
	for i := 0; i < 24; i++ {
		items = append(items, &category{
			ID:       id,
			ParentID: lastID,
			Name:     "Test",
		})

		if id%5 == 0 {
			lastID = id
		}
		id++
	}

	for i := 0; i < b.N; i++ {
		sorter := topo.New()
		err := sorter.PushAll(
			items[0], items[1], items[2], items[3], items[4], items[5],
			items[6], items[7], items[8], items[9], items[10], items[11],
			items[12], items[13], items[14], items[15], items[16], items[17],
			items[18], items[19], items[20], items[21], items[22], items[23],
		)
		if err != nil && err != topo.ErrVertexDefined {
			b.Fatal(err)
		}
	}
}

package scorecard

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReducerSetNodes(t *testing.T) {
	node1 := Node{
		ID:       1,
		ParentID: nil,
		Score:    0,
	}
	node2 := Node{
		ID:       2,
		ParentID: &node1.ID,
		Score:    100,
	}

	r := NewReducer()
	r.SetNodes([]*Node{&node2, &node1})

	m := map[int]*Node{}

	node1.children = []*Node{&node2}
	m[node1.ID] = &node1

	m[node2.ID] = &node2

	assert.Equal(t, m, r.m)
}

func TestReducerReduce(t *testing.T) {
	node1 := Node{
		ID:       1,
		ParentID: nil,
		Score:    0,
	}
	node2 := Node{
		ID:       2,
		ParentID: &node1.ID,
		Score:    100,
	}
	node3 := Node{
		ID:       3,
		ParentID: &node1.ID,
		Score:    50,
	}

	r := NewReducer()
	r.SetNodes([]*Node{&node3, &node1, &node2})
	r.Reduce()

	m := map[int]*Node{}

	node1.Score = (node2.Score + node3.Score) / 2
	node1.filled = true
	node1.children = []*Node{&node2}
	m[node1.ID] = &node1

	node2.filled = true
	node1.children = []*Node{&node3}
	m[node2.ID] = &node2

	node3.filled = true
	m[node3.ID] = &node3

	assert.Equal(t, m, r.m)
}

func TestReducerGetRoots(t *testing.T) {
	node1 := Node{
		ID:       1,
		ParentID: nil,
		Score:    0,
	}
	node2 := Node{
		ID:       2,
		ParentID: &node1.ID,
		Score:    100,
	}

	r := NewReducer()
	r.SetNodes([]*Node{&node2, &node1})

	node1.children = []*Node{&node2}

	assert.Equal(t, []*Node{&node1}, r.GetRoots())
}

func TestReducerGetByParentID(t *testing.T) {
	node1 := Node{
		ID:       1,
		ParentID: nil,
		Score:    0,
	}
	node2 := Node{
		ID:       2,
		ParentID: &node1.ID,
		Score:    100,
	}

	r := NewReducer()
	r.SetNodes([]*Node{&node2, &node1})

	assert.Equal(t, []*Node{&node2}, r.GetByParentID(node1.ID))
}

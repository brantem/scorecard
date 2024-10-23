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

	assert.Len(t, r.m, 2)
	assert.Equal(t, r.m[node1.ID].ID, node1.ID)
	assert.Len(t, r.m[node1.ID].children, 1)
	assert.Equal(t, r.m[node1.ID].children[0].ID, node2.ID)
	assert.Equal(t, r.m[node2.ID].ID, node2.ID)
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

	assert.Equal(t, (node2.Score+node3.Score)/2, node1.Score)

	// The following scenario should never occur
	originalScore := node1.Score
	node3.Score = 100

	r.Reduce() // Shouldn't change anything
	assert.Equal(t, originalScore, node1.Score)
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

	roots := r.GetRoots()
	assert.Len(t, roots, 1)
	assert.Equal(t, roots[0].ID, node1.ID)
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

	nodes := r.GetByParentID(node1.ID)
	assert.Len(t, nodes, 1)
	assert.Equal(t, nodes[0].ID, node2.ID)
}

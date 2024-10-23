package scorecard

type Node struct {
	ID       int
	ParentID *int `db:"parent_id"`
	Score    float64
	filled   bool
	children []*Node
}

type Reducer struct {
	m map[int]*Node
}

func NewReducer() *Reducer {
	return &Reducer{
		m: make(map[int]*Node),
	}
}

func (r *Reducer) SetNodes(nodes []*Node) {
	r.m = make(map[int]*Node, len(nodes))

	for _, node := range nodes {
		r.m[node.ID] = node
	}

	for _, node := range nodes {
		if node.ParentID != nil {
			if parent, ok := r.m[*node.ParentID]; ok {
				parent.children = append(parent.children, node)
			}
		}
	}
}

func (r *Reducer) Reduce() {
	for _, node := range r.m {
		if node.ParentID == nil {
			r.fillScore(node)
		}
	}
}

func (r *Reducer) GetRoots() []*Node {
	var nodes []*Node
	for _, node := range r.m {
		if node.ParentID == nil {
			nodes = append(nodes, node)
		}
	}
	return nodes
}

func (r *Reducer) GetByParentID(parentID int) []*Node {
	var nodes []*Node
	for _, node := range r.m {
		if node.ParentID != nil && *node.ParentID == parentID {
			nodes = append(nodes, node)
		}
	}
	return nodes
}

func (r *Reducer) fillScore(parent *Node) {
	if parent.filled {
		return
	}

	if len(parent.children) == 0 {
		parent.filled = true
		return
	}

	weight := 100 / float64(len(parent.children))
	for _, child := range parent.children {
		r.fillScore(child)
		parent.Score += (child.Score * weight) / 100
	}

	parent.filled = true
}

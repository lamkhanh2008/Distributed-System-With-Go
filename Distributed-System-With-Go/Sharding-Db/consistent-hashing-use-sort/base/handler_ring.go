package consistent

import (
	"errors"
	"sort"
)

func (r *Ring) AddNode(id string) {
	r.Lock()
	defer r.Unlock()
	node := NewNode(id)
	r.Nodes = append(r.Nodes, node)
	sort.Sort(r.Nodes)
}

func (r *Ring) RemoveNode(id string) error {
	r.Lock()
	defer r.Unlock()
	i := r.search(id)
	if i >= r.Nodes.Len() || r.Nodes[i].Id != id {
		return errors.New("node not found")
	}

	r.Nodes = append(r.Nodes[:i], r.Nodes[i+1:]...)
	return nil
}

func (r *Ring) Get(id string) string {
	i := r.search(id)
	if i >= r.Nodes.Len() {
		i = 0
	}
	return r.Nodes[i].Id
}

func (r *Ring) search(id string) int {
	searchfn := func(i int) bool {
		return r.Nodes[i].HashId >= hashId(id)
	}

	return sort.Search(r.Nodes.Len(), searchfn)
}

package consistent

func (n Nodes) Len() int {
	return len(n)
}

func (n Nodes) Less(i, j int) bool {
	return n[i].HashId < n[j].HashId
}
func (n Nodes) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

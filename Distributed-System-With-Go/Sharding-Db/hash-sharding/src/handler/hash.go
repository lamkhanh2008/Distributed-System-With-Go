package handler

func FnvHash(data string) uint32 {
	const prime = 16777619
	hash := uint32(2166136261)
	for i := 0; i < len(data); i++ {
		hash ^= uint32(data[i])
		hash *= prime
	}
	return hash
}

package cache

type ByteView struct {
	b []byte
}

func (b ByteView) Len() int {
	return len(b.b)
}

func (b ByteView) ByteSlice() []byte {
	return cloneBytes(b.b)
}

func (b ByteView) String() string {
	return string(b.b)
}

func cloneBytes(bytes []byte) []byte {
	var c = make([]byte, len(bytes))

	copy(c, bytes)
	return c
}

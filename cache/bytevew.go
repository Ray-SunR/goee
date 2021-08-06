package cache

type ByteView struct {
	bytes []byte
}

func (byteView ByteView) SizeInBytes() int64 {
	return (int64) (len(byteView.bytes))
}

func cloneBytes(bytes []byte) []byte {
	clonedBytes := make([]byte, len(bytes))
	copy(clonedBytes, bytes)

	return clonedBytes
}

func (byteView ByteView) ByteSlice() []byte {
	return cloneBytes(byteView.bytes)
}

func (byteView ByteView) String() string {
	return string(byteView.bytes)
}

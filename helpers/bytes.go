package helpers


// copy src to new []byte and return it
func Copy(src []byte) []byte {
	dst := make([]byte, len(src))
	copy(dst,src)
	return dst
}




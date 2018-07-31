package tools

//LE2Int unmarshals the given []byte which is in little endian to an int
func LE2Int(bts []byte) uint {
	var res uint
	for k := range bts {
		res |= uint(bts[k]) << uint(8*k)
	}
	return res
}

//Int2LE marshals the given int to a []byte in little endian order
func Int2LE(i uint) [4]byte {
	var bts = [4]byte{byte(i & 0xff)}

	for ind := 1; ; ind++ {
		i >>= 8
		if 0 == i {
			break
		}
		bts[ind] = byte(i & 0xff)
	}
	return bts
}

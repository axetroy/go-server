package smallendian

import (
	"encoding/binary"
)

// helper function which converts a uint64 to a []byte in Small Endian
func ToUint64(n uint64) [8]byte {
	s := make([]byte, 8)
	binary.LittleEndian.PutUint64(s, n)
	a := [8]byte{}
	copy(a[:], s[:8])
	return a
}

// helper function which converts a small endian []byte to a uint64
func FromUint64(data [8]byte) uint64 {
	ui64 := binary.LittleEndian.Uint64(data[:])
	return ui64
}

// helper function which converts a int to a []byte in Small Endian
func ToInt(n int) [4]byte {
	s := make([]byte, 4)
	binary.LittleEndian.PutUint32(s, uint32(n))
	a := [4]byte{}
	copy(a[:], s[:4])
	return a
}

// helper function which converts a small endian []byte to a int
func FromInt(data [4]byte) int {
	ui32 := binary.LittleEndian.Uint32(data[:])
	return int(ui32)
}

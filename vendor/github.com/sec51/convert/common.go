package convert

import (
	"math"
)

// Helper function which rounds the float to the nearest integet
func Round(n float64) uint64 {
	if n < 0 {
		return uint64(math.Ceil(n - 0.5))
	}
	return uint64(math.Floor(n + 0.5))
}

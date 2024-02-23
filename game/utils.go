package game

import (
	"math"
	"sort"

	"github.com/google/uuid"
)

func GetOrderedIds[T any](m map[uuid.UUID]T) []uuid.UUID {

	keys := make([]uuid.UUID, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i].String() < keys[j].String()
	})
	return keys
}

func EuclidianDistance(a, b Vector) float64 {
	return math.Sqrt((b.X-a.X)*(b.X-a.X) + (b.Y-a.Y)*(b.Y-a.Y))
}
func QuickDistance(a, b Vector) float64 {
	return math.Abs(b.X-a.X) + math.Abs(b.Y-a.Y)
}

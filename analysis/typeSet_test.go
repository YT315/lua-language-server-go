package analysis

import (
	"math/rand"
	"testing"
)

func TestSet(t *testing.T) {
	set := NewTypeSet()
	a := []TypeInfo{}
	for i := 0; i < 100; i++ {
		ty := &TypeNumber{Value: float64(i)}
		set.Add(ty)
		a = append(a, ty)
	}
	for index, val := range a {
		if !set.Contain(val) {
			t.Error("not contain")
		}
		if set.IndexOf(val) != index {
			t.Error("indexerr")
		}
	}
	index := rand.Intn(99)
	val := a[index]
	set.Del(val)
	for i := 0; i < set.Count; i++ {
		if set.Types[i] == val {
			t.Error("del err")
		}
	}
}

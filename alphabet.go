package shortuuid

import (
	"fmt"
	"math"
	"slices"
	"unicode/utf8"
)

// DefaultAlphabet is the default alphabet used.
const (
	DefaultAlphabet = "23456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
)

type alphabet struct {
	chars    []rune
	len      int64
	encLen   int64
	maxBytes int64
	indexMap map[rune]int64
}

// Remove duplicates and sort it to ensure reproducibility.
func newAlphabet(s string) alphabet {
	a := newEncodeAlphabet(s)
	a.indexMap = make(map[rune]int64, len(a.chars))
	for index, c := range a.chars {
		a.indexMap[c] = int64(index)
	}
	return a
}

func newEncodeAlphabet(s string) alphabet {
	abc := []rune(s)
	slices.Sort(abc)
	abc = slices.Compact(abc)

	if len(abc) < 2 {
		panic("encoding alphabet must be at least two characters")
	}

	a := alphabet{
		chars:    abc,
		len:      int64(len(abc)),
		encLen:   int64(math.Ceil(128 / math.Log2(float64(len(abc))))),
		maxBytes: 1,
	}
	for _, c := range a.chars {
		var b int64
		switch i := uint32(c); {
		case i <= rune1Max:
			b = 1
		case i <= rune2Max:
			b = 2
		case i < surrogateMin, surrogateMax < i && i <= rune3Max:
			b = 3
		case i > rune3Max && i <= utf8.MaxRune:
			b = 4
		default:
			b = 3
		}
		if b > a.maxBytes {
			a.maxBytes = b
		}
	}

	return a
}

func (a *alphabet) Length() int64 {
	return a.len
}

// Index returns the index of the first instance of t in the alphabet, or an
// error if t is not present.
func (a *alphabet) Index(t rune) (int64, error) {
	if a.indexMap == nil {
		return 0, fmt.Errorf("index map is not initialized")
	}
	i, ok := a.indexMap[t]
	if !ok {
		return 0, fmt.Errorf("element '%v' is not part of the alphabet", t)
	}
	return i, nil
}

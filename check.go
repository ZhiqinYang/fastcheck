package fastcheck

import (
	"math"
	"strings"
	"sync"
)

const (
	length = 1 << 18
)

var empty = struct{}{}

type FastCheck struct {
	hashSet    map[string]struct{}
	fastCheck  [length]byte
	fastLength [length]byte
	charCheck  [length]bool
	endCheck   [length]bool
	maxWordLen float64
	minWordLen float64
	ignoreCase bool
	sync.RWMutex
}

func NewFastCheck(ignoreCase bool) *FastCheck {
	return &FastCheck{hashSet: make(map[string]struct{}), minWordLen: math.MaxFloat64, ignoreCase: ignoreCase}
}

func (fc *FastCheck) AddWord(text string) {
	var runes []rune
	if fc.ignoreCase {
		text = strings.ToUpper(text)
	}

	fc.Lock()
	defer fc.Unlock()
	if _, ok := fc.hashSet[text]; ok {
		return
	}

	runes = []rune(text)
	sz := float64(len(runes))
	fc.maxWordLen = math.Max(fc.maxWordLen, sz)
	fc.minWordLen = math.Min(fc.minWordLen, sz)
	for i := uint(0); i < 7 && i < uint(sz); i++ {
		fc.fastCheck[runes[i]] |= byte(1 << i)
	}

	for i := uint(7); i < uint(sz); i++ {
		fc.fastCheck[runes[i]] |= 0x80
	}

	if sz == 1 {
		fc.charCheck[runes[0]] = true
	} else {
		fc.fastLength[runes[0]] |= byte(1 << uint(math.Min(7, sz-2)))
		fc.endCheck[runes[int(sz)-1]] = true
		fc.hashSet[text] = empty
	}
}

func (fc *FastCheck) HasWord(text string) bool {
	var idx int
	var runes []rune
	if fc.ignoreCase {
		runes = []rune(strings.ToUpper(text))
	} else {
		runes = []rune(text)
	}

	fc.RLock()
	defer fc.RUnlock()
	sz := len(runes)
	for idx < sz {
		var count = 1
		if idx > 0 || fc.fastCheck[runes[idx]]&1 == 0 {
			for ; idx < (sz-1) && (fc.fastCheck[runes[idx+1]]&1) == 0; idx++ {
			}

			if idx < sz-1 {
				idx++
			}

		}

		var begin = runes[idx]
		if fc.minWordLen == 1 && fc.charCheck[begin] {
			return true
		}

		for j := 1; j <= int(math.Min(fc.maxWordLen, float64(sz-idx-1))); j += 1 {
			var current = runes[idx+j]
			if fc.fastCheck[current]&1 == 0 {
				count++
			}

			if fc.fastCheck[current]&(1<<uint(math.Min(float64(j), 7))) == 0 {
				break
			}

			if float64(j+1) >= fc.minWordLen {
				if (fc.fastLength[begin]&(1<<uint(math.Min(float64(j-1), 7))) > 0) && fc.endCheck[current] {
					if _, ok := fc.hashSet[string(runes[idx:idx+j+1])]; ok {
						return true
					}
				}
			}
		}
		idx += count
	}
	return false
}

func (fc *FastCheck) ReplaceWith(text string, char rune) string {
	var idx int
	runes := []rune(strings.ToUpper(text))
	original := []rune(text)
	sz := len(runes)

	fc.RLock()
	defer fc.RUnlock()

	for idx < sz {
		var count = 1
		if idx > 0 || fc.fastCheck[runes[idx]]&1 == 0 {
			for ; idx < (sz-1) && (fc.fastCheck[runes[idx+1]]&1) == 0; idx++ {
			}

			if idx < sz-1 {
				idx++
			}

		}

		var begin = runes[idx]
		if fc.minWordLen == 1 && fc.charCheck[begin] {
			original[idx] = char
			idx++
			continue
		}

		for j := 1; j <= int(math.Min(fc.maxWordLen, float64(sz-idx-1))); j += 1 {
			var current = runes[idx+j]
			if fc.fastCheck[current]&1 == 0 {
				count++
			}

			if fc.fastCheck[current]&(1<<uint(math.Min(float64(j), 7))) == 0 {
				break
			}

			if float64(j+1) >= fc.minWordLen {
				if (fc.fastLength[begin]&(1<<uint(math.Min(float64(j-1), 7))) > 0) && fc.endCheck[current] {
					if _, ok := fc.hashSet[string(runes[idx:idx+j+1])]; ok {
						for x := idx; x < idx+j+1; x++ {
							original[x] = char
						}
						break
					}
				}
			}
		}
		idx += count
	}
	return string(original)
}

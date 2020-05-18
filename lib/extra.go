package vegeta

import (
	"bufio"
	"bytes"
	"errors"
	"sort"
	"strings"
)

// Extra represents extra properties that should be passed from targets to results.
type Extra map[string]string

// Equal returns true if the given Extra is equal to the other.
func (e Extra) Equal(other Extra) bool {
	if len(e) != len(other) {
		return false
	}
	if e == nil || other == nil {
		return e == nil && other == nil
	}
	for key, value1 := range e {
		if value2, ok := other[key]; !ok || value1 != value2 {
			return false
		}
	}

	return true
}

// Clone returns a copy of the extra properties.
func (e Extra) Clone() Extra {
	if e == nil {
		return nil
	}

	clone := make(Extra, len(e))
	for k, v := range e {
		clone[k] = v
	}
	return clone
}

// Serialize the extra properties to a byte slice.
func (e Extra) Serialize() []byte {
	if e == nil {
		return nil
	}

	keys := make([]string, 0, len(e))
	for k := range e {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var b bytes.Buffer
	for i, k := range keys {
		b.WriteString(k)
		b.WriteString("=")
		b.WriteString(e[k])
		if i < len(keys)-1 {
			b.WriteByte('\n')
		}
	}
	return b.Bytes()
}

func deserializeExtra(serializedState []byte) (Extra, error) {
	if serializedState == nil {
		return nil, nil
	}

	extra := Extra{}
	s := bufio.NewScanner(bytes.NewReader(serializedState))
	for s.Scan() {
		line := s.Text()
		tokens := strings.SplitN(line, "=", 2)
		if len(tokens) != 2 {
			return nil, errors.New("bad extra bytes")
		}
		extra[tokens[0]] = tokens[1]
	}
	return extra, nil
}

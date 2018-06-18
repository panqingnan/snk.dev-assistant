package maputil

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type fnCallSignature struct {
	Value  interface{}
	Path   []string
	IsLeaf bool
}

type walkTestStruct struct {
	Name   string
	Value  int64
	Flags  []bool
	Submap map[string]string
}

func TestMapWalkScalar(t *testing.T) {
	assert := require.New(t)

	assert.Nil(Walk(nil, nil))

	i := 0
	assert.Nil(Walk(42, func(value interface{}, path []string, isLeaf bool) error {
		i += 1

		assert.Equal(42, value)
		assert.Nil(path)
		assert.True(isLeaf)

		return nil
	}))

	assert.Equal(1, i)
}

func TestMapWalkFlatMap(t *testing.T) {
	assert := require.New(t)

	input := map[string]interface{}{
		`a`: 1,
		`b`: true,
		`c`: `Three`,
	}

	checkAnswers := func(callSignatures map[string]fnCallSignature) {
		v, ok := callSignatures[``]
		assert.True(ok)
		assert.Equal(fnCallSignature{input, nil, false}, v)

		v, ok = callSignatures[`a`]
		assert.True(ok)
		assert.Equal(fnCallSignature{1, []string{`a`}, true}, v)

		v, ok = callSignatures[`b`]
		assert.True(ok)
		assert.Equal(fnCallSignature{true, []string{`b`}, true}, v)

		v, ok = callSignatures[`c`]
		assert.True(ok)
		assert.Equal(fnCallSignature{`Three`, []string{`c`}, true}, v)
	}

	callSignatures := make(map[string]fnCallSignature)
	assert.Nil(Walk(input, func(value interface{}, path []string, isLeaf bool) error {
		callSignatures[strings.Join(path, `.`)] = fnCallSignature{
			Value:  value,
			Path:   path,
			IsLeaf: isLeaf,
		}

		return nil
	}))

	checkAnswers(callSignatures)

	callSignatures = make(map[string]fnCallSignature)
	assert.Nil(Walk(&input, func(value interface{}, path []string, isLeaf bool) error {
		callSignatures[strings.Join(path, `.`)] = fnCallSignature{
			Value:  value,
			Path:   path,
			IsLeaf: isLeaf,
		}

		return nil
	}))

	checkAnswers(callSignatures)
}

func TestMapWalkNestedMap(t *testing.T) {
	assert := require.New(t)

	callSignatures := make(map[string]fnCallSignature)

	b2a_map := map[string]interface{}{
		`a`: true,
	}

	b2b_map := map[string]interface{}{
		`a`: 42,
	}

	b2_slice := []map[string]interface{}{b2a_map, b2b_map}
	b_map := map[string]interface{}{
		`b1`: 11,
		`b2`: b2_slice,
	}

	input := map[string]interface{}{
		`a`: 1,
		`b`: b_map,
	}

	Walk(input, func(value interface{}, path []string, isLeaf bool) error {
		callSignatures[strings.Join(path, `.`)] = fnCallSignature{
			Value:  value,
			Path:   path,
			IsLeaf: isLeaf,
		}

		return nil
	})

	v, ok := callSignatures[``]
	assert.True(ok)
	assert.Equal(fnCallSignature{input, nil, false}, v)

	v, ok = callSignatures[`a`]
	assert.True(ok)
	assert.Equal(fnCallSignature{1, []string{`a`}, true}, v)

	v, ok = callSignatures[`b`]
	assert.True(ok)
	assert.Equal(fnCallSignature{b_map, []string{`b`}, false}, v)

	v, ok = callSignatures[`b.b1`]
	assert.True(ok)
	assert.Equal(fnCallSignature{11, []string{`b`, `b1`}, true}, v)

	v, ok = callSignatures[`b.b2`]
	assert.True(ok)
	assert.Equal(fnCallSignature{b2_slice, []string{`b`, `b2`}, false}, v)

	v, ok = callSignatures[`b.b2.0`]
	assert.True(ok)
	assert.Equal(fnCallSignature{b2a_map, []string{`b`, `b2`, `0`}, false}, v)

	v, ok = callSignatures[`b.b2.0.a`]
	assert.True(ok)
	assert.Equal(fnCallSignature{true, []string{`b`, `b2`, `0`, `a`}, true}, v)

	v, ok = callSignatures[`b.b2.1.a`]
	assert.True(ok)
	assert.Equal(fnCallSignature{42, []string{`b`, `b2`, `1`, `a`}, true}, v)
}

func TestMapWalkStruct(t *testing.T) {
	assert := require.New(t)

	input := walkTestStruct{
		Name:  `First`,
		Value: 42,
		Flags: []bool{true, true, false, true},
		Submap: map[string]string{
			`a`: `1`,
			`b`: `true`,
			`c`: `Three`,
		},
	}

	checkAnswers := func(callSignatures map[string]fnCallSignature) {
		v, ok := callSignatures[``]
		assert.True(ok)
		assert.Equal(fnCallSignature{input, nil, false}, v)

		v, ok = callSignatures[`Name`]
		assert.True(ok)
		assert.Equal(fnCallSignature{`First`, []string{`Name`}, true}, v)

		v, ok = callSignatures[`Value`]
		assert.True(ok)
		assert.Equal(fnCallSignature{int64(42), []string{`Value`}, true}, v)

		v, ok = callSignatures[`Flags`]
		assert.True(ok)
		assert.Equal(fnCallSignature{input.Flags, []string{`Flags`}, false}, v)

		v, ok = callSignatures[`Flags.0`]
		assert.True(ok)
		assert.Equal(fnCallSignature{true, []string{`Flags`, `0`}, true}, v)

		v, ok = callSignatures[`Flags.1`]
		assert.True(ok)
		assert.Equal(fnCallSignature{true, []string{`Flags`, `1`}, true}, v)

		v, ok = callSignatures[`Flags.2`]
		assert.True(ok)
		assert.Equal(fnCallSignature{false, []string{`Flags`, `2`}, true}, v)

		v, ok = callSignatures[`Flags.3`]
		assert.True(ok)
		assert.Equal(fnCallSignature{true, []string{`Flags`, `3`}, true}, v)

		v, ok = callSignatures[`Submap`]
		assert.True(ok)
		assert.Equal(fnCallSignature{input.Submap, []string{`Submap`}, false}, v)

		v, ok = callSignatures[`Submap.a`]
		assert.True(ok)
		assert.Equal(fnCallSignature{`1`, []string{`Submap`, `a`}, true}, v)

		v, ok = callSignatures[`Submap.b`]
		assert.True(ok)
		assert.Equal(fnCallSignature{`true`, []string{`Submap`, `b`}, true}, v)

		v, ok = callSignatures[`Submap.c`]
		assert.True(ok)
		assert.Equal(fnCallSignature{`Three`, []string{`Submap`, `c`}, true}, v)
	}

	callSignatures := make(map[string]fnCallSignature)
	assert.Nil(Walk(input, func(value interface{}, path []string, isLeaf bool) error {
		callSignatures[strings.Join(path, `.`)] = fnCallSignature{
			Value:  value,
			Path:   path,
			IsLeaf: isLeaf,
		}

		return nil
	}))

	checkAnswers(callSignatures)

	callSignatures = make(map[string]fnCallSignature)
	assert.Nil(Walk(&input, func(value interface{}, path []string, isLeaf bool) error {
		callSignatures[strings.Join(path, `.`)] = fnCallSignature{
			Value:  value,
			Path:   path,
			IsLeaf: isLeaf,
		}

		return nil
	}))

	checkAnswers(callSignatures)
}

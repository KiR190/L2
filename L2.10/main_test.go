package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetKey(t *testing.T) {
	line := "apple\tbanana\tcherry"

	assert.Equal(t, "apple", getKey(line, 1))
	assert.Equal(t, "banana", getKey(line, 2))
	assert.Equal(t, "cherry", getKey(line, 3))
	assert.Equal(t, "", getKey(line, 4))
	assert.Equal(t, line, getKey(line, 0))
}

func TestParseHuman(t *testing.T) {
	assert.Equal(t, 1024.0, parseHuman("1K"))
	assert.Equal(t, 1048576.0, parseHuman("1M"))
	assert.Equal(t, 1073741824.0, parseHuman("1G"))
	assert.Equal(t, 1536.0, parseHuman("1.5K"))
	assert.Equal(t, 0.0, parseHuman(""))
	assert.Equal(t, 42.0, parseHuman("42"))
}

func TestUniq(t *testing.T) {
	lines := []string{"a", "a", "b", "b", "c"}
	expected := []string{"a", "b", "c"}
	assert.Equal(t, expected, uniq(lines))

	assert.Equal(t, []string{}, uniq([]string{}))
	assert.Equal(t, []string{"x"}, uniq([]string{"x"}))
}

func TestCompareNumeric(t *testing.T) {
	a := "10"
	b := "2"

	assert.False(t, compare(a, b, 0, true, false, false, false))
	assert.True(t, compare(a, b, 0, true, false, false, true)) // reverse
}

func TestCompareLexicographic(t *testing.T) {
	a := "apple"
	b := "banana"

	assert.True(t, compare(a, b, 0, false, false, false, false))
	assert.False(t, compare(a, b, 0, false, false, false, true)) // reverse
}

func TestCompareMonth(t *testing.T) {
	a := "Jan"
	b := "Feb"

	assert.True(t, compare(a, b, 0, false, true, false, false))
	assert.False(t, compare(a, b, 0, false, true, false, true))
}

func TestCompareHuman(t *testing.T) {
	a := "1K"
	b := "1M"

	assert.True(t, compare(a, b, 0, false, false, true, false))
	assert.False(t, compare(a, b, 0, false, false, true, true))
}

func TestIsSorted(t *testing.T) {
	linesSorted := []string{"1", "2", "3"}
	linesUnsorted := []string{"3", "1", "2"}

	assert.True(t, isSorted(linesSorted, 0, true, false, false, false))
	assert.False(t, isSorted(linesUnsorted, 0, true, false, false, false))
}

func TestExpandShortFlags(t *testing.T) {
	args := []string{"-nr", "--help", "-k"}
	expected := []string{"-n", "-r", "--help", "-k"}
	assert.Equal(t, expected, expandShortFlags(args))
}

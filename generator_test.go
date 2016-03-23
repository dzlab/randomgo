package random

import (
	"strconv"
	"strings"
	"testing"
	"time"
)

// test for date genetor
func TestDateGenerator(t *testing.T) {
	date1 := "2012-01-01 00:00:00"
	date2 := "2016-02-29 00:00:00"
	format := "2006-01-02 15:04:05" //"yyyy-MM-dd hh:mm:ss"
	generator, err := NewDateGenerator(date1, date2, format)
	CheckTrue(t, err == nil, "Valid inputs should not cause constructor to fail")
	min, err := time.Parse(format, date1)
	CheckTrue(t, err == nil, "Should not fail to parse valid date")
	max, err := time.Parse(format, date2)
	CheckTrue(t, err == nil, "Should not fail to parse valid date")
	// generate some dates data
	c := generator.Generate()
	actual := []string{string(<-c), string(<-c), string(<-c), string(<-c), string(<-c), string(<-c), string(<-c), string(<-c)}
	// check all generated dates are between two input dates
	for _, d := range actual {
		date, err := time.Parse(format, d)
		if err != nil {
			t.Error("Should not fail to parse valid date")
		}
		if date.Before(min) {
			t.Error("Generated date", date, "should not be before", min)
		}
		if date.After(max) {
			t.Error("Generated date", date, "should not be after", max)
		}
	}
}

// test for incremental integer generator
func TestIncrementGenerator(t *testing.T) {
	g, err := NewIncrementGenerator(5)
	CheckTrue(t, err == nil, "Should not fail to convert integer")
	c := g.Generate()
	actual := []string{string(<-c), string(<-c), string(<-c), string(<-c), string(<-c), string(<-c), string(<-c), string(<-c)}
	for idx, s := range actual {
		i, err := strconv.Atoi(s)
		CheckTrue(t, err == nil, "Should generate valid integers")
		CheckTrue(t, 5+idx+1 == i, "Should generate ancremental values")
	}
}

// test for floating point number generator
func TestFloatGenerator(t *testing.T) {
	g, err := NewFloatGenerator(-28.0168595, 52.5388779)
	CheckTrue(t, err == nil, "Should not fail to convert integer")
	c := g.Generate()
	actual := []string{string(<-c), string(<-c), string(<-c), string(<-c), string(<-c), string(<-c), string(<-c), string(<-c)}
	// check all generated floats are between two input min/max
	fmin, _ := strconv.ParseFloat("-28.0168595", 64)
	fmax, _ := strconv.ParseFloat("52.5388779", 64)
	for _, f := range actual {
		float, err := strconv.ParseFloat(string(f), 64)
		CheckTrue(t, err == nil, "Should generate valid floating point numbers")
		CheckTrue(t, float >= fmin, "Should generate a number greater or equals to min")
		CheckTrue(t, float < fmax, "Should generate a number lower than max")
	}
}

/*
 * Test pick values generator
 */
func TestPickGenerator(t *testing.T) {
	// pick from file
	g, err := NewPickFromFileGenerator("/non/existant.file")
	CheckTrue(t, err != nil, "Should return an error when reading from non existant file")
	CheckTrue(t, g == nil, "Should return a nil generator when reading from non existant file")
	// pick from values
	values := []string{"A", "B", "C", "D", "E"}
	g, err = NewPickFromValuesGenerator(values)
	CheckTrue(t, err == nil, "Should not fail to create a generator from valid input")
	c := g.Generate()
	actual := []string{string(<-c), string(<-c), string(<-c), string(<-c), string(<-c), string(<-c), string(<-c), string(<-c)}
	// check all generated values are withing the provided slice
	for _, s := range actual {
		CheckTrue(t, len(string(s)) == 1, "Should pick a valid element")
		CheckTrue(t, strings.Index("ABCDE", string(s)) > -1, "Should generate a string picked from the generated values")
	}
}

/*
 * Check if two string slices are equals otherwise report an error
 */
func Equals(t *testing.T, actual, expected []string) {
	for i := 0; i < len(actual); i++ {
		if actual[i] != expected[i] {
			t.Error("Should be equal", actual[i], expected[i])
		}
	}
}

/*
 * Check if a condition is true otherwise report an error
 */
func CheckTrue(t *testing.T, condition bool, msg string) {
	if !condition {
		t.Error(msg)
	}
}

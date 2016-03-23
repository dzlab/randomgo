package random

import (
	"bytes"
	"io/ioutil"
	"math/rand"
	"strconv"
	"text/template"
	"time"
)

var (
	chars = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@#$^&*(){}][:<>.")
)

type RecordGenerator interface {
	Generate() <-chan []byte
}

/*
 * A generator of boolean data
 */
type BoolGenerator struct {
	random *rand.Rand
}

/*
 * Create a boolean data generator
 */
func NewBoolGenerator() (*BoolGenerator, error) {
	source := rand.NewSource(time.Now().UnixNano())
	return &BoolGenerator{random: rand.New(source)}, nil
}

/*
 * Generates a random boolean
 */
func (this *BoolGenerator) Generate() <-chan []byte {
	channel := make(chan []byte)
	go func() {
		for {
			f := this.random.Float64()
			channel <- []byte(strconv.FormatBool(f <= 0.5))
		}
	}()
	return channel
}

/*
 * A generator of dates data
 */
type DateGenerator struct {
	random   *rand.Rand
	min      time.Time
	duration int64
	format   string
}

/*
 * Create a new date generators
 */
func NewDateGenerator(min, max, format string) (*DateGenerator, error) {
	// parse min date
	pmin, err := time.Parse(format, min)
	if err != nil {
		return nil, err
	}
	// parse max date
	pmax, err := time.Parse(format, max)
	if err != nil {
		return nil, err
	}
	elapsed := pmax.Sub(pmin)
	source := rand.NewSource(time.Now().UnixNano())
	return &DateGenerator{random: rand.New(source), min: pmin, duration: int64(elapsed), format: format}, nil
}

func (this *DateGenerator) Generate() <-chan []byte {
	channel := make(chan []byte)
	go func() {
		for {
			elapsed := time.Duration(this.random.Int63n(this.duration))
			date := this.min.Add(elapsed)
			formatted := date.Format(this.format)
			channel <- []byte(formatted)
		}
	}()
	return channel
}

/*
 * A generator of float data
 */
type FloatGenerator struct {
	random   *rand.Rand
	min      float64
	interval float64
}

/*
 * Create a new float generator
 */
func NewFloatGenerator(min, max float64) (*FloatGenerator, error) {
	source := rand.NewSource(time.Now().UnixNano())
	return &FloatGenerator{random: rand.New(source), min: min, interval: (max - min)}, nil
}

/*
 *
 */
func (this *FloatGenerator) Generate() <-chan []byte {
	channel := make(chan []byte)
	go func() {
		for {
			float := this.random.Float64()*this.interval + this.min
			formatted := strconv.FormatFloat(float, 'f', 7, 64)
			channel <- []byte(formatted)
		}
	}()
	return channel
}

/*
 * A generator of random strings from a File
 */
type PickGenerator struct {
	random     *rand.Rand
	collection [][]byte
}

/*
 * Create a pick generator from lines of a file
 */
func NewPickFromFileGenerator(filename string) (*PickGenerator, error) {
	bb, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	collection := bytes.Split(bb, []byte("\n"))
	source := rand.NewSource(time.Now().UnixNano())
	return &PickGenerator{random: rand.New(source), collection: collection}, nil
}

/*
 * Create a pick generator from a collection of strings
 */
func NewPickFromValuesGenerator(values []string) (*PickGenerator, error) {
	var collection [][]byte
	for _, val := range values {
		collection = append(collection, []byte(val))
	}
	source := rand.NewSource(time.Now().UnixNano())
	return &PickGenerator{random: rand.New(source), collection: collection}, nil
}

/*
 * Generate random bytes by selecting randomly from a defined collection
 */
func (this *PickGenerator) Generate() <-chan []byte {
	channel := make(chan []byte)
	go func() {
		size := len(this.collection)
		for {
			index := this.random.Intn(size)
			channel <- this.collection[index]
		}
	}()
	return channel
}

/*
 * A generator of incremental integers
 */
type IncrementGenerator struct {
	current int
}

/*
 * Create a new incremental values generator
 */
func NewIncrementGenerator(initial int) (*IncrementGenerator, error) {
	return &IncrementGenerator{current: initial}, nil
}

/*
 * Generate random bytes by selecting randomly from a defined collection
 */
func (this *IncrementGenerator) Generate() <-chan []byte {
	channel := make(chan []byte)
	go func() {
		for {
			this.current += 1
			channel <- []byte(strconv.Itoa(this.current))
		}
	}()
	return channel
}

/*
 * A generator of random bytes data
 */
type BytesGenerator struct {
	random   *rand.Rand
	size     int
	min      int
	interval int
}

/*
 * Create a generator for fixed size strings
 */
func NewFixedSizeStringGenerator(size int) (*BytesGenerator, error) {
	source := rand.NewSource(time.Now().UnixNano())
	return &BytesGenerator{random: rand.New(source), size: size, min: -1, interval: -1}, nil
}

/*
 * Create a generator for variable size strings
 */
func NewVariableSizeStringGenerator(min, max int) (*BytesGenerator, error) {
	source := rand.NewSource(time.Now().UnixNano())
	return &BytesGenerator{random: rand.New(source), size: -1, min: min, interval: (max - min)}, nil
}

/*
 * Calcaulate the rigth size to use whether fixed or variable
 */
func (this *BytesGenerator) Size() int {
	if this.size > -1 {
		return this.size
	}
	return this.min + this.random.Intn(this.interval)
}

// Returns a random message generated from the chars byte slice.
// Message length of m bytes as defined by msgSize.
func (this *BytesGenerator) Generate() <-chan []byte {
	channel := make(chan []byte)
	go func() {
		// serve data for ever
		for {
			m := make([]byte, this.Size())
			for i := range m {
				m[i] = chars[this.random.Intn(len(chars))]
			}
			channel <- m
		}
	}()
	return channel
}

func NewJsonGenerator(name, templ string, data interface{}) (*JsonGenerator, error) {
	t := template.New(name)
	t, err := t.Parse(templ)
	if err != nil {
		return nil, err
		//buff := bytes.NewBufferString("")
		//t.Execute(buff, person)
	}
	return &JsonGenerator{templ: t, data: data}, nil
}

type JsonGenerator struct {
	templ *template.Template
	data  interface{}
}

func (this *JsonGenerator) Generate() <-chan []byte {
	channel := make(chan []byte)
	go func() {
		// serve data fro ever
		for {
			buff := bytes.NewBufferString("")
			this.templ.Execute(buff, this.data)
			channel <- buff.Bytes()
		}
	}()
	return channel
}

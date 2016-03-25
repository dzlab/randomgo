package random

import (
	"encoding/json"
	"log"
	"strconv"
)

type Encoder interface {
	Encode(obj Object) string
}

type Decoder interface {
	Decode(str string) *Object
}

type KVEncoder struct {
	sep1 string
	sep2 string
}

func NewKVEncoder(sep1, sep2 string) *KVEncoder {
	return &KVEncoder{sep1: sep1, sep2: sep2}
}

/*
 * Generate a key-value string of this object's attributes using the provided separator
 * sep1 separates the key from its value
 * sep2 separates key-value pairs
 */
func (this *KVEncoder) Encode(obj *Object) string {
	result := ""
	for _, a := range obj.Attributes {
		// check if we have to ignore this attribute or not
		if a.Optional > 0 && obj.random.Float64() > a.Optional {
			continue
		}
		value := string(<-a.Channel)
		result += a.Name + this.sep1 + value + this.sep2
	}
	result = result[:len(result)-len(this.sep2)]
	return result
}

type JsonEncoder struct {
}

func NewJsonEncoder() *JsonEncoder {
	return &JsonEncoder{}
}

/*
 * Generate a Jason string of this object's attributes
 */
func (this *JsonEncoder) Encode(obj *Object) string {
	o := make(map[string]interface{})
	for _, a := range obj.Attributes {
		if a.Optional > 0 && obj.random.Float64() > a.Optional {
			continue
		}
		value := string(<-a.Channel)
		// if value is integer or float than parse it, otherwise use string
		if f, err := strconv.ParseFloat(value, 64); err == nil {
			o[a.Name] = f
		} else if i, err := strconv.Atoi(value); err == nil {
			o[a.Name] = i
		} else {
			o[a.Name] = value
		}
	}
	// convert the map into json
	result, err := json.Marshal(o)
	if err != nil {
		log.Println(err)
		return ""
	}
	return string(result)
}

/*
 * Generate infinitely JSON data
 */
func (this *JsonEncoder) Generator(obj *Object) <-chan []byte {
	channel := make(chan []byte)
	go func() {
		for {
			data := this.Encode(obj)
			channel <- []byte(data)
		}
	}()
	return channel
}

/*
 * Generate infinitely KV data
 */
func (this *KVEncoder) Generator(obj *Object) <-chan []byte {
	channel := make(chan []byte)
	go func() {
		for {
			data := this.Encode(obj)
			channel <- []byte(data)
		}
	}()
	return channel
}

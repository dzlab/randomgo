randomgo
==============
[![Build Status](https://travis-ci.org/dzlab/randomgo.png)](https://travis-ci.org/dzlab/randomgo)
[![License](http://img.shields.io/:license-mit-blue.svg)](http://doge.mit-license.org)

randomgo is a tool for generating complex objects with attributes of random values. It can be used for instance to create datasets for benchmarking databases or queue systems. It currently has support for:
* generate primitive types (int, float, string)
* pick randomly values from a collection or a file (e.g. list of cities name)
* optional attributes (i.e. attributes that appear for a given probability)
* output objects to formats like CSV/TSV or JSON
* ... more to come

### Installation
```go get github.com/dzlab/randomgo```

### Documentation
http://godoc.org/github.com/dzlab/randomgo

### Usage
Define a YML schema file that describe the generators to use, along with the target object (with possible optional attributes).
Example: `schema.yml` (more details under the examples folder)
```
---
- generators:
  - name: a
    type: increment
    initial: 0
  - name: b
    type: pick
    file: cities.txt
- schema:
  - name: a1
    value: $a
  - name: a2
    value: $b
    optional: 0.7
```
Then, parse this schema file and start receiving objects
```
p := random.NewParser()
object, err := p.Parse("schema.yml")
kv := random.NewKVEncoder("=", "&")
json := random.NewJsonEncoder()
for i := 0; i < 10; i++ {
  log.Println(i, ">", kv.Encode(object))
  log.Println(i, ">", json.Encode(object))
}
```

### Contribute
This tool is still under very active development. Any contribution is welcome.

Some planned features:

* Collections (e.g. arrays)
* Nested objects (i.e. objects inside an other object)
* Add more output types (e.g. avro)
* Constraints or dependable values (e.g. city depends on country)
* ...

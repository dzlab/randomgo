random.go
==============
[![Build Status](https://travis-ci.org/dzlab/random.go.png)](https://travis-ci.org/dzlab/random.go)

random.go is a tool for generating complex objects with attributes of random values. It can be used for instance to create datasets for benchmarking databases or queue systems. It currently has support for:
* generate primitive types (int, float, string)
* pick randomly values from a collection or a file (e.g. list of cities name)
* optional attributes (i.e. attributes that appear for a given probability)
* output objects to formats like CSV/TSV or JSON
* ... more to come

### Installation
```go get github.com/dzlab/random.go```

### Documentation
http://godoc.org/github.com/dzlab/random.go

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
for i := 0; i < 10; i++ {
  log.Println(i, ">", object.GetKV("=", "&"))
  log.Println(i, ">", object.GetJSON())
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

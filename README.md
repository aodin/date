# date [![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/aodin/date)
Golang's missing date package, including ranges

Date builds on Golang's `time.Time` package to provide a [ISO 8601](https://en.wikipedia.org/wiki/ISO_8601) Date type

Create a new `Date`:

```go
import "github.com/aodin/date"

func main() {
    march1st := date.New(2015, 3, 1)
    fmt.Println(march1st) // 2015-03-01
}
```

Parse a date or build it from a time:

```go
date.Parse("2015-03-01")
date.FromTime(time.Now())
```

Ranges, including `Union` and `Intersection` operations:

```go
date.NewRange(date.Today(), date.Today().AddDays(7))
```

```go
date.EntireYear(2014).Union(date.EntireYear(2015))
```

By default, the `Date` type uses the `time.UTC` location. It can be passed to functions requiring the `time.Time` type using the embedded `Time` field:

```go
jan1 := time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)
jan1.Before(march1.Time)
```

Happy Hacking!

aodin, 2015-2016

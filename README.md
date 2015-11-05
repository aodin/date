# date
Golang's missing date package, including ranges

Date builds on Golang's `time.Time` package to provide a ISO 8601 Date type

Create a new `Date`:

```go
import "github.com/aodin/date"

func main() {
    march1st := date.New(2015, 3, 1)
    fmt.Println(march1st) // 2015-03-01
}
```

Parse a date:

```go
Parse("2015-03-01")
```

Range operations:

```go
Range(Today(), Today.AddDays(1)).Intersection(Today())
```

```go
EntireMonth().Union(EntireYear())
```

By default, the `Date` type uses the `time.UTC` location. It can be passed to functions requiring the `time.Time` type using the embedded `Time` field:

```go
jan1 := time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)
jan1.Before(march1.Time)
```

Happy Hacking!

aodin, 2015

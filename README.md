# Project goutils
[![Go](https://github.com/m-szalik/goutils/actions/workflows/go.yml/badge.svg)](https://github.com/m-szalik/goutils/actions/workflows/go.yml)
[![REUSE status](https://api.reuse.software/badge/github.com/m-szalik/goutils)](https://api.reuse.software/info/github.com/m-szalik/goutils)
[![Go Reference](https://pkg.go.dev/badge/github.com/m-szalik/goutils.svg)](https://pkg.go.dev/github.com/m-szalik/goutils)

Simple and useful functions for Go. Requires Go 1.22 or newer.
Full API documentation: https://pkg.go.dev/github.com/m-szalik/goutils

## Package: goutils
```shell
go get github.com/m-szalik/goutils
```
### Converters
 * `BoolToStr(b bool, trueVal, falseVal string) string`
   returns trueVal or falseVal depending on b
 * `BoolTo[T any](b bool, trueVal, falseVal T) T`
   returns trueVal or falseVal depending on b
 * `HexToInt(hex string) (int, error)`
   converts a hexadecimal string (optional sign and `0x` prefix) to int
 * `ParseBool(str string) (bool, error)`
   true is one of "true", "1", "on", false is one of "false", "0", "off"
 * `ParseValue(str string) interface{}`
   converts a string to bool, nil, int64 or float64;
   the string itself is returned when no conversion applies
 * `AsFloat64(input any) (float64, error)`
   converts float32, float64, int, int32, int64, string, []byte and pointers to them to float64
 * `RoundFloat(val float64, precision uint) float64`
   rounds to the given number of decimal places

### Slices
 * `SliceIndexOf[T comparable](slice []T, e T) int`
   index of e in slice or -1 when not found
 * `SliceContains[T comparable](slice []T, e T) bool`
 * `SliceRemove[T comparable](slice []T, e any) ([]T, int)`
   new slice without e and the number of removed elements; the input slice is unchanged
 * `SliceMap[I any, O any](input []I, mapper func(I) O) []O`
   maps a slice to a slice of another type
 * `SlicesEq[T comparable](a, b []T) bool`
 * `DistrictValues[T comparable](input []T) []T`
   unique values in the order of first occurrence
 * `Filter`, `FindFirst`, `AllMatch`, `AnyMatch`, `CountMatch`
   (also as `SliceAllMatch`, `SliceAnyMatch`, `SliceCountMatch`) - predicate helpers
 * `IterateOver[T any](ctx context.Context, elements []T, callback func(index int, element T) error) error`
   calls callback for every element until it fails or ctx is done

Example:
```go
sliceOfStrings := goutils.SliceMap[int, string]([]int{2, 7, -11}, func(i int) string { return fmt.Sprint(i) })
```

### Maps
 * `MapKeys(m map[string]any) []string`
 * `MapMerge[K comparable, V any](conflictFunc func(key K, v0, v1 V) V, initMap map[K]V, maps ...map[K]V) map[K]V`
   merges maps into initMap, conflicts are resolved by conflictFunc
 * `MapMergeOverride[K comparable, V any](initMap map[K]V, maps ...map[K]V) map[K]V`
   merges maps into initMap, the incoming value wins

### Environment and process
 * `Env[T string | int | bool | float32 | float64 | time.Duration](name string, def T) T`
   environment variable converted to T, or def when missing or invalid
 * `EnvRequired[T ...](name string) T`
   like `Env` but panics when the variable is missing or invalid
 * `EnvInt(name string, def int) int` - deprecated, use `Env`
 * `ExitNow(code int, message string, messageArgs ...interface{})`
   logs the message, appends it to `TERMINATION_MESSAGE_PATH` (default `/dev/termination-log`) when that file exists, and exits
 * `ExitOnError(err error, code int)`
 * `ExitOnErrorf(err error, code int, message string, messageArgs ...interface{})`
 * `FileExists(path string) bool`, `DirExists(path string) bool`
 * `CloseQuietly(closer io.Closer)`

### Errors and channels
 * `NewJoinErrorHelper(errs ...error) *JoinErrorHelper`
   collects errors: `Append`, `ErrorsCount`, `Iterate` and `AsError`
   (nil, the single error, or `errors.Join` of all of them)
 * `SafeSend[T any](ch chan<- T, value T) error`
   sends to a channel and returns `ErrSafeSendChannelClosed` instead of panicking when it is closed

### Structs and reflection
 * `CmpWalkStructAreEqual(a, b interface{}) CmpError`
   deep comparison that reports the first mismatch together with its field path
 * `CopyStructAll(src, dst interface{}) error`,
   `CopyStructSelected(src, dst interface{}, selectedPaths ...string) error`,
   `CopyStructAllExcept(src, dst interface{}, excludedPaths ...string) error`,
   `CopyStruct(src, dst interface{}, acceptFunc AcceptFunc) error`
   copy non-zero values between two values of the same type; dst must be a pointer
 * `IterateDeep(element interface{}, callback IteratorCallback)`
   walks nested structs, maps, slices and pointers and calls back for every leaf value with its path

### Math
 * `RootN(x, n float64) float64`
   approximation of the n-th root of x

### Time
`TimeProvider` is an abstraction for `time.Now()` that allows independent testing:
```go
type TimeProvider interface {
	Now() time.Time
}
tpMock := goutils.NewMockTimeProvider() // tpMock.Add(delta) moves the clock
tpSystem := goutils.SystemTimeProvider()
```

`StopWatch` measures elapsed time:
```go
sw := goutils.NewStopWatch()
sw.Start()
// some task here
sw.Stop()
fmt.Printf("execution of a task took %s", sw)
```

`TimeCounter` accumulates the time between `Start()` and `Stop()` calls;
`Value()` returns the total and `Reset()` returns it and starts from zero.

### Other
 * `NewStringWriter() StringWriter` - an `io.Writer` whose `String()` returns everything written so far

## Package: collector
```shell
go get github.com/m-szalik/goutils/collector
```
Few implementations of collections including:
 * `NewRollingCollection` - collection that keeps the last N added elements.
 * `NewSimpleCollection` - collection that keeps all elements, slice that grows when needed.
 * `NewTimedCollection` - collection that keeps elements for the defined duration only.
 * `NewDataPointsCollector` - collection that can calculate Avg, Max or Min over a time window.
 * `NewSet` - set of unique elements.
 * `NewStack` - implementation of the stack data structure.

## Package: throttle
```shell
go get github.com/m-szalik/goutils/throttle
```
Few implementations of throttling:
 * `NewPeriodicThrottler` - forwards the latest event once per period.
 * `NewMinDelayThrottler` - forwards events no more often than the minimal delay, events in between are dropped.
 * `NewMinStable` - forwards a changed value only after it stayed unchanged for the given duration.

## Package: pubsub
```shell
go get github.com/m-szalik/goutils/pubsub
```
Simple Publish-Subscribe implementation based on channels.
Implementation allows to have multiple subscribers as well as multiple publishers.

[Example](./pubsub/example.go)

## Package: dbfile
```shell
go get github.com/m-szalik/goutils/dbfile
```
`NewKeyValueDBFile(file)` - minimal key-value store (`Put`, `Get`, `Remove`, `Keys`) persisted as a JSON file.

## License
Apache License 2.0, see [LICENSE](./LICENSE).
The repository follows the [REUSE](https://reuse.software) specification; copyright and license
information for every file is declared in [REUSE.toml](./REUSE.toml).

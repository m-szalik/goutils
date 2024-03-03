# Project goutils
[![Go](https://github.com/m-szalik/goutils/actions/workflows/go.yml/badge.svg)](https://github.com/m-szalik/goutils/actions/workflows/go.yml)
Simple and useful functions for go.

## Package: goutils
```shell
go get github.com/m-szalik/goutils 
```
### Converters
 * `BoolToStr(b bool, trueVal, falseVal string) string`  
    BoolToStr - return string for true or false bool value
 * `BoolTo[T interface{}](b bool, trueVal, falseVal T) T`  
    BoolTo - return T object for true or false bool value
 *  `HexToInt(hex string) (int, error)`  
    HexToInt convert hex representation to int
 * `ParseBool(str string) (bool, error)`  
    ParseBool - return bool. True is one of "true", "1", "on", false is one of "false", "0", "off"
 * `ParseValue(str string) interface{}`  
   ParseValue - converts string to one of int64, flot64, string, bool, nil.  
   If impossible to covert the same string is returned.
 * `AsFloat64(i interface{}) (float64, error)`  
   AsFloat64 - convert multiple types (float32,int,int32,int64,string,[]byte) to float64.
 * `RoundFloat(val float64, precision uint) float64`  
    RoundFloat round float64 number

### Slices
 * `SliceIndexOf(slice []interface{}, e interface{}) int`
   SliceIndexOf find an element in a slice.  
   Returns index of the element in a slice or -1 if not found.
 * `SliceContains(slice []interface{}, e interface{}) bool`
   SliceContains check if slice contains the element
 * `SliceRemove[T comparable](slice []T, e any) ([]T, int)`
    SliceRemove remove the element from a slice.
    Returns []T = new slice, int number of removed elements.
 * `SliceMap[I interface{}, O interface{}](inputData []I, mapper func(I) O) []O`
    Map slice to slice of different object.  
    Example:
```go
    sliceOfStrings := goutils.SliceMap[int, string]([]int{2, 7, -11}, func(i int) string { return fmt.Sprint(i) })
```

### Other utils
 * `ExitOnError(err error, code int, message string, messageArgs ...interface{})`
 * `Env(name string, def string) string`
 * `EnvInt(name string, def int) int`
 * `CloseQuietly(closer io.Closer)`
 * `Root(x, n float64) float64`
    Root - this function calculate approximation of n-th root of a number.
    
### TimeProvider
An abstraction for `time.Now()` that allows independent testing.
```go
type TimeProvider interface {
	Now() time.Time
}
```

Create TimeProvider:
```go
tpMock := goutils.NewMockTimeProvider()
tpSystem := goutils.SystemTimeProvider()
```

### StopWatch
```go
sw := goutils.NewStopWatch()
sw.Start()
// some task here
sw.Stop()
fmt.Printf("execution of a task took %s", sw)
```


## Package: collector
```shell
go get github.com/m-szalik/goutils/collector 
```
Few implementation of collections including:
 * `NewRollingCollection` - Collection that keeps last N added elements.
 * `NewSimpleCollection` - Collection that keeps all elements, slice that grows when needed.
 * `NewTimedCollection` - Collection that keeps elements for defined duration only.
 * `NewDataPointsCollector` - Collection that can calculate Avg, Max or Min over a time window.
 * `NewStack` - implementation of stack data structure.


## Package: throttle
```shell
go get github.com/m-szalik/goutils/throttle 
```
Few implementation throttling.

## Package: pubsub
```shell
go get github.com/m-szalik/goutils/pubsub 
```
Simple Publish-Subscribe implementation based on channels.


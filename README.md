# optional
Optional ("nullable") value library for Go

## Requirements
- Go 1.18 or newer

## Installation
```bash
go get github.com/heucuva/optional
```

## How to use
Import this library into your code:

```go
import "github.com/heucuva/optional"
```

Then, you can easily add an optional value like so:
```go
var testValue optional.Value[string]

func Example() {
    // testValue should be 'nil', as it's initially unset
    // the expected result here is that `value` is "" and `set` is false
    value, set := testValue.Get()
    
    // this will set the testValue variable to a new value
    // NOTE: setting an optional.Value to nil does not cause it to return `set` as false
    //       the preferred behavior for clearing a Value is to call Reset, as shown below
    testValue.Set("Hello world!")

    // testValue should now be set
    // the expected result here is that `value is "Hello world!" and `set` is true
    value, set = testValue.Get()

    // this will unset the testValue variable
    // essentially setting the value to 'nil'
    // do not confuse this with setting a Value to actual `nil`
    testValue.Reset()

    // testValue should be once again be 'nil', as it's been unset by the Reset function
    // the expected result here is that `value` is "" and `set` is false
    value, set := testValue.Get()
}
```

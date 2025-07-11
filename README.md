# durago
Implementation of ISO-8601 duration parsing

Inspired by [sosodev's implementation](https://github.com/sosodev/duration)

# Installation
`go get github.com/MeatAndBlood/durago`

# Purposes
I needed this parser for work purposes.  
sosodev's solution was not suitable for me, because it did not handle the parsing errors I needed and allows you to specify much more than I needed in the values.

# Accuracy and values
The package provides stronger validation than sosodev's solution, emphasizing restricting values ​​to integers except seconds.  
Duration values ​​with years or months will be converted with a slight inaccuracy, as the values ​​vary.
Similar to sosodev's solution, `2.628e+15` nanoseconds for a month and `3.154e+16` nanoseconds for a year are used.

# Usage
```golang
package main

import (
	"fmt"
	"time"

	"github.com/MeatAndBlood/durago"
)

func main() {
	d, err := durago.ParseDuration("P3Y6M2W4DT12H30M5S")
	if err != nil {
		panic(err)
	}

	fmt.Println(d.GetTimeDuration()) // 31104h30m5s

	d, err = durago.ParseDuration("PT12.5S")
	if err != nil {
		panic(err)
	}

	fmt.Println(d.GetTimeDuration() == time.Second*12+time.Millisecond*500) // true
}
```

# Restrictions
Anything larger than `P292Y5M2W5DT21H47M16.854775807S` will be converted incorrectly, as an int64 overflow will occur.
You will still be able to get the correct value from the `String` method, but the value in `GetTimeDuration` will be incorrect.

```golang
package main

import (
	"fmt"
	"math"
	"time"

	"github.com/MeatAndBlood/durago"
)

func main() {
	d := time.Duration(math.MaxInt64)
	fmt.Println(durago.FromTimeDuration(d).GetTimeDuration() == math.MaxInt64) // true

	dur, _ := durago.ParseDuration("P292Y5M2W5DT21H47M17S")
	fmt.Println(dur.GetTimeDuration()) // -2562047h47m16.709551616s

}
```
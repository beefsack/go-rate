go-rate-limiter
===============

[![Build Status](https://travis-ci.org/beefsack/go-rate-limiter.svg?branch=master)](https://travis-ci.org/beefsack/go-rate-limiter)

**go-rate-limiter** is a rate limiter designed for a range of use cases,
including server side spam protection and preventing saturation of APIs you
consume.

It is used in production at
[LangTrend](http://langtrend.com/l/Java,PHP,JavaScript) to adhere to the GitHub
API rate limits.

Examples
--------

### Simple rate limiting with blocking

This example demonstrates limiting the rate to 3 times per second.

```Go
package main

import (
	"fmt"
	"time"

	"github.com/beefsack/go-rate-limiter"
)

func main() {
	rl := rate_limiter.New(3, time.Second) // 3 times per second
	begin := time.Now()
	for i := 1; i <= 10; i++ {
		rl.Wait()
		fmt.Printf("%d started at %s\n", i, time.Now().Sub(begin))
	}
	// Output:
	// 1 started at 12.584us
	// 2 started at 40.13us
	// 3 started at 44.92us
	// 4 started at 1.000125362s
	// 5 started at 1.000143066s
	// 6 started at 1.000144707s
	// 7 started at 2.000224641s
	// 8 started at 2.000240751s
	// 9 started at 2.00024244s
	// 10 started at 3.000314332s
}
```

### Multi rate simiting with blocking

This example demonstrates combining rate limiters, one limiting at once per
second, the other limiting at 2 times per 3 seconds.

```Go
package main

import (
	"fmt"
	"time"

	"github.com/beefsack/go-rate-limiter"
)

func main() {
	begin := time.Now()
	rl1 := rate_limiter.New(1, time.Second)   // Once per second
	rl2 := rate_limiter.New(2, time.Second*3) // 2 times per 3 seconds
	for i := 1; i <= 10; i++ {
		rl1.Wait()
		rl2.Wait()
		fmt.Printf("%d started at %s\n", i, time.Now().Sub(begin))
	}
	// Output:
	// 1 started at 11.197us
	// 2 started at 1.00011941s
	// 3 started at 3.000105858s
	// 4 started at 4.000210639s
	// 5 started at 6.000189578s
	// 6 started at 7.000289992s
	// 7 started at 9.000289942s
	// 8 started at 10.00038286s
	// 9 started at 12.000386821s
	// 10 started at 13.000465465s
}
```

### Non-blocking rate limiting

This example demonstrates non-blocking rate limiting, such as would be used to
limit spam in a chat client.

```Go
package main

import (
	"fmt"
	"time"

	"github.com/beefsack/go-rate-limiter"
)

var rl = rate_limiter.New(3, time.Second) // 3 times per second

func say(message string) {
	if ok, remaining := rl.Try(); ok {
		fmt.Printf("You said: %s\n", message)
	} else {
		fmt.Printf("Spam filter triggered, please wait %s\n", remaining)
	}
}

func main() {
	for i := 1; i <= 5; i++ {
		say(fmt.Sprintf("Message %d", i))
	}
	// Output:
	// You said: Message 1
	// You said: Message 2
	// You said: Message 3
	// Spam filter triggered, please wait 999.980816ms
	// Spam filter triggered, please wait 999.976704ms
}
```

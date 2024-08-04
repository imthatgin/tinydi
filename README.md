# tinydi

## Usage
You can either create a new context, using `tinydi.New()`, or use `nil` in place of it, to use the default
context.

### Tags
Make sure to annotate your sub-dependencies with *\`di:"inject"\`*

### Code example
```go
package main

import "fmt"
import "github.com/imthatgin/tinydi"

type HelloService struct {
	World    *WorldService `di:"inject"`
	HelloMap map[string]string
}

func (s *HelloService) Hello() string {
	return fmt.Sprintf("Hello %s", s.World.World())
}

type WorldService struct {
	value   string
}

func (s *WorldService) World() string {
	return s.value
}

func main() {
    i := tinydi.New() // use a nil value in place of i if you wish to use the global, default context.
    
    tinydi.Add[HelloService](i, func(i *tinydi.Injector) *HelloService {
        return &HelloService{}
    })
    tinydi.Add[WorldService](i, func(i *tinydi.Injector) *WorldService {
        return &WorldService{
            value: "World",
        }
    })
    
    h := tinydi.MustGet[HelloService](i)
    
    h.Hello()
}

```
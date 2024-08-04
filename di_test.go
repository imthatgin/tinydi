package tinydi_test

import (
	"fmt"
	"github.com/imthatgin/tinydi"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

type HelloService struct {
	World    *WorldService `di:"inject"`
	HelloMap map[string]string
}

func (s *HelloService) Hello() string {
	return fmt.Sprintf("Hello %s", s.World.World())
}

type WorldService struct {
	value   string
	Numbers *NumberService
}

func (s *WorldService) World() string {
	return s.value
}

type NumberService struct{}

func (s *NumberService) Number() float64 {
	return rand.Float64()
}

func TestAdd(t *testing.T) {
	i := tinydi.New()

	tinydi.Add[HelloService](i, func(i *tinydi.Injector) *HelloService {
		return &HelloService{}
	})
	tinydi.Add[WorldService](i, func(i *tinydi.Injector) *WorldService {
		return &WorldService{
			value: "World",
		}
	})
	tinydi.Add[NumberService](i, func(i *tinydi.Injector) *NumberService {
		return &NumberService{}
	})

	h := tinydi.MustGet[HelloService](i)
	h2 := tinydi.MustGet[HelloService](i)
	assert.NotNil(t, h.World)
	assert.Nil(t, h.World.Numbers)
	assert.Nil(t, h.HelloMap)
	assert.NotSame(t, h, h2)
	assert.NotSame(t, h.World, h2.World)
	assert.Equal(t, "Hello World", h.Hello())
}

func TestAddSingleton(t *testing.T) {
	i := tinydi.New()
	tinydi.AddSingleton[HelloService](i, func(i *tinydi.Injector) *HelloService {
		return &HelloService{}
	})
	tinydi.Add[WorldService](i, func(i *tinydi.Injector) *WorldService {
		return &WorldService{
			value: "World",
		}
	})

	h := tinydi.MustGet[HelloService](i)
	h2 := tinydi.MustGet[HelloService](i)
	assert.NotNil(t, h.World)
	assert.Nil(t, h.HelloMap)
	assert.Same(t, h, h2)
	assert.Same(t, h.World, h2.World)
	assert.Equal(t, "Hello World", h.Hello())
}

func TestAddSingletonChild(t *testing.T) {
	i := tinydi.New()

	tinydi.Add[HelloService](i, func(i *tinydi.Injector) *HelloService {
		return &HelloService{}
	})
	tinydi.AddSingleton[WorldService](i, func(i *tinydi.Injector) *WorldService {
		return &WorldService{
			value: "World",
		}
	})

	h := tinydi.MustGet[HelloService](i)
	h2 := tinydi.MustGet[HelloService](i)
	h3 := tinydi.MustGet[HelloService](i)
	assert.NotNil(t, h.World)
	assert.Nil(t, h.World.Numbers)
	assert.Nil(t, h.HelloMap)
	assert.NotSame(t, h, h2)
	assert.NotSame(t, h, h3)
	assert.NotSame(t, h2, h3)
	assert.Same(t, h.World, h2.World)
	assert.Same(t, h2.World, h3.World)
	assert.Same(t, h.World, h3.World)
	assert.Equal(t, "Hello World", h.Hello())
}

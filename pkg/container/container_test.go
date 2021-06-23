package container

import (
	"testing"
)

type MockStruck struct {
	value int
}

func TestNewContainer(t *testing.T) {
	container := NewContainer()

	if len(container.registered) != 1 {
		t.Errorf("Container was created with wrong quantity of registered services in it")
	}

	if len(container.services) != 0 {
		t.Error("There are cached services in container")
	}

	containerFromService := container.Get("container.Container")
	if containerFromService == nil {
		t.Error("There is no container instance")
	}

	if len(container.services) != 1 {
		t.Error("Container instance hasn't cached")
	}

	if containerFromService != container {
		t.Error("Container returned wrong instance")
	}
}

func TestRegister(t *testing.T) {
	container := NewContainer()

	container.Register("MockService", func(c *Container) interface{} {
		return MockStruck{value: 1}
	})

	if len(container.registered) != 2 {
		t.Errorf("Wrong quantity of registered services in container")
	}

	if len(container.services) != 0 {
		t.Error("There are cached services in container")
	}
}

func TestGet(t *testing.T) {
	container := NewContainer()

	container.Register("MockService", func(c *Container) interface{} {
		return MockStruck{value: 1}
	})

	service := container.Get("MockService").(MockStruck)
	if service.value != 1 {
		t.Error("Container returned wrong instance of MockService")
	}

	if len(container.services) != 1 {
		t.Error("Container didn't cache the service")
	}

	service.value = 111
	serviceCached := container.Get("MockService").(MockStruck)
	if serviceCached.value != 1 {
		t.Error("Container returned non-cached instance of MockService")
	}
}

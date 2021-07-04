package container

import (
	"fmt"
	"testing"
)

type MockStruck struct {
	value int
}

type LastFatal struct {
	value string
}

func TestContainer(t *testing.T) {
	LastFatal := &LastFatal{}
	handler := func(format string, v ...interface{}) {
		LastFatal.value = fmt.Sprintf(format, v...)
	}

	container := NewContainer(handler)

	t.Run("Constructor", func(t *testing.T) {
		if len(container.registered) != 1 {
			t.Error("Container was created with wrong amount of registered services in it")
		}

		if len(container.services) != 0 {
			t.Error("There are cached services in container")
		}
	})

	t.Run("Check container", func(t *testing.T) {
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
	})

	t.Run("Register", func(t *testing.T) {
		container.Register("MockService", func(c *Container) interface{} {
			return &MockStruck{value: 1}
		})

		if len(container.services) != 1 {
			t.Error("Registered service is cached without calling")
		}
	})

	t.Run("Register duplicate", func(t *testing.T) {
		container.Register("MockService", func(c *Container) interface{} {
			return &MockStruck{value: 1}
		})

		if LastFatal.value != "Service with name `MockService` already registered in container" {
			t.Error("fatal error didn't occure")
		}
	})

	t.Run("Get", func(t *testing.T) {
		service := container.Get("MockService").(*MockStruck)
		if service.value != 1 {
			t.Error("Container returned wrong instance of MockService")
		}

		if len(container.services) != 2 {
			t.Error("Container didn't cache the service")
		}

		service.value = 111
		serviceCached := container.Get("MockService").(*MockStruck)
		if serviceCached.value == 1 {
			t.Error("Container returned non-cached instance of MockService")
		}
	})

	t.Run("Get non-registered", func(t *testing.T) {
		_ = container.Get("NonRegistered")
		if LastFatal.value != "Service with name `NonRegistered` hasn't registered in container" {
			t.Error("fatal error didn't occure")
		}
	})
}

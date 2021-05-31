package container

import (
	"fmt"
	"log"
)

type Container struct {
	registered map[string]func(c *Container) interface{}
	services   map[string]interface{}
}

func (c *Container) Register(service string, handler func(c *Container) interface{}) {
	if c.registered[service] != nil {
		log.Fatal(fmt.Sprintf("Service with name `%s` already registered in container", service))
	}

	c.registered[service] = handler
}

func (c *Container) Get(service string) interface{} {
	if c.services[service] != nil {
		return c.services[service]
	}

	if c.registered[service] == nil {
		log.Fatal(fmt.Sprintf("Service with name `%s` hasn't registered in container", service))
	}

	c.services[service] = c.registered[service](c)

	return c.services[service]
}

func NewContainer() *Container {
	instance := &Container{
		services:   make(map[string]interface{}),
		registered: make(map[string]func(c *Container) interface{}),
	}

	instance.Register("container.Container", func(c *Container) interface{} {
		return c
	})

	return instance
}

package container

type Container struct {
	fatalHanlder func(format string, v ...interface{})
	registered   map[string]func(c *Container) interface{}
	services     map[string]interface{}
}

// NewContainer - the container constructor
func NewContainer(fatalHanlder func(format string, v ...interface{})) *Container {
	instance := &Container{
		fatalHanlder: fatalHanlder,
		services:     make(map[string]interface{}),
		registered:   make(map[string]func(c *Container) interface{}),
	}

	instance.Register("container.Container", func(c *Container) interface{} {
		return c
	})

	return instance
}

// Register - the function that registers services in container
func (c *Container) Register(service string, resolver func(c *Container) interface{}) {
	if c.registered[service] != nil {
		c.fatalHanlder("Service with name `%s` already registered in container", service)
		return
	}

	c.registered[service] = resolver
}

// Get - the function that extracts services from container by its service name
func (c *Container) Get(service string) interface{} {
	if c.services[service] != nil {
		return c.services[service]
	}

	if c.registered[service] == nil {
		c.fatalHanlder("Service with name `%s` hasn't registered in container", service)
		return nil
	}

	c.services[service] = c.registered[service](c)

	return c.services[service]
}

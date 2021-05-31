package template

func GetRegisterServiceViaConstructorTemplate() string {
	return `
	compiledContainer.Register("%s", func(c *container.Container) interface{} {
		return %s(%s)
	})
	`
}

func GetRegisterServiceViaFactoryTemplate() string {
	return `
	compiledContainer.Register("%s", func(c *container.Container) interface{} {
		return c.Get("%s").(%s).%s(%s)
	})
	`
}

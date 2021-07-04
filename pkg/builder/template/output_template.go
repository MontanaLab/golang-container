package template

func GetOutputTemplate() string {
	return `package %s

import (
	%s
)

// =====================================================================================================
//
// BuildContainer is the autogenerated function, you must not change it. If you need another realisation
// of this function you have to:
// 1. change the container configuration if it neccessary
// 2. run rebuild cmd to rebuild this function
//
// Date of generation: %s
//
// =====================================================================================================
func BuildContainer(fatalHanlder func(format string, v ...interface{})) *container.Container {
	compiledContainer := container.NewContainer(fatalHanlder)
	%s
	return compiledContainer
}
	`
}

package template

func GetArgumentTemplate() string {
	return `c.Get("%s").(%s),`
}

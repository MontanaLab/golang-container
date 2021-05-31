##### Navigation:
* [Features](./FEATURES.md)
* [Getting Started](./GETTING_STARTED.md)
* [Configuration examples](./CONFIGURATION.md)
* [License](../LICENSE)

# Getting started
#### Installation
At first get container to you project:
```
go get github.com/MontanaLab/golang-container
```
#### Using builder
1. Create `cmd/build.go` command with same code:
```
package main

import (
	"flag"

	"github.com/MontanaLab/golang-container/pkg/builder"
)

func main() {
	var containerPath string
	var outputFile string

	flag.StringVar(&containerPath, "container", "config/container.json", "The path to container json file")
	flag.StringVar(&outputFile, "output", "internal/builder/builder.go", "The path to compiled container builder")

	flag.Parse()

	builder := builder.NewContainerBuilder()
	builder.Compile(containerPath, outputFile)
}
```
2. Create output folder: `mkdir -p internal/builder`
3. Create configuration file `config/container.json`. See the [Configuration examples](./CONFIGURATION.md).
4. Run command `go run cmd/build.go -config="config/container.json" -output="internal/builder/builder.go"`
5. Now you can see your compiled container in file `internal/builder/builder.go` and use function `BuildContainer()` from it that returns configurated container.
6. Now you can get registered services with `get` function. For example: `service := container.get("github.com/package/foo/bar/PublicServiceName").(*bar.PublicServiceName)`
7. Enjoy =)

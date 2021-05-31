package builder

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"time"

	"github.com/MontanaLab/golang-container/pkg/builder/dto"
	"github.com/MontanaLab/golang-container/pkg/builder/template"
)

type ContainerBuilder struct{}

func NewContainerBuilder() *ContainerBuilder {
	return &ContainerBuilder{}
}

type innerService struct {
	alias     string
	usage     string
	pkg       string
	factory   []string
	arguments []string
}

func (b *ContainerBuilder) Compile(configPath string, outputPath string) {
	file, err := os.Open(configPath)
	if err != nil {
		panic(err)
	}

	defer file.Close()
	bytes, _ := ioutil.ReadAll(file)

	var config dto.Config
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		panic(err)
	}

	services := make([]string, 0)
	servicesMapping := make(map[string]innerService)
	pkgAliases := make(map[string]string)
	for _, pkg := range config.Packages {
		alias := b.getPkgAlias(pkg.Name)
		pkgAliases[pkg.Name] = alias

		for _, service := range pkg.Services {
			serviceName := fmt.Sprintf("%s/%s", pkg.Name, service.Name)
			servicesMapping[serviceName] = b.getInnerService(
				config.Defaults,
				pkg,
				alias,
				service,
				serviceName,
			)

			services = append(services, serviceName)
		}
	}

	containerPkgName := "github.com/MontanaLab/golang-container/pkg/container"

	pkgAliases[containerPkgName] = ""
	servicesMapping[containerPkgName] = innerService{
		alias:     containerPkgName,
		pkg:       containerPkgName,
		usage:     "*container.Container",
		factory:   make([]string, 0),
		arguments: make([]string, 0),
	}

	registeredServices := b.getRegister(services, servicesMapping, pkgAliases)

	output := fmt.Sprintf(
		template.GetOutputTemplate(),
		"builder",
		b.getImport(pkgAliases, registeredServices),
		time.Now().String(),
		registeredServices,
	)

	err = ioutil.WriteFile(outputPath, []byte(output), 0644)
	if err != nil {
		panic(err)
	}
}

func (b *ContainerBuilder) getInnerService(defaults dto.Defaults, pkg dto.Package, pkgAlias string, service dto.Service, serviceName string) innerService {
	pointer := defaults.Pointer
	public := defaults.Public

	if service.Pointer != nil {
		pointer = *service.Pointer
	}

	if service.Public != nil {
		public = *service.Public
	}

	pointerValue := "*"
	if pointer == false {
		pointerValue = ""
	}

	serviceAlias := serviceName
	if public == false {
		serviceAlias = b.getHash(serviceAlias)
	}

	return innerService{
		alias:     serviceAlias,
		pkg:       pkg.Name,
		usage:     fmt.Sprintf("%s%s.%s", pointerValue, pkgAlias, service.Name),
		factory:   service.Factory,
		arguments: service.Arguments,
	}
}

func (b *ContainerBuilder) getPkgAlias(packageName string) string {
	return "pfx_" + b.getHash(packageName)
}

func (b *ContainerBuilder) getHash(value string) string {
	hash := sha1.New()
	hash.Write([]byte(value))

	return hex.EncodeToString(hash.Sum(nil))
}

func (b *ContainerBuilder) getImport(pkgAliases map[string]string, registeredServices string) string {
	pkgCounter := 0
	importString := ""
	for pkgName, alias := range pkgAliases {
		var value string
		if alias == "" {
			value = fmt.Sprintf(`"%s"`, pkgName)
		} else {
			value = fmt.Sprintf(`%s "%s"`, alias, pkgName)
			rule, _ := regexp.Compile(fmt.Sprintf("%s.", alias))
			if !rule.MatchString(registeredServices) {
				continue
			}
		}

		if pkgCounter != len(pkgAliases)-1 {
			value += "\n"
		}

		if pkgCounter > 0 {
			importString += "\t" + value
		} else {
			importString += value
		}

		pkgCounter++
	}

	return importString
}

func (b *ContainerBuilder) getRegister(services []string, servicesMapping map[string]innerService, pkgAliases map[string]string) string {
	if len(services) == 0 {
		return ""
	}

	register := ""
	for _, serviceName := range services {
		service := servicesMapping[serviceName]
		count := len(service.factory)
		if count == 0 || count > 2 {
			panic(fmt.Sprintf("Invalid factory value for service \"%s\"", serviceName))
		}

		switch count {
		case 1:
			if service.factory[0] == "" {
				panic(fmt.Sprintf("Factory first element cannot be empty for \"%s\"", serviceName))
			}

			pkgAlias := pkgAliases[service.pkg]
			constructor := fmt.Sprintf("%s.%s", pkgAlias, service.factory[0])

			registerString := b.getRegisterViaConstructor(
				constructor,
				service,
				servicesMapping,
			)

			register += registerString

			break
		case 2:
			if service.factory[0] == "" {
				panic(fmt.Sprintf("Factory first element cannot be empty for \"%s\"", serviceName))
			}

			if service.factory[0] == "" {
				panic(fmt.Sprintf("Factory second element cannot be empty for \"%s\"", serviceName))
			}

			if string([]rune(service.factory[0])[0]) != "@" {
				pkgAlias := pkgAliases[service.factory[1]]
				constructor := fmt.Sprintf("%s.%s", pkgAlias, service.factory[0])

				registerString := b.getRegisterViaConstructor(
					constructor,
					service,
					servicesMapping,
				)

				register += registerString
			} else {
				registerString := b.getRegisterViaFactory(
					service,
					string([]rune(service.factory[0])[1:]),
					service.factory[1],
					servicesMapping,
				)

				register += registerString
			}

			break
		}
	}

	return register
}

func (b *ContainerBuilder) getRegisterViaConstructor(
	constructor string,
	service innerService,
	servicesMapping map[string]innerService,
) string {
	return fmt.Sprintf(
		template.GetRegisterServiceViaConstructorTemplate(),
		service.alias,
		constructor,
		b.getArgumentsString(service.arguments, servicesMapping),
	)
}

func (b *ContainerBuilder) getRegisterViaFactory(
	service innerService,
	factory string,
	method string,
	servicesMapping map[string]innerService,
) string {
	factoryService, ok := servicesMapping[factory]
	if !ok {
		panic(fmt.Sprintf("There is no service \"%s\" in container", factory))
	}

	return fmt.Sprintf(
		template.GetRegisterServiceViaFactoryTemplate(),
		service.alias,
		factoryService.alias,
		factoryService.usage,
		method,
		b.getArgumentsString(service.arguments, servicesMapping),
	)
}

func (b *ContainerBuilder) getArgumentsString(
	arguments []string,
	servicesMapping map[string]innerService,
) string {
	argumentsString := ""
	for counter, argument := range arguments {
		argument = string([]rune(argument)[1:])
		if counter == 0 {
			argumentsString += "\n"
		}

		argumentService, ok := servicesMapping[argument]
		if !ok {
			panic(fmt.Sprintf("Arguments service \"%s\" does not exist in container", argument))
		}

		argumentsString += "\t\t\t" + fmt.Sprintf(
			template.GetArgumentTemplate(),
			argumentService.alias,
			argumentService.usage,
		)
		argumentsString += "\n"
	}

	if argumentsString != "" {
		argumentsString += "\t\t"
	}

	return argumentsString
}

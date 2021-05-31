##### Navigation:
* [Features](./FEATURES.md)
* [Getting Started](./GETTING_STARTED.md)
* [Configuration examples](./CONFIGURATION.md)
* [License](../LICENCE)

## Configuration examples
The structure of the configuration file must be like this:
```
{
    "defaults": {
        "pointer": true,
        "public": false
    },
    "packages": [
        {
            "name": "github.com/package/foo/bar",
            "services": [
                {
                    "name": "ServiceName",
                    "factory": ["ServiceNameConstructor"],
                    "public": true,
                    "pointer": true,
                    "arguments": [
                        "@github.com/package/foo/bar/SomeService"
                    ]
                }
            ]
        }
    ]
}
```
There are two sections in the configuration file:
1. `defaults` **_(required)_** - in this section you can set up default values for your container. There two options:
    - `pointer` `bool` **_(required)_** - this option is necessary to say to the container that all registered services are returning as a pointer by default
    - `public` `bool` **_(required)_** - the role of this option is to set up access to registered services. If you set up the value of it to `true`, the container will understand you want all services can be extracted from it with the container function `get`. In another hand, if you set up the value of it to `false` container will hide all services for function `get` and exclusion will be services that you will set up as `public: true` in the service configuration block.
    - Example:
    ```
    "defaults": {
        "pointer": true, // expects that all services will be returned as a pointer
        "public": false // will hide all services for function get
    },
    ```
2. `packages` **_(required)_** - the array of all packages that you want to register in your container. The structure of it is:
    - `name` `string` **_(required)_** - the name of the package that you want to add to the container
    - `services` `[]Service` **_(required)_** - the array of services from the current package that you want to register in the container. The structure of service you can find below.
    2.1. `Service` - the item of package services block
        - `name` `string` **_(required)_** - the name of registering service
        - `factory` `[]string` **_(required)_** - the way that container will create an instance of your service. There three ways to create a new instance of service. About these ways, you can read below.
        - `public` `bool` **_(optional)_** - this flag says that there are no ways to access the service with the container function `get`
        - `pointer` `bool` **_(optional)_** - this flag says that the service will be returned as a pointer from the container
        - `arguments` `[]string` **_(optional)_** - the list of services that you want to inject into your service (all these arguments will be passed in the constructor function). All you need is to write the full name of the service that you want to inject, but you must add the `@` symbol before. For example: `@github.com/package/foo/bar/ServiceName`

## Service create ways:
1. Create via a constructor that is located in the same package as the service. You can do it like that:
```
"services": [
    {
        "name": "ServiceName",
        "factory": ["ServiceNameConstructor"]
    }
]
```
2. Create via a constructor that is located in another package. You can do it like that:
```
"services": [
    {
        "name": "ServiceName",
        "factory": ["ServiceNameConstructor", "github.com/some/another/package/name"]
    }
]
3. Create via a factory that is already registered in the container. You must write this code:
```
"services": [
    {
        "name": "ServiceName",
        "factory": ["@github.com/some/another/package/name/ServiceFactory", "CreateFunction"]
    }
]

here CreateFunction - is the function that returns an instance of the desired service
```
    
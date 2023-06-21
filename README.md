# provider-simplejsonapp

`provider-simplejsonapp` is a minimal [Crossplane](https://crossplane.io/) Provider
that is meant to be used as a simplejsonapp for implementing new Providers. It comes
with the following features that are meant to be refactored:

- A `ProviderConfig` type that only points to a credentials `Secret`.
- A `MyType` resource type that serves as an example managed resource.
- A managed resource controller that reconciles `MyType` objects and simply
  prints their configuration in its `Observe` method.

## Developing

1. Use this repository as a simplejsonapp to create a new one.
1. Run `make submodules` to initialize the "build" Make submodule we use for CI/CD.
1. Rename the provider by running the follwing command:

```
  make provider.prepare provider={PascalProviderName}
```

4. Add your new type by running the following command:

```
make provider.addtype provider={PascalProviderName} group={group} kind={type}
```

5. Replace the _sample_ group with your new group in apis/{provider}.go
6. Replace the _mytype_ type with your new type in internal/controller/{provider}.go
7. Replace the default controller and ProviderConfig implementations with your own
8. Run `make reviewable` to run code generation, linters, and tests.
9. Run `make build` to build the provider.

Refer to Crossplane's [CONTRIBUTING.md] file for more information on how the
Crossplane community prefers to work. The [Provider Development][provider-dev]
guide may also be of use.

[CONTRIBUTING.md]: https://github.com/crossplane/crossplane/blob/master/CONTRIBUTING.md
[provider-dev]: https://github.com/crossplane/crossplane/blob/master/contributing/guide-provider-development.md

## Local testing

Below images show Crossplane installed on a local cluster and provider-simplejsonapp running against it and resource managed by APIs exposed by [SimpleJsonApp](https://github.com/chaitanyakolluru/go-works/tree/main/simpleJsonApp) managed using k8s api.

// to add documentation on spec
![record resource](./images/Record-resource.png)

By taking this Record spec and applying it:

```apiVersion: records.simplejsonapp.crossplane.io/v1alpha1
kind: Record
metadata:
  name: example-record
  namespace: default
spec:
  forProvider:
    name: "chaitanya"
    age: 22
    location: "something"
    designation: "something"
    todos:
      - "see if running in kind woks"
      - set it up with local k8s
```

### Record creation:

![record create](./images/record-create.png)

### Record updation:

![record update](./images/record-update.png)

![record update](./images/record-update-on-simplejsonapp.png)

### Record deletion:

![record deletion](./images/record-delete.png)

![record deletion](./images/record-delete-on-simplejsonapp.png)

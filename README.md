# provider-simplejsonapp

`provider-simplejsonapp` is a minimal [Crossplane](https://crossplane.io/) Provider,
using which one can express external json records as a k8s `record` resource.
This creates json records in an external application called [simple json app](https://github.com/chaitanyakolluru/go-works/tree/main/simpleJsonApp). Below installation process helps with setting up the api server and Crossplane, along with this provider locally.

## Installation

### Simple json app api server

- Clone [simple json app](https://github.com/chaitanyakolluru/go-works/tree/main/simpleJsonApp) locally and run

```
$ go run main.go
```

to install dependencies and run the app locally on port 8081. Refer to it's README.md for more details on how to authenticate to the app and make requests to it.

### Provider

- Install Crossplane:

  ```
  $ kc  apply -k https://github.com/crossplane/crossplane//cluster\?ref\=master
  ```

  You can also install it using helm as explained [here.](https://docs.crossplane.io/latest/software/install/)

- Create a secret with simple json app's jwt token like so:

  ```
  $ kc create secret generic simplejsonapp-secret --from-literal=token=<jwt token from simple json app's swagger ui here>
  ```

- Next, with your kube context pointing to cluster, or where ever you wish to install Crossplane and the provider, run:
  ```
  $ kubectl apply -f package/testYml/provider-simplejsonapp.yml
  ```
  which installs Provider `provider-simplejsonapp` and ProviderConfig `provider-simplejsonapp-config` (which references the secret we created in the previous step)

## Create and manage records

With Crossplane, Provider, ProviderConfig all installed we can now create `records` and manage them using Crossplane's Control plane.

Image showing the newly available `Record` managed resource

![record resource](./images/Record-resource.png)

Run this command to create a record MR:

```
$ kubctl apply -f package/testYml/record.yml
```

This file contains spec for the record Managed Resource. Details on some fields added below:

```
apiVersion: records.simplejsonapp.crossplane.io/v1alpha1
kind: Record
metadata:
  name: example-record
  namespace: default
spec:
  providerConfigRef:
    name: provider-simplejsonapp-config ## Every MR can point to the Provider Config its Provider could use to authenticate to the external system
  forProvider: ## contains object with details on record properties
    name: chaitanyaSomething
    age: 22
    location: something
    designation: something
    todos:
      - see if running in kind woks
      - set it up with local k8s
      - causing an update
      - testing provider image setup
```

Once you apply record.yml file, you can then see the record created as part of k8s api (Managed resource) and can also verify
record being created from simple json app's swagger page (external system)

### Record creation:

![record create](./images/record-create.png)

### Record deletion:

![record deletion](./images/record-delete.png)

![record deletion](./images/record-delete-on-simplejsonapp.png)

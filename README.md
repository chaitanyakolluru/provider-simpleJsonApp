# provider-simplejsonapp

[simple json app](https://gitlab.com/heb-engineering/teams/platform-engineering/gke-hybrid-cloud/kon/crossplane/simplejsonapp/simplejsonapp) is a simple api server which exposes some api endpoints and will act as the external resource for which we'll build a provider for.

`provider-simplejsonapp` is a minimal [Crossplane](https://crossplane.io/) Provider,
using which one can express external json records as a k8s `record` resource.
This creates json records in simple json app. Once the provider is installed, we will be able to manage records in the api server using Managed Reosurce `Record`.

Once that is done, we will devleop `Composition`, `XRD` and `Claim` to expose record object as a ConfigMap k8s resource and have a deployment use it, to prove we can build Composites of Managed Resources from more than one provider, with output (stored in MR's status) of one MR being fed into a ConfigMap MR and subsequently, to the deployment MR within the Composite reosurce.

# Installation

## Simple json app api server

- [simple json app](https://gitlab.com/heb-engineering/teams/platform-engineering/gke-hybrid-cloud/kon/crossplane/simplejsonapp/simplejsonapp) can be installed by this yaml:

  ```
  $ kubectl apply -f testYml/simplejsonapp.yml
  ```

  which installs simple json app api server as a deployment and exposes it with a k8s service resource.

## Provider

- Install Crossplane:

  ```
  $ helm repo add crossplane-stable https://charts.crossplane.io/stable
  $ helm repo update
  $ helm install crossplane \
  --namespace crossplane-system \
  --create-namespace crossplane-stable/crossplane
  ```

- Get auth token from running app:

  Port-forward into the app and get the auth token from Swagger UI.

  ```
  $ kc port-forward <simplejsonapp-pod> -n provider-simplejsonapp 8081:8081
  ```

  Now navigate to app's Swagger UI [here](http://localhost:8081/swagger/index.html) and execute request /auth/token to get your token

  ![swagger ui auth](./images/swagger-auth.png)

- Create a secret with simple json app's jwt token like so:

  ```
  $ kubectl create secret generic simplejsonapp-secret --from-literal=token=<jwt token from simple json app's swagger ui here>
  ```

- Install Simple json app CRDs for Record Managed Resource and other ProviderConfig CRDs:

  ```
  $ kubectl apply -f package/crds
  ```

- Next, with your kube context pointing to cluster, or where ever you wish to install Crossplane and the provider, run:

  ```
  $ kubectl apply -f testYml/provider-simplejsonapp.yml
  ```

  which installs Provider as a pod and sets up necessary permissions for it to be able to manage k8s resources and ProviderConfig `provider-simplejsonapp-config` (which references the secret we created in the previous step)

## Defining Compositions, Composite Resources, Composite Resource Deifnitions and Claims:

We have a Managed Resource called "Record" that's provided to us by `provider-simplejsonapp`. Let's setup composite resources over it now.

We will be setting up a `Composition` and a `Composite Resource Definition(XRD)` to setup a composite set of managed resources. This will be made
available to app teams as a `Composite Resource Claim`, created from a `XRD`. A `Claim` when applied against the k8s api creates a `Composite Reosource (XR)` and details on which resources make up the `XR` come from the `Composition`.

So essentially, we create the blueprint for all the resources that compose up to make a `Composite Resource(XR)` using our `Composition` and define a `XRD` which creates a `Claim`
that can be used to create `XR` in an application namespace, and therefore all the resources composed by the `XR`.

More details on the relationship of the resources and other terminology [here](https://docs.crossplane.io/latest/concepts/terminology/)

### Install

Install [provider-kubernetes](https://github.com/crossplane-contrib/provider-kubernetes) using below command

```

$ kubectl apply -f testYml/provider-kubernetes.yml

```

which creates k8s provider as seen below:

```

➤ kc get providers
NAME INSTALLED HEALTHY PACKAGE AGE
provider-kubernetes crossplane/provider-kubernetes:main 5s

```

### Apply Composition

### Apply Composite Resource Definition

### Apply Claim

### Verify if XRs and its Managed Resources are created

```

```

```

```

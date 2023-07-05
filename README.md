# provider-simplejsonapp

[simple json app](https://gitlab.com/heb-engineering/teams/platform-engineering/gke-hybrid-cloud/kon/crossplane/simplejsonapp/simplejsonapp) is a simple api server which exposes some api endpoints and will act as the external resource for which we'll build a Custom Crossplane provider.

`provider-simplejsonapp` is a minimal [Crossplane](https://crossplane.io/) Provider using which one can express external json records as a k8s `record` resource.
This creates json records in simple json app. Once the provider is installed, we will be able to manage records in the api server using a Managed Reosurce called `Record`.

Once that is done, we will devleop `Composition`, `XRD` and `Claim` to expose record object as a ConfigMap k8s resource and have a deployment use it, to prove we can build Composites of Managed Resources from more than one provider, with the output (stored in MR's status) of one MR being fed the definition of another MR from a different provider.

## Install Crossplane

```
$ helm repo add crossplane-stable https://charts.crossplane.io/stable
$ helm repo update
$ helm install crossplane \
--namespace crossplane-system \
--create-namespace crossplane-stable/crossplane
```

## Deploy and setup api server for our provider

- [simple json app](https://gitlab.com/heb-engineering/teams/platform-engineering/gke-hybrid-cloud/kon/crossplane/simplejsonapp/simplejsonapp) can be installed by this yaml:

  ```
  $ kubectl apply -f testYml/simplejsonapp.yml
  ```

  which installs simple json app api server as a deployment and exposes it with a k8s service resource.

- Get auth token from running app:

  Port-forward into the app and get the auth token from its Swagger UI.

  ```
  $ kc port-forward <simplejsonapp-pod> -n provider-simplejsonapp 8081:8081
  ```

  Now navigate to app's Swagger UI [here](http://localhost:8081/swagger/index.html) and execute request /auth/token to get your token

  ![swagger ui auth](./images/swagger-auth.png)

- Create a secret with simple json app's jwt token like so:

  ```
  $ kubectl create secret generic simplejsonapp-secret --from-literal=token=<jwt token from simple json app's swagger ui here> -n crossplane-system
  ```

## Provider-simplejsonapp installation

- Install Simple json app CRDs for Record Managed Resource and other ProviderConfig CRDs:

  ```
  $ kubectl apply -f package/crds
  ```

- Next, with your kube context pointing to cluster, or where ever you wish to install Crossplane and the provider, run:

  ```
  $ kubectl apply -f testYml/provider-simplejsonapp.yml
  ```

  which installs Provider as a pod and sets up necessary permissions for it to be able to manage k8s resources and ProviderConfig `provider-simplejsonapp-config` (which references the secret we created in the previous step)

## Provider-kubernetes installation

## Todo: Best way for provider to auth into the k8s cluster.

- Install provider-kubernetes using below command:

  ```
  $ kubectl apply -f testYml/provider-kubernetes.yml
  ```

  which installs the Provider and its associated ProviderConfig (used by Object MRs later for kubeconfig)

## Record MR from provider-simplejsonapp in action

We have a Managed Resource called "Record" that's provided to us by `provider-simplejsonapp`. Below section shows that you can using the Managed Resource `Record` create and manage those record objects:

Below command to create a Record MR:

```
$ kubectl apply -f testYml/record.yml
```

This creates a MR called record with properties as shown below in spec.forProvider:

```
---
apiVersion: records.simplejsonapp.crossplane.io/v1alpha1
kind: Record
metadata:
  name: example-record-new
  namespace: default
spec:
  providerConfigRef:
    name: provider-simplejsonapp-config
  forProvider:
    name: testingNow
    age: 22
    location: something
    designation: something
    todos:
      - see if running in kind woks
      - new record added
```

When this yaml is applied, it creates a MR and makes an api call to create the json record in the `simplejsonapp` deployed in `provider-simplejsonapp` namespace. It uses ProviderConfig `provider-simplejsonapp-config` as referenced by spec.providerConfigRef in the yaml, which contains the auth token
to be able to make calls against the simplejsonapp api server. Below are images showing successful creation of the record within the api server.

![record created](./images/record-created.png)

Below image showing simplejsonapp api service logging a post call to create the resource, and a json item added to it's local storage in data.json file:

![record verification](./images/record-verification.png)

So with Crossplane Custom Providers, we can easily incorporate an existing external system that has a RESTful API into a k8s resource called a Managed Resource. With Crossplane, one can go one step ahead and make `Compositions` of multiple resources and create schema to provide creation and management of all resources defined by the composite and provide the same for app teams for consumption.

## Defining Compositions, Composite Resources, Composite Resource Definitions and Claims:

We will be setting up a `Composition` and a `Composite Resource Definition(XRD)` to setup a composite set of managed resources. This will be made
available to app teams as a `Composite Resource Claim(or simply Claim)`, created from a `XRD`. A `Claim` when applied against the k8s api creates a `Composite Reosource (XR)`; the resources to be created when a `XR` is created come from `Composition`.

So essentially, we create the blueprint for all the resources that compose up to make a `Composite Resource(XR)` using our `Composition` and define a `XRD` which creates a `Claim` that can be used to create `XR` in an application namespace.

More details on the relationship of the resources and other terminology [here](https://docs.crossplane.io/latest/concepts/terminology/)

## Composition

```
$ kubectl apply -f testYml/composition.yml
```

Please refer to the yaml in question. This kind `Composition` essentially creates a blueprint for the `Composite Resources(XR)` that get created when one creates a `XR` using a `Claim`. The `XR` it creates is defined by `spec.compositeRef` object in it's definition. Note that the kind has a prefix 'X' to denote that it creates a `Composite Resource`. `spec.resources` becomes the rest of it's definition and as mentioned previously, it defines the resources that this `Composition` is made up of and dictates the resources that will be composed by the `XR` this `Composition` creates.

You see three resources in this example, one for the `Record` managed resource coming from our custom `provider-simplejsonapp` provider, and two coming from `provider-kubernetes` called `Object` managed resource, which is basically an abstraction within the provider to allow management of any k8s resource, in this case we use them to setup a `ConfigMap` to store our `Record` managed resource data and a `Deployment` resource to get record data from the `ConfigMap`.

Other things to note with the `Record` managed resource definition in the `Composition`:

- It references the same `ProviderConfig` that was used to create a single `Record` managed resource, and contains auth token to make api calls against the simple json app api server.
- Patches are used to make paches either to the `Composed Resource(Managed Resource)` by copying some attribute value from it's `XR` or to the `Composite Reosource(XR)` by copying some attribute value from it's managed Resource.

  This resource object has three patches, first two apply patches to the Managed Resource, the third applies a patch to the Composite Resource so it can later pass record details down to `ConfigMap` managed resource.

  First two update Record MR's name and namespace with labels added to the XR based on the Claim it gets created from, and the third patch sets `status.record` on the XR from record's `spec.froProvider` section.

  ```
  patches:
    - type: FromCompositeFieldPath
      fromFieldPath: metadata.labels[crossplane.io/claim-namespace]
      toFieldPath: metadata.namespace
      policy:
        fromFieldPath: Required
    - type: FromCompositeFieldPath
      fromFieldPath: spec.parameters
      toFieldPath: spec.forProvider
      policy:
        fromFieldPath: Required
    - type: ToCompositeFieldPath
      fromFieldPath: spec.forProvider
      toFieldPath: status.record
      policy:
        fromFieldPath: Required
  ```

Rest two Object MRs for k8s' Deployment and ConfigMap resources:

These are defined using `provider-kubernetes`'s `Object` managed resource and is simply an abstraction over an existing k8s resource and allows for a regular k8s like resource to be embedded into a `Crossplane Composition` and managed as a `XR`.

Things to note in second resource, `configmap`:

- Has a dummy `data` section under `spec.forProvider.manifest`, which will be soon replaced by actual data using patches.
- Speaking of, patches from index `2` show each record attribute being picked from `XR`'s `status.record` section and pached into `data` on the ConfigMap.
- Like the previous resource, patch indexes `0 & 1` patch the name and namespace of the resulting k8s ConfigMap resource.

Note patch types `FromCompositeFieldPath` to patch a `MR` using spec data from `XR`, `ToCompositeFieldPath` (used in Record resource definition here: spec.resources[0].patches[2]) to patch `XR` using spec data from `MR` and `CombineFromComposite`, generally used to combine spec attributes from multiple sources and create a composite patch. For more details on the types see [here.](https://docs.crossplane.io/v1.10/reference/composition/#patch-types)

Things to note with the third resource, `deployment`:

- This deployment simply exposes the contents of the ConfigMap from previous step and proves that we can access data obtained from a managed resource from a different provider (in this case `Record` MR from `provider-simplejsonapp`)
- As far as patches go, the only notable one appears as the second patch item patching configMap name to be used by the Deployment to refer while creating mount path in its pod.

To summarize, this yaml essentially creates the structure of `XRs` that will be created by referring to this Composition.

## Composite Resource Definition (XRD)

To apply XRD, see below:

```
$ kubectl apply -f testYml/xrd.yml
```

Please refer to the yaml in question.

This defines the schema of the `Claim` resources that we can use to create `XR`, and using it other managed resources. `spec.names` defines the kind of
`XRs` that get created when one uses this `XRD` to create a claim with. `spec.claimNames` defines the claim kind that will be created by Crossplane once this yaml is applied.

`spec.versions` list can contain more than one version of this `XRD` definition; in our case, we are using `v1alpha1`. What follows is a schema definition
for the spec that `Claims` resource can use to define `XR` and status to store `Record` managed resource data, in `spec.versions[0].schema.openAPIV3Schema.properties.spec` and `spec.versions[0].schema.openAPIV3Schema.properties.status`

The properties for both spec and status in this case has been chosen to be same, so as to store record attributes in `XR`'s status and propagate them to `ConfigMap` managed resource and contain a spec as shown below:

```
name:
  type: string
age:
  type: integer
designation:
  type: string
location:
  type: string
todos:
  type: array
  items:
    type: string
```

When the xrd is applied, we see a new k8s resource called `JsonApp` being created, for more details of which see next section.

## Composite Resource Claim (or Claim)

To apply a Claim, run below command:

```
$ kubectl apply -f testYml/claim.yml
```

yaml shown below:

```
apiVersion: simplejsonapp.crossplane.io/v1alpha1
kind: JsonApp
metadata:
  name: example-json-app
spec:
  parameters:
    name: testingNow
    age: 22
    location: something
    designation: something
    todos:
      - see if running in kind woks
      - new record added
  compositionRef:
    name: simplejsonapp-composition
```

As said previously, once an `XRD`is applied, it creates a new kind defined by `spec.claimNames.kind` in `XRD` definition (in our case `JsonApp`), which are namespaced resources, unlike `XRs` and can act as namespaced claims we can provide to consumers of Crossplane XRs. For a user of Crossplane, they are only required to know how to use a Claim and apply it to create XRs of their choosing.

In our case, once we apply the Claim, we see a XR corresponding to the Claim be created, and managed resources, all created and managed as part of this XR, with managed resource's name and namespaces being the same as the name and namespace of the Claim being used.

Below images show all resources being created:

### `JsonApp` claim being created:

![jsonapp claim](./images/jsonapp-claim.png)

### `XJsonApp` XR being created:

![jsonapp claim](./images/jsonapp-claim.png)

### `Record` MR from provider-simplejsonapp being created:

![record mr](./images/record-mr.png)

### `Object` MRs from provider-k8s being created storing definition for Deployment and ConfigMap:

![object mrs](./images/both-object-mrs.png)

### ConfigMap containing record attributes:

![cm-record-data](./images/cm-with-record-data.png)

### Deployment with pod containing record attributes:

![dep-record-data](./images/deploy-with-pod-containing-record-data.png)

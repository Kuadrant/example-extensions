# example-extensions

A collection of example [Kuadrant](https://kuadrant.io) policy extensions,
showing how to implement, build, and deploy a Kuadrant extension using the
[Extension SDK](https://github.com/Kuadrant/kuadrant-operator/blob/main/doc/extensions/extension-sdk-developer-guide.md)
provided by `kuadrant-operator`.

A Kuadrant extension defines a new policy CRD (e.g. `ThreatPolicy`) that
targets Gateway API resources (`Gateway`, `HTTPRoute`) and, on reconcile,
registers HTTP request/response behavior with the Kuadrant operator over
gRPC. The extension runs as its own controller process, separate from the
operator itself.

Extensions in this repo are meant to be read as reference material - each one
is a small, complete, runnable example rather than a library.

## Repository layout

Each extension lives in its own top-level directory and is a standalone Go
module:

```
example-extensions/
└── threat-policy/       # example extension
    ├── go.mod
    ├── main.go
    ├── Dockerfile
    ├── Makefile
    ├── api/v1alpha1/     # CRD Go types
    ├── internal/controller/ # reconciler
    └── config/
        ├── crd/          # generated CRD manifests
        ├── rbac/         # ClusterRole for the extension's own resources
        └── deploy/       # kustomize bases for deploying the extension
```

## Anatomy of an extension

An extension is laid out like an ordinary
[kubebuilder](https://book.kubebuilder.io/)/
[controller-runtime](https://github.com/kubernetes-sigs/controller-runtime)
project, and the pieces will look familiar if you've worked with one:

- **`api/v1alpha1/`** - the policy CRD's Go types and scheme registration
  boilerplate, generated with `controller-gen` the same way as any kubebuilder
  API package.
- **`internal/controller/`** - the reconciler for the policy type.
- **`config/crd/`** and **`config/rbac/`** - generated CRD and RBAC
  manifests, produced with `make manifests` as usual.
- **`main.go`** - wires up the scheme and starts the controller.
- **`config/deploy/`** - kustomize bases for running the extension (see
  Deploying, below).
- **`Dockerfile`** / **`Makefile`** - build the extension binary/image.

Where it differs from a plain kubebuilder project is in how the reconciler is
built and how it talks to Kuadrant: rather than a `controller-runtime.Manager`
reconciling directly against the cluster, the reconciler is wired up through
`kuadrant-operator`'s [Extension SDK](https://github.com/Kuadrant/kuadrant-operator/blob/main/doc/extensions/extension-sdk-developer-guide.md)
(source: [`pkg/extension`](https://github.com/Kuadrant/kuadrant-operator/tree/main/pkg/extension)),
which supplies the `Reconcile` entrypoint and a `KuadrantCtx` for registering
gRPC-backed action methods and request/response pipelines against the target
`Gateway`/`HTTPRoute`. See the SDK developer guide for details -
`threat-policy`'s `internal/controller/threatpolicy_reconciler.go` is a
worked example of using it.

This repo isn't a framework you install a dependency on - it's reference
material. To build your own extension, copy `threat-policy/` (or the layout
above) into your own project and adapt it.

## Deploying an extension

An extension runs as its own `Deployment`, with its own `ServiceAccount`, and
authenticates to the Kuadrant operator's extension gRPC service using a
projected `ServiceAccount` token.

Prerequisites:

- A cluster with the Kuadrant operator v1.6.0 or newer installed and running
  in the `kuadrant-system` namespace, with its extensions gRPC service
  reachable at `kuadrant-operator-extensions.kuadrant-system.svc:50052`.
- Gateway API CRDs installed.

Steps (from an extension's directory, e.g. `threat-policy/`):

1. Build and publish the extension image (or load it into the cluster
   directly, e.g. `kind load docker-image` for a local kind cluster):

   ```sh
   make docker-build
   ```

2. Apply the deploy overlay:

   ```sh
   kubectl apply -k config/deploy/standalone
   ```

   This creates, in order:
   - the extension's namespace (e.g. `threat-policy-system`) and
     `ServiceAccount`;
   - the extension's own CRD (`config/crd`) and `ClusterRole`
     (`config/rbac`), bound to that `ServiceAccount`;
   - a `ClusterRole`/`ClusterRoleBinding` granting the `register` verb on
     `policyregistrations` for this policy kind, so the operator accepts the
     extension registering itself;
   - the extension `Deployment`, which mounts a projected `ServiceAccount`
     token (audience `kuadrant-extensions`) and points
     `KUADRANT_EXTENSION_ADDRESS` at the operator's extensions service;
   - a `NetworkPolicy` in `kuadrant-system` allowing ingress from the
     extension's pods to the operator's controller-manager on port `50052`.

3. Create an instance of the policy CRD targeting a `Gateway` or
   `HTTPRoute`, and check its `status.conditions` for `Accepted`/`Enforced`.

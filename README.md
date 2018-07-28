# DC/OS Network object controller

This custom controller is able to apply network policies from different plugin vendors on DC/OS. At the moment it is a PoC that only works with Project Calico.
]
This example is based on Kubernetes [apiextensions-apiserver](https://github.com/kubernetes/apiextensions-apiserver) example  

## Organization

* crd      - define and register our CRD class
* client   - client library to create and use our CRD (CRUD)
* kube-crd - main part, demonstrate how to create, use, and watch our CRD
* examples - JSON files with custum object examples
* deployment - files for deploying packages and services on DC/OS

## Running

```
# assumes you have a working kubeconfig, not required if operating in-cluster
go run *.go -kubeconf=$HOME/.kube/config
```

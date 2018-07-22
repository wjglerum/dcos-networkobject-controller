package plugins

import "github.com/wjglerum/kube-crd/crd"

type PluginInterface interface {
	ListPolicies() ([]crd.SecurityPolicy, error)
	AddPolicy(policy crd.SecurityPolicy) (crd.SecurityPolicy, error)
	DeletePolicy(name string) error
}

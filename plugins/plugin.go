package plugins

import "github.com/wjglerum/kube-crd/crd"

type PluginInterface interface {
	ListPolicies() ([]crd.NetworkPolicy, error)
	AddPolicy(policy crd.NetworkPolicy) (crd.NetworkPolicy, error)
	DeletePolicy(name string) error
}

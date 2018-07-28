package plugins

import (
	"github.com/projectcalico/libcalico-go/lib/client"
	"github.com/projectcalico/libcalico-go/lib/api"

	"github.com/wjglerum/kube-crd/crd"
	"github.com/projectcalico/libcalico-go/lib/numorstring"
)

// CalicoPlugin implements PluginInterface
type CalicoPlugin struct{}

var newClient, _ = client.NewFromEnv()
var policies = newClient.Policies()

func (c CalicoPlugin) ListPolicies() ([]crd.NetworkPolicy, error) {
	policies, err := policies.List(api.PolicyMetadata{})
	return convertPoliciesTo(policies), err
}

func (c CalicoPlugin) AddPolicy(policy crd.NetworkPolicy) (crd.NetworkPolicy, error) {
	p, err := policies.Create(convertPolicy(policy))
	return convertPolicyTo(p), err
}

func (c CalicoPlugin) DeletePolicy(name string) error {
	meta := api.PolicyMetadata{
		Name: name,
	}
	return policies.Delete(meta)
}

func convertPolicy(policy crd.NetworkPolicy) *api.Policy {
	protocol := numorstring.ProtocolFromString(policy.Port[0].Protocol)
	p := api.NewPolicy()
	// TODO fix for array of rules
	rule := api.Rule{
		Action:    "deny",
		IPVersion: nil,
		Protocol: &protocol,
		Destination: api.EntityRule{
			Ports:    []numorstring.Port{numorstring.SinglePort(uint16(policy.Port[0].Port))},
		},
	}
	p.Spec = api.PolicySpec{
		IngressRules: []api.Rule{rule},
		EgressRules:  []api.Rule{rule},
		Selector:     policy.Selector[0].Matcher,
		Types:        []api.PolicyType{api.PolicyTypeIngress},
	}
	p.Metadata.Name = policy.Name
	return p
}

func convertPolicyTo(policy *api.Policy) crd.NetworkPolicy {
	return crd.NetworkPolicy{
		Type: policy.Kind,
		Name: policy.Metadata.Name,
		Selector: []crd.Selector{{
			Type:    "label",
			Matcher: policy.Spec.Selector,
		}},
		// TODO Fix for entire array of ports
		Port: []crd.Port{{
			Protocol: policy.Spec.EgressRules[0].Protocol.StrVal,
			Port:     int(policy.Spec.EgressRules[0].Destination.Ports[0].MaxPort),
		}},
	}
}

func convertPoliciesTo(policies *api.PolicyList) []crd.NetworkPolicy {
	newPolicies := make([]crd.NetworkPolicy, len(policies.Items))
	for i, v := range policies.Items {
		newPolicies[i] = convertPolicyTo(&v)
	}
	return newPolicies
}

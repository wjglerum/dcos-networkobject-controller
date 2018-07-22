// Copyright (c) 2017 Tigera, Inc. All rights reserved.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package model_test

import (
	. "github.com/projectcalico/libcalico-go/lib/backend/model"

	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/projectcalico/libcalico-go/lib/net"
)

var _ = DescribeTable(
	"key parsing",
	func(strKey string, expected Key, shouldFail bool) {
		key := KeyFromDefaultPath(strKey)
		if shouldFail {
			Expect(key).To(BeNil())
		} else {
			Expect(key).To(Equal(expected))
			serialized, err := KeyToDefaultPath(expected)
			Expect(err).ToNot(HaveOccurred())
			Expect(serialized).To(Equal(strKey))
		}
	},
	Entry(
		"profile rules with a /",
		"/calico/v1/policy/profile/foo%2fbar/rules",
		ProfileRulesKey{ProfileKey: ProfileKey{Name: "foo/bar"}},
		false,
	),
	Entry(
		"profile tags with a /",
		"/calico/v1/policy/profile/foo%2fbar/tags",
		ProfileTagsKey{ProfileKey: ProfileKey{Name: "foo/bar"}},
		false,
	),
	Entry(
		"profile labels with a /",
		"/calico/v1/policy/profile/foo%2fbar/labels",
		ProfileLabelsKey{ProfileKey: ProfileKey{Name: "foo/bar"}},
		false,
	),
	Entry(
		"policy with a /",
		"/calico/v1/policy/tier/default/policy/biff%2fbop",
		PolicyKey{Name: "biff/bop"},
		false,
	),
	Entry(
		"workload with a /",
		"/calico/v1/host/foobar/workload/open%2fstack/work%2fload/endpoint/end%2fpoint",
		WorkloadEndpointKey{
			Hostname:       "foobar",
			OrchestratorID: "open/stack",
			WorkloadID:     "work/load",
			EndpointID:     "end/point",
		},
		false,
	),
	Entry(
		"host endpoint with a /",
		"/calico/v1/host/foobar/endpoint/end%2fpoint",
		HostEndpointKey{
			Hostname:   "foobar",
			EndpointID: "end/point",
		},
		false,
	),
	Entry(
		"host IP",
		"/calico/v1/host/foobar/bird_ip",
		HostIPKey{Hostname: "foobar"},
		false,
	),
	Entry(
		"IP pool",
		"/calico/v1/ipam/v4/pool/10.0.0.0-8",
		IPPoolKey{CIDR: mustParseCIDR("10.0.0.0/8")},
		false,
	),
	Entry(
		"poorly formated IP pool",
		"/calico/v1/ipam/v4/pool/577559",
		nil,
		true,
	),
	Entry(
		"global felix config",
		"/calico/v1/config/foo",
		GlobalConfigKey{Name: "foo"},
		false,
	),
	Entry(
		"host config",
		"/calico/v1/host/hostname/config/foo",
		HostConfigKey{Hostname: "hostname", Name: "foo"},
		false,
	),
	Entry(
		"ready flag",
		"/calico/v1/Ready",
		ReadyFlagKey{},
		false,
	),
)

func mustParseCIDR(s string) net.IPNet {
	_, ipNet, err := net.ParseCIDR(s)
	if err != nil {
		panic(err)
	}
	return *ipNet
}

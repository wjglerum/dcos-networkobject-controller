/*
Copyright 2016 Iguazio Systems Ltd.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"fmt"
	"time"

	"github.com/wjglerum/kube-crd/client"
	"github.com/wjglerum/kube-crd/crd"

	apiextcs "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"flag"
	"github.com/google/go-cmp/cmp"
	"github.com/wjglerum/kube-crd/plugins"
)

// return rest config, if path not specified assume in cluster config
func GetClientConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	return rest.InClusterConfig()
}

func main() {

	kubeconf := flag.String("kubeconf", "admin.conf", "Path to a kube config. Only required if out-of-cluster.")
	flag.Parse()

	config, err := GetClientConfig(*kubeconf)
	if err != nil {
		panic(err.Error())
	}

	// create clientset and create our CRD, this only need to run once
	clientset, err := apiextcs.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// note: if the CRD exist our CreateCRD function is set to exit without an error
	err = crd.CreateCRD(clientset)
	if err != nil {
		panic(err)
	}

	// Wait for the CRD to be created before we use it (only needed if its a new one)
	time.Sleep(3 * time.Second)

	// Create a new clientset which include our CRD schema
	crdcs, scheme, err := crd.NewClient(config)
	if err != nil {
		panic(err)
	}

	// Create a CRD client interface
	crdclient := client.CrdClient(crdcs, scheme, "default")

	// Create a new Example object and write to k8s
	port := crd.Port{
		Protocol: "tcp",
		Port:     80,
	}

	selector := crd.Selector{
		Type:    "label",
		Matcher: "webserver",
	}

	rule := crd.NetworkPolicy{
		Type:     "ingress",
		Name:     "test",
		Selector: []crd.Selector{selector},
		Port: []crd.Port{port},
	}

	example := &crd.NetworkObject{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:   "example123",
			Labels: map[string]string{"mylabel": "test"},
		},
		VirtualNetwork: crd.VirtualNetwork{
			Name:      "test-network",
			Namespace: "default",
			Driver:    "calico",
			Subnet:    []string{"192.168.1.0/24"},
			Service:   []string{},
			Policy:    []crd.NetworkPolicy{rule},
		},
		NetworkDriver:  []crd.NetworkDriver{},
		NetworkService: []crd.NetworkService{},
	}

	result, err := crdclient.Create(example)
	if err == nil {
		fmt.Printf("CREATED: %#v\n", result)
	} else if apierrors.IsAlreadyExists(err) {
		fmt.Printf("ALREADY EXISTS: %#v\n", result)
	} else {
		panic(err)
	}

	// List all Example objects
	items, err := crdclient.List(meta_v1.ListOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("List:\n%s\n", items)

	calico := &plugins.CalicoPlugin{}
	// Example Controller
	// Watch for changes in Example objects and fire Add, Delete, Update callbacks
	_, controller := cache.NewInformer(
		crdclient.NewListWatch(),
		&crd.NetworkObject{},
		time.Minute*10,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				add(obj, calico)
			},
			DeleteFunc: func(obj interface{}) {
				delete(obj, calico)
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				update(oldObj, newObj, calico)
			},
		},
	)

	stop := make(chan struct{})
	go controller.Run(stop)

	// Wait forever
	select {}
}

func add(obj interface{}, plugin *plugins.CalicoPlugin) {
	nom := obj.(*crd.NetworkObject)
	pl, err := plugin.AddPolicy(nom.VirtualNetwork.Policy[0])
	fmt.Printf("add: %s \n", obj)
	fmt.Print(pl)
	fmt.Print(err)
}

func delete(obj interface{}, plugin *plugins.CalicoPlugin) {
	err := plugin.DeletePolicy(obj.(*crd.NetworkObject).VirtualNetwork.Policy[0].Name)
	fmt.Printf("delete: %s \n", obj)
	fmt.Print(err)
}

func update(oldObj, newObj interface{}, plugin *plugins.CalicoPlugin) {
	fmt.Printf("Update old: %s \n      New: %s\n", oldObj, newObj)
	fmt.Printf("Equal: %t\n", cmp.Equal(oldObj, newObj))
	fmt.Printf("Diff: %s\n", cmp.Diff(oldObj, newObj))
}

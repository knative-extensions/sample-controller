/*
Copyright 2019 The Knative Authors

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
	// The set of controllers this controller process runs.
	"flag"
	"log"

	"k8s.io/client-go/discovery"
	"knative.dev/sample-controller/pkg/apis/samples/install"
	"knative.dev/sample-controller/pkg/client/clientset/internalversion/scheme"
	"knative.dev/sample-controller/pkg/reconciler/addressableservice"

	// This defines the shared main for injected controllers.
	"knative.dev/pkg/injection/sharedmain"
	"knative.dev/pkg/signals"
)

func main() {
	var (
		masterURL  = flag.String("master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
		kubeconfig = flag.String("kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	)

	flag.Parse()

	// Set up signals so we handle the first shutdown signal gracefully.
	ctx := signals.NewContext()
	cfg, err := sharedmain.GetConfig(*masterURL, *kubeconfig)
	if err != nil {
		log.Fatal("Error building kubeconfig:", err)
	}

	dclient, err := discovery.NewDiscoveryClientForConfig(cfg)
	if err != nil {
		panic(err)
	}

	install.DiscoverAndUpdateVersionPriority(dclient, scheme.Scheme)

	sharedmain.MainWithConfig(ctx, "controller", cfg,
		addressableservice.NewController,
	)
}

module knative.dev/sample-controller

go 1.13

require (
	github.com/gobuffalo/envy v1.7.1 // indirect
	github.com/google/licenseclassifier v0.0.0-20200402202327-879cb1424de0
	github.com/grpc-ecosystem/grpc-gateway v1.12.2 // indirect
	go.uber.org/zap v1.10.0
	gomodules.xyz/jsonpatch/v2 v2.1.0 // indirect
	k8s.io/api v0.17.2
	k8s.io/apiextensions-apiserver v0.17.2 // indirect
	k8s.io/apimachinery v0.17.2
	k8s.io/client-go v11.0.1-0.20190805182717-6502b5e7b1b5+incompatible
	k8s.io/code-generator v0.18.0
	k8s.io/kube-openapi v0.0.0-20200121204235-bf4fb3bd569c
	knative.dev/pkg v0.0.0-20200410152005-2a1db869228c
	knative.dev/test-infra v0.0.0-20200427164251-3205a40d6171
)

replace (
	github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.2
	k8s.io/api => k8s.io/api v0.16.4
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.16.4
	k8s.io/apimachinery => k8s.io/apimachinery v0.16.4
	k8s.io/client-go => k8s.io/client-go v0.16.4
	k8s.io/code-generator => k8s.io/code-generator v0.16.4
	knative.dev/pkg => github.com/chizhg/pkg v0.0.0-20200420011907-9117cd5cf224
)

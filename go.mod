module knative.dev/sample-controller

go 1.13

require (
	contrib.go.opencensus.io/exporter/ocagent v0.6.0 // indirect
	contrib.go.opencensus.io/exporter/prometheus v0.1.0 // indirect
	contrib.go.opencensus.io/exporter/stackdriver v0.13.1 // indirect
	github.com/gobuffalo/envy v1.7.1 // indirect
	github.com/google/licenseclassifier v0.0.0-20200402202327-879cb1424de0
	github.com/grpc-ecosystem/grpc-gateway v1.12.2 // indirect
	github.com/markbates/inflect v1.0.4 // indirect
	github.com/openzipkin/zipkin-go v0.2.2 // indirect
	go.uber.org/zap v1.10.0
	gomodules.xyz/jsonpatch/v2 v2.1.0 // indirect
	k8s.io/api v0.17.2
	k8s.io/apiextensions-apiserver v0.17.2 // indirect
	k8s.io/apimachinery v0.17.2
	k8s.io/client-go v11.0.1-0.20190805182717-6502b5e7b1b5+incompatible
	k8s.io/code-generator v0.18.0
	k8s.io/kube-openapi v0.0.0-20200121204235-bf4fb3bd569c
	knative.dev/pkg v0.0.0-20200427190051-6b9ee63b4aad
	knative.dev/test-infra v0.0.0-20200429211942-f4c4853375cf
)

replace (
	github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.2
	k8s.io/api => k8s.io/api v0.16.4
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.16.4
	k8s.io/apimachinery => k8s.io/apimachinery v0.16.4
	k8s.io/client-go => k8s.io/client-go v0.16.4
	k8s.io/code-generator => k8s.io/code-generator v0.16.4
)

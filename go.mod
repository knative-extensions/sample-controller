module knative.dev/sample-controller

go 1.14

require (
	contrib.go.opencensus.io/exporter/stackdriver v0.13.1 // indirect
	github.com/google/licenseclassifier v0.0.0-20200402202327-879cb1424de0
	github.com/grpc-ecosystem/grpc-gateway v1.12.2 // indirect
	go.uber.org/zap v1.14.1
	k8s.io/api v0.17.6
	k8s.io/apimachinery v0.17.6
	k8s.io/client-go v11.0.1-0.20190805182717-6502b5e7b1b5+incompatible
	k8s.io/code-generator v0.18.0
	k8s.io/kube-openapi v0.0.0-20200410145947-bcb3869e6f29
	knative.dev/pkg v0.0.0-20200624210428-eb05e8dd5b5b
	knative.dev/test-infra v0.0.0-20200625142528-763d57db2487
)

replace (
	github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.2
	k8s.io/api => k8s.io/api v0.17.6
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.17.6
	k8s.io/apimachinery => k8s.io/apimachinery v0.17.6
	k8s.io/client-go => k8s.io/client-go v0.17.6
	k8s.io/code-generator => k8s.io/code-generator v0.17.6
)

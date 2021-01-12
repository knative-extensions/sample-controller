module knative.dev/sample-controller

go 1.14

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/google/licenseclassifier v0.0.0-20200708223521-3d09a0ea2f39
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/prometheus/procfs v0.0.11 // indirect
	github.com/sergi/go-diff v1.1.0 // indirect
	go.uber.org/zap v1.16.0
	k8s.io/api v0.18.12
	k8s.io/apimachinery v0.18.12
	k8s.io/client-go v11.0.1-0.20190805182717-6502b5e7b1b5+incompatible
	k8s.io/code-generator v0.18.12
	k8s.io/kube-openapi v0.0.0-20200410145947-bcb3869e6f29
	knative.dev/hack v0.0.0-20201214230143-4ed1ecb8db24
	knative.dev/pkg v0.0.0-20210107022335-51c72e24c179
)

replace (
	github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.2

	k8s.io/api => k8s.io/api v0.18.8
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.18.8
	k8s.io/apimachinery => k8s.io/apimachinery v0.18.8
	k8s.io/apiserver => k8s.io/apiserver v0.18.8
	k8s.io/client-go => k8s.io/client-go v0.18.8
	k8s.io/code-generator => k8s.io/code-generator v0.18.8
	k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20200410145947-61e04a5be9a6
)

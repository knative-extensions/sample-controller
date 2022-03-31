module knative.dev/sample-controller

go 1.15

require (
	go.uber.org/zap v1.19.1
	k8s.io/api v0.23.5
	k8s.io/apimachinery v0.23.5
	k8s.io/client-go v0.23.5
	k8s.io/code-generator v0.23.5
	k8s.io/kube-openapi v0.0.0-20211115234752-e816edb12b65
	knative.dev/hack v0.0.0-20220331040044-9c0ea69d9b4d
	knative.dev/hack/schema v0.0.0-20220331040044-9c0ea69d9b4d
	knative.dev/pkg v0.0.0-20220329144915-0a1ec2e0d46c
)

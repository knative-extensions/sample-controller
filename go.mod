module knative.dev/sample-controller

go 1.15

require (
	go.uber.org/zap v1.19.1
	k8s.io/api v0.23.5
	k8s.io/apimachinery v0.23.5
	k8s.io/client-go v0.23.5
	k8s.io/code-generator v0.23.5
	k8s.io/kube-openapi v0.0.0-20211115234752-e816edb12b65
	knative.dev/hack v0.0.0-20220524153203-12d3e2a7addc
	knative.dev/hack/schema v0.0.0-20220524153203-12d3e2a7addc
	knative.dev/pkg v0.0.0-20220524202603-19adf798efb8
)

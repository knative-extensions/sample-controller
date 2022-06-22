module knative.dev/sample-controller

go 1.15

require (
	go.uber.org/zap v1.19.1
	k8s.io/api v0.23.8
	k8s.io/apimachinery v0.23.8
	k8s.io/client-go v0.23.8
	k8s.io/code-generator v0.23.8
	k8s.io/kube-openapi v0.0.0-20211115234752-e816edb12b65
	knative.dev/hack v0.0.0-20220610014127-dc6c287516dc
	knative.dev/hack/schema v0.0.0-20220610014127-dc6c287516dc
	knative.dev/pkg v0.0.0-20220621173822-9c5a7317fa9d
)

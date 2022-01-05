module knative.dev/sample-controller

go 1.15

require (
	go.uber.org/zap v1.19.1
	k8s.io/api v0.22.5
	k8s.io/apimachinery v0.22.5
	k8s.io/client-go v0.22.5
	k8s.io/code-generator v0.22.5
	k8s.io/kube-openapi v0.0.0-20211109043538-20434351676c
	knative.dev/hack v0.0.0-20211222071919-abd085fc43de
	knative.dev/hack/schema v0.0.0-20211222071919-abd085fc43de
	knative.dev/pkg v0.0.0-20220104185830-52e42b760b54
)

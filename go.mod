module knative.dev/sample-controller

go 1.15

require (
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.16.0
	k8s.io/api v0.19.7
	k8s.io/apimachinery v0.19.7
	k8s.io/client-go v0.19.7
	k8s.io/code-generator v0.19.7
	k8s.io/kube-openapi v0.0.0-20200805222855-6aeccd4b50c6
	knative.dev/hack v0.0.0-20210427190353-86f9adc0c8e2
	knative.dev/hack/schema v0.0.0-20210427190353-86f9adc0c8e2
	knative.dev/pkg v0.0.0-20210428023153-5a308fa62139
)

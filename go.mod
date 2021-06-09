module knative.dev/sample-controller

go 1.15

require (
	github.com/aws/aws-sdk-go v1.31.12 // indirect
	go.uber.org/zap v1.17.0
	k8s.io/api v0.19.7
	k8s.io/apimachinery v0.19.7
	k8s.io/client-go v0.19.7
	k8s.io/code-generator v0.19.7
	k8s.io/kube-openapi v0.0.0-20200805222855-6aeccd4b50c6
	knative.dev/hack v0.0.0-20210609124042-e35bcb8f21ec
	knative.dev/hack/schema v0.0.0-20210609124042-e35bcb8f21ec
	knative.dev/pkg v0.0.0-20210608193741-f19eef192438
)

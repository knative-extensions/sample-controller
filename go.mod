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
	knative.dev/hack v0.0.0-20210601210329-de04b70e00d0
	knative.dev/hack/schema v0.0.0-20210601210329-de04b70e00d0
	knative.dev/pkg v0.0.0-20210602095030-0e61d6763dd6
)

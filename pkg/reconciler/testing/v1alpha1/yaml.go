package v1alpha1

import (
	"io"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func ParseYAML(yamls io.Reader) ([]unstructured.Unstructured, error) {
	decoder := yaml.NewYAMLToJSONDecoder(yamls)
	var objs []unstructured.Unstructured
	var err error
	for {
		out := unstructured.Unstructured{}
		err = decoder.Decode(&out)
		if err != nil {
			break
		}
		objs = append(objs, out)
	}
	if err != io.EOF {
		return nil, err
	}
	return objs, nil
}

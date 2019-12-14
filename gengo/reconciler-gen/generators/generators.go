/*
Copyright 2015 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package generators

import (
	"fmt"
	"path/filepath"
	"strings"

	"k8s.io/gengo/examples/set-gen/sets"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/namer"
	"k8s.io/gengo/types"
	"knative.dev/sample-controller/gengo/reconciler-gen/args"

	"k8s.io/klog"
)

// This is the comment tag that carries parameters for deep-copy generation.
const (
	tagEnabledName  = "genreconciler"
	kindTagName     = tagEnabledName + ":kind"
	stubsTagName    = tagEnabledName + ":stubs"
	clientTagName   = tagEnabledName + ":client"
	informerTagName = tagEnabledName + ":informer"
)

// enabledTagValue holds parameters from a tagName tag.
type tagValue struct {
	stubs    bool
	kind     string
	client   string
	informer string
}

func extractTag(comments []string) *tagValue {
	tags := types.ExtractCommentTags("+", comments)
	if tags[tagEnabledName] == nil {
		return nil
	}

	// If there are multiple values, abort.
	if len(tags[tagEnabledName]) > 1 {
		klog.Fatalf("Found %d %s tags: %q", len(tags[tagEnabledName]), tagEnabledName, tags[tagEnabledName])
	}

	// If we got here we are returning something.
	tag := &tagValue{}

	if v := tags[kindTagName]; v != nil {
		tag.kind = v[0]
	}

	if v := tags[stubsTagName]; v != nil {
		tag.stubs = true
	}

	if v := tags[clientTagName]; v != nil {
		tag.client = v[0]
	}

	if v := tags[informerTagName]; v != nil {
		tag.informer = v[0]
	}

	return tag
}

// NameSystems returns the name system used by the generators in this package.
func NameSystems() namer.NameSystems {
	return namer.NameSystems{
		"raw": namer.NewRawNamer("", nil),
	}
}

// DefaultNameSystem returns the default name system for ordering the types to be
// processed by the generators in this package.
func DefaultNameSystem() string {
	return "public"
}

func Packages(context *generator.Context, arguments *args.GeneratorArgs) generator.Packages {
	boilerplate, err := arguments.LoadGoBoilerplate()
	if err != nil {
		klog.Fatalf("Failed loading boilerplate: %v", err)
	}

	inputs := sets.NewString(context.Inputs...)
	packages := generator.Packages{}
	header := append([]byte(fmt.Sprintf("// +build !%s\n\n", arguments.GeneratedBuildTag)), boilerplate...)

	editHeader, err := arguments.LoadEditGoBoilerplate()
	if err != nil {
		klog.Fatalf("Failed loading edit boilerplate: %v", err)
	}

	for i := range inputs {
		klog.V(5).Infof("Considering pkg %q", i)
		pkg := context.Universe[i]
		if pkg == nil {
			// If the input had no Go files, for example.
			continue
		}

		ptag := extractTag(pkg.Comments)
		if ptag != nil {
			klog.V(5).Infof("  tag: %+v", ptag)
		} else {
			klog.V(5).Infof("  no tag")
			continue
		}

		packages = append(packages,
			&generator.DefaultPackage{
				PackageName: strings.Split(filepath.Base(pkg.Path), ".")[0],
				PackagePath: pkg.Path,
				HeaderText:  header,
				GeneratorFunc: func(c *generator.Context) (generators []generator.Generator) {
					return []generator.Generator{
						NewGenController(arguments.OutputFileBaseName+".controller", pkg.Path, ptag.kind, ptag.client, ptag.informer),
						NewGenReconciler(arguments.OutputFileBaseName+".reconciler", pkg.Path, ptag.kind, ptag.client, ptag.informer),
					}
				},
				FilterFunc: func(c *generator.Context, t *types.Type) bool {
					return false
				},
			})

		if ptag.stubs {

			name := ptag.kind[strings.LastIndex(ptag.kind, ".")+1:]

			packages = append(packages,
				&generator.DefaultPackage{
					PackageName: strings.Split(filepath.Base(pkg.Path), ".")[0],
					PackagePath: pkg.Path,
					HeaderText:  editHeader,
					GeneratorFunc: func(c *generator.Context) (generators []generator.Generator) {
						return []generator.Generator{
							NewGenStubController("controller", pkg.Path, ptag.kind, ptag.client, ptag.informer),
							NewGenStubReconciler(strings.ToLower(name), pkg.Path, ptag.kind, ptag.client, ptag.informer),
						}
					},
					FilterFunc: func(c *generator.Context, t *types.Type) bool {
						return false
					},
				},
			)
		}

	}
	return packages
}

func UnsafeGuessKindToResource(kind string) string {
	if len(kind) == 0 {
		return ""
	}

	switch string(kind[len(kind)-1]) {
	case "s":
		return kind + "es"
	case "y":
		return strings.TrimSuffix(kind, "y") + "ies"
	}

	return kind + "s"
}

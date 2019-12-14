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
	"io"
	"strings"

	"k8s.io/gengo/generator"
	"k8s.io/gengo/namer"
	"k8s.io/gengo/types"
	"k8s.io/klog"
)

// genStubReconciler
type genStubReconciler struct {
	generator.DefaultGen
	targetPackage string

	kind     string
	client   string
	informer string

	imports      namer.ImportTracker
	typesForInit []*types.Type
}

func NewGenStubReconciler(sanitizedName, targetPackage string, kind, client, informer string) generator.Generator {
	return &genStubReconciler{
		DefaultGen: generator.DefaultGen{
			OptionalName: sanitizedName,
		},
		targetPackage: targetPackage,
		kind:          kind,
		client:        client,
		informer:      informer,
		imports:       generator.NewImportTracker(),
		typesForInit:  make([]*types.Type, 0),
	}
}

func (g *genStubReconciler) Namers(c *generator.Context) namer.NameSystems {
	return NameSystems()
}

func (g *genStubReconciler) Filter(c *generator.Context, t *types.Type) bool {
	return false
}

func (g *genStubReconciler) isOtherPackage(pkg string) bool {
	if pkg == g.targetPackage {
		return false
	}
	if strings.HasSuffix(pkg, "\""+g.targetPackage+"\"") {
		return false
	}
	return true
}

func (g *genStubReconciler) Imports(c *generator.Context) (imports []string) {
	importLines := []string{}
	for _, singleImport := range g.imports.ImportLines() {
		if g.isOtherPackage(singleImport) {
			importLines = append(importLines, singleImport)
		}
	}
	return importLines
}

func (g *genStubReconciler) Init(c *generator.Context, w io.Writer) error {
	kind := g.kind
	klog.Infof("Generating genreconciler function for kind %v", kind)

	sw := generator.NewSnippetWriter(w, c, "{{", "}}")

	klog.V(5).Infof("processing kind %v", g.kind)

	pkg := kind[:strings.LastIndex(kind, ".")]
	name := kind[strings.LastIndex(kind, ".")+1:]

	m := map[string]interface{}{
		"resourceName":   c.Universe.Type(types.Name{Name: strings.ToLower(name), Package: g.targetPackage}),
		"resourceNames":  c.Universe.Type(types.Name{Name: UnsafeGuessKindToResource(name), Package: g.targetPackage}),
		"resource":       c.Universe.Type(types.Name{Package: pkg, Name: name}),
		"controllerImpl": c.Universe.Type(types.Name{Package: "knative.dev/pkg/controller", Name: "Impl"}),
		"loggingFromContext": c.Universe.Function(types.Name{
			Package: "knative.dev/pkg/logging",
			Name:    "FromContext",
		}),
		"clientGet": c.Universe.Function(types.Name{
			Package: g.client,
			Name:    "Get",
		}),
		"informerGet": c.Universe.Function(types.Name{
			Package: g.informer,
			Name:    "Get",
		}),
		"corev1EventSource": c.Universe.Function(types.Name{
			Package: "k8s.io/api/core/v1",
			Name:    "EventSource",
		}),
	}

	sw.Do(stubReconcilerFactory, m)

	return sw.Error()
}

func (g *genStubReconciler) GenerateType(c *generator.Context, t *types.Type, w io.Writer) error {
	return nil
}

var stubReconcilerFactory = `
// Reconciler implements controller.Reconciler for {{.resource|raw}} resources.
type Reconciler struct {
	Core
}

// Check that our Reconciler implements reconciler.Interface
var _ Interface = (*Reconciler)(nil)

// ReconcileKind implements Interface.ReconcileKind.
func (r *Reconciler) ReconcileKind(ctx context.Context, o *{{.resource|raw}}) error {
	if o.GetDeletionTimestamp() != nil {
		// Check for a DeletionTimestamp.  If present, elide the normal reconcile logic.
		// When a controller needs finalizer handling, it would go here.
		return nil
	}	
	o.Status.InitializeConditions()

	// TODO: add custom reconciliation logic here.

	o.Status.ObservedGeneration = o.Generation
	return nil
}

`
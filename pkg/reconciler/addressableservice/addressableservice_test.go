package addressableservice

import (
	"context"
	"flag"
	"fmt"
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/cucumber/messages-go/v10"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	clientgotesting "k8s.io/client-go/testing"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/logging"
	"knative.dev/pkg/tracker"
	client "knative.dev/sample-controller/pkg/client/injection/client"
	"knative.dev/sample-controller/pkg/client/injection/reconciler/samples/v1alpha1/addressableservice"
	. "knative.dev/sample-controller/pkg/reconciler/testing/v1alpha1"
	"os"
	"strings"
	"testing"

	pkgtest "knative.dev/pkg/reconciler/testing"
)

var opt = godog.Options{
	Output: colors.Colored(os.Stdout),
}

var testStatus int

func TestMain(m *testing.M) {
	flag.Parse()

	if len(flag.Args()) > 0 {
		opt.Paths = flag.Args()
	} else {
		opt.Paths = []string{
			"./features/",
		}
	}

	format := "progress"
	for _, arg := range os.Args[1:] {
		if arg == "-test.v=true" { // go test transforms -v option
			format = "pretty"
			break
		}
	}

	opt.Format = format

	os.Exit(m.Run())
}

func TestReconcile(t *testing.T) {
	status := godog.RunWithOptions("AddressableService", func(s *godog.Suite) {
		AddressableServiceFeatureContext(t, s)
	}, opt)

	if status != 0 {
		t.Fail()
	}
}

func AddressableServiceFeatureContext(t *testing.T, s *godog.Suite) {
	ctx := context.Background()

	rt := &ReconcilerTest{
		t: t,
		row: pkgtest.TableRow{
			Name: "AddressableService",
			Ctx:  ctx,
			//OtherTestData:           nil,
			//Objects:                 nil,
			//Key:                     "",
			//WantErr:                 false,
			//WantCreates:             nil,
			//WantUpdates:             nil,
			//WantStatusUpdates:       nil,
			//WantDeletes:             nil,
			//WantDeleteCollections:   nil,
			//WantPatches:             nil,
			//WantEvents:              nil,
			//WithReactors:            nil,
			//SkipNamespaceValidation: false,
			//PostConditions:          nil,
			//Reconciler:              nil,
		}}

	s.Step(`^the following objects:$`, rt.theFollowingObjects)
	s.Step(`^an AddressableService reconciler$`, rt.anAddressableServiceReconciler)
	s.Step(`^reconciling "([^"]*)"$`, rt.reconcilingKey)
	s.Step(`^expect nothing$`, rt.expectNothing)
	s.Step(`^expect status updates:$`, rt.expectStatusUpdates)
	s.Step(`^expect Kubernetes Events:$`, rt.expectKubernetesEvents)

	s.AfterScenario(func(pickle *messages.Pickle, err error) {
		originObjects := make([]runtime.Object, 0, len(rt.row.Objects))
		for _, obj := range rt.row.Objects {
			originObjects = append(originObjects, obj.DeepCopyObject())
		}

		// Reconcile.
		rt.row.Name = pickle.Name
		rt.t.Run(pickle.Name, func(t *testing.T) {
			t.Helper()
			rt.row.Test(t, rt.factory)
		})

		// Validate cached objects do not get soiled after controller loops
		if diff := cmp.Diff(originObjects, rt.row.Objects, safeDeployDiff, cmpopts.EquateEmpty()); diff != "" {
			rt.t.Errorf("Unexpected objects in test %s (-want, +got): %v", rt.row.Name, diff)
		}
	})

	_ = rt
}

type ReconcilerTest struct {
	t   *testing.T
	row pkgtest.TableRow

	factory pkgtest.Factory
}

func (rt *ReconcilerTest) theFollowingObjects(y *messages.PickleStepArgument_PickleDocString) error {
	objs, err := ParseYAML(strings.NewReader(y.Content))
	if err != nil {
		return err
	}

	rt.row.Objects = make([]runtime.Object, 0, len(objs))

	for _, obj := range objs {
		rt.row.Objects = append(rt.row.Objects, obj.DeepCopyObject())
	}

	rt.row.Objects = ToKnownObjects(rt.row.Objects)

	return nil
}

func (rt *ReconcilerTest) anAddressableServiceReconciler() error {

	rt.factory = MakeFactory(func(ctx context.Context, listers *Listers, watcher configmap.Watcher) controller.Reconciler {
		r := &Reconciler{
			Tracker:       tracker.New(func(types.NamespacedName) {}, 0),
			ServiceLister: listers.GetK8sServiceLister(),
		}

		return addressableservice.NewReconciler(ctx, logging.FromContext(ctx),
			client.Get(ctx), listers.GetAddressableServiceLister(),
			controller.GetEventRecorder(ctx), r)
	})

	return nil
}

func (rt *ReconcilerTest) reconcilingKey(key string) error {
	// Set the reconciler key
	rt.row.Key = key

	return nil
}

func (rt *ReconcilerTest) expectNothing() error {
	return nil
}

func (rt *ReconcilerTest) expectStatusUpdates(y *messages.PickleStepArgument_PickleDocString) error {
	objs, err := ParseYAML(strings.NewReader(y.Content))
	if err != nil {
		return err
	}

	updates := make([]runtime.Object, 0, len(objs))
	for _, obj := range objs {
		updates = append(updates, obj.DeepCopyObject())
	}
	updates = ToKnownObjects(updates)

	rt.row.WantStatusUpdates = make([]clientgotesting.UpdateActionImpl, 0)
	for _, u := range updates {
		updateAction := clientgotesting.UpdateActionImpl{
			Object: u,
		}
		rt.row.WantStatusUpdates = append(rt.row.WantStatusUpdates, updateAction)
	}
	return nil
}

func (rt *ReconcilerTest) expectKubernetesEvents(attributes *messages.PickleStepArgument_PickleTable) error {

	rt.row.WantEvents = make([]string, 0)

	for _, row := range attributes.Rows {
		eventType := row.Cells[0].Value
		reason := row.Cells[1].Value
		message := row.Cells[2].Value

		if eventType == "Type" {
			// ignore the headers
			continue
		}

		rt.row.WantEvents = append(rt.row.WantEvents, fmt.Sprintf(eventType+" "+reason+" "+message))
	}
	return nil
}

var (
	safeDeployDiff = cmpopts.IgnoreUnexported(resource.Quantity{})
)

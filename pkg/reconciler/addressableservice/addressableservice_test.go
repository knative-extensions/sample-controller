package addressableservice

import (
	"context"
	"flag"
	"github.com/cucumber/messages-go/v10"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
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

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"

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
	s.Step(`^an unmodified cache$`, rt.anUnmodifiedCache)
	_ = rt
}

type ReconcilerTest struct {
	originObjects []runtime.Object
	row           pkgtest.TableRow

	t *testing.T

	r       controller.Reconciler
	listers *Listers
	cmw     configmap.Watcher
}

func (rt *ReconcilerTest) theFollowingObjects(y *messages.PickleStepArgument_PickleDocString) error {
	objs, err := ParseYAML(strings.NewReader(y.Content))
	if err != nil {
		return err
	}

	originObjects := make([]runtime.Object, 0, len(objs))
	rt.row.Objects = make([]runtime.Object, 0, len(objs))

	for _, obj := range objs {
		originObjects = append(originObjects, obj.DeepCopyObject())
		rt.row.Objects = append(rt.row.Objects, obj.DeepCopyObject())
	}

	return nil
}

func (rt *ReconcilerTest) anAddressableServiceReconciler() error {

	factory := MakeFactory(func(ctx context.Context, listers *Listers, watcher configmap.Watcher) controller.Reconciler {
		r := &Reconciler{
			Tracker:       tracker.New(func(types.NamespacedName) {}, 0),
			ServiceLister: rt.listers.GetK8sServiceLister(),
		}

		return addressableservice.NewReconciler(ctx, logging.FromContext(ctx),
			client.Get(ctx), rt.listers.GetAddressableServiceLister(),
			controller.GetEventRecorder(ctx), r)
	})

	rt.t.Run(rt.row.Name, func(t *testing.T) {
		t.Helper()
		rt.row.Test(t, factory)
	})

	return nil
}

func (rt *ReconcilerTest) reconcilingKey(key string) error {
	return godog.ErrPending
}

func (rt *ReconcilerTest) expectNothing() error {
	return godog.ErrPending
}

func (rt *ReconcilerTest) anUnmodifiedCache() error {
	// Validate cached objects do not get soiled after controller loops
	if diff := cmp.Diff(rt.originObjects, rt.row.Objects, safeDeployDiff, cmpopts.EquateEmpty()); diff != "" {
		rt.t.Errorf("Unexpected objects in test %s (-want, +got): %v", rt.row.Name, diff)
	}
	return nil
}

var (
	//ignoreLastTransitionTime = cmp.FilterPath(func(p cmp.Path) bool {
	//	return strings.HasSuffix(p.String(), "LastTransitionTime.Inner.Time")
	//}, cmp.Ignore())

	safeDeployDiff = cmpopts.IgnoreUnexported(resource.Quantity{})
)

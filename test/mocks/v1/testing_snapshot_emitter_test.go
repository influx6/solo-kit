// Code generated by solo-kit. DO NOT EDIT.

// +build solokit

package v1

import (
	"context"
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/go-utils/kubeutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	kuberc "github.com/solo-io/solo-kit/pkg/api/v1/clients/kube"
	"github.com/solo-io/solo-kit/pkg/utils/log"
	"github.com/solo-io/solo-kit/test/helpers"
	"github.com/solo-io/solo-kit/test/setup"
	"k8s.io/client-go/rest"

	// Needed to run tests in GKE
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	// From https://github.com/kubernetes/client-go/blob/53c7adfd0294caa142d961e1f780f74081d5b15f/examples/out-of-cluster-client-configuration/main.go#L31
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/cache"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

var _ = Describe("V1Emitter", func() {
	if os.Getenv("RUN_KUBE_TESTS") != "1" {
		log.Printf("This test creates kubernetes resources and is disabled by default. To enable, set RUN_KUBE_TESTS=1 in your env.")
		return
	}
	var (
		namespace1                string
		namespace2                string
		name1, name2              = "angela" + helpers.RandString(3), "bob" + helpers.RandString(3)
		cfg                       *rest.Config
		emitter                   TestingEmitter
		mockResourceClient        MockResourceClient
		fakeResourceClient        FakeResourceClient
		anotherMockResourceClient AnotherMockResourceClient
		clusterResourceClient     ClusterResourceClient
	)

	BeforeEach(func() {
		namespace1 = helpers.RandString(8)
		namespace2 = helpers.RandString(8)
		var err error
		cfg, err = kubeutils.GetConfig("", "")
		Expect(err).NotTo(HaveOccurred())
		err = setup.SetupKubeForTest(namespace1)
		Expect(err).NotTo(HaveOccurred())
		err = setup.SetupKubeForTest(namespace2)
		Expect(err).NotTo(HaveOccurred())
		var kube kubernetes.Interface
		// MockResource Constructor
		mockResourceClientFactory := &factory.KubeResourceClientFactory{
			Crd:         MockResourceCrd,
			Cfg:         cfg,
			SharedCache: kuberc.NewKubeCache(context.TODO()),
		}
		mockResourceClient, err = NewMockResourceClient(mockResourceClientFactory)
		Expect(err).NotTo(HaveOccurred())
		// FakeResource Constructor
		kube, err = kubernetes.NewForConfig(cfg)
		Expect(err).NotTo(HaveOccurred())

		kcache, err := cache.NewKubeCoreCache(context.TODO(), kube)
		Expect(err).NotTo(HaveOccurred())
		fakeResourceClientFactory := &factory.KubeConfigMapClientFactory{
			Clientset: kube,
			Cache:     kcache,
		}
		fakeResourceClient, err = NewFakeResourceClient(fakeResourceClientFactory)
		Expect(err).NotTo(HaveOccurred())
		// AnotherMockResource Constructor
		anotherMockResourceClientFactory := &factory.KubeResourceClientFactory{
			Crd:         AnotherMockResourceCrd,
			Cfg:         cfg,
			SharedCache: kuberc.NewKubeCache(context.TODO()),
		}
		anotherMockResourceClient, err = NewAnotherMockResourceClient(anotherMockResourceClientFactory)
		Expect(err).NotTo(HaveOccurred())
		// ClusterResource Constructor
		clusterResourceClientFactory := &factory.KubeResourceClientFactory{
			Crd:         ClusterResourceCrd,
			Cfg:         cfg,
			SharedCache: kuberc.NewKubeCache(context.TODO()),
		}
		clusterResourceClient, err = NewClusterResourceClient(clusterResourceClientFactory)
		Expect(err).NotTo(HaveOccurred())
		emitter = NewTestingEmitter(mockResourceClient, fakeResourceClient, anotherMockResourceClient, clusterResourceClient)
	})
	AfterEach(func() {
		setup.TeardownKube(namespace1)
		setup.TeardownKube(namespace2)
		clusterResourceClient.Delete(name1, clients.DeleteOpts{})
		clusterResourceClient.Delete(name2, clients.DeleteOpts{})
	})
	It("tracks snapshots on changes to any resource", func() {
		ctx := context.Background()
		err := emitter.Register()
		Expect(err).NotTo(HaveOccurred())

		snapshots, errs, err := emitter.Snapshots([]string{namespace1, namespace2}, clients.WatchOpts{
			Ctx:         ctx,
			RefreshRate: time.Second,
		})
		Expect(err).NotTo(HaveOccurred())

		var snap *TestingSnapshot

		/*
			MockResource
		*/

		assertSnapshotMocks := func(expectMocks MockResourceList, unexpectMocks MockResourceList) {
		drain:
			for {
				select {
				case snap = <-snapshots:
					for _, expected := range expectMocks {
						if _, err := snap.Mocks.List().Find(expected.Metadata.Ref().Strings()); err != nil {
							continue drain
						}
					}
					for _, unexpected := range unexpectMocks {
						if _, err := snap.Mocks.List().Find(unexpected.Metadata.Ref().Strings()); err == nil {
							continue drain
						}
					}
					break drain
				case err := <-errs:
					Expect(err).NotTo(HaveOccurred())
				case <-time.After(time.Second * 10):
					nsList1, _ := mockResourceClient.List(namespace1, clients.ListOpts{})
					nsList2, _ := mockResourceClient.List(namespace2, clients.ListOpts{})
					combined := MocksByNamespace{
						namespace1: nsList1,
						namespace2: nsList2,
					}
					Fail("expected final snapshot before 10 seconds. expected " + log.Sprintf("%v", combined))
				}
			}
		}
		mockResource1a, err := mockResourceClient.Write(NewMockResource(namespace1, name1), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		mockResource1b, err := mockResourceClient.Write(NewMockResource(namespace2, name1), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotMocks(MockResourceList{mockResource1a, mockResource1b}, nil)
		mockResource2a, err := mockResourceClient.Write(NewMockResource(namespace1, name2), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		mockResource2b, err := mockResourceClient.Write(NewMockResource(namespace2, name2), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotMocks(MockResourceList{mockResource1a, mockResource1b, mockResource2a, mockResource2b}, nil)

		err = mockResourceClient.Delete(mockResource2a.Metadata.Namespace, mockResource2a.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		err = mockResourceClient.Delete(mockResource2b.Metadata.Namespace, mockResource2b.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotMocks(MockResourceList{mockResource1a, mockResource1b}, MockResourceList{mockResource2a, mockResource2b})

		err = mockResourceClient.Delete(mockResource1a.Metadata.Namespace, mockResource1a.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		err = mockResourceClient.Delete(mockResource1b.Metadata.Namespace, mockResource1b.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotMocks(nil, MockResourceList{mockResource1a, mockResource1b, mockResource2a, mockResource2b})

		/*
			FakeResource
		*/

		assertSnapshotFakes := func(expectFakes FakeResourceList, unexpectFakes FakeResourceList) {
		drain:
			for {
				select {
				case snap = <-snapshots:
					for _, expected := range expectFakes {
						if _, err := snap.Fakes.List().Find(expected.Metadata.Ref().Strings()); err != nil {
							continue drain
						}
					}
					for _, unexpected := range unexpectFakes {
						if _, err := snap.Fakes.List().Find(unexpected.Metadata.Ref().Strings()); err == nil {
							continue drain
						}
					}
					break drain
				case err := <-errs:
					Expect(err).NotTo(HaveOccurred())
				case <-time.After(time.Second * 10):
					nsList1, _ := fakeResourceClient.List(namespace1, clients.ListOpts{})
					nsList2, _ := fakeResourceClient.List(namespace2, clients.ListOpts{})
					combined := FakesByNamespace{
						namespace1: nsList1,
						namespace2: nsList2,
					}
					Fail("expected final snapshot before 10 seconds. expected " + log.Sprintf("%v", combined))
				}
			}
		}
		fakeResource1a, err := fakeResourceClient.Write(NewFakeResource(namespace1, name1), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		fakeResource1b, err := fakeResourceClient.Write(NewFakeResource(namespace2, name1), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotFakes(FakeResourceList{fakeResource1a, fakeResource1b}, nil)
		fakeResource2a, err := fakeResourceClient.Write(NewFakeResource(namespace1, name2), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		fakeResource2b, err := fakeResourceClient.Write(NewFakeResource(namespace2, name2), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotFakes(FakeResourceList{fakeResource1a, fakeResource1b, fakeResource2a, fakeResource2b}, nil)

		err = fakeResourceClient.Delete(fakeResource2a.Metadata.Namespace, fakeResource2a.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		err = fakeResourceClient.Delete(fakeResource2b.Metadata.Namespace, fakeResource2b.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotFakes(FakeResourceList{fakeResource1a, fakeResource1b}, FakeResourceList{fakeResource2a, fakeResource2b})

		err = fakeResourceClient.Delete(fakeResource1a.Metadata.Namespace, fakeResource1a.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		err = fakeResourceClient.Delete(fakeResource1b.Metadata.Namespace, fakeResource1b.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotFakes(nil, FakeResourceList{fakeResource1a, fakeResource1b, fakeResource2a, fakeResource2b})

		/*
			AnotherMockResource
		*/

		assertSnapshotAnothermockresources := func(expectAnothermockresources AnotherMockResourceList, unexpectAnothermockresources AnotherMockResourceList) {
		drain:
			for {
				select {
				case snap = <-snapshots:
					for _, expected := range expectAnothermockresources {
						if _, err := snap.Anothermockresources.List().Find(expected.Metadata.Ref().Strings()); err != nil {
							continue drain
						}
					}
					for _, unexpected := range unexpectAnothermockresources {
						if _, err := snap.Anothermockresources.List().Find(unexpected.Metadata.Ref().Strings()); err == nil {
							continue drain
						}
					}
					break drain
				case err := <-errs:
					Expect(err).NotTo(HaveOccurred())
				case <-time.After(time.Second * 10):
					nsList1, _ := anotherMockResourceClient.List(namespace1, clients.ListOpts{})
					nsList2, _ := anotherMockResourceClient.List(namespace2, clients.ListOpts{})
					combined := AnothermockresourcesByNamespace{
						namespace1: nsList1,
						namespace2: nsList2,
					}
					Fail("expected final snapshot before 10 seconds. expected " + log.Sprintf("%v", combined))
				}
			}
		}
		anotherMockResource1a, err := anotherMockResourceClient.Write(NewAnotherMockResource(namespace1, name1), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		anotherMockResource1b, err := anotherMockResourceClient.Write(NewAnotherMockResource(namespace2, name1), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotAnothermockresources(AnotherMockResourceList{anotherMockResource1a, anotherMockResource1b}, nil)
		anotherMockResource2a, err := anotherMockResourceClient.Write(NewAnotherMockResource(namespace1, name2), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		anotherMockResource2b, err := anotherMockResourceClient.Write(NewAnotherMockResource(namespace2, name2), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotAnothermockresources(AnotherMockResourceList{anotherMockResource1a, anotherMockResource1b, anotherMockResource2a, anotherMockResource2b}, nil)

		err = anotherMockResourceClient.Delete(anotherMockResource2a.Metadata.Namespace, anotherMockResource2a.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		err = anotherMockResourceClient.Delete(anotherMockResource2b.Metadata.Namespace, anotherMockResource2b.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotAnothermockresources(AnotherMockResourceList{anotherMockResource1a, anotherMockResource1b}, AnotherMockResourceList{anotherMockResource2a, anotherMockResource2b})

		err = anotherMockResourceClient.Delete(anotherMockResource1a.Metadata.Namespace, anotherMockResource1a.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		err = anotherMockResourceClient.Delete(anotherMockResource1b.Metadata.Namespace, anotherMockResource1b.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotAnothermockresources(nil, AnotherMockResourceList{anotherMockResource1a, anotherMockResource1b, anotherMockResource2a, anotherMockResource2b})

		/*
			ClusterResource
		*/

		assertSnapshotClusterresources := func(expectClusterresources ClusterResourceList, unexpectClusterresources ClusterResourceList) {
		drain:
			for {
				select {
				case snap = <-snapshots:
					for _, expected := range expectClusterresources {
						if _, err := snap.Clusterresources.Find(expected.Metadata.Ref().Strings()); err != nil {
							continue drain
						}
					}
					for _, unexpected := range unexpectClusterresources {
						if _, err := snap.Clusterresources.Find(unexpected.Metadata.Ref().Strings()); err == nil {
							continue drain
						}
					}
					break drain
				case err := <-errs:
					Expect(err).NotTo(HaveOccurred())
				case <-time.After(time.Second * 10):
					nsList, _ := clusterResourceClient.List(clients.ListOpts{})
					combined := ClusterresourcesByNamespace{
						"": nsList,
					}
					Fail("expected final snapshot before 10 seconds. expected " + log.Sprintf("%v", combined))
				}
			}
		}
		clusterResource1a, err := clusterResourceClient.Write(NewClusterResource(namespace1, name1), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotClusterresources(ClusterResourceList{clusterResource1a}, nil)
		clusterResource2a, err := clusterResourceClient.Write(NewClusterResource(namespace1, name2), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotClusterresources(ClusterResourceList{clusterResource1a, clusterResource2a}, nil)

		err = clusterResourceClient.Delete(clusterResource2a.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotClusterresources(ClusterResourceList{clusterResource1a}, ClusterResourceList{clusterResource2a})

		err = clusterResourceClient.Delete(clusterResource1a.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotClusterresources(nil, ClusterResourceList{clusterResource1a, clusterResource2a})
	})
	It("tracks snapshots on changes to any resource using AllNamespace", func() {
		ctx := context.Background()
		err := emitter.Register()
		Expect(err).NotTo(HaveOccurred())

		snapshots, errs, err := emitter.Snapshots([]string{""}, clients.WatchOpts{
			Ctx:         ctx,
			RefreshRate: time.Second,
		})
		Expect(err).NotTo(HaveOccurred())

		var snap *TestingSnapshot

		/*
			MockResource
		*/

		assertSnapshotMocks := func(expectMocks MockResourceList, unexpectMocks MockResourceList) {
		drain:
			for {
				select {
				case snap = <-snapshots:
					for _, expected := range expectMocks {
						if _, err := snap.Mocks.List().Find(expected.Metadata.Ref().Strings()); err != nil {
							continue drain
						}
					}
					for _, unexpected := range unexpectMocks {
						if _, err := snap.Mocks.List().Find(unexpected.Metadata.Ref().Strings()); err == nil {
							continue drain
						}
					}
					break drain
				case err := <-errs:
					Expect(err).NotTo(HaveOccurred())
				case <-time.After(time.Second * 10):
					nsList1, _ := mockResourceClient.List(namespace1, clients.ListOpts{})
					nsList2, _ := mockResourceClient.List(namespace2, clients.ListOpts{})
					combined := MocksByNamespace{
						namespace1: nsList1,
						namespace2: nsList2,
					}
					Fail("expected final snapshot before 10 seconds. expected " + log.Sprintf("%v", combined))
				}
			}
		}
		mockResource1a, err := mockResourceClient.Write(NewMockResource(namespace1, name1), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		mockResource1b, err := mockResourceClient.Write(NewMockResource(namespace2, name1), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotMocks(MockResourceList{mockResource1a, mockResource1b}, nil)
		mockResource2a, err := mockResourceClient.Write(NewMockResource(namespace1, name2), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		mockResource2b, err := mockResourceClient.Write(NewMockResource(namespace2, name2), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotMocks(MockResourceList{mockResource1a, mockResource1b, mockResource2a, mockResource2b}, nil)

		err = mockResourceClient.Delete(mockResource2a.Metadata.Namespace, mockResource2a.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		err = mockResourceClient.Delete(mockResource2b.Metadata.Namespace, mockResource2b.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotMocks(MockResourceList{mockResource1a, mockResource1b}, MockResourceList{mockResource2a, mockResource2b})

		err = mockResourceClient.Delete(mockResource1a.Metadata.Namespace, mockResource1a.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		err = mockResourceClient.Delete(mockResource1b.Metadata.Namespace, mockResource1b.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotMocks(nil, MockResourceList{mockResource1a, mockResource1b, mockResource2a, mockResource2b})

		/*
			FakeResource
		*/

		assertSnapshotFakes := func(expectFakes FakeResourceList, unexpectFakes FakeResourceList) {
		drain:
			for {
				select {
				case snap = <-snapshots:
					for _, expected := range expectFakes {
						if _, err := snap.Fakes.List().Find(expected.Metadata.Ref().Strings()); err != nil {
							continue drain
						}
					}
					for _, unexpected := range unexpectFakes {
						if _, err := snap.Fakes.List().Find(unexpected.Metadata.Ref().Strings()); err == nil {
							continue drain
						}
					}
					break drain
				case err := <-errs:
					Expect(err).NotTo(HaveOccurred())
				case <-time.After(time.Second * 10):
					nsList1, _ := fakeResourceClient.List(namespace1, clients.ListOpts{})
					nsList2, _ := fakeResourceClient.List(namespace2, clients.ListOpts{})
					combined := FakesByNamespace{
						namespace1: nsList1,
						namespace2: nsList2,
					}
					Fail("expected final snapshot before 10 seconds. expected " + log.Sprintf("%v", combined))
				}
			}
		}
		fakeResource1a, err := fakeResourceClient.Write(NewFakeResource(namespace1, name1), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		fakeResource1b, err := fakeResourceClient.Write(NewFakeResource(namespace2, name1), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotFakes(FakeResourceList{fakeResource1a, fakeResource1b}, nil)
		fakeResource2a, err := fakeResourceClient.Write(NewFakeResource(namespace1, name2), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		fakeResource2b, err := fakeResourceClient.Write(NewFakeResource(namespace2, name2), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotFakes(FakeResourceList{fakeResource1a, fakeResource1b, fakeResource2a, fakeResource2b}, nil)

		err = fakeResourceClient.Delete(fakeResource2a.Metadata.Namespace, fakeResource2a.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		err = fakeResourceClient.Delete(fakeResource2b.Metadata.Namespace, fakeResource2b.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotFakes(FakeResourceList{fakeResource1a, fakeResource1b}, FakeResourceList{fakeResource2a, fakeResource2b})

		err = fakeResourceClient.Delete(fakeResource1a.Metadata.Namespace, fakeResource1a.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		err = fakeResourceClient.Delete(fakeResource1b.Metadata.Namespace, fakeResource1b.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotFakes(nil, FakeResourceList{fakeResource1a, fakeResource1b, fakeResource2a, fakeResource2b})

		/*
			AnotherMockResource
		*/

		assertSnapshotAnothermockresources := func(expectAnothermockresources AnotherMockResourceList, unexpectAnothermockresources AnotherMockResourceList) {
		drain:
			for {
				select {
				case snap = <-snapshots:
					for _, expected := range expectAnothermockresources {
						if _, err := snap.Anothermockresources.List().Find(expected.Metadata.Ref().Strings()); err != nil {
							continue drain
						}
					}
					for _, unexpected := range unexpectAnothermockresources {
						if _, err := snap.Anothermockresources.List().Find(unexpected.Metadata.Ref().Strings()); err == nil {
							continue drain
						}
					}
					break drain
				case err := <-errs:
					Expect(err).NotTo(HaveOccurred())
				case <-time.After(time.Second * 10):
					nsList1, _ := anotherMockResourceClient.List(namespace1, clients.ListOpts{})
					nsList2, _ := anotherMockResourceClient.List(namespace2, clients.ListOpts{})
					combined := AnothermockresourcesByNamespace{
						namespace1: nsList1,
						namespace2: nsList2,
					}
					Fail("expected final snapshot before 10 seconds. expected " + log.Sprintf("%v", combined))
				}
			}
		}
		anotherMockResource1a, err := anotherMockResourceClient.Write(NewAnotherMockResource(namespace1, name1), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		anotherMockResource1b, err := anotherMockResourceClient.Write(NewAnotherMockResource(namespace2, name1), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotAnothermockresources(AnotherMockResourceList{anotherMockResource1a, anotherMockResource1b}, nil)
		anotherMockResource2a, err := anotherMockResourceClient.Write(NewAnotherMockResource(namespace1, name2), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		anotherMockResource2b, err := anotherMockResourceClient.Write(NewAnotherMockResource(namespace2, name2), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotAnothermockresources(AnotherMockResourceList{anotherMockResource1a, anotherMockResource1b, anotherMockResource2a, anotherMockResource2b}, nil)

		err = anotherMockResourceClient.Delete(anotherMockResource2a.Metadata.Namespace, anotherMockResource2a.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		err = anotherMockResourceClient.Delete(anotherMockResource2b.Metadata.Namespace, anotherMockResource2b.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotAnothermockresources(AnotherMockResourceList{anotherMockResource1a, anotherMockResource1b}, AnotherMockResourceList{anotherMockResource2a, anotherMockResource2b})

		err = anotherMockResourceClient.Delete(anotherMockResource1a.Metadata.Namespace, anotherMockResource1a.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		err = anotherMockResourceClient.Delete(anotherMockResource1b.Metadata.Namespace, anotherMockResource1b.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotAnothermockresources(nil, AnotherMockResourceList{anotherMockResource1a, anotherMockResource1b, anotherMockResource2a, anotherMockResource2b})

		/*
			ClusterResource
		*/

		assertSnapshotClusterresources := func(expectClusterresources ClusterResourceList, unexpectClusterresources ClusterResourceList) {
		drain:
			for {
				select {
				case snap = <-snapshots:
					for _, expected := range expectClusterresources {
						if _, err := snap.Clusterresources.Find(expected.Metadata.Ref().Strings()); err != nil {
							continue drain
						}
					}
					for _, unexpected := range unexpectClusterresources {
						if _, err := snap.Clusterresources.Find(unexpected.Metadata.Ref().Strings()); err == nil {
							continue drain
						}
					}
					break drain
				case err := <-errs:
					Expect(err).NotTo(HaveOccurred())
				case <-time.After(time.Second * 10):
					nsList, _ := clusterResourceClient.List(clients.ListOpts{})
					combined := ClusterresourcesByNamespace{
						"": nsList,
					}
					Fail("expected final snapshot before 10 seconds. expected " + log.Sprintf("%v", combined))
				}
			}
		}
		clusterResource1a, err := clusterResourceClient.Write(NewClusterResource(namespace1, name1), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotClusterresources(ClusterResourceList{clusterResource1a}, nil)
		clusterResource2a, err := clusterResourceClient.Write(NewClusterResource(namespace1, name2), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotClusterresources(ClusterResourceList{clusterResource1a, clusterResource2a}, nil)

		err = clusterResourceClient.Delete(clusterResource2a.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotClusterresources(ClusterResourceList{clusterResource1a}, ClusterResourceList{clusterResource2a})

		err = clusterResourceClient.Delete(clusterResource1a.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotClusterresources(nil, ClusterResourceList{clusterResource1a, clusterResource2a})
	})
})

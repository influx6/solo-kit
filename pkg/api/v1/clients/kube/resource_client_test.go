package kube_test

import (
	"context"
	"os"
	"time"

	"github.com/solo-io/solo-kit/pkg/api/v1/resources"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/go-utils/kubeutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/crd/client/clientset/versioned"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/crd/client/clientset/versioned/fake"
	solov1 "github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/crd/solo.io/v1"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/log"
	"github.com/solo-io/solo-kit/test/helpers"
	"github.com/solo-io/solo-kit/test/mocks/util"
	v1 "github.com/solo-io/solo-kit/test/mocks/v1"
	"github.com/solo-io/solo-kit/test/setup"
	"github.com/solo-io/solo-kit/test/tests/generic"
	apiext "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	errors2 "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/testing"

	// Needed to run tests in GKE
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

var _ = Describe("Test Kube ResourceClient", func() {

	const (
		namespace1 = "test-ns-1"
		namespace2 = "test-ns-2"
		resource1  = "res-name-1"
		data       = "some data"
		dumbValue  = "I'm dumb"
	)

	var (
		mockResourceCrd = &solov1.Resource{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "testing.solo.io/v1",
				Kind:       "MockResource",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      resource1,
				Namespace: namespace1,
			},
			Spec: &solov1.Spec{
				"data":          data,
				"someDumbField": dumbValue,
			},
		}
	)

	Context("integrations tests", func() {

		if os.Getenv("RUN_KUBE_TESTS") != "1" {
			log.Printf("This test creates kubernetes resources and is disabled by default. To enable, set RUN_KUBE_TESTS=1 in your env.")
			return
		}
		var (
			namespace string
			cfg       *rest.Config
			client    *kube.ResourceClient
		)
		BeforeEach(func() {
			namespace = helpers.RandString(8)
			err := setup.SetupKubeForTest(namespace)
			Expect(err).NotTo(HaveOccurred())

			cfg, err = kubeutils.GetConfig("", "")
			Expect(err).NotTo(HaveOccurred())

			clientset, err := versioned.NewForConfig(cfg, v1.MockResourceCrd)
			Expect(err).NotTo(HaveOccurred())

			// Create the CRD in the cluster
			apiExts, err := apiext.NewForConfig(cfg)
			Expect(err).NotTo(HaveOccurred())

			err = v1.MockResourceCrd.Register(apiExts)
			Expect(err).NotTo(HaveOccurred())

			client = kube.NewResourceClient(v1.MockResourceCrd, clientset, kube.NewKubeCache(context.TODO()), &v1.MockResource{}, []string{metav1.NamespaceAll}, 0)
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			if err := setup.TeardownKube(namespace); err != nil {
				panic(err)
			}
			if err := setup.DeleteCrd(v1.MockResourceCrd.FullName()); err != nil {
				panic(err)
			}
		})

		It("CRUDs resources", func() {
			generic.TestCrudClient(namespace, client, time.Minute)
		})
	})

	Context("multiple namespaces", func() {
		var (
			ns1, ns2 string
			cfg      *rest.Config
			client   *kube.ResourceClient
		)
		BeforeEach(func() {
			ns1 = helpers.RandString(8)
			ns2 = helpers.RandString(8)
			err := setup.SetupKubeForTest(ns1)
			Expect(err).NotTo(HaveOccurred())
			err = setup.SetupKubeForTest(ns2)
			Expect(err).NotTo(HaveOccurred())

			cfg, err = kubeutils.GetConfig("", "")
			Expect(err).NotTo(HaveOccurred())

			clientset, err := versioned.NewForConfig(cfg, v1.MockResourceCrd)
			Expect(err).NotTo(HaveOccurred())

			// Create the CRD in the cluster
			apiExts, err := apiext.NewForConfig(cfg)
			Expect(err).NotTo(HaveOccurred())

			err = v1.MockResourceCrd.Register(apiExts)
			Expect(err).NotTo(HaveOccurred())

			client = kube.NewResourceClient(v1.MockResourceCrd, clientset, kube.NewKubeCache(context.TODO()), &v1.MockResource{}, []string{metav1.NamespaceAll}, 0)
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			setup.TeardownKube(ns1)
			setup.TeardownKube(ns2)
			setup.DeleteCrd(v1.MockResourceCrd.FullName())
		})
		It("can watch resources across namespaces when using NamespaceAll", func() {
			namespace := ""
			boo := "hoo"
			goo := "goo"

			err := client.Register()
			Expect(err).NotTo(HaveOccurred())

			w, errs, err := client.Watch(namespace, clients.WatchOpts{Ctx: context.TODO()})
			Expect(err).NotTo(HaveOccurred())

			var r1, r2 resources.Resource
			wait := make(chan struct{})
			go func() {
				defer GinkgoRecover()
				defer func() {
					close(wait)
				}()
				r1, err = client.Write(&v1.MockResource{
					Data: data,
					Metadata: core.Metadata{
						Name:      boo,
						Namespace: ns1,
					},
				}, clients.WriteOpts{})
				Expect(err).NotTo(HaveOccurred())

				r2, err = client.Write(&v1.MockResource{
					Data: data,
					Metadata: core.Metadata{
						Name:      goo,
						Namespace: ns2,
					},
				}, clients.WriteOpts{})
				Expect(err).NotTo(HaveOccurred())
			}()
			select {
			case <-wait:
			case <-time.After(time.Second * 5):
				Fail("expected wait to be closed before 5s")
			}

			list, err := client.List(namespace, clients.ListOpts{})
			Expect(err).NotTo(HaveOccurred())
			Expect(list).To(ContainElement(r1))
			Expect(list).To(ContainElement(r2))

			select {
			case err := <-errs:
				Expect(err).NotTo(HaveOccurred())
			case list = <-w:
			case <-time.After(time.Millisecond * 5):
				Fail("expected a message in channel")
			}

			var timesDrained int
		drain:
			for {
				select {
				case list = <-w:
					timesDrained++
					if timesDrained > 50 {
						Fail("drained the watch channel 50 times, something is wrong")
					}
				case err := <-errs:
					Expect(err).NotTo(HaveOccurred())
				case <-time.After(time.Second / 4):
					break drain
				}
			}

			Expect(list).To(ContainElement(r1))
			Expect(list).To(ContainElement(r2))
		})
	})

	Context("unit tests", func() {

		var (
			clientset *fake.Clientset
			cache     kube.SharedCache
			rc        *kube.ResourceClient
		)

		BeforeEach(func() {
			clientset = fake.NewSimpleClientset(v1.MockResourceCrd)
			cache = kube.NewKubeCache(context.TODO())
			rc = kube.NewResourceClient(v1.MockResourceCrd, clientset, cache, &v1.MockResource{}, []string{namespace1}, 0)
		})

		It("return the expected kind name", func() {
			Expect(rc.Kind()).To(BeEquivalentTo("*v1.MockResource"))
		})

		It("can be registered", func() {
			Expect(rc.Register()).NotTo(HaveOccurred())
		})

		Describe("invoking operations on non-allowed namespaces causes an error", func() {

			It("call read", func() {
				_, err := rc.Read(namespace2, "test", clients.ReadOpts{})
				Expect(err).To(HaveOccurred())
			})

			It("call write", func() {
				_, err := rc.Write(&v1.MockResource{Metadata: core.Metadata{Namespace: namespace2}}, clients.WriteOpts{})
				Expect(err).To(HaveOccurred())
			})

			It("call list", func() {
				_, err := rc.List(namespace2, clients.ListOpts{})
				Expect(err).To(HaveOccurred())
			})

			It("call delete", func() {
				err := rc.Delete(namespace2, "test", clients.DeleteOpts{})
				Expect(err).To(HaveOccurred())
			})

			It("call watch", func() {
				_, _, err := rc.Watch(namespace2, clients.WatchOpts{})
				Expect(err).To(HaveOccurred())
			})
		})

		Describe("reading a resource", func() {

			var (
				clientset             *fake.Clientset
				malformedResourceName = "malformed-res"
				malformedResourceCrd  = &solov1.Resource{
					TypeMeta: metav1.TypeMeta{
						APIVersion: "testing.solo.io/v1",
						Kind:       "MockResource",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      resource1,
						Namespace: namespace1,
					},
					Spec: &solov1.Spec{
						"unexpectedField": data,
					},
				}
			)

			BeforeEach(func() {
				clientset = fake.NewSimpleClientset(v1.MockResourceCrd)
				clientset.PrependReactor("get", "mocks", func(action testing.Action) (handled bool, ret runtime.Object, err error) {
					switch action := action.(type) {
					case testing.GetActionImpl:
						if action.GetName() == resource1 {
							return true, mockResourceCrd, nil
						}
						if action.GetName() == malformedResourceName {
							return true, malformedResourceCrd, nil
						}
					}
					return true, nil, &errors2.StatusError{ErrStatus: metav1.Status{
						Status: metav1.StatusFailure,
						Reason: metav1.StatusReasonNotFound,
					}}
				})
				rc = kube.NewResourceClient(v1.MockResourceCrd, clientset, cache, &v1.MockResource{}, []string{namespace1}, 0)
				Expect(rc.Register()).NotTo(HaveOccurred())
			})

			It("correctly retrieves an existing resource", func() {
				res, err := rc.Read(namespace1, resource1, clients.ReadOpts{})
				Expect(err).NotTo(HaveOccurred())

				mockRes, ok := res.(*v1.MockResource)
				Expect(ok).To(BeTrue())
				Expect(mockRes.Metadata.Name).To(BeEquivalentTo(mockResourceCrd.Name))
				Expect(mockRes.Metadata.Namespace).To(BeEquivalentTo(mockResourceCrd.Namespace))
				Expect(mockRes.Data).To(BeEquivalentTo((*mockResourceCrd.Spec)["data"]))
				Expect(mockRes.SomeDumbField).To(BeEquivalentTo((*mockResourceCrd.Spec)["someDumbField"]))
			})
			It("return an error when retrieving a non-existing resource", func() {
				_, err := rc.Read(namespace1, "non-existing", clients.ReadOpts{})
				Expect(err).To(HaveOccurred())
				Expect(errors.IsNotExist(err)).To(BeTrue())
			})

			It("return an error when receiving a malformed resource", func() {
				_, err := rc.Read(namespace1, malformedResourceName, clients.ReadOpts{})
				Expect(err).To(HaveOccurred())
				Expect(errors.IsNotExist(err)).To(BeFalse())
			})
		})

		Describe("writing a resource", func() {

			var (
				clientset *fake.Clientset

				resourceToCreate = &v1.MockResource{
					Metadata: core.Metadata{
						Name:      "to-create",
						Namespace: namespace1,
					},
					Data:          data,
					SomeDumbField: dumbValue,
				}
				resourceToUpdate = &v1.MockResource{
					Metadata: core.Metadata{
						Name:      "mock-1",
						Namespace: namespace1,
					},
					Data:          data,
					SomeDumbField: dumbValue,
				}
			)

			BeforeEach(func() {
				clientset = fake.NewSimpleClientset(v1.MockResourceCrd)

				// Create an initial resource with the name of resourceToUpdate
				err := util.CreateMockResource(clientset, namespace1, resourceToUpdate.Metadata.Name, "to-be-updated")
				Expect(err).NotTo(HaveOccurred())

				rc = kube.NewResourceClient(v1.MockResourceCrd, clientset, cache, &v1.MockResource{}, []string{namespace1}, 0)
				Expect(rc.Register()).NotTo(HaveOccurred())
			})

			Context("resource does not exist", func() {

				It("correctly creates the resource", func() {
					res, err := rc.Write(resourceToCreate, clients.WriteOpts{})
					Expect(err).NotTo(HaveOccurred())
					Expect(res).NotTo(BeNil())
					mockRes, ok := res.(*v1.MockResource)
					Expect(ok).To(BeTrue())
					Expect(mockRes.Metadata.Name).To(BeEquivalentTo(resourceToCreate.Metadata.Name))
					Expect(mockRes.Metadata.Namespace).To(BeEquivalentTo(resourceToCreate.Metadata.Namespace))
					Expect(mockRes.Data).To(BeEquivalentTo(resourceToCreate.Data))
					Expect(mockRes.SomeDumbField).To(BeEquivalentTo(resourceToCreate.SomeDumbField))
				})
			})

			Context("resource exists and we want to overwrite", func() {

				It("correctly updates the resource", func() {
					res, err := rc.Write(resourceToUpdate, clients.WriteOpts{OverwriteExisting: true})
					Expect(err).NotTo(HaveOccurred())
					Expect(res).NotTo(BeNil())

					checkRes, err := rc.Read(namespace1, resourceToUpdate.Metadata.Name, clients.ReadOpts{})
					Expect(err).NotTo(HaveOccurred())
					Expect(checkRes).NotTo(BeNil())
					checkMockRes, ok := res.(*v1.MockResource)
					Expect(ok).To(BeTrue())
					Expect(checkMockRes.SomeDumbField).To(BeEquivalentTo(resourceToUpdate.SomeDumbField))
				})
			})

			Context("resource exists and we don't want to overwrite", func() {

				It("returns the appropriate error", func() {
					_, err := rc.Write(resourceToUpdate, clients.WriteOpts{OverwriteExisting: false})
					Expect(err).To(HaveOccurred())
					Expect(errors.IsExist(err)).To(BeTrue())
				})
			})
		})

		Describe("listing resources", func() {

			var clientset *fake.Clientset

			BeforeEach(func() {
				clientset = fake.NewSimpleClientset(v1.MockResourceCrd)

				// Create initial resources
				Expect(util.CreateMockResource(clientset, namespace1, "res-1", "val-1")).NotTo(HaveOccurred())
				Expect(util.CreateMockResource(clientset, namespace1, "res-2", "val-2")).NotTo(HaveOccurred())
				Expect(util.CreateMockResource(clientset, namespace1, "res-3", "val-3")).NotTo(HaveOccurred())
				Expect(util.CreateMockResource(clientset, namespace2, "res-4", "val-4")).NotTo(HaveOccurred())

				rc = kube.NewResourceClient(v1.MockResourceCrd, clientset, cache, &v1.MockResource{}, []string{namespace1, namespace2, "empty"}, 0)
				Expect(rc.Register()).NotTo(HaveOccurred())
			})

			It("lists the correct resources for the given namespace", func() {
				list, err := rc.List(namespace1, clients.ListOpts{})
				Expect(err).NotTo(HaveOccurred())
				Expect(list).To(HaveLen(3))

				list, err = rc.List(namespace2, clients.ListOpts{})
				Expect(err).NotTo(HaveOccurred())
				Expect(list).To(HaveLen(1))

				list, err = rc.List("empty", clients.ListOpts{})
				Expect(err).NotTo(HaveOccurred())
				Expect(list).To(HaveLen(0))
			})
		})

		Describe("deleting resources", func() {

			var clientset *fake.Clientset

			BeforeEach(func() {
				clientset = fake.NewSimpleClientset(v1.MockResourceCrd)

				// Create initial resource
				Expect(util.CreateMockResource(clientset, namespace1, "res-1", "val-1")).NotTo(HaveOccurred())

				rc = kube.NewResourceClient(v1.MockResourceCrd, clientset, cache, &v1.MockResource{}, []string{namespace1}, 0)
				Expect(rc.Register()).NotTo(HaveOccurred())
			})

			Context("resource exists", func() {

				It("correctly deletes an existing resource", func() {
					err := rc.Delete(namespace1, "res-1", clients.DeleteOpts{})
					Expect(err).NotTo(HaveOccurred())

					// Verify whether resource was actually deleted
					_, err = rc.Read(namespace1, "res-1", clients.ReadOpts{})
					Expect(errors.IsNotExist(err)).To(BeTrue())
				})
			})

			Context("resource does not exist", func() {

				It("returns error when trying to delete", func() {
					err := rc.Delete(namespace1, "res-X", clients.DeleteOpts{})
					Expect(err).To(HaveOccurred())
					Expect(errors.IsNotExist(err)).To(BeTrue())
				})

				It("does not error when passing the correspondent option", func() {
					err := rc.Delete(namespace1, "res-X", clients.DeleteOpts{IgnoreNotExist: true})
					Expect(err).NotTo(HaveOccurred())
				})
			})
		})

		Describe("watching resources", func() {

			var clientset *fake.Clientset

			BeforeEach(func() {
				clientset = fake.NewSimpleClientset(v1.MockResourceCrd)

				rc = kube.NewResourceClient(v1.MockResourceCrd, clientset, cache, &v1.MockResource{}, []string{namespace1, namespace2}, 0)
				Expect(rc.Register()).NotTo(HaveOccurred())
			})

			It("correctly receives notifications for resources in the given namespace", func() {
				resources, errors, err := rc.Watch(namespace1, clients.WatchOpts{})
				Expect(err).NotTo(HaveOccurred())

				// Create a resource
				go Expect(util.CreateMockResource(clientset, namespace1, "res-1", "val-1")).NotTo(HaveOccurred())

				skippedInitialRead := false
				after := time.After(200 * time.Millisecond)
			LOOP:
				for {
					select {
					case res := <-resources:
						if skippedInitialRead {
							Expect(res).To(HaveLen(1))
							Expect(res[0].GetMetadata().Name).To(BeEquivalentTo("res-1"))
							break LOOP
						}
						Expect(res).To(HaveLen(0))
						skippedInitialRead = true
						continue
					case <-errors:
						Fail("unexpected error on watch error channel")
					case <-after:
						Fail("timed out waiting for event notification")
					}
				}
			})

			It("does not receives notifications for resources other namespaces", func() {
				resources, errors, err := rc.Watch(namespace1, clients.WatchOpts{})
				Expect(err).NotTo(HaveOccurred())

				// Create a resource
				go Expect(util.CreateMockResource(clientset, namespace2, "res-1", "val-1")).NotTo(HaveOccurred())

				skippedInitialRead := false
				after := time.After(200 * time.Millisecond)
			LOOP:
				for {
					select {
					case res := <-resources:
						if skippedInitialRead {
							Fail("timed out waiting for event notification")
						}
						Expect(res).To(HaveLen(0))
						skippedInitialRead = true
						continue
					case <-errors:
						Fail("unexpected error on watch error channel")
					case <-after:
						break LOOP
					}
				}
			})
		})
	})
})

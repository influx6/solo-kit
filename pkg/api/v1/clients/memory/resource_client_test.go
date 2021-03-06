package memory_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	. "github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	core "github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	v1 "github.com/solo-io/solo-kit/test/mocks/v1"
	"github.com/solo-io/solo-kit/test/tests/generic"
)

var _ = Describe("Base", func() {
	var (
		client *ResourceClient
	)
	BeforeEach(func() {
		client = NewResourceClient(NewInMemoryResourceCache(), &v1.MockResource{})
	})
	AfterEach(func() {
	})
	It("CRUDs resources", func() {
		generic.TestCrudClient("test", client, time.Minute)
	})
	It("should not return pointer to internal object", func() {
		obj := &v1.MockResource{
			Metadata: core.Metadata{
				Namespace: "ns",
				Name:      "n",
			},
			Data: "test",
		}
		client.Write(obj, clients.WriteOpts{})
		ret, err := client.Read("ns", "n", clients.ReadOpts{})
		Expect(err).NotTo(HaveOccurred())
		Expect(ret).NotTo(BeIdenticalTo(obj))

		ret2, err := client.Read("ns", "n", clients.ReadOpts{})
		Expect(err).NotTo(HaveOccurred())
		Expect(ret).NotTo(BeIdenticalTo(ret2))

		listret, err := client.List("ns", clients.ListOpts{})
		Expect(err).NotTo(HaveOccurred())
		Expect(listret[0]).NotTo(BeIdenticalTo(obj))

		listret2, err := client.List("ns", clients.ListOpts{})
		Expect(err).NotTo(HaveOccurred())
		Expect(listret[0]).NotTo(BeIdenticalTo(listret2[0]))

	})
})

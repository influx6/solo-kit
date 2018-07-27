package mocks

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
)

var _ = Describe("MocksEventLoop", func() {
	var (
		namespace string
		cache     Cache
		err       error
	)

	BeforeEach(func() {
		mockResourceClientFactory := factory.NewResourceClientFactory(&factory.MemoryResourceClientOpts{})
		mockResourceClient := NewMockResourceClient(mockResourceClientFactory)
		fakeResourceClientFactory := factory.NewResourceClientFactory(&factory.MemoryResourceClientOpts{})
		fakeResourceClient := NewFakeResourceClient(fakeResourceClientFactory)
		cache = NewCache(mockResourceClient, fakeResourceClient)
	})
	It("runs sync function on a new snapshot", func() {
		_, err = cache.MockResource().Write(NewMockResource(namespace, "jerry"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())
		_, err = cache.FakeResource().Write(NewFakeResource(namespace, "jerry"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())
		sync := &mockSyncer{}
		el := NewEventLoop(cache, sync)
		go func() {
			defer GinkgoRecover()
			err := el.Run(clients.WatchOpts{Namespace: namespace})
			Expect(err).NotTo(HaveOccurred())
		}()
		Eventually(func() bool { return sync.synced }, time.Second).Should(BeTrue())
	})
})

type mockSyncer struct {
	synced bool
}

func (s *mockSyncer) Sync(snap *Snapshot) error {
	s.synced = true
	return nil
}
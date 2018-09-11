package v1

import (
	"sync"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/errutils"
)

type DiscoveryEmitter interface {
	Register() error
	Upstream() UpstreamClient
	Snapshots(watchNamespaces []string, opts clients.WatchOpts) (<-chan *DiscoverySnapshot, <-chan error, error)
}

func NewDiscoveryEmitter(upstreamClient UpstreamClient) DiscoveryEmitter {
	return NewDiscoveryEmitterWithEmit(upstreamClient, make(chan struct{}))
}

func NewDiscoveryEmitterWithEmit(upstreamClient UpstreamClient, emit <-chan struct{}) DiscoveryEmitter {
	return &discoveryEmitter{
		upstream:  upstreamClient,
		forceEmit: emit,
	}
}

type discoveryEmitter struct {
	forceEmit <-chan struct{}
	upstream  UpstreamClient
}

func (c *discoveryEmitter) Register() error {
	if err := c.upstream.Register(); err != nil {
		return err
	}
	return nil
}

func (c *discoveryEmitter) Upstream() UpstreamClient {
	return c.upstream
}

func (c *discoveryEmitter) Snapshots(watchNamespaces []string, opts clients.WatchOpts) (<-chan *DiscoverySnapshot, <-chan error, error) {
	errs := make(chan error)
	var done sync.WaitGroup
	/* Create channel for Upstream */
	type upstreamListWithNamespace struct {
		list      UpstreamList
		namespace string
	}
	upstreamChan := make(chan upstreamListWithNamespace)

	for _, namespace := range watchNamespaces {
		/* Setup watch for Upstream */
		upstreamNamespacesChan, upstreamErrs, err := c.upstream.Watch(namespace, opts)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "starting Upstream watch")
		}

		done.Add(1)
		go func(namespace string) {
			defer done.Done()
			errutils.AggregateErrs(opts.Ctx, errs, upstreamErrs, namespace+"-upstreams")
		}(namespace)

		/* Watch for changes and update snapshot */
		go func(namespace string) {
			for {
				select {
				case <-opts.Ctx.Done():
					return
				case upstreamList := <-upstreamNamespacesChan:
					select {
					case <-opts.Ctx.Done():
						return
					case upstreamChan <- upstreamListWithNamespace{list: upstreamList, namespace: namespace}:
					}
				}
			}
		}(namespace)
	}

	snapshots := make(chan *DiscoverySnapshot)
	go func() {
		currentSnapshot := DiscoverySnapshot{}
		sync := func(newSnapshot DiscoverySnapshot) {
			if currentSnapshot.Hash() == newSnapshot.Hash() {
				return
			}
			currentSnapshot = newSnapshot
			sentSnapshot := currentSnapshot.Clone()
			snapshots <- &sentSnapshot
		}
		for {
			select {
			case <-opts.Ctx.Done():
				close(snapshots)
				done.Wait()
				close(errs)
				return
			case <-c.forceEmit:
				sentSnapshot := currentSnapshot.Clone()
				snapshots <- &sentSnapshot
			case upstreamNamespacedList := <-upstreamChan:
				namespace := upstreamNamespacedList.namespace
				upstreamList := upstreamNamespacedList.list

				newSnapshot := currentSnapshot.Clone()
				newSnapshot.Upstreams.Clear(namespace)
				newSnapshot.Upstreams.Add(upstreamList...)
				sync(newSnapshot)
			}
		}
	}()
	return snapshots, errs, nil
}
package thunder

import (
	"context"
	"errors"
	"golang.org/x/sync/errgroup"
	"io"
	"sync/atomic"
	"time"
)

func (p *Processor) StreamDocuments(ctx context.Context, exporter Exporter, w io.Writer, limit uint64) error {
	// Preload exporter
	if err := exporter.Load(w); err != nil {
		return err
	}

	// Start indexing
	eg, egCtx := errgroup.WithContext(ctx)

	// Start source
	var inChan = make(chan *Document)
	eg.Go(func() error {
		defer close(inChan)
		return p.Source.Driver.GetDocumentsForProcessor(p, inChan, egCtx, limit)
	})

	// Start exporter broadcasting
	eg.Go(func() error {
		var position atomic.Uint64
		for {
			select {
			case <-egCtx.Done():
				return egCtx.Err()
			case <-time.After(time.Second * 10):
				return context.DeadlineExceeded
			case doc, open := <-inChan:
				if !open {
					if position.Load() >= 1 {
						return exporter.AfterAll()
					}
					return nil
				}
				position.Add(1)
				if position.Load() == 1 {
					if err := exporter.BeforeAll(); err != nil {
						return err
					}
				}

				if err := exporter.WriteDocument(doc, position.Load()); err != nil {
					return err
				}
			}
		}
	})

	err := eg.Wait()
	if err != nil {
		if errors.Is(err, context.Canceled) {
			err = context.Cause(egCtx)
		}
		return err
	}
	return nil
}

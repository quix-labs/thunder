package thunder

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

type Condition struct {
	Column   string `json:"column"`
	Operator string `json:"operator"`
	Value    string `json:"value,omitempty"`
}

type Processor struct {
	ID int

	Indexing atomic.Bool

	Listening     atomic.Bool
	ContextCancel context.CancelFunc

	Source  *Source
	Targets []*Target

	Table       string
	PrimaryKeys []string
	Conditions  []Condition

	Mapping Mapping

	Index string

	Enabled bool
}

type Document struct {
	Pkey string `json:"_pkey"`
	Json []byte `json:"json"`
}

func (p *Processor) Start() error {
	if p.Listening.Load() {
		return errors.New("processor is already listening")
	}

	p.Listening.Store(true)
	GetBroadcaster().Dispatch("processor-updated", p.ID)

	defer func() {
		p.Listening.Store(false)
		GetBroadcaster().Dispatch("processor-updated", p.ID)
	}()

	// Create the context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	p.ContextCancel = cancel
	defer p.ContextCancel()

	var eventsChan = make(chan DbEvent)
	defer close(eventsChan)

	// Bootstrap target channels
	var targetChannels = make([]chan TargetEvent, len(p.Targets))
	for i, _ := range p.Targets {
		targetChannels[i] = make(chan TargetEvent, 1)
	}

	// Start target in parallel
	var wg sync.WaitGroup
	for i, target := range p.Targets {
		wg.Add(1)
		idx := i
		go func() {
			defer wg.Done()
			defer close(targetChannels[idx])
			if err := target.Driver.HandleEvents(p, targetChannels[idx], ctx); err != nil && !errors.Is(err, context.Canceled) {
				cancel()
			}
		}()
	}

	// Start broadcasting events
	go func() {
		for event := range eventsChan {
			for _, channel := range targetChannels {
				switch typedEvent := event.(type) {
				case *DbInsertEvent:
					fmt.Println("inser event received")
					// TODO FETCH FULL USING PRIMARY KEYS

					//channel <- &TargetInsertEvent{
					//	PrimaryKeys: typedEvent.PrimaryKeys,
					//	Json:        typedEvent.Json,
					//}
				case *DbPatchEvent:
					channel <- &TargetPatchEvent{
						Path:      typedEvent.Path,
						Pkey:      typedEvent.Pkey,
						JsonPatch: typedEvent.JsonPatch,
					}

				case *DbDeleteEvent:
					channel <- &TargetDeleteEvent{
						Pkey: typedEvent.Pkey,
					}

				case *DbTruncateEvent:
					channel <- &TargetTruncateEvent{}
				}

			}
		}
	}()

	// Start driver event handling
	err := p.Source.Driver.Start(p, eventsChan, ctx)
	if err != nil && !errors.Is(err, context.Canceled) {
		cancel()
	}
	wg.Wait() // Wait for processor stopped

	return err
}

func (p *Processor) Stop() error {
	if p.Listening.Load() {
		p.ContextCancel()
		GetBroadcaster().Dispatch("processor-updated", p.ID)
	}
	return nil
}

func (p *Processor) FullIndex() error {
	if p.Indexing.Load() {
		return errors.New("processor is already indexing")
	}
	p.Indexing.Store(true)
	GetBroadcaster().Dispatch("processor-updated", p.ID)
	defer func() {
		p.Indexing.Store(false)
		GetBroadcaster().Dispatch("processor-indexed", p.ID)
		GetBroadcaster().Dispatch("processor-updated", p.ID)
		//TODO RESTART IF STOPPED NOT WORK
	}()

	// Start indexing
	docChan := make(chan *Document, 1)
	errChan := make(chan error, 1)

	go func() {
		defer close(docChan)
		p.Source.Driver.GetDocumentsForProcessor(p, docChan, errChan, 0)
	}()

	// Initialize channels
	var targetEventsChans = make([]chan TargetEvent, len(p.Targets))
	for idx, _ := range p.Targets {
		targetEventsChans[idx] = make(chan TargetEvent, 1)
	}

	// Start targets in parallel
	go func() {
		defer close(errChan)
		var wg sync.WaitGroup
		for idx, target := range p.Targets {
			wg.Add(1)
			go func() {
				defer wg.Done()
				if err := target.Driver.HandleEvents(p, targetEventsChans[idx], context.Background()); err != nil {
					errChan <- err
				}

			}()
		}
		wg.Wait()
	}()
	for {
		select {
		case doc, open := <-docChan:
			if !open {
				// Send end signal
				for _, targetEventChan := range targetEventsChans {
					close(targetEventChan)
				}

				// Wait for all closed or err
				return <-errChan
			}

			// Dispatch across different targets
			for _, targetEventChan := range targetEventsChans {
				event := &TargetInsertEvent{
					Pkey: doc.Pkey,
					Json: doc.Json,
				}
				targetEventChan <- event
			}

		case err, opened := <-errChan:
			if !opened {
				return nil
			}
			if err != nil {
				// TODO STOP docChan
				for _, targetEventChan := range targetEventsChans {
					close(targetEventChan)
				}
				<-errChan
				return err
			}
		}
	}
}

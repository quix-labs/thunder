package thunder

import (
	"context"
	"errors"
	"fmt"
	"github.com/oklog/ulid/v2"
	"github.com/quix-labs/thunder/utils"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

type Condition struct {
	Column   string `json:"column"`
	Operator string `json:"operator"`
	Value    string `json:"value,omitempty"`
}

var Processors = utils.NewRegistry[*Processor]("processor").SetIdGenerator(func(item **Processor) (string, error) {
	ulid := ulid.Make().String()
	(*item).ID = ulid
	return ulid, nil
})

type Processor struct {
	ID string

	Indexing atomic.Bool

	Listening         atomic.Bool
	ListenBroadcaster *utils.Broadcaster[DbEvent, TargetEvent]

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

	defer func() {
		p.Listening.Store(false)
	}()

	// Initialize broadcaster
	p.ListenBroadcaster = utils.NewBroadcaster[DbEvent, TargetEvent](func(event DbEvent) TargetEvent {
		switch typedEvent := event.(type) {
		case *DbInsertEvent:
			fmt.Println("insert event received")
			// TODO FETCH FULL USING PRIMARY KEYS

			//return &TargetInsertEvent{
			//	PrimaryKeys: typedEvent.PrimaryKeys,
			//	Json:        typedEvent.Json,
			//}
		case *DbPatchEvent:
			return &TargetPatchEvent{
				Relation:  typedEvent.Relation,
				Pkey:      typedEvent.Pkey,
				JsonPatch: typedEvent.JsonPatch,
			}

		case *DbDeleteEvent:
			return &TargetDeleteEvent{
				Pkey:     typedEvent.Pkey,
				Relation: typedEvent.Relation,
			}

		case *DbTruncateEvent:
			return &TargetTruncateEvent{
				Relation: typedEvent.Relation,
			}
		}

		return nil
	})
	p.ListenBroadcaster.Start()
	defer p.ListenBroadcaster.Close()

	// Start targets in parallel
	for _, target := range p.Targets {
		go func() {
			listenChan, stopListening := p.ListenBroadcaster.NewListenChan()
			if err := target.Driver.HandleEvents(p, listenChan, context.TODO()); err != nil {
				stopListening()
				//				broadcaster.Close() Uncomment to stop emission
				//TODO error
			}
		}()
	}

	// Start driver event handling
	if err := p.Source.Driver.Start(p, p.ListenBroadcaster.In()); err != nil {
		return err
	}

	return nil
}

func (p *Processor) Stop() error {
	if p.Listening.Load() {
		p.ListenBroadcaster.Close()
		// Wait for closed
		for p.Listening.Load() {
			time.Sleep(10 * time.Millisecond)
		}
	}

	return nil
}

func (p *Processor) FullIndex() error {
	if p.Indexing.Load() {
		return errors.New("processor is already indexing")
	}
	p.Indexing.Store(true)
	defer p.Indexing.Store(false)

	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(nil)

	eg, egCtx := errgroup.WithContext(ctx)

	// Initialize each channel individually
	var targetChans = make([]chan TargetEvent, len(p.Targets))
	for i := range targetChans {
		targetChans[i] = make(chan TargetEvent) // create an actual channel for each slice element
	}

	// Start targets in parallel
	for i, target := range p.Targets {
		targetChan := targetChans[i] // capture channel locally
		eg.Go(func() error {
			return target.Driver.HandleEvents(p, targetChan, egCtx)
		})
	}

	// Start source in background
	inChan := make(chan *Document)
	eg.Go(func() error {
		defer close(inChan)
		return p.Source.Driver.GetDocumentsForProcessor(p, inChan, egCtx, 0)
	})

	var docCount atomic.Uint64

	// Start dispatcher
	eg.Go(func() error {
		defer func() {
			for _, targetChan := range targetChans {
				close(targetChan) // close target channels when done
			}
		}()

		for {
			select {
			case <-egCtx.Done():
				return egCtx.Err()
			case <-time.After(time.Second * 10):
				cancel(context.DeadlineExceeded)
			case doc, ok := <-inChan:
				if !ok {
					return nil
				}
				docCount.Add(1)

				// Send event across all channels in parallel
				var wg sync.WaitGroup
				for _, targetChan := range targetChans {
					wg.Add(1)
					go func(chanToSend chan TargetEvent) {
						defer wg.Done()
						chanToSend <- &TargetInsertEvent{
							Pkey: doc.Pkey,
							Json: doc.Json,
						}
					}(targetChan)
				}
				wg.Wait()
			}
		}
	})

	//TODO REMOVE STATS
	start := time.Now()
	err := eg.Wait()
	totalTime := time.Since(start)
	totalTimeMs := totalTime.Milliseconds()
	totalSeconds := totalTime.Seconds()
	docsPerSec := float64(docCount.Load()) / totalSeconds
	meanDocTimeMs := float64(totalTimeMs) / float64(docCount.Load())

	fmt.Println("Total Time (s):", totalSeconds)
	fmt.Println("Documents Processed:", docCount.Load())
	fmt.Println("Documents per second:", docsPerSec)
	fmt.Println("Mean Time per Document (ms):", meanDocTimeMs)
	fmt.Println("Total Time (ms):", totalTimeMs)

	if errors.Is(err, context.Canceled) {
		return context.Cause(egCtx)
	}
	return err
}

func GetLoggerForProcessor(processor *Processor) *zerolog.Logger {
	logger := zerolog.New(os.Stdout).Level(zerolog.DebugLevel).With().Str("processor", processor.ID).Stack().Timestamp().Logger()
	return &logger
}

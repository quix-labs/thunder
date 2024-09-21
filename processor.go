package thunder

import (
	"errors"
	"fmt"
	"github.com/quix-labs/thunder/utils"
	"sync/atomic"
	"time"
)

type Condition struct {
	Column   string `json:"column"`
	Operator string `json:"operator"`
	Value    string `json:"value,omitempty"`
}

type Processor struct {
	ID int

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
	GetBroadcaster().Dispatch("processor-updated", p.ID)

	defer func() {
		p.Listening.Store(false)
		GetBroadcaster().Dispatch("processor-updated", p.ID)
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
			if err := target.Driver.HandleEvents(p, listenChan); err != nil {
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
	broadcaster := utils.NewBroadcaster[*Document, TargetEvent](func(doc *Document) TargetEvent {
		return &TargetInsertEvent{
			Pkey: doc.Pkey,
			Json: doc.Json,
		}
	})
	broadcaster.Start()
	defer broadcaster.Close()

	// Start targets
	for _, target := range p.Targets {
		go func() {
			listenChan, stopListening := broadcaster.NewListenChan()
			defer stopListening()
			if err := target.Driver.HandleEvents(p, listenChan); err != nil {
				stopListening()
				//				broadcaster.Close() Uncomment to stop emission
				//TODO error handling
			}
		}()
	}

	// Start source
	go func() {
		defer broadcaster.In().Finish()
		if err := p.Source.Driver.GetDocumentsForProcessor(p, broadcaster.In(), 0); err != nil {
			fmt.Println("error in getting documents for processor")
		}
	}()

	broadcaster.Wait()
	return nil
}

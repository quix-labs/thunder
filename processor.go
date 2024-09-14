package thunder

import (
	"errors"
	"sync"
)

type ProcessorStatus string

var (
	ProcessorIndexing  = ProcessorStatus("indexing")
	ProcessorListening = ProcessorStatus("listening")
	ProcessorInactive  = ProcessorStatus("inactive")
)

type Condition struct {
	Column   string `json:"column"`
	Operator string `json:"operator"`
	Value    string `json:"value,omitempty"`
}

type Processor struct {
	ID int

	Status ProcessorStatus

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
	PrimaryKeys []string `json:"primary_keys"`
	Json        []byte   `json:"json"`
}

func (p *Processor) FullIndex() error {
	if p.Status != ProcessorInactive {
		return errors.New("processor is active")
	}

	p.Status = ProcessorIndexing
	GetBroadcaster().Dispatch("processor-updated", p.ID)
	defer func() {
		p.Status = ProcessorInactive
		GetBroadcaster().Dispatch("processor-indexed", p.ID)
		GetBroadcaster().Dispatch("processor-updated", p.ID)
	}()

	// Start indexing
	docChan := make(chan *Document, 1)
	errChan := make(chan error, 1)

	go func() {
		defer close(docChan)
		p.Source.Driver.GetDocumentsForProcessor(p, docChan, errChan, 0)
	}()

	// Initialize channels
	var targetsDocChan = make([]chan *Document, len(p.Targets))
	for idx, _ := range p.Targets {
		targetsDocChan[idx] = make(chan *Document, 1)
	}

	// Start targets in parallel
	go func() {
		defer close(errChan)
		var wg sync.WaitGroup
		for idx, target := range p.Targets {
			wg.Add(1)
			go func() {
				defer wg.Done()
				target.Driver.IndexDocumentsForProcessor(p, targetsDocChan[idx], errChan)
			}()
		}
		wg.Wait()
	}()
	for {
		select {
		case doc, open := <-docChan:
			if !open {
				// Send end signal
				for _, targetDocChan := range targetsDocChan {
					close(targetDocChan)
				}

				// Wait for all closed or err
				return <-errChan
			}

			// Dispatch across different targets
			for _, targetDocChan := range targetsDocChan {
				targetDocChan <- doc
			}

		case err, opened := <-errChan:
			if !opened {
				return nil
			}
			if err != nil {
				// TODO STOP docChan
				for _, targetDocChan := range targetsDocChan {
					close(targetDocChan)
				}
				<-errChan
				return err
			}
		}
	}
}

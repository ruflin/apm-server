package beater

import (
	"errors"
	"sync"

	"github.com/elastic/beats/libbeat/beat"
)

type publisher struct {
	events chan []beat.Event
	client beat.Client
	wg     sync.WaitGroup
}

var (
	errFull = errors.New("Queue is full")
)

func newPublisher(pipeline beat.Pipeline) (*publisher, error) {
	client, err := pipeline.ConnectWith(beat.ClientConfig{
		PublishMode: beat.GuaranteedSend,

		// We want to wait for events in pipeline on shutdown?
		// If set `Close` will block for the duration or until pipeline is empty
		WaitClose: 0,
	})
	if err != nil {
		return nil, err
	}

	p := &publisher{
		events: make(chan []beat.Event, 20),
		client: client,
	}

	p.wg.Add(1)
	go p.run()
	return p, nil
}

func (p *publisher) Stop() {
	close(p.events)
	p.client.Close()
	p.wg.Wait()
}

func (p *publisher) Send(batch []beat.Event) error {
	select {
	case p.events <- batch:
		return nil
	default:
		return errFull
	}
}

func (p *publisher) run() {
	defer p.wg.Done()
	for batch := range p.events {
		p.client.PublishAll(batch)
	}
}

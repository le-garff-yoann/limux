package broadcaster

import (
	"limux/processor"
	"sync"

	uuid "github.com/satori/go.uuid"
)

// Broadcaster broadcast messages sent by the publisher to the receivers.
type Broadcaster struct {
	Pub   chan processor.Event
	recvs map[uuid.UUID]chan processor.Event
	mux   sync.Mutex
}

// New is the constructor to `Broadcaster`.
func New() *Broadcaster {
	return &Broadcaster{Pub: make(chan processor.Event), recvs: make(map[uuid.UUID]chan processor.Event)}
}

// Recv append (create) and returns a receiver.
func (s *Broadcaster) Recv() (<-chan processor.Event, uuid.UUID) {
	s.mux.Lock()
	defer s.mux.Unlock()

	uuid := uuid.NewV4()

	s.recvs[uuid] = make(chan processor.Event)

	return s.recvs[uuid], uuid
}

// RemoveRecv remove and close a receiver.
func (s *Broadcaster) RemoveRecv(uuid uuid.UUID) bool {
	_, ok := s.recvs[uuid]

	if ok {
		s.mux.Lock()
		defer s.mux.Unlock()

		close(s.recvs[uuid])
		delete(s.recvs, uuid)
	}

	return ok
}

// RemoveAllRecv close and revmove all the receivers.
func (s *Broadcaster) RemoveAllRecv() {
	s.mux.Lock()
	defer s.mux.Unlock()

	for _, rs := range s.recvs {
		close(rs)
	}

	s.recvs = make(map[uuid.UUID]chan processor.Event)
}

// Run start broadcasting messages.
func (s *Broadcaster) Run() {
	for {
		m := <-s.Pub

		var wg sync.WaitGroup

		s.mux.Lock()

		for _, rs := range s.recvs {
			wg.Add(1)

			go func(rs chan processor.Event) {
				rs <- m

				wg.Done()
			}(rs)
		}

		wg.Wait()
		s.mux.Unlock()
	}
}

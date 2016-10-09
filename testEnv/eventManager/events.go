package eventManager

import (
	"github.com/op/go-logging"
	connectInit "hyperchain/protos"
	//"reflect"
)

var logger *logging.Logger // package-level logger

func init() {
	logger = logging.MustGetLogger("consensus/events")
}

const (
	ConsensusMsgEventConst="ConsensusMsgEvent"
)

// Event is a type meant to clearly convey that the return type or parameter to a function will be supplied to/from an events.Manager
type Event interface{}

// Receiver is a consumer of events, ProcessEvent will be called serially
type Receiver interface {
	// ProcessEvent delivers an event to the Receiver, if it returns non-nil, the return is the next processed event
	ProcessEvent(e Event) Event
}

/////////////////////////////////////////////////////
//                                                 //
//                 Threaded object                 //
//                                                 //
/////////////////////////////////////////////////////


// threaded holds an exit channel to allow threads to break from a select
type threaded struct {
	exit chan struct{}
}

// halt tells the threaded object's thread to exit
func (t *threaded) Halt() {
	select {
	case <-t.exit:
		logger.Warning("Attempted to halt a threaded object twice")
	default:
		close(t.exit)
	}
}

/////////////////////////////////////////////////////
//                                                 //
//                 Event Manager                   //
//                                                 //
/////////////////////////////////////////////////////

// Manager provides a serialized interface for submitting events to
// a Receiver on the other side of the queue
type Manager interface {
	Inject(Event)         // A temporary interface to allow the event manager thread to skip the queue
	Queue() chan<- Event  // Get a write-only reference to the queue, to submit events
	Start()               // Starts the Manager thread TODO, these thread management things should probably go away
	Halt()                // Stops the Manager thread
	RegistReceiver(e string,receiver Receiver)
}

// managerImpl is an implementation of Manger
type managerImpl struct {
	threaded
	receiver map[string] Receiver
	events   chan Event
}

type ConsensusMsgEvent struct {
	Msg *connectInit.Message
}
type StateUpdatedEvent struct {
	Msg *connectInit.Message
}

// NewManagerImpl creates an instance of managerImpl
func NewManagerImpl() Manager {
	return &managerImpl{
		receiver: make(map[string] Receiver),
		events:   make(chan Event),
		threaded: threaded{make(chan struct{})},
	}
}


// SetReceiver sets the destination for events
func (em *managerImpl) RegistReceiver(e string,receiver Receiver) {
	em.receiver[e]=receiver
}

// Start creates the go routine necessary to deliver events
func (em *managerImpl) Start() {
	go em.eventLoop()
}

// queue returns a write only reference to the event queue
func (em *managerImpl) Queue() chan<- Event {
	return em.events
}

// SendEvent performs the event loop on a receiver to completion
func (em *managerImpl) SendEvent(event Event) {

	switch event.(type) {
	case *ConsensusMsgEvent:
		//logger.Info("------>SendEvent",reflect.TypeOf(em.receiver[ConsensusMsgEventConst]),len(em.receiver),em.receiver[ConsensusMsgEventConst]==nil)
		em.receiver[ConsensusMsgEventConst].ProcessEvent(event)
	case *StateUpdatedEvent:
		logger.Info("comming a StateUpdatedEvent")
		em.receiver[ConsensusMsgEventConst].ProcessEvent(event)
	default:
		logger.Info("there id no default now!")
	}

}

// Inject can only safely be called by the managerImpl thread itself, it skips the queue
func (em *managerImpl) Inject(event Event) {
	if em.receiver != nil {
		go em.SendEvent(event)
	}
}

// eventLoop is where the event thread loops, delivering events
func (em *managerImpl) eventLoop() {
	for {
		select {
		case next := <-em.events:
			//logger.Info("*********> comming event *********> ",reflect.TypeOf(next))
			em.Inject(next)
		case <-em.exit:
			logger.Debug("eventLoop told to exit")
			return
		}
	}
}

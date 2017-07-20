package pubnub

import (
	"net/http"

	"github.com/pubnub/go/pnerr"
)

// Default constants
const (
	Version     = "4.0.0-alpha"
	MaxSequence = 65535
)

// Errors
var (
	ErrMissingPubKey  = pnerr.NewValidationError("pubnub: Missing Publish Key")
	ErrMissingSubKey  = pnerr.NewValidationError("pubnub: Missing Subscribe Key")
	ErrMissingChannel = pnerr.NewValidationError("pubnub: Missing Channel")
	ErrMissingMessage = pnerr.NewValidationError("pubnub: Missing Message")
)

// TODO: pn.UnsubscribeAll() to be deferred

// No server connection will be established when you create a new PubNub object.
// To establish a new connection use Subscribe() function of PubNub type.
type PubNub struct {
	Config              *Config
	publishSequence     chan int
	subscriptionManager *SubscriptionManager
	client              *http.Client
	subscribeClient     *http.Client
}

// TODO: replace result with a pointer
func (pn *PubNub) Publish(opts *PublishOpts) (PublishResponse, error) {
	res, err := PublishRequest(pn, opts)
	return res, err
}

// TODO: replace result with a pointer
func (pn *PubNub) PublishWithContext(ctx Context,
	opts *PublishOpts) (PublishResponse, error) {

	return PublishRequestWithContext(ctx, pn, opts)
}

func (pn *PubNub) History(opts *HistoryOpts) (*HistoryResponse, error) {
	return HistoryRequest(pn, opts)
}

func (pn *PubNub) HistoryWithContext(ctx Context,
	opts *HistoryOpts) (*HistoryResponse, error) {

	return HistoryRequestWithContext(ctx, pn, opts)
}

// TODO: use builder instead plain arguments
func (pn *PubNub) Subscribe(subOperation *SubscribeOperation) {
	pn.subscriptionManager.adaptSubscribe(subOperation)
}

func (pn *PubNub) AddListener(listener *Listener) {
	pn.subscriptionManager.AddListener(listener)
}

// Set a client for transactional requests
func (pn *PubNub) SetClient(c *http.Client) {
	pn.client = c
}

// Set a client for transactional requests
func (pn *PubNub) GetClient() *http.Client {
	if pn.client == nil {
		pn.client = NewHttpClient(pn.Config.ConnectTimeout,
			pn.Config.NonSubscribeRequestTimeout)
	}

	return pn.client
}

// Set a client for transactional requests
func (pn *PubNub) GetSubscribeClient() *http.Client {
	if pn.subscribeClient == nil {
		pn.subscribeClient = NewHttpClient(pn.Config.ConnectTimeout,
			pn.Config.SubscribeRequestTimeout)
	}

	return pn.subscribeClient
}

func NewPubNub(pnconf *Config) *PubNub {
	publishSequence := make(chan int)

	go runPublishSequenceManager(MaxSequence, publishSequence)

	pn := &PubNub{
		Config:          pnconf,
		publishSequence: publishSequence,
	}

	pn.subscriptionManager = newSubscriptionManager(pn)

	return pn
}

func NewPubNubDemo() *PubNub {
	return &PubNub{
		Config: NewDemoConfig(),
	}
}

func runPublishSequenceManager(maxSequence int, ch chan int) {
	for i := 1; ; i++ {
		if i == maxSequence {
			i = 1
		}

		ch <- i
	}
}
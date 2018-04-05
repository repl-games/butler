package option

import (
	"errors"
	"net/http"
	"time"

	"github.com/itchio/httpkit/timeout"
	"github.com/itchio/wharf/state"
)

type EOSSettings struct {
	HTTPClient *http.Client
	Consumer   *state.Consumer
	MaxTries   int
}

var defaultConsumer *state.Consumer

func DefaultSettings() *EOSSettings {
	return &EOSSettings{
		HTTPClient: defaultHTTPClient(),
		Consumer:   defaultConsumer,
		MaxTries:   2,
	}
}

func SetDefaultConsumer(consumer *state.Consumer) {
	defaultConsumer = consumer
}

func defaultHTTPClient() *http.Client {
	client := timeout.NewClient(time.Second*time.Duration(30), time.Second*time.Duration(15))
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if len(via) >= 10 {
			return errors.New("stopped after 10 redirects")
		}

		// forward initial request headers
		// see https://github.com/itchio/itch/issues/965
		ireq := via[0]
		for key, values := range ireq.Header {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}

		return nil
	}
	return client
}

//////////////////////////////////////

type Option interface {
	Apply(*EOSSettings)
}

//

type httpClientOption struct {
	client *http.Client
}

func (o *httpClientOption) Apply(settings *EOSSettings) {
	settings.HTTPClient = o.client
}

func WithHTTPClient(client *http.Client) Option {
	return &httpClientOption{client}
}

//

type consumerOption struct {
	consumer *state.Consumer
}

func (o *consumerOption) Apply(settings *EOSSettings) {
	settings.Consumer = o.consumer
}
func WithConsumer(consumer *state.Consumer) Option {
	return &consumerOption{consumer}
}

//

type maxTriesOption struct {
	maxTries int
}

func (o *maxTriesOption) Apply(settings *EOSSettings) {
	settings.MaxTries = o.maxTries
}

func WithMaxTries(maxTries int) Option {
	return &maxTriesOption{maxTries}
}

package publisher

import (
	"bytes"
	"encoding/json"
	"fmt"

	"net/http"
	"time"

	"github.co.za/PandaxZA/hayvn/logs"
	"github.co.za/PandaxZA/hayvn/models"
)

// The publisher interface has a send method, allowing us to use multiple publishing interfaces in the future, and not bound to just REST.
type Publisher interface {
	Send(models.AggregatedmessagesBody)
}

type RestPublisher struct {
	baseURL  string
	basePort int
	logger   *logs.Logger
}

func NewRestPublisher(baseURL string, basePort int, logger *logs.Logger) Publisher {

	return &RestPublisher{
		baseURL:  baseURL,
		basePort: basePort,
		logger:   logger,
	}
}

func (p *RestPublisher) Send(message models.AggregatedmessagesBody) {
	// Default HTTP Transport
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}

	json_data, err := json.Marshal(message)
	if err != nil {
		p.logger.Error().Msgf("%v", err)
	}
	url := fmt.Sprintf("%s:%d/aggregated-messages", p.baseURL, p.basePort)

	p.logger.Info().Msgf("Sending to: %s with body: %s", url, string(json_data))
	_, err = client.Post(url, "application/json", bytes.NewBuffer(json_data))
	if err != nil {
		p.logger.Error().Msgf("%v", err)
	}
}

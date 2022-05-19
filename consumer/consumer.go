package consumer

import (
	"time"

	"github.co.za/PandaxZA/hayvn/logs"
	"github.co.za/PandaxZA/hayvn/models"
)

type Consumer interface {
	StoreMessage(msg models.MessageBody)
	FlushMessages() models.AggregatedmessagesBody
}

type Batcher struct {
	memory        map[string][]models.AggregatedmessagesMessages
	timerRunning  bool
	timerChan     chan bool
	logger        *logs.Logger
	rateLimitSecs int
}

func NewBatcher(memory map[string][]models.AggregatedmessagesMessages, timerChan chan bool, logger *logs.Logger, rateLimitSecs int) *Batcher {
	return &Batcher{
		memory:        memory,
		timerRunning:  false,
		timerChan:     timerChan,
		logger:        logger,
		rateLimitSecs: rateLimitSecs,
	}
}

func (b *Batcher) StoreMessage(msg models.MessageBody) {
	b.logger.Info().Msg("StoreMessage")
	if val, ok := b.memory[msg.Destination]; ok {
		b.memory[msg.Destination] = append(val, models.AggregatedmessagesMessages{
			Text:      msg.Text,
			Timestamp: msg.Timestamp,
		})
	} else {
		a := []models.AggregatedmessagesMessages{{
			Text:      msg.Text,
			Timestamp: msg.Timestamp,
		}}
		b.memory[msg.Destination] = a
	}
	go b.StartTimer()
}

func (b *Batcher) FlushMessages() models.AggregatedmessagesBody {
	b.logger.Info().Msg("Flush Messages")
	batches := []models.AggregatedmessagesBatches{}
	for key, value := range b.memory {
		batches = append(batches, models.AggregatedmessagesBatches{
			Destination: key,
			Messages:    value,
		})
	}

	resp := models.AggregatedmessagesBody{
		Batches: batches,
	}

	b.memory = make(map[string][]models.AggregatedmessagesMessages)

	return resp
}

func (b *Batcher) StartTimer() {
	b.logger.Info().Msg("Start batching")
	if b.timerRunning {
		return
	}

	b.timerRunning = true
	duration := b.rateLimitSecs * int(time.Second)
	time.Sleep(time.Duration(duration))
	b.timerChan <- true
	b.timerRunning = false
}

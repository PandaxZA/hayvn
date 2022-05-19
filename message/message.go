package message

import (
	"context"
	"errors"

	"github.co.za/PandaxZA/hayvn/logs"
	"github.co.za/PandaxZA/hayvn/models"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

func Message(logger *logs.Logger, channel chan models.MessageBody) usecase.IOInteractor {
	messageError := status.Wrap(errors.New("error receiving message"), status.Internal)
	u := usecase.NewIOI(new(models.MessageBody), nil, func(ctx context.Context, input, output interface{}) error {
		var (
			in = input.(*models.MessageBody)
		)

		logger.Info().Msgf("Received message for %s", in.Destination)

		go func() {
			channel <- *in
		}()

		return nil

	})

	u.SetExpectedErrors(messageError)
	u.SetDescription("Receive Message")

	return u

}

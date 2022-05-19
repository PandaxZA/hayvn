package main

import (
	"context"
	"fmt"
	"net/http"

	"github.co.za/PandaxZA/hayvn/consumer"
	"github.co.za/PandaxZA/hayvn/logs"
	"github.co.za/PandaxZA/hayvn/message"
	"github.co.za/PandaxZA/hayvn/models"
	"github.co.za/PandaxZA/hayvn/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/swaggest/rest"
	"github.com/swaggest/rest/chirouter"
	"github.com/swaggest/rest/jsonschema"
	"github.com/swaggest/rest/nethttp"
	"github.com/swaggest/rest/openapi"
	"github.com/swaggest/rest/request"
	"github.com/swaggest/rest/response"
	"github.com/swaggest/rest/response/gzip"
	"github.com/swaggest/swgui/v3cdn"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

func main() {
	// Load up config and env vars
	config := utils.MakeEnvironment()
	logger := logs.New(true, config.IS_LOCAL)
	logger.Printf("Env loaded %+v \n", config.ENVIRONMENT_NAME)
	r := chirouter.NewWrapper(chi.NewRouter())
	logger.Printf("Router created")

	// Init API documentation schema.
	apiSchema := &openapi.Collector{}
	apiSchema.Reflector().SpecEns().Info.Title = "HAYVN interview test"
	apiSchema.Reflector().SpecEns().Info.WithDescription("API exposing endpoints for testing")
	apiSchema.Reflector().SpecEns().Info.Version = "v1.0.1"

	validatorFactory := jsonschema.NewFactory(apiSchema, apiSchema)
	decoderFactory := request.NewDecoderFactory()
	decoderFactory.ApplyDefaults = true
	decoderFactory.SetDecoderFunc(rest.ParamInPath, chirouter.PathToURLValues)
	decoderFactory.SetDecoderFunc(rest.ParamIn("wwwForm"), utils.WWWFormUrlEncodedToURLValues)

	// Default CORS options
	corsOptions := cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}
	// Setup middlewares.
	r.Use(
		middleware.Recoverer,                          // Panic recovery.
		nethttp.OpenAPIMiddleware(apiSchema),          // Documentation collector.
		request.DecoderMiddleware(decoderFactory),     // Request decoder setup.
		request.ValidatorMiddleware(validatorFactory), // Request validator setup.
		response.EncoderMiddleware,                    // Response encoder setup.
		gzip.Middleware,                               // Response compression with support for direct gzip pass through.
		cors.Handler(corsOptions),
	)

	//  In-memory storage
	m := make(map[string][]models.AggregatedmessagesMessages)

	//  Channels:
	c := make(chan models.MessageBody)
	p := make(chan models.AggregatedmessagesBody)
	t := make(chan bool, 1)

	// Services
	consumer := consumer.NewBatcher(m, t, logger, config.RATE_LIMIT_SECONDS)

	// Workers
	go messageWorker(c, t, consumer, p)

	// Routes:
	r.Method(http.MethodGet, "/", nethttp.NewHandler(HealthCheck()))
	r.Method(http.MethodGet, "/docs/openapi.json", apiSchema)
	r.Mount("/docs", v3cdn.NewHandler(apiSchema.Reflector().Spec.Info.Title, "/docs/openapi.json", "/docs"))
	r.Method(http.MethodPost, "/message", nethttp.NewHandler(message.Message(logger, c)))

	// Start server.
	logger.Printf("Server started on %s:%d/docs", config.HOST_NAME, config.HOST_PORT)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", config.HOST_PORT), r); err != nil {
		logger.Error().Msgf("error occurred while starting up gateway - %v", err.Error())
	}

}

func HealthCheck() usecase.IOInteractor {
	type healthCheckOutput struct {
		Status string `json:"status"`
		Data   string `json:"data"`
	}

	u := usecase.NewIOI(nil, new(healthCheckOutput), func(ctx context.Context, input, output interface{}) error {
		var (
			out = output.(*healthCheckOutput)
		)
		out.Data = "Server is up and running"
		out.Status = status.OK.String()
		return nil
	})
	u.SetTitle("Health Check")

	return u
}

func messageWorker(c chan models.MessageBody, t chan bool, consumer *consumer.Batcher, publishChan chan models.AggregatedmessagesBody) {
	for {
		select {
		case <-t:
			fmt.Println("Data in timer")
			message := consumer.FlushMessages()
			publishChan <- message

		case x := <-c:
			fmt.Println("Data in messageWorker is: ", x)
			consumer.StoreMessage(x)

		}
	}
}

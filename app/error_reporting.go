package app

import (
	"log"
	"net/http"

	"github.com/getsentry/sentry-go"
)

func initErrorReportingClient() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:         "https://895159d46c5d4ae3974884179de6f7a7@o226064.ingest.sentry.io/6390031",
		Environment: ENV,
		BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			// dont report errors on localhost
			if IsDEV() {
				return nil
			}
			return event
		},
	})
	if err != nil {
		log.Fatalf("unable to init sentry : %s", err)
	}
}

func ReportError(r *http.Request, err error) {
	if ENV != DEV {
		if hub := sentry.GetHubFromContext(r.Context()); hub != nil {
			hub.CaptureException(err)
		}
	}
}

package m4m

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/webhook"
)

func init() {
	cfg := GetConfig()
	stripe.Key = cfg.Stripe.Secret
	stripe.DefaultLeveledLogger = &stripe.LeveledLogger{
		Level: stripe.LevelWarn,
	}
}

// StripeWebHookHandler handles
func StripeWebHookHandler(w http.ResponseWriter, r *http.Request) {
	var event stripe.Event

	// don't verify signature in debug mode
	if os.Getenv("DEBUG") != "" {
		if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
			LogError("cannot unmarshal stripe's web hook: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		// verify event's signature
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			LogError("cannot read request body: %v", err)
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		cfg := GetConfig()
		event, err = webhook.ConstructEvent(body, r.Header.Get("Stripe-Signature"), cfg.Stripe.SigningSecret)
		if err != nil {
			LogError("cannot verify stripe's web hook signature: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	if err := dispatchStripeEvent(&event); err != nil {
		LogError("cannot process event: %v", err)
	}

	w.WriteHeader(http.StatusOK)
}

// dispatchStripeEvent dispatches Stripe event for processing
func dispatchStripeEvent(event *stripe.Event) error {
	switch event.Type {
	case "customer.subscription.created":
		return syncMemberState(event.Data.Object["customer"].(string))
	case "customer.subscription.trial_will_end":
		return syncMemberState(event.Data.Object["customer"].(string))
	case "customer.subscription.updated":
		return syncMemberState(event.Data.Object["customer"].(string))
	}

	return nil
}

package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"ride-sharing/services/api-gateway/grpc_clients"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/env"
	"ride-sharing/shared/messaging"
	"ride-sharing/shared/tracing"
	"time"

	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/webhook"
)

var tripServiceURL = getTripServiceURL()

var tracer = tracing.GetTracer("api-gateway")

func getTripServiceURL() string {
	if url := os.Getenv("TRIP_SERVICE_URL"); url != "" {
		return url
	}
	return "http://trip-service:8083/preview"
}

func handleStripeWebhook(w http.ResponseWriter, r *http.Request, rabbitmq *messaging.RabbitMQ) {
	ctx, span := tracer.Start(r.Context(), "handleStripeWebhook")
	defer span.End()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	webhookKey := env.GetString("STRIPE_WEBHOOK_KEY", "")
	if webhookKey == "" {
		log.Println("STRIPE_WEBHOOK_KEY is not set, skipping webhook processing")
		http.Error(w, "Webhook key not configured", http.StatusInternalServerError)
		return
	}

	event, err := webhook.ConstructEventWithOptions(
		body,
		r.Header.Get("Stripe-Signature"),
		webhookKey,
		webhook.ConstructEventOptions{
			IgnoreAPIVersionMismatch: true,
		},
	)
	if err != nil {
		log.Printf("Failed to verify webhook signature: %v", err)
		http.Error(w, "Invalid webhook signature", http.StatusBadRequest)
		return
	}

	log.Printf("Received strip event: %v", event)
	switch event.Type {
	case "checkout.session.completed":
		var session stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &session)
		if err != nil {
			log.Printf("Failed to parse webhook event data: %v", err)
			http.Error(w, "Failed to parse event data", http.StatusBadRequest)
			return
		}

		payload := messaging.PaymentStatusUpdateData{
			TripID:   session.Metadata["trip_id"],
			UserID:   session.Metadata["user_id"],
			DriverID: session.Metadata["driver_id"],
		}

		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			log.Printf("Failed to marshal payment status update data: %v", err)
			http.Error(w, "Failed to marshal payload", http.StatusInternalServerError)
			return
		}

		message := contracts.AmqpMessage{
			OwnerID: session.Metadata["user_id"],
			Data:    payloadBytes,
		}

		if err := rabbitmq.PublishMessage(
			ctx,
			contracts.PaymentEventSuccess,
			message,
		); err != nil {
			log.Printf("Error publising payment event: %v", err)
			http.Error(w, "Failed to publish payment event", http.StatusInternalServerError)
			return
		}
	}
}

func handleTripPreview(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "handleTripPreview")
	defer span.End()

	if os.Getenv("ENABLE_DELAY") == "true" {
		time.Sleep(2 * time.Second)
	}

	var reqBody previewTripRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		writeJSONError(w, http.StatusBadRequest, "DECODE_ERROR", "failed to parse JSON data")
		return
	}

	defer r.Body.Close()

	if reqBody.UserID == "" {
		writeJSONError(w, http.StatusBadRequest, "VALIDATION_ERROR", "user ID is required")
		return
	}

	// Why we need to create a new client for each connection?
	// Because if a service is down, we don't want to block the whole application
	// so we createa a new client for each connection
	tripService, err := grpc_clients.NewTripServiceClient()
	if err != nil {
		log.Fatal(err)
	}

	defer tripService.Close()

	tripPreview, err := tripService.Client.PreviewTrip(ctx, reqBody.toProto())
	if err != nil {
		log.Printf("Failed to preview a trip: %v", err)
		http.Error(w, "Failed to preview a trip", http.StatusInternalServerError)
		return
	}

	response := contracts.APIResponse{Data: toPreviewTripResponse(tripPreview)}
	writeJSON(w, http.StatusCreated, response)
}

func handleTripStart(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "handleTripStart")
	defer span.End()

	var reqBody createTripRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		writeJSONError(w, http.StatusBadRequest, "DECODE_ERROR", "failed to parse JSON data")
		return
	}

	defer r.Body.Close()
	if reqBody.UserID == "" {
		writeJSONError(w, http.StatusBadRequest, "VALIDATION_ERROR", "user ID is required")
		return
	}

	tripService, err := grpc_clients.NewTripServiceClient()
	if err != nil {
		log.Fatal(err)
	}

	defer tripService.Close()

	trip, err := tripService.Client.CreateTrip(ctx, reqBody.toProto())
	if err != nil {
		log.Printf("Failed to create a trip: %v", err)
		http.Error(w, "Failed to create a trip", http.StatusInternalServerError)
		return
	}

	response := contracts.APIResponse{Data: toCreateTripResponse(trip)}
	writeJSON(w, http.StatusCreated, response)
}

package rabbitmq

import "go.opentelemetry.io/otel/propagation"

// HeaderCarrier adapts amqp.Table to propagation.TextMapCarrier
type HeaderCarrier map[string]interface{}

// Get returns the value associated with the passed key.
func (h HeaderCarrier) Get(key string) string {
	if val, ok := h[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// Set stores the key-value pair.
func (h HeaderCarrier) Set(key string, value string) {
	h[key] = value
}

// Keys lists the keys stored in this carrier.
func (h HeaderCarrier) Keys() []string {
	keys := make([]string, 0, len(h))
	for k := range h {
		keys = append(keys, k)
	}
	return keys
}

// NewHeaderCarrier creates a new HeaderCarrier
func NewHeaderCarrier(headers map[string]interface{}) propagation.TextMapCarrier {
	if headers == nil {
		headers = make(map[string]interface{})
	}
	return HeaderCarrier(headers)
}

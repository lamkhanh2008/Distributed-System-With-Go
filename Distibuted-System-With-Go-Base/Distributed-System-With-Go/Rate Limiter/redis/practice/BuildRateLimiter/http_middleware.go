package redis_rate_limiter

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	rateLimitingTotalRequests = "Rate-Limiting-Total-Requests"
	rateLimitingState         = "Rate-Limiting-State"
	rateLimitingExpiresAt     = "Rate-Limiting-Expires-At"
)

var (
	_            http.Handler = &httpRateLimiterHandler{}
	_            Extractor    = &httpHeaderExtractor{}
	stateStrings              = map[State]string{
		Allow: "Allow",
		Deny:  "Deny",
	}
)

type Extractor interface {
	Extract(r *http.Request) (string, error)
}

type httpHeaderExtractor struct {
	headers []string
}

func (h *httpHeaderExtractor) Extract(r *http.Request) (string, error) {
	values := make([]string, 0, len(h.headers))
	for _, key := range h.headers {
		if value := strings.TrimSpace(r.Header.Get(key)); value == "" {
			return "", fmt.Errorf("the header %v must have a value set", key)
		} else {
			values = append(values, value)
		}
	}

	return strings.Join(values, "-"), nil
}

func NewHttpHeadersExtractor(headers ...string) Extractor {
	return &httpHeaderExtractor{headers: headers}
}

type RateLimiterConfig struct {
	Extractor  Extractor
	Strategy   Strategy
	Expiration time.Duration
	MaxRequest uint64
}

func NewHTTPRateLimiterHandler(originalHandler http.Handler, config *RateLimiterConfig) http.Handler {
	return &httpRateLimiterHandler{
		handler: originalHandler,
		config:  config,
	}
}

type httpRateLimiterHandler struct {
	handler http.Handler
	config  *RateLimiterConfig
}

func (h *httpRateLimiterHandler) writeResponse(writer http.ResponseWriter, status int, msg string, args ...interface{}) {
	writer.Header().Set("Content-Type", "text/plain")
	writer.WriteHeader(status)
	if _, err := writer.Write([]byte(fmt.Sprintf(msg, args...))); err != nil {
		fmt.Printf("failed to write body to http request %v", err)
	}
}

func (h *httpRateLimiterHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	key, err := h.config.Extractor.Extract(request)
	if err != nil {
		h.writeResponse(writer, http.StatusBadRequest, "faied to collect rate limiting")
		return
	}

	result, err := h.config.Strategy.Run(request.Context(), &Request{
		Key:      key,
		Limit:    h.config.MaxRequest,
		Duration: h.config.Expiration,
	})

	if err != nil {
		h.writeResponse(writer, http.StatusInternalServerError, "failed to run rate limiting")
		return
	}

	writer.Header().Set(rateLimitingTotalRequests, strconv.FormatUint(result.TotalRequests, 10))
	writer.Header().Set(rateLimitingState, stateStrings[result.State])
	writer.Header().Set(rateLimitingExpiresAt, result.ExpiresAt.Format(time.RFC3339))
	if result.State == Deny {
		h.writeResponse(writer, http.StatusTooManyRequests, "you have sent too many requests to this service, slow down please")
		return
	}

	h.handler.ServeHTTP(writer, request)
}

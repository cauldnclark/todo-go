package middleware

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/cauldnclark/todo-go/internal/ratelimit"
)

func RateLimitMiddleware(rateLimiter *ratelimit.RateLimiter, maxReq int, window time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := r.Context().Value(UserIdContextKey).(int)
			if !ok {
				// Not authenticated — you may want to rate limit by IP instead
				next.ServeHTTP(w, r)
				return
			}

			allowed, remaining, resetIn, err := rateLimiter.Allow(r.Context(), userID, maxReq, window)
			if err != nil {
				// Fail open — don't block on Redis error
				next.ServeHTTP(w, r)
				return
			}

			// Set headers (RFC 6585)
			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(maxReq))
			w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
			w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(time.Now().Unix()+resetIn, 10))

			if !allowed {
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("Retry-After", strconv.FormatInt(resetIn, 10))
				w.WriteHeader(http.StatusTooManyRequests)
				response := map[string]string{
					"error":       "Rate limit exceeded",
					"retry_after": strconv.FormatInt(resetIn, 10) + " seconds",
				}
				json.NewEncoder(w).Encode(response)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

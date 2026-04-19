package exchange

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func setupMockServer(statusCode int, responseBody []byte, delay time.Duration) *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if delay > 0 {
			time.Sleep(delay)
		}
		w.WriteHeader(statusCode)
		w.Write(responseBody)
	})
	return httptest.NewServer(handler)
}

func TestGetRate(t *testing.T) {
	t.Run("Valid Response Scenario", func(t *testing.T) {
		srv := setupMockServer(http.StatusOK, []byte(`{"base":"USD","target":"EUR","rate":0.92}`), 0)
		defer srv.Close()

		client := NewExchangeService(srv.URL)
		rate, err := client.GetRate("USD", "EUR")

		require.NoError(t, err)
		require.Equal(t, 0.92, rate)
	})

	t.Run("Business Error Format", func(t *testing.T) {
		srv := setupMockServer(http.StatusBadRequest, []byte(`{"error":"invalid currency pair"}`), 0)
		defer srv.Close()

		client := NewExchangeService(srv.URL)
		_, err := client.GetRate("USD", "RUB")

		require.Error(t, err)
		require.Contains(t, err.Error(), "api error: invalid currency pair")
	})

	t.Run("Broken JSON Body", func(t *testing.T) {
		srv := setupMockServer(http.StatusOK, []byte(`{"base":"USD", "target":"EUR"`), 0)
		defer srv.Close()

		client := NewExchangeService(srv.URL)
		_, err := client.GetRate("USD", "EUR")

		require.ErrorContains(t, err, "decode error")
	})

	t.Run("Timeout Simulation", func(t *testing.T) {
		srv := setupMockServer(http.StatusOK, []byte(`{}`), 50*time.Millisecond)
		defer srv.Close()

		client := NewExchangeService(srv.URL)
		client.Client.Timeout = 10 * time.Millisecond 

		_, err := client.GetRate("USD", "EUR")
		require.ErrorContains(t, err, "network error")
	})

	t.Run("Server 500 Panic", func(t *testing.T) {
		srv := setupMockServer(http.StatusInternalServerError, []byte(`{}`), 0)
		defer srv.Close()

		client := NewExchangeService(srv.URL)
		_, err := client.GetRate("USD", "EUR")

		require.EqualError(t, err, "unexpected status: 500")
	})

	t.Run("Absolutely Empty Body", func(t *testing.T) {
		srv := setupMockServer(http.StatusOK, []byte(""), 0)
		defer srv.Close()

		client := NewExchangeService(srv.URL)
		_, err := client.GetRate("USD", "EUR")

		require.ErrorContains(t, err, "decode error")
	})
}
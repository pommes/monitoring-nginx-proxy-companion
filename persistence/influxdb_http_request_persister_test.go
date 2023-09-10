package persistence

import (
	"github.com/influxdata/influxdb1-client/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"nginx-proxy-metrics/config"
	"os"
	"testing"
	"time"
)

type MockClient struct {
	mock.Mock
}

func (m *MockClient) Ping(timeout time.Duration) (time.Duration, string, error) {
	panic("not implemented in mock")
}

func (m *MockClient) Write(bp client.BatchPoints) error {
	panic("not implemented in mock")
}

func (m *MockClient) QueryAsChunk(q client.Query) (*client.ChunkedResponse, error) {
	panic("not implemented in mock")
}

func (m *MockClient) Close() error {
	panic("not implemented in mock")
}

// Hier sollten Sie alle Methoden mocken, die von dem InfluxDB-Client aufgerufen werden könnten
func (m *MockClient) Query(q client.Query) (*client.Response, error) {
	args := m.Called(q)
	return args.Get(0).(*client.Response), args.Error(1)
}

func TestMain(m *testing.M) {
	os.Setenv(config.EnvProxyContainerName, "-")
	os.Setenv(config.EnvInfluxUrl, "https://xxx")
	os.Setenv(config.EnvInfluxDbName, "-")
	os.Setenv(config.EnvInfluxDbRetentionDuration, "1d")
	os.Setenv(config.EnvInfluxDbTagInstance, "-")
	os.Setenv(config.EnvInfluxDbTagSourceIpsLocal, "-")
	// ... setze andere Variablen hier ...

	config.LoadConfig()

	code := m.Run() // führe alle Tests aus

	os.Exit(code)
}

func TestSetup(t *testing.T) {

	// Erstellen Sie einen neuen MockClient
	mockClient := new(MockClient)

	// Mock-Erwartungen definieren
	mockClient.On("Query", mock.Anything).Return(&client.Response{}, nil)

	// Erstellen Sie eine Instanz des zu testenden Persisters
	persister := InfluxdbHttpRequestPersister{
		influxClient: mockClient,
	}

	// Setup aufrufen
	persister.NewInfluxdbHttpRequestPersister(mockClient)

	// Überprüfen Sie die Erwartungen
	mockClient.AssertExpectations(t)
}

func TestGetSourceIpArea(t *testing.T) {
	tests := []struct {
		name          string
		ip            string
		expectedArea  string
		localSourceIP string
	}{
		{"Test Local IP One Entry", "10.0.0.1", "local", "10."},
		{"Test Local IP Multiple Entries", "127.0.0.1", "local", "10., 192., 127."},
		{"Test Public IP One Entry", "192.168.1.1", "public", "10."},
		{"Test Public IP Multiple Entries", "192.168.1.1", "public", "10., 127., fe80::"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.InfluxDbTagSourceIpsLocal = tt.localSourceIP
			area := getSourceIpArea(tt.ip)
			assert.Equal(t, tt.expectedArea, area)
		})
	}
}

// Sie können ähnliche Tests für andere Funktionen hinzufügen

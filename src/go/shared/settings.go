package shared

import "os"

type ServiceSettings struct {
	ServiceName string
	Host        string
	Port        string
}

func LoadSettings(serviceName string) ServiceSettings {
	host := envOr("SERVICE_HOST", "0.0.0.0")
	port := envOr("SERVICE_PORT", "8000")
	return ServiceSettings{
		ServiceName: serviceName,
		Host:        host,
		Port:        port,
	}
}

func envOr(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

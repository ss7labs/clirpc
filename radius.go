package clirpc

import (
	"context"
	"fmt"

	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
)

// RadiusConfig holds configuration for RADIUS client
type RadiusConfig struct {
	Host     string
	Port     int
	Secret   string
	Username string
}

// SendRadiusDisconnect sends a RADIUS Disconnect-Request message
// Format: "User-Name=0:0:130" to server:port with given secret
func SendRadiusDisconnect(ctx context.Context, config RadiusConfig) error {
	// Create a new RADIUS packet
	packet := radius.New(radius.CodeDisconnectRequest, []byte(config.Secret))

	// Add User-Name attribute
	if err := rfc2865.UserName_SetString(packet, config.Username); err != nil {
		return fmt.Errorf("failed to set User-Name attribute: %w", err)
	}

	// Create the client
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	client := radius.Client{
		Retry: 3,
	}

	// Send the packet
	response, err := client.Exchange(ctx, packet, addr)
	if err != nil {
		return fmt.Errorf("failed to send RADIUS packet: %w", err)
	}

	// Check the response
	if response.Code != radius.CodeDisconnectACK {
		return fmt.Errorf("unexpected response code: %v", response.Code)
	}

	return nil
}

// NewRadiusConfig creates a new RadiusConfig with default values
func NewRadiusConfig(host string, port int, secret string) RadiusConfig {
	return RadiusConfig{
		Host:     host,
		Port:     port,
		Secret:   secret,
		Username: "0:0:130", // Default username as per the radclient example
	}
}

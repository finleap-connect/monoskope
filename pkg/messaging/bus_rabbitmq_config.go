package messaging

import "time"

const (
	DefaultExchangeName   = "m8_events"      // name of the monoskope exchange
	DefaultReconnectDelay = 10 * time.Second // When reconnecting to the server after connection failure
	DefaultReInitDelay    = 5 * time.Second  // When setting up the channel after a channel exception
	DefaultResendDelay    = 3 * time.Second  // When resending messages the server didn't confirm
	DefaultMaxResends     = 10               // How many times resending messages the server didn't confirm
)

type rabbitEventBusConfig struct {
	Name             string        // Name of the client, required
	Url              string        // Connection string, required
	RoutingKeyPrefix string        // Prefix for routing of messages
	ExchangeName     string        // Name of the exchange to initialize/use
	ReconnectDelay   time.Duration // When reconnecting to the server after connection failure
	ResendDelay      time.Duration // When resending messages the server didn't confirm
	MaxResends       int           // How many times resending messages the server didn't confirm
	ReInitDelay      time.Duration // When setting up the channel after a channel exception
}

// NewRabbitEventBusConfig creates a new RabbitEventBusConfig with defaults.
func NewRabbitEventBusConfig(name, url string) *rabbitEventBusConfig {
	return &rabbitEventBusConfig{
		Name:             name,
		Url:              url,
		RoutingKeyPrefix: "m8",
		ExchangeName:     DefaultExchangeName,
		ReconnectDelay:   DefaultReconnectDelay,
		ResendDelay:      DefaultResendDelay,
		MaxResends:       DefaultMaxResends,
		ReInitDelay:      DefaultReInitDelay,
	}
}

// Validate validates that all required fields have been provided
func (conf *rabbitEventBusConfig) Validate() error {
	if conf.Name == "" {
		return ErrConfigNameRequired
	}
	if conf.Url == "" {
		return ErrConfigNameRequired
	}
	return nil
}

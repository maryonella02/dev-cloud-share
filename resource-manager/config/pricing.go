package config

type PricingConfig struct {
	ResourceType string  `json:"resource_type"`
	PricePerCore float64 `json:"price_per_core"`
	PricePerMB   float64 `json:"price_per_mb"`
	PricePerGB   float64 `json:"price_per_gb"`
}

var PricingModel = []PricingConfig{
	{
		ResourceType: "CPU",
		PricePerCore: 0.10,
		PricePerMB:   0.01,
		PricePerGB:   0.02,
	},
	// Add other resource types and their pricing if needed
}

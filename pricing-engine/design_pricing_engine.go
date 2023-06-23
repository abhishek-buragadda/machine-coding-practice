package pricing_engine

import (
	"fmt"
	"math"
	"time"
)

type PricingEngine struct {
	basePrice     float64
	pricePerMile  float64
	pricePerMinute float64
}

type Ride struct {
	distance   float64
	duration   time.Duration
	startTime  time.Time
}

func NewPricingEngine(basePrice, pricePerMile, pricePerMinute float64) *PricingEngine {
	return &PricingEngine{
		basePrice:     basePrice,
		pricePerMile:  pricePerMile,
		pricePerMinute: pricePerMinute,
	}
}

func (pe *PricingEngine) CalculatePrice(ride *Ride, demandMultiplier float64) float64 {
	basePrice := pe.basePrice * demandMultiplier
	distancePrice := ride.distance * pe.pricePerMile
	durationPrice := ride.duration.Minutes() * pe.pricePerMinute

	price := math.Round((basePrice + distancePrice + durationPrice) * 100) / 100 // round to 2 decimal places
	return price
}

func main() {
	pricingEngine := NewPricingEngine(2.0, 1.5, 0.3)

	ride1 := &Ride{
		distance: 10.0,
		duration: time.Minute * 30,
		startTime: time.Now(),
	}

	ride2 := &Ride{
		distance: 20.0,
		duration: time.Minute * 45,
		startTime: time.Now().Add(-time.Hour),
	}

	// Calculate prices for the rides at different demand multipliers
	price1 := pricingEngine.CalculatePrice(ride1, 1.0)
	price2 := pricingEngine.CalculatePrice(ride2, 1.5)

	fmt.Printf("Ride 1 price: $%.2f\n", price1)
	fmt.Printf("Ride 2 price: $%.2f\n", price2)
}

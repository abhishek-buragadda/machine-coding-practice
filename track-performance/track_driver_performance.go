package track_performance

import (
	"fmt"
	"sync"
)

type Driver struct {
	id         int
	numRides   int
	totalEarnings float64
	rating     float64
}

type DriverStats struct {
	sync.Mutex
	drivers map[int]*Driver
}

func NewDriverStats() *DriverStats {
	return &DriverStats{
		drivers: make(map[int]*Driver),
	}
}

func (s *DriverStats) AddDriver(id int) {
	s.Lock()
	defer s.Unlock()

	if _, ok := s.drivers[id]; !ok {
		s.drivers[id] = &Driver{id: id}
	}
}

func (s *DriverStats) AddRide(id int, earnings float64, rating float64) {
	s.Lock()
	defer s.Unlock()

	if driver, ok := s.drivers[id]; ok {
		driver.numRides++
		driver.totalEarnings += earnings
		driver.rating = (driver.rating*float64(driver.numRides-1) + rating) / float64(driver.numRides)
	}
}

func (s *DriverStats) GetDriverStats(id int) *Driver {
	s.Lock()
	defer s.Unlock()

	if driver, ok := s.drivers[id]; ok {
		return driver
	}
	return nil
}

func main() {
	driverStats := NewDriverStats()

	// add drivers
	driverStats.AddDriver(1)
	driverStats.AddDriver(2)
	driverStats.AddDriver(3)

	// add rides
	driverStats.AddRide(1, 10.0, 4.5)
	driverStats.AddRide(1, 12.0, 4.0)
	driverStats.AddRide(2, 8.0, 3.0)
	driverStats.AddRide(3, 15.0, 4.8)

	// get driver stats
	driver1 := driverStats.GetDriverStats(1)
	driver2 := driverStats.GetDriverStats(2)
	driver3 := driverStats.GetDriverStats(3)

	fmt.Printf("Driver %d stats: numRides=%d, totalEarnings=%.2f, rating=%.2f\n", driver1.id, driver1.numRides, driver1.totalEarnings, driver1.rating)
	fmt.Printf("Driver %d stats: numRides=%d, totalEarnings=%.2f, rating=%.2f\n", driver2.id, driver2.numRides, driver2.totalEarnings, driver2.rating)
	fmt.Printf("Driver %d stats: numRides=%d, totalEarnings=%.2f, rating=%.2f\n", driver3.id, driver3.numRides, driver3.totalEarnings, driver3.rating)
}

package ride_sharing

import (
	"errors"
	"strconv"
	"time"
)

//Rider struct
type Rider struct {
	ID     string
	Name   string
	Mobile string
}

//Driver struct
type Driver struct {
	ID     string
	Name   string
	Mobile string
	CarNo  string
}

//Ride struct
type Ride struct {
	ID            string
	RiderID       string
	DriverID      string
	FromLocation  string
	ToLocation    string
	StartTime     time.Time
	EndTime       time.Time
	Fare          float64
	RideCompleted bool
}

//RideService interface
type RideService interface {
	RequestRide(rider Rider, fromLocation string, toLocation string) (Ride, error)
	GetRideDetails(rideID string) (Ride, error)
	CompleteRide(rideID string) error
}

//RideServiceImpl struct
type RideServiceImpl struct {
	rides      map[string]Ride
	drivers    map[string]Driver
	nextRideID int
}

//NewRideService function to create a new ride service
func NewRideService() RideService {
	return &RideServiceImpl{
		rides:      make(map[string]Ride),
		drivers:    make(map[string]Driver),
		nextRideID: 1,
	}
}

//RequestRide function to request a ride
func (rs *RideServiceImpl) RequestRide(rider Rider, fromLocation string, toLocation string) (Ride, error) {
	driverID := ""
	for _, driver := range rs.drivers {
		if driver.CarNo != "" {
			driverID = driver.ID
			break
		}
	}

	if driverID == "" {
		return Ride{}, errors.New("no driver available")
	}

	rideID := strconv.Itoa(rs.nextRideID)
	rs.nextRideID++

	startTime := time.Now()

	ride := Ride{
		ID:            rideID,
		RiderID:       rider.ID,
		DriverID:      driverID,
		FromLocation:  fromLocation,
		ToLocation:    toLocation,
		StartTime:     startTime,
		EndTime:       time.Time{},
		Fare:          0,
		RideCompleted: false,
	}

	rs.rides[rideID] = ride

	return ride, nil
}

//GetRideDetails function to get the details of a ride
func (rs *RideServiceImpl) GetRideDetails(rideID string) (Ride, error) {
	ride, ok := rs.rides[rideID]
	if !ok {
		return Ride{}, errors.New("ride not found")
	}

	return ride, nil
}

//CompleteRide function to complete a ride
func (rs *RideServiceImpl) CompleteRide(rideID string) error {
	ride, ok := rs.rides[rideID]
	if !ok {
		return errors.New("ride not found")
	}

	if ride.RideCompleted {
		return errors.New("ride already completed")
	}

	endTime := time.Now()

	duration := endTime.Sub(ride.StartTime)

	fare := duration.Seconds() / 60.0 * 0.5

	ride.EndTime = endTime
	ride.Fare = fare
	ride.RideCompleted = true

	rs.rides[rideID] = ride

	return nil
}

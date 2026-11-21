package util

import (
	"math"
)

const earthRadiusMeters = 6371000.0 // Earth's mean radius in meters

// degreesToRadians converts degrees to radians.
func degreesToRadians(deg float64) float64 {
	return deg * math.Pi / 180.0
}

// CalculateDistanceMeters calculates the great-circle distance between two GPS coordinates
// using the Haversine formula, returning distance in meters.
func CalculateDistanceMeters(lat1, lon1, lat2, lon2 float64) float64 {
	dLat := degreesToRadians(lat2 - lat1)
	dLon := degreesToRadians(lon2 - lon1)

	rLat1 := degreesToRadians(lat1)
	rLat2 := degreesToRadians(lat2)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(rLat1)*math.Cos(rLat2)*math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusMeters * c
}

// IsWithinRadius checks if a user coordinate is within the specified radius (in meters)
// from the outlet reference coordinate. Returns true if within radius, and the exact distance.
func IsWithinRadius(userLat, userLon, outletLat, outletLon, maxRadiusMeters float64) (bool, float64) {
	dist := CalculateDistanceMeters(userLat, userLon, outletLat, outletLon)
	return dist <= maxRadiusMeters, math.Round(dist*10) / 10
}

package util

import (
	"math"
	"testing"
)

func TestCalculateDistanceMeters(t *testing.T) {
	// Point A: FAN Laundry Kalibaru (-6.113425, 106.914230)
	outletLat := -6.113425
	outletLon := 106.914230

	// Point B: Same location
	distSame := CalculateDistanceMeters(outletLat, outletLon, outletLat, outletLon)
	if distSame != 0 {
		t.Errorf("expected 0 meters for same point, got %f", distSame)
	}

	// Point C: ~15 meters away
	userLat := -6.113520
	userLon := 106.914230
	dist := CalculateDistanceMeters(outletLat, outletLon, userLat, userLon)
	if dist < 8 || dist > 15 {
		t.Errorf("expected ~10-12m distance, got %f", dist)
	}

	// Test IsWithinRadius
	within, d := IsWithinRadius(userLat, userLon, outletLat, outletLon, 25.0)
	if !within {
		t.Errorf("expected within 25m radius, got within=false, dist=%f", d)
	}

	withinFar, dFar := IsWithinRadius(-6.120000, 106.914230, outletLat, outletLon, 25.0)
	if withinFar {
		t.Errorf("expected false for point 700m away, got withinFar=true, dist=%f", dFar)
	}
	if math.Abs(dFar-730) > 50 {
		t.Errorf("expected ~730m, got %f", dFar)
	}
}

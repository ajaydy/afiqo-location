package maps

import (
	"afiqo-location/maps/models"
	"context"
)

func GetDistanceBetweenTwoLocations(ctx context.Context, origin, destination string) (models.DistanceMatrix, error) {

	extraURL := "/maps/api/distancematrix/json"

	var distance models.DistanceMatrix

	param := map[string]string{
		"origins":      origin,
		"destinations": destination,
	}

	err := Get(ctx, extraURL, param, &distance)
	if err != nil {
		return models.DistanceMatrix{}, nil
	}

	return distance, nil
}

package maps

import (
	"context"
	uuid "github.com/satori/go.uuid"
)

func GetDistanceBetweenTwoLocations(ctx context.Context, origin, destination string, ID uuid.UUID) (DistanceMatrix, error) {

	extraURL := "/maps/api/distancematrix/json"

	var distance DistanceMatrix

	distance.ID = ID

	param := map[string]string{
		"origins":      origin,
		"destinations": destination,
	}

	err := Get(ctx, extraURL, param, &distance)
	if err != nil {
		return DistanceMatrix{}, nil
	}

	return distance, nil
}

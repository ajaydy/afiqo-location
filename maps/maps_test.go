package maps

import (
	"afiqo-location/maps/models"
	"context"
	"fmt"
	"log"
	"testing"
)

func TestGetDistanceBetweenTwoLocations(t *testing.T) {

	ctx := context.Background()

	maps := Maps{
		URL:    "https://maps.googleapis.com",
		ApiKey: "",
	}

	Init(maps)

	extraURL := "/maps/api/distancematrix/json"

	var distance models.DistanceMatrix

	param := map[string]string{
		"origins":      "229, Jalan 2, Kampung Subang Baru, 40150 Shah Alam, Selangor",
		"destinations": "1, Jalan SS 7/26a, Ss 7, 47301 Petaling Jaya, Selangor",
	}

	err := Get(ctx, extraURL, param, &distance)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(distance)

}

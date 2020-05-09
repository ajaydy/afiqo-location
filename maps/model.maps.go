package maps

import uuid "github.com/satori/go.uuid"

type DistanceMatrix struct {
	ID                   uuid.UUID `json:"id"`
	DestinationAddresses []string  `json:"destination_addresses"`
	OriginAddresses      []string  `json:"origin_addresses"`
	Rows                 []struct {
		Elements []struct {
			Distance struct {
				Text  string `json:"text"`
				Value int    `json:"value"`
			} `json:"distance"`
			Duration struct {
				Text  string `json:"text"`
				Value int    `json:"value"`
			} `json:"duration"`
			Status string `json:"status"`
		} `json:"elements"`
	} `json:"rows"`
	Status string `json:"status"`
}

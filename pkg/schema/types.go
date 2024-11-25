package schema

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/jamieyoung5/kwikmedical-eventstream/pb"
)

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (loc Location) Value() (driver.Value, error) {
	bytes, err := json.Marshal(loc)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal location to JSON: %w", err)
	}
	return string(bytes), nil
}

func (loc *Location) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal location: expected []byte, got %T", value)
	}

	if err := json.Unmarshal(bytes, loc); err != nil {
		return fmt.Errorf("failed to unmarshal location: %w", err)
	}

	return nil
}

func LocationFromPb(loc *pb.Location) Location {
	return Location{
		Latitude:  loc.Latitude,
		Longitude: loc.Longitude,
	}
}

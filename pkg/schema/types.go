package schema

import (
	"database/sql/driver"
	"fmt"
	"github.com/jamieyoung5/kwikmedical-eventstream/pb"
)

type Point struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// Implement the driver.Valuer interface for converting Point to a PostgreSQL-compatible format
func (p Point) Value() (driver.Value, error) {
	// Convert the Point to the PostgreSQL POINT syntax
	return fmt.Sprintf("POINT(%f %f)", p.Latitude, p.Longitude), nil
}

// Implement the sql.Scanner interface for converting a PostgreSQL POINT to a Point struct
func (p *Point) Scan(value interface{}) error {
	// Convert the database value to a Point
	b, ok := value.(string)
	if !ok {
		return fmt.Errorf("Point: unable to scan type %T into Point", value)
	}

	// Parse the POINT format (e.g., "POINT(1.23 4.56)")
	var lat, lon float64
	_, err := fmt.Sscanf(b, "POINT(%f %f)", &lat, &lon)
	if err != nil {
		return fmt.Errorf("Point: unable to parse value '%s': %v", b, err)
	}

	p.Latitude = lat
	p.Longitude = lon
	return nil
}

func pointFromPb(loc *pb.Location) Point {
	return Point{
		Latitude:  loc.Latitude,
		Longitude: loc.Longitude,
	}
}

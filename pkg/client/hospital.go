package client

import (
	"fmt"
	"github.com/jamieyoung5/kwikmedical-db-lib/pkg/schema"
	pbSchema "github.com/jamieyoung5/kwikmedical-eventstream/pb"
)

func (db *KwikMedicalDBClient) GetNearestHospital(location *pbSchema.Location) (*schema.RegionalHospital, error) {
	point := schema.LocationFromPb(location)

	var nearestHospital schema.RegionalHospital
	err := db.gormDb.Raw(`
	SELECT * FROM regional_hospitals
	ORDER BY ST_DistanceSphere(ST_MakePoint((location->>'longitude')::float, (location->>'latitude')::float), ST_MakePoint(?, ?)) ASC
	LIMIT 1
`, point.Longitude, point.Latitude).Scan(&nearestHospital).Error

	if err != nil {
		return nil, fmt.Errorf("error fetching nearest hospital: %w", err)
	}

	return &nearestHospital, nil
}

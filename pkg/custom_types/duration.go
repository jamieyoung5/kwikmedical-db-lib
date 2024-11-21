package custom_types

import (
	"database/sql/driver"
	"fmt"
	"google.golang.org/protobuf/types/known/durationpb"
	"strings"
	"time"
)

type Duration struct {
	time.Duration
}

func (d *Duration) ToProto() *durationpb.Duration {
	return durationpb.New(d.Duration)
}

func (d *Duration) FromProto(pb *durationpb.Duration) {
	d.Duration = pb.AsDuration()
}

func (d *Duration) Scan(value interface{}) error {
	strValue, ok := value.(string)
	if !ok {
		bytes, ok := value.([]uint8)
		if !ok {
			return fmt.Errorf("invalid type for Duration: %T", value)
		}
		strValue = string(bytes)
	}

	parsed, err := parsePostgresInterval(strValue)
	if err != nil {
		return err
	}

	d.Duration = parsed
	return nil
}

func (d *Duration) Value() (driver.Value, error) {
	return d.String(), nil
}

func (d *Duration) String() string {
	return fmt.Sprintf("%d:%02d:%02d", int64(d.Hours()), int64(d.Minutes())%60, int64(d.Seconds())%60)
}

func parsePostgresInterval(interval string) (time.Duration, error) {
	parts := strings.Split(interval, ":")
	if len(parts) != 3 {
		return 0, fmt.Errorf("invalid INTERVAL format: %s", interval)
	}

	hours := parts[0]
	minutes := parts[1]
	seconds := parts[2]

	return time.ParseDuration(fmt.Sprintf("%sh%sm%ss", hours, minutes, seconds))
}

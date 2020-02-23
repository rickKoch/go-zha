package rtmapi

import (
	"strconv"
	"strings"
	"time"
)

// TimeStamp represents slack time representation.
type TimeStamp struct {
	Time          time.Time
	OriginalValue string
}

// UnmarshalText parses a given slack timestamp to time.Time
func (ts *TimeStamp) UnmarshalText(b []byte) error {
	str := string(b)
	ts.OriginalValue = str

	i, err := strconv.ParseInt(strings.Split(str, ".")[0], 10, 64)
	if err != nil {
		return err
	}
	ts.Time = time.Unix(i, 0)

	return nil
}

// String returuns the original timestamp
func (ts *TimeStamp) String() string {
	return ts.OriginalValue
}

// MarshalText returns the stringified value of slack timestamp.
func (ts *TimeStamp) MarshalText() ([]byte, error) {
	return []byte(ts.String()), nil
}
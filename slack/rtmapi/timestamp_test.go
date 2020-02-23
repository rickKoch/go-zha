package rtmapi

import (
	"bytes"
	"strconv"
	"testing"
	"time"
)

func TestUnmarshalTimeStampText(t *testing.T) {
	ts := &TimeStamp{}
	if err := ts.UnmarshalText([]byte("1355517536.000001")); err != nil {
		t.Errorf("error on unmarshal slack timestamp %s", err.Error())
		return
	}

	expectedTime := time.Unix(1355517536, 0)

	if !ts.Time.Equal(expectedTime) {
		t.Errorf("unmarshaled time is wrong %s. expected %s", ts.Time.String(), expectedTime.String())
	}
}

func TestUnmarshalInvalidTimeStampText(t *testing.T) {
	invalidInput := "testing"
	ts := &TimeStamp{}
	if err := ts.UnmarshalText([]byte(invalidInput)); err == nil {
		t.Errorf("error should appear %s", invalidInput)
	}
}

func TestMarshalTimeStampText(t *testing.T) {
	now := time.Now()
	slackTS := strconv.Itoa(now.Second()) + ".123"
	ts := &TimeStamp{Time: now, OriginalValue: slackTS}
	b, e := ts.MarshalText()
	if e != nil {
		t.Errorf("error on marshal slack timestamp. %s.", e.Error())
		return
	}

	if !bytes.Equal(b, []byte(slackTS)) {
		t.Errorf("marshaled value is wrong %s. expected %s", string(b), slackTS)
	}

}

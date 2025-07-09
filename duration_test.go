package durago

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

const (
	timeDay   = time.Hour * 24
	timeWeek  = timeDay * 7
	timeYear  = timeDay * 365
	timeMonth = timeYear / 12
)

func TestParseDuration(t *testing.T) {
	cases := []struct {
		Duration    string
		Expected    time.Duration
		ExpectedErr string
	}{
		{
			Duration: "PT1H",
			Expected: time.Hour,
		},
		{
			Duration: "-PT1H",
			Expected: -time.Hour,
		},
		{
			Duration: "+P00DT01H30M00S",
			Expected: time.Hour + time.Minute*30,
		},
		{
			Duration: "P3Y6M4DT12H30M5S",
			Expected: timeYear*3 + timeMonth*6 + timeDay*4 + time.Hour*12 + time.Minute*30 + time.Second*5,
		},
		{
			Duration: "P2WT4H",
			Expected: timeWeek*2 + time.Hour*4,
		},
		{
			Duration: "P3Y6M2W4DT12H30M5S",
			Expected: timeYear*3 + timeMonth*6 + timeWeek*2 + timeDay*4 + time.Hour*12 + time.Minute*30 + time.Second*5,
		},
		{
			Duration:    "P3Y6M6M2W4DT12H30M5S",
			ExpectedErr: "invalid format: unexpected month designator",
		},
		{
			Duration: "P0Y0M0W00DT00H00M05S",
			Expected: time.Second * 5,
		},
		{
			Duration: "P0Y0M0W0DT0H00M05.5S",
			Expected: time.Second*5 + time.Millisecond*500,
		},
		{
			Duration:    "P6",
			ExpectedErr: "invalid format: missing designator",
		},
		{
			Duration:    "P6Y4",
			ExpectedErr: "invalid format: missing designator",
		},
	}

	for _, c := range cases {
		d, err := ParseDuration(c.Duration)
		if err != nil || c.ExpectedErr != "" {
			if err.Error() != c.ExpectedErr {
				t.Fatalf("expecting error '%s'; got '%s'", c.ExpectedErr, err.Error())
			}
			continue
		}

		if c.Expected != d.GetTimeDuration() {
			t.Fatalf("expected duration %d; got %d", c.Expected, d.GetTimeDuration())
		}
	}
}

func TestFromTimeDuration(t *testing.T) {
	cases := []struct {
		Duration time.Duration
		Expected *Duration
	}{
		{
			Duration: time.Hour,
			Expected: &Duration{
				d:     time.Hour,
				hours: 1,
			},
		},
		{
			Duration: timeYear + timeWeek + timeDay + time.Hour*2,
			Expected: &Duration{
				d:     timeYear + timeWeek + timeDay + time.Hour*2,
				years: 1,
				weeks: 1,
				days:  1,
				hours: 2,
			},
		},
		{
			Duration: -(timeMonth*2 + time.Second),
			Expected: &Duration{
				d:        timeMonth*2 + time.Second,
				negative: true,
				months:   2,
				seconds:  1,
			},
		},
		{
			Duration: time.Second + time.Millisecond*500,
			Expected: &Duration{
				d:       time.Second + time.Millisecond*500,
				seconds: 1.5,
			},
		},
	}

	for _, c := range cases {
		got := FromTimeDuration(c.Duration)
		if !reflect.DeepEqual(got, c.Expected) {
			t.Fatalf("expected duration %v; got %v", c.Expected, got)
		}
	}
}

func TestDuration_GetTimeDuration(t *testing.T) {
	cases := []struct {
		Duration *Duration
		Expected time.Duration
	}{
		{
			Duration: &Duration{
				d: time.Hour + time.Second,
			},
			Expected: time.Hour + time.Second,
		},
		{
			Duration: &Duration{
				d:        time.Hour + time.Second,
				negative: true,
			},
			Expected: -(time.Hour + time.Second),
		},
	}

	for _, c := range cases {
		got := c.Duration.GetTimeDuration()
		if c.Expected != got {
			t.Fatalf("expected duration %d; got %d", c.Expected, got)
		}
	}
}

func TestDuration_String(t *testing.T) {
	cases := []struct {
		Expected string
	}{
		{
			Expected: "P1Y1M1W1DT1H1M1.5S",
		},
		{
			Expected: "-P1DT1H1M",
		},
		{
			Expected: "PT0.001S",
		},
	}

	for _, c := range cases {
		d, err := ParseDuration(c.Expected)
		if err != nil {
			t.Fatalf("expected to parse duration; got %v", err)
		}

		got := d.String()
		if got != c.Expected {
			t.Fatalf("expected duration %s; got %s", c.Expected, got)
		}
	}
}

func TestDuration_MarshalJSON(t *testing.T) {
	d, err := ParseDuration("P3Y6M4DT12H30M5.5S")
	if err != nil {
		t.Fatalf("expected to parse duration; got %v", err)
	}

	jsoned, err := json.Marshal(struct {
		Duration *Duration `json:"duration"`
	}{Duration: d})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	if string(jsoned) != `{"duration":"P3Y6M4DT12H30M5.5S"}` {
		t.Fatalf("expected duration %s; got %s", `{"duration":"P3Y6M4DT12H30M5.5S"}`, string(jsoned))
	}

	jsoned, err = json.Marshal(struct {
		Duration Duration `json:"duration"`
	}{Duration: *d})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	if string(jsoned) != `{"duration":"P3Y6M4DT12H30M5.5S"}` {
		t.Fatalf("expected duration %s; got %s", `{"duration":"P3Y6M4DT12H30M5.5S"}`, string(jsoned))
	}
}

func TestDuration_UnmarshalJSON(t *testing.T) {
	jsoned := `{"duration":"P3Y6M4DT12H30M5.5S"}`

	expected, err := ParseDuration("P3Y6M4DT12H30M5.5S")
	if err != nil {
		t.Fatalf("expected to parse duration; got %v", err)
	}

	var ptrStruct struct {
		Duration *Duration `json:"duration"`
	}

	if err := json.Unmarshal([]byte(jsoned), &ptrStruct); err != nil {
		t.Fatalf("expected to unmarshal; got %v", err)
	}

	if !reflect.DeepEqual(ptrStruct.Duration, expected) {
		t.Fatalf("expected duration %s; got %s", expected, ptrStruct.Duration)
	}

	var plainStruct struct {
		Duration Duration `json:"duration"`
	}

	if err := json.Unmarshal([]byte(jsoned), &plainStruct); err != nil {
		t.Fatalf("expected to unmarshal; got %v", err)
	}

	if !reflect.DeepEqual(plainStruct.Duration, *expected) {
		t.Fatalf("expected duration %s; got %s", expected, &plainStruct.Duration)
	}
}

func BenchmarkParseDuration(b *testing.B) {
	duration := "+P3Y6M1W4DT12H30M5S"

	for b.Loop() {
		ParseDuration(duration)
	}
}

func BenchmarkDuration_String(b *testing.B) {
	duration := "+P99Y11M4W30DT23H59M59S"
	d, _ := ParseDuration(duration)

	for b.Loop() {
		_ = d.String()
	}
}

func BenchmarkFromTimeDuration(b *testing.B) {
	duration := timeYear + timeMonth + timeWeek + timeDay + time.Hour + time.Minute + time.Second + time.Millisecond*500

	for b.Loop() {
		FromTimeDuration(duration)
	}
}

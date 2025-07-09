package durago

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
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
			Duration: "P0Y0M0W00DT00H00M05.5S",
			Expected: time.Second*5 + time.Millisecond*500,
		},
	}

	for _, c := range cases {
		d, err := ParseDuration(c.Duration)
		if err != nil {
			require.EqualError(t, err, c.ExpectedErr)
			continue
		}

		require.Equal(t, c.Expected, d.GetTimeDuration())
	}
}

func BenchmarkParseDuration(b *testing.B) {
	duration := "+P3Y6M1W4DT12H30M5S"

	for b.Loop() {
		ParseDuration(duration)
	}
}

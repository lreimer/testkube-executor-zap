package runner

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshal(t *testing.T) {

	t.Run("Parse Baseline Config YAML", func(t *testing.T) {
		yaml, err := ioutil.ReadFile("../../examples/zap-baseline.yaml")
		if err != nil {
			assert.FailNow(t, "Unable to read ZAP Baseline scan YAML")
		}

		args := Options{}
		err = args.UnmarshalYAML(yaml)

		assert.NoError(t, err)
		assert.Equal(t, "https://www.example.com/", args.Baseline.Target)
	})

	t.Run("Parse API Config YAML", func(t *testing.T) {
		yaml, err := ioutil.ReadFile("../../examples/zap-api.yaml")
		if err != nil {
			assert.FailNow(t, "Unable to read ZAP API scan YAML")
		}

		args := Options{}
		err = args.UnmarshalYAML(yaml)

		assert.NoError(t, err)
		assert.Equal(t, "https://www.example.com/openapi.json", args.API.Target)
	})

	t.Run("Parse Full Config YAML", func(t *testing.T) {
		yaml, err := ioutil.ReadFile("../../examples/zap-full.yaml")
		if err != nil {
			assert.FailNow(t, "Unable to read ZAP Full scan YAML")
		}

		args := Options{}
		err = args.UnmarshalYAML(yaml)

		assert.NoError(t, err)
		assert.Equal(t, "https://www.example.com/", args.Full.Target)
		assert.Equal(t, -1, args.Full.Minutes)
	})
}

func TestArgs(t *testing.T) {
	t.Run("Baseline Scan Args", func(t *testing.T) {
		args := Options{
			Baseline: BaselineOptions{
				Target:  "https://www.example.com/",
				Config:  "examples/zap-baseline.conf",
				Minutes: 3,
				Delay:   -1,
				Debug:   true,
				Level:   "INFO",
				Ajax:    false,
				Short:   true,
			},
		}

		cmd := args.ToBaselineScanArgs("baseline.xml")

		assert.Equal(t, "-t", cmd[0])
		assert.Equal(t, "https://www.example.com/", cmd[1])
		assert.Equal(t, "-c", cmd[2])
		assert.Equal(t, "examples/zap-baseline.conf", cmd[3])
		assert.Equal(t, "-m", cmd[4])
		assert.Equal(t, "3", cmd[5])
		assert.Equal(t, "-d", cmd[6])
		assert.Equal(t, "-l", cmd[7])
		assert.Equal(t, "INFO", cmd[8])
		assert.Equal(t, "-s", cmd[9])
		assert.Equal(t, "-x", cmd[10])
		assert.Equal(t, "baseline.xml", cmd[11])
		assert.Equal(t, "--auto", cmd[12])
	})

	t.Run("Full Scan Args", func(t *testing.T) {
		args := Options{
			Full: FullOptions{
				Target:  "https://www.example.com/",
				Config:  "examples/zap-baseline.conf",
				Minutes: -1,
				Debug:   false,
				Level:   "FAIL",
				Ajax:    true,
				Short:   true,
			},
		}

		cmd := args.ToFullScanArgs("full.xml")

		assert.Equal(t, "-t", cmd[0])
		assert.Equal(t, "https://www.example.com/", cmd[1])
		assert.Equal(t, "-c", cmd[2])
		assert.Equal(t, "examples/zap-baseline.conf", cmd[3])
		assert.Equal(t, "-j", cmd[4])
		assert.Equal(t, "-l", cmd[5])
		assert.Equal(t, "FAIL", cmd[6])
		assert.Equal(t, "-s", cmd[7])
		assert.Equal(t, "-x", cmd[8])
		assert.Equal(t, "full.xml", cmd[9])
	})

	t.Run("API Scan Args", func(t *testing.T) {
		args := Options{
			API: ApiOptions{
				Target:     "https://www.example.com/openapi.json",
				Format:     "openapi",
				Safe:       true,
				Config:     "https://www.example.com/zap-api.conf",
				Debug:      true,
				Short:      false,
				Level:      "PASS",
				User:       "anonymous",
				Delay:      5,
				Time:       60,
				Hostname:   "https://www.example.com",
				ZapOptions: "-config aaa=bbb",
			},
		}

		cmd := args.ToApiScanArgs("report.xml")

		assert.Equal(t, "-t", cmd[0])
		assert.Equal(t, "https://www.example.com/openapi.json", cmd[1])
		assert.Equal(t, "-f", cmd[2])
		assert.Equal(t, "openapi", cmd[3])
		assert.Equal(t, "-u", cmd[4])
		assert.Equal(t, "https://www.example.com/zap-api.conf", cmd[5])
		assert.Equal(t, "-d", cmd[6])
		assert.Equal(t, "-D", cmd[7])
		assert.Equal(t, "5", cmd[8])
		assert.Equal(t, "-l", cmd[9])
		assert.Equal(t, "PASS", cmd[10])
		assert.Equal(t, "-S", cmd[11])
		assert.Equal(t, "-T", cmd[12])
		assert.Equal(t, "60", cmd[13])
		assert.Equal(t, "-U", cmd[14])
		assert.Equal(t, "anonymous", cmd[15])
		assert.Equal(t, "-O", cmd[16])
		assert.Equal(t, "https://www.example.com", cmd[17])
		assert.Equal(t, "-z", cmd[18])
		assert.Equal(t, "-config aaa=bbb", cmd[19])
	})
}

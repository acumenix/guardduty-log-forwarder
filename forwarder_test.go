package main_test

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/google/uuid"
	forwarder "github.com/m-mizutani/guardduty-log-forwarder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testPreference struct {
	s3Bucket string `json:"s3_bucket"`
	s3Prefix string `json:"s3_prefix"`
	s3Region string `json:"s3_region"`
}

func loadTestPref() testPreference {
	raw, err := ioutil.ReadFile("testpref.json")
	if err != nil {
		panic(err)
	}

	cfg := testPreference{}
	err = json.Unmarshal(raw, &cfg)
	if err != nil {
		panic(err)
	}

	return cfg
}

func TestLoadEvent(t *testing.T) {
	jdata, err := ioutil.ReadFile("testdata/gd_sample.json")
	assert.Nil(t, err)
	assert.NotEqual(t, 0, len(jdata))

	var ev forwarder.Event
	err = json.Unmarshal(jdata, &ev)
	assert.Nil(t, err)
	assert.Equal(t, "0", ev.Version)
}

func TestHandler(t *testing.T) {
	pref := loadTestPref()
	jdata, err := ioutil.ReadFile("testdata/gd_sample.json")
	require.Nil(t, err)

	var ev forwarder.Event
	err = json.Unmarshal(jdata, &ev)
	require.Nil(t, err)

	ev.ID = uuid.New().String()

	args := forwarder.Args{
		Event:    ev,
		S3Bucket: pref.s3Bucket,
		S3Prefix: pref.s3Prefix,
		S3Region: pref.s3Region,
	}
	res, err := forwarder.Handler(&args)

	assert.NotEqual(t, 0, len(res))
	assert.Nil(t, err)

}

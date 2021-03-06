package dash

import (
	"bufio"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestStaticParse(t *testing.T) {
	f, err := os.Open("./testdata/static.mpd")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	mpd := &MPD{}
	if err := mpd.Parse(bufio.NewReader(f)); err != nil {
		t.Fatal(err)
	}

	if len(mpd.Periods) != 3 {
		t.Errorf("Expecting 3 Period elements, but got %d", len(mpd.Periods))
	}

	if len(mpd.Periods[0].AdaptationSets) != 2 {
		t.Errorf("Expecting 2 AdaptationSet element on first Period, but got %d", len(mpd.Periods[0].AdaptationSets))
	}
}

func TestDynamicParse(t *testing.T) {
	f, err := os.Open("./testdata/dynamic.mpd")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	mpd := &MPD{}

	if err := mpd.Parse(bufio.NewReader(f)); err != nil {
		t.Fatal(err)
	}

	pt, _ := time.Parse(time.RFC3339Nano, "2013-08-10T22:03:00Z")
	if !reflect.DeepEqual(mpd.PublishTime.Time, pt) {
		t.Errorf("Expected PublishTime %v, but got %v", pt, mpd.PublishTime)
	}
	d, _ := time.ParseDuration("0h10m54.00s")
	if !reflect.DeepEqual(mpd.MediaPresDuration.Duration, d) {
		t.Errorf("Expecting MediaPresDuration to be %v, but got %v", d, mpd.MediaPresDuration)
	}

	if len(mpd.Periods) != 1 {
		t.Errorf("Expecting 1 Period element, but got %d", len(mpd.Periods))
	}

	if len(mpd.Periods[0].AdaptationSets) != 2 {
		t.Errorf("Expecting 2 AdaptationSets, but got %d", len(mpd.Periods[0].AdaptationSets))
	}

	if len(mpd.Periods[0].AdaptationSets[0].Representations) != 3 {
		t.Errorf("Expecting 3 Representations of AdaptationSets[0], but got %d", len(mpd.Periods[0].AdaptationSets[0].Representations))
	}

	if len(mpd.Periods[0].AdaptationSets[1].Representations) != 1 {
		t.Errorf("Expecting 3 Representations of AdaptationSets[1], but got %d", len(mpd.Periods[0].AdaptationSets[1].Representations))
	}

	if len(mpd.Periods[0].AdaptationSets[1].Representations[0].AudioChannelConfig) != 1 {
		t.Errorf("Expecting 1 AudioChannelConfig, but got %d", len(mpd.Periods[0].AdaptationSets[1].Representations[0].AudioChannelConfig))
	}
}

func TestEventMessage(t *testing.T) {
	f, err := os.Open("./testdata/eventmessage.mpd")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	mpd := &MPD{}
	if err := mpd.Parse(bufio.NewReader(f)); err != nil {
		t.Fatal(err)
	}
	event := mpd.Periods[0].EventStream[0]
	if event.SchemeIDURI != "urn:uuid:XYZY" {
		t.Errorf("Expecting SchemeIdURI urn:uuid:XYZY, but got %s", event.SchemeIDURI)
	}
	for i, e := range event.Event {
		if e.Message == "" {
			t.Errorf("Expecting Message to not be empty.")
		}
		if e.Duration != 10000 {
			t.Errorf("Expecting Duration to be 10000, but got %d", e.Duration)
		}
		if e.ID != i {
			t.Errorf("Expecting ID to be %d, but got %d", i, e.ID)
		}
	}

	rep := mpd.Periods[0].AdaptationSets[1].Representations[0]
	for i, ie := range rep.InbandEventStream {
		if ie.SchemeIDURI == "" {
			t.Errorf("Expecting %d SchemeIdURI to be set.", i)
		}
		if ie.Value == "" {
			t.Errorf("Expecting %d Value to be set.", i)
		}
	}
}

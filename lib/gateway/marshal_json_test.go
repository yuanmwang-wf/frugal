package gateway_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/Workiva/frugal/examples/go/gen-go/v1/music"
	"github.com/Workiva/frugal/lib/gateway"
)

func TestJSONMarshal(t *testing.T) {
	var m gateway.JSON
	msg := music.Album{
		ASIN:     "c54d385a-5024-4f3f-86ef-6314546a7e7f",
		Duration: 1200,
		Tracks: []*music.Track{{
			Title:     "Comme des enfants",
			Artist:    "Coeur de pirate",
			Publisher: "Grosse Boîte",
			Composer:  "Béatrice Martin",
			Duration:  169,
			Pro:       music.PerfRightsOrg_ASCAP,
		}},
	}

	buf, err := m.Marshal(&msg)
	if err != nil {
		t.Errorf("m.Marshal(%v) failed with %v; want success", &msg, err)
	}

	var got music.Album
	if err := json.Unmarshal(buf, &got); err != nil {
		t.Errorf("json.Unmarshal(%q, &got) failed with %v; want success", buf, err)
	}
	if want := msg; !reflect.DeepEqual(got, want) {
		t.Errorf("got = %v; want %v", &got, &want)
	}
}

package postcodes

import (
	"fmt"
	"reflect"
	"testing"

	"google.golang.org/genproto/googleapis/type/latlng"
)

type testAPI struct {
	err bool
}

func (api testAPI) Get(postcode string) (LatLong, error) {
	if api.err {
		return LatLong{}, fmt.Errorf("error")
	}
	return LatLong{1, 2}, nil
}

func TestGetLatLong(t *testing.T) {
	type args struct {
		postcode string
	}
	type tc struct {
		name    string
		args    args
		want    *latlng.LatLng
		apiErr  bool
		wantErr bool
	}
	tests := []tc{
		tc{
			name: "should return location data for the postcode",
			args: args{"HP70JE"},
			want: &latlng.LatLng{
				Longitude: -0.615432,
				Latitude:  51.665516,
			},
			apiErr:  false,
			wantErr: false,
		},
		tc{
			name:    "should error if the postcode is invalid",
			args:    args{"invalidpcode"},
			want:    &latlng.LatLng{},
			apiErr:  true,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Postcodes{}
			got, err := p.GetLatLong(tt.args.postcode)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLatLong() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLatLong() = %v, want %v", got, tt.want)
			}
		})
	}
}

package location

import (
	"reflect"
	"testing"
)

func TestGetLatLong(t *testing.T) {
	type args struct {
		postcode string
	}
	type tc struct {
		name    string
		args    args
		want    LatLong
		wantErr bool
	}
	tests := []tc{
		tc{
			name: "should return location data for the postcode",
			args: args{"HP70JE"},
			want: LatLong{
				Longitude: -0.615432,
				Latitude:  51.665516,
			},
			wantErr: false,
		},
		tc{
			name:    "should error if the postcode is invalid",
			args:    args{"invalidpcode"},
			want:    LatLong{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetLatLong(tt.args.postcode)
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

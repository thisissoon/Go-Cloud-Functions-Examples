package location

import (
	"context"
	"fmt"
	"testing"

	"cloud.google.com/go/firestore"
	"github.com/stretchr/testify/assert"
	"google.golang.org/genproto/googleapis/type/latlng"
)

var fixedLatlng = &latlng.LatLng{
	Latitude:  1,
	Longitude: 2,
}

type testGeolocation struct {
	err bool
}

func (gl testGeolocation) GetLatLong(postcode string) (*latlng.LatLng, error) {
	if gl.err {
		return &latlng.LatLng{}, fmt.Errorf("error")
	}
	return fixedLatlng, nil
}

type update struct {
	col   string
	doc   string
	field string
	value *latlng.LatLng
}

type testStore struct {
	err            bool
	expectedUpdate update
	t              *testing.T
}

func (s testStore) UpdateDoc(ctx context.Context, collection string, doc string, update firestore.Update) error {
	if s.err {
		return fmt.Errorf("error")
	}
	assert.Equal(s.t, s.expectedUpdate.col, collection)
	assert.Equal(s.t, s.expectedUpdate.doc, doc)
	assert.Equal(s.t, s.expectedUpdate.field, update.Path)
	assert.Equal(s.t, s.expectedUpdate.value, update.Value)
	return nil
}

func Test_runUpdateLocation(t *testing.T) {
	type StringValue struct {
		StringValue string `json:"stringValue"`
	}
	pcode1 := &Fields{
		Postcode: StringValue{"AB1 2AB"},
	}
	pcode2 := &Fields{
		Postcode: StringValue{"CD3 4CD"},
	}
	type args struct {
		e FirestoreEvent
	}
	tests := []struct {
		name           string
		args           args
		expectedUpdate update
		latlngErr      bool
		storeErr       bool
		wantErr        bool
	}{
		{
			name: "should exit the function early as nothing to do",
			args: args{
				e: FirestoreEvent{
					OldValue: FirestoreValue{},
					Value:    FirestoreValue{},
				},
			},
			storeErr:       false,
			latlngErr:      false,
			wantErr:        false,
			expectedUpdate: update{},
		},
		{
			name: "should do something",
			args: args{
				e: FirestoreEvent{
					OldValue: FirestoreValue{
						Name:   "projects/gcp-id/databases/(default)/documents/people/lNobStdJlbMR23hRZXFL",
						Fields: pcode1,
					},
					Value: FirestoreValue{
						Name:   "projects/gcp-id/databases/(default)/documents/people/lNobStdJlbMR23hRZXFL",
						Fields: pcode2,
					},
				},
			},
			latlngErr: false,
			storeErr:  false,
			expectedUpdate: update{
				col:   "people",
				doc:   "lNobStdJlbMR23hRZXFL",
				field: "location",
				value: fixedLatlng,
			},
			wantErr: false,
		},
		{
			name: "should do nothing if postcodes match",
			args: args{
				e: FirestoreEvent{
					OldValue: FirestoreValue{
						Name:   "projects/gcp-id/databases/(default)/documents/people/lNobStdJlbMR23hRZXFL",
						Fields: pcode1,
					},
					Value: FirestoreValue{
						Name:   "projects/gcp-id/databases/(default)/documents/people/lNobStdJlbMR23hRZXFL",
						Fields: pcode1,
					},
				},
			},
			latlngErr:      false,
			storeErr:       false,
			expectedUpdate: update{},
			wantErr:        false,
		},
		{
			name: "should error out if storage fails",
			args: args{
				e: FirestoreEvent{
					OldValue: FirestoreValue{
						Name:   "projects/gcp-id/databases/(default)/documents/people/lNobStdJlbMR23hRZXFL",
						Fields: pcode1,
					},
					Value: FirestoreValue{
						Name:   "projects/gcp-id/databases/(default)/documents/people/lNobStdJlbMR23hRZXFL",
						Fields: pcode2,
					},
				},
			},
			latlngErr:      false,
			storeErr:       true,
			expectedUpdate: update{},
			wantErr:        true,
		},
		{
			name: "should error when fetching location fails",
			args: args{
				e: FirestoreEvent{
					OldValue: FirestoreValue{
						Name:   "projects/gcp-id/databases/(default)/documents/people/lNobStdJlbMR23hRZXFL",
						Fields: pcode1,
					},
					Value: FirestoreValue{
						Name:   "projects/gcp-id/databases/(default)/documents/people/lNobStdJlbMR23hRZXFL",
						Fields: pcode2,
					},
				},
			},
			latlngErr:      true,
			storeErr:       false,
			expectedUpdate: update{},
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := testStore{
				err:            tt.storeErr,
				expectedUpdate: tt.expectedUpdate,
				t:              t,
			}
			gl := testGeolocation{
				err: tt.latlngErr,
			}
			ctx := context.Background()
			if err := runUpdateLocation(
				ctx,
				tt.args.e,
				gl,
				s,
			); (err != nil) != tt.wantErr {
				t.Errorf("runUpdateLocation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

package addPerson

import (
	"reflect"
	"testing"
	"time"
)

func TestPerson_DobTimestamp(t *testing.T) {
	type fields struct {
		Firstname string
		Lastname  string
		Dob       string
	}
	tests := []struct {
		name    string
		fields  fields
		want    time.Time
		wantErr bool
	}{
		{
			name: "should return correct date",
			fields: fields{
				Firstname: "Maurice",
				Lastname:  "Moss",
				Dob:       "1989-06-24T00:00:00.000Z",
			},
			want:    time.Date(1989, time.June, 24, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Person{
				Firstname: tt.fields.Firstname,
				Lastname:  tt.fields.Lastname,
				Dob:       tt.fields.Dob,
			}
			got, err := p.DobTimestamp()
			if (err != nil) != tt.wantErr {
				t.Errorf("Person.DobTimestamp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Person.DobTimestamp() = %v, want %v", got, tt.want)
			}
		})
	}
}

package location

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github/thisissoon/Go-Cloud-Functions-Examples/functions/events/location/updateLocation/postcodes"
	"github/thisissoon/Go-Cloud-Functions-Examples/functions/events/location/updateLocation/storage"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/genproto/googleapis/type/latlng"
)

// FirestoreEvent is the payload of a Firestore event.
type FirestoreEvent struct {
	OldValue   FirestoreValue `json:"oldValue"`
	Value      FirestoreValue `json:"value"`
	UpdateMask struct {
		FieldPaths []string `json:"fieldPaths"`
	} `json:"updateMask"`
}

// FirestoreValue holds Firestore fields.
type FirestoreValue struct {
	CreateTime time.Time `json:"createTime"`
	// Fields is the data for this value. The type depends on the format of your
	// database. Log an interface{} value and inspect the result to see a JSON
	// representation of your database fields.
	Fields     *Fields   `json:"fields"`
	Name       string    `json:"name"`
	UpdateTime time.Time `json:"updateTime"`
}

// Fields represents values from Firestore. The type definition depends on the format of your database.
type Fields struct {
	CreatedAt struct {
		TimestampValue time.Time `json:"timestampValue"`
	} `json:"createdAt"`
	Dob struct {
		TimestampValue time.Time `json:"timestampValue"`
	} `json:"dob"`
	Firstname struct {
		StringValue string `json:"stringValue"`
	} `json:"firstname"`
	Lastname struct {
		StringValue string `json:"stringValue"`
	} `json:"lastname"`
	Postcode struct {
		StringValue string `json:"stringValue"`
	} `json:"postcode"`
	Location struct {
		GeopointValue latlng.LatLng `json:"geopointValue"`
	} `json:"location"`
}

// GCLOUD_PROJECT is automatically set by the Cloud Functions runtime.
var projectID = os.Getenv("GCLOUD_PROJECT")

// client is a Firestore client, reused between function invocations.
var client *firestore.Client

func init() {
	// Use the application default credentials.
	conf := &firebase.Config{ProjectID: projectID}

	// Use context.Background() because the app/client should persist across
	// invocations.
	ctx := context.Background()

	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalf("firebase.NewApp: %v", err)
	}

	client, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalf("app.Firestore: %v", err)
	}
}

func runUpdateLocation(
	ctx context.Context,
	e FirestoreEvent,
	l postcodes.Location,
	s storage.Storage,
) error {
	// Check if there's any data to process
	// this works because fields is a pointer
	if e.Value.Fields == nil {
		return nil
	}
	// if old values don't exist initialize to empty values
	if e.OldValue.Fields == nil {
		e.OldValue.Fields = &Fields{}
	}
	fullPath := strings.Split(e.Value.Name, "/documents/")[1]
	pathParts := strings.Split(fullPath, "/")
	collection := pathParts[0]
	doc := strings.Join(pathParts[1:], "/")
	// check if postcode has changed
	oldPostcode := e.OldValue.Fields.Postcode.StringValue
	newPostcode := e.Value.Fields.Postcode.StringValue
	if oldPostcode != newPostcode {
		// get location data for postcode
		newLocation, err := l.GetLatLong(newPostcode)
		if err != nil {
			return err
		}
		fmt.Printf("for postcode: %v got lat long: %+v", newPostcode, newLocation)
		// update database with new location could use Set but that would overwrite everything
		update := firestore.Update{
			Path:  "location",
			Value: newLocation,
		}
		err = s.UpdateDoc(ctx, collection, doc, update)
		if err != nil {
			return fmt.Errorf("Update error: %v", err)
		}
	}
	return nil
}

func UpdateLocation(ctx context.Context, e FirestoreEvent) error {
	p := postcodes.Postcodes{}
	s := storage.Store{
		Client: client,
	}
	return runUpdateLocation(ctx, e, p, s)
}

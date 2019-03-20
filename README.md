# Google Cloud Functions using GO

This project has examples for using Go Cloud Function and Firestore

The project is separated into 2 directories
- /functions - for Cloud Function
- //events - Event triggered functions
- //http   - HTTP triggered functions
- /firestore - for non Cloud Function code i.e. can be run on a local machine

Function triggers are:
* HTTP
* Events (firestore)

Also includes a demonstration of interacting with Firestore outside of Cloud Functions

The functions are designed to work together to create a simple people application

In the database people are located in a root collection `/people` and a person is added as a document in the collection e.g. `/people/SYd2omJqJnoahIgqtdgs/`. A person includes the follow firestore data:
| Key       | Type      |
|-----------|-----------|
| createdAt | timestamp |
| dob       | timestamp |
| firstname | string    |
| lastname  | string    |
| postcode  | string    |
| location  | geopoint  |

The http function `AddPerson` takes the following json payload and creates a new person in the DB
```json
{
	"fname": "Maurice",
	"lname": "Moss",
	"dob": "1980-04-20T00:00:00.000Z",
	"postcode": "E1 1AB"
}
```

The event function `UpdateLocation` will respond to write events on the person document and fetch the Latitude/Longitude for the postcode using postcodes api

The firestore application can be run locally and provides an example of querying firestore and extracting a document by ID
you will need to set 2 environment variables
- GCLOUD_PROJECT: (GCP project id)
- GCLOUD_SERVICE_ACCOUNT_PATH: (path to local json service account which has firestore access)
build the program `go build`
and run providing a query id and a string to query lastname for `./firestore iLI2lKSyha8HMGpNCyO7 Moss`




package rideshare_tracker

import (
	"html/template"
	"net/http"

	"appengine"
	"appengine/datastore"
	"time"
	"strconv"
	"encoding/json"
)

type RideshareCar struct {
	state string
	plateNumber	string
	Make	string
	Model string
	color string
}
type RideshareCarView struct {
	whenSeen    time.Time
	seenBy	string
	location string
}

func init() {
	http.HandleFunc("/", showRideShareCars)
	http.HandleFunc("/getCarsJson", getRideShareCars)
	http.HandleFunc("/showRideShareCars", showRideShareCars)
	http.HandleFunc("/addCar", addCar)
}

func getRideShareCars(responseWriter http.ResponseWriter, request *http.Request) {
	context := appengine.NewContext(request)
	// Ancestor queries, as shown here, are strongly consistent with the High
	// Replication Datastore. Queries that span entity groups are eventually
	// consistent. If we omitted the .Ancestor from this query there would be
	// a slight chance that Greeting that had just been written would not
	// show up in a query.
	query := datastore.NewQuery("RideshareCar")
	//	q := datastore.NewQuery("Greeting").Order("-Date").Limit(10)

	cars := make([]RideshareCar, 0, 50)

	if _, err := query.GetAll(context, &cars); err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	carsJson, _ := json.Marshal(cars)
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.Write(carsJson)
}

func showRideShareCars(responseWriter http.ResponseWriter, request *http.Request) {
	context := appengine.NewContext(request)
	// Ancestor queries, as shown here, are strongly consistent with the High
	// Replication Datastore. Queries that span entity groups are eventually
	// consistent. If we omitted the .Ancestor from this query there would be
	// a slight chance that Greeting that had just been written would not
	// show up in a query.
	query := datastore.NewQuery("RideshareCar")
	//	q := datastore.NewQuery("Greeting").Order("-Date").Limit(10)

	cars := make([]RideshareCar, 0, 50)

	if _, err := query.GetAll(context, &cars); err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	context.Infof("car list read. Len = " + strconv.Itoa(len(cars)))
	for i := 0; i < len(cars); i++ {
		context.Infof("Car was read")
	}

	var rideShareCarTemplate, _ = template.ParseFiles("listCars.html");
	if err := rideShareCarTemplate.Execute(responseWriter, cars); err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
	}
}

func addCar(writer http.ResponseWriter, request *http.Request) {
	context := appengine.NewContext(request);
	// [END new_context]
	car := RideshareCar{
		Make: request.FormValue("make"),
		Model: request.FormValue("model"),
	}
	context.Infof("car created. Make = " + car.Make)

	// We set the same parent key on every Greeting entity to ensure each Greeting
	// is in the same entity group. Queries across the single entity group
	// will be consistent. However, the write rate to a single entity group
	// should be limited to ~1/second.
	key := datastore.NewIncompleteKey(context, "RideshareCar", nil)
	_, err := datastore.Put(context, key, &car)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(writer, request, "/showRideShareCars", http.StatusFound)
	// [END if_user]
}




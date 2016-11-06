package rideshare_tracker

import (
	"html/template"
	"net/http"

	"appengine"
	"appengine/datastore"
	"time"
	"encoding/json"
	"strings"
	"strconv"
)

type RideshareCar struct {
	Key string
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
	http.HandleFunc("/cars", carsWebService)
	http.HandleFunc("/showRideShareCars", showRideShareCars)
	http.HandleFunc("/addCar", addCar)
}

func carsWebService(responseWriter http.ResponseWriter, request *http.Request) {
	context := appengine.NewContext(request)

	context.Infof("HTTP method is " + request.Method)
	if (request.Method == http.MethodGet) {
		query := datastore.NewQuery("RideshareCar")

		cars := make([]RideshareCar, 0, 50)

		t := query.Run(context)
		for {
			var car RideshareCar
			key, err := t.Next(&car)
			if err == datastore.Done {
				break // No further entities match the query.
			}
			if err != nil {
				context.Errorf("fetching next car: %v", err)
				break
			}
			// Do something with Person p and Key k
			car.Key = key.String()
			cars = append(cars, car)
		}

		carsJson, _ := json.Marshal(cars)
		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.Write(carsJson)
		return
	}
	if (request.Method == http.MethodPut) {
		request.ParseForm()
		for k, v := range request.Form {
			context.Infof("key[%s] value[%s]\n", k, v)
		}
		keyString := request.Form.Get("Key")
		keyStrings := strings.Split(keyString, ",")
		kind := strings.Trim(keyStrings[0], "/ ")
		stringId := strings.TrimSpace(keyStrings[1])
		intId, _ := strconv.ParseInt(stringId, 10, 64)
		key := datastore.NewKey(context, kind, "", intId, nil)
		context.Infof("kind = %s, stringId = %s, intId = %d, key = %s", kind, stringId, intId, key.String())
		var car RideshareCar

		car.Make = request.Form.Get("Make")
		car.Model = request.Form.Get("Model")
		if _, err  := datastore.Put(context, key, &car); err != nil {
			context.Errorf("Error in put : %v", err)
			return
		}
	}
}

func showRideShareCars(responseWriter http.ResponseWriter, _ *http.Request) {
	var rideShareCarTemplate, _ = template.ParseFiles("listCars.html");
	if err := rideShareCarTemplate.Execute(responseWriter, nil); err != nil {
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




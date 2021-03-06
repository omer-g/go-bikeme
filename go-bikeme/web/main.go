package web

import (
	"fmt"
	"appengine"
	"net/http"
	"go-bikeme/station"
    "go-bikeme/bikeshareservice"
	"appengine/datastore"
    "html/template"
    "time"
)

func init() {
	http.HandleFunc("/update_stations", updateStations)
	http.HandleFunc("/list_stations", listStations)
	http.HandleFunc("/", handler)
}

func listStations(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    q := datastore.NewQuery("Station").Limit(10)
    stations := make([]station.Station, 0, 10)
    if _, err := q.GetAll(c, &stations); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    if err := stationsTemplate.Execute(w, stations); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

var stationsTemplate = template.Must(template.New("listing").Parse(stationsTemplateHTML))

const stationsTemplateHTML = `
<html>
  <body>
    {{range .}}
        <p><b>{{.StationName}}</b>.</p>
    {{end}}
  </body>
</html>
`

func updateStations(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	services := []bikeshareservice.Service{bikeshareservice.NewBicingService(c), bikeshareservice.NewCapitalBikeShareService(c), bikeshareservice.NewTelOFunService(c)}
	for _, service := range services {
		stations, err := service.Stations()
		if err != nil {
			fmt.Fprintf(w, "%T received an error: '%s'\n", service, err.Error())
			continue
		}
		updateTime := time.Now()
		// keys = datastore.Key []datastore.Key
		for _, station := range stations  {
			station.LastUpdate = updateTime
			datastore.Put(c, datastore.NewKey(c, "Station", station.StationId, 0, nil), &station)
		}
		// Future Code	
		//	keys = append(keys, datastore.NewIncompleteKey(c, "Station", station.StationId))
		//}
		//datastore.PutMulti(c, keys, stations)
		fmt.Fprintf(w,"There are %d stations in the %T system!\n", len(stations), service)
	}
	fmt.Fprint(w, "done")
}


func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world!")
}

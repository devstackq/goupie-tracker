package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

//use Api, Marshall, json parser, struct type, web server -Api, client-side html
var (
	temp = template.Must(template.ParseFiles("templates/index.html", "templates/search.html", "templates/locations.html", "templates/artist.html", "templates/members.html", "templates/groupcity.html"))
)

type Locations struct {
	Index []struct {
		ID        int
		Locations []string
	}
}

//give locations and date artist
type Relations struct {
	Index []struct {
		ID             int
		DatesLocations map[string][]string
	}
}

var API struct {
	ID            int
	IDS           int
	GroupWasCity  []int
	IDWC          int
	Artist        []Singers
	LocationsHtml Locations
	RelationHtml  Relations
}

type FindCity struct {
	CityName string
}

type GroupWasCity struct {
	Name string
}
type AlbumDate struct {
	Name string
}
type CreationsDate struct {
	Name string
}

type Singers struct {
	ID           int
	Image        string
	Name         string
	Members      []string
	CreationDate int
	FirstAlbum   string
}

func main() {

	// create each singer, location struct -> then put data
	artists, _ := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	artistsBytes, _ := ioutil.ReadAll(artists.Body)
	artists.Body.Close()
	json.Unmarshal(artistsBytes, &API.Artist)

	locations, _ := http.Get("https://groupietrackers.herokuapp.com/api/locations")
	locationsBytes, _ := ioutil.ReadAll(locations.Body)
	locations.Body.Close()
	json.Unmarshal(locationsBytes, &API.LocationsHtml)

	relations, _ := http.Get("https://groupietrackers.herokuapp.com/api/relation")
	relationsBytes, _ := ioutil.ReadAll(relations.Body)
	relations.Body.Close()
	json.Unmarshal(relationsBytes, &API.RelationHtml)

	//static data, css, js
	static := http.FileServer(http.Dir("public"))
	//secure, not access another files
	http.Handle("/public/", http.StripPrefix("/public/", static))

	http.HandleFunc("/", getAllArtists)
	http.HandleFunc("/artist", getArtist)

	http.HandleFunc("/searchz", MainSearch)

	// http.HandleFunc("/geo", GeoLocation)

	err := http.ListenAndServe(":6969", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

// func isUnique(arr []string) bool {

// 	var newcity []string
// 	fmt.Println(arr)
// 	for i, v := range arr {
// 		if v != newcity[i] {
// 			newcity = append(newcity, v)
// 		}
// 	}
// 	fmt.Println(newcity)
// 	return true
// }

//bakcend send all data - artist input artist
// then client search input value input
func getAllArtists(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		if r.URL.Path != "/" {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		//var arr []string
		// for _, v := range API.LocationsHtml.Index {
		// 	if isUnique(v.Locations) {
		// 		// for _, g := range v.Locations {
		// 		// 	arr = append(arr, string(g))
		// 		// }

		// 	}
		// }

		temp.ExecuteTemplate(w, "index", API)
	}
}

//get define artist by ID, id take - client -> post(name="uid"), -> find Artist by id, return all data -> client
func getArtist(w http.ResponseWriter, r *http.Request) {
	// temp, _! := template.ParseFiles("templates/artist.html")
	if r.URL.Path != "/artist" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}

	if r.Method == "GET" {
		temp.ExecuteTemplate(w, "index", "")

	}
	if r.Method == "POST" {
		ID, _ := strconv.Atoi(r.FormValue("uid"))
		API.ID = ID - 1
		temp.ExecuteTemplate(w, "artist", API)
	}
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	temp = template.Must(template.ParseGlob("templates/*.html"))

	if status == 404 {
		temp.ExecuteTemplate(w, "error.html", nil)
	}
	if status == 500 {
		temp.ExecuteTemplate(w, "error1.html", nil)
	}
}

// 2 input, 2 page. 1 1 artist, 2 more 2 artist and city etc

func MainSearch(w http.ResponseWriter, r *http.Request) {

	key := r.FormValue("input")
	k := r.FormValue("value")
	if key == "members" {
		SearchMembers(w, r, k)
	} else if key == "artist" {
		Search(w, r, k)
	} else if key == "locations" {
		SearchByLocation(w, r, k)
	} else if key == "city" {
		SearchCityGroup(w, r, k)
	} else if key == "albumdate" {
		SearchAlbumDate(w, r, k)
	} else if key == "creation" {
		SearchByCreationDate(w, r, k)
	}

	if key == "" || len(key) > 9 {
		if !SearchAll(w, r) {
			errorHandler(w, r, 404)
		}
	}
}

func SearchAll(w http.ResponseWriter, r *http.Request) bool {

	var id int

	k := r.FormValue("value")
	num := false
	for i, v := range k {
		if v >= '0' && v <= '9' && i > 0 && i < 3 {
			num = true
		}
	}

	key := strings.Trim(k, " ")
	flag := true
	artist := true
	groupwas := true
	member := true

	country := false

	city := true
	if !num {

		if !SearchMembers(w, r, key) {
			member = false
		}

		if !member {
			if !Search(w, r, k) {
				artist = false
			}
		}

		for _, v := range k {
			if v == '-' {
				country = true
				break
			}
		}

		if !artist && !member && !country {
			if !SearchByLocation(w, r, key) {
				city = false
			}
		}

		if !artist && !member && country {
			if !SearchCityGroup(w, r, k) {
				groupwas = false
			}
		}

		if !artist && !member && !groupwas && !city {
			errorHandler(w, r, http.StatusNotFound)
			flag = false
		}

		API.IDS = id

		if id > 0 {
			temp.ExecuteTemplate(w, "search", API)
		}

	} else {
		if !SearchAlbumDate(w, r, k) {
			if !SearchByCreationDate(w, r, k) {
				flag = false
				errorHandler(w, r, http.StatusNotFound)
			}
		}
	}
	return flag

}

//find 1 member -> artist, else -> members.html
func SearchMembers(w http.ResponseWriter, r *http.Request, k string) bool {

	key := strings.Trim(k, " ")
	fmt.Print(key)
	fm := 0
	for i, v := range API.Artist {
		for _, m := range v.Members {
			if strings.Contains(strings.ToLower(m), strings.ToLower(key)) {
				API.IDS = i - 1
				API.ID = i
				fm++
			}
		}
	}
	flag := false
	if fm == 0 {
		flag = false
	}
	if fm == 1 {
		temp.ExecuteTemplate(w, "artist", API)
		flag = true
	} else if fm > 1 {
		temp.ExecuteTemplate(w, "members", API)
		flag = true

	}
	return flag
}

func SearchByLocation(w http.ResponseWriter, r *http.Request, k string) bool {

	key := strings.Trim(k, " ")

	//japan
	flag := false
	var cit FindCity
	var cities []FindCity

	for _, v := range API.LocationsHtml.Index {
		for _, city := range v.Locations {
			fmt.Println(city)
			if strings.Contains(strings.ToLower(city), strings.ToLower(key)) {

				cit.CityName = city
				// API.Cities = append(API.Cities, city)
				cities = append(cities, cit)
			}
		}
	}

	if cities != nil {
		temp.ExecuteTemplate(w, "all", cities)
		flag = true
	}
	return flag

}

func SearchAlbumDate(w http.ResponseWriter, r *http.Request, key string) bool {
	flag := false
	var album AlbumDate

	for _, v := range API.Artist {

		if strings.Contains(v.FirstAlbum, key) {
			album.Name = v.Name

		}
	}

	if album.Name != "" {
		flag = true
	}

	temp.ExecuteTemplate(w, "all", album)
	return flag
}

// api.artist
func SearchByCreationDate(w http.ResponseWriter, r *http.Request, input string) bool {
	k, _ := strconv.Atoi(input)
	fmt.Print(k)
	flag := false

	var createDate CreationsDate
	var creDates []CreationsDate

	for _, v := range API.Artist {
		if v.CreationDate == k {
			createDate.Name = v.Name
			creDates = append(creDates, createDate)
		}
	}
	if creDates != nil {
		flag = true
	}

	temp.ExecuteTemplate(w, "all", creDates)
	return flag
}

func SearchCityGroup(w http.ResponseWriter, r *http.Request, k string) bool {

	var group GroupWasCity
	var groups []GroupWasCity

	key := strings.Trim(k, " ")
	flag := false

	for i, v := range API.LocationsHtml.Index {
		for _, city := range v.Locations {

			if strings.Contains(strings.ToLower(city), strings.ToLower(key)) {

				for j, a := range API.Artist {

					if i == j {
						group.Name = a.Name
						groups = append(groups, group)

					}
				}
			}
		}
	}

	if groups != nil {

		temp.ExecuteTemplate(w, "all", groups)
		flag = true
	}

	return flag

}

// web server, get api, by endpoint  -  handler -> get -> show all artists  - json -> get Data -> send struct
func Search(w http.ResponseWriter, r *http.Request, k string) bool {

	var id int

	key := strings.Trim(k, " ")
	temps := ""
	for i, v := range API.Artist {

		if strings.Contains(strings.ToLower(v.Name), strings.ToLower(key)) {
			// fmt.Println(v.Name)
			id = i
			temps = v.Name
		}
	}
	flag := false
	if id == 0 && temps != "queen" {
		fmt.Print("not found artist ")
	}
	if id >= 0 && temps != "" {
		API.IDS = id
		flag = true
		temp.ExecuteTemplate(w, "search", API)
	}

	return flag
}

//write data from, api []singers, middleware DB,
//then search array, by keword, and return find result -> client

//trim value, before after artist

//all handler -> if artist || member || date || album - search all in array artist an return result || not filter checkbox
//style - london uk, japan page, data page

// func Search(w http.ResponseWriter, r *http.Request, flma bool) {
// 	var key string
// 	var notArtist bool
// 	var id int

// 	if !flma {
// 		key = r.FormValue("value")
// 	} else if flma {
// 		k := r.Form["memart"]
// 		for _, v := range k {
// 			key += string(v)
// 		}
// 	}

// 	fmt.Print(key, flma)

// 	//1 lower case, do Big

// 	for i, v := range API.Artist {

// 		if strings.Contains(strings.ToLower(v.Name), strings.ToLower(key)) {
// 			// fmt.Println(v.Name)
// 			id = i
// 		// } else {
// 		// 	notArtist = true
// 		// 	break
// 		// }
// 	}
// }

// 	// fmt.Println(notArtist)
// 	// if notArtist {
// 	// 	for i, v := range API.Artist {
// 	// 		for _, m := range v.Members {
// 	// 			if strings.Contains(strings.ToLower(m), strings.ToLower(key)) {
// 	// 				id = i
// 	// 				break
// 	// 			}
// 	// 		}
// 	// 	}
// 	// 	fmt.Println(id)
// 	// }
// 	if id == 0 && key != "queen" {
// 		id = 52
// 		return
// 	}

// 	API.IDS = id

// 	//write data from, api []singers, middleware DB,
// 	//then search array, by keword, and return find result -> client

// 	temp.ExecuteTemplate(w, "search", API)

// }


// web server, get api, by endpoint  -  handler -> get -> show all artists  - json -> get Data -> send struct
func Search(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/search" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}

	if r.Method == "GET" {
		temp.ExecuteTemplate(w, "search", API.Artist)
	}

	if r.Method == "POST" {
		k := r.FormValue("keyword")

		var find string
		nw := []rune(k)
		//1 lower case, do Big
		lower := false
		key := k
		for i, v := range k {
			if v >= 'a' && v <= 'z' && i == 0 {
				lower = true
			}
		}
		if lower {

			for i, v := range k {
				if v >= 'a' && v <= 'z' && i == 0 {
					// tl = v - 32
					nw[0] = v - 32
					// nw = nw[:]
					key = string(nw)
				}
			}
		}

		// ..68413822c

		id := 0
		for i, v := range API.Artist {

			if strings.Contains(v.Name, key) {
				fmt.Println(v.Name)
				id = i
				break
			}
		}

		fmt.Println(id, "id conatin artist")

		if id == 52 {
			fmt.Println("not found lol")
			//redirect not found page
			return
		}

		//contains

		API.Name = find
		API.IDS = id
		//write data from, api []singers, middleware DB,
		//then search array, by keword, and return find result -> client
		temp.ExecuteTemplate(w, "search", API)
	}
}

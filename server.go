package main

import (
	"encoding/json"
	"fmt"
	"groupie/groupie"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Artists struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	DeezerName     string
	Image          string   `json:"image"`
	Members        []string `json:"members"`
	NbMembers      int
	Preview        string
	Locations      []string
	DatesLocations []DatesLocations `json:"datesLocations"`
	CreationDate   int              `json:"creationDate"`
	FirstAlbum     string           `json:"firstAlbum"`
	NbFans         int
	DeezerLink     string
}

type Locations struct {
	Locations []string `json:"locations"`
}

type DatesLocations struct {
	Locations []string `json:"locations"`
}

type DeezerData struct {
	Data []Data `json:"data"`
}

type Data struct {
	Preview string `json:"preview"`
	Artist  Artist `json:"artist"`
}

type Artist struct {
	Link string `json:"link"`
}

type DeezerArtist struct {
	NbFollowers int `json:"nb_fan"`
}

var artists []Artists

var loca = regexp.MustCompile(`/\d`)

func Groupie(w http.ResponseWriter, r *http.Request, artists []Artists) {
	var artistslist []Artists
	var MemberFilter string = r.FormValue("FilterMb")
	var FanFilter string = r.FormValue("FilterFan")
	var CreationFilter string = r.FormValue("FilterCreation")
	var FormArtist string = r.FormValue("SearchBar")
	MbFilter, _ := strconv.Atoi(MemberFilter)
	FanFilters, _ := strconv.Atoi(FanFilter)
	CreationFilters, _ := strconv.Atoi(CreationFilter)
	detId := strings.TrimPrefix(r.URL.Path, "/")
	marks, _ := strconv.Atoi(detId)
	template, err := template.ParseFiles("./pages/index.html", "./pages/details.html", "./templates/footer.html", "./templates/header.html")
	if err != nil {
		log.Fatal(err)
	}
	switch {
	case FormArtist != "":
		for g := 0; g < len(artists); g++ {
			if groupie.Capitalize(FormArtist) == artists[g].Name || groupie.ToHigher(FormArtist) == artists[g].Name || FormArtist == artists[g].Name {
				Details(w, r, artists[g])
			}
		}
	case MbFilter != 0:
		artistslist = FilterMember(artistslist, MbFilter)
		template.Execute(w, FilterMember(artists, MbFilter))
	case FanFilters != 0:
		artistslist = FilterFan(artistslist, FanFilters)
		template.Execute(w, FilterFan(artists, FanFilters))
	case CreationFilters != 0:
		artistslist = FilterCreationDate(artistslist, CreationFilters)
		template.Execute(w, FilterCreationDate(artists, CreationFilters))
	// case AlbumFilters != 0:
	// 	artistslist = FilterAlbumDate(artistslist, AlbumFilters)
	// 	template.Execute(w, FilterAlbumDate(artists, AlbumFilters))
	case loca.MatchString(r.URL.Path):
		Details(w, r, artists[marks-1])
	default:
		template.Execute(w, artists)
	}
}

func FilterMember(artists []Artists, NbMember int) []Artists {
	var artistss []Artists
	for i := 0; i < len(artists); i++ {
		if artists[i].NbMembers <= NbMember {
			artistss = append(artistss, artists[i])
		}
	}
	return artistss
}

func FilterFan(artists []Artists, NbFans int) []Artists {
	var artistss []Artists
	for i := 0; i < len(artists); i++ {
		if artists[i].NbFans <= NbFans {
			artistss = append(artistss, artists[i])
		}
	}
	return artistss
}

func FilterCreationDate(artists []Artists, CreationDate int) []Artists {
	var artistss []Artists
	for i := 0; i < len(artists); i++ {
		if artists[i].CreationDate <= CreationDate {
			artistss = append(artistss, artists[i])
		}
	}
	return artistss
}

func FilterAlbumDate(artists []Artists, FirstAlbum string) []Artists {
	var artistss []Artists
	for i := 0; i < len(artists); i++ {
		if artists[i].FirstAlbum <= FirstAlbum {
			artistss = append(artistss, artists[i])
		}
	}
	return artistss
}

func Details(w http.ResponseWriter, r *http.Request, artists Artists) {
	template, err := template.ParseFiles("./pages/details.html", "./templates/footer.html", "./templates/header.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, artists)
}

func GetFromGroupAPI() []Artists {
	response, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(responseData, &artists)
	return (artists)
}

func SetLocationsToArt(artists []Artists) []Artists {
	for i := 0; i < len(artists); i++ {
		response, err := http.Get("https://groupietrackers.herokuapp.com/api/locations/" + strconv.Itoa(artists[i].Id))
		if err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		var responseObject Locations

		json.Unmarshal(responseData, &responseObject)

		if len(responseObject.Locations) != 0 {
			artists[i].Locations = responseObject.Locations
		}
	}
	return artists
}

func replaceSpaces(s string) string {
	return strings.ReplaceAll(s, " ", "%20")
}

func SetDeezerName(artists []Artists) []Artists {
	for i := 0; i < len(artists); i++ {
		artists[i].DeezerName = replaceSpaces(artists[i].Name)
	}
	return artists
}

func GetDeezerPreviews(artists []Artists) []Artists {
	for i := 0; i < len(artists); i++ {
		nom := artists[i].DeezerName
		response, err := http.Get("https://api.deezer.com/search?q=" + nom)
		if err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		var responseObject DeezerData

		json.Unmarshal(responseData, &responseObject)

		if len(responseObject.Data) != 0 {
			artists[i].Preview = responseObject.Data[0].Preview
			artists[i].DeezerLink = responseObject.Data[0].Artist.Link
		}
	}
	return artists
}

func DeezerLinkToNbFans(artists []Artists) []Artists {
	for i := 0; i < len(artists); i++ {
		linkId := strings.Split(artists[i].DeezerLink, "/")
		countj := 0
		for j := 0; j < len(linkId)-1; j++ {
			countj++
		}
		response, err := http.Get("https://api.deezer.com/artist/" + linkId[countj])
		countj = 0
		if err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		var responseObject DeezerArtist

		json.Unmarshal(responseData, &responseObject)

		if responseObject.NbFollowers != 0 {
			artists[i].NbFans = responseObject.NbFollowers
		}
	}
	return artists
}

func TestPrint(artists []Artists) {
	for i := 0; i < len(artists); i++ {
		fmt.Println(artists[i].Id)
		fmt.Println(artists[i].Name)
		fmt.Println(artists[i].Image)
		fmt.Println(artists[i].Members)
		fmt.Println(artists[i].Preview)
		fmt.Println(artists[i].Locations)
		fmt.Println(artists[i].CreationDate)
		fmt.Println(artists[i].FirstAlbum)
		fmt.Println(artists[i].NbMembers)
		fmt.Println(artists[i].DeezerLink)
		fmt.Println(artists[i].NbFans)
	}
}

func CountMembers(artists []Artists) []Artists {
	count := 0
	for j := 0; j < len(artists); j++ {
		for i := 0; i < len(artists[j].Members); i++ {
			count++
		}
		artists[j].NbMembers = count
		count = 0
	}
	return artists
}

func main() {
	GetFromGroupAPI()
	SetDeezerName(artists)
	GetDeezerPreviews(artists)
	SetLocationsToArt(artists)
	CountMembers(artists)
	DeezerLinkToNbFans(artists)
	TestPrint(artists)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		Groupie(w, r, artists)
	})
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.Handle("/Ressources/", http.StripPrefix("/Ressources/", http.FileServer(http.Dir("./Ressources"))))
	fmt.Println("Server running on port :8080")
	http.ListenAndServe(":8080", nil)
}

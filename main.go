package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
)

type Response struct {
	Results []any
	Count uint
}

type Car struct {
	Id uint
	Brand string
	Model string
}

const indexLocation string = "./index.html";

var indexFile []byte;
var cars []Car = []Car{
	{
		Id: 0,
		Brand: "Chevrolet",
		Model: "Onix",
	},
	{
		Id: 1,
		Brand: "Chevrolet",
		Model: "Cruze",
	},
	{
		Id: 2,
		Brand: "Ford",
		Model: "Mustang",
	},
	{
		Id: 3,
		Brand: "Ford",
		Model: "Fiesta",
	},
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func getListParams(r *http.Request)(page uint, size uint){
	page = 0;
	size = 2;
	q := r.URL.Query()
	//fmt.Println(q)
	if(q.Has("page")){
		tmp := q.Get("page")
		tmp2, err := strconv.ParseUint(tmp, 10, 32);
		if(err == nil){
			page = uint(tmp2)
		}
	}
	if(q.Has("size")){
		tmp := q.Get("size")
		tmp2, err := strconv.ParseUint(tmp, 10, 32);
		if(err == nil){
			size = uint(tmp2)
		}
	}
	return
}

func minUint(num1 uint, num2 uint)(uint){
	if(num1 < num2){
		return num1
	}
	return num2
}

func handleCars(w http.ResponseWriter, r *http.Request) {
	page, size := getListParams(r);
	//fmt.Println(page, size)
	var rdata []any = []any{};
	for _, v := range cars[(page*size):minUint(((page+1)*size), uint(len(cars)))] {
		rdata = append(rdata, v)
	}
	byteArr, err := json.Marshal(Response{Results: rdata, Count: uint(len(cars))})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(byteArr)
}

func handleCar(w http.ResponseWriter, r *http.Request){
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	for _, v := range cars {
		if(v.Id == uint(id)){
			byteArr, err := json.Marshal(v)
			if(err != nil){
				w.WriteHeader(http.StatusInternalServerError)
			}else{
				w.Header().Set("Content-Type", "application/json")
				w.Write(byteArr)
			}
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func handleCarDelete(w http.ResponseWriter, r *http.Request){}

func handleCarAdd(w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()
	body, err1 := ioutil.ReadAll(r.Body);
	if(err1 != nil){
		//fmt.Println(err1)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err1.Error()))
	}
	//fmt.Println(string(body), uint(rand.Uint64()))
	var newCar Car
	err2 := json.Unmarshal(body, &newCar)
	if(err2 != nil){
		//fmt.Println(err2)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err2.Error()))
	}
	fmt.Println(newCar)
	newCar.Id = uint(rand.Uint64());
	cars = append(cars, newCar)
	w.WriteHeader(http.StatusCreated)
}

func handleDefault(w http.ResponseWriter, r *http.Request) {
	//lee el file cada vez! bueno para debug pero en prod cargar en meme!
	http.ServeFile(w, r, indexLocation)
}

func handleDefault2(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(indexFile)
}

func init() {
	var err error
	indexFile, err = os.ReadFile(indexLocation)
	check(err)
}

func main() {
	mux := http.NewServeMux()
	if(os.Getenv("ENV") == "PROD"){
		mux.HandleFunc("/", handleDefault2)
	}else{
		mux.HandleFunc("/", handleDefault)
	}
	mux.HandleFunc("GET /api/cars", handleCars)
	mux.HandleFunc("POST /api/cars", handleCarAdd)
	mux.HandleFunc("GET /api/cars/{id}", handleCar)
	mux.HandleFunc("DELETE /api/cars/{id}", handleCarDelete)
	mux.HandleFunc("/favicon.ico", http.NotFound)
	err := http.ListenAndServe("0.0.0.0:8080", mux)
	check(err)
}

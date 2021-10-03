package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {

	// TODO:make concurrently safe
	h := act{
		Action:  make(map[int]JsAct),
		counter: 0,
	}

	http.HandleFunc("/add", h.action)
	http.HandleFunc("/del", h.delete)
	http.HandleFunc("/rew", h.rewrite)
	port := ":9090"
	println("Server listen on port", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("ListernAndServe", err)
	}
}

type act struct {
	Action  map[int]JsAct
	counter int
}

type JsAct struct {
	St string
}

type Del struct {
	Id int
}

type Calcul struct {
	Cltr []int
}

type Rew struct {
	Num int
	Str string
}

func (h *act) action(w http.ResponseWriter, r *http.Request) {
	var typ JsAct

	str, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(str, &typ)

	if err != nil {
		fmt.Println("Error", err)
	}

	if typ.St == "" {
		http.Error(w, "Missing Field \"St\"", http.StatusBadRequest)
		return
	}

	h.counter++
	h.Action[h.counter] = typ
	js, _ := json.Marshal(h.Action)
	w.Write(js)
}

func (h *act) delete(w http.ResponseWriter, r *http.Request) {
	var n Del
	var s Calcul
	s.Cltr = append(s.Cltr, 0)

	str, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(str, &n)

	if err != nil {
		http.Error(w, "Missing Field \"Id\"", http.StatusBadRequest)
		return
	}

	for _, v := range s.Cltr {
		if v == n.Id {
			http.Error(w, "Missing Field \"Id\"", http.StatusBadRequest)
			return
		}
	}

	s.Cltr = append(s.Cltr, n.Id)

	for a := range h.Action {
		if a == n.Id {
			delete(h.Action, n.Id)
		}
	}

	js, err := json.Marshal(h.Action)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Write(js)
}

func (h *act) rewrite(w http.ResponseWriter, r *http.Request) {
	var mp Rew

	str, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(str, &mp)

	if err != nil {
		http.Error(w, "Missing Fields \"Num\", \"Str\"", http.StatusBadRequest)
		return
	}

	for a := range h.Action {
		if a == mp.Num {
			h.Action[a] = JsAct{mp.Str}
		}
	}

	js, err := json.Marshal(h.Action)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Write(js)
}

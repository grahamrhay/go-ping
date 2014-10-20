package main

import (
	"container/ring"
	"html/template"
	"log"
	"net/http"
)

func startServer(ring **ring.Ring) {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	http.HandleFunc("/", makeHandler(ring))
	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()
}

type Point struct {
	X int64
	Y float64
}

type Data struct {
	Avg []Point
}

func makeHandler(ring **ring.Ring) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := &Data{Avg: []Point{}}
		ring.Do(func(value interface{}) {
			if value == nil {
				return
			}

			res := value.(*PingResult)
			data.Avg = append(data.Avg, *&Point{X: res.Time.Unix(), Y: res.Avg})
		})
		t, _ := template.ParseFiles("index.html")
		t.Execute(w, data)
	}
}

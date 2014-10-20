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
	Min  []Point
	Avg  []Point
	Max  []Point
	Mdev []Point
}

func makeHandler(ring **ring.Ring) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := &Data{Avg: []Point{}}
		ring.Do(func(value interface{}) {
			if value == nil {
				return
			}

			res := value.(*PingResult)
			unix := res.Time.Unix()
			data.Min = append(data.Min, *&Point{X: unix, Y: res.Min})
			data.Avg = append(data.Avg, *&Point{X: unix, Y: res.Avg})
			data.Max = append(data.Max, *&Point{X: unix, Y: res.Max})
			data.Mdev = append(data.Mdev, *&Point{X: unix, Y: res.Mdev})
		})
		t, _ := template.ParseFiles("index.html")
		t.Execute(w, data)
	}
}

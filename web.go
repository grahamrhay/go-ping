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

type Data struct {
	Times []string
	Avg   []float64
}

func makeHandler(ring **ring.Ring) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := &Data{Times: []string{}, Avg: []float64{}}
		ring.Do(func(value interface{}) {
			if value == nil {
				return
			}

			res := value.(*PingResult)
            data.Times = append(data.Times, res.Time.Format("15:04"))
			data.Avg = append(data.Avg, res.Avg)
		})
		t, _ := template.ParseFiles("index.html")
		t.Execute(w, data)
	}
}

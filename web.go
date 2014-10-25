package main

import (
	"container/ring"
	"html/template"
	"log"
	"net/http"
)

func startServer(startOfLast3Hours **ring.Ring, startOfLast30Hours **ring.Ring) {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	http.HandleFunc("/", makeHandler(startOfLast3Hours, startOfLast30Hours))
	go func() {
		port := "8080"
		log.Println("Listening on:", port)
		log.Fatal(http.ListenAndServe(":"+port, nil))
	}()
}

type Chart struct {
	Times []string
	Avg   []float64
}

type Data struct {
	Last3Hours  *Chart
	Last30Hours *Chart
}

func makeHandler(startOfLast3Hours **ring.Ring, startOfLast30Hours **ring.Ring) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chart := &Chart{Times: []string{}, Avg: []float64{}}
		chart2 := &Chart{Times: []string{}, Avg: []float64{}}
		data := &Data{Last3Hours: chart, Last30Hours: chart2}

		startOfLast3Hours.Do(func(value interface{}) {
			if value == nil {
				return
			}

			res := value.(*PingResult)
			chart.Times = append(chart.Times, res.Time.Format("15:04"))
			chart.Avg = append(chart.Avg, res.Avg)
		})

		startOfLast30Hours.Do(func(value interface{}) {
			if value == nil {
				return
			}

			res := value.(*PingResult)
			chart2.Times = append(chart2.Times, res.Time.Format("Mon 15:04"))
			chart2.Avg = append(chart2.Avg, res.Avg)
		})

		t, _ := template.ParseFiles("index.html")
		t.Execute(w, data)
	}
}

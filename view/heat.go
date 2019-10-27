package view

import (
	"encoding/json"
	"fuzz_debug_platform/sqlfuzz"
	"log"
	"net/http"
)

type HeatView struct {
	Number int  `json:"number"`
	Alter  int  `json:"alter"`
	Heat   float32  `json:"heat"`
}


func Heat() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		heatViews := make([]*HeatView, 0)

		sqlfuzz.ACounter.Mu.RLock()
		for pos, counter := range sqlfuzz.ACounter.AlterMap {
			heatViews = append(heatViews, &HeatView{
				Number:pos[0],
				Alter:pos[1],
				Heat:float32(counter.Fail) / float32(counter.Succ + counter.Fail),
			})
		}
		sqlfuzz.ACounter.Mu.RUnlock()

		jsonBytes, err := json.Marshal(heatViews)
		if err != nil {
			log.Fatalf("should not hear %v \n", err)
		}

		writer.Write(jsonBytes)
	}
}

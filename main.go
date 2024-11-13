package main

import (
	"encoding/json"
	"fmt"
	"gotools2/modbus2"
	"log"
	"net/http"
	"time"

	"github.com/simonvetter/modbus"
)

type Data struct {
	Signal []string `json:"signal"`
	Values []int    `json:"values"`
}

var (
	res = modbus2.Res{}

	conf = modbus2.Conf{}
)

func serveIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	data := Data{
		Signal: res.Name,
		Values: res.Res,
	}

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func updateValues(mc *modbus.ModbusClient) {
	for {
		time.Sleep(1 * time.Second)
		conf.Read(mc, &res)
	}
}

func main() {
	mc := modbus2.CreateModbusClient("localhost", "1502")
	conf.Decode()
	fmt.Println(conf)
	go updateValues(mc)

	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/data", dataHandler)

	log.Println("Serveur démarré sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

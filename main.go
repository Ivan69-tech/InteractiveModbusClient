package main

import (
	"gotools2/modbus2"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/simonvetter/modbus"
)

type Data struct {
	Keys   []string
	Values []int
}

var res map[string]int

func viewHandler(w http.ResponseWriter, r *http.Request, data map[string]int) {
	tmpl, err := template.ParseFiles("index.html") // Charger le fichier HTML
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	keys := make([]string, 0, len(data))
	values := make([]int, 0, len(data))
	for k, v := range data {
		keys = append(keys, k)
		values = append(values, v)
	}

	err = tmpl.Execute(w, Data{
		Keys:   keys,
		Values: values,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func updateValues(mc *modbus.ModbusClient, data map[string]int) {
	for {
		time.Sleep(2 * time.Second)
		res = modbus2.Read(mc, data)
	}
}

func main() {

	mc := modbus2.CreateModbusClient("localhost", "1502")
	a := modbus2.Decode()
	go updateValues(mc, a)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		viewHandler(w, r, res)
	})

	log.Println("Serveur démarré sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

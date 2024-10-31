package main

import (
	"fmt"
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

var (
	res = modbus2.Res{}

	conf = modbus2.Conf{}
)

func viewHandler(w http.ResponseWriter, r *http.Request, keys []string, values []int) {
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, Data{
		Keys:   keys,
		Values: values,
	})
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
	mc := modbus2.CreateModbusClient("192.168.0.190", "5502")
	conf.Decode()
	fmt.Println(conf)
	go updateValues(mc)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		viewHandler(w, r, res.Name, res.Res)
	})

	log.Println("Serveur démarré sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

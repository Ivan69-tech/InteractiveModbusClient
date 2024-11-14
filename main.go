package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"gotools2/logs"
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

type DataReceived struct {
	Host string `json:"host"`
	Path string `json:"path"`
}

type LogData struct {
	Log string `json:"log"`
}

var (
	res = modbus2.Res{}

	conf = modbus2.Conf{}
)

//go:embed html/index.html
//go:embed html/logs.html
var content embed.FS

func serveIndex(w http.ResponseWriter, r *http.Request) {
	// Serve le fichier index.html depuis l'embed
	data, err := content.ReadFile("html/index.html")
	if err != nil {
		http.Error(w, "Erreur lors du chargement du fichier", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(data)
}

func serveLogs(w http.ResponseWriter, r *http.Request) {

	data, err := content.ReadFile("html/logs.html")
	if err != nil {
		http.Error(w, "Erreur lors du chargement du fichier", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(data)
}

func logsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	logger := LogData{
		Log: logs.LogsBuffer.String(),
	}

	err := json.NewEncoder(w).Encode(logger)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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

func sendDataHandler(w http.ResponseWriter, r *http.Request) {

	var dataReceived DataReceived

	err := json.NewDecoder(r.Body).Decode(&dataReceived)
	if err != nil {
		http.Error(w, "Erreur lors du décodage de la donnée", http.StatusBadRequest)
		return
	}

	fmt.Println("Donnée reçue:", dataReceived)

	response := map[string]string{
		"status":  "success",
		"message": "Donnée reçue avec succès",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	mc, err := modbus2.CreateModbusClient(dataReceived.Host)

	// pour ne pas quitter le programme si l'utilisateur se trompe dans le host
	if err != nil {
		return
	}

	conf.Decode("conf-copy.csv")
	go updateValues(mc)
}

func updateValues(mc *modbus.ModbusClient) {
	for {
		time.Sleep(1 * time.Second)
		conf.Read(mc, &res)
	}
}

func main() {

	logs.StartLogging()
	go logs.ResetBuffer()
	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/data", dataHandler)
	http.HandleFunc("/sendData", sendDataHandler)
	http.HandleFunc("/logs", serveLogs)
	http.HandleFunc("/getlogs", logsHandler)

	log.Println("Serveur démarré sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"gotools2/database"
	"gotools2/logs"
	"gotools2/modbus2"
	"gotools2/server"
	"log"
	"net/http"
	"sync"
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

type GetRequest struct {
	Register int    `json:"register"`
	DataSize string `json:"dataSize"`
	DataType string `json:"dataType"`
	Value    int    `json:"value"`
}

type LogData struct {
	Log string `json:"log"`
}

type Db struct {
	Time    []time.Time `json:"time"`
	Signaux []string    `json:"signaux"`
	Data    [][]int     `json:"data"`
}

var (
	res = modbus2.Res{}

	conf = modbus2.Conf{}

	db = database.Database{}

	// Ajout d'un mutex pour synchroniser l'accès à `res`
	resMutex sync.Mutex

	mcGlobal *modbus.ModbusClient

	stopChannel chan bool
)

//go:embed static/html/*.html
//go:embed static/js/*.js
//go:embed README.md

var content embed.FS

func serveReadme(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "README.md")
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
		fmt.Println(err)
		return
	}

	fmt.Println("Donnée reçue:", dataReceived)

	response := map[string]string{
		"status":  "success",
		"message": "Donnée reçue avec succès",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	// Arrêter la goroutine en cours si elle existe
	if stopChannel != nil {
		stopChannel <- true
	}

	mc, err := modbus2.CreateModbusClient(dataReceived.Host)
	mcGlobal = mc

	// Pour ne pas quitter le programme si l'utilisateur se trompe dans le host
	if err != nil {
		fmt.Println("Erreur lors de la création du client Modbus:", err)
		return
	}

	conf = modbus2.Conf{}
	conf.Decode("conf-copy.csv")

	// Créer un nouveau channel pour arrêter la nouvelle goroutine
	stopChannel = make(chan bool)

	// Démarrer une nouvelle goroutine pour mettre à jour les valeurs
	go updateValues(mc, stopChannel)
}

func GetModbusWrite(w http.ResponseWriter, r *http.Request) {

	var getRequest GetRequest

	err := json.NewDecoder(r.Body).Decode(&getRequest)
	if err != nil {
		http.Error(w, "Erreur lors du décodage de la donnée", http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	//fmt.Println(getRequest)

	var writeReq = modbus2.WriteReq{
		Register: getRequest.Register,
		DataSize: getRequest.DataSize,
		DataType: getRequest.DataType,
		Value:    getRequest.Value,
	}

	//fmt.Println("mc= ", mcGlobal)

	modbus2.Write(mcGlobal, writeReq)

	return
}

func sendDbHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	data := Db{
		Time:    db.Time,
		Signaux: db.Signaux,
		Data:    db.Data,
	}

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func updateValues(mc *modbus.ModbusClient, stop chan bool) {
	for {
		select {
		case <-stop:
			return // Arrêter la goroutine
		default:
			time.Sleep(1 * time.Second)
			res = modbus2.Res{}
			conf.Read(mc, &res)
			db.Save(res)
		}
	}
}

func main() {
	logs.StartLogging()
	go logs.ResetBuffer()
	go server.Server()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	fs := http.FileServer(http.Dir("./static/html"))
	http.Handle("/", http.StripPrefix("/", fs))

	http.HandleFunc("/data", dataHandler)
	http.HandleFunc("/sendData", sendDataHandler)
	http.HandleFunc("/getlogs", logsHandler)
	http.HandleFunc("/getdb", sendDbHandler)
	http.HandleFunc("/readme", serveReadme)
	http.HandleFunc("/writemodbus", GetModbusWrite)

	log.Println("Serveur démarré sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

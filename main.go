package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
)

func getenvStr(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

// FIXME: We can't use a int value with os.Getenv ???
// func getenvInt(key string, defaultValue int) int {
// 	value := os.Getenv(key)

// 	if len(value) == 0 {
// 		return defaultValue
// 	}
// 	return value
// }

var host = getenvStr("DB_HOST", "localhost")
var port = 5432
var user = getenvStr("DB_USER", "postgres")
var password = getenvStr("DB_PASSWORD", "566245")
var dbname = getenvStr("DB_PORT", "density")

type Sensor struct {
	Document_id string `json:"id"`
	Sensor_id   string `json:"sensor_id"`
	Room_id     string `json:"room_id"`
	Floor_id    string `json:"floor_id"`
	Building_id string `json:"building_id"`
	Created_at  string `json:"created_at"`
}

var upgrader = websocket.Upgrader{} // use default options

func socketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade our raw HTTP connection to a websocket based one
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Error during connection upgradation:", err)
		return
	}
	defer conn.Close()

	// The event loop
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error during message reading:", err)
			break
		}

		log.Printf("Received: %s", message)

		var jsonMap map[string]interface{}
		json.Unmarshal([]byte(message), &jsonMap)

		fmt.Println(jsonMap)

		// connection string
		psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		// open database
		db, err := sql.Open("postgres", psqlconn)
		if err != nil {
			panic(err)
		}

		// close database
		defer db.Close()

		// check db
		err = db.Ping()
		CheckError(err)

		fmt.Println("Connected!")

		_, err = db.Exec("INSERT INTO sensor.europe_sensors (sensor_id, room_id, floor_id, building_id) VALUES ($1,$2,$3,$4)", jsonMap["sensor_id"], jsonMap["room_id"], jsonMap["floor_id"], jsonMap["building_id"])

		if err != nil {
			log.Println("Error during insert into sensor:", err)
		}

		err = conn.WriteMessage(messageType, message)
		if err != nil {
			log.Println("Error during message writing:", err)
			break
		}
	}
}

func list(w http.ResponseWriter, r *http.Request) {
	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	var sensors []Sensor
	// open database
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		panic(err)
	}

	// close database
	defer db.Close()

	// check db
	err = db.Ping()
	CheckError(err)

	rows, err := db.Query("SELECT * FROM sensor.europe_sensors")
	if err != nil {
		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.Execute(w, sensors)
		return
	}

	tmpl := template.Must(template.ParseFiles("index.html"))
	for rows.Next() {
		var s Sensor
		//&s.current_time
		err = rows.Scan(&s.Document_id, &s.Sensor_id, &s.Room_id, &s.Floor_id, &s.Building_id, &s.Created_at)
		if err != nil {
			fmt.Println(err)
			panic(err)

		}
		sensors = append(sensors, s)
	}
	fmt.Println(sensors)
	tmpl.Execute(w, sensors)
	defer rows.Close()

}

func main() {
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.HandleFunc("/socket", socketHandler)
	http.HandleFunc("/", list)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

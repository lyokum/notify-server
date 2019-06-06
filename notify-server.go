package main

import (
	"encoding/json"
	"gitlab.com/lyokum/update"
	"io/ioutil"
	"log"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
)

var (
	toggleMutex      = &sync.Mutex{}
	notifFlag   bool = false
)

func toggleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// flip flag
		toggleMutex.Lock()
		notifFlag = !notifFlag
		toggleMutex.Unlock()
	}
}

func notifyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// check flag
		toggleMutex.Lock()
		displayNotif := notifFlag
		toggleMutex.Unlock()

		if !displayNotif {
			return
		}

		// read body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		// create update from body
		var update update.Update
		err = json.Unmarshal([]byte(body), &update)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		// send notification
		update.Notify()
	} else {
		http.Error(w, "Invalid method: must use GET", http.StatusSeeOther)
	}

}

/* -----Main Execution----- */

func main() {
	// ignore SIGHUP
	signal.Ignore(syscall.SIGHUP)

	// assign handlers
	http.HandleFunc("/", notifyHandler)
	http.HandleFunc("/toggle", toggleHandler)

	// run server on port 8080
	log.Fatal(http.ListenAndServe(":8080", nil))
}

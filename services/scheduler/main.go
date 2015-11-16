package main

import (
	"../../lib/liboct"
	"../../lib/routes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const SchedularCacheDir = "/tmp/.test_schedular_cache"

type SchedulerConfig struct {
	Port           int
	ServerListFile string
	CacheDir       string
	Debug          bool
}

func ReceiveTaskCommand(w http.ResponseWriter, r *http.Request) {
	var ret liboct.HttpRet
	id := r.URL.Query().Get(":ID")
	sInterface, ok := liboct.DBGet(liboct.DBScheduler, id)
	if !ok {
		ret.Status = liboct.RetStatusFailed
		ret.Message = "Invalid task id"
		ret_string, _ := json.MarshalIndent(ret, "", "\t")
		w.Write([]byte(ret_string))
		return
	}
	s, _ := liboct.SchedulerFromString(sInterface.String())

	result, _ := ioutil.ReadAll(r.Body)
	fmt.Println("Receive task Command ", string(result))
	r.Body.Close()
	/* Donnot use this now FIXME
	var cmd liboct.TestActionCommand
	json.Unmarshal([]byte(result), &cmd)
	*/
	action, ok := liboct.TestActionFromString(string(result))
	if !ok {
		ret.Status = liboct.RetStatusFailed
		ret.Message = "Invalid action"
	} else {
		if s.Command(action) {
			ret.Status = liboct.RetStatusOK
		} else {
			ret.Status = liboct.RetStatusFailed
			ret.Message = "Failed to run the action"
		}
	}

	ret_string, _ := json.MarshalIndent(ret, "", "\t")
	w.Write([]byte(ret_string))
}

func ReceiveTask(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ReceiveTask begin")
	var ret liboct.HttpRet
	var tc liboct.TestCase
	realURL, _ := liboct.ReceiveFile(w, r, SchedularCacheDir)
	tc, err := liboct.CaseFromTar(realURL, "")
	fmt.Println(realURL)
	if err != nil {
		ret.Status = liboct.RetStatusFailed
		ret.Message = err.Error()
		ret_string, _ := json.MarshalIndent(ret, "", "\t")
		w.Write([]byte(ret_string))
		return
	} else {
	}

	s, _ := liboct.SchedulerNew(tc)

	if s.Command(liboct.TestActionApply) {
		if id, ok := liboct.DBAdd(liboct.DBScheduler, s); ok {
			fmt.Println("Add ok ", id)
			ret.Status = liboct.RetStatusOK
			ret.Message = id
			ret_string, _ := json.MarshalIndent(ret, "", "\t")
			w.Write([]byte(ret_string))
			return
		} else {
			fmt.Println("Failed to add db ", s)
		}
	}
	ret.Status = liboct.RetStatusFailed
	ret.Message = "Failed in allocate the resource"
	ret_string, _ := json.MarshalIndent(ret, "", "\t")
	w.Write([]byte(ret_string))
	return
}

// Will use DB in the future, (mongodb for example)
func init() {
	sf, err := os.Open("scheduler.conf")
	if err != nil {
		fmt.Errorf("Cannot find scheduler.conf.")
		return
	}
	defer sf.Close()

	if err = json.NewDecoder(sf).Decode(&pubConfig); err != nil {
		fmt.Errorf("Error in parse scheduler.conf")
		return
	}

	liboct.DBRegist(liboct.DBResource)
	liboct.DBRegist(liboct.DBScheduler)
	if len(pubConfig.ServerListFile) == 0 {
		return
	}

	slf, err := os.Open(pubConfig.ServerListFile)
	if err != nil {
		return
	}
	defer slf.Close()

	var rl []liboct.Resource
	if err = json.NewDecoder(slf).Decode(&rl); err != nil {
		return
	}

	for index := 0; index < len(rl); index++ {
		if _, ok := liboct.DBAdd(liboct.DBResource, rl[index]); ok {
			fmt.Println(rl[index])
		}
	}
}

var pubConfig SchedulerConfig

func main() {
	mux := routes.New()

	mux.Get("/resource", GetResource)
	mux.Post("/resource", PostResource)
	mux.Get("/resource/:ID/status", GetResourceStatus)

	mux.Post("/task", ReceiveTask)
	mux.Post("/task/:ID", ReceiveTaskCommand)

	http.Handle("/", mux)
	port := fmt.Sprintf(":%d", pubConfig.Port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

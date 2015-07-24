package main

import (
	"github.com/ant0ine/go-json-rest/rest"
	"log"
	"net/http"
	"sync"
	"strconv"
)

func main() {
	init_db ()
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/os", GetOS),
		rest.Post("/os", PostOS),
		rest.Delete("/os/:distribution", DeleteOS),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Fatal(http.ListenAndServe(":8080", api.MakeHandler()))
}

type OS struct {
	ID string
	Distribution string
	Version string
	Arch string
	CPU int64
	Memory int64
	IP string
	locked bool
}

var store = map[string]*OS{}

var lock = sync.RWMutex{}

func GetOSQuery(r *rest.Request) (os OS) {
	os.Distribution = r.URL.Query().Get("Distribution")
	os.Version = r.URL.Query().Get("Version")
	os.Arch = r.URL.Query().Get("Arch")

	var cpu string
	cpu = r.URL.Query().Get("CPU")
	if len(cpu) > 0 {
        	cpu_count, cpu_err := strconv.ParseInt(cpu, 10, 64)
 	        if cpu_err != nil {
                	//TODO, should report the err
		} else {
			os.CPU = cpu_count
		}
        } else {
		os.CPU = 0
	}

	var memory string
	memory = r.URL.Query().Get("Memory")
	if len(memory) > 0 {
        	memory_count, memory_err := strconv.ParseInt(cpu, 10, 64)
 	        if memory_err != nil {
                	//TODO, should report the err
		} else {
			os.Memory = memory_count
		}
        } else {
		os.Memory = 0
	}

	log.Println(os)
	return os
}

// Will use sql to seach, for now, just
func GetAvaliableResource(os_query OS) (ID string) {
	for _, os := range store {
		if len(os_query.Distribution) > 1 {
			if os_query.Distribution != (*os).Distribution {
				continue
			}
		}
		if len(os_query.Version) > 1 {
			if os_query.Version != (*os).Version {
				continue
			}
		}
		if len(os_query.Arch) > 1 {
			if os_query.Arch != (*os).Arch {
				continue
			}
		}
		if os_query.CPU >  (*os).CPU {
			log.Println("not enough CPU")
			continue
		}
		if os_query.Memory > (*os).Memory {
			log.Println("not enough Memory")
			continue
		}
		ID = (*os).ID
		return ID
	}
	return ""
}

func GetOS(w rest.ResponseWriter, r *rest.Request) {
	var os_query OS
	os_query = GetOSQuery (r)
	if len(os_query.Distribution) < 1 {
		GetAllOS(w, r)
		return
	}
	lock.RLock()

	var ID string
	ID = GetAvaliableResource(os_query)
	lock.RUnlock()

	log.Println(ID)
	if len(ID) < 1 {
		rest.NotFound(w, r)
		return
	}
	w.WriteJson(ID)
}

func GetAllOS(w rest.ResponseWriter, r *rest.Request) {
	lock.RLock()
	os_list := make([]OS, len(store))
	i := 0
	for _, os := range store {
		os_list[i] = *os
		i++
	}
	lock.RUnlock()
	w.WriteJson(&os_list)
}

func PostOS(w rest.ResponseWriter, r *rest.Request) {
	os := OS{}
	err := r.DecodeJsonPayload(&os)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if os.Distribution == "" {
		rest.Error(w, "os distribution required", 400)
		return
	}
	if os.Version == "" {
		rest.Error(w, "os version required", 400)
		return
	}
	if os.Arch == "" {
		rest.Error(w, "os arch required", 400)
		return
	}
	lock.Lock()
	store[os.Distribution] = &os
	lock.Unlock()
	w.WriteJson(&os)
}

func DeleteOS(w rest.ResponseWriter, r *rest.Request) {
	Distribution := r.PathParam("Distribution")
	lock.Lock()
	delete(store, Distribution)
	lock.Unlock()
	w.WriteHeader(http.StatusOK)
}

// Will use DB in the future, (mongodb for example)
// for now, just two demo hosts
func init_db () {
	var os [2]OS
	os[0].Distribution = "Ubuntu"
	os[0].Version = "12.04"
	os[0].Arch = "x86_64"
	os[0].ID = "0001"
	os[0].CPU = 2
	os[0].Memory = 4
	os[0].IP = "192.168.0.1"
	os[0].locked = false
	store[os[0].ID] = &os[0]

	os[1].Distribution = "CentOS"
	os[1].Version = "7"
	os[1].Arch = "x86_64"
	os[1].ID = "0002"
	os[1].CPU = 1
	os[1].Memory = 3
	os[1].IP = "192.168.0.2"
	os[1].locked = false
	store[os[1].ID] = &os[1]
}

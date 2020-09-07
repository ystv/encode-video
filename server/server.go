package server

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"

	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/tasks"
	task "github.com/ystv/encode-video/tasks"
)

// StartServer creates a new server
func StartServer(taskserver *machinery.Server) {
	r := http.NewServeMux()
	r.HandleFunc("/encode/single", func(w http.ResponseWriter, r *http.Request) {
		p := new(task.Payload)
		err := json.NewDecoder(r.Body).Decode(p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		reqJSON, err := json.Marshal(p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		b64EncodedReq := base64.StdEncoding.EncodeToString(reqJSON)
		t := tasks.Signature{
			Name: "encode_video",
			Args: []tasks.Arg{
				{
					Type:  "string",
					Value: b64EncodedReq,
				},
			},
		}
		res, err := taskserver.SendTask(&t)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		resp := struct {
			UUID string `json:"uuid"`
		}{
			UUID: res.GetState().TaskUUID,
		}
		resJSON, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(resJSON)
	})
	r.HandleFunc("/encode/multi", func(w http.ResponseWriter, r *http.Request) {
		p := new([]task.Payload)
		err := json.NewDecoder(r.Body).Decode(p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		t := make([]*tasks.Signature, len(*p))
		for i := range t {
			reqJSON, err := json.Marshal((*p)[i])
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			t[i] = &tasks.Signature{
				Name: "encode_video",
				Args: []tasks.Arg{
					{
						Type:  "string",
						Value: base64.StdEncoding.EncodeToString(reqJSON),
					},
				},
			}
		}
		group, err := tasks.NewGroup(t...)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res, err := taskserver.SendGroup(group, 1)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		resp := struct {
			UUID string `json:"uuid"`
		}{
			UUID: res[0].GetState().TaskUUID,
		}
		resJSON, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(resJSON)
	})
	r.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		task := r.FormValue("uuid")
		state, err := taskserver.GetBackend().GetState(task)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		resp := struct {
			State string `json:"state"`
		}{
			State: state.State,
		}
		resJSON, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(resJSON)
	})
	r.HandleFunc("/pending", func(w http.ResponseWriter, r *http.Request) {
		t, err := taskserver.GetBroker().GetPendingTasks("")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res, err := json.Marshal(t)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Write(res)
	})

	log.Println("listening on port 8082")
	log.Fatal(http.ListenAndServe(":8082", r))
}

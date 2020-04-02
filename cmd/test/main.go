package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"time"

	"github.com/dotmesh-io/terraform-provider-dotscience/pkg/types"
	"github.com/gorilla/mux"
)

func Terraform(args []string) {
	c := exec.Command("terraform", args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err := c.Run()
	if err != nil {
		fmt.Printf("----------------------------------------------------------\n")
		fmt.Printf("Error running %s\n", args)
		fmt.Printf("----------------------------------------------------------\n")
		os.Exit(1)
	}
}

func response(obj interface{}, statusCode int, err error, resp http.ResponseWriter, req *http.Request) {
	if err != nil {
		code := 500
		resp.WriteHeader(code)
		resp.Write([]byte(err.Error()))
		return
	}
	if obj != nil {
		resp.Header().Set("Content-Type", "application/json")
		resp.WriteHeader(statusCode)
		pr, pw := io.Pipe()
		go func() {
			pw.CloseWithError(json.NewEncoder(pw).Encode(obj))
		}()

		io.Copy(resp, pr)
	} else {
		resp.WriteHeader(statusCode)
	}
}

func main() {
	runners := []*types.Runner{
		&types.Runner{
			ID:          "testrunner",
			AccountID:   "testaccount",
			Name:        "Test Runner",
			Status:      "online",
			ServerState: "online",
			Tasks: []*types.Task{
				&types.Task{
					ID:            "1",
					DesiredStatus: "running",
					Status:        "running",
				},
			},
		},
	}

	deletedRunners := []*types.Runner{}

	requestLogger := func(handler func(resp http.ResponseWriter, req *http.Request)) func(resp http.ResponseWriter, req *http.Request) {
		return func(resp http.ResponseWriter, req *http.Request) {
			fmt.Printf("GOT request: %s %s\n", req.Method, req.URL)
			handler(resp, req)
		}
	}

	versionHandler := func(resp http.ResponseWriter, req *http.Request) {
		response("ok", 200, nil, resp, req)
	}

	listHandler := func(resp http.ResponseWriter, req *http.Request) {
		response(runners, 200, nil, resp, req)
	}

	actionHandler := func(resp http.ResponseWriter, req *http.Request) {
		runner := runners[0]
		task := runner.Tasks[0]
		task.DesiredStatus = "terminated"
		response("ok", 200, nil, resp, req)
		// set the task to terminated after 10 seconds
		go func() {
			time.Sleep(time.Second * 12)
			task.Status = "terminated"
		}()
	}

	deleteHandler := func(resp http.ResponseWriter, req *http.Request) {
		runner := runners[0]
		task := runner.Tasks[0]
		if task.Status != "terminated" {
			response("task is not terminated", 414, nil, resp, req)
			return
		}
		runner.Status = "deleting"
		runner.ServerState = "deleting"
		response("ok", 200, nil, resp, req)
		// remove the runner after 10 seconds to simulate pb removing the runner
		// in the background
		go func() {
			time.Sleep(time.Second * 12)
			deletedRunners = append(deletedRunners, runner)
			runners = []*types.Runner{}
		}()
	}

	r := mux.NewRouter()
	r.HandleFunc("/v2/version", requestLogger(versionHandler)).Methods("GET")
	r.HandleFunc("/admin/v1/runners", requestLogger(listHandler)).Methods("GET")
	r.HandleFunc("/admin/v1/runners/{account}/{id}", requestLogger(deleteHandler)).Methods("DELETE")
	r.HandleFunc("/admin/v1/runners/{account}/{id}/action", requestLogger(actionHandler)).Methods("POST")

	ts := httptest.NewServer(r)
	defer ts.Close()

	os.Setenv("TF_VAR_hub_public_url", ts.URL)
	os.Setenv("TF_VAR_hub_admin_username", "admin")
	os.Setenv("TF_VAR_hub_admin_password", "password")

	fmt.Printf("api server url: %s\n", ts.URL)

	os.Chdir("example")

	Terraform([]string{"init"})
	Terraform([]string{"apply", "-auto-approve"})
	Terraform([]string{"destroy", "-auto-approve"})

	if len(runners) != 0 {
		fmt.Printf("----------------------------------------------------------\n")
		fmt.Printf("Expected 0 runners got %d\n", len(runners))
		fmt.Printf("----------------------------------------------------------\n")
		os.Exit(1)
	}

	if len(deletedRunners) != 1 {
		fmt.Printf("----------------------------------------------------------\n")
		fmt.Printf("Expected 1 deleted runners got %d\n", len(deletedRunners))
		fmt.Printf("----------------------------------------------------------\n")
		os.Exit(1)
	}

	fmt.Printf("----------------------------------------------------------\n")
	fmt.Printf("test complete\n")
	fmt.Printf("----------------------------------------------------------\n")
	os.Exit(0)
}

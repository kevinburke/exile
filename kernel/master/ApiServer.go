package master

import (
	"encoding/json"
	"fmt"
	"github.com/m9rco/exile/kernel/common"
	. "github.com/m9rco/exile/kernel/utils"
	"net"
	"net/http"
	"os"
	"time"
)

// Job http service
type ApiServer struct {
	httpServer *http.Server
}

var (
	err error
	job common.Job
)

// Save the jobs
/*

POST
{
   name: "Job",
   command: "echo Job",
   cronExpr: '* * * * * '
}
 */
func handleJobSave(writer http.ResponseWriter, request *http.Request) {
	var (
		oldJob    *common.Job
		bytes     []byte
		jobManageSev JobManager
	)
	if err = request.ParseForm(); err != nil {
		goto ERROR
	}
	if err = json.Unmarshal([]byte(request.PostForm.Get("job")), &job); err != nil {
		goto ERROR
	}

	jobManageSev = common.Manage.GetSingleton("JobManager").(JobManager)
	if oldJob, err = jobManageSev.SaveJob(&job); err != nil {
		goto ERROR
	}
	if bytes, err = common.BuildResponse(0, "success", oldJob); err == nil {
		writer.Write(bytes)
	}
	return

ERROR:
	// return to the front anomalies. errno -1
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		writer.Write(bytes)
	}
	return
}

// Initialize the service
func InitApiServer() (err error) {
	var (
		mux             *http.ServeMux
		httpServer      *http.Server
		listener        net.Listener
		configureSource interface{}
	)
	if configureSource, err = common.Manage.GetPrototype("configure"); err != nil {
		fmt.Printf("fail to read file: %v", err)
		os.Exit(1)
	}
	configure := configureSource.(IniParser)
	// configure the routers
	mux = http.NewServeMux()
	mux.HandleFunc("/job/save", handleJobSave)

	// Start TCP listener
	if listener, err = net.Listen(
		configure.GetString("server", "protocol"), configure.GetString("server", "port"));
		err != nil {
		return
	}
	// Create http servers
	httpServer = &http.Server{
		Addr:              "",
		Handler:           mux,
		TLSConfig:         nil,
		ReadTimeout:       time.Duration(configure.GetInt64("server", "read_timeout")) * time.Millisecond,
		ReadHeaderTimeout: 0,
		WriteTimeout:      time.Duration(configure.GetInt32("server", "write_timeout")) * time.Millisecond,
		IdleTimeout:       0,
		MaxHeaderBytes:    0,
		TLSNextProto:      nil,
		ConnState:         nil,
		ErrorLog:          nil,
	}
	common.Manage.SetSingleton("ApiServer", &ApiServer{
		httpServer: httpServer,
	})
	go httpServer.Serve(listener)
	return
}
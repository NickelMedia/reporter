/*
   Copyright 2016 Vastech SA (PTY) LTD

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package main

import (
	"log"
	"net/http"
	"os"

	"github.com/IzakMarais/reporter/grafana"
	"github.com/IzakMarais/reporter/report"
	"github.com/gorilla/mux"
	"github.com/namsral/flag"
)

var proto = flag.String("proto", "http://", "Grafana Protocol")
var ip = flag.String("ip", "localhost:3000", "Grafana IP and port")
var port = flag.String("port", ":8686", "Port to serve on")
var templateDir = flag.String("templates", "templates/", "Directory for custom TeX templates")

var (
	dbHost = flag.String("dbHost", "", "Reporting metadata database host")
	dbPort = flag.String("dbPort", "", "Reporting metadata database port")
	username = flag.String("username", "", "Reporting metadata database username")
	password = flag.String("password", "", "Reporting metadata database password")
	database = flag.String("database", "", "Reporting metadata database name")
	queryStr = flag.String("query", "SELECT title, content FROM sections WHERE report_id IN (?) ORDER BY FIELD(title, 'Project', 'Overview', 'Targets', 'Method', 'Results');", "Reporting metadata query")
)

func main() {
	flag.Parse()
	log.SetOutput(os.Stdout)

	//'generated*'' variables injected from build.gradle: task 'injectGoVersion()'
	log.Printf("grafana reporter, version: %s.%s-%s hash: %s", generatedMajor, generatedMinor, generatedRelease, generatedGitHash)
	log.Printf("serving at '%s' and using grafana at '%s'", *port, *ip)
	if len(*dbHost) > 0 {
		log.Printf("using report metadata database at '%s'", *dbHost)
	}

	router := mux.NewRouter()
	RegisterHandlers(
		router,
		ServeReportHandler{grafana.NewV4Client, report.New},
		ServeReportHandler{grafana.NewV5Client, report.New},
	)

	log.Fatal(http.ListenAndServe(*port, router))
}

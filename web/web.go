// Copyright [2024] [Argus]
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package web provides the web server for Argus.
package web

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/release-argus/Argus/config"
	"github.com/release-argus/Argus/util"
	v1 "github.com/release-argus/Argus/web/api/v1"
)

var jLog *util.JLog

// NewRouter serves Prometheus metrics, WebSocket, and Node.js frontend at RoutePrefix.
func NewRouter(cfg *config.Config, hub *v1.Hub) *mux.Router {
	// Go
	api := v1.NewAPI(cfg, jLog)

	// Prometheus metrics
	api.Router.Handle("/metrics", promhttp.Handler())

	// WebSocket
	api.Router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Connection", "close")
		defer r.Body.Close()
		v1.ServeWs(hub, w, r)
	})

	// HTTP API
	api.SetupRoutesAPI()
	// Node.js
	api.SetupRoutesNodeJS()

	return api.BaseRouter
}

// newWebUI will set up everything web-related for Argus.
func newWebUI(cfg *config.Config) *mux.Router {
	hub := v1.NewHub()
	go hub.Run()
	router := NewRouter(cfg, hub)

	// Hand out the broadcast channel
	cfg.HardDefaults.Service.Status.AnnounceChannel = &hub.Broadcast
	for _, service := range cfg.Service {
		service.Status.SetAnnounceChannel(&hub.Broadcast)
	}

	return router
}

// Run the web server.
func Run(cfg *config.Config, log *util.JLog) {
	// Only set if unset (avoid RACE condition in tests).
	if log != nil && jLog == nil {
		jLog = log
	}

	router := newWebUI(cfg)

	listenAddress := fmt.Sprintf("%s:%s", cfg.Settings.WebListenHost(), cfg.Settings.WebListenPort())
	jLog.Info("Listening on "+listenAddress+cfg.Settings.WebRoutePrefix(), util.LogFrom{}, true)

	srv := &http.Server{
		Addr:         listenAddress,
		Handler:      router,
		ReadTimeout:  10 * time.Second, // Max time to read request headers and body.
		WriteTimeout: 10 * time.Second, // Max time to write response.
		IdleTimeout:  0,                // Disable to keep websocket connections open.
	}

	if cfg.Settings.WebCertFile() != "" && cfg.Settings.WebKeyFile() != "" {
		jLog.Fatal(
			srv.ListenAndServeTLS(cfg.Settings.WebCertFile(), cfg.Settings.WebKeyFile()),
			util.LogFrom{}, true)
	} else {
		jLog.Fatal(
			srv.ListenAndServe(),
			util.LogFrom{}, true)
	}
}

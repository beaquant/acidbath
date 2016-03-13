/*
   AcidBath - framework for your trading
   Copyright (C) 2016 Mark Laczynski

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/marklaczynski/acidbath/broker/factory"
	"github.com/marklaczynski/acidbath/lib/mjlog"

	"github.com/marklaczynski/acidbath/web/handlers"

	"github.com/gorilla/mux"
)

var (
	logInfo  = log.New(mjlog.CreateInfoFile(), "INFO  [main]: ", log.LstdFlags|log.Lshortfile)
	logDebug = log.New(mjlog.CreateDebugFile(), "DEBUG [main]: ", log.LstdFlags|log.Lshortfile)
	logError = log.New(mjlog.CreateErrorFile(), "ERROR [main]: ", log.LstdFlags|log.Lshortfile)
)

func main() {
	tdSession := factory.CreateBroker(factory.TD)

	logInfo.Printf("Starting up...\n")

	gmMux := mux.NewRouter()
	gmMux.Host("localhost:1111")

	// basic ui requests
	gmMux.HandleFunc("/", handlers.MakeHandler(handlers.RootHandler, tdSession))
	gmMux.HandleFunc("/login", handlers.MakeHandler(handlers.LoginHandler, tdSession))
	gmMux.HandleFunc("/logout", handlers.MakeHandler(handlers.LogoutHandler, tdSession))
	gmMux.HandleFunc("/reqOptChain", handlers.MakeHandler(handlers.ReqOptChainHandler, tdSession))
	gmMux.HandleFunc("/reqOrderBook", handlers.MakeHandler(handlers.ReqOrderBookHandler, tdSession))
	gmMux.HandleFunc("/trackOption", handlers.MakeHandler(handlers.TrackOptionHandler, tdSession))
	gmMux.HandleFunc("/untrackOption", handlers.MakeHandler(handlers.UntrackOptionHandler, tdSession))

	// event handlers that push data to ui
	gmMux.HandleFunc("/portfolioUpdateEvent", handlers.MakeHandler(handlers.PortfolioUpdateEvent, tdSession))
	gmMux.HandleFunc("/orderUpdateEvent", handlers.MakeHandler(handlers.OrderUpdateEvent, tdSession))
	gmMux.HandleFunc("/optionUpdateEvent", handlers.MakeHandler(handlers.OptionUpdateEvent, tdSession))

	// "under the covers" api
	gmMux.HandleFunc("/releaseOptionUpdatesEvents", handlers.MakeHandler(handlers.ReleaseOptionUpdatesEventsHandler, tdSession))

	// test phase
	gmMux.HandleFunc("/testOrderHandler", handlers.MakeHandler(handlers.TestOrderHandler, tdSession))
	gmMux.HandleFunc("/testCancelOrderHandler", handlers.MakeHandler(handlers.TestCancelOrderHandler, tdSession))

	//file handler
	gmMux.PathPrefix("/web/").Handler(http.StripPrefix("/web/", http.FileServer(http.Dir("./web/"))))

	fmt.Printf("Listening...\n")
	err := http.ListenAndServeTLS(":1111", "web/certificates/cert.pem", "web/certificates/key.pem", gmMux)
	if err != nil {
		logError.Printf("Error %s\n", err)
	}
}

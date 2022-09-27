/*
 * wars of warp - an implementation of warpwar
 *
 * Copyright (c) 2022 Michael D Henderson
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package cli

import (
	"context"
	"errors"
	"github.com/mdhender/wow/internal/june"
	"github.com/mdhender/wow/pkg/server"
	"github.com/spf13/cobra"
	"log"
	"net"
	"net/http"
	"time"
)

// cmdServer starts a server that returns an SVG map.
var cmdServer = &cobra.Command{
	Use:   "server",
	Short: "serve map api",
	Run: func(cmd *cobra.Command, args []string) {
		host, port := "", "8080"

		s, err := server.New()
		if err != nil {
			log.Fatal(err)
		}

		// server assumes that it is exposed to the internet.
		// it sets timeouts to avoid simple DOS attacks.
		// this may not be needed (or desired) if the server
		// hides behind a proxy like nginx.
		srv := &http.Server{
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  60 * time.Second,
			Addr:         net.JoinHostPort(host, port),
			Handler:      s.Routes(),
		}

		// start the server as a go routine.
		// this allows us to catch signals to shut it down gracefully.
		go func(srv *http.Server) {
			log.Printf("[server] address %q\n", srv.Addr)
			if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatal(err)
			}
			log.Println("[server] stopped")
		}(srv)

		// create the signal catchers, then wait on a signal.
		// the catchers don't need to know about the server;
		// they are only interested in trapping the signals.
		stopCh, closeCh := june.CreateChannel()
		defer closeCh()
		log.Println("[server] notified:", <-stopCh)

		// shut the server down via a context request.
		june.ShutdownServer(context.Background(), srv, 5*time.Second)

		// and we're done.
		log.Println("[server] terminating")
	},
}

func init() {
	cmdBase.AddCommand(cmdServer)
}

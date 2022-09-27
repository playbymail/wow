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

package server

import (
	"encoding/json"
	"errors"
	"github.com/mdhender/wow/pkg/board"
	"io"
	"net/http"
	"unicode/utf8"
)

// handlePostMapData accepts map data as CSV and returns an SVG or an error page.
func (s *Server) handlePostMapData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type errorObject struct {
			Code   int    `json:"code,omitempty"`
			Detail string `json:"detail,omitempty"`
		}
		type errResponse struct {
			Status string        `json:"status"`
			Errors []errorObject `json:"errors,omitempty"`
		}
		type okResponse struct {
			Status string      `json:"status"`
			Data   interface{} `json:"data"`
		}

		type node struct {
			Name      string   `json:"name"`
			Col       int      `json:"col"`
			Row       int      `json:"row"`
			EconValue int      `json:"econ-value"` // non-zero only if hasStar
			Warps     []string `json:"warps"`
		}

		var input struct {
			Mono  bool   `json:"mono,omitempty"`
			Nodes []node `json:"nodes,omitempty"`
		}

		contentType := r.Header.Get("Content-type")
		switch contentType {
		case "application/json":
			// enforce a maximum read of 10kb from the response body
			r.Body = http.MaxBytesReader(w, r.Body, 10*1024)
			// create a json decoder that will accept only our specific fields
			dec := json.NewDecoder(r.Body)
			dec.DisallowUnknownFields()
			if err := dec.Decode(&input); err != nil {
				response := errResponse{
					Status: "error",
					Errors: []errorObject{{
						Code:   http.StatusBadRequest,
						Detail: "invalid json object",
					}},
				}
				w.Header().Set("Content-Type", "application/vnd.api+json")
				w.WriteHeader(http.StatusOK)
				_ = json.NewEncoder(w).Encode(response)
				return
			}
			// call decode again to confirm that the request contained only a single JSON object
			if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
				response := errResponse{
					Status: "error",
					Errors: []errorObject{{
						Code:   http.StatusBadRequest,
						Detail: "request body must only contain a single json object",
					}},
				}
				w.Header().Set("Content-Type", "application/vnd.api+json")
				w.WriteHeader(http.StatusOK)
				_ = json.NewEncoder(w).Encode(response)
				return
			}
		case "application/x-www-form-urlencoded":
			if err := r.ParseForm(); err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			for k, v := range r.Form {
				switch k {
				case "mono":
					if len(v) != 1 || !utf8.ValidString(v[0]) {
						http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
						return
					}
					input.Mono = v[0] == "true"
				}
			}
		case "text/html":
			if err := r.ParseForm(); err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			for k, v := range r.Form {
				switch k {
				case "mono":
					if len(v) != 1 || !utf8.ValidString(v[0]) {
						http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
						return
					}
					input.Mono = v[0] == "true"
				}
			}
		default:
			http.Error(w, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
			return
		}

		if len(input.Nodes) == 0 {
			response := errResponse{
				Status: "error",
				Errors: []errorObject{{
					Code:   http.StatusBadRequest,
					Detail: "missing map data",
				}},
			}
			w.Header().Set("Content-Type", "application/vnd.api+json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(response)
			return
		} else if len(input.Nodes) > 40 {
			response := errResponse{
				Status: "error",
				Errors: []errorObject{{
					Code:   http.StatusBadRequest,
					Detail: "maximum number of nodes is 40",
				}},
			}
			w.Header().Set("Content-Type", "application/vnd.api+json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(response)
			return
		}

		// max row and col determine the size of the board
		minRow, maxRow, minCol, maxCol := 0, 0, 0, 0
		for i, h := range input.Nodes {
			if h.Row < minRow || i == 0 {
				minRow = h.Row
			}
			if h.Row > maxRow || i == 0 {
				maxRow = h.Row
			}
			if h.Col < minCol || i == 0 {
				minCol = h.Col
			}
			if h.Col > maxCol || i == 0 {
				maxCol = h.Col
			}
		}

		// sanity and performance checks
		if minCol < 1 || minRow < 1 {
			response := errResponse{
				Status: "error",
				Errors: []errorObject{{
					Code:   http.StatusBadRequest,
					Detail: "col and row must be at least 1",
				}},
			}
			w.Header().Set("Content-Type", "application/vnd.api+json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(response)
			return
		} else if maxCol > 40 || maxRow > 40 {
			response := errResponse{
				Status: "error",
				Errors: []errorObject{{
					Code:   http.StatusBadRequest,
					Detail: "col and row must be at least 1",
				}},
			}
			w.Header().Set("Content-Type", "application/vnd.api+json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(response)
			return
		}

		// create the board, add all the stars, then add the wormholes
		gb := board.NewBoard(maxRow, maxCol)
		for _, h := range input.Nodes {
			gb.AddStar(h.Name, h.Row, h.Col, h.EconValue)
		}
		for _, h := range input.Nodes {
			for _, target := range h.Warps {
				gb.AddWormHole(h.Name, target)
			}
		}

		// save the board as an SVG file
		w.Header().Set("content-type", "image/svg+xml")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(gb.AsSVG(input.Mono))
	}
}

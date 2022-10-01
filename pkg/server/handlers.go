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
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mdhender/wow/pkg/board"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode/utf8"
)

// handleIndex does that
func (s *Server) handleIndex(public string, maxCols, maxRows int) http.HandlerFunc {
	index, err := os.ReadFile(filepath.Join(public, "index.html"))
	if err != nil {
		log.Printf("[server] %+v\n", err)
		return func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(index)
	}
}

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
				case "data":
					if len(v) != 1 || !utf8.ValidString(v[0]) {
						http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
						return
					}
					r := csv.NewReader(strings.NewReader(v[0]))
					r.FieldsPerRecord = -1 // allow variable number of fields per line
					for {
						record, err := r.Read()
						if err != nil {
							break
						} else if len(record) < 5 {
							continue
						}
						n := node{
							Name:      strings.TrimSpace(record[0]),
							Col:       atoi(record[1]),
							Row:       atoi(record[2]),
							EconValue: atoi(record[3]),
						}
						for _, dest := range record[4:] {
							if dest = strings.TrimSpace(dest); len(dest) != 0 {
								n.Warps = append(n.Warps, dest)
							}
						}
						input.Nodes = append(input.Nodes, n)
					}
				case "fill-type":
					if len(v) != 1 || !utf8.ValidString(v[0]) {
						http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
						return
					}
					input.Mono = v[0] == "mono"
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
				if err := gb.AddWormHole(h.Name, target); err != nil {
					w.Header().Set("content-type", "text/html")
					w.WriteHeader(http.StatusOK)
					_, _ = w.Write([]byte(fmt.Sprintf(`<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8"><title>Wars of Warp</title></head><body><p>Sorry, but there was an error with the input</p><pre><code>%+v</code></pre>`, err)))
					return
				}
			}
		}

		// save the board as an SVG file
		w.Header().Set("content-type", "image/svg+xml")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(gb.AsSVG(input.Mono))
	}
}

// handleRandomMap does that
func (s *Server) handleRandomMap() http.HandlerFunc {
	type node struct {
		Name      string   `json:"name"`
		Col       int      `json:"col"`
		Row       int      `json:"row"`
		EconValue int      `json:"econ-value"` // non-zero only if hasStar
		Warps     []string `json:"warps"`
		distance  int
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// names returns a shuffled list of starry sounding names for stars.
		var names = []string{
			"Afak", "Agrab", "Akkad", "Al-Diniye", "Al-Esotam", "Al-Hafriyat", "Annah", "Arbela", "Arbīl", "Arrapkha",
			"Ashur", "Assur", "Athína", "Awan", "Babil", "Babylon", "Baghdad", "Borsippa", "Corinth", "Kurigalzu", "El-Ana",
			"El-Is", "En-Aasar", "En-Amitat", "En-Shubat", "Erech", "Erétria", "Eshnunna", "Gubba", "Hafriyat", "Haradum",
			"Hillah", "Kassite", "Khirbit", "Khūzestān", "Kirkūk", "Kutha", "Kórinthos", "Lagash", "Mari", "Mashkan", "Nagar",
			"Neribtum", "Nimrud", "Nineveh", "Nippur", "Nuffar", "Nuzi", "Opis", "Ramad", "Rapiqum", "Riblah", "Ródos",
			"Shaduppum", "Shapir", "Shushan", "Shūsh", "Sippar", "Siracusa", "Sirpurla", "Sparta", "Spárti", "Susa", "Tayma",
			"Te Ashyia", "Te Brak", "Te Ishchali", "Te Leilan", "Thebes", "Thíva", "Tuttul", "Tutub", "Umm", "Uqair", "Ur",
			"Urhai", "Urkesh", "Uruk", "Árgos", "Égina", "Şanlıurfa",
		}
		rand.Shuffle(len(names), func(i, j int) {
			names[i], names[j] = names[j], names[i]
		})

		var baseMap [22][22]*node
		for col := 1; col <= 20; col++ {
			for row := 1; row <= 20; row++ {
				// each hex has a 1 in 12 chance of containing a star
				if rand.Intn(12) != 1 {
					continue
				}
				// can't have a neighbor
				if baseMap[col-1][row-1] != nil {
					continue
				} else if baseMap[col-1][row] != nil {
					continue
				} else if baseMap[col-1][row+1] != nil {
					continue
				} else if baseMap[col][row-1] != nil {
					continue
				} else if baseMap[col][row+1] != nil {
					continue
				} else if baseMap[col+1][row-1] != nil {
					continue
				} else if baseMap[col+1][row] != nil {
					continue
				} else if baseMap[col+1][row+1] != nil {
					continue
				}
				// a good range for econ values is 0..5 with higher values being rarer
				var econValue int
				switch rand.Intn(23) {
				case 0:
					econValue = 5
				case 1, 2:
					econValue = 4
				case 3, 4, 5, 6:
					econValue = 3
				case 7, 8, 9, 10, 11, 12:
					econValue = 2
				case 13, 14, 15, 16, 17, 18, 19, 20:
					econValue = 1
				default:
					econValue = 0
				}
				var name string
				if len(names) == 0 {
					name = fmt.Sprintf("N%02d%02d", col, row)
				} else {
					name, names = names[0], names[1:]
				}
				baseMap[col][row] = &node{Name: name, Col: col, Row: row, EconValue: econValue}
			}
		}

		var nodes []*node
		for col := 1; col <= 20; col++ {
			for row := 1; row <= 20; row++ {
				if baseMap[col][row] != nil {
					nodes = append(nodes, baseMap[col][row])
				}
			}
		}

		// each star has a 1 in 4 chance of having a warp to each of the 4 nearest stars
		for _, n := range nodes {
			// fetch eligible stars (can't already have 4 warps out) and sort them by distance
			var neighbors []*node
			for _, x := range nodes {
				if x != n {
					x.distance = (n.Col-x.Col)*(n.Col-x.Col) + (n.Row-x.Row)*(n.Row-x.Row)
					neighbors = append(neighbors, x)
				}
			}
			// sort the neighbors by distance
			for i := 0; i < len(neighbors); i++ {
				for j := i + 1; j < len(neighbors); j++ {
					if neighbors[i].distance > neighbors[j].distance {
						neighbors[i], neighbors[j] = neighbors[j], neighbors[i]
					}
				}
			}
			// then check up to the first four neighbors warp lines
			for i := 0; i < 4 && i < len(neighbors) && len(n.Warps) < 4; i++ {
				// 1 in 4 chance of having a warp to this neighbor
				if len(n.Warps) < 4 && len(neighbors[i].Warps) < 4 && rand.Intn(4) == 1 {
					n.Warps = append(n.Warps, neighbors[i].Name)
					neighbors[i].Warps = append(neighbors[i].Warps, n.Name)
				}
			}

			if len(n.Warps) == 0 {
				// all stars must have a warp. if we didn't get one above, force it.
				for _, x := range neighbors {
					if len(x.Warps) < 4 {
						n.Warps = append(n.Warps, x.Name)
						x.Warps = append(x.Warps, n.Name)
						break
					}
				}
			}
		}

		// create the board, add all the stars, then add the wormholes

		// board will always be 20 x 20
		gb := board.NewBoard(20, 20)

		// add stars
		for _, n := range nodes {
			gb.AddStar(n.Name, n.Row, n.Col, n.EconValue)
		}

		// add wormholes
		for _, n := range nodes {
			for _, target := range n.Warps {
				if err := gb.AddWormHole(n.Name, target); err != nil {
					w.Header().Set("content-type", "text/html")
					w.WriteHeader(http.StatusOK)
					_, _ = w.Write([]byte(fmt.Sprintf(`<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8"><title>Wars of Warp</title></head><body><p>Sorry, but there was an error with the input</p><pre><code>%+v</code></pre>`, err)))
					return
				}
			}
		}

		// save the board as an SVG file
		w.Header().Set("content-type", "image/svg+xml")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(gb.AsSVG(true))
	}
}

// handleStandardMap does that
func (s *Server) handleStandardMap(static string, color bool) http.HandlerFunc {
	var filename string
	if color {
		filename = filepath.Join(static, "svg-test.svg")
	} else {
		filename = filepath.Join(static, "svg-mono-test.svg")
	}
	index, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("[server] %+v\n", err)
		return func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "image/svg+xml")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(index)
	}
}

func atoi(s string) int {
	i, _ := strconv.Atoi(strings.TrimSpace(s))
	return i
}

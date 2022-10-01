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

// Package board implements a hexagonal game board.
package board

// see https://www.redblobgames.com/grids/hexagons/ for background and overview

import (
	"bytes"
	"fmt"
	"github.com/mdhender/wow/pkg/hexes"
	"math"
)

func NewBoard(rows, cols int) *Board {
	rows, cols = rows+1, cols+1
	b := &Board{
		Rows:  rows + 1,
		Cols:  cols + 1,
		Stars: make(map[string]*Hex),
	}

	b.Hexes = make([][]*Hex, b.Rows)
	for row := 0; row < b.Rows; row++ {
		b.Hexes[row] = make([]*Hex, b.Cols)
		for col := 0; col < b.Cols; col++ {
			b.Hexes[row][col] = &Hex{Coords: Coords{Row: row, Col: col}}
		}
	}

	return b
}

func (b *Board) AddStar(name string, row, col int, econValue int) {
	hex := &Hex{
		Coords:    Coords{Row: row, Col: col},
		Name:      name,
		HasStar:   true,
		EconValue: econValue,
	}
	b.Hexes[row][col] = hex
	b.Stars[name] = hex
}

func (b *Board) AddWormHole(sourceStar, targetStar string) error {
	// lookup both ends of the wormhole
	from, ok := b.Stars[sourceStar]
	if !ok {
		return fmt.Errorf("board: invalid source star: %q", sourceStar)
	}
	to, ok := b.Stars[targetStar]
	if !ok {
		return fmt.Errorf("board: invalid target star: %q", targetStar)
	}

	// add to exits if not there already
	from.AddWormHole(to)
	to.AddWormHole(from)

	return nil
}

func (b *Board) AsHTML(mono bool) []byte {
	// create the svg for the board
	buf := &bytes.Buffer{}
	_, _ = fmt.Fprintln(buf, `<!doctype html>`)
	_, _ = fmt.Fprintln(buf, `<html lang="en">`)
	_, _ = fmt.Fprintln(buf, `<head>`)
	_, _ = fmt.Fprintln(buf, `<meta charset="utf-8">`)
	_, _ = fmt.Fprintln(buf, `<title>SVG Test</title>`)
	//_, _ = fmt.Fprintln(b, `<style>div.scroll {background-color: #fed9ff;width: 95%;height: 95%;overflow: auto;text-align: justify;padding: 1%;}</style>`)
	_, _ = fmt.Fprintf(buf, `<style>`)
	_, _ = fmt.Fprintf(buf, `svg{background-color:hsl(197, 18%%, 95%%);padding:50px 50px 50px 50px;}`)
	//_, _ = fmt.Fprintf(buf, `rect{fill:red;stroke:none;shape-rendering:crispEdges;}`)
	_, _ = fmt.Fprintf(buf, `</style>`)
	_, _ = fmt.Fprintln(buf, `</head>`)
	_, _ = fmt.Fprintln(buf, `<body>`)
	//_, _ = fmt.Fprintln(b, `<div class="scroll">`)
	_, _ = fmt.Fprintln(buf, b.asSVG(mono).String())
	//_, _ = fmt.Fprintln(b, `</div>`)
	_, _ = fmt.Fprintln(buf, "</body>")
	_, _ = fmt.Fprintln(buf, "</html>")

	return buf.Bytes()
}

func (b *Board) AsSVG(mono bool) []byte {
	return []byte(b.asSVG(mono).String())
}

func (b *Board) asSVG(mono bool) *svg {
	size := 55.0
	width, height := 2*size, math.Sqrt(3)*size
	layout := hexes.NewFlatLayout(hexes.NewPoint(size, size), hexes.NewPoint(height, width))

	// "hsl(39, 100%, 50%)" // "LightBlue" // "hsl(197, 78%, 85%)"
	var hexFill, starFill string
	if mono {
		hexFill, starFill = "none", "White"
	} else {
		hexFill, starFill = "hsl(197, 78%, 85%)", "hsl(53, 100%, 94%)"
	}

	// svg has 0,0 in the upper left.
	s := &svg{id: "s"}

	// create the hexes
	for row := 0; row < b.Rows; row++ {
		for col := 0; col < b.Cols; col++ {
			// assumes flat with even-q layout
			h := hexes.QOffsetToCube(col, row, hexes.EVEN)

			cx, cy := layout.CenterPoint(h).Coords()
			poly := &polygon{col: col, row: row, cx: cx, cy: cy, radius: height / 2.0}

			poly.style.stroke = "Grey"
			poly.style.fill = hexFill
			if poly.style.fill == poly.style.stroke {
				poly.style.stroke = "Black"
			}
			poly.style.strokeWidth = "2px"

			for _, p := range layout.PolygonCorners(h) {
				px, py := p.Coords()
				if width := int(px); width > s.viewBox.width {
					s.viewBox.width = width
				}
				if height := int(py); height > s.viewBox.height {
					s.viewBox.height = height
				}
				poly.points = append(poly.points, point{x: px, y: py})
			}

			s.hexes = append(s.hexes, poly)
		}
	}

	// create the stars
	for row := 0; row < b.Rows; row++ {
		for col := 0; col < b.Cols; col++ {
			hex := b.Hexes[row][col]
			if hex.Name == "" {
				continue // not a star
			}

			// assumes flat with even-q layout
			h := hexes.QOffsetToCube(col, row, hexes.EVEN)
			// qq, rr, ss := h.Coords()

			cx, cy := layout.CenterPoint(h).Coords()
			poly := &polygon{cx: cx, cy: cy, radius: height / 2.0}

			poly.text = []string{hex.Name, fmt.Sprintf("( %d )", hex.EconValue)}
			poly.addCircle = true

			poly.style.stroke = "Grey"
			poly.style.fill = starFill
			if poly.style.fill == poly.style.stroke {
				poly.style.stroke = "Black"
			}
			poly.style.strokeWidth = "2px"

			for _, p := range layout.PolygonCorners(h) {
				px, py := p.Coords()
				if width := int(px); width > s.viewBox.width {
					s.viewBox.width = width
				}
				if height := int(py); height > s.viewBox.height {
					s.viewBox.height = height
				}
				poly.points = append(poly.points, point{x: px, y: py})
			}

			s.polygons = append(s.polygons, poly)

			for _, star := range hex.WormHoleExits {
				sx, sy := layout.CenterPoint(hexes.QOffsetToCube(star.Coords.Col, star.Coords.Row, hexes.EVEN)).Coords()
				s.lines = append(s.lines, [4]float64{cx, cy, sx, sy})
			}
		}
	}

	return s
}

type Board struct {
	Rows, Cols int
	Hexes      [][]*Hex
	Stars      map[string]*Hex
}

// AddWormHole adds a new exit to the hex.
// Caller should call this for both ends of the wormhole.
func (h *Hex) AddWormHole(to *Hex) {
	for _, hex := range h.WormHoleExits {
		if hex == to {
			return
		}
	}
	h.WormHoleExits = append(h.WormHoleExits, to)
}

// Fill returns the color that should be used to fill the hex on the map
func (h *Hex) Fill() string {
	if h.Name == "" {
		return "hsl(197, 78%, 85%)"
	}
	return "hsl(53, 100%, 94%)"
}

type Coords struct {
	Col, Row int
}

func (c Coords) Less(d Coords) bool {
	if c.Row < d.Row {
		return true
	} else if c.Row > d.Row {
		return false
	}
	return c.Col < d.Col
}

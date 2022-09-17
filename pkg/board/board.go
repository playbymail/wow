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

func (b *Board) AddWormHole(sourceStar, targetStar string) {
	// lookup both ends of the wormhole
	from, ok := b.Stars[sourceStar]
	if !ok {
		panic(fmt.Sprintf("board: invalid source star: %q", sourceStar))
	}
	to, ok := b.Stars[targetStar]
	if !ok {
		panic(fmt.Sprintf("board: invalid target star: %q", targetStar))
	}

	// add to exits if not there already
	from.AddWormHole(to)
	to.AddWormHole(from)
}

func (b *Board) AsHTML() []byte {
	// create the svg for the board
	buf := &bytes.Buffer{}
	_, _ = fmt.Fprintln(buf, `<!doctype html>`)
	_, _ = fmt.Fprintln(buf, `<html lang="en">`)
	_, _ = fmt.Fprintln(buf, `<head>`)
	_, _ = fmt.Fprintln(buf, `<meta charset="utf-8">`)
	_, _ = fmt.Fprintln(buf, `<title>SVG Test</title>`)
	//_, _ = fmt.Fprintln(b, `<style>div.scroll {background-color: #fed9ff;width: 95%;height: 95%;overflow: auto;text-align: justify;padding: 1%;}</style>`)
	_, _ = fmt.Fprintln(buf, `</head>`)
	_, _ = fmt.Fprintln(buf, `<body>`)
	//_, _ = fmt.Fprintln(b, `<div class="scroll">`)
	_, _ = fmt.Fprintln(buf, b.asSVG().String())
	//_, _ = fmt.Fprintln(b, `</div>`)
	_, _ = fmt.Fprintln(buf, "</body>")
	_, _ = fmt.Fprintln(buf, "</html>")

	return buf.Bytes()
}

func (b *Board) asSVG() *svg {
	radius := 30.0 * 1.5 // mdhender: scaled for star names
	offset := (math.Sqrt(3) * radius) / 2

	// why 40.0 here?
	maxX := 40.0 + offset*float64(b.Cols*2)
	maxY := 40.0 + offset*float64(b.Rows)*math.Sqrt(3)

	s := &svg{}
	s.id = "s"
	s.viewBox.minX, s.viewBox.width = 0, int(maxX+radius/2)
	s.viewBox.minY, s.viewBox.height = 0, int(maxY+radius/2)

	for row := 0; row < b.Rows; row++ {
		for col := 0; col < b.Cols; col++ {
			hex := b.Hexes[row][col]

			x := 40.0 + offset*float64(hex.Coords.Col*2)
			if hex.Coords.Row%2 == 0 {
				x += offset
			}
			y := 40.0 + offset*float64(hex.Coords.Row)*math.Sqrt(3)

			// the board has 0,0 in the lower left but svg has 0,0 in the upper left.
			// we have to change y from [0..maxY] to [maxY..0].
			y = maxY - y

			poly := &polygon{x: x, y: y, radius: radius}
			if hex.Name == "" {
				poly.label = fmt.Sprintf("%02d%02d", hex.Coords.Row, hex.Coords.Col)
			} else {
				poly.label = fmt.Sprintf("%s (%d)", hex.Name, hex.EconValue)
			}
			poly.style.fill = hex.Fill()
			poly.style.stroke = "LightGrey"
			poly.style.stroke = "Grey"
			if poly.style.fill == poly.style.stroke {
				poly.style.stroke = "Black"
			}
			poly.style.strokeWidth = "2px"
			for _, p := range poly.hexPoints() {
				poly.points = append(poly.points, point{x: p.x, y: p.y})
			}

			s.polygons = append(s.polygons, poly)
		}
	}

	return s
}

type Board struct {
	Rows, Cols int
	Hexes      [][]*Hex
	Stars      map[string]*Hex
}

type Hex struct {
	Coords        Coords
	Name          string
	HasStar       bool
	EconValue     int
	WormHoleExits []*Hex
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

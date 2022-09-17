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

// Package hex implements https://www.redblobgames.com/grids/hexagons/codegen/output/lib.cpp
package hex

import (
	"math"
)

type Point struct {
	x, y float64
}

func NewPoint(x, y float64) Point {
	return Point{x: x, y: y}
}

type Hex struct {
	q, r, s int
}

func NewHex(q, r, s int) Hex {
	if q+r+s != 0 {
		panic("assert(q + r + s == 0)")
	}
	return Hex{q: q, r: r, s: s}
}

type FractionalHex struct {
	q, r, s float64
}

func NewFractionalHex(q, r, s float64) FractionalHex {
	if math.Round(q+r+s) != 0 {
		panic("assert(q + r + s == 0)")
	}
	return FractionalHex{q: q, r: r, s: s}
}

type OffsetCoord struct {
	col, row int
}

func NewOffsetCoord(col, row int) OffsetCoord {
	return OffsetCoord{col: col, row: row}
}

type DoubledCoord struct {
	col, row int
}

func NewDoubledCoord(col, row int) DoubledCoord {
	return DoubledCoord{col: col, row: row}
}

type Orientation struct {
	f0, f1, f2, f3 float64
	b0, b1, b2, b3 float64
	start_angle    float64
}

func NewOrientation(f0, f1, f2, f3, b0, b1, b2, b3, start_angle float64) Orientation {
	return Orientation{
		f0: f0, f1: f1, f2: f2, f3: f3,
		b0: b0, b1: b1, b2: b2, b3: b3,
		start_angle: start_angle,
	}
}

type Layout struct {
	orientation  Orientation
	size, origin Point
}

func NewLayout(orientation Orientation, size, origin Point) Layout {
	return Layout{orientation: orientation, size: size, origin: origin}
}

func hex_add(a, b Hex) Hex {
	return NewHex(a.q+b.q, a.r+b.r, a.s+b.s)
}

func hex_subtract(a, b Hex) Hex {
	return NewHex(a.q-b.q, a.r-b.r, a.s-b.s)
}

func hex_scale(a Hex, k int) Hex {
	return NewHex(a.q*k, a.r*k, a.s*k)
}

func hex_rotate_left(a Hex) Hex {
	return NewHex(-a.s, -a.q, -a.r)
}

func hex_rotate_right(a Hex) Hex {
	return NewHex(-a.r, -a.s, -a.q)
}

var hex_directions = []Hex{
	NewHex(1, 0, -1),
	NewHex(1, -1, 0),
	NewHex(0, -1, 1),
	NewHex(-1, 0, 1),
	NewHex(-1, 1, 0),
	NewHex(0, 1, -1),
}

func hex_direction(direction int) Hex {
	return hex_directions[direction]
}

func hex_neighbor(hex Hex, direction int) Hex {
	return hex_add(hex, hex_direction(direction))
}

var hex_diagonals = []Hex{
	NewHex(2, -1, -1),
	NewHex(1, -2, 1),
	NewHex(-1, -1, 2),
	NewHex(-2, 1, 1),
	NewHex(-1, 2, -1),
	NewHex(1, 1, -2),
}

func hex_diagonal_neighbor(hex Hex, direction int) Hex {
	return hex_add(hex, hex_diagonals[direction])
}

// abs is a helper function to get the absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func hex_length(hex Hex) int {
	return (abs(hex.q) + abs(hex.r) + abs(hex.s)) / 2
}

func hex_distance(a, b Hex) int {
	return hex_length(hex_subtract(a, b))
}

func hex_round(h FractionalHex) Hex {
	qi := int(math.Round(h.q))
	q_diff := math.Abs(float64(qi) - h.q)
	ri := int(math.Round(h.r))
	r_diff := math.Abs(float64(ri) - h.r)
	si := int(math.Round(h.s))
	s_diff := math.Abs(float64(si) - h.s)

	if q_diff > r_diff && q_diff > s_diff {
		qi = -ri - si
	} else if r_diff > s_diff {
		ri = -qi - si
	} else {
		si = -qi - ri
	}

	return NewHex(qi, ri, si)
}

func hex_lerp(a, b FractionalHex, t float64) FractionalHex {
	return NewFractionalHex(a.q*(1.0-t)+b.q*t, a.r*(1.0-t)+b.r*t, a.s*(1.0-t)+b.s*t)
}

func hex_linedraw(a, b Hex) (results []Hex) {
	N := hex_distance(a, b)
	a_nudge := NewFractionalHex(float64(a.q)+1e-06, float64(a.r)+1e-06, float64(a.s)-2e-06)
	b_nudge := NewFractionalHex(float64(b.q)+1e-06, float64(b.r)+1e-06, float64(b.s)-2e-06)
	step := 1.0 / math.Max(float64(N), 1.0)

	for i := 0; i <= N; i++ {
		results = append(results, hex_round(hex_lerp(a_nudge, b_nudge, step*float64(i))))
	}

	return results
}

type OFFSET int

const (
	EVEN OFFSET = 1
	ODD  OFFSET = -1
)

func qoffset_from_cube(offset OFFSET, h Hex) OffsetCoord {
	col := h.q
	row := h.r + (h.q+int(offset)*(h.q&1))/2

	return NewOffsetCoord(col, row)
}

func qoffset_to_cube(offset OFFSET, h OffsetCoord) Hex {
	q := h.col
	r := h.row - (h.col+int(offset)*(h.col&1))/2
	s := -q - r

	return NewHex(q, r, s)
}

func roffset_from_cube(offset OFFSET, h Hex) OffsetCoord {
	col := h.q + (h.r+int(offset)*(h.r&1))/2
	row := h.r

	return NewOffsetCoord(col, row)
}

func roffset_to_cube(offset OFFSET, h OffsetCoord) Hex {
	q := h.col - ((h.row + int(offset)*(h.row&1)) / 2)
	r := h.row
	s := -q - r

	return NewHex(q, r, s)
}

func qdoubled_from_cube(h Hex) DoubledCoord {
	col := h.q
	row := 2*h.r + h.q

	return NewDoubledCoord(col, row)
}

func qdoubled_to_cube(h DoubledCoord) Hex {
	q := h.col
	r := (h.row - h.col) / 2
	s := -q - r

	return NewHex(q, r, s)
}

func rdoubled_from_cube(h Hex) DoubledCoord {
	col := 2*h.q + h.r
	row := h.r

	return NewDoubledCoord(col, row)
}

func rdoubled_to_cube(h DoubledCoord) Hex {
	q := (h.col - h.row) / 2
	r := h.row
	s := -q - r

	return NewHex(q, r, s)
}

var layout_pointy = NewOrientation(math.Sqrt(3.0), math.Sqrt(3.0)/2.0, 0.0, 3.0/2.0, math.Sqrt(3.0)/3.0, -1.0/3.0, 0.0, 2.0/3.0, 0.5)

var layout_flat = NewOrientation(3.0/2.0, 0.0, math.Sqrt(3.0)/2.0, math.Sqrt(3.0), 2.0/3.0, 0.0, -1.0/3.0, math.Sqrt(3.0)/3.0, 0.0)

func hex_to_pixel(layout Layout, h Hex) Point {
	M := layout.orientation
	size := layout.size
	origin := layout.origin

	x := (M.f0*float64(h.q) + M.f1*float64(h.r)) * size.x
	y := (M.f2*float64(h.q) + M.f3*float64(h.r)) * size.y

	return NewPoint(x+origin.x, y+origin.y)
}

func pixel_to_hex(layout Layout, p Point) FractionalHex {
	M := layout.orientation
	size := layout.size
	origin := layout.origin

	pt := NewPoint((p.x-origin.x)/size.x, (p.y-origin.y)/size.y)

	q := M.b0*pt.x + M.b1*pt.y
	r := M.b2*pt.x + M.b3*pt.y

	return NewFractionalHex(q, r, -q-r)
}

func hex_corner_offset(layout Layout, corner int) Point {
	M := layout.orientation
	size := layout.size

	angle := 2.0 * math.Pi * (M.start_angle - float64(corner)) / 6.0

	return NewPoint(size.x*math.Cos(angle), size.y*math.Sin(angle))
}

func polygon_corners(layout Layout, h Hex) (corners []Point) {
	center := hex_to_pixel(layout, h)
	for i := 0; i < 6; i++ {
		offset := hex_corner_offset(layout, i)
		corners = append(corners, NewPoint(center.x+offset.x, center.y+offset.y))
	}

	return corners
}

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

// Package hexes implements a hex grid library.
// Lifted almost as-is from https://www.redblobgames.com/grids/hexagons/codegen/output/lib.cpp
package hexes

import (
	"math"
)

type Hex struct {
	q, r, s int
}

func NewHex(q, r, s int) Hex {
	if q+r+s != 0 {
		panic("assert(q + r + s == 0)")
	}
	return Hex{q: q, r: r, s: s}
}

func (h Hex) Add(b Hex) Hex {
	return NewHex(h.q+b.q, h.r+b.r, h.s+b.s)
}

func (h Hex) DiagonalNeighbor(direction int) Hex {
	// direction = mod(direction, 6)
	direction = (6 + (direction % 6)) % 6
	return h.Add(hex_diagonals[direction])
}

func (h Hex) Distance(b Hex) int {
	return h.Subtract(b).Length()
}

func (h Hex) Equals(b Hex) bool {
	return h.q == b.q && h.s == b.s && h.r == b.r
}

func (h Hex) Length() int {
	return (abs(h.q) + abs(h.r) + abs(h.s)) / 2
}

func (h Hex) LineDraw(b Hex) (results []Hex) {
	N := h.Distance(b)

	a_nudge := NewFractionalHex(float64(h.q)+1e-06, float64(h.r)+1e-06, float64(h.s)-2e-06)
	b_nudge := NewFractionalHex(float64(b.q)+1e-06, float64(b.r)+1e-06, float64(b.s)-2e-06)
	step := 1.0 / math.Max(float64(N), 1.0)

	for i := 0; i <= N; i++ {
		results = append(results, a_nudge.Lerp(b_nudge, step*float64(i)).Round())
	}

	return results
}

func (h Hex) Multiply(k int) Hex {
	return NewHex(h.q*k, h.r*k, h.s*k)
}

func (h Hex) Neighbor(direction int) Hex {
	return h.Add(hex_direction(direction))
}

func (h Hex) RotateLeft() Hex {
	return NewHex(-h.s, -h.q, -h.r)
}

func (h Hex) RotateRight() Hex {
	return NewHex(-h.r, -h.s, -h.q)
}

func (h Hex) Scale(k int) Hex {
	return NewHex(h.q*k, h.r*k, h.s*k)
}

func (h Hex) Subtract(b Hex) Hex {
	return NewHex(h.q-b.q, h.r-b.r, h.s-b.s)
}

var hex_directions = []Hex{
	NewHex(1, 0, -1),
	NewHex(1, -1, 0),
	NewHex(0, -1, 1),
	NewHex(-1, 0, 1),
	NewHex(-1, 1, 0),
	NewHex(0, 1, -1),
}

var hex_diagonals = []Hex{
	NewHex(2, -1, -1),
	NewHex(1, -2, 1),
	NewHex(-1, -1, 2),
	NewHex(-2, 1, 1),
	NewHex(-1, 2, -1),
	NewHex(1, 1, -2),
}

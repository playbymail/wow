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

package hex

// this file implements the tests from https://www.redblobgames.com/grids/hexagons/codegen/output/lib.cpp

import "testing"

func TestHexArithmetic(t *testing.T) {
	a := NewHex(4, -10, 6)
	if !a.Equals(hex_add(NewHex(1, -3, 2), NewHex(3, -7, 4))) {
		t.Error("hex_add")
	}

	b := NewHex(-2, 4, -2)
	if !b.Equals(hex_subtract(NewHex(1, -3, 2), NewHex(3, -7, 4))) {
		t.Error("hex_subtract")
	}
}

func TestHexDirection(t *testing.T) {
	a := NewHex(0, -1, 1)
	if !a.Equals(hex_direction(2)) {
		t.Error("hex_direction")
	}
}

func TestHexNeighbor(t *testing.T) {
	a := NewHex(1, -3, 2)
	if !a.Equals(hex_neighbor(NewHex(1, -2, 1), 2)) {
		t.Error("hex_neighbor")
	}
}

func TestHexDiagonal(t *testing.T) {
	a := NewHex(-1, -1, 2)
	if !a.Equals(hex_diagonal_neighbor(NewHex(1, -2, 1), 3)) {
		t.Error("hex_diagonal")
	}
}

func TestHexDistance(t *testing.T) {
	if hex_distance(NewHex(3, -7, 4), NewHex(0, 0, 0)) != 7 {
		t.Error("hex_distance")
	}
}

func TestHexRotateRight(t *testing.T) {
	a := hex_rotate_right(NewHex(1, -3, 2))
	if !a.Equals(NewHex(3, -2, -1)) {
		t.Error("hex_rotate_right")
	}
}

func TestHexRotateLeft(t *testing.T) {
	a := hex_rotate_left(NewHex(1, -3, 2))
	if !a.Equals(NewHex(-2, -1, 3)) {
		t.Error("hex_rotate_left")
	}
}

func TestHexRound(t *testing.T) {
	a := NewFractionalHex(0.0, 0.0, 0.0)
	b := NewFractionalHex(1.0, -1.0, 0.0)
	c := NewFractionalHex(0.0, -1.0, 1.0)
	equal_hex(t, "hex_round 1", NewHex(5, -10, 5), hex_round(hex_lerp(NewFractionalHex(0.0, 0.0, 0.0), NewFractionalHex(10.0, -20.0, 10.0), 0.5)))
	equal_hex(t, "hex_round 2", hex_round(a), hex_round(hex_lerp(a, b, 0.499)))
	equal_hex(t, "hex_round 3", hex_round(b), hex_round(hex_lerp(a, b, 0.501)))
	equal_hex(t, "hex_round 4", hex_round(a), hex_round(NewFractionalHex(a.q*0.4+b.q*0.3+c.q*0.3, a.r*0.4+b.r*0.3+c.r*0.3, a.s*0.4+b.s*0.3+c.s*0.3)))
	equal_hex(t, "hex_round 5", hex_round(c), hex_round(NewFractionalHex(a.q*0.3+b.q*0.3+c.q*0.4, a.r*0.3+b.r*0.3+c.r*0.4, a.s*0.3+b.s*0.3+c.s*0.4)))
}

func TestHexLinedraw(t *testing.T) {
	equal_hex_array(t, "hex_linedraw", []Hex{NewHex(0, 0, 0), NewHex(0, -1, 1), NewHex(0, -2, 2), NewHex(1, -3, 2), NewHex(1, -4, 3), NewHex(1, -5, 4)}, hex_linedraw(NewHex(0, 0, 0), NewHex(1, -5, 4)))
}

func TestLayout(t *testing.T) {
	h := NewHex(3, 4, -7)

	flat := NewLayout(layout_flat, NewPoint(10.0, 15.0), NewPoint(35.0, 71.0))
	if !h.Equals(hex_round(pixel_to_hex(flat, hex_to_pixel(flat, h)))) {
		t.Error("layout")
	}

	pointy := NewLayout(layout_pointy, NewPoint(10.0, 15.0), NewPoint(35.0, 71.0))
	if !h.Equals(hex_round(pixel_to_hex(pointy, hex_to_pixel(pointy, h)))) {
		t.Error("layout")
	}
}

func TestOffsetRoundtrip(t *testing.T) {
	a := NewHex(3, 4, -7)
	if !a.Equals(qoffset_to_cube(EVEN, qoffset_from_cube(EVEN, a))) {
		t.Error("conversion_from_to even-q")
	}
	if !a.Equals(qoffset_to_cube(ODD, qoffset_from_cube(ODD, a))) {
		t.Error("conversion_from_to odd-q")
	}
	if !a.Equals(roffset_to_cube(EVEN, roffset_from_cube(EVEN, a))) {
		t.Error("conversion_from_to even-r")
	}
	if !a.Equals(roffset_to_cube(ODD, roffset_from_cube(ODD, a))) {
		t.Error("conversion_from_to odd-r")
	}

	b := NewOffsetCoord(1, -3)
	if !b.Equals(qoffset_from_cube(EVEN, qoffset_to_cube(EVEN, b))) {
		t.Error("conversion_to_from even-q")
	}
	if !b.Equals(qoffset_from_cube(ODD, qoffset_to_cube(ODD, b))) {
		t.Error("conversion_to_from odd-q")
	}
	if !b.Equals(roffset_from_cube(EVEN, roffset_to_cube(EVEN, b))) {
		t.Error("conversion_to_from even-r")
	}
	if !b.Equals(roffset_from_cube(ODD, roffset_to_cube(ODD, b))) {
		t.Error("conversion_to_from odd-r")
	}
}

func TestOffsetFromCube(t *testing.T) {
	a := NewOffsetCoord(1, 3)
	if !a.Equals(qoffset_from_cube(EVEN, NewHex(1, 2, -3))) {
		t.Error("offset_from_cube even-q")
	}

	b := NewOffsetCoord(1, 2)
	if !b.Equals(qoffset_from_cube(ODD, NewHex(1, 2, -3))) {
		t.Error("offset_from_cube odd-q")
	}
}

func TestOffsetToCube(t *testing.T) {
	a := NewHex(1, 2, -3)
	if !a.Equals(qoffset_to_cube(EVEN, NewOffsetCoord(1, 3))) {
		t.Error("offset_to_cube even-")
	}

	b := NewHex(1, 2, -3)
	if !b.Equals(qoffset_to_cube(ODD, NewOffsetCoord(1, 2))) {
		t.Error("offset_to_cube odd-q")
	}
}

func TestDoubledRoundtrip(t *testing.T) {
	a := NewHex(3, 4, -7)
	if !a.Equals(qdoubled_to_cube(qdoubled_from_cube(a))) {
		t.Error("conversion_from_to doubled-q")
	}
	if !a.Equals(rdoubled_to_cube(rdoubled_from_cube(a))) {
		t.Error("conversion_from_to doubled-r")
	}

	b := NewDoubledCoord(1, -3)
	if !b.Equals(qdoubled_from_cube(qdoubled_to_cube(b))) {
		t.Error("conversion_to_from doubled-q")
	}
	if !b.Equals(rdoubled_from_cube(rdoubled_to_cube(b))) {
		t.Error("conversion_to_from doubled-r")
	}
}

func TestDoubledFromCube(t *testing.T) {
	a := NewDoubledCoord(1, 5)
	if !a.Equals(qdoubled_from_cube(NewHex(1, 2, -3))) {
		t.Error("doubled_from_cube doubled-q")
	}

	b := NewDoubledCoord(4, 2)
	if !b.Equals(rdoubled_from_cube(NewHex(1, 2, -3))) {
		t.Error("doubled_from_cube doubled-r")
	}
}

func TestDoubledToCube(t *testing.T) {
	a := NewHex(1, 2, -3)
	if !a.Equals(qdoubled_to_cube(NewDoubledCoord(1, 5))) {
		t.Error("doubled_to_cube doubled-q")
	}

	b := NewHex(1, 2, -3)
	if !b.Equals(rdoubled_to_cube(NewDoubledCoord(4, 2))) {
		t.Error("doubled_to_cube doubled-r")
	}
}

////////////////////////////////////////////////////
// helper functions for testing

func equal_hex(t *testing.T, name string, a, b Hex) {
	if !a.Equals(b) {
		t.Error(name)
	}
}

func equal_offsetcoord(t *testing.T, name string, a, b OffsetCoord) {
	if !a.Equals(b) {
		t.Error(name)
	}
}

func equal_doubledcoord(t *testing.T, name string, a, b DoubledCoord) {
	if !a.Equals(b) {
		t.Error(name)
	}
}

func equal_int(t *testing.T, name string, a, b int) {
	if !(a == b) {
		t.Error(name)
	}
}

func equal_hex_array(t *testing.T, name string, a, b []Hex) {
	equal_int(t, name, len(a), len(b))
	for i := 0; i < len(a); i++ {
		equal_hex(t, name, a[i], b[i])
	}
}

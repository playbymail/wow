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

package hexes

import "math"

// abs is a helper function to get the absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// mod is a helper function to get the modulus of an integer
// (as opposed to %, which is the remainder operator)
func mod(a, b int) int {
	// you can check for b == 0 separately and do what you want
	if b < 0 {
		return -mod(-a, -b)
	}
	m := a % b
	if m < 0 {
		m += b
	}
	return m
}

func hex_add(a, b Hex) Hex {
	return a.Add(b)
}

func hex_corner_offset(layout Layout, corner int) Point {
	M := layout.orientation
	size := layout.size

	angle := 2.0 * math.Pi * (M.start_angle - float64(corner)) / 6.0

	return NewPoint(size.x*math.Cos(angle), size.y*math.Sin(angle))
}

func hex_diagonal_neighbor(hex Hex, direction int) Hex {
	return hex.DiagonalNeighbor(direction)
}

func hex_direction(direction int) Hex {
	// direction = mod(direction, 6)
	direction = (6 + (direction % 6)) % 6
	return hex_directions[direction]
}

func hex_distance(a, b Hex) int {
	return a.Distance(b)
}

func hex_length(hex Hex) int {
	return hex.Length()
}

// hex_lerp does a linear interpolation of
func hex_lerp(a, b FractionalHex, t float64) FractionalHex {
	return NewFractionalHex(a.q*(1.0-t)+b.q*t, a.r*(1.0-t)+b.r*t, a.s*(1.0-t)+b.s*t)
}

func hex_linedraw(a, b Hex) (results []Hex) {
	return a.LineDraw(b)
}

func hex_neighbor(hex Hex, direction int) Hex {
	return hex.Neighbor(direction)
}

func hex_rotate_left(a Hex) Hex {
	return a.RotateLeft()
}

func hex_rotate_right(a Hex) Hex {
	return a.RotateRight()
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

func hex_scale(a Hex, k int) Hex {
	return a.Scale(k)
}

func hex_subtract(a, b Hex) Hex {
	return a.Subtract(b)
}

func hex_to_pixel(layout Layout, h Hex) Point {
	return h.ToPixel(layout)
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

func polygon_corners(layout Layout, h Hex) (corners []Point) {
	return h.PolygonCorners(layout)
}

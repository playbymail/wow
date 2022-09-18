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

type FractionalHex struct {
	q, r, s float64
}

// NewFractionalHex returns an initialized FractionalHex
func NewFractionalHex(q, r, s float64) FractionalHex {
	if math.Round(q+r+s) != 0 {
		panic("assert(q + r + s == 0)")
	}
	return FractionalHex{q: q, r: r, s: s}
}

// Lerp does a linear interpolation of
func (fh FractionalHex) Lerp(b FractionalHex, t float64) FractionalHex {
	return NewFractionalHex(fh.q*(1.0-t)+b.q*t, fh.r*(1.0-t)+b.r*t, fh.s*(1.0-t)+b.s*t)
}

// Round returns the hex that the fractional hex is located in.
func (fh FractionalHex) Round() Hex {
	qi := int(math.Round(fh.q))
	q_diff := math.Abs(float64(qi) - fh.q)

	ri := int(math.Round(fh.r))
	r_diff := math.Abs(float64(ri) - fh.r)

	si := int(math.Round(fh.s))
	s_diff := math.Abs(float64(si) - fh.s)

	if q_diff > r_diff && q_diff > s_diff {
		qi = -ri - si
	} else if r_diff > s_diff {
		ri = -qi - si
	} else {
		si = -qi - ri
	}

	return NewHex(qi, ri, si)
}

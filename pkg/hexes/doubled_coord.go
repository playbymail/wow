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

type DoubledCoord struct {
	col, row int
}

func NewDoubledCoord(col, row int) DoubledCoord {
	return DoubledCoord{col: col, row: row}
}

func (a DoubledCoord) Equals(b DoubledCoord) bool {
	return a.col == b.col && a.row == b.row
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

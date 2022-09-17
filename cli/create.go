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

package cli

import (
	"github.com/mdhender/wow/pkg/board"
	"github.com/spf13/cobra"
	"io/ioutil"
)

// cmdCreateMap creates a map
var cmdCreateMap = &cobra.Command{
	Use:   "create",
	Short: "create a new map",
	Run: func(cmd *cobra.Command, args []string) {
		type hex struct {
			name      string
			col, row  int
			econValue int // non-zero only if hasStar
			warps     []string
		}

		// hexes is an untested guess as to the layout of the original map
		hexes := []hex{
			{"Adab", 10, 11, 0, []string{"Erech", "Khafa", "Byblos"}},
			{"Akkad", 19, 20, 3, []string{"Kish"}},
			{"Assur", 17, 11, 2, []string{"Nippur", "Lagash"}},
			{"Babylon", 22, 23, 4, []string{"Sumer"}},
			{"Byblos", 8, 13, 3, []string{"Adab"}},
			{"Calah", 9, 8, 1, []string{"Nippur"}},
			{"Elam", 16, 16, 5, []string{"Lagash"}},
			{"Erech", 7, 10, 3, []string{"Ur", "Adab"}},
			{"Eridu", 23, 17, 1, []string{"Kish", "Ugarit"}},
			{"Girsu", 18, 16, 1, []string{"Umma"}},
			{"Jarmo", 18, 14, 3, []string{"Kish"}},
			{"Isin", 16, 22, 1, []string{"Nineveh"}},
			{"Khafa", 13, 12, 2, []string{"Adab"}},
			{"Kish", 20, 18, 0, []string{"Jarmo", "Eridu"}},
			{"Lagash", 15, 14, 1, []string{"Assur"}},
			{"Larsu", 8, 4, 2, []string{"Susa"}},
			{"Mari", 14, 14, 1, []string{"Ubaid", "Umma"}},
			{"Mosul", 3, 7, 2, []string{"Sippur"}},
			{"Nineveh", 21, 25, 2, []string{"Isin"}},
			{"Nippur", 13, 10, 1, []string{"Calah", "Susa", "Assur", "Lagash"}},
			{"Sippur", 6, 1, 1, []string{"Mosul"}},
			{"Sumarra", 13, 19, 2, []string{"Ubaid", "Umma"}},
			{"Sumer", 19, 20, 0, []string{"Umma", "Babylon"}},
			{"Susa", 12, 7, 0, []string{"Larsu", "Nippur"}},
			{"Ubaid", 10, 14, 5, []string{"Mari", "Sumarra"}},
			{"Ugarit", 26, 22, 2, []string{"Eridu"}},
			{"Umma", 16, 19, 2, []string{"Sumarra", "Mari", "Girsu", "Sumer"}},
			{"Ur", 6, 6, 4, []string{"Erech"}},
		}

		// max row and col determine the size of the board
		maxRow, maxCol := 0, 0
		for _, h := range hexes {
			if h.row > maxRow {
				maxRow = h.row
			}
			if h.col > maxCol {
				maxCol = h.col
			}
		}

		// create the board, add all the stars, then add the wormholes
		gb := board.NewBoard(maxRow, maxCol)
		for _, h := range hexes {
			gb.AddStar(h.name, h.row, h.col, h.econValue)
		}
		for _, h := range hexes {
			for _, target := range h.warps {
				gb.AddWormHole(h.name, target)
			}
		}

		// save the board as an HTML file
		cobra.CheckErr(ioutil.WriteFile("svg-test.html", gb.AsHTML(), 0644))
	},
}

func init() {
	cmdBase.AddCommand(cmdCreateMap)
}

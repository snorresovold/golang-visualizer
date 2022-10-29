package main

import (
	"fmt"
	"reflect"
)

type cell struct {
	switching bool
	height    int
	start     int
	end       int
}

type cellList struct {
	cells []cell
}

func (cl *cellList) AddCell(c cell) []cell {
	cl.cells = append(cl.cells, c)
	return cl.cells
}

func main() {
	cl := cellList{}
	c := cell{false, 1, 2, 3}
	cl.AddCell(c)
	c2 := cell{true, 100, 2123, 3342}
	cl.AddCell(c2)
	fmt.Println(cl.cells)

	swapC := reflect.Swapper(cl.cells)
	cl.cells[0].switching, cl.cells[1].switching = true, true
	swapC(0, 1)
	fmt.Println(cl.cells)
	cl.cells[0].switching, cl.cells[1].switching = false, false
	fmt.Println(cl.cells)
}

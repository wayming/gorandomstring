package main

import "fmt"

type Point struct {
	X float64
	Y float64
}

func (p *Point) Init() {
	p.X = 100
	p.Y = 101
	fmt.Println(p)
}

func (p Point) Dump() {
	fmt.Println(p)
}

func main() {
	fmt.Println("method test")
	// p1 := Point{1, 2}
	// p1 := new(Point)
	var p1 *Point
	p1 = new(Point)
	p1.Init()
	p1.Dump()
	(*p1).Dump()
	(*p1).Init()

}

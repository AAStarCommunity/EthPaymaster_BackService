package model

import "fmt"

func main() {
	impl1 := &Impl1{Ver: "1"}
	impl2 := &Impl2{
		Impl1{
			Base: nil,
			Ver:  "2",
		}}
	Hello(impl1)
	Hello(impl2)
}

type Base interface {
	GetVer() string
}

type Impl1 struct {
	Base
	Ver string
}

func (a *Impl1) GetVer() string {
	fmt.Println(a.Ver)
	return a.Ver
}

type Impl2 struct {
	Impl1
}

func Hello(b Base) {
	if b1, ok := b.(*Impl1); ok {
		b1.GetVer()
	} else {
		b2, _ := b.(*Impl2)
		b2.GetVer()
	}
}

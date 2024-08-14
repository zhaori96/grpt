package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/zhaori96/grpt"
	"github.com/zhaori96/grpt/temp"
)

func main() {
	//SomeTests()
	now := time.Now()
	x := temp.BoletoPrecito{}
	x.Imprimir()
	fmt.Printf("Total: %s\n", time.Since(now))

}

func SomeTests() {
	doc := grpt.Document{
		PageSize: grpt.PageSizeA4,
		Padding:  grpt.DefaultPagePadding,
		Body: grpt.DocumentBody{
			Elements: grpt.Elements{
				&grpt.Visible{
					Visible:           false,
					AlwaysOccupySpace: true,
					Child:             &grpt.Text{Value: "Value", Size: grpt.NewSize(50, 50)},
				},
				&grpt.Visible{
					Visible: true,
					Child:   &grpt.Text{Value: "Value2", Size: grpt.NewSize(50, 50)},
				},
				&grpt.Selector{
					Selector: func(length int) int {
						return rand.Intn(length - 1)
					},
					Elements: grpt.Elements{
						&grpt.Text{Value: "ValueA", Size: grpt.NewSize(50, 50)},
						&grpt.Text{Value: "ValueB", Size: grpt.NewSize(50, 50)},
						&grpt.Text{Value: "ValueC", Size: grpt.NewSize(50, 50)},
						&grpt.Text{Value: "ValueD", Size: grpt.NewSize(50, 50)},
						&grpt.Text{Value: "ValueE", Size: grpt.NewSize(50, 50)},
						&grpt.Text{Value: "ValueF", Size: grpt.NewSize(50, 50)},
					},
				},
			},
		},
	}

	doc.Build("temp/some_tests.pdf")
}

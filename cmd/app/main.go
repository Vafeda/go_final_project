package main

import (
	"github.com/Vafeda/go_final_project/internal/app"
	"github.com/Vafeda/go_final_project/internal/nextdate"
)

func main() {
	//now, _ := time.Parse("20060102", "20240126")
	nextdate.Test()
	//fmt.Println(nextdate.NextDate(now, "20240113", "d 7"))
	app.Run()

}

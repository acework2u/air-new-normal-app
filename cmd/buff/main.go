package main

import (
	"fmt"
	"time"
)

func main() {

	YYYYMMDD := "2006-01-02 3:4:5 pm"
	current_time := time.Now()

	fmt.Println(YYYYMMDD)
	fmt.Println(current_time)

	fmt.Println(current_time.Format(YYYYMMDD))
}

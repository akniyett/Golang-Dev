package main
import (
  "net/http"
  "fmt"
  "time"
  "math"
  
)

type СoeffOfBody struct {
  weight float64
  height float64
}

// type SleepNow struct {
//   name string
//   time string
// }

func main() {

  

  http.HandleFunc("/", timeForHealth)

  fmt.Println(http.ListenAndServe(":8080", nil));
}

func timeForHealth(w http.ResponseWriter, r *http.Request) {
    coeff := СoeffOfBody{weight:50, height:170}
    // water := SleepNow{"Eren", time.Now().Format(time.Stamp)}
    timein := time.Now().Local().Add(time.Hour * time.Duration(9) +
                                 time.Minute * time.Duration(0) +
                                 time.Second * time.Duration(0))
    fmt.Fprintf(w, "Eren should wake up at %v\n",timein )
    fmt.Fprintf(w, "Your body coefficient is %v", coeff.weight/(math.Pow(coeff.height/100, 2)))

}
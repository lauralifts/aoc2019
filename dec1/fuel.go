package main

import "bufio"
import "fmt"
import "log"
import "math"
import "os"
import "strconv"

func main(){

  moduleMasses := getModuleMasses()
  fuelNeeded := 0

  for _, mass := range moduleMasses {
    fuelNeeded += getFuel(mass)
  }

  fmt.Printf("Fuel needed is %d\n", fuelNeeded)
}

func getFuel(mass int) int {
  fuelNeeded := getFuelInner(mass)

  // adjust for weight of fuel itself
  fuelForFuel := getFuelInner(fuelNeeded)
  for fuelForFuel >= 0 {
    lastFFF := fuelForFuel
    fuelNeeded += fuelForFuel
    fmt.Printf("Adding %d to cover fuel\n", fuelForFuel)
    fuelForFuel = getFuelInner(lastFFF)
  }

  return fuelNeeded
}

func getFuelInner(mass int) int {
  result := math.Floor(float64(mass)/float64(3)) - float64(2)
  return int(result)
}

func getModuleMasses() []int {
  file, err := os.Open("./input.txt")
  if err != nil {
    log.Fatal(err)
  }
  defer file.Close()
  result := make([]int, 0)

  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    line := scanner.Text()
    //fmt.Println(line)
    num, err := strconv.Atoi(line)
    if err != nil {
      log.Fatal(err)
    }

    result = append(result, num)
  }

  if err := scanner.Err(); err != nil {
    log.Fatal(err)
  }

  return result
}

package main

import "fmt"

const (
  MAXPS int = 4
)

func main() {
  text := getText("input.txt")

  // phase settings go from 0 to 4
  phaseSettings := []int{0, 0, 0, 0, 0}
  lastPS := false

  var out int
  var err error

  maxseen := 0
  for ;lastPS != true; {
    out, err = runAmpSequence(phaseSettings, text)
    if (out > maxseen) {
      maxseen = out
      fmt.Printf("Best phase settings so far: %v with output %d\n", phaseSettings, out)
    }
    if err != nil {
      fmt.Printf("Error %v\n", err)
    }
    phaseSettings, lastPS = incPhaseSettings(phaseSettings)
  }

  fmt.Printf("Program done, final outputs are %v\n", maxseen)
}

// new settings, return true if can't increment anymore
func incPhaseSettings(ps []int) ([]int, bool) {
  for i := 0; i < len(ps); i++ {
    if ps[i] < MAXPS {
      ps[i]++
      return ps, false
    } else {
      ps[i] = 0
    }
  }
  return ps, true
}

func runAmpSequence(phaseSettings []int, text string) (int, error) {
    curInput := 0
    var out []int
    var err error
    for i := 0; i < len(phaseSettings); i++ {
      out, err = runAmplifier(text, phaseSettings[i], curInput)
      if err != nil {
        return 0, err
      }

      //fmt.Printf("Amp run %d, phase setting %d, input %d, output %d\n",
      //  i, phaseSettings[i], curInput, out[0])
      curInput = out[0]
    }

    //fmt.Printf("Program done, final outputs are %v\n", out)
    return out[0],  nil
}

func runAmplifier(text string, phaseCode int, input int) ([]int, error) {
    program := getIntCodes(text)
    inputs := make([]int, 0)
    inputs = append(inputs, phaseCode, input)

    outputs, err := runProgram(program, inputs)
    if err != nil {
      fmt.Printf("%v", err)
    }
    return outputs, nil
}

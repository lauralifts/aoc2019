package main

import "bufio"
import "fmt"
import "log"
import "os"
import "strconv"
import "strings"

const (
  ADD = 1
  MUL = 2
  HALT = 99
)

func main() {
  part2()
}

func part2() {
  text := getText()

  noun := 0
  verb := 0


  for noun = 0; noun < 100; noun++ {
    for verb = 0; verb < 100; verb++ {
      program := getIntCodes(text)
      result, err := runProgram(program, noun, verb)

      fmt.Printf("%d, %d -> %d, %v\n", noun, verb, result, err)
      if result == 19690720 {
        log.Fatal("done")
      }
    }
  }
}

func part1() {
  text := getText()
  program := getIntCodes(text)
  //printProgram(program)

  result, _ := runProgram(program, 12, 2)
  fmt.Printf("Output: %d\n", result)
}

func runProgram(program []*int, noun int, verb int) (int, error) {
  program[1] = &noun
  program[2] = &verb

  ptr := 0
  var err error
  err = nil
  for *program[ptr] != HALT {
    switch *program[ptr] {
    case ADD:
      ptr, err = doADD(program, ptr)
      if err != nil {
          return 0, err
      }
    case MUL:
      ptr, err = doMUL(program, ptr)
      if err != nil {
          return 0, err
      }
    default:
      log.Fatalf("Unrecognised opcode: %d\n", *program[ptr])
    }

    //fmt.Printf("New program state:")
    //printProgram(program)
  }

  //fmt.Printf("HALT at %d\n", ptr)

  return *program[0], nil
}

// Returns new position
func doADD(program []*int, pos int) (int, error) {
  op1 := *program[pos+1]
  op2 := *program[pos+2]
  out := *program[pos+3]

  if (out >= len(program) || op1 >= len(program) || op2 >= len(program)) {
      return 0, fmt.Errorf("Out of range")
  }

  val := *program[op1] + *program[op2]
  program[out] = &val
  //fmt.Printf("ADD stored %d at %d\n", *program[out], out)
  return pos + 4, nil
}

// Returns new position
func doMUL(program []*int, pos int) (int, error) {
  op1 := *program[pos+1]
  op2 := *program[pos+2]
  out := *program[pos+3]
  if (out >= len(program) || op1 >= len(program) || op2 >= len(program)) {
      return 0, fmt.Errorf("Out of range")
  }

  val := *program[op1] * *program[op2]
  program[out] = &val
  //fmt.Printf("MUL stored %d at %d\n", *program[out], out)
  return pos + 4, nil
}

func getText() string {
  file, err := os.Open("./input_real.txt")
  if err != nil {
    log.Fatal(err)
  }
  defer file.Close()

  result := ""

  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    line := scanner.Text()
    result = result + line
  }
  return result
}

func getIntCodes(line string) []*int {
  result := make([]*int, 0, 1000)

  fields := strings.Split(line, ",")
  for _, f := range fields {
    num, err := strconv.Atoi(f)
    if err != nil {
      log.Fatal(err)
    }

    result = append(result, &num)
   }

  return result
}

func printProgram(program []*int){
  for k, ptr := range program {
    if k%4 == 0 {
      switch *program[k] {
      case ADD:
        fmt.Printf("\n%d: ADD, ", k)
      case MUL:
        fmt.Printf("\n%d: MUL, ", k)
      case HALT:
        fmt.Printf("\n%d: HALT, ", k)
      default:
        fmt.Printf("%d, ", *ptr)
      }
    } else {
      fmt.Printf("%d, ", *ptr)
    }
  }
  fmt.Printf("\n")
}

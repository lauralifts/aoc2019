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
  INPUT = 3
  OUTPUT = 4
  HALT = 99

  POS_MODE = 0
  IMM_MODE = 1
)

type op func(int, int) int

func main() {
  text := getText()
  program := getIntCodes(text)
  printProgram(program)

  err := runProgram(program)
  if err != nil {
    fmt.Printf("%v", err)
  }
}

func runProgram(program []*int) error {
  ptr := 0
  var err error
  err = nil
  for *program[ptr] != HALT {
    switch opcode(*program[ptr]) {
    case ADD:
      ptr, err = doADD(program, ptr)
      if err != nil {
          return err
      }
    case MUL:
      ptr, err = doMUL(program, ptr)
      if err != nil {
          return err
      }
    case INPUT:
      ptr, err = doInput(program, ptr)
      if err != nil {
          return err
      }
    case OUTPUT:
      ptr, err = doOutput(program, ptr)
      if err != nil {
          return err
      }
    default:
      log.Fatalf("Unrecognised opcode: %d at position %d\n", *program[ptr], ptr)
    }
  }

  return nil
}


// Returns new position
func doInput(program []*int, pos int) (int, error) {
  reader := bufio.NewReader(os.Stdin)
  fmt.Print("Enter input: ")
  text, _ := reader.ReadString('\n')
  val, err := strconv.Atoi(strings.TrimSpace(text))
  if err != nil {
      return 0, err
  }

  out := *program[pos+1]    // Can't be in immediate mode
  if (out >= len(program) || out < 0) {
      return 0, fmt.Errorf("Out of range: position %d, out %d", pos, out)
  }
  program[out] = &val
  return pos + 2, nil
}

// Returns new position
func doOutput(program []*int, pos int) (int, error) {
  mode1, _ := paramModes(*program[pos])
  op1, err := param(program, pos+1, mode1)
  if err != nil {
    return 0, err
  }

  fmt.Printf("\n***Output: %d***\n", op1)

  return pos + 2, nil
}

// Returns new position
func doMUL(program []*int, pos int) (int, error) {
  return doFunc(program, pos, mul)
}

// Returns new position
func doADD(program []*int, pos int) (int, error) {
  return doFunc(program, pos, add)
}

func doFunc(program []*int, pos int, fn op) (int, error){
  mode1, mode2 := paramModes(*program[pos])
  op1, err := param(program, pos+1, mode1)
  if err != nil {
    return 0, err
  }
  op2, err := param(program, pos+2, mode2)
  if err != nil {
    return 0, err
  }

  out := *program[pos+3]    // Can't be in immediate mode

  if (out >= len(program) || out < 0) {
      return 0, fmt.Errorf("Out of range: position %d, out %d", pos, out)
  }
  val := fn(op1, op2)

  program[out] = &val
  return pos + 4, nil
}

func param(program []*int, pos int, mode int) (int, error) {
    if mode == POS_MODE {
      p := *program[pos]
      if p >= len(program) || p < 0 {
        return 0, fmt.Errorf("Out of range: position %d, out %d", pos, p)
      }
      return *program[p], nil
    } else {
      return *program[pos], nil
    }
}

func add(x int, y int) int {
    return x + y
}

func mul(x int, y int) int {
    return x * y
}

func getText() string {
  file, err := os.Open("./input.txt")
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

func opcode(instr int) int {
  return instr % 100
}

func paramModes(instr int) (int, int) {
  mode1, mode2 := POS_MODE, POS_MODE
  if instr >= 10000 {
    instr -= 10000
  }
  if instr >= 1000 {
    mode2 = IMM_MODE
    instr -= 1000
  }
  if instr >= 100 {
    mode1 = IMM_MODE
  }
  return mode1, mode2
}


func printProgram(program []*int){
  k := 0
  for k < len(program) {
    switch opcode(*program[k]) {
      case ADD:
        fmt.Printf("\n%d: ADD, ", *program[k])
        k += 4
      case MUL:
        fmt.Printf("\n%d: MUL, ", *program[k])
        k += 4
      case INPUT:
        fmt.Printf("\n%d: INPUT, ", *program[k])
        k += 2
      case OUTPUT:
        fmt.Printf("\n%d: OUTPUT, ", *program[k])
        k += 2
      case HALT:
        fmt.Printf("\n%d: HALT, ", *program[k])
        k += 1
      default:
        fmt.Printf("%d, ", *program[k])
        k++
      }
  }
  fmt.Printf("\n")
}

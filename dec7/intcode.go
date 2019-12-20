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
  JUMPIFTRUE = 5
  JUMPIFFALSE = 6
  LESSTHAN = 7
  EQUALS = 8
  HALT = 99

  POS_MODE = 0
  IMM_MODE = 1
)

type op func(int, int) int
type lop func(int) bool


func runProgram(program []*int, inputs []int) ([]int, error) {
  ptr := 0
  inptr := 0
  outputs := make([]int, 0)

  var err error
  err = nil
  for *program[ptr] != HALT {
    switch opcode(*program[ptr]) {
    case ADD:
      ptr, err = doADD(program, ptr)
      if err != nil {
          return nil, err
      }
    case MUL:
      ptr, err = doMUL(program, ptr)
      if err != nil {
          return nil, err
      }
    case INPUT:
      ptr, err = doInput(program, ptr, inputs[inptr]) // might error out
      inptr++
      if err != nil {
          return nil, err
      }
    case OUTPUT:
      var out int
      ptr, out, err = doOutput(program, ptr)
      if err != nil {
          return nil, err
      }
      outputs = append(outputs, out)
    case JUMPIFTRUE:
      ptr, err = doJumpTrue(program, ptr)
      if err != nil {
          return nil, err
      }
    case JUMPIFFALSE:
      ptr, err = doJumpFalse(program, ptr)
      if err != nil {
          return nil, err
      }
    case LESSTHAN:
      ptr, err = doLessThan(program, ptr)
      if err != nil {
          return nil, err
      }
    case EQUALS:
      ptr, err = doEquals(program, ptr)
      if err != nil {
          return nil, err
      }
    default:
      log.Fatalf("Unrecognised opcode: %d at position %d\n", *program[ptr], ptr)
    }
  }

  return outputs, nil
}

// Returns new position
func doInput(program []*int, pos int, input int) (int, error) {
  out := *program[pos+1]    // Can't be in immediate mode
  if (out >= len(program) || out < 0) {
      return 0, fmt.Errorf("Out of range: position %d, out %d", pos, out)
  }
  program[out] = &input
  return pos + 2, nil
}

// Returns new position
func doOutput(program []*int, pos int) (int, int, error) {
  mode1, _, _ := paramModes(*program[pos])
  op1, err := param(program, pos+1, mode1)
  if err != nil {
    return 0, 0, err
  }

  //fmt.Printf("\n***Output: %d***\n", op1)

  return pos + 2, op1, nil
}

// Returns new position
func doLessThan(program []*int, pos int) (int, error) {
  return doFunc(program, pos, lessthan)
}

// Returns new position
func doEquals(program []*int, pos int) (int, error) {
  return doFunc(program, pos, equals)
}

// Returns new position
func doJumpTrue(program []*int, pos int) (int, error) {
  return doJumpFunc(program, pos, truefn)
}

// Returns new position
func doJumpFalse(program []*int, pos int) (int, error) {
  return doJumpFunc(program, pos, falsefn)
}

func doJumpFunc(program []*int, pos int, fn lop) (int, error){
  mode1, mode2, _ := paramModes(*program[pos])
  op1, err := param(program, pos+1, mode1)
  if err != nil {
    return 0, err
  }
  op2, err := param(program, pos+2, mode2)
  if err != nil {
    return 0, err
  }

  if fn(op1) {
    return op2, nil
  }

  return pos + 3, nil
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
  mode1, mode2, _ := paramModes(*program[pos])
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

func truefn(x int) bool {
    return x != 0
}

func falsefn(x int) bool {
    return x == 0
}

func lessthan(x int, y int) int {
    if x < y {
      return 1
    }
    return 0
}

func equals(x int, y int) int {
    if x == y {
      return 1
    }
    return 0
}

func getText(filename string) string {
  file, err := os.Open(filename)
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

func paramModes(instr int) (int, int, int) {
  mode1, mode2, mode3 := POS_MODE, POS_MODE, POS_MODE
  if instr >= 10000 {
    instr -= 10000
    mode3 = IMM_MODE
  }
  if instr >= 1000 {
    mode2 = IMM_MODE
    instr -= 1000
  }
  if instr >= 100 {
    mode1 = IMM_MODE
  }
  return mode1, mode2, mode3
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

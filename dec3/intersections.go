package main

import "bufio"
import "fmt"
import "log"
import "os"
import "strconv"
import "strings"

type pt struct {
  x int
  y int
}

var  LEFT = pt{-1, 0}
var  RIGHT = pt{1, 0}
var  UP = pt{0, 1}
var  DOWN = pt{0, -1}

type segment struct {
  direction pt
  len int
}

type wire struct {
  segments []segment
}

func main() {
  wires := getWires()
  fmt.Printf("Wires: %v\n", wires)

  ints := wires[0].getIntersections(wires[1])

  if len(ints) == 0 {
    log.Fatal("No intersections")
  }

  minsteps := wires[0].steps(ints[0]) + wires[1].steps(ints[0])
  minpt := ints[0]
  for i := 1; i < len(ints); i++ {
    steps := wires[0].steps(ints[i]) + wires[1].steps(ints[i])
    if steps < minsteps {
      minsteps = steps
      minpt = ints[i]
    }
  }

  fmt.Printf("Nearest intersection is %v, steps %d\n", minpt, minsteps)

}

func part1() {
  wires := getWires()
  fmt.Printf("Wires: %v\n", wires)

  ints := wires[0].getIntersections(wires[1])

  if len(ints) == 0 {
    log.Fatal("No intersections")
  }

  mindist := ints[0].dist()
  minpt := ints[0]
  for i := 1; i < len(ints); i++ {
    if ints[i].dist() < mindist {
      mindist = ints[i].dist()
      minpt = ints[i]
    }
  }

  fmt.Printf("Nearest intersection is %v, distance %d\n", minpt, mindist)

}

func NewWire() wire {
  res := wire {
      segments : make([]segment, 0),
  }
  return res
}

func add(a pt, b pt) pt {
  return pt{a.x + b.x, a.y + b.y}
}

// dist from origin
func (p pt) dist() int {
  sum := 0
  if p.x < 0 {
    sum += p.x * -1
  } else {
    sum += p.x
  }
  if p.y < 0 {
    sum += p.y * -1
  } else {
    sum += p.y
  }
  return sum
}

func (w wire) steps(p pt) int {
  result := 0

  last := pt{0, 0}
  for _, seg := range w.segments {
    for i := 0; i < seg.len; i++ {
      cur := add(last, seg.direction)
      result ++

      if cur.x == p.x && cur.y == p.y {
        return result
      }

      last = cur
    }
  }

  log.Fatalf("Can't get to point %v", p)
  return 0
}


func (w wire) getIntersections(other wire) []pt {
  mypts := w.getPoints()
  othpts := other.getPoints()

  result := make([]pt, 0)

  for i := 1; i < len(mypts); i++ {
    for j := 1; j < len(othpts); j++ {
      if mypts[i].x == othpts[j].x && mypts[i].y == othpts[j].y {
        result = append(result, mypts[i])
      }
    }
  }

  return result
}

func (w wire) getPoints() []pt {
  result := make([]pt, 0)
  last := pt{0, 0}
  result = append(result, last)

  for _, seg := range w.segments {
    for i := 0; i < seg.len; i++ {
      cur := add(last, seg.direction)
      result = append(result, cur)
      last = cur
    }
  }

  return result
}

func getWires() []wire {
  file, err := os.Open("./input_real.txt")
  if err != nil {
    log.Fatal(err)
  }
  defer file.Close()

  wires := make([]wire, 0)
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    line := scanner.Text()
    wire := NewWire()

    fields := strings.Split(line, ",")
    for _, f := range fields {
      var seg segment
      switch f[0] {
      case 'L':
        seg.direction = LEFT
      case 'R':
        seg.direction = RIGHT
      case 'U':
        seg.direction = UP
      case 'D':
        seg.direction = DOWN
      default:
        log.Fatal("Line invalid: %s at field %v", line, f)
      }

      v, err := strconv.Atoi(f[1:])
      if err != nil {
        log.Fatal(err)
      }
      seg.len = v
      wire.segments = append(wire.segments, seg)
    }
    wires = append(wires, wire)
  }
  return wires
}

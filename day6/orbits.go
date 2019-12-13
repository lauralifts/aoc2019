package main

import "bufio"
import "fmt"
import "log"
import "os"
import "strings"

type node struct {
  name string
  children []*node
  parent *node
}


func main() {
  orbits := getOrbits()
  fmt.Printf("Tree is rooted at %s\n", orbits.name)

  result := orbits.countChildRelationships(0)
  fmt.Printf("Total orbits %d\n", result)
}

func newNode(name string) node {
  n := node {
    name: name,
  }
  //fmt.Printf("Create node %s\n", name)



  return n
}

func (n *node) addChild(c *node){
  n.children = append(n.children, c)
}


func (n *node) countChildRelationships(depthToHere int) int {
  result := len(n.children) * depthToHere
  for _, c := range n.children {
    result += c.countChildRelationships(depthToHere+1) + 1
  }

  fmt.Printf("Node %s, my child relationships count is %d\n", n.name, result)
  return result
}

func (n *node) setParent(c *node){
  if n.parent != nil {
    log.Fatalf("oops, node %s has two parents", n.name)
  }
  n.parent = c
}

func (n *node) getRoot() *node {
  cur := n
  for cur.parent != nil {
    cur = cur.parent
  }
  return cur
}

func getOrbits() *node {
  file, err := os.Open("./input.txt")
  if err != nil {
    log.Fatal(err)
  }
  defer file.Close()

  dict := make(map[string]*node)
  var lastnode *node

  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    line := scanner.Text()

    if err != nil {
      log.Fatal(err)
    }
    parts := strings.Split(line, ")")
    if len(parts) != 2 {
     log.Fatalf("bad input %s", line)
    }
    bodystr := parts[0]
    satstr := parts[1]

    if dict[bodystr] == nil {
      n := newNode(bodystr)
      dict[bodystr] = &n
    }
    if dict[satstr] == nil {
      n := newNode(satstr)
      dict[satstr] = &n
    }
    dict[bodystr].addChild(dict[satstr])
    dict[satstr].setParent(dict[bodystr])
    lastnode = dict[bodystr]
    // fmt.Printf("Adding node: %s orbits %s\n", satstr, bodystr)
  }

  if err := scanner.Err(); err != nil {
    log.Fatal(err)
  }

  result := lastnode.getRoot()
  return result
}

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
  orbits, dict := getOrbits()
  fmt.Printf("Tree is rooted at %s\n", orbits.name)

  santa := dict["SAN"]
  you := dict["YOU"]

  //fmt.Printf("Santas %v, path to root %v\n", santa, santa.pathToRoot())
  //fmt.Printf("You %v, path to root %v\n", you, you.pathToRoot())

  sp := santa.pathToRoot()
  yp := you.pathToRoot()

  i := 0
  for ;i < len(sp) && i < len(yp); i++{
    if sp[i].name != yp[i].name {
      break
    }
  }

  i--
  fmt.Printf("First common ancestor is %s at %d\n", sp[i].name, i)

  // subtract 1 for the common ancestor, and 2 because we're transferring
  // between parents of YOU and SAN
  result := len(sp) + len(yp) - 2*i -4
  fmt.Printf("Transfer orbits: %d\n", result)
}

func part1() {
  orbits, _ := getOrbits()
  fmt.Printf("Tree is rooted at %s\n", orbits.name)

  result := orbits.countChildRelationships(0)
  fmt.Printf("Total orbits %d\n", result)
}


func newNode(name string) node {
  n := node {
    name: name,
  }

  return n
}

func (n *node) pathToRoot() []*node {
  cur := n
  path := make([]*node, 0)
  path = append(path, n)
  for cur.parent != nil {
    path = append(path, cur.parent)
    cur = cur.parent
  }

  //reverse
  for i := len(path)/2-1; i >= 0; i-- {
	  opp := len(path)-1-i
	  path[i], path[opp] = path[opp], path[i]
  }
  return path
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

func getOrbits() (*node,map[string]*node){
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
  return result, dict
}

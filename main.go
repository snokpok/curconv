package main

import (
	"fmt"
	"os"
)

type AdjListNode struct {
  endpoint string
  value float64
}

type AdjList map[string][]AdjListNode

type CurrencyPair struct {
  left string
  right string
  value float64
}

func findExchangeRate(left, right string, adjList AdjList) (float64, error) {
  if _, ok := adjList[left]; !ok {
    return -1, fmt.Errorf("No connecting currency pair provided for %s", left)
  }
  if _, ok := adjList[right]; !ok {
    return -1, fmt.Errorf("No connecting currency pair provided for %s", right)
  }
  fwdtable := map[string]AdjListNode{}
  visited := map[string]bool{}
  stack := []string{}
  stack = append(stack, left)
  for len(stack)!=0 {
    top := stack[len(stack)-1]
    stack = stack[:len(stack)-1]

    if top == right {
      break
    }

    if vsted, ok := visited[top]; !ok || !vsted {
      visited[top] = true
      for _, neighbor := range adjList[top] {
        stack = append(stack, neighbor.endpoint)
        fwdtable[top] = AdjListNode{neighbor.endpoint, neighbor.value}
      }
    }
  }


  excRate := 1.0
  for {
    nexthop := fwdtable[left]
    excRate = excRate*nexthop.value
    if nexthop.endpoint == right {
      break
    }
    left = nexthop.endpoint
  }

  return excRate, nil
}

func main() {
  pairs := []CurrencyPair{
    {"USD", "CAD", 1.35},
    {"CHF", "CAD", 1.53},
  }

  adjList := AdjList{}
  for _, p := range pairs {
    if _, ok := adjList[p.left]; !ok {
      adjList[p.left] = []AdjListNode{}
    }
    adjList[p.left] = append(adjList[p.left], AdjListNode{
      endpoint: p.right, 
      value: p.value,
    })
    if _, ok := adjList[p.right]; !ok {
      adjList[p.right] = []AdjListNode{}
    }
    adjList[p.right] = append(adjList[p.right], AdjListNode{
      endpoint: p.left,
      value: 1/p.value,
    });
  }

  fmt.Println(adjList)
  
  left := "USD"
  right := "CHF"
  excRate, err := findExchangeRate(left, right, adjList)
  if err !=nil {
    fmt.Println(err.Error())
    os.Exit(1)
  }
  fmt.Printf("1 %s = %f %s\n", left, excRate, right)
}

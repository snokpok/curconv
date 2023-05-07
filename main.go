package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/urfave/cli/v2"
)

type AdjListNode struct {
	endpoint string
	value    float64
}

type AdjList map[string][]AdjListNode

type CurrencyPair struct {
	left  string
	right string
	value float64
}

func findExchangeRate(left, right string, adjList AdjList, showStepByStep bool) (float64, *[]CurrencyPair, error) {
	if _, ok := adjList[left]; !ok {
		return -1, nil, fmt.Errorf("No connecting currency pair provided for %s", left)
	}
	if _, ok := adjList[right]; !ok {
		return -1, nil, fmt.Errorf("No connecting currency pair provided for %s", right)
	}

	// DFS down the currency graph
	fwdtable := map[string]AdjListNode{}
	visited := map[string]bool{}
	stack := []string{}
	stack = append(stack, left)
	for len(stack) != 0 {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if vsted, ok := visited[top]; !ok || !vsted {
			visited[top] = true
			for _, neighbor := range adjList[top] {
				if visited[neighbor.endpoint] {
					continue
				}
				stack = append(stack, neighbor.endpoint)
				fwdtable[top] = AdjListNode{neighbor.endpoint, neighbor.value}
			}
		}
	}

	if showStepByStep {
		log.Printf("Done traversing currency pairs; successor table: %v\n", fwdtable)
	}

	// trace back from left -> right via predecessors
	excRate := 1.0
	interconvs := []CurrencyPair{}
	for {
		nexthop := fwdtable[left]
		excRate = excRate * nexthop.value
		if showStepByStep {
			interconvs = append(interconvs, CurrencyPair{left, nexthop.endpoint, nexthop.value})
		}
		if nexthop.endpoint == right {
			break
		}
		left = nexthop.endpoint
	}

	if showStepByStep {
		return excRate, &interconvs, nil
	}

	return excRate, nil, nil
}

// construct adjacency list from the currency pair
func constructAdjListFromCurrencyPairs(pairs []CurrencyPair) AdjList {
	adjList := AdjList{}
	for _, p := range pairs {
		if _, ok := adjList[p.left]; !ok {
			adjList[p.left] = []AdjListNode{}
		}
		adjList[p.left] = append(adjList[p.left], AdjListNode{
			endpoint: p.right,
			value:    p.value,
		})
		if _, ok := adjList[p.right]; !ok {
			adjList[p.right] = []AdjListNode{}
		}
		adjList[p.right] = append(adjList[p.right], AdjListNode{
			endpoint: p.left,
			value:    1 / p.value,
		})
	}
	return adjList
}

func readCurrenciesCSV(path string) (*[]CurrencyPair, error) {
	pairs := []CurrencyPair{}

	f, err := os.Open(path)
	defer f.Close()
	reader := bufio.NewScanner(f)
	reader.Split(bufio.ScanLines)

	linenum := 1
	for reader.Scan() {
		line := reader.Text()
		if err != nil {
			break
		}
		linenum++
		splitted := strings.Split(line, ",")
		if len(splitted) != 3 {
			return nil, fmt.Errorf("Line %d is invalid; must have exactly 3 values between commas; got '%s'\n", linenum, line)
		}
		left := strings.Trim(splitted[0], " ")
		right := strings.Trim(splitted[1], " ")
		value, err := strconv.ParseFloat(splitted[2], 64)
		if err != nil {
			return nil, fmt.Errorf("Line %d is invalid; value not an integer; got '%s'\n", linenum, splitted[2])
		}
		pairs = append(pairs, CurrencyPair{left, right, value})
	}
	return &pairs, nil
}

func main() {
	app := &cli.App{
		Name:    "curconv",
		Usage:   "A currency converter CLI tool that can find the rate between any two pairs in your given currency list",
		Suggest: true,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Usage:   "Do it verbosely",
			},
			&cli.StringFlag{
				Name:     "file",
				Aliases:  []string{"f"},
				Usage:    "Path to CSV file containing comma-delimited currency pairs",
				Required: true,
			},
		},
		UsageText: "./curconv [global options] <currency from> <currency to> [args]",
		Action: func(ctx *cli.Context) error {
			left := ctx.Args().Get(0)
			right := ctx.Args().Get(1)
			verbose := ctx.Bool("v")

			pairs, err := readCurrenciesCSV(ctx.String("file"))
			if err != nil {
				return err
			}

			adjList := constructAdjListFromCurrencyPairs(*pairs)

			if verbose {
				log.Printf("Created adjacency list of currency pairs: %v\n", adjList)
			}

			if left == "" {
				return fmt.Errorf("Must provide the currency to convert from")
			}
			if right == "" {
				return fmt.Errorf("Must provide the currency to convert to")
			}
			excRate, interconvs, err := findExchangeRate(left, right, adjList, verbose)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			if verbose {
				fmt.Printf("Steps:\n")
				for _, pair := range *interconvs {
					fmt.Printf("1 %s = %f %s\n", pair.left, pair.value, pair.right)
				}
				fmt.Printf("--------------\n")
			}

			fmt.Printf("1 %s = %f %s\n", left, excRate, right)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}

package engine

import (
	"bufio"
	"context"
	"fmt"
	. "github.com/mhib/combusken/backend"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type positionTest struct {
	Position
	bestMoves []string
	id        string
}

func loadEPD(fileLocation string) (res []positionTest) {
	path, _ := filepath.Abs(fileLocation)
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		var entry positionTest
		positionAndRest := strings.Split(line, " bm ")
		entry.Position = ParseFen(positionAndRest[0])
		bestMovesAndID := strings.Split(positionAndRest[1], ";")
		entry.bestMoves = strings.Split(bestMovesAndID[0], " ")
		entry.id = bestMovesAndID[1]
		res = append(res, entry)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return
}

func Benchmark() {
	engine := NewEngine()
	engine.Threads.Val = 1
	engine.Hash.Val = 256
	start := time.Now()
	nodes := 0
	for _, entry := range loadEPD("./engine/test_positions/WinAtChess.epd") {
		tmpNodes := 0
		engine.NewGame()
		engine.Update = func(si SearchInfo) {
			tmpNodes = si.Nodes
		}
		engine.Search(context.Background(), SearchParams{Positions: []Position{entry.Position}, Limits: LimitsType{Depth: 5}})
		nodes += tmpNodes
	}
	duration := time.Now().Sub(start)
	fmt.Printf("Time\t:\t%d\n", duration.Nanoseconds()/1e6)
	fmt.Printf("Nodes\t:\t%d\n", nodes)
	fmt.Printf("NPS\t:\t%d\n", int64(nodes)/(duration.Nanoseconds()/1e9))
}

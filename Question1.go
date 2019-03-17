package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)


var wg sync.WaitGroup

type Cell struct {
	cost int
	quantity int
	valid bool
	factory int
	warehouse int
}

type Path struct {
	cells []Cell
	cost int
}

func (p *Path) append(c Cell) {
	p.cells = append(p.cells, c)
}



func steppingStone(problem *[][] Cell, m int, n int) {

	chanSize := m*n -(m+n-1)
	result := make(chan Path,chanSize)
	wg.Add(chanSize)

	for i:=0; i < m; i++ {
		for j:=0; j < n; j++ {
			if (*problem)[i][j].quantity == 0 { //check if empty
				//fmt.Println()

				go marginalCost((*problem)[i][j], problem, m, n, result)

				//fmt.Println()


			} else {
				continue // if non-empty check next cell
			}
		}

	}

	wg.Wait()

	var optimalPath Path
	for k:=0 ; k <chanSize; k++{
		path:= <- result
		if path.cost < optimalPath.cost {
			optimalPath = path
		}
	}

	allocation := 9999999
	for s,elem := range optimalPath.cells{
		if elem.quantity < allocation && elem.quantity != 0 && s%2 !=0{
			allocation = elem.quantity
		}
	}
	if allocation == 9999999 {
		return
	}
	for p:=0; p < len(optimalPath.cells); p++{
		if p % 2 == 0 {
			(*problem)[optimalPath.cells[p].factory][optimalPath.cells[p].warehouse].quantity += allocation
		}else{
			(*problem)[optimalPath.cells[p].factory][optimalPath.cells[p].warehouse].quantity -= allocation
		}
	}

	fmt.Println(*problem)
	steppingStone(problem,m,n)

}

func marginalCost(cell Cell,problem *[][]Cell, m int, n int, result chan Path) { //ADD MARGINAL COST TO THE PATH ??
	//copy problem array
	temp := make([][]Cell, len(*problem))
	for i := range *problem {
		temp[i] = make([] Cell, len((*problem)[i]))
		copy(temp[i], (*problem)[i])
	}

	path := Path{make([]Cell, 0),0}
	path.append(cell)
	path.cost += cell.cost

	for { //break when closed
		closed := true
		for i:=0; i < m; i++ {
			for j := 0; j < n; j++ {
				if !hasNeighbours(temp, cell, temp[i][j], m, n) && temp[i][j].valid{
					temp[i][j].valid = false
					closed = false
				}
			}
		}

		if closed {
			break
		}

	}

	temp[cell.factory][cell.warehouse].valid = true
	//current := Cell{cell.cost,cell.quantity,cell.valid,cell.visited,cell.factory,cell.warehouse}
	var f int
	var w int

	f = cell.factory
	w = cell.warehouse

	for {
		completed := false
		//horizontal

		for i:=0 ; i< n ; i++{

			if temp[f][i].valid && w != i{
				path.append(temp[f][i])
				path.cost -= temp[f][i].cost
				w = i
				if w == cell.warehouse {
					completed = true
				}
				break
			}
		}

		if completed{
			break
		}

		for j := 0; j < m; j++ {
			if temp[j][w].valid && f != j {
				path.append(temp[j][w])
				path.cost += temp[j][w].cost
				f = j
				break
			}
		}

	}

	/*for i:=0; i < m; i++ {
		for j := 0; j < n; j++ {
			if temp[i][j].valid{
				fmt.Printf("%d-%d ",temp[i][j].quantity, temp[i][j].cost)
			}
		}
	}

	fmt.Println()
	for _,elem := range path.cells {
		fmt.Printf("%d-%d ",elem.quantity, elem.cost)
	}

	fmt.Println(path.cost)*/

	result <- path
	wg.Done()


}

func hasNeighbours(temp [][]Cell, start Cell, c Cell, m int, n int) bool{
	hasHorizontal := false
	hasVertical := false

	if c.quantity == 0 {
		return false
	}
	//horizontal
	if start.factory == c.factory && c.warehouse != start.warehouse {
		hasHorizontal = true
	} else {
		for i := 0; i < n; i++ {
			if temp[c.factory][i].quantity > 0 && i != c.warehouse && temp[c.factory][i].valid {
				hasHorizontal = true
				break
			}
		}
	}

	//vertical
	if start.warehouse == c.warehouse && c.factory != start.factory {
		hasVertical = true
	} else {
		for j := 0; j < m; j++ {
			if temp[j][c.warehouse].quantity > 0 && j != c.factory && temp[j][c.warehouse].valid{ //CHECK IF SAME COLUMN ??
				hasVertical = true
				break
			}
		}
	}

	return hasVertical && hasHorizontal


}



func main() {
	/*problem := [][]Cell{{Cell{6,0,true,false,0,0},
	Cell{8,25,true,false,0,1},
	Cell{10,125,true,false,0,2}},

	{Cell{7,0,true,false,1,0},
	Cell{11,0,true,false,1,1},
	Cell{11,175,true,false,1,2}},

	{Cell{4,200,true,false,2,0},
	Cell{5,75,true,false,2,1},
	Cell{12,0,true,false,2,2}} }

	bool := hasNeighbours( problem , problem[0][0],problem[2][0], 3,3)
	fmt.Println(bool)

	steppingStone(&problem,3,3)*/

	var fileName string

	fmt.Print("Enter name of cost file: ")
	fmt.Scan(&fileName)
	fmt.Println()
	file ,err := os.Open(fileName)

	if err != nil {
		panic(err)
	}

	defer file.Close()


	scanner := bufio.NewScanner(file)
	scanner.Scan()
	n := len(strings.Split(scanner.Text()," ")) - 2

	m := -1
	for scanner.Scan(){
		m++
	}


	supply := make([]int, m)
	demand := make([]int, n)
	problem := make([][]Cell,m)
	for i := range problem{
		problem[i] = make([]Cell, n)
	}

	costFile ,err := os.Open(fileName)

	if err != nil {
		panic(err)
	}
	defer costFile.Close()


	var fileName1 string
	fmt.Print("Enter name of initial solution file: ")
	fmt.Scan(&fileName1)
	fmt.Println()
	initialFile, err1 := os.Open(fileName1)

	if err1 != nil {
		panic(err1)
	}
	defer initialFile.Close()

	scannerCost := bufio.NewScanner(costFile)
	scannerCost.Scan()

	scannerInitial := bufio.NewScanner(initialFile)
	scannerInitial.Scan()

	for i:= 0; i < m ; i++{
		scannerCost.Scan()
		scannerInitial.Scan()

		row := strings.Split(scannerCost.Text(), " ")
		temp,_ := strconv.Atoi(row[ len(row) - 1 ])

		row1 := strings.Split(scannerInitial.Text(), " ")

		supply[i] = temp
		for j:= 0; j < n; j++{
			problem[i][j].cost,_ = strconv.Atoi(row[j+1])
			problem[i][j].factory = i
			problem[i][j].warehouse = j
			problem[i][j].valid = true
			qty, err := strconv.Atoi(row1[j+1])
			if err == nil {
				problem[i][j].quantity = qty
			}

		}
	}
	scannerCost.Scan()
	row := strings.Split(scannerCost.Text(), " ")
	for j:= 0; j < n; j++{
		demand[j],_ = strconv.Atoi(row[j+1])
	}


	fmt.Println(problem)

	steppingStone(&problem,m,n)






}
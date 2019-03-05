package main

import (
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

/*
1a + 2b + 3c + 4d = target. a,b,c,d
wa * xb - yc / zd = target. a,b,c,d

600 - target = fitness
*/

type Stack struct {
	operands  []int
	operators []string
}

func (stk *Stack) addOperand(v int) {
	stk.operands = append(stk.operands, v)
}

func (stk *Stack) addOperator(s string) {
	stk.operators = append(stk.operators, s)
}

func (stk *Stack) peekOperand() (int, bool) {
	if len(stk.operands) > 0 {
		value := stk.operands[0]
		return value, true
	}
	return -1, false  // todo: return something better here. Create error struct to handle these?
}

func (stk *Stack) popOperand() int {
	value := stk.operands[0]
	stk.operands = stk.operands[1:]
	return value
}

func (stk *Stack) peekOperator() (string, bool) {
	if len(stk.operators) > 0 {
		operator := stk.operators[0]
		return operator, true
	}
	return "nil", false
}

func (stk *Stack) popOperator() string {
	operator := stk.operators[0]
	stk.operators = stk.operators[1:]
	return operator
}

func (stk *Stack) operandsEmpty() bool {
	return len(stk.operands) == 0
}

func (stk *Stack) operatorsEmpty() bool {
	return len(stk.operators) == 0
}

type Chromosome struct {
	gene    []int
	fitness int
}

func (chromosome *Chromosome) calculateFitness(target int) {
	chromosome.fitness = Abs(chromosome.getGeneValue() - target)
}

func (chromosome *Chromosome) getGeneValue() int {
	//return w*chromosome.gene[0]*x*chromosome.gene[1] - y*chromosome.gene[2]/z*chromosome.gene[3]
	return -1
}

func Abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}

type Population struct {
	chromosomes []Chromosome
}

func (population *Population) initializePopulation(size, geneLen, target int) {
	index := 0
	for index < size {
		var gene []int
		for i := 0; i < geneLen; i++ {
			gene[i] = rand.Intn(target)
		}
		chromosome := Chromosome{gene: gene}
		chromosome.calculateFitness(target)
		population.chromosomes = append(population.chromosomes, chromosome)
		index++
	}
}

func (population *Population) sortPopulation() {
	sort.Slice(population.chromosomes, func(i, j int) bool {
		return population.chromosomes[i].fitness < population.chromosomes[j].fitness
	})
}

func (population *Population) getTopKChromosomes(k int) []Chromosome {
	population.sortPopulation()
	var result []Chromosome
	var index int
	for index < k {
		result = append(result, population.chromosomes[index])
		index++
	}
	return result
}

func (population *Population) crossover(topKChromosomes []Chromosome, target int) {
	var newPopulation []Chromosome
	newPopulation = append(newPopulation, topKChromosomes...)
	topKLen := len(topKChromosomes)
	noOfChildren := len(population.chromosomes) - topKLen
	index := 0
	for index < noOfChildren {
		randomParent1 := topKChromosomes[rand.Intn(topKLen)].gene
		randomParent2 := topKChromosomes[rand.Intn(topKLen)].gene
		var childGene []int
		for childIndex := 0; childIndex < len(randomParent1); childIndex++ {
			randomSwitch := rand.Intn(6)
			if randomSwitch <= 2 {
				childGene[childIndex] = randomParent1[childIndex]
			} else if randomSwitch <= 4 {
				childGene[childIndex] = randomParent2[childIndex]
			} else {
				childGene[childIndex] = rand.Intn(target)
			}
		}
		childChromosome := Chromosome{gene: childGene}
		childChromosome.calculateFitness(target)
		newPopulation = append(newPopulation, childChromosome)
		index++
	}
	population.chromosomes = newPopulation
}

// returns the gene length, target
func getGeneLen(dataArr []string) (int, int) {
	var geneLen int
	if target, err := strconv.Atoi(dataArr[len(dataArr) - 1]); err == nil {
		if len(dataArr) == 0 {
			return -1, target
		}
		geneLen = len(dataArr) / 2 // todo: use of magic number here. Add validation logic, maybe special parser?
		return geneLen, target
	} else {
		fmt.Println("An error occurred ", err)
	}
	return -1, -1
}

//1a + 4b + 5c + 50d + 6x = 200
// solves the equation with genes provided by a chromosome and returns it's value
func solveEquation(dataArr []string, gene []int, operands string) int {
	dataArr = dataArr[:len(dataArr) - 2]
	var result int
	s := new(Stack)
	fmt.Println(len(dataArr), len(gene))
	var geneCounter int
	for i := 0; i < len(dataArr); i++ {
		data := dataArr[i]
		if !strings.Contains(operands, data) {
			reg, err := regexp.Compile("[^0-9]")
			if err != nil {
				log.Fatal(err)
			}
			data := reg.ReplaceAllString(data, "")
			currentGene := gene[geneCounter]
			geneCounter++
			if operator, err := strconv.Atoi(data); err == nil {
				s.addOperand(operator * currentGene)
			}
		} else {
			s.addOperator(data)
		}
	}
	fmt.Println(s)
	return result
}

func hasHigherPrecedence(peekedOperator, currentOperator string) bool {
	if currentOperator == "*" && (peekedOperator == "+" || peekedOperator == "-") {
		return true
	}
	if currentOperator == "/" {
		return true
	}
	if currentOperator == "+" && peekedOperator == "-" {
		return true
	}
	if currentOperator == "-" && peekedOperator == "+" {
		return true
	}
	return false
}

func solve(lhs, rhs int, operand string) (int, bool) {
	switch operand {
	case "-":
		return lhs - rhs, true
	case "+":
		return lhs + rhs, true
	case "*":
		return lhs * rhs, true
	case "/":
		return lhs / rhs, true
	}
	return -1, false
}



func main() {
	//rand.Seed(time.Now().UTC().UnixNano())
	//population := new(Population)
	//population.initializePopulation(100)
	//noOfIterations := 200
	//for index := 0; index < noOfIterations; index++ {
	//	topKChromosomes := population.getTopKChromosomes(10)
	//	bestChromosome := topKChromosomes[0]
	//	fmt.Println("Best chromosome has fitness score of ", bestChromosome.fitness, " with gene ",
	//		bestChromosome.gene, " and gene value ", bestChromosome.getGeneValue())
	//	if bestChromosome.fitness == 0 {
	//		break
	//	}
	//	population.crossover(topKChromosomes)
	//}
	//reader := bufio.NewReader(os.Stdin)
	//fmt.Print("Enter text: ")
	//text, _ := reader.ReadString('\n')
	//text = strings.TrimSuffix(text, "\n")
	text := "1a + 4b + 5c + 50d + 6x = 200"
	operands := "+-/*"
	dataArr := strings.Split(text, " ")
	geneLen, target := getGeneLen(dataArr)
	fmt.Println(geneLen, target)
	var sampleGene []int
	sampleGene = append(sampleGene, 20, 5, 38, 4, 8)
	solveEquation(dataArr, sampleGene, operands)
	//s := new(Stack)
	//s.addOperator("(")
	//s.addValue(5)
	//s.addValue(50)
	//s.addValue(15)
	//fmt.Println(s)
	//s.popValue()
	//fmt.Println(s)
}

package main

import (
	"flag"
	"io/ioutil"
	"log"
	"math"
	"sort"
	"strconv"
	"strings"
)

var resistorListFile string
var Uin float64
var Uout float64
var numResults int

type Resistor struct {
	packageType string
	resistance  int
}

type UDiv struct {
	r1, r2 Resistor
}

func (uDiv UDiv) voltageDivision(uin float64) (uout float64) {
	return uin * (float64(uDiv.r2.resistance) / float64(uDiv.r1.resistance+uDiv.r2.resistance))
}

type ByDeviation []UDiv

func (uDiv ByDeviation) Len() int      { return len(uDiv) }
func (uDiv ByDeviation) Swap(i, j int) { uDiv[i], uDiv[j] = uDiv[j], uDiv[i] }
func (uDiv ByDeviation) Less(i, j int) bool {
	return math.Abs(uDiv[i].voltageDivision(Uin)-Uout) < math.Abs(uDiv[j].voltageDivision(Uin)-Uout)
}

func init() {
	flag.StringVar(&resistorListFile, "resistors", "", "File containing availiable resistors:"+
		" Each line should contain the package description and,"+
		" seperated by at least one space, the resistance in ohms.")
	flag.Float64Var(&Uin, "Uin", 1.0, "Input voltage in volts")
	flag.Float64Var(&Uout, "Uout", 0.5, "Output voltag in volts")
	flag.IntVar(&numResults, "n", -1, "Number of dividers to calculate")
	flag.Parse()
}
func main() {
	resistors, err := ioutil.ReadFile(resistorListFile)
	if err != nil {
		log.Fatal(err)
	}

	resistorsString := strings.Replace(string(resistors), "\r", "", -1)
	resistorValues := parseResistorList(resistorsString)

	var dividers []UDiv
	for _, r1 := range resistorValues {
		for _, r2 := range resistorValues {
			dividers = append(dividers, UDiv{r1, r2})
		}
	}

	sort.Sort(ByDeviation(dividers))

	if numResults == -1 || numResults > len(dividers) {
		numResults = len(dividers)
	}
	for i := 0; i < numResults; i++ {
		log.Print("Divider " + strconv.Itoa(i) + ": ")
		log.Println(dividers[i])
		log.Println(dividers[i].voltageDivision(Uin))
	}
}

func parseResistorList(list string) (resistorValues []Resistor) {
	resistorList := strings.Split(list, "\n")
	for i := 0; i < len(resistorList); i++ {
		if strings.HasPrefix(resistorList[i], "#") || strings.HasPrefix(resistorList[i], "//") || len(resistorList[i]) == 0 {
			continue
		}
		words := strings.Split(resistorList[i], " ")
		var valueString string
		var packageString string
		if len(words) == 1 {
			valueString = words[0]
		} else if len(words) >= 2 {
			valueString = words[1]
			packageString = words[0]
		}
		value, err := strconv.Atoi(valueString)
		if err != nil {
			continue
		}
		resistorValues = append(resistorValues, Resistor{packageString, value})
	}
	return resistorValues
}

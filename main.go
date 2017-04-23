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
var uIn float64
var uOut float64
var numResults int

type resistor struct {
	packageType string
	resistance  int
}

type uDiv struct {
	r1, r2 resistor
}

func (u uDiv) voltageDivision(uin float64) (uout float64) {
	return uin * (float64(u.r2.resistance) / float64(u.r1.resistance+u.r2.resistance))
}

type byDeviation []uDiv

func (u byDeviation) Len() int      { return len(u) }
func (u byDeviation) Swap(i, j int) { u[i], u[j] = u[j], u[i] }
func (u byDeviation) Less(i, j int) bool {
	return math.Abs(u[i].voltageDivision(uIn)-uOut) < math.Abs(u[j].voltageDivision(uIn)-uOut)
}

func init() {
	flag.StringVar(&resistorListFile, "resistors", "", "File containing available resistors:"+
		" Each line should contain the package description and,"+
		" separated by at least one space, the resistance in ohms.")
	flag.Float64Var(&uIn, "Uin", 1.0, "Input voltage in volts")
	flag.Float64Var(&uOut, "Uout", 0.5, "Output voltag in volts")
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

	var dividers []uDiv
	for _, r1 := range resistorValues {
		for _, r2 := range resistorValues {
			dividers = append(dividers, uDiv{r1, r2})
		}
	}

	sort.Sort(byDeviation(dividers))

	if numResults == -1 || numResults > len(dividers) {
		numResults = len(dividers)
	}
	for i := 0; i < numResults; i++ {
		log.Print("Divider " + strconv.Itoa(i) + ": ")
		log.Println(dividers[i])
		log.Println(dividers[i].voltageDivision(uIn))
	}
}

func parseResistorList(list string) (resistorValues []resistor) {
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
		resistorValues = append(resistorValues, resistor{packageString, value})
	}
	return resistorValues
}

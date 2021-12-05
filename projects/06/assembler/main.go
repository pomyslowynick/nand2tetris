package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {

	sourceFile, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println("Error reading file")
		panic(err)
	}
	defer sourceFile.Close()
	// Windows variant of slash
	lastSlash := strings.LastIndex(os.Args[1], "\\")
	compiledFileName := strings.Split(os.Args[1][lastSlash+1:], ".")[0]
	outputFile, err := os.Create(compiledFileName + ".hack")
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	w := bufio.NewWriter(outputFile)
	scanner := bufio.NewScanner(sourceFile)
	lineCounter := 0
	labelsTable := make(map[string]int)

	// First pass labels scan
	for scanner.Scan() {
		text := scanner.Text()
		newLineCounter := ParseFirstPass(text, lineCounter, labelsTable)
		lineCounter += newLineCounter
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	sourceFile2, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println("Error reading file")
		panic(err)
	}
	scanner2 := bufio.NewScanner(sourceFile2)
	symbolCounter := 16
	for scanner2.Scan() {
		text := scanner2.Text()
		val, symbolVal := Parse(text, lineCounter, labelsTable, symbolCounter)
		if val != (Instruction{}) {
			translatedCode := TranslateAssembly(val)
			_, err = fmt.Fprintf(w, "%s\n", translatedCode)
			if err != nil {
				panic(err)
			}
		}
		symbolCounter += symbolVal
	}
	w.Flush()
	if err := scanner2.Err(); err != nil {
		log.Fatal(err)
	}

}

func ParseFirstPass(line string, lineCounter int, labelsTable map[string]int) int {
	SetupTable(labelsTable)
	trimmedLine := strings.TrimSpace(line)
	if strings.HasPrefix(trimmedLine, "//") || len(trimmedLine) == 0 {
		return 0
	} else if strings.HasPrefix(trimmedLine, "(") {
		AddLabel(trimmedLine, labelsTable, lineCounter)
		return 0
	} else {
		return 1
	}
}

func AddLabel(line string, labelsTable map[string]int, lineCounter int) {
	lineNoFrontBracket := strings.SplitAfter(line, "(")[1]
	lineFinal := lineNoFrontBracket[0 : len(lineNoFrontBracket)-1]
	labelsTable[lineFinal] = lineCounter
}

func TranslateAssembly(instruction Instruction) string {
	switch instruction.instType {
	case 0:
		return TranslateAInstruction(instruction)
	case 1:
		return TranslateCInstruction(instruction)
	default:
		return "unknown"
	}
}

func TranslateCInstruction(instruction Instruction) (translated string) {
	translated += instruction.controlsBits
	translated += TranslateComp(instruction.comp)
	translated += TranslateDest(instruction.dest)
	translated += TranslateJump(instruction.jump)
	return translated
}

func TranslateComp(comp string) (translated string) {
	if strings.Index(comp, "M") != -1 {
		translated += "1"
	} else {
		translated += "0"
	}

	switch comp {
	case "0":
		translated += "101010"
	case "1":
		translated += "111111"
	case "-1":
		translated += "111010"
	case "D":
		translated += "001100"
	case "A":
		fallthrough
	case "M":
		translated += "110000"
	case "!D":
		translated += "001111"
	case "!A":
		fallthrough
	case "!M":
		translated += "110001"
	case "-D":
		translated += "001111"
	case "-A":
		fallthrough
	case "-M":
		translated += "110011"
	case "D+1":
		translated += "011111"
	case "A+1":
		fallthrough
	case "M+1":
		translated += "110111"
	case "D-1":
		translated += "001110"
	case "A-1":
		fallthrough
	case "M-1":
		translated += "110010"
	case "D+A":
		fallthrough
	case "D+M":
		translated += "000010"
	case "D-A":
		fallthrough
	case "D-M":
		translated += "010011"
	case "A-D":
		fallthrough
	case "M-D":
		translated += "000111"
	case "D&A":
		fallthrough
	case "D&M":
		translated += "000000"
	case "D|A":
		fallthrough
	case "D|M":
		translated += "010101"
	}

	return translated
}

func TranslateDest(dest string) (translated string) {
	switch dest {
	default:
		translated += "000"
	case "M":
		translated += "001"
	case "D":
		translated += "010"
	case "MD":
		translated += "011"
	case "A":
		translated += "100"
	case "AM":
		translated += "101"
	case "AD":
		translated += "110"
	case "AMD":
		translated += "111"
	}

	return translated
}

func TranslateJump(jump string) (translated string) {
	switch jump {
	default:
		translated += "000"
	case "JGT":
		translated += "001"
	case "JEQ":
		translated += "010"
	case "JGE":
		translated += "011"
	case "JLT":
		translated += "100"
	case "JNE":
		translated += "101"
	case "JLE":
		translated += "110"
	case "JMP":
		translated += "111"
	}

	return translated
}

func TranslateAInstruction(instruction Instruction) (translated string) {
	translated += "0"
	s := strconv.FormatInt(int64(instruction.value), 2)

	for i := 0; i < 15-len(s); i++ {
		translated += "0"
	}
	translated += s
	return translated
}

type Instruction struct {
	instType     int
	controlsBits string
	comp         string
	dest         string
	jump         string
	value        int
}

func SetupTable(labelsTable map[string]int) {
	for i := 0; i < 16; i++ {
		labelsTable["R"+strconv.Itoa(i)] = i
	}

	labelsTable["SCREEN"] = 16384
	labelsTable["KBD"] = 24576
	labelsTable["SP"] = 0
	labelsTable["LCL"] = 1
	labelsTable["ARG"] = 2
	labelsTable["THIS"] = 3
	labelsTable["THAT"] = 4
}

func Parse(line string, lineCounter int, labelsTable map[string]int, symbolCounter int) (Instruction, int) {
	trimmedLine := strings.TrimSpace(line)
	if strings.HasPrefix(trimmedLine, "//") || len(trimmedLine) == 0 {
		return Instruction{}, 0
	} else if strings.HasPrefix(trimmedLine, "@") {
		instruction, symbolVal := ParseAinstruction(trimmedLine, labelsTable, symbolCounter)
		return instruction, symbolVal
	} else if strings.HasPrefix(trimmedLine, "(") {
		return Instruction{}, 0
	} else {
		return ParseCinstruction(trimmedLine), 0
	}
}

func ParseCinstruction(line string) Instruction {
	instruction := new(Instruction)
	instruction.controlsBits = "111"
	instruction.instType = 1
	splitInLineComment := strings.Split(line, "//")[0]
	dest := strings.Split(splitInLineComment, "=")
	if len(dest) == 2 {
		instruction.dest = strings.TrimSpace(dest[0])
	}

	jump := strings.SplitAfter(splitInLineComment, ";")
	if len(jump) == 2 {
		instruction.jump = strings.TrimSpace(jump[1])
	}

	if len(dest) == 2 && len(jump) == 2 {
		instruction.comp = strings.SplitAfter(dest[1], ";")[0]
	} else if len(dest) == 2 {
		instruction.comp = strings.TrimSpace(dest[1])
	} else if len(jump) == 2 {
		tempComp := strings.TrimSpace(jump[0])
		instruction.comp = tempComp[:len(tempComp)-1]
	} else {
		instruction.comp = splitInLineComment
	}
	return *instruction
}

func ParseAinstruction(line string, labelsTable map[string]int, symbolCounter int) (Instruction, int) {
	var symbolVal int
	instruction := new(Instruction)
	instruction.instType = 0
	instruction.controlsBits = "0"
	instrValue := strings.SplitAfter(line, "@")[1]
	instrValue = strings.TrimSpace(instrValue)
	val, err := strconv.Atoi(instrValue)
	if err != nil {
		symbolVal = AddSymbol(instrValue, labelsTable, symbolCounter)
		instruction.value = labelsTable[instrValue]
	} else {
		instruction.value = val
	}
	return *instruction, symbolVal
}

func AddSymbol(line string, labelsTable map[string]int, symbolCounter int) int {
	if _, ok := labelsTable[line]; ok {
		return 0
	}
	labelsTable[line] = symbolCounter
	return 1
}

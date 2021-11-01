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

	compiledFileName := strings.Split(os.Args[1], ".")[0]
	outputFile, err := os.Create(compiledFileName + ".hack")
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	w := bufio.NewWriter(outputFile)
	scanner := bufio.NewScanner(sourceFile)
	for scanner.Scan() {
		text := scanner.Text()
		fmt.Println(text)
		val := Parse(text)
		if val != (Instruction{}) {
			fmt.Printf("%+v\n", val)
			translatedCode := TranslateAssembly(val)
			_, err = fmt.Fprintf(w, "%s\n", translatedCode)
			if err != nil {
				panic(err)
			}
		}
	}
	w.Flush()

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

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
	fmt.Println(instruction.value)
	fmt.Println(int64(instruction.value))
	s := strconv.FormatInt(int64(instruction.value), 2)
	fmt.Println(s)

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

func Parse(line string) Instruction {
	trimmedLine := strings.TrimSpace(line)
	if strings.HasPrefix(line, "//") || len(trimmedLine) == 0 {
		return Instruction{}
	} else if strings.HasPrefix(line, "@") {
		return ParseAinstruction(line)
	} else if strings.HasPrefix(line, "(") {
		return ParseSymbol(line)
	} else {
		return ParseCinstruction(line)
	}
}

func ParseCinstruction(line string) Instruction {
	instruction := new(Instruction)
	instruction.controlsBits = "111"
	instruction.instType = 1
	splitInLineComment := strings.SplitAfter(line, "//")[0]
	dest := strings.SplitAfter(splitInLineComment, "=")
	if len(dest) == 2 {
		instruction.dest = strings.TrimSpace(dest[0])[:len(dest)-1]
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
		instruction.comp = strings.TrimSpace(jump[1])
	} else {
		instruction.comp = splitInLineComment
	}
	return *instruction
}

func ParseAinstruction(line string) Instruction {
	instruction := new(Instruction)
	instruction.instType = 0
	instruction.controlsBits = "0"
	val, err := strconv.Atoi(strings.SplitAfter(line, "@")[1])
	if err != nil {
		fmt.Println("Error converting to A instruction value to Integer")
	}
	instruction.value = val
	return *instruction
}

func ParseSymbol(line string) Instruction {
	return Instruction{}
}

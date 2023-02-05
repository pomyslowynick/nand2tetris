package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// type algebraCommand struct {
// 	command string
// }
type VMCommand struct {
	commandType  string
	command      string
	memoryRegion string
	value        string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Provide the filename for translation")
		os.Exit(1)
	}
	var vmCommands = make([]VMCommand, 0)

	ParseFile(&vmCommands)
	// fmt.Println(vmCommands)
	WriteCode(&vmCommands)
}

func WriteCode(vmCommands *[]VMCommand) {
	// Windows variant of slash
	lastSlash := strings.LastIndex(os.Args[1], "\\")
	compiledFileName := strings.Split(os.Args[1][lastSlash+1:], ".")[0]
	outputFile, err := os.Create(compiledFileName + ".asm")
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	w := bufio.NewWriter(outputFile)
	for _, command := range *vmCommands {

		_, err = fmt.Fprintf(w, "%s\n", TranslateCommandToAssembly(command))
		if err != nil {
			panic(err)
		}
	}
	w.Flush()

}

func TranslateCommandToAssembly(vmCommand VMCommand) (translatedCommand string) {
	translatedCommand += fmt.Sprintf("// %s %s %s \n", vmCommand.command, vmCommand.memoryRegion, vmCommand.value)
	switch vmCommand.commandType {
	case "memory":
		translatedCommand += TranslateMemoryVMCommand(vmCommand)
	case "algebra":
		translatedCommand += TranslateAlgebraVMCommand(vmCommand)
	}
	return
}

func TranslateMemoryVMCommand(vmCommand VMCommand) (translatedCommand string) {
	switch vmCommand.command {
	case "push":
		translatedCommand += TranslatePushMemoryRegion(vmCommand)
	case "pop":
		translatedCommand += TranslatePopMemoryRegion(vmCommand)
	default:
		translatedCommand += "ERROR: Unknown Memory Command"
	}
	return
}

func TranslatePushMemoryRegion(vmCommand VMCommand) string {
	switch vmCommand.memoryRegion {
	case "local":
		return TranslatePushLATT("LCL", vmCommand.value)
	case "argument":
		return TranslatePushLATT("ARG", vmCommand.value)
	case "this":
		return TranslatePushLATT("THIS", vmCommand.value)
	case "that":
		return TranslatePushLATT("THAT", vmCommand.value)
	case "constant":
		return TranslatePushConstant(vmCommand.value)
	case "static":
		return TranslatePushStatic(vmCommand.value)
	case "temp":
		return TranslatePushTemp(vmCommand.value)
	case "pointer":
		return TranslatePushPointer(vmCommand.value)
	}

	return ""
}

func TranslatePopMemoryRegion(vmCommand VMCommand) string {
	switch vmCommand.memoryRegion {
	case "local":
		return TranslatePopLATT("LCL", vmCommand.value)
	case "argument":
		return TranslatePopLATT("ARG", vmCommand.value)
	case "this":
		return TranslatePopLATT("THIS", vmCommand.value)
	case "that":
		return TranslatePopLATT("THAT", vmCommand.value)
	case "static":
		return TranslatePopStatic(vmCommand.value)
	case "temp":
		return TranslatePopTemp(vmCommand.value)
	case "pointer":
		return TranslatePopPointer(vmCommand.value)
	}

	return ""
}

func TranslatePushPointer(value string) (translatedCommand string) {

	// Pop stack value on the memory region
	if value == "0" {
		translatedCommand += "@THIS\n"
	} else if value == "1" {
		translatedCommand += "@THAT\n"
	}
	translatedCommand += "D=M\n"
	translatedCommand += "@SP\n"
	translatedCommand += "A=M\n"
	translatedCommand += "M=D\n"

	// Decrement the stack pointer
	translatedCommand += "@SP\n"
	translatedCommand += "M=M+1\n"
	return
}

func TranslatePopPointer(value string) (translatedCommand string) {
	// Decrement the stack pointer
	translatedCommand += "@SP\n"
	translatedCommand += "M=M-1\n"

	// Pop stack value on the memory region
	translatedCommand += "@SP\n"
	translatedCommand += "A=M\n"
	translatedCommand += "D=M\n"
	if value == "0" {
		translatedCommand += "@THIS\n"
	} else if value == "1" {
		translatedCommand += "@THAT\n"
	}
	translatedCommand += "M=D\n"
	return

}

func TranslatePopLATT(memoryRegion string, value string) (translatedCommand string) {
	// Add the value to the memoryRegion pointer
	translatedCommand += fmt.Sprintf("@%s\n", value)
	translatedCommand += "D=A\n"
	translatedCommand += fmt.Sprintf("@%s\n", memoryRegion)
	translatedCommand += "A=M\n"
	translatedCommand += "D=D+A\n"
	translatedCommand += "@R13\n"
	translatedCommand += "M=D\n"

	// Decrement the stack pointer
	translatedCommand += "@SP\n"
	translatedCommand += "M=M-1\n"

	// Pop stack value on the memory region
	translatedCommand += "A=M\n"
	translatedCommand += "D=M\n"
	translatedCommand += "@R13\n"
	translatedCommand += "A=M\n"
	translatedCommand += "M=D\n"
	return
}

func TranslatePopStatic(value string) (translatedCommand string) {
	// Decrement the stack pointer
	translatedCommand += "@SP\n"
	translatedCommand += "M=M-1\n"

	// Pop stack value on the memory region
	translatedCommand += "A=M\n"
	translatedCommand += "D=M\n"
	translatedCommand += fmt.Sprintf("@Foo.%s\n", value)
	translatedCommand += "M=D\n"
	return
}

func TranslatePopTemp(value string) (translatedCommand string) {
	convertedValue, err := strconv.Atoi(value)
	if err != nil {
		panic(err)
	}
	memoryRegion := fmt.Sprintf("%s", strconv.Itoa(5+convertedValue))

	// Decrement the stack pointer
	translatedCommand += "@SP\n"
	translatedCommand += "M=M-1\n"

	// Pop stack value on the memory region
	translatedCommand += "A=M\n"
	translatedCommand += "D=M\n"
	translatedCommand += fmt.Sprintf("@%s\n", memoryRegion)
	translatedCommand += "M=D\n"
	return
}

func TranslatePushLATT(memoryRegion string, value string) (translatedCommand string) {
	// Add the value to the memoryRegion pointer
	translatedCommand += fmt.Sprintf("@%s\n", value)
	translatedCommand += "D=A\n"
	translatedCommand += fmt.Sprintf("@%s\n", memoryRegion)
	translatedCommand += "D=D+M\n"

	// Push memoryRegion value on the stack
	translatedCommand += "A=D\n"
	translatedCommand += "D=M\n"
	translatedCommand += "@SP\n"
	translatedCommand += "A=M\n"
	translatedCommand += "M=D\n"

	// Increment the stack pointer
	translatedCommand += "@SP\n"
	translatedCommand += "M=M+1\n"
	return
}

func TranslatePushStatic(value string) (translatedCommand string) {
	translatedCommand += fmt.Sprintf("@Foo.%s\n", value)
	translatedCommand += "A=M\n"
	translatedCommand += "D=M\n"
	translatedCommand += "@SP\n"
	translatedCommand += "A=M\n"
	translatedCommand += "M=D\n"
	translatedCommand += "@SP\n"
	translatedCommand += "M=M+1\n"
	return
}

func TranslatePushTemp(value string) (translatedCommand string) {
	convertedValue, err := strconv.Atoi(value)
	if err != nil {
		panic(err)
	}
	tempAddr := fmt.Sprintf("%s", strconv.Itoa(5+convertedValue))

	translatedCommand += fmt.Sprintf("@%s\n", tempAddr)
	// translatedCommand += "A=M\n"
	translatedCommand += "D=M\n"
	translatedCommand += "@SP\n"
	translatedCommand += "A=M\n"
	translatedCommand += "M=D\n"
	translatedCommand += "@SP\n"
	translatedCommand += "M=M+1\n"
	return
}

func TranslatePushConstant(value string) (translatedCommand string) {
	translatedCommand += fmt.Sprintf("@%s\n", value)
	translatedCommand += "D=A\n"
	translatedCommand += "@SP\n"
	translatedCommand += "A=M\n"
	translatedCommand += "M=D\n"
	translatedCommand += "@SP\n"
	translatedCommand += "M=M+1\n"
	return
}

func TranslateAlgebraVMCommand(vmCommand VMCommand) (translatedCommand string) {

	switch vmCommand.command {
	case "add":
		// Pop the first value into D
		translatedCommand += "@SP\n"
		translatedCommand += "M=M-1\n"
		translatedCommand += "A=M\n"
		translatedCommand += "D=M\n"

		// Add D to SP-1
		translatedCommand += "@SP\n"
		translatedCommand += "M=M-1\n"
		translatedCommand += "A=M\n"
		translatedCommand += "D=D+M\n"
		translatedCommand += "M=D\n"
		// Increment the stack pointer
		translatedCommand += "@SP\n"
		translatedCommand += "M=M+1\n"
	case "sub":
		// Pop the first value into D
		translatedCommand += "@SP\n"
		translatedCommand += "M=M-1\n"
		translatedCommand += "A=M\n"
		translatedCommand += "D=M\n"

		// Add D to SP-1
		translatedCommand += "@SP\n"
		translatedCommand += "M=M-1\n"
		translatedCommand += "A=M\n"
		translatedCommand += "D=M-D\n"
		translatedCommand += "M=D\n"
		// Increment the stack pointer
		translatedCommand += "@SP\n"
		translatedCommand += "M=M+1\n"
	case "neg":
		// Pop the first value into D
		translatedCommand += "@SP\n"
		translatedCommand += "M=M-1\n"
		translatedCommand += "A=M\n"
		translatedCommand += "M=-M\n"

		// Increment the stack pointer
		translatedCommand += "@SP\n"
		translatedCommand += "M=M+1\n"
	case "eq":
		// Pop the first value into D
		translatedCommand += "@SP\n"
		translatedCommand += "M=M-1\n"
		translatedCommand += "A=M\n"
		translatedCommand += "D=M\n"

		// Add D to SP-1
		translatedCommand += "@SP\n"
		translatedCommand += "M=M-1\n"
		translatedCommand += "A=M\n"
		translatedCommand += "D=D&M\n"
		translatedCommand += "M=D\n"
		// Increment the stack pointer
		translatedCommand += "@SP\n"
		translatedCommand += "M=M+1\n"
	case "gt":
		// Pop the first value into D
		translatedCommand += "@SP\n"
		translatedCommand += "M=M-1\n"
		translatedCommand += "A=M\n"
		translatedCommand += "D=M\n"

		// Add D to SP-1
		translatedCommand += "@SP\n"
		translatedCommand += "M=M-1\n"
		translatedCommand += "A=M\n"
		translatedCommand += "D=DM\n"
		translatedCommand += "M=D\n"
		// Increment the stack pointer
		translatedCommand += "@SP\n"
		translatedCommand += "M=M+1\n"
	case "lt":
	case "and":
		// Pop the first value into D
		translatedCommand += "@SP\n"
		translatedCommand += "M=M-1\n"
		translatedCommand += "A=M\n"
		translatedCommand += "D=M\n"

		translatedCommand += "@SP\n"
		translatedCommand += "M=M-1\n"
		translatedCommand += "A=M\n"
		translatedCommand += "D=D&M\n"
		translatedCommand += "M=D\n"
		// Increment the stack pointer
		translatedCommand += "@SP\n"
		translatedCommand += "M=M+1\n"
	case "or":
		// Pop the first value into D
		translatedCommand += "@SP\n"
		translatedCommand += "M=M-1\n"
		translatedCommand += "A=M\n"
		translatedCommand += "D=M\n"

		translatedCommand += "@SP\n"
		translatedCommand += "M=M-1\n"
		translatedCommand += "A=M\n"
		translatedCommand += "D=D|M ; \n"
		translatedCommand += "M=D\n"
		// Increment the stack pointer
		translatedCommand += "@SP\n"
		translatedCommand += "M=M+1\n"
	case "not":
		// Pop the first value into D
		translatedCommand += "@SP\n"
		translatedCommand += "M=M-1\n"
		translatedCommand += "A=M\n"
		translatedCommand += "M=!M\n"

		// Increment the stack pointer
		translatedCommand += "@SP\n"
		translatedCommand += "M=M+1\n"
	default:
		return translatedCommand + "Unknown Command"
	}
	return
}

func ParseFile(vmCommands *[]VMCommand) []VMCommand {
	// Open the file to process
	sourceFile, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println("Error reading file")
		panic(err)
	}
	scanner := bufio.NewScanner(sourceFile)

	for scanner.Scan() {
		text := scanner.Text()
		parseCommand := ParseLine(text)
		if parseCommand != (VMCommand{}) {
			*vmCommands = append(*vmCommands, parseCommand)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	defer sourceFile.Close()
	return *vmCommands
}

func ParseLine(line string) VMCommand {
	if strings.TrimSpace(line) == "" {
		return VMCommand{}
	} else if strings.HasPrefix(strings.TrimLeft(line, " "), "//") {
		return VMCommand{}
	} else {
		return GetVMCommand(line)
	}
}

func GetVMCommand(line string) VMCommand {
	splitLine := strings.Split(line, " ")
	switch splitLine[0] {
	case "push":
		return VMCommand{commandType: "memory", command: "push", memoryRegion: splitLine[1], value: splitLine[2]}
	case "pop":
		return VMCommand{commandType: "memory", command: "pop", memoryRegion: splitLine[1], value: splitLine[2]}
	case "add":
		return VMCommand{commandType: "algebra", command: "add", memoryRegion: "", value: ""}
	case "sub":
		return VMCommand{commandType: "algebra", command: "sub", memoryRegion: "", value: ""}
	case "eq":
		return VMCommand{commandType: "algebra", command: "eq", memoryRegion: "", value: ""}
	default:
		return VMCommand{commandType: "Unknown"}
	}
}
func WritePushPop(line string) string {
	return "Placeholder"
}

func WriteArithmetic(line string) string {
	return "Placeholder"
}

func AddLabel(line string, labelsTable map[string]int, lineCounter int) {
	lineNoFrontBracket := strings.SplitAfter(line, "(")[1]
	lineFinal := lineNoFrontBracket[0 : len(lineNoFrontBracket)-1]
	labelsTable[lineFinal] = lineCounter
}

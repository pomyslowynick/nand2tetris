package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
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
	file         string
	name         string
}

var LabelGlobalVar = 0
var DirectoryName = ""
var CurrentFunction = "Sys.init"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Provide the filename for translation")
		os.Exit(1)
	}
	var vmCommands = make([]VMCommand, 0)

	ParseAndOpenVMFile(&vmCommands)
	WriteCode(&vmCommands)
}

func ParseAndOpenVMFile(vmCommands *[]VMCommand) []VMCommand {
	// Open the file to process
	sourceFile, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println("Error reading file")
		panic(err)
	}
	defer sourceFile.Close()

	fileOrDir, err := sourceFile.Stat()
	if err != nil {
		fmt.Println("Error reading file")
		panic(err)
	}

	if fileOrDir.IsDir() {
		fmt.Println("[Parser] Traversing a directory", fileOrDir.Name())
		fmt.Println("[Parser] VM files found:")

		// Set the name of the output file, I know this shouldn't be set from here
		// Fix if possible, should set it with a dedicated method
		DirectoryName = fileOrDir.Name()

		items, _ := ioutil.ReadDir(os.Args[1])
		for _, item := range items {
			if strings.Split(item.Name(), ".")[1] == "vm" {
				fmt.Println("\t", item.Name())

				vmSourceFile, err := os.Open(os.Args[1] + item.Name())
				if err != nil {
					fmt.Println("Error reading file")
					panic(err)
				}
				filename := strings.Split(item.Name(), ".")[0]
				scanner := bufio.NewScanner(vmSourceFile)

				for scanner.Scan() {
					text := scanner.Text()
					parseCommand := ParseLine(text, filename)
					if parseCommand != (VMCommand{}) {
						*vmCommands = append(*vmCommands, parseCommand)
					}
				}
				if err := scanner.Err(); err != nil {
					log.Fatal(err)
				}
				defer vmSourceFile.Close()
				return *vmCommands
			}
		}
		return *vmCommands
	} else {
		scanner := bufio.NewScanner(sourceFile)
		for scanner.Scan() {
			text := scanner.Text()
			filename := strings.Split(sourceFile.Name(), ".")[0]
			parseCommand := ParseLine(text, filename)
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
}

func GetOutputFileName() string {
	if DirectoryName == "" {
		lastSlash := strings.LastIndex(os.Args[1], "\\")
		return strings.Split(os.Args[1][lastSlash+1:], ".")[0]
	} else {
		return DirectoryName
	}
}

func WriteCode(vmCommands *[]VMCommand) {
	// Windows variant of slash
	outputFileName := GetOutputFileName()
	outputFile, err := os.Create(outputFileName + ".asm")

	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	w := bufio.NewWriter(outputFile)
	// _, err = fmt.Fprintf(w, "%s\n", WriteInit())
	if err != nil {
		panic(err)
	}

	for _, command := range *vmCommands {

		_, err = fmt.Fprintf(w, "%s\n", TranslateCommandToAssembly(command))
		if err != nil {
			panic(err)
		}
	}
	w.Flush()

}

func WriteInit() (translatedCommand string) {
	// Set SP to 256
	translatedCommand += "@256\n"
	translatedCommand += "A=D\n"
	translatedCommand += "@SP\n"
	translatedCommand += "M=D\n"

	// Call Sys.init
	translatedCommand += "M=D\n"
	return
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
	case "label":
		translatedCommand += TranslateLabelCommand(vmCommand)
	case "goto":
		translatedCommand += TranslateGotoCommand(vmCommand)
	case "if-goto":
		translatedCommand += TranslateIfGotoCommand(vmCommand)
	case "function":
		translatedCommand += TranslateFunctionCommand(vmCommand)
	case "call":
		translatedCommand += TranslateCallCommand(vmCommand)
	case "return":
		translatedCommand += TranslateReturnCommand(vmCommand)
	default:
		translatedCommand += "ERROR: Unknown Memory Command"
	}
	return
}

func TranslateCallCommand(vmCommand VMCommand) (translatedCommand string) {
	// Set return address to the top of the stack
	translatedCommand += fmt.Sprintf("@%s\n", fmt.Sprintf("RETURN_%s", CurrentFunction))
	translatedCommand += "D=A\n"
	translatedCommand += "@SP\n"
	translatedCommand += "A=M\n"
	translatedCommand += "M=D\n"
	translatedCommand += "@SP\n"
	translatedCommand += "M=M+1\n"

	// Save the LCL pointer value
	translatedCommand += "@LCL\n"
	translatedCommand += "A=M\n"
	translatedCommand += "D=A\n"
	translatedCommand += "@SP\n"
	translatedCommand += "A=M\n"
	translatedCommand += "M=D\n"
	translatedCommand += "@SP\n"
	translatedCommand += "M=M+1\n"

	// Save ARG pointer value
	translatedCommand += "@ARG\n"
	translatedCommand += "A=M\n"
	translatedCommand += "D=A\n"
	translatedCommand += "@SP\n"
	translatedCommand += "A=M\n"
	translatedCommand += "M=D\n"
	translatedCommand += "@SP\n"
	translatedCommand += "M=M+1\n"

	// Save the THIS pointer value
	translatedCommand += "@THIS\n"
	translatedCommand += "A=M\n"
	translatedCommand += "D=A\n"
	translatedCommand += "@SP\n"
	translatedCommand += "A=M\n"
	translatedCommand += "M=D\n"
	translatedCommand += "@SP\n"
	translatedCommand += "M=M+1\n"

	// Save the THAT pointer value
	translatedCommand += "@THAT\n"
	translatedCommand += "A=M\n"
	translatedCommand += "D=A\n"
	translatedCommand += "@SP\n"
	translatedCommand += "A=M\n"
	translatedCommand += "M=D\n"
	translatedCommand += "@SP\n"
	translatedCommand += "M=M+1\n"

	// Set the ARG value
	translatedCommand += "@SP\n"
	translatedCommand += "A=M\n"
	translatedCommand += "D=A\n"
	translatedCommand += fmt.Sprintf("@%s\n", vmCommand.value)
	translatedCommand += "D=D-A\n"
	translatedCommand += "@5\n"
	translatedCommand += "D=D-A\n"
	translatedCommand += "@ARG\n"
	translatedCommand += "M=D\n"

	// Set the LCL value
	translatedCommand += "@SP\n"
	translatedCommand += "A=M\n"
	translatedCommand += "D=A\n"
	translatedCommand += "@LCL\n"
	translatedCommand += "M=D\n"

	// Jump to the function
	translatedCommand += fmt.Sprintf("@%s\n", vmCommand.name)
	translatedCommand += "0;JMP\n"

	// Save the return address
	translatedCommand += fmt.Sprintf("(RETURN_%s)", CurrentFunction)
	CurrentFunction = fmt.Sprintf("%s", vmCommand.name)

	return
}

func TranslateReturnCommand(vmCommand VMCommand) (translatedCommand string) {

	// Store LCL in GPR 14
	translatedCommand += "@LCL\n"
	translatedCommand += "A=M\n"
	translatedCommand += "D=A\n"
	translatedCommand += "@14\n"
	translatedCommand += "M=D\n"

	// Get the return address
	translatedCommand += "@14\n"
	translatedCommand += "A=M\n"
	translatedCommand += "D=A\n"
	translatedCommand += "@5\n"
	translatedCommand += "D=D-A\n"
	translatedCommand += "A=D\n"
	translatedCommand += "D=M\n"

	// Storing the retAddr in GPR 15
	translatedCommand += "@15\n"
	translatedCommand += "M=D\n"

	// Set ARG to be the return value
	translatedCommand += "@SP\n"
	translatedCommand += "M=M-1\n"
	translatedCommand += "@SP\n"
	translatedCommand += "A=M\n"
	translatedCommand += "D=M\n"
	translatedCommand += "@ARG\n"
	translatedCommand += "A=M\n"
	translatedCommand += "M=D\n"

	// Restore SP to before the call state
	translatedCommand += "@ARG\n"
	translatedCommand += "A=M\n"
	translatedCommand += "D=A+1\n"
	translatedCommand += "@SP\n"
	translatedCommand += "M=D\n"

	// Restore the memory region pointers
	// Restore THAT
	translatedCommand += "@14\n"
	translatedCommand += "A=M-1\n"
	translatedCommand += "D=M\n"
	translatedCommand += "@THAT\n"
	translatedCommand += "M=D\n"

	// Restore THIS
	translatedCommand += "@14\n"
	translatedCommand += "A=M-1\n"
	translatedCommand += "A=A-1\n"
	translatedCommand += "D=M\n"
	translatedCommand += "@THIS\n"
	translatedCommand += "M=D\n"

	// Restore ARG
	translatedCommand += "@14\n"
	translatedCommand += "A=M-1\n"
	translatedCommand += "A=A-1\n"
	translatedCommand += "A=A-1\n"
	translatedCommand += "D=M\n"
	translatedCommand += "@ARG\n"
	translatedCommand += "M=D\n"

	// Restore LCL
	translatedCommand += "@14\n"
	translatedCommand += "A=M-1\n"
	translatedCommand += "A=A-1\n"
	translatedCommand += "A=A-1\n"
	translatedCommand += "A=A-1\n"
	translatedCommand += "D=M\n"
	translatedCommand += "@LCL\n"
	translatedCommand += "M=D\n"

	// Jump back to caller
	translatedCommand += "@15\n"
	translatedCommand += "A=M\n"
	translatedCommand += "0;JMP\n"
	return
}

func TranslateFunctionCommand(vmCommand VMCommand) (translatedCommand string) {
	translatedCommand += fmt.Sprintf("(%s)\n", vmCommand.name)
	noLocalVars, err := strconv.Atoi(vmCommand.value)
	if err != nil {
		panic(err)
	}
	for i := 0; i < noLocalVars; i++ {
		translatedCommand += TranslatePushConstant("0")
	}

	return
}

func TranslateIfGotoCommand(vmCommand VMCommand) (translatedCommand string) {
	// Decrement the stack pointer
	translatedCommand += "@SP\n"
	translatedCommand += "M=M-1\n"

	// Pop stack value
	translatedCommand += "@SP\n"
	translatedCommand += "A=M\n"
	translatedCommand += "D=M\n"

	translatedCommand += fmt.Sprintf("@%s\n", vmCommand.value)
	translatedCommand += "D;JGT\n"
	return
}
func TranslateGotoCommand(vmCommand VMCommand) (translatedCommand string) {
	translatedCommand += fmt.Sprintf("@%s\n", vmCommand.value)
	translatedCommand += "0;JMP\n"
	return
}

func TranslateLabelCommand(vmCommand VMCommand) string {
	return fmt.Sprintf("(%s)\n", vmCommand.value)
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
	LabelGlobalVar += 1
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
		translatedCommand += "D=D-M\n"
		translatedCommand += fmt.Sprintf("@ARE_EQUAL_%d\n", LabelGlobalVar)
		translatedCommand += "D;JEQ\n"
		translatedCommand += fmt.Sprintf("@NOT_EQUAL_%d\n", LabelGlobalVar)
		translatedCommand += "0;JMP\n"

		// Both stack entries are equal
		translatedCommand += fmt.Sprintf("(ARE_EQUAL_%d)\n", LabelGlobalVar)
		translatedCommand += "@SP\n"
		translatedCommand += "A=M\n"
		translatedCommand += "M=-1\n"
		translatedCommand += fmt.Sprintf("@END_EQUAL_%d\n", LabelGlobalVar)
		translatedCommand += "0;JMP\n"
		// Increment the stack pointer
		translatedCommand += fmt.Sprintf("(NOT_EQUAL_%d)\n", LabelGlobalVar)
		translatedCommand += "@SP\n"
		translatedCommand += "A=M\n"
		translatedCommand += "M=0\n"

		translatedCommand += fmt.Sprintf("(END_EQUAL_%d)\n", LabelGlobalVar)
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
		translatedCommand += "D=M-D\n"
		translatedCommand += fmt.Sprintf("@IS_GREATER_%d\n", LabelGlobalVar)
		translatedCommand += "D;JGT\n"
		translatedCommand += fmt.Sprintf("@NOT_GREATER_%d\n", LabelGlobalVar)
		translatedCommand += "0;JMP\n"

		// Both stack entries are equal
		translatedCommand += fmt.Sprintf("(IS_GREATER_%d)\n", LabelGlobalVar)
		translatedCommand += "@SP\n"
		translatedCommand += "A=M\n"
		translatedCommand += "M=-1\n"
		translatedCommand += fmt.Sprintf("@END_GREATER_%d\n", LabelGlobalVar)
		translatedCommand += "0;JMP\n"
		// Increment the stack pointer
		translatedCommand += fmt.Sprintf("(NOT_GREATER_%d)\n", LabelGlobalVar)
		translatedCommand += "@SP\n"
		translatedCommand += "A=M\n"
		translatedCommand += "M=0\n"

		translatedCommand += fmt.Sprintf("(END_GREATER_%d)\n", LabelGlobalVar)
		translatedCommand += "@SP\n"
		translatedCommand += "M=M+1\n"
	case "lt":
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
		translatedCommand += fmt.Sprintf("@IS_LESSER_%d\n", LabelGlobalVar)
		translatedCommand += "D;JLT\n"
		translatedCommand += fmt.Sprintf("@NOT_LESSER_%d\n", LabelGlobalVar)
		translatedCommand += "0;JMP\n"

		// Both stack entries are equal
		translatedCommand += fmt.Sprintf("(IS_LESSER_%d)\n", LabelGlobalVar)
		translatedCommand += "@SP\n"
		translatedCommand += "A=M\n"
		translatedCommand += "M=-1\n"
		translatedCommand += fmt.Sprintf("@END_LESSER_%d\n", LabelGlobalVar)
		translatedCommand += "0;JMP\n"
		// Increment the stack pointer
		translatedCommand += fmt.Sprintf("(NOT_LESSER_%d)\n", LabelGlobalVar)
		translatedCommand += "@SP\n"
		translatedCommand += "A=M\n"
		translatedCommand += "M=0\n"

		translatedCommand += fmt.Sprintf("(END_LESSER_%d)\n", LabelGlobalVar)
		translatedCommand += "@SP\n"
		translatedCommand += "M=M+1\n"
	case "and":
		translatedCommand += "@SP\n"
		translatedCommand += "M=M-1\n"
		translatedCommand += "A=M\n"
		translatedCommand += "D=M\n"

		translatedCommand += "@SP\n"
		translatedCommand += "M=M-1\n"
		translatedCommand += "A=M\n"
		translatedCommand += "D=M&D\n"
		translatedCommand += "@SP\n"
		translatedCommand += "A=M\n"
		translatedCommand += "M=D\n"

		translatedCommand += "@SP\n"
		translatedCommand += "M=M+1\n"

	case "or":
		translatedCommand += "@SP\n"
		translatedCommand += "M=M-1\n"
		translatedCommand += "A=M\n"
		translatedCommand += "D=M\n"

		translatedCommand += "@SP\n"
		translatedCommand += "M=M-1\n"
		translatedCommand += "A=M\n"
		translatedCommand += "D=M|D\n"
		translatedCommand += "@SP\n"
		translatedCommand += "A=M\n"
		translatedCommand += "M=D\n"

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

func ParseLine(line string, file string) VMCommand {
	if strings.TrimSpace(line) == "" {
		return VMCommand{}
	} else if strings.HasPrefix(strings.TrimLeft(line, " "), "//") {
		return VMCommand{}
	} else {
		return GetVMCommand(line, file)
	}
}

func DiscardComments(splitLine string) string {
	if strings.Contains(splitLine, "//") {
		splitLine = strings.Split(splitLine, "//")[0]
		return strings.Trim(splitLine, " \t")
	}
	return splitLine
}
func GetVMCommand(line string, file string) VMCommand {
	line = DiscardComments(line)
	splitLine := strings.Split(line, " ")
	switch splitLine[0] {
	case "push":
		return VMCommand{commandType: "memory", command: "push", memoryRegion: splitLine[1], value: splitLine[2]}
	case "pop":
		return VMCommand{commandType: "memory", command: "pop", memoryRegion: splitLine[1], value: splitLine[2]}
	case "function":
		return VMCommand{commandType: "memory", command: "function", name: splitLine[1], value: splitLine[2], file: file}
	case "return":
		return VMCommand{commandType: "memory", command: "return"}
	case "label":
		return VMCommand{commandType: "memory", command: "label", value: splitLine[1]}
	case "call":
		return VMCommand{commandType: "memory", command: "call", name: splitLine[1], value: splitLine[2], file: file}
	case "goto":
		return VMCommand{commandType: "memory", command: "goto", value: splitLine[1]}
	case "if-goto":
		return VMCommand{commandType: "memory", command: "if-goto", value: splitLine[1]}
	case "add":
		return VMCommand{commandType: "algebra", command: "add"}
	case "sub":
		return VMCommand{commandType: "algebra", command: "sub"}
	case "eq":
		return VMCommand{commandType: "algebra", command: "eq"}
	case "neg":
		return VMCommand{commandType: "algebra", command: "neg"}
	case "gt":
		return VMCommand{commandType: "algebra", command: "gt"}
	case "lt":
		return VMCommand{commandType: "algebra", command: "lt"}
	case "and":
		return VMCommand{commandType: "algebra", command: "and"}
	case "or":
		return VMCommand{commandType: "algebra", command: "or"}
	case "not":
		return VMCommand{commandType: "algebra", command: "not"}
	default:
		return VMCommand{commandType: "Unknown"}
	}
}

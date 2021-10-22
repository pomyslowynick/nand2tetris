// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/04/Mult.asm

// Multiplies R0 and R1 and stores the result in R2.
// (R0, R1, R2 refer to RAM[0], RAM[1], and RAM[2], respectively.)
//
// This program only needs to handle arguments that satisfy
// R0 >= 0, R1 >= 0, and R0*R1 < 32768.

    // Set initial variables
    @R1
    D=M
    @n
    M=D
    @R2
    M=0
    @i
    M=0

    // Loop over i to add R0 times R1 to the sum
    (LOOP)
    @i
    D=M
    @n
    D=D-M
    // if n - i == 0 then END
    @END
    D;JEQ

    @R0
    D=M
    @R2
    M=D+M

    // i++
    @i
    M=M+1

    @LOOP
    0;JMP

    (END)
    @END
    0;JMP

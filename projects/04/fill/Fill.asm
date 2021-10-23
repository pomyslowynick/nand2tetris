// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/04/Fill.asm

// Runs an infinite loop that listens to the keyboard input.
// When a key is pressed (any key), the program blackens the screen,
// i.e. writes "black" in every pixel;
// the screen should remain fully black as long as the key is pressed. 
// When no key is pressed, the program clears the screen, i.e. writes
// "white" in every pixel;
// the screen should remain fully clear as long as no key is pressed.

// Put your code here.

// pseudo code
// if KBR > 0 then while(i < SCREEN(MAX)) { RAM[addr] = -1 }

// Set the initial variables
@SCREEN
D=A
@addr
M=D
@i
M=0
@32
D=A

// Busy loop, wait until key is pressed
(WAIT)
@KBD
D=M
@BLACKEN
D;JGT
@WHITEN
D;JEQ
@WAIT
0;JMP

// Turn the screen black on key press
(BLACKEN)
@addr
D=M
@KBD
D=A-D
@WAIT
D;JEQ   // If we are on KBD address (one right after last SCREEN) jump back to busy loop

@addr
A=M
M=-1
@addr
M=M+1
@WAIT
0;JMP

// Turn screen white when no key is pressed
(WHITEN)
@addr
D=M
D=D+1   // Add one to cover last screen address
@SCREEN
D=A-D
@WAIT
D;JEQ   // If we are on SCREEN address jump back to busy loop

@addr
A=M
M=0
@addr
M=M-1
@WAIT
0;JMP
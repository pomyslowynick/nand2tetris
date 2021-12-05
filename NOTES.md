## 15/10/2021
Started on the course, working on project 1 at the moment.
I have in the past completed parts of the course, to my astonishment I still remember parts of this assignment. It seems so much easier after few years and exposure to professional dev work.

## 19/10/2021
Getting good hang of the material, proceeding actually quite quickly through the course. It's funny how I struggle with exactly same problems now after 2-3 years, like the usage of DFF in Bit implementation. I remember what I struggled with, but not the solutions, or only vaguelly. I am approaching week 4, which is the one at which I stopped the last time. Can't wait to learn more!

BTW: I am quite impressed that I made it so far few years ago with absurdly less knowledge about computing systems.

I loved the notion of the universal machine introduced by Turing, haven't really thought about it before. It makes sense that computers in prnciple have the capability to calculate anything that a Turing machine can. 

## 22/10/2021
In video 4.7 there is a great quote from Donald Knuth that I haven't seen before. It boils down to "code as if you were explaining what you are doing to a fellow programmer".
In Hack machine language there is only one type of variable, 16 bits of data.

Pointers seem quite easy from the hardware perspective, or maybe it's because I have already learnt about them in C?

## 23/10/2021
Von Neumann architecture.
Universal Turing Machine.

## 27/10/2021
It took a good while to complete project 5, but it felt extremely satisfying to connect the Memory, CPU and ROM together and have the tests pass.

Time to tackle the last week: Assembler.

Assembler is a program which translates assembly to machine language.
It's the first software layer above the hardware.
Two passes when translating code to machine language: basically build a table of all the labels before translating the rest into instructions. Other approach is to keep the variables in the table and if the are not defined just update their value when they are.

## 30/10/2021

Parts of assembler I need to develop:
* Reading and parsing commands - Parser
    * Read a file with a given name
    * Get the next command in the file
    * Get the fields of the current command
* Converting mnemonics to code
* Handling symbols

## 11/1/2021
Finally finished part I :D Took me about 17 days to finish it all, I think that's quite good.
It was a great course, even though I have done a lot of it before I enjoyed going over some of this stuff again.
Assembler took me a few hours with all those different cases, I don't like the code I've written all that much. 
Could refactor it one day, but it's not really likely.

Watched first video of Nand 2 Tetris Part 2, so this is the official beginning! Let's how long it will take to finish it.

# 11/2/2021
Still in week 1, taking to heart what Shimon has said about the course being much harder, spamegg said so too.
And he had PhD in maths when taking it, so I guess it really must be.

* Two tier compilation - Java employs it to "compile code once and run it anywhere"
* Java is translated into bytecode (more general VM code)
* JVM - Java Virtual Machine - program that takes bytecode and translates it into target code.
* C++ compiler - compiles from very top to bottom, doesn't use a Virtual Machine.
* Virtualization - one of the most important ideas in whole of Computer Science
* "Program that computes another program" - the idea of Turings UM has contributed to bring computing onto a much higher level of sophistication.
* "Thinking about thinking is the hallmark of intelligence" 
* Stack machine abstraction - Architecture and set of commands
* Stack arithmetic, pop the arguments from the stack, compute f, push it back on the stack.
* Segmentation of memory allows us to maintain the descriptors of variables like static or local (scope?).
Pointer manipulation as a way to divide memory into segments

VM Parser components:
* Parser
* CodeWriter
* Main

# 12/5/2021

Back on it after almost a month, spent almost all of my time on Programming Languages course. I have forgotten few things about how the stack machine is organized, but I am going to revise the material and get the first homework done :) 


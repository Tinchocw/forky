# Forky Interpreter

A concurrent interpreter for the Forky programming language, designed to provide simple fork-join parallelism primitives. The language is designed for ease of use with unique string literal syntax and built-in support for parallel execution.

## Overview

Forky is a programming language that combines sequential and parallel execution paradigms. It features:

- **Parallel Scanning**: Efficient tokenization that can run in parallel
- **Fork Primitives**: Simple constructs for concurrent execution
- **Array Support**: Multi-dimensional arrays with intuitive syntax
- **Functional Programming**: First-class functions with lexical scoping
- **Unique Syntax**: String literals delimited by `"` and `'`, negative operator `~`

## Installation

### Prerequisites

- Go 1.19 or later

### Building

```bash
git clone https://github.com/Tinchocw/forky.git
cd forky
go build -o forky main.go
```

### Usage

```bash
./forky [<filename>.forky] [options]
```

#### Command Line Options

Forky supports several command-line flags to control its behavior:

- `-debug`: Enable debug output for troubleshooting
- `-mode <mode>`: Set the run mode
  - `normal`: Full execution (default)
  - `scanning`: Only perform lexical analysis
  - `parsing`: Only perform parsing (no execution)
- `-workers <number>`: Number of workers for parallel scanning (default: 4)

#### Examples

```bash
# Run a Forky program
./forky examples/fundamentals/basic.forky

# Run with debug output
./forky -debug examples/fundamentals/basic.forky

# Only scan the file (no execution)
./forky -mode scanning examples/fundamentals/basic.forky

# Run parsing only
./forky -mode parsing examples/fundamentals/basic.forky

# Use 8 workers for parallel scanning
./forky -workers 8 examples/fundamentals/basic.forky

# Start REPL mode (interactive)
./forky
```

#### REPL Mode

When run without arguments, Forky starts an interactive Read-Eval-Print Loop (REPL) for experimenting with code:

```
Forky - REPL with arrow key support. Ctrl-C or Ctrl-D (on empty line) to exit.
Use ↑↓ arrows for history, ←→ for line editing.

> var x = 42;
> print(x);
42
> func add(a, b) { return a + b; }
> print(add(5, 3));
8
>
```

**REPL Features:**
- **Interactive execution**: Type Forky statements and see results immediately
- **Arrow key support**: Use ↑↓ arrows to navigate command history, ←→ for line editing
- **Multi-line support**: Enter complex statements across multiple lines
- **Exit commands**: Press Ctrl-C to abort current input, or Ctrl-D (on empty line) to exit

The REPL maintains a history of your commands and supports all Forky language features including functions, variables, and parallel execution.

## Language Syntax

### Comments

Forky does not support comments. The `//` comments shown in the examples below are for documentation purposes only and are not part of the actual language syntax.

### Data Types

- **Numbers**: Integer literals (e.g., `42`)
- **Strings**: Delimited by `"` at start and `'` at end (e.g., `"hello'`)
- **Booleans**: `true`, `false`
- **None**: `none` (null value)
- **Arrays**: Multi-dimensional arrays

### Variables

#### Declaration

```forky
var x = 5;
var name = "hello';
var empty;  // initialized to none
```

#### Assignment

```forky
set x = 10;
set name = "world';
```

### Arrays

#### Declaration

```forky
var arr[5];           // 1D array of length 5, filled with none
var matrix[3][3];     // 3x3 matrix, filled with none
var cube[2][2][2] = 1; // 3D array filled with 1
```

#### Accessing

```forky
var value = arr[0];
var element = matrix[1][2];
```

#### Assignment

```forky
set arr[0] = 42;
set matrix[1][2] = "value';
```

#### Array Literals

```forky
var arr = [1, 2, 3];
var matrix = [[1, 2], [3, 4]];
```

### Operators

#### Arithmetic

- Addition: `+`
- Subtraction: `-`
- Multiplication: `*`
- Division: `/`
- Negation: `~` (unary)

#### Comparison

- Equal: `==`
- Not equal: `!=`
- Less than: `<`
- Less than or equal: `<=`
- Greater than: `>`
- Greater than or equal: `>=`

#### Logical

- And: `and`
- Or: `or`
- Not: `!`

#### Concatenation

The `+` operator concatenates values of different types:

```forky
print(5 + " items');  // Outputs: 5 items
```

### Control Flow

#### If-Else

```forky
if (condition) {
    // code
} else if (other_condition) {
    // code
} else {
    // code
}
```

#### While Loops

```forky
while (condition) {
    // code
}
```

#### Break

```forky
break;
```

### Functions

#### Definition

```forky
func add(a, b) {
    return a + b;
}
```

#### Call

```forky
var result = add(5, 3);
```

### Parallel Execution

#### Fork Block

Executes multiple blocks of code in parallel:

```forky
fork {
    {
        print("Task 1');
        // more statements
    }
    {
        print("Task 2');
        // more statements
    }
    {
        print("Task 3');
        // more statements
    }
}
```

Each inner block runs in its own environment concurrently.

#### Fork Array

Iterates over array elements in parallel:

```forky
var numbers = [1, 2, 3, 4, 5];

fork numbers elem {
    print(elem);
}
```

Or with index:

```forky
fork numbers index, elem {
    print(index + ": " + elem);
}
```

If only one identifier is provided, it defaults to the element.

### Print Statement

```forky
print(expression);
```

### Truthiness

- Arrays are truthy if they are not empty
- Numbers: non-zero is true
- Strings: non-empty is true
- `none` is false
- `false` is false

## Error Handling

Forky includes runtime error checking for common programming mistakes:

#### Division by Zero

```forky
var result = 10 / 0;  // Runtime error: Division by zero
```

#### Array Access Out of Bounds

```forky
var arr[3];
print(arr[5]);  // Runtime error: Array index out of bounds
```

#### Accessing Non-Array as Array

```forky
var x = 5;
print(x[0]);  // Runtime error: Cannot access index on non-array value
```

#### Calling Non-Function

```forky
var x = 42;
x();  // Runtime error: Cannot call non-function value
```

#### Incompatible Types in Operations

```forky
var a = "hello';
var b = 5;
print(a - b);  // Runtime error: Incompatible types for subtraction
```

#### Undefined Variable Access

```forky
print(undefined_var);  // Runtime error: Undefined variable
```

#### Function Call with Wrong Number of Arguments

```forky
func add(a, b) {
    return a + b;
}
add(1);  // Runtime error: Wrong number of arguments
```

When a runtime error occurs, the interpreter will display an error message and halt execution.

## Examples

### Basic Arithmetic

```forky
var a = 10;
var b = 5;
print(a + b);  // 15
print(a - b);  // 5
print(a * b);  // 50
print(a / b);  // 2
```

### String Manipulation

```forky
var greeting = "Hello';
var name = "World';
print(greeting + " " + name + "!');  // Hello World!
```

### Arrays

```forky
var arr[3] = 0;
set arr[0] = 1;
set arr[1] = 2;
set arr[2] = 3;

print(arr[0]);  // 1
print(arr[1]);  // 2
print(arr[2]);  // 3
```

### Functions

```forky
func greet(name) {
    print("Hola ' + name + "!');
}

func add(x, y) {
    return x + y;
}

func factorial(n) {
    if (n <= 1) {
        return 1;
    } else {
        return n * factorial(n - 1);
    }
}

greet("Mundo');
var result = add(5, 3);
print(result);

var fact = factorial(5);
print(fact);

func power(base, exp) {
    if (exp == 0) {
        return 1;
    } else {
        return base * power(base, exp - 1);
    }
}
```

### Control Flow

```forky
var counter = 0;
while (counter < 5) {
    print("Contador: ' + counter);
    set counter = counter + 1;
}

var i = 0;
while (i < 3) {
    var j = 0;
    while (j < 3) {
        print("i: ' + i + ", j: ' + j);
        set j = j + 1;
    }
    set i = i + 1;
}

var n = 0;
while (true) {
    if (n >= 5) {
        break;
    }
    print("n: ' + n);
    set n = n + 1;
}

print("Countdown: ');

func countdown(x) {
    while (true) {
        print(x);
        set x = x - 1;

        if (x == 0) {
            return "Despegue';
        }
    }
}

print(countdown(5));
```

### Parallel Execution

#### Fork Block

```forky
func worker(id) {
    print("Worker " + id + " starting');
    // simulate work
    var i = 0;
    while (i < 1000000) {
        set i = i + 1;
    }
    print("Worker " + id + " done');
}

fork {
    worker(1);
    worker(2);
    worker(3);
}
```

#### Fork Array

```forky
var data = [10, 20, 30, 40, 50];

func process(value) {
    return value * 2;
}

fork data elem {
    print(process(elem));
}
```

#### Nested Parallel Processing

```forky
var size = 3;
var arr[size][size] = 0;

fork arr i, e {
    fork e j, f {
        set arr[i][j] = i + j + f;
        print("Row: ' + i + ", Col: ' + j + ", Value: ' + arr[i][j]);  
    }
}
print(arr);
```

### Boolean Operations

```forky
var a = 10;
var b = 20;
var c = 15;

print(a < b);
print(b > c);
print(a == 10);
print(b != c);

var result1 = a < b and c > a;
print(result1);

var result2 = a > b or c < b;
print(result2);

var complex = (a < b) and (c > 5) and !(b == 15);
print(complex);
```

### Variable Scoping

```forky
var globalVar = "soy global';

func testScope() {
    var localVar = "soy local';
    set globalVar = "modificado desde función';
    
    print(localVar);
    print(globalVar);
}

print("Antes de la función:');
print(globalVar);

testScope();

print("Después de la función:');
print(globalVar);

if (true) {
    var blockVar = "variable de bloque';
    var globalVar = "variable global en bloque';
    print("Dentro del bloque:');
    print(blockVar);
    print(globalVar);
}

print("Fuera del bloque:');
print(globalVar);
```

## Architecture

The interpreter consists of several components:

- **Scanner**: Parallel tokenization of source code
- **Parser**: Recursive descent parser building AST
- **Interpreter**: Tree-walking interpreter with concurrent execution support
- **Resolver**: Static analysis for variable resolution

## License

This project is licensed under the MIT License.

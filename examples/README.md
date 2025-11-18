# Forky Examples

This folder contains examples demonstrating the features of the Forky programming language. The examples are organized into two main categories:

- **fundamentals/**: Basic language features and syntax
- **usecases/**: Practical applications and data structures

## Fundamentals

The fundamentals folder contains examples that introduce core language concepts in a logical progression.

### 1. `basic.forky`
- Basic variable declaration with `var`
- Printing values with `print()`
- String literals

### 2. `math.forky`
- Arithmetic operations: `+`, `-`, `*`, `/`
- Operator precedence
- Variable declarations and assignments

### 3. `strings.forky`
- String literals and concatenation with `+`
- Combining strings and numbers

### 4. `arrays.forky`
- Array declaration with `var name[size]`
- Array access with `arr[index]`
- Array assignment with `set arr[index] = value`
- Multi-dimensional arrays

### 5. `booleans.forky`
- Comparison operators: `<`, `>`, `==`, `!=`, `<=`, `>=`
- Logical operators: `and`, `or`
- Negation with `!`

### 6. `variables.forky`
- Variable scoping (global, local, block)
- Variable shadowing
- Variable reassignment with `set`

### 7. `functions.forky`
- Function definition with `func name(params)`
- Function calls
- Return statements
- Recursion

### 8. `first_class_functions.forky`
- First-class functions
- Passing functions as arguments
- Function expressions

### 9. `conditionals.forky`
- `if` statements
- `else if` and `else` clauses
- Nested conditionals

### 10. `loops.forky`
- `while` loops
- `break` statements
- Nested loops
- Loop control

### 11. `fork.forky`
- Parallel execution with `fork { { block1 } { block2 } ... }`
- Concurrent block execution

### 12. `forkArray.forky`
- Parallel array iteration with `fork arr var { ... }`
- Nested fork statements
- Array processing in parallel

### 13. `complex.forky`
- Complex program combining multiple features
- Integration of functions, loops, and conditionals

### 14. `errors.forky`
- Common error cases and runtime errors
- Examples of what causes errors

## Use Cases

The usecases folder contains practical examples showing how to implement common data structures and algorithms.

### 1. `dynamic_vector.forky`
- Implementation of a dynamic vector (array list)
- Methods: append, get, pop, size, resize
- Error handling for bounds checking
- Demonstration of encapsulation with functions

### 2. `multi_dim_sums.forky`
- Parallel computation of sums in multi-dimensional arrays
- Use of nested fork statements for parallel processing
- Array manipulation and aggregation

### 3. `parallel_fork_usecase.forky`
- Practical example of parallel processing with fork
- Concurrent execution of independent tasks
- Performance demonstration with parallel blocks

## Syntax Notes

- Variable declarations: `var name = value;`
- Assignments: `set name = value;`
- Array assignments: `set arr[index] = value;`
- Functions: `func name(params) { statements }`
- Conditionals: `if (condition) { statements } else if (condition) { statements } else { statements }`
- Loops: `while (condition) { statements }`
- Parallel execution: `fork { { block1 } { block2 } }`
- Strings: `"text'` (starts with `"` and ends with `'`)
- Negation: `~value`
- All statements end with `;`

## Running Examples

```bash
# Build the interpreter
make build

# Run a fundamental example
./forky examples/fundamentals/basic.forky

# Run a use case example
./forky examples/usecases/dynamic_vector.forky

# Run in parsing mode to see AST
./forky --mode parsing examples/fundamentals/basic.forky

# Run in scanning mode to see tokens
./forky --mode scanning examples/fundamentals/basic.forky

# Inject an example and continue in REPL for interactive experimentation
make inject FILE=examples/usecases/dynamic_vector.forky
```

The inject mode loads the example code and then starts the REPL, allowing you to interact with the loaded functions and variables. This is particularly useful for the usecases examples like the dynamic vector, where you can test the implemented functions interactively.
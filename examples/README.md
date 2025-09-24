# Forky Examples

Esta carpeta contiene ejemplos simples del lenguaje Forky para demostrar sus características.

## Ejemplos disponibles:

### 1. `basic.forky`
- Asignación de variables básica
- Operaciones aritméticas
- Uso de `print()`

### 2. `conditionals.forky`
- Declaraciones `if/else`
- Condicionales anidados
- Comparaciones

### 3. `boolean.forky`
- Operaciones booleanas (`and`, `or`, `not`)
- Comparaciones (`<`, `>`, `==`, `!=`)
- Expresiones complejas

### 4. `functions.forky`
- Definición de funciones con `func`
- Parámetros y `return`
- Recursión (factorial)

### 5. `loops.forky`
- Bucles `while`
- Bucles anidados
- Contadores

### 6. `strings.forky`
- Literales de cadena con formato `"texto'`
- Manejo de strings
- Caracteres especiales

### 7. `math.forky`
- Operaciones matemáticas
- Precedencia de operadores
- Uso de paréntesis

### 8. `variables.forky`
- Declaración de variables con `var`
- Alcance (scope) de variables
- Variables locales y globales

### 9. `complex.forky`
- Ejemplo complejo combinando múltiples características
- Fibonacci recursivo
- Verificación de números primos

## Cómo usar los ejemplos:

```bash
# Modo parsing (para ver el AST)
./forky --mode parsing examples/basic.forky

# Modo scanning (para ver los tokens)
./forky --mode scanning examples/conditionals.forky

# O simplemente usando make
make build
./forky examples/math.forky
```

## Formato de strings

Nota que Forky usa un formato especial para strings: `"texto'` (comienza con `"` y termina con `'`).
# Distribución

Como vamos a hacer nuestro interprete concurrente para mejorar la performance.

## Etapas

### Escaneo

Tranformar un archivo de texto o un stream de bytes en una secuencia de tokens.

ejemplo:

concat ("hola', )

index = 14
token_index = 6
parenthesis_index = -2

tokens = [
    {
        SIMBOLO_MAS: [1],
        LITERAL: [0],
        LITERAL_START_INCOMPLETE
    },
    {
        SIMBOLO_ASTERISCO: [3],
        LITERAL: [2],
    },
    {
        SIMBOLO_MAS: [5],
        LITERAL: [4, 6],
    }
]

THREAD 0

tokens = {
    0: {},
    1 : {}
    2 : {}
}

index = 0 -> 7
parenthesis_index = 2

----------------------------------

THREAD 1

tokens = {
    -2: {}
    -1 : {}
    0 : {}
}

index = 0 -> 5
parenthesis_index = -2

### Parsing

Tranformar una secuencia de tokens en un árbol de sintáxis abstracta (AST).

"""md
expression     → equality ;
equality       → comparison ( ( "!=" | "==" ) comparison )* ;
comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
term           → factor ( ( "-" | "+" ) factor )* ;
factor         → unary ( ( "/" | "*" ) unary )* ;
unary          → ( "!" | "-" ) unary | primary ;
primary        → NUMBER | STRING | "true" | "false" | "nil"
                 | "(" expression ")" ;
"""

TOKENs = [LITERAL[1], SIMBOLO_ASTERISCO, LITERAL[3], SIMBOLO_MAS, LITERAL[4], SIMBOLO_ASTERISCO, LITERAL[2]]

swithc TOKEN {
    case LITERAL:
        AST.push(NODO_LITERAL)
    case SIMBOLO_MAS:
        AST.push(NODO_SUMA)
    case ...
}

TOKENs = {
    SIMBOLO_MAS: [3],
    SIMBOLO_ASTERISCO: [1, 5],
    LITERAL: [0, 2, 4, 6]
    LITERAL_VALOR: [1, 3, 4, 2]
}

"""plaintext
3 * (3 + (1 + 2) * 3 == 2) + (2 + 5)
"""

"""plaintext
3 * (1 + 1)
"""

TOKENS = {
    PARENTESIS_ABRE: [0],
    PARENTESIS_CIERRA: [4],
    SIMBOLO_IGUAL_IGUAL: [7]
    SIMBOLO_MAS: [2],
    SIMBOLO_ASTERISCO: [5],
    LITERAL: [1, 3, 6, 8],
    LITERAL_VALOR: [1, 2, 3, 2]
}

----------------

DIC: Map<Token, []Posicion>

TOKENs: []Dic

-----------------

ejemplo: 1 + (2 * (3 + 1)) + (3 + 1)

THREAD: 0 (0-10)

TOKENS = [
    {
        SIMBOLO_MAS: [1, 7],
        LITERAL: [0]
    },
    {
        SIMBOLO_MAS: [9],
        SIMBOLO_ASTERISCO: [3],
        LITERAL: [2, 8, 10]
    },
    {
        SIMBOLO_MAS: [5],
        LITERAL: [4, 6]
    }
]

THREAD: 0.0 (0-0)

TOKENS = [
    {
        SIMBOLO_MAS: [1, 7],
        LITERAL: [0]
    },
    {
        SIMBOLO_MAS: [9],
        SIMBOLO_ASTERISCO: [3],
        LITERAL: [2, 8, 10]
    },
    {
        SIMBOLO_MAS: [5],
        LITERAL: [4, 6]
    }
]

THREAD: 0.1 (2-10)

TOKENS = [
    {
        SIMBOLO_MAS: [1, 7],
        LITERAL: [0]
    },
    {
        SIMBOLO_MAS: [9],
        SIMBOLO_ASTERISCO: [3],
        LITERAL: [2, 8, 10]
    },
    {
        SIMBOLO_MAS: [5],
        LITERAL: [4, 6]
    }
]

THREAD: 0.1.0 (2-6)

TOKENS = [
    {
        SIMBOLO_MAS: [1, 7],
        LITERAL: [0]
    },
    {
        SIMBOLO_MAS: [9],
        SIMBOLO_ASTERISCO: [3],
        LITERAL: [2, 8, 10]
    },
    {
        SIMBOLO_MAS: [5],
        LITERAL: [4, 6]
    }
]

THREAD: 0.1.1 (8-10)

TOKENS = [
    {
        SIMBOLO_MAS: [1, 7],
        LITERAL: [0]
    },
    {
        SIMBOLO_MAS: [9],
        SIMBOLO_ASTERISCO: [3],
        LITERAL: [2, 8, 10]
    },
    {
        SIMBOLO_MAS: [5],
        LITERAL: [4, 6]
    }
]

ejemplo:  3 + !false + 3 * 2

0 - 8
0 - 3 y 5 - 8
0 - 0 y 2 - 3

### Ejecución

TBD

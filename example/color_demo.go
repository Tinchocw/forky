package example

import (
	"fmt"

	"github.com/Tinchocw/Interprete-concurrente/common"
)

func main() {
	fmt.Println("🎨 Demostración de funciones de colores 🎨\n")

	// Colores básicos
	fmt.Println("=== Colores básicos ===")
	fmt.Println(common.Red("Este texto es rojo"))
	fmt.Println(common.Green("Este texto es verde"))
	fmt.Println(common.Blue("Este texto es azul"))
	fmt.Println(common.Yellow("Este texto es amarillo"))
	fmt.Println(common.Magenta("Este texto es magenta"))
	fmt.Println(common.Cyan("Este texto es cian"))
	fmt.Println()

	// Estilos
	fmt.Println("=== Estilos ===")
	fmt.Println(common.Bold("Este texto es negrita"))
	fmt.Println(common.Italic("Este texto es cursiva"))
	fmt.Println(common.Underline("Este texto está subrayado"))
	fmt.Println(common.Strike("Este texto está tachado"))
	fmt.Println()

	// Combinaciones
	fmt.Println("=== Combinaciones ===")
	fmt.Println(common.BoldRed("Rojo en negrita"))
	fmt.Println(common.BoldGreen("Verde en negrita"))
	fmt.Println(common.UnderlineBlue("Azul subrayado"))
	fmt.Println()

	// Funciones divertidas
	fmt.Println("=== Funciones especiales ===")
	fmt.Println(common.Rainbow("¡Texto arco iris!"))
	fmt.Println(common.Gradient("Gradiente", common.COLOR_RED, common.COLOR_BLUE))
	fmt.Println()

	// Box
	fmt.Println("=== Cajas ===")
	fmt.Println(common.Box("¡Hola mundo!", common.COLOR_GREEN, common.BG_BLACK))
	fmt.Println()

	// Printf con colores
	fmt.Println("=== Printf con colores ===")
	common.PrintfRed("Error: %s\n", "Algo salió mal")
	common.PrintfGreen("Éxito: %s\n", "Todo funcionó bien")
	common.PrintfYellow("Advertencia: %d elementos pendientes\n", 5)
	common.PrintfBlue("Info: Versión %s\n", "1.0.0")
	fmt.Println()

	// Uso manual con Colorize
	fmt.Println("=== Colorize personalizado ===")
	fancy := common.Colorize("¡Súper fancy!", common.STYLE_BOLD, common.STYLE_UNDERLINE, common.COLOR_MAGENTA)
	fmt.Println(fancy)
}

package main
import	(
	`fmt`
	`github.com/nsf/termbox-go`
)

func main() {
	fore := termbox.ColorGreen | termbox.AttrUnderline
	fmt.Printf("f: %08b\n", fore)
	fore = termbox.ColorGreen | termbox.AttrUnderline | termbox.AttrReverse
	fmt.Printf("f: %08b\n", fore)
	fore &= ^termbox.AttrReverse
	fmt.Printf("f: %08b\n", fore)
}

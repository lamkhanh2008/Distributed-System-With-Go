package main

type Computer interface {
	print()
	setPrinter(printer)
}

type printer interface {
	printFile()
}

type Windows struct {
}

func (w *Windows) print() {

}

func (w *Windows) setPrinter(p printer) {

}

type HPprinter struct {
}

func (h *HPprinter) printFile() {

}
func main() {
	w := Windows{}
	HPprinter := HPprinter{}
	w.setPrinter(&HPprinter)
}

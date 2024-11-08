package main

import "fmt"

type Inode interface {
	print(string)
	clone() Inode
}

type File struct {
	name string
}

func (f *File) print(indentation string) {
	fmt.Println(indentation + f.name)
}

func (f *File) clone() *File {
	return &File{name: f.name + "_clone"}
}

type Folder struct {
	children []File
	name     string
}

func (f *Folder) print(indentation string) {
	fmt.Println(indentation + f.name)
	for _, i := range f.children {
		i.print(indentation + indentation)
	}
}

func (f *Folder) clone() *Folder {
	cloneFolder := &Folder{name: f.name + "_clone"}
	var tempChildren []File
	for _, i := range f.children {
		copy := i.clone()
		tempChildren = append(tempChildren, *copy)
	}
	cloneFolder.children = tempChildren
	return cloneFolder
}

type Person struct {
	name string
}

func TestPointer(p *Person) *Person {
	return &Person{
		name: p.name,
	}
}

func main() {
	file1 := &File{name: "File1"}
	file2 := &File{name: "File2"}
	file3 := &File{name: "File3"}

	folder1 := &Folder{
		children: []File{*file1},
		name:     "Folder1",
	}
	fmt.Println(folder1)
	folder2 := &Folder{
		children: []File{*file2, *file3},
		name:     "Folder2",
	}
	fmt.Println("\nPrinting hierarchy for Folder2")
	folder2.print("  ")

	cloneFolder := folder2.clone()
	fmt.Println("\nPrinting hierarchy for clone Folder")
	cloneFolder.print("  ")

	fmt.Printf("Original file address: %p\n", file2)
	fmt.Printf("Cloned file address: %p\n", &cloneFolder.children[0])
	// p := Person{name: "sss"}
	// fmt.Printf("--- %p", &p)
	// new := TestPointer(&p)
	// fmt.Printf("--- %p", new)

}

// 1 folder có nhiều file hoặc folder khác. Hãy đếm số lượng keyword có trong folder đó, key word có thể nằm trong tệp file.

package main

type component interface {
	search(key string)
}

type Folder struct {
	components []component
}

func (f *Folder) search(key string) {
	for _, component := range f.components {
		component.search(key)
	}
}

func (f *Folder) add(comp component) {
	f.components = append(f.components, comp)
}

type File struct {
}

func (f *File) search(key string) {

}

func main() {
	folder1 := Folder{}
	file1 := File{}
	folder2 := Folder{}
	folder1.add(&file1)
	folder2.add(&folder1)

	folder2.search("hehe")
}

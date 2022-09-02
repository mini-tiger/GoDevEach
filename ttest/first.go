package main

func main() {
	var m map[string]interface{} = make(map[string]interface{}, 0)
	m["1"] = 1
	a(m)
}
func a(a interface{}) {
	aa := make(map[string]interface{}, 2)
	aa["a"] = 1
	print(len(aa))
}

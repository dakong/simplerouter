package simplerouter

import (
	"fmt"
	"net/http"
	"testing"
)

func printTree(n *node) {
	fmt.Print(n.component + " { ")
	for _, child := range n.children {
		printTree(child)
	}
	fmt.Print(" }")
}

func assert(t *testing.T, value string, actual string) {
	if value != actual {
		t.Error("Expected", value, ", got", actual)
	}
}

// TestAddNodes tests creations of nodes
func TestAddNodes(t *testing.T) {
	emptyHandler := func(res http.ResponseWriter, req *http.Request) {}

	root := node{component: "/", isNamedParam: false, methods: make(map[string]Handle)}

	root.addNode("GET", "/blog", emptyHandler)
	assert(t, "/{blog{}}", root.String())

	root.addNode("GET", "/human", emptyHandler)
	assert(t, "/{blog{}human{}}", root.String())

	root.addNode("POST", "/human/student", emptyHandler)
	assert(t, "/{blog{}human{student{}}}", root.String())

	root.addNode("PUT", "/human/student/biology", emptyHandler)
	assert(t, "/{blog{}human{student{biology{}}}}", root.String())

	root.addNode("DELETE", "/human/professor", emptyHandler)
	assert(t, "/{blog{}human{student{biology{}}professor{}}}", root.String())

	root.addNode("GET", "/human/professor/engineering", emptyHandler)
	assert(t, "/{blog{}human{student{biology{}}professor{engineering{}}}}", root.String())

	root.addNode("GET", "/blog/techcrunch", emptyHandler)
	assert(t, "/{blog{techcrunch{}}human{student{biology{}}professor{engineering{}}}}", root.String())

	root.addNode("GET", "/blog/human/biology", emptyHandler)
	assert(t, "/{blog{techcrunch{}human{biology{}}}human{student{biology{}}professor{engineering{}}}}", root.String())

	root.addNode("GET", "/animal/dog/retriever", emptyHandler)
	assert(t, "/{blog{techcrunch{}human{biology{}}}human{student{biology{}}professor{engineering{}}}animal{dog{retriever{}}}}", root.String())

	rootTwo := node{component: "/", isNamedParam: false, methods: make(map[string]Handle)}

	rootTwo.addNode("GET", "/animal/mammal/llama", emptyHandler)
	assert(t, "/{animal{mammal{llama{}}}}", rootTwo.String())

	rootTwo.addNode("GET", "/animal/mammal", emptyHandler)
	assert(t, "/{animal{mammal{llama{}}}}", rootTwo.String())

}

func TestRegisterHandler(t *testing.T) {
	emptyHandler := func(res http.ResponseWriter, req *http.Request) {}
	root := node{component: "/", isNamedParam: false, methods: make(map[string]Handle)}
	root.addNode("GET", "/animal/dog/retriever", emptyHandler)

	components := []string{"animal", "dog", "retriever"}
	leaf, _ := root.traverseTree(components)

	if leaf.methods["GET"] == nil {
		t.Error("Expected a Hanlder, got", nil)
	}

	components = []string{"animal", "dog"}
	leaf, _ = root.traverseTree(components)

	if leaf.methods["GET"] != nil {
		t.Error("Expected a nil, got a handler")
	}

	root.addNode("POST", "/animal", emptyHandler)

	components = []string{"animal"}
	leaf, _ = root.traverseTree(components)

	if leaf.methods["POST"] == nil {
		t.Error("Expected a Hanlder, got", nil)
	}
}

func TestPathParam(t *testing.T) {
	emptyHandler := func(res http.ResponseWriter, req *http.Request) {}
	root := node{component: "/", isNamedParam: false, methods: make(map[string]Handle)}

	root.addNode("GET", "/animal/dog/:id", emptyHandler)
	root.addNode("GET", "/animal/dog/:breed/:id", emptyHandler)
	assert(t, "/{animal{dog{id{}breed{id{}}}}}", root.String())

	_, params := root.search("GET", "/animal/dog/5")
	val, found := params.GetValue("id")

	if !found || val != "5" {
		t.Error("Expected id to be 5, got", val)
	}

	_, params = root.search("GET", "/animal/dog/retriever/2")
	breed, breedFound := params.GetValue("breed")
	id, idFound := params.GetValue("id")

	if !idFound || id != "2" {
		t.Error("Expected id to be 2, got", id)
	}

	if !breedFound || breed != "retriever" {
		t.Error("Expected breed to be retriever, got", breed)
	}
}

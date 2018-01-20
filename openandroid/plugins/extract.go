package main

import(
	"plugin"
	"fmt"
	"log"
)

func main() {
	p, err := plugin.Open("plugin.so")
	if err != nil {
		panic(err)
	}

	k, err := p.Lookup("GetKey")
	if err != nil {
		panic(err)
	}
	keyfunc, ok := k.(func() string)
	if !ok {
		log.Fatal("malformed key function")
	}
	key := keyfunc()

	v, err := p.Lookup("GetValue")
	if err != nil {
		panic(err)
	}
	valuefunc, ok := v.(func() interface{})
	if !ok {
		log.Fatal("malformed value function")
	}
	result := valuefunc()
	value, ok := result.([]string)
	if !ok {
		log.Fatal("malformed value type")
	}

	fmt.Println("key: " + key)
	fmt.Printf("value: %v, %v\n", value[0], value[1])
}
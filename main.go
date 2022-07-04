package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

type Arguments map[string]string

// type Person map[string]interface{}

// func (p Person) String() string {
// 	return fmt.Sprintf("%v %v %v", p["id"], p["email"], p["age"])
// }

// type People []Person

// func (p People) String() string {
// 	return fmt.Sprint("{" + p.String() + "}")
// }

type Item struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Age   uint   `json:"age"`
}

func (i Item) String() string {
	return fmt.Sprintf(`{"id":"%v","email":"%v","age":%v}`, i.Id, i.Email, i.Age)
}

type Book []Item

func (b Book) String() string {
	s := "["
	for i, val := range b {
		s += fmt.Sprint(val)
		if i < len(b)-1 {
			s += ","
		}
	}
	s += "]"
	return s
}

var (
	op   *bool
	item *bool
	file *bool
)

func parseArgs() Arguments {
	op = flag.Bool("operation", false, "choose operation: add, list, findById, remove")
	item = flag.Bool("item", false, "item to work with")
	file = flag.Bool("fileName", false, "file to work with")
	flag.Parse()
	fmt.Println(*op, *item, *file)
	fmt.Println(flag.Args())
	args := make(Arguments)
	if *op {
		args["operation"] = ""
	}
	if *file {
		args["fileName"] = ""
	}
	return args
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func Perform(args Arguments, writer io.Writer) error {
	if args == nil {
		writer.Write([]byte("-operation flag has to be specified"))
		return errors.New("-operation flag has to be specified")
	}
	if args["operation"] == "" {
		writer.Write([]byte("-operation flag has to be specified"))
		return errors.New("-operation flag has to be specified")
	}
	if args["fileName"] == "" {
		writer.Write([]byte("-fileName flag has to be specified"))
		return errors.New("-fileName flag has to be specified")
	}
	switch args["operation"] {
	case "add":
		if args["item"] == "" {
			writer.Write([]byte("-item flag has to be specified"))
			return errors.New("-item flag has to be specified")
		}
		var itemToAdd Item
		err := json.Unmarshal([]byte(args["item"]), &itemToAdd)
		check(err)
		dat, err := os.ReadFile(args["fileName"])
		var book Book
		isIdFound := false
		if err == nil {
			fmt.Println("MAIN-START" + string(dat) + "MAIN-END")
			err = json.Unmarshal(dat, &book)
			check(err)
			for _, item := range book {
				fmt.Println("MAIN-START" + fmt.Sprint(item) + "MAIN-END")
				if item.Id == itemToAdd.Id {
					isIdFound = true
					break
				}
			}
		}
		if isIdFound {
			writer.Write([]byte("Item with id " + itemToAdd.Id + " already exists"))
		} else {
			book = append(book, itemToAdd)
			dat, err = json.Marshal(book)
			check(err)
			err = os.WriteFile(args["fileName"], dat, 0666)
			check(err)
		}
	case "list":
		dat, err := os.ReadFile(args["fileName"])
		check(err)
		fmt.Println("MAIN-START" + string(dat) + "MAIN-END")
		var book Book
		err = json.Unmarshal(dat, &book)
		check(err)
		writer.Write([]byte(fmt.Sprint(book)))
	case "findById":
		if args["id"] == "" {
			writer.Write([]byte("-id flag has to be specified"))
			return errors.New("-id flag has to be specified")
		}
		dat, err := os.ReadFile(args["fileName"])
		check(err)
		fmt.Println("MAIN-START" + string(dat) + "MAIN-END")
		var book Book
		err = json.Unmarshal(dat, &book)
		check(err)
		for _, item := range book {
			if item.Id == args["id"] {
				writer.Write([]byte(fmt.Sprint(item)))
				break
			}
		}
	case "remove":
		if args["id"] == "" {
			writer.Write([]byte("-id flag has to be specified"))
			return errors.New("-id flag has to be specified")
		}
		dat, err := os.ReadFile(args["fileName"])
		check(err)
		fmt.Println("MAIN-START" + string(dat) + "MAIN-END")
		var book Book
		err = json.Unmarshal(dat, &book)
		check(err)
		out := make(Book, 0)
		isIdFound := false
		for _, item := range book {
			if item.Id != args["id"] {
				out = append(out, item)
			} else {
				isIdFound = true
			}
		}
		if !isIdFound {
			writer.Write([]byte("Item with id " + args["id"] + " not found"))
		} else {
			dat, err = json.Marshal(out)
			check(err)
			err = os.WriteFile(args["fileName"], dat, 0666)
			check(err)
		}
	default:
		writer.Write([]byte("Operation " + args["operation"] + " not allowed!"))
		return errors.New("Operation " + args["operation"] + " not allowed!")
	}
	return nil
}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}

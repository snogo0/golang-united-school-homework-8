package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

const (
	fOperation = "operation"
	fItem      = "item"
	fId        = "id"
	fFileName  = "fileName"
)

const (
	opAdd      = "add"
	opList     = "list"
	opFindById = "findById"
	opRemove   = "remove"
)

const (
	sSpecified = " flag has to be specified"
)

type Arguments map[string]string

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
	op = flag.Bool(fOperation, false, fmt.Sprintf("choose operation: %s, %s, %s, %s", opAdd, opList, opFindById, opRemove))
	item = flag.Bool(fItem, false, "item to work with")
	file = flag.Bool(fFileName, false, "file to work with")
	flag.Parse()
	fmt.Println(*op, *item, *file)
	fmt.Println(flag.Args())
	args := make(Arguments)
	if *op {
		args[fOperation] = ""
	}
	if *file {
		args[fFileName] = ""
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
		s := "-" + fOperation + sSpecified
		writer.Write([]byte(s))
		return errors.New(s)
	}
	if args[fOperation] == "" {
		s := "-" + fOperation + sSpecified
		writer.Write([]byte(s))
		return errors.New(s)
	}
	if args[fFileName] == "" {
		s := "-" + fFileName + sSpecified
		writer.Write([]byte(s))
		return errors.New(s)
	}
	switch args[fOperation] {
	case "add":
		if args[fItem] == "" {
			s := "-" + fItem + sSpecified
			writer.Write([]byte(s))
			return errors.New(s)
		}
		var itemToAdd Item
		err := json.Unmarshal([]byte(args[fItem]), &itemToAdd)
		check(err)
		dat, err := os.ReadFile(args[fFileName])
		var book Book
		isIdFound := false
		if err == nil {
			err = json.Unmarshal(dat, &book)
			check(err)
			for _, item := range book {
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
			err = os.WriteFile(args[fFileName], dat, 0666)
			check(err)
		}
	case "list":
		dat, err := os.ReadFile(args[fFileName])
		check(err)
		var book Book
		err = json.Unmarshal(dat, &book)
		check(err)
		writer.Write([]byte(fmt.Sprint(book)))
	case "findById":
		if args[fId] == "" {
			s := "-" + fId + sSpecified
			writer.Write([]byte(s))
			return errors.New(s)
		}
		dat, err := os.ReadFile(args[fFileName])
		check(err)
		var book Book
		err = json.Unmarshal(dat, &book)
		check(err)
		for _, item := range book {
			if item.Id == args[fId] {
				writer.Write([]byte(fmt.Sprint(item)))
				break
			}
		}
	case "remove":
		if args[fId] == "" {
			s := "-" + fId + sSpecified
			writer.Write([]byte(s))
			return errors.New(s)
		}
		dat, err := os.ReadFile(args[fFileName])
		check(err)
		var book Book
		err = json.Unmarshal(dat, &book)
		check(err)
		out := make(Book, 0)
		isIdFound := false
		for _, item := range book {
			if item.Id != args[fId] {
				out = append(out, item)
			} else {
				isIdFound = true
			}
		}
		if !isIdFound {
			writer.Write([]byte("Item with id " + args[fId] + " not found"))
		} else {
			dat, err = json.Marshal(out)
			check(err)
			err = os.WriteFile(args[fFileName], dat, 0666)
			check(err)
		}
	default:
		s := "Operation " + args[fOperation] + " not allowed!"
		writer.Write([]byte(s))
		return errors.New(s)
	}
	return nil
}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}

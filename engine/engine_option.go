package engine

import (
	"errors"
	"fmt"
	"strconv"
)

type EngineOption interface {
	ToUci() string
	GetName() string
	SetValue(value string) error
}

type IntOption struct {
	Name string
	Min  int
	Max  int
	Val  int
}

func (option *IntOption) ToUci() string {
	return fmt.Sprintf("option name %v type %v default %v min %v max %v",
		option.Name, "spin", option.Val, option.Min, option.Max)
}

func (option *IntOption) GetName() string {
	return option.Name
}

func (option *IntOption) SetValue(value string) error {
	v, err := strconv.Atoi(value)
	if err != nil {
		return errors.New("Invalid setoption arguments")
	}
	if v < option.Min || v > option.Max {
		return errors.New("argument out of range")
	}
	option.Val = v
	return nil
}

type StringOption struct {
	Name  string
	Val   string
	Dirty bool
}

func (option *StringOption) ToUci() string {
	var inspectValue string
	if option.Val == "" {
		inspectValue = "<empty>"
	} else {
		inspectValue = option.Val
	}
	return fmt.Sprintf("option name %v type %v default %v",
		option.Name, "string", inspectValue)
}

func (option *StringOption) GetName() string {
	return option.Name
}

func (option *StringOption) Clean() {
	option.Dirty = false
}

func (option *StringOption) SetValue(value string) error {
	option.Val = value
	option.Dirty = true
	return nil
}

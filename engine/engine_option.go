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
	Name    string
	Min     int
	Max     int
	Val     int
	Default int
}

func (option *IntOption) ToUci() string {
	return fmt.Sprintf("option name %s type spin default %d min %d max %d",
		option.Name, option.Default, option.Min, option.Max)
}

func (option *IntOption) GetName() string {
	return option.Name
}

func (option *IntOption) SetValue(value string) error {
	v, err := strconv.Atoi(value)
	if err != nil {
		return errors.New("invalid setoption arguments")
	}
	if v < option.Min || v > option.Max {
		return errors.New("argument out of range")
	}
	option.Val = v
	return nil
}

type StringOption struct {
	Name    string
	Val     string
	Default string
	Dirty   bool
}

func (option *StringOption) ToUci() string {
	var inspectValue string
	if option.Default == "" {
		inspectValue = "<empty>"
	} else {
		inspectValue = option.Default
	}
	return fmt.Sprintf("option name %s type string default %s", option.Name, inspectValue)
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

type CheckOption struct {
	Name    string
	Val     bool
	Default bool
}

func (option *CheckOption) ToUci() string {
	var inspectValue = "false"
	if option.Default {
		inspectValue = "true"
	}
	return fmt.Sprintf("option name %s type check default %s", option.Name, inspectValue)
}

func (option *CheckOption) GetName() string {
	return option.Name
}

func (option *CheckOption) SetValue(value string) error {
	option.Val = (value == "true")
	return nil
}

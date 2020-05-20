package main

import (
	"strconv"
	"time"
)

type IntCell struct {
	value string
}

func (ic *IntCell) IsValid() {
	if _, err := strconv.Atoi(ic.value); err == nil {
		logger.Debug("Valid integer: ", ic.value)
	} else {
		logger.Error("Invalid integer: ", ic.value)
	}
}
func (ic *IntCell) GetStrValue() string {
	return ic.value
}

type VarCharCell struct {
	value  string
	length int
}

func (vcc *VarCharCell) IsValid() {
	if len(vcc.value) <= vcc.length {
		logger.Debug("Valid string length: ", len(vcc.value), " <= ", vcc.length)
	} else {
		logger.Error("Invalid string length: ", len(vcc.value), " > ", vcc.length)
	}
}
func (vcc *VarCharCell) GetStrValue() string {
	return vcc.value
}

type DateCell struct {
	value string
}

func (dc *DateCell) IsValid() {
	if _, err := time.Parse("2006-01-02", dc.value); err == nil {
		logger.Debug("Valid date format: ", dc.value)
	} else {
		logger.Error("Invalid date format: ", dc.value)
	}
}
func (dc *DateCell) GetStrValue() string {
	return dc.value
}

type TextCell struct {
	value string
}

func (tc *TextCell) IsValid() {
	logger.Debug("Text type (No validation)")
}
func (tc *TextCell) GetStrValue() string {
	return tc.value
}

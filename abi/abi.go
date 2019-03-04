// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package abi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"reflect"
	"strconv"

	"github.com/LambdaIM/lambda-libs/common"
)

// The ABI holds information about a contract's context and available
// invokable methods. It will allow you to type check function calls and
// packs data accordingly.
type ABI struct {
	Constructor Method
	Methods     map[string]Method
	Events      map[string]Event
}

// JSON returns a parsed ABI interface and error if it failed.
func JSON(reader io.Reader) (ABI, error) {
	dec := json.NewDecoder(reader)

	var abi ABI
	if err := dec.Decode(&abi); err != nil {
		return ABI{}, err
	}

	return abi, nil
}

// Pack the given method name to conform the ABI. Method call's data
// will consist of method_id, args0, arg1, ... argN. Method id consists
// of 4 bytes and arguments are all 32 bytes.
// Method ids are created from the first 4 bytes of the hash of the
// methods string signature. (signature = baz(uint32,string32))
func (abi ABI) Pack(name string, args ...interface{}) ([]byte, error) {
	// Fetch the ABI of the requested method
	if name == "" {
		// constructor
		arguments, err := abi.Constructor.Inputs.Pack(args...)
		if err != nil {
			return nil, err
		}
		return arguments, nil

	}
	method, exist := abi.Methods[name]
	if !exist {
		return nil, fmt.Errorf("method '%s' not found", name)
	}

	arguments, err := method.Inputs.Pack(args...)
	if err != nil {
		return nil, err
	}
	// Pack up the method ID too if not a constructor and return
	return append(method.Id(), arguments...), nil
}

// Unpack output in v according to the abi specification
func (abi ABI) Unpack(v interface{}, name string, output []byte) (err error) {
	if len(output) == 0 {
		return fmt.Errorf("abi: unmarshalling empty output")
	}
	// since there can't be naming collisions with contracts and events,
	// we need to decide whether we're calling a method or an event
	if method, ok := abi.Methods[name]; ok {
		if len(output)%32 != 0 {
			return fmt.Errorf("abi: improperly formatted output")
		}
		return method.Outputs.Unpack(v, output)
	} else if event, ok := abi.Events[name]; ok {
		return event.Inputs.Unpack(v, output)
	}
	return fmt.Errorf("abi: could not locate named method or event")
}

// UnmarshalJSON implements json.Unmarshaler interface
func (abi *ABI) UnmarshalJSON(data []byte) error {
	var fields []struct {
		Type      string
		Name      string
		Constant  bool
		Anonymous bool
		Inputs    []Argument
		Outputs   []Argument
	}

	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}

	abi.Methods = make(map[string]Method)
	abi.Events = make(map[string]Event)
	for _, field := range fields {
		switch field.Type {
		case "constructor":
			abi.Constructor = Method{
				Inputs: field.Inputs,
			}
		// empty defaults to function according to the abi spec
		case "function", "":
			abi.Methods[field.Name] = Method{
				Name:    field.Name,
				Const:   field.Constant,
				Inputs:  field.Inputs,
				Outputs: field.Outputs,
			}
		case "event":
			abi.Events[field.Name] = Event{
				Name:      field.Name,
				Anonymous: field.Anonymous,
				Inputs:    field.Inputs,
			}
		}
	}

	return nil
}

// MethodById looks up a method by the 4-byte id
// returns nil if none found
func (abi *ABI) MethodById(sigdata []byte) (*Method, error) {
	for _, method := range abi.Methods {
		if bytes.Equal(method.Id(), sigdata[:4]) {
			return &method, nil
		}
	}
	return nil, fmt.Errorf("no method with id: %#x", sigdata[:4])
}

func (abi *ABI) Encode(name string, args string) ([]byte, error) {
	method, exist := abi.Methods[name]
	if !exist {
		return nil, fmt.Errorf("method '%s' not found", name)
	}

	argStack := parseArgs(args)
	var iArgs []interface{}
	for _, v := range method.Inputs {
		switch v.Type.T {
		case IntTy, UintTy:
			_sv, err := argStack.Pop()
			expectNoErr(err)
			_v, err := strconv.ParseInt(_sv.(string), 10, 0)
			expectNoErr(err)
			iArgs = append(iArgs, new(big.Int).SetInt64(_v))
		case BoolTy:
			_sv, err := argStack.Pop()
			expectNoErr(err)
			_v, err := strconv.ParseBool(_sv.(string))
			expectNoErr(err)
			iArgs = append(iArgs, _v)
		case StringTy:
			_sv, err := argStack.Pop()
			expectNoErr(err)
			iArgs = append(iArgs, _sv)
		case SliceTy:
			baseTy, _, err := parseBaseTy(v.Type)
			expectNoErr(err)
			switch baseTy {
			case "int", "uint":
				arr, err := parseNums(argStack)
				expectNoErr(err)
				iArgs = append(iArgs, arr)
			default:
				return nil, fmt.Errorf("unsupported type %v", v.Type.String())
			}
		case ArrayTy:
			baseTy, _, err := parseBaseTy(v.Type)
			expectNoErr(err)
			switch baseTy {
			case "int", "uint":
				arr, err := parseNums(argStack)
				expectNoErr(err)
				iArgs = append(iArgs, arr)
			default:
				return nil, fmt.Errorf("unsupported type %v", v.Type.String())
			}
		case AddressTy:
			_sv, err := argStack.Pop()
			expectNoErr(err)
			iArgs = append(iArgs, common.BytesToAddress(common.FromHex(_sv.(string))))
		case FixedBytesTy:
			_sv, err := argStack.Pop()
			expectNoErr(err)
			_nv := reflect.New(reflect.ArrayOf(v.Type.Size, reflect.TypeOf(uint8(0)))).Elem()
			_nv1 := []byte(_sv.(string))
			expectTrue(_nv.Len() > len(_nv1))
			iv := _nv.Len() - 1
			for index := len(_nv1) - 1; index >= 0; index = index - 1 {
				_nv.Index(iv).Set(reflect.ValueOf(uint8(_nv1[index])))
				iv = iv - 1
			}
			iArgs = append(iArgs, _nv.Interface())
		case BytesTy:
			_sv, err := argStack.Pop()
			expectNoErr(err)

			_bv := []byte(_sv.(string))
			_bvSize := len(_bv)
			_nv := reflect.MakeSlice(reflect.TypeOf([]byte{}), _bvSize, _bvSize)
			_inv := _nv.Len() - 1
			for index := len(_bv) - 1; index >= 0; index = index - 1 {
				_nv.Index(_inv).Set(reflect.ValueOf(uint8(_bv[index])))
				_inv = _inv - 1
			}
			iArgs = append(iArgs, _nv.Interface())
		default:
			return nil, fmt.Errorf("unsupported type %v", v.Type.String())
		}
	}

	return abi.Pack(name, iArgs...)
}

func (abi *ABI) Decode(name string, output []byte) ([]interface{}, error) {
	if len(output) == 0 {
		return nil, fmt.Errorf("abi: unmarshalling empty output")
	}
	// since there can't be naming collisions with contracts and events,
	// we need to decide whether we're calling a method or an event
	if method, ok := abi.Methods[name]; ok {
		if len(output)%32 != 0 {
			return nil, fmt.Errorf("abi: improperly formatted output")
		}
		return method.Outputs.UnpackValues(output)
	} else if event, ok := abi.Events[name]; ok {
		return event.Inputs.UnpackValues(output)
	}

	return nil, fmt.Errorf("abi: could not locate named method or event")
}

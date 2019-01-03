package abi

import (
	"fmt"
	"math/big"
	"regexp"
	"strconv"
	"strings"
)

func parseArgs(args string) *Stack {
	s := NewStack()

	parasRegex := regexp.MustCompile("\\[|(\\w+)|\\]+")
	argsString := parasRegex.FindAllString(args, -1)

	for index := len(argsString) - 1; index >= 0; index = index - 1 {
		s.Push(argsString[index])
	}

	return s
}

func expect(val string, exp string) {
	if strings.Compare(val, exp) != 0 {
		panic(fmt.Errorf("expected value : %v %v", val, exp))
	}
}

func expectTrue(exp bool) {
	if !exp {
		panic(fmt.Errorf("expected  true : %v", exp))
	}
}

func expectNoErr(err error) {
	if err != nil {
		panic(err)
	}
}

func parseBaseTy(ty Type) (string, int64, error) {
	tElem := ty.Elem.String()
	// we do not support multi-array/slice type
	if strings.Count(tElem, "[") != 0 || strings.Count(tElem, "]") != 0 {
		return "", 0, fmt.Errorf("invalid type in abi %v", tElem)
	}

	// grab the array/slice size with regexp
	re := regexp.MustCompile("^(bytes|int|uint)([0-9]*)")
	tyz := re.FindAllStringSubmatch(tElem, -1)
	if tyz == nil {
		return "", 0, fmt.Errorf("invalid type in abi %v", tElem)
	}
	baseTy := tyz[0][1]
	baseTySize := tyz[0][2]

	if len(baseTySize) == 0 {
		return baseTy, 0, nil
	}

	size, err := strconv.ParseInt(baseTySize, 10, 0)
	expectNoErr(err)

	return baseTy, size, nil
}

func parseNums(stack *Stack) ([]*big.Int, error) {
	_sv, err := stack.Pop()
	expectNoErr(err)
	expect("[", _sv.(string))
	arr := make([]*big.Int, 0)
	for {
		_sv, err := stack.Pop()
		expectNoErr(err)
		if strings.Compare(_sv.(string), "]") == 0 {
			break
		}
		_v, err := strconv.ParseInt(_sv.(string), 10, 0)
		expectNoErr(err)
		arr = append(arr, new(big.Int).SetInt64(_v))
	}

	return arr, nil
}

func parseBytes(stack *Stack, size int64) ([]string, error) {
	_sv, err := stack.Pop()
	expectNoErr(err)
	expect("[", _sv.(string))

	arr := make([]string, 0)
	for {
		_sv, err := stack.Pop()
		expectNoErr(err)
		if strings.Compare(_sv.(string), "]") == 0 {
			break
		}
		arr = append(arr, _sv.(string))
	}

	return arr, nil
}

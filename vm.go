package main

import "fmt"

type Op byte

const (
	OpConstant Op = iota
	OpNil
	OpFalse
	OpTrue
	OpNegate
	OpNot
	OpAdd
	OpSubtract
	OpMultiply
	OpDivide
	OpEqual
	OpGreater
	OpLess
	OpReturn
)

type Stack struct {
	vals []Value
}

func newStack() *Stack {
	return &Stack{}
}

func (s *Stack) push(val Value) {
	s.vals = append(s.vals, val)
}

func (s *Stack) pop() Value {
	n := len(s.vals) - 1
	val := s.vals[n]
	s.vals = s.vals[:n]
	return val
}

type Chunk struct {
	code []byte
	vals []Value
}

func (c *Chunk) addByte(b byte) {
	c.code = append(c.code, b)
}

func (c *Chunk) addOp(op Op) {
	c.addByte(byte(op))
}

func (c *Chunk) addVal(val Value) int {
	c.vals = append(c.vals, val)
	return len(c.vals) - 1
}

func dumpChunk(c *Chunk, title string) {
	fmt.Printf("== %s\n", title)
	for i := 0; i < len(c.code); {
		i += dumpOp(c, i)
	}
}

func dumpOp(c *Chunk, offset int) int {
	op := Op(c.code[offset])

	fmt.Printf("%04d %s", offset, op)
	defer fmt.Println()

	switch op {
	case OpConstant:
		val := c.code[offset+1]
		fmt.Printf(" %3d [%s]", val, c.vals[val])
		return 2
	}

	return 1
}

type VM interface {
	run(chunk *Chunk) error
}

type vm struct{}

func newVM() VM {
	return vm{}
}

func (vm vm) run(chunk *Chunk) error {
	stack := newStack()

	literal := func(v Value) error {
		stack.push(v)
		return nil
	}

	unary := func(fn func(Value) (Value, error)) error {
		v := stack.pop()
		res, err := fn(v)
		if err == nil {
			stack.push(res)
		}
		return err
	}

	binary := func(fn func(Value, Value) (Value, error)) error {
		b := stack.pop()
		a := stack.pop()
		res, err := fn(a, b)
		if err == nil {
			stack.push(res)
		}
		return err
	}

	for ip := 0; ip < len(chunk.code); ip++ {
		dumpOp(chunk, ip)
		op := Op(chunk.code[ip])

		var err error

		switch op {
		case OpConstant:
			ip++
			err = literal(chunk.vals[chunk.code[ip]])
		case OpNil:
			err = literal(nilValue())
		case OpFalse:
			err = literal(boolValue(false))
		case OpTrue:
			err = literal(boolValue(true))
		case OpNegate:
			err = unary(negateValue)
		case OpNot:
			err = unary(notValue)
		case OpAdd:
			err = binary(addValues)
		case OpSubtract:
			err = binary(subtractValues)
		case OpMultiply:
			err = binary(multiplyValues)
		case OpDivide:
			err = binary(divideValues)
		case OpEqual:
			err = binary(valuesEqual)
		case OpGreater:
			err = binary(valueGreater)
		case OpLess:
			err = binary(valueLess)
		case OpReturn:
			fmt.Println(stack.pop())
		default:
			err = fmt.Errorf("unknown op: %q\n", op)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

//go:generate stringer -type=Op

package interpreter

import (
	"fmt"
	"strconv"
	"strings"

	"ksm/parser"
)

type Interpreter struct {
	variables map[string]string
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		variables: make(map[string]string),
	}
}

func (i *Interpreter) Interpret(node *parser.Node) error {
	if node == nil {
		return nil
	}

	switch node.Type {
	case parser.NodeVarDecl:
		parts := strings.Split(node.Literal, "=")
		if len(parts) != 2 {
			return fmt.Errorf("invalid variable declaration: %s", node.Literal)
		}
		varName := strings.TrimSpace(parts[0])
		varValue := strings.TrimSpace(parts[1])
		i.variables[varName] = varValue
		fmt.Printf("Variable Declaration: %s = %s\n", varName, varValue)

	case parser.NodePrint:
		value := strings.TrimPrefix(node.Literal, "print ")
		fmt.Printf("Print Statement: %s\n", i.evaluateExpression(value))

	case parser.NodeIf:
		condition := strings.TrimPrefix(node.Literal, "if ")
		if i.evaluateCondition(condition) {
			fmt.Println("If Statement (True):", condition)
			for _, child := range node.Children {
				if err := i.Interpret(child); err != nil {
					return err
				}
			}
		} else {
			fmt.Println("If Statement (False):", condition)
		}

	case parser.NodeOtherwise:
		fmt.Println("Otherwise Statement")
		for _, child := range node.Children {
			if err := i.Interpret(child); err != nil {
				return err
			}
		}

	case parser.NodeBlock:
		for _, child := range node.Children {
			if err := i.Interpret(child); err != nil {
				return err
			}
		}
	}

	return nil
}

func (i *Interpreter) evaluateExpression(expr string) string {
	if val, ok := i.variables[expr]; ok {
		return val
	}
	return expr
}

func (i *Interpreter) evaluateCondition(condition string) bool {
	parts := strings.Split(condition, " ")
	if len(parts) != 3 {
		return false
	}

	left := i.evaluateExpression(parts[0])
	operator := parts[1]
	right := i.evaluateExpression(parts[2])

	leftNum, leftErr := strconv.Atoi(left)
	rightNum, rightErr := strconv.Atoi(right)

	if leftErr == nil && rightErr == nil {
		switch operator {
		case ">":
			return leftNum > rightNum
		case "<":
			return leftNum < rightNum
		case "==":
			return leftNum == rightNum
		}
	}

	return false
}

package evaluator

import (
	"fmt"

	"pidgin-lang/ast"
	"pidgin-lang/object"
)

// Singleton objects for efficiency
var (
	NOTHING = &object.Nothing{}
	TRU     = &object.Boolean{Value: true}
	LIE     = &object.Boolean{Value: false}
)

// Eval evaluates an AST node and returns an object
func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {

	// Statements
	case *ast.Program:
		return evalProgram(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.MakeStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
		return val

	case *ast.BringStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.NothingLiteral:
		return NOTHING

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)

	case *ast.SupposeExpression:
		return evalSupposeExpression(node, env)

	case *ast.WhileExpression:
		return evalWhileExpression(node, env)

	case *ast.DoExpression:
		return evalDoExpression(node, env)

	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args)
	}

	return nil
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	// Check user-defined variables first
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	// Check built-in functions
	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError("I no sabi dis one: %s", node.Value)
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

// =============================================================================
// Prefix Expressions: -5, !tru, no be x
// =============================================================================

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	case "no be":
		return evalBangOperatorExpression(right)
	default:
		return newError("I no understand dis operator: %s%s", operator, right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRU:
		return LIE
	case LIE:
		return TRU
	case NOTHING:
		return TRU
	default:
		return LIE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("I no fit minus dis one: %s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

// =============================================================================
// Infix Expressions: 5 + 5, x big pass y, etc.
// =============================================================================

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	case left.Type() == object.STRING_OBJ || right.Type() == object.STRING_OBJ:
		// Allow string concatenation with + when one side is a string
		if operator == "+" {
			return evalStringConcatenation(left, right)
		}
		return newError("I no fit do %s wit %s and %s", operator, left.Type(), right.Type())
	case operator == "be":
		return nativeBoolToBooleanObject(left == right)
	case operator == "na":
		return nativeBoolToBooleanObject(left == right)
	case operator == "no be":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("I no fit compare %s wit %s", left.Type(), right.Type())
	default:
		return newError("I no understand dis operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		if rightVal == 0 {
			return newError("Omo! You no fit divide by zero o!")
		}
		return &object.Integer{Value: leftVal / rightVal}
	case "big pass":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "no reach":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "be":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "na":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	default:
		return newError("I no understand dis operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	switch operator {
	case "+":
		return &object.String{Value: leftVal + rightVal}
	case "be":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "na":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	default:
		return newError("I no understand dis operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringConcatenation(left, right object.Object) object.Object {
	var leftVal, rightVal string

	switch l := left.(type) {
	case *object.String:
		leftVal = l.Value
	case *object.Integer:
		leftVal = fmt.Sprintf("%d", l.Value)
	case *object.Boolean:
		leftVal = l.Inspect()
	default:
		leftVal = l.Inspect()
	}

	switch r := right.(type) {
	case *object.String:
		rightVal = r.Value
	case *object.Integer:
		rightVal = fmt.Sprintf("%d", r.Value)
	case *object.Boolean:
		rightVal = r.Inspect()
	default:
		rightVal = r.Inspect()
	}

	return &object.String{Value: leftVal + rightVal}
}

// =============================================================================
// Control Flow: suppose/abi, dey do while
// =============================================================================

func evalSupposeExpression(se *ast.SupposeExpression, env *object.Environment) object.Object {
	condition := Eval(se.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(se.Consequence, env)
	} else if se.Alternative != nil {
		return Eval(se.Alternative, env)
	} else {
		return NOTHING
	}
}

func evalWhileExpression(we *ast.WhileExpression, env *object.Environment) object.Object {
	var result object.Object = NOTHING

	for {
		condition := Eval(we.Condition, env)
		if isError(condition) {
			return condition
		}

		if !isTruthy(condition) {
			break
		}

		result = Eval(we.Body, env)

		// Check for return or error
		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

// =============================================================================
// Functions
// =============================================================================

func evalDoExpression(de *ast.DoExpression, env *object.Environment) object.Object {
	fn := &object.Function{
		Parameters: de.Parameters,
		Body:       de.Body,
		Env:        env,
	}

	// If function has a name, bind it to the environment
	if de.Name != nil {
		fn.Name = de.Name.Value
		env.Set(de.Name.Value, fn)
	}

	return fn
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)

	case *object.Builtin:
		return fn.Fn(args...)

	default:
		return newError("Dis one no be function: %s", fn.Type())
	}
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		if paramIdx < len(args) {
			env.Set(param.Value, args[paramIdx])
		}
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

// =============================================================================
// Built-in Functions
// =============================================================================

var builtins = map[string]*object.Builtin{
	"yarn": {
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			return NOTHING
		},
	},
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("len wan make one argument, you give am %d", len(args))
			}
			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return newError("I no fit check length of %s", args[0].Type())
			}
		},
	},
	"type": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("type wan make one argument, you give am %d", len(args))
			}
			return &object.String{Value: string(args[0].Type())}
		},
	},
}

// =============================================================================
// Helpers
// =============================================================================

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRU
	}
	return LIE
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NOTHING:
		return false
	case TRU:
		return true
	case LIE:
		return false
	default:
		return true
	}
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}
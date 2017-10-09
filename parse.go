package gval

import (
	"fmt"
	"strconv"
	"text/scanner"
)

//ParseExpression scans an expression into an Evaluable.
func (p *Parser) ParseExpression() (eval Evaluable, err error) {
	stack := stageStack{}
	for {
		eval, err = p.ParseNextExpression()
		if err != nil {
			return nil, err
		}

		if stage, err := p.parseOperator(&stack, eval); err != nil {
			return nil, err
		} else if err = stack.push(stage); err != nil {
			return nil, err
		}

		if stack.peek().infixBuilder == nil {
			return stack.pop().Evaluable, nil
		}
	}
}

//ParseNextExpression scans the expression ignoring following operators
func (p *Parser) ParseNextExpression() (eval Evaluable, err error) {
	scan := p.Scan()
	ex, ok := p.prefixes[scan]
	if !ok {
		return nil, p.Expected("extensions")
	}
	return ex(p)
}

func parseString(p *Parser) (Evaluable, error) {
	s, err := strconv.Unquote(p.TokenText())
	if err != nil {
		return nil, fmt.Errorf("could not parse string: %s", err)
	}
	return p.Const(s), nil
}

func parseNumber(p *Parser) (Evaluable, error) {
	n, err := strconv.ParseFloat(p.TokenText(), 64)
	if err != nil {
		return nil, err
	}
	return p.Const(n), nil
}

func parseParentheses(p *Parser) (Evaluable, error) {
	eval, err := p.ParseExpression()
	if err != nil {
		return nil, err
	}
	switch p.Scan() {
	case ')':
		return eval, nil
	default:
		return nil, p.Expected("parentheses", ')')
	}
}

func (p *Parser) parseOperator(stack *stageStack, eval Evaluable) (st stage, err error) {
	for {
		c := p.Scan()
		op := p.TokenText()
		mustOp := false
		if isSymbolOperation(c) {
			c = p.Peek()
			for isSymbolOperation(c) {
				mustOp = true
				op += string(c)
				p.Next()
				c = p.Peek()
			}
		} else if c != scanner.Ident {
			p.Camouflage("operator")
			return stage{Evaluable: eval}, nil
		}
		operator, _ := p.operators[op]
		switch operator := operator.(type) {
		case *infix:
			return stage{
				Evaluable:          eval,
				infixBuilder:       operator.builder,
				operatorPrecedence: operator.operatorPrecedence,
			}, nil
		case directInfix:
			return stage{
				Evaluable:          eval,
				infixBuilder:       operator.infixBuilder,
				operatorPrecedence: operator.operatorPrecedence,
			}, nil
		case postfix:
			if err = stack.push(stage{
				operatorPrecedence: operator.operatorPrecedence,
				Evaluable:          eval,
			}); err != nil {
				return stage{}, err
			}
			eval, err = operator.f(p, stack.pop().Evaluable, operator.operatorPrecedence)
			if err != nil {
				return
			}
			continue
		}

		if !mustOp {
			p.Camouflage("operator")
			return stage{Evaluable: eval}, nil
		}
		return stage{}, fmt.Errorf("unknown operator %s", op)
	}
}

func parseIdent(p *Parser) (call string, alternative func() (Evaluable, error), err error) {
	token := p.TokenText()
	return token,
		func() (Evaluable, error) {
			fullname := token

			keys := []Evaluable{p.Const(token)}
			for {
				scan := p.Scan()
				switch scan {
				case '.':
					scan = p.Scan()
					switch scan {
					case scanner.Ident:
						token = p.TokenText()
						keys = append(keys, p.Const(token))
					default:
						return nil, p.Expected("field", scanner.Ident)
					}
				case '(':
					args, err := p.parseArguments()
					if err != nil {
						return nil, err
					}
					return p.callEvaluable(fullname, p.Var(keys...), args...), nil
				case '[':
					key, err := p.ParseExpression()
					if err != nil {
						return nil, err
					}
					switch p.Scan() {
					case ']':
						keys = append(keys, key)
					default:
						return nil, p.Expected("array key", ']')
					}
				default:
					p.Camouflage("variable", '.', '(', '[')
					return p.Var(keys...), nil
				}
			}
		}, nil

}

func (p *Parser) parseArguments() (args []Evaluable, err error) {
	if p.Scan() == ')' {
		return
	}
	p.Camouflage("scan arguments", ')')
	for {
		arg, err := p.ParseExpression()
		args = append(args, arg)
		if err != nil {
			return nil, err
		}
		switch p.Scan() {
		case ')':
			return args, nil
		case ',':
		default:
			return nil, p.Expected("arguments", ')', ',')
		}
	}
}

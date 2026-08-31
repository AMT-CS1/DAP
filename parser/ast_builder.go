package parser

import (
	"dap/ast"
	"dap/scanner"
	"fmt"
	"strings"
)

type ASTBuilder struct {
	tokens []scanner.Token
	pos    int
}

// BuildAST parses the token stream into an AST ListNode representation.
func BuildAST(t *scanner.Token) (*ast.ListNode, error) {
	var tokenList []scanner.Token

	for {
		typ := t.Next()
		tokCopy := scanner.Token{
			Typ:   t.Typ,
			Val:   t.Val,
			Lno:   t.Lno,
			Cno:   t.Cno,
			First: t.First,
		}
		tokenList = append(tokenList, tokCopy)

		if typ == "$ENDPROG" || typ == "" {
			break
		}
	}

	b := &ASTBuilder{
		tokens: tokenList,
		pos:    0,
	}

	return b.parseProgram()
}

func (b *ASTBuilder) cur() scanner.Token {
	if b.pos < len(b.tokens) {
		return b.tokens[b.pos]
	}
	return scanner.Token{Typ: "$ENDPROG", Val: "END OF FILE"}
}

func (b *ASTBuilder) nextToken() scanner.Token {
	c := b.cur()
	if b.pos < len(b.tokens) {
		b.pos++
	}
	return c
}

func (b *ASTBuilder) peekToken() scanner.Token {
	return b.cur()
}

func (b *ASTBuilder) lookahead(n int) scanner.Token {
	idx := b.pos + n
	if idx < len(b.tokens) {
		return b.tokens[idx]
	}
	return scanner.Token{Typ: "$ENDPROG", Val: "END OF FILE"}
}

func (b *ASTBuilder) parseProgram() (*ast.ListNode, error) {
	programNode := &ast.ListNode{
		Type:    "ListNode",
		Program: []interface{}{},
	}

	if b.peekToken().Typ == "$PROGRAM" {
		b.nextToken() // consume $PROGRAM
		if b.peekToken().Typ == "$NAME" {
			b.nextToken() // consume program name
		}
	}

	for b.pos < len(b.tokens) {
		tok := b.peekToken()
		if tok.Typ == "$ENDPROG" || tok.Typ == "" {
			break
		}

		if tok.Typ == "$DICT" || tok.Typ == "$LOCAL" || tok.Typ == "$GLOBAL" {
			b.nextToken() // consume dict header
			dictNode := b.parseDictionary()
			if len(dictNode.Variables) > 0 {
				programNode.Program = append(programNode.Program, dictNode)
			}
			continue
		}

		if tok.Typ == "$CODE" {
			b.nextToken() // consume $CODE
			continue
		}

		if strings.EqualFold(tok.Val, "endprogram") || (strings.EqualFold(tok.Val, "end") && b.lookahead(1).Typ == "$PROGRAM") {
			break
		}

		stmt := b.parseStatement()
		if stmt != nil {
			programNode.Program = append(programNode.Program, stmt)
		} else {
			b.nextToken() // advance on unhandled token to avoid infinite loop
		}
	}

	return programNode, nil
}

func (b *ASTBuilder) parseDictionary() *ast.DictionaryNode {
	dictNode := &ast.DictionaryNode{
		Type:      "DictionaryNode",
		Variables: []interface{}{},
	}

	for b.pos < len(b.tokens) {
		tok := b.peekToken()
		if tok.Typ == "$CODE" || tok.Typ == "$ENDPROG" || tok.Typ == "" {
			break
		}

		if tok.Typ == "$CONST" {
			b.nextToken() // consume $CONST
			if b.peekToken().Typ == "$NAME" {
				cname := b.nextToken().Val
				if b.peekToken().Typ == "$COLON" {
					b.nextToken()
					if b.peekToken().Typ != "" {
						b.nextToken() // type
					}
				}
				if b.peekToken().Typ == "$ASSG" || b.peekToken().Typ == "$EQ" {
					b.nextToken()
				}
				valExpr := b.parseExpression()
				dictNode.Variables = append(dictNode.Variables, ast.VarDeclNode{
					Type:          "VarDeclNode",
					Name:          cname,
					Value:         valExpr,
					IsConst:       true,
					IsDeclaration: false,
				})
			}
			continue
		}

		if tok.Typ == "$VAR" {
			b.nextToken() // consume $VAR
		}

		if b.peekToken().Typ == "$NAME" {
			var names []string
			for {
				if b.peekToken().Typ == "$NAME" {
					names = append(names, b.nextToken().Val)
				}
				if b.peekToken().Typ == "$COMMA" {
					b.nextToken() // consume comma
				} else {
					break
				}
			}

			if b.peekToken().Typ == "$COLON" {
				b.nextToken() // consume colon
				typeTok := b.nextToken()
				typeStr := b.mapTypeName(typeTok.Val, typeTok.Typ)

				var initExpr interface{}
				if b.peekToken().Typ == "$ASSG" || b.peekToken().Typ == "$EQ" {
					b.nextToken()
					initExpr = b.parseExpression()
				}

				for _, name := range names {
					var valVal interface{}
					if initExpr != nil {
						valVal = initExpr
					} else {
						valVal = ast.VarAccessToken{
							Type: "VarAccessToken",
							Name: typeStr,
						}
					}
					dictNode.Variables = append(dictNode.Variables, ast.VarDeclNode{
						Type:          "VarDeclNode",
						Name:          name,
						Value:         valVal,
						IsConst:       false,
						IsDeclaration: true,
					})
				}
			}
			continue
		}

		break
	}

	return dictNode
}

func (b *ASTBuilder) mapTypeName(raw string, typ string) string {
	switch typ {
	case "$INT":
		return "integer"
	case "$REAL":
		return "real"
	case "$CHAR":
		return "character"
	case "$BOOL":
		return "boolean"
	case "$CHARRAY":
		return "string"
	default:
		if raw != "" {
			return raw
		}
		return "integer"
	}
}

func (b *ASTBuilder) parseStatementList(stopTokens ...string) *ast.ListNode {
	listNode := &ast.ListNode{
		Type:    "ListNode",
		Program: []interface{}{},
	}

	for b.pos < len(b.tokens) {
		tok := b.peekToken()
		if tok.Typ == "" || tok.Typ == "$ENDPROG" {
			break
		}

		isStop := false
		for _, stop := range stopTokens {
			if tok.Typ == stop || strings.EqualFold(tok.Val, stop) {
				isStop = true
				break
			}
		}
		if isStop {
			break
		}

		if strings.EqualFold(tok.Val, "end") {
			peekNext := b.lookahead(1).Typ
			if peekNext == "$IF" || peekNext == "$FOR" || peekNext == "$WHILE" || peekNext == "$CASE" {
				break
			}
		}

		stmt := b.parseStatement()
		if stmt != nil {
			listNode.Program = append(listNode.Program, stmt)
		} else {
			b.nextToken()
		}
	}

	return listNode
}

func (b *ASTBuilder) parseStatement() interface{} {
	tok := b.peekToken()

	switch tok.Typ {
	case "$INPUT", "$OUTPUT", "$CALL":
		b.nextToken()
		funcName := strings.ToLower(tok.Val)
		if tok.Typ == "$INPUT" {
			funcName = "input"
		} else if tok.Typ == "$OUTPUT" {
			funcName = "output"
		}

		args := []interface{}{}
		if b.peekToken().Typ == "$LEFTPAR" {
			b.nextToken() // consume '('
			for {
				if b.peekToken().Typ == "$RIGHTPAR" || b.peekToken().Typ == "$ENDPROG" || b.peekToken().Typ == "" {
					b.nextToken()
					break
				}
				arg := b.parseExpression()
				if arg != nil {
					args = append(args, arg)
				}
				if b.peekToken().Typ == "$COMMA" {
					b.nextToken()
				} else if b.peekToken().Typ == "$RIGHTPAR" {
					b.nextToken()
					break
				} else {
					break
				}
			}
		}
		return ast.CallNode{
			Type: "CallNode",
			Call: ast.VarAccessToken{
				Type: "VarAccessToken",
				Name: funcName,
			},
			Args: args,
		}

	case "$NAME":
		varName := b.nextToken().Val
		nextTok := b.peekToken()
		if nextTok.Typ == "$ASSG" || nextTok.Typ == "$EQ" {
			b.nextToken() // consume assign
			valExpr := b.parseExpression()
			return ast.VarAssignNode{
				Type:          "VarAssignNode",
				Name:          varName,
				Value:         valExpr,
				IsConst:       false,
				IsDeclaration: false,
			}
		} else if nextTok.Typ == "$LEFTPAR" {
			b.nextToken() // consume '('
			args := []interface{}{}
			for {
				if b.peekToken().Typ == "$RIGHTPAR" || b.peekToken().Typ == "$ENDPROG" || b.peekToken().Typ == "" {
					b.nextToken()
					break
				}
				arg := b.parseExpression()
				if arg != nil {
					args = append(args, arg)
				}
				if b.peekToken().Typ == "$COMMA" {
					b.nextToken()
				} else if b.peekToken().Typ == "$RIGHTPAR" {
					b.nextToken()
					break
				} else {
					break
				}
			}
			return ast.CallNode{
				Type: "CallNode",
				Call: ast.VarAccessToken{
					Type: "VarAccessToken",
					Name: varName,
				},
				Args: args,
			}
		}

	case "$IF":
		b.nextToken() // consume $IF
		cond := b.parseExpression()
		if b.peekToken().Typ == "$THEN" {
			b.nextToken()
		}

		thenBody := b.parseStatementList("$ELIF", "$ELSE", "$ENDIF")

		ifNode := &ast.IfNode{
			Type: "IfNode",
			Cases: []ast.IfCase{
				{
					Condition: cond,
					Body:      thenBody,
				},
			},
		}

		tokNext := b.peekToken()
		if tokNext.Typ == "$ELIF" {
			b.nextToken()
			elseIfNode := b.parseStatement()
			ifNode.ElseCase = &ast.ElseCase{
				Body: elseIfNode,
			}
		} else if tokNext.Typ == "$ELSE" {
			b.nextToken()
			elseBody := b.parseStatementList("$ENDIF")
			ifNode.ElseCase = &ast.ElseCase{
				Body: elseBody,
			}
		}

		if b.peekToken().Typ == "$ENDIF" {
			b.nextToken()
		} else if strings.EqualFold(b.peekToken().Val, "end") && b.lookahead(1).Typ == "$IF" {
			b.nextToken()
			b.nextToken()
		}

		return ifNode

	case "$WHILE":
		b.nextToken() // consume $WHILE
		cond := b.parseExpression()
		if b.peekToken().Typ == "$DO" {
			b.nextToken()
		}
		body := b.parseStatementList("$ENDWHILE")

		if b.peekToken().Typ == "$ENDWHILE" {
			b.nextToken()
		} else if strings.EqualFold(b.peekToken().Val, "end") && b.lookahead(1).Typ == "$WHILE" {
			b.nextToken()
			b.nextToken()
		}

		return ast.WhileNode{
			Type:      "WhileNode",
			Condition: cond,
			Body:      body,
		}

	case "$FOR":
		b.nextToken() // consume $FOR
		varName := ""
		if b.peekToken().Typ == "$NAME" {
			varName = b.nextToken().Val
		}
		if b.peekToken().Typ == "$ASSG" || b.peekToken().Typ == "$EQ" {
			b.nextToken()
		}
		startExpr := b.parseExpression()

		direction := "to"
		dirTok := b.nextToken()
		if dirTok.Typ == "$DOWNTO" {
			direction = "downto"
		} else if dirTok.Typ == "$TO" {
			direction = "to"
		}

		endExpr := b.parseExpression()
		if b.peekToken().Typ == "$DO" {
			b.nextToken()
		}

		body := b.parseStatementList("$ENDFOR")

		if b.peekToken().Typ == "$ENDFOR" {
			b.nextToken()
		} else if strings.EqualFold(b.peekToken().Val, "end") && b.lookahead(1).Typ == "$FOR" {
			b.nextToken()
			b.nextToken()
		}

		return ast.ForNode{
			Type:      "ForNode",
			Variable:  varName,
			Start:     startExpr,
			End:       endExpr,
			Direction: direction,
			Body:      body,
		}

	case "$REPEAT":
		b.nextToken() // consume $REPEAT
		body := b.parseStatementList("$UNTIL")
		if b.peekToken().Typ == "$UNTIL" {
			b.nextToken()
		}
		cond := b.parseExpression()

		return ast.RepeatNode{
			Type:      "RepeatNode",
			Condition: cond,
			Body:      body,
		}
	}

	return nil
}

func (b *ASTBuilder) parseExpression() interface{} {
	return b.parseOr()
}

func (b *ASTBuilder) parseOr() interface{} {
	left := b.parseAnd()
	for {
		tok := b.peekToken()
		if tok.Typ == "$OR" {
			b.nextToken()
			right := b.parseAnd()
			left = ast.BinOpNode{
				Type:     "BinOpNode",
				Left:     left,
				Operator: "or",
				Right:    right,
			}
		} else {
			break
		}
	}
	return left
}

func (b *ASTBuilder) parseAnd() interface{} {
	left := b.parseRelational()
	for {
		tok := b.peekToken()
		if tok.Typ == "$AND" {
			b.nextToken()
			right := b.parseRelational()
			left = ast.BinOpNode{
				Type:     "BinOpNode",
				Left:     left,
				Operator: "and",
				Right:    right,
			}
		} else {
			break
		}
	}
	return left
}

func (b *ASTBuilder) parseRelational() interface{} {
	left := b.parseAdditive()
	for {
		tok := b.peekToken()
		op := ""
		switch tok.Typ {
		case "$EQ":
			op = "=="
		case "$NE":
			op = "!="
		case "$LT":
			op = "<"
		case "$GT":
			op = ">"
		case "$LE":
			op = "<="
		case "$GE":
			op = ">="
		}

		if op != "" {
			b.nextToken()
			right := b.parseAdditive()
			left = ast.BinOpNode{
				Type:     "BinOpNode",
				Left:     left,
				Operator: op,
				Right:    right,
			}
		} else {
			break
		}
	}
	return left
}

func (b *ASTBuilder) parseAdditive() interface{} {
	left := b.parseMultiplicative()
	for {
		tok := b.peekToken()
		op := ""
		if tok.Typ == "$PLUS" {
			op = "+"
		} else if tok.Typ == "$MINUS" {
			op = "-"
		}

		if op != "" {
			b.nextToken()
			right := b.parseMultiplicative()
			left = ast.BinOpNode{
				Type:     "BinOpNode",
				Left:     left,
				Operator: op,
				Right:    right,
			}
		} else {
			break
		}
	}
	return left
}

func (b *ASTBuilder) parseMultiplicative() interface{} {
	left := b.parseUnary()
	for {
		tok := b.peekToken()
		op := ""
		switch tok.Typ {
		case "$MULT":
			op = "*"
		case "$RDIV":
			op = "/"
		case "$DIV":
			op = "div"
		case "$MOD":
			op = "mod"
		}

		if op != "" {
			b.nextToken()
			right := b.parseUnary()
			left = ast.BinOpNode{
				Type:     "BinOpNode",
				Left:     left,
				Operator: op,
				Right:    right,
			}
		} else {
			break
		}
	}
	return left
}

func (b *ASTBuilder) parseUnary() interface{} {
	tok := b.peekToken()
	if tok.Typ == "$MINUS" || tok.Typ == "$NOT" {
		b.nextToken()
		op := "-"
		if tok.Typ == "$NOT" {
			op = "not"
		}
		subNode := b.parseUnary()
		return ast.UnaryOpNode{
			Type:     "UnaryOpNode",
			Operator: op,
			Node:     subNode,
		}
	}
	return b.parsePrimary()
}

func (b *ASTBuilder) parsePrimary() interface{} {
	tok := b.nextToken()

	switch tok.Typ {
	case "$INT", "$REAL":
		return ast.NumberNode{
			Type:  "NumberNode",
			Value: tok.Val,
		}
	case "$STRING", "$CHARRAY":
		val := strings.Trim(tok.Val, "\"")
		return ast.StringNode{
			Type:  "StringNode",
			Value: val,
		}
	case "$CHAR":
		val := strings.Trim(tok.Val, "'")
		return ast.CharNode{
			Type:  "CharNode",
			Value: val,
		}
	case "$TRUE", "$FALSE":
		return ast.BoolNode{
			Type:  "BoolNode",
			Value: fmt.Sprintf("%v", tok.Typ == "$TRUE"),
		}
	case "$NAME":
		return ast.VarAccessToken{
			Type: "VarAccessToken",
			Name: tok.Val,
		}
	case "$LEFTPAR":
		expr := b.parseExpression()
		if b.peekToken().Typ == "$RIGHTPAR" {
			b.nextToken()
		}
		return expr
	}

	return ast.VarAccessToken{
		Type: "VarAccessToken",
		Name: tok.Val,
	}
}

package ast

// ListNode represents the top-level program or a block of statements.
type ListNode struct {
	Type    string        `json:"type"`
	Program []interface{} `json:"program"`
}

// DictionaryNode represents variable and constant declarations.
type DictionaryNode struct {
	Type      string        `json:"type"`
	Variables []interface{} `json:"variables"`
}

// VarDeclNode represents a variable or constant declaration entry.
type VarDeclNode struct {
	Type          string      `json:"type"`
	Name          string      `json:"name"`
	Value         interface{} `json:"value"`
	IsConst       bool        `json:"is_const"`
	IsDeclaration bool        `json:"is_declaration"`
}

// VarAssignNode represents variable assignment.
type VarAssignNode struct {
	Type          string      `json:"type"`
	Name          string      `json:"name"`
	Value         interface{} `json:"value"`
	IsConst       bool        `json:"is_const"`
	IsDeclaration bool        `json:"is_declaration"`
}

// CallNode represents procedure/function calls (e.g. input, output).
type CallNode struct {
	Type string        `json:"type"`
	Call VarAccessToken `json:"call"`
	Args []interface{} `json:"args"`
}

// BinOpNode represents binary operations (+, -, *, /, div, mod, ==, !=, <, >, <=, >=, and, or).
type BinOpNode struct {
	Type     string      `json:"type"`
	Left     interface{} `json:"left"`
	Operator string      `json:"operator"`
	Right    interface{} `json:"right"`
}

// UnaryOpNode represents unary operations (-, not).
type UnaryOpNode struct {
	Type     string      `json:"type"`
	Operator string      `json:"operator"`
	Node     interface{} `json:"node"`
}

// NumberNode represents integer or real literal numbers.
type NumberNode struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// StringNode represents string literal constants.
type StringNode struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// CharNode represents char literal constants.
type CharNode struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// BoolNode represents boolean literal constants (true, false).
type BoolNode struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// VarAccessToken represents variable/type identifiers.
type VarAccessToken struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

// IfCase represents a single condition-body branch in an if statement.
type IfCase struct {
	Condition interface{} `json:"condition"`
	Body      *ListNode   `json:"body"`
}

// ElseCase represents the optional else branch in an if statement.
type ElseCase struct {
	Body interface{} `json:"body"`
}

// IfNode represents conditional execution (if / else if / else).
type IfNode struct {
	Type     string    `json:"type"`
	Cases    []IfCase  `json:"cases"`
	ElseCase *ElseCase `json:"else_case,omitempty"`
}

// WhileNode represents while-do loops.
type WhileNode struct {
	Type      string      `json:"type"`
	Condition interface{} `json:"condition"`
	Body      *ListNode   `json:"body"`
}

// ForNode represents for loops.
type ForNode struct {
	Type      string      `json:"type"`
	Variable  string      `json:"variable"`
	Start     interface{} `json:"start"`
	End       interface{} `json:"end"`
	Step      interface{} `json:"step,omitempty"`
	Direction string      `json:"direction,omitempty"`
	Body      *ListNode   `json:"body"`
}

// RepeatNode represents repeat-until loops.
type RepeatNode struct {
	Type      string      `json:"type"`
	Condition interface{} `json:"condition"`
	Body      *ListNode   `json:"body"`
}

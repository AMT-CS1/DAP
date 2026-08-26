package parser

/**
2024-11-13: Fix minus unary operation
2024-11-16: Fix multiple arithmetic operations (a+b+c+d...)
            Split INT and REAL, remove NUMBER
	    Split CHAR and CHARRAY
*/

import (
	em "dap/emulator"
	"dap/scanner"
	"log"
	"strconv"
	"strings"
)

type nameattr struct {
	parent string
	typ    string
	val    string
	loc    int
}

var (
	errcount = 0
	parent   = ""
	varcoll  = map[string]nameattr{}
	loc      = 0
	token    *scanner.Token
)

/*
skip until typ is found

	alternate stopping tokens are stop and dead
*/
func sync(typ string, stop string, dead string) {
	for t := token.Next(); t != typ && t != stop && t != dead; t = token.Next() {
		// log.Println(" waiting", typ, "skipping", t)
	}
	token.PushBack()
}

/* skip if typ is found
 */
func skip(typ string) bool {
	// log.Print("skip ", typ)
	matched := token.Next() == typ
	if !matched {
		token.PushBack()
	}
	return matched
}

func expect(tok, msg string) {
	if token.Next() != tok {
		l, c := token.GetLineCol()
		log.Printf("DAP.p %v:%v -- %s", l, c, msg)
		token.PushBack()
		errcount++
	}
}

// expectEnd: cek kata kunci penutup gabungan (mis. "$ENDFOR") ATAU bentuk 2 kata
// "end xxx" yang ditulis terpisah nempel di baris yang sama (mis. "end" + "for").
// secondword itu TIPE token kata kedua (mis. "$FOR"). Kalau dua-duanya gak cocok,
// lapor error dan errcount ditambah (persis kayak expect()).
func expectEnd(oneword, secondword, msg string) {
	if token.Next() == oneword {
		return
	}
	if strings.EqualFold(token.Val, "end") && token.Peek() == secondword && !token.First {
		token.Next() // consumsi kata kedua ("for"/"while"/"if"/"case")
		return
	}
	l, c := token.GetLineCol()
	log.Printf("DAP.p %v:%v -- %s", l, c, msg)
	errcount++
	token.PushBack()
}

/*
variable { , variable }*
*/
func variable_list(gen func(string, string, int)) (namelist []string) {
	// log.Print("variables", token.Val)
	namelist = []string{}
	for {
		if typ := token.Next(); typ == "$NAME" {
			// found a declared name
			namelist = append(namelist, token.Val)
			attr := varcoll[parent+":"+token.Val]
			gen(attr.typ, attr.parent, attr.loc)
			// log.Print(" ", token.Val)
		}
		if token.Next() != "$COMMA" {
			break
		}
	}
	token.PushBack()
	return
}

/*
{ [var] variable_list : type }+
*/
func declaration() {
	// log.Print("declaring ")
	totdecl := 0
	type pendingInit struct {
		typ string
		loc int
		val string
	}
	pending := []pendingInit{}
	for typ := token.Next(); typ == "$DICT" || typ == "$LOCAL" || typ == "$GLOBAL"; typ = token.Next() {
	}
	for {
		// em.GenLine(token.GetLineCol())
		if token.Typ == "$CODE" || token.Typ == "$ENDPROG" {
			break
		}

		if token.Typ == "$CONST" {
			skip("$CONST")
			expect("$NAME", "name expected")
			clabel := token.Val
			// skip("$MEQ")
			expect("$COLON", "':' and type expected after constant name")
			dtyp := token.Next()
			// log.Print("declaring const expr ", token.Typ)
			if dtyp != "$INT" && dtyp != "$REAL" && dtyp != "$CHAR" && dtyp != "$BOOL" && dtyp != "$CHARRAY" {
				l, c := token.GetLineCol()
				log.Printf("DAP.p %v:%v -- invalid or missing type for constant %v", l, c, clabel)
				errcount++
				token.PushBack()
			}
			expect("$ASSG", "'=' expected")
			typ, val := expression()
			// log.Print("declaring done const")
			if typ == "$INT" && dtyp == "$REAL" && val != em.EMPTY {
				// auto-promote literal int ke real, mis. constant PI : real = 3
				if nval, err := strconv.Atoi(val); err == nil {
					val = strconv.FormatFloat(float64(nval), 'E', -1, 64)
					typ = "$REAL"
				}
			}
			if typ != dtyp {
				l, c := token.GetLineCol()
				log.Printf("DAP.p %v:%v -- constant %v declared as %v but value is %v", l, c, clabel, dtyp, typ)
				errcount++
			}
			if _, exists := varcoll[parent+":"+clabel]; exists {
				l, c := token.GetLineCol()
				log.Printf("DAP.p %v:%v -- Variable/constant %v already declared", l, c, clabel)
				errcount++
			} else {
				varcoll[parent+":"+clabel] = nameattr{parent: parent, typ: dtyp, val: val}
			}
			// log.Print(varcoll[parent+":"+clabel], typ, val, clabel)
		} else {
			token.PushBack()
			skip("$VAR")
			namelist := variable_list(em.GenNil)
			totdecl += len(namelist)
			skip("$COLON")
			typ := token.Next()
			if typ == "$INT" || typ == "$REAL" || typ == "$CHAR" || typ == "$BOOL" || typ == "$CHARRAY" {
				/* 20241116 I think now CHAR and CHARRAY are different
				if typ == "$CHAR" {
					typ = "$CHARRAY" // for now, later will split again
				}
				*/
				var locs []int
				for _, v := range namelist {
					if _, exists := varcoll[parent+":"+v]; exists {
						l, c := token.GetLineCol()
						log.Printf("DAP.p %v:%v -- Variable/constant %v already declared", l, c, v)
						errcount++
					} else {
						loc++
						varcoll[parent+":"+v] = nameattr{parent: parent, typ: typ, val: em.EMPTY, loc: loc}
						locs = append(locs, loc)
					}
				}
				// inisialisasi opsional: nama1, nama2 : tipe <- nilai
				if pk := token.Peek(); pk == "$ASSG" {
					token.Next() // consumsi <-/=/:=
					etyp, eval := expression()
					if etyp == "$REAL" && typ == "$INT" && eval != em.EMPTY {
						if fval, err := strconv.ParseFloat(eval, 64); err == nil {
							eval = strconv.Itoa(int(fval))
						}
					} else if etyp == "$INT" && typ == "$REAL" && eval != em.EMPTY {
						if ival, err := strconv.Atoi(eval); err == nil {
							eval = strconv.FormatFloat(float64(ival), 'E', -1, 64)
						}
					} else if etyp != typ {
						l, c := token.GetLineCol()
						log.Printf("DAP.p %v:%v -- Type mismatch in initialization", l, c)
						errcount++
					}
					if eval == em.EMPTY && len(locs) > 1 {
						l, c := token.GetLineCol()
						log.Printf("DAP.p %v:%v -- Multiple variables cannot share a non-constant initializer", l, c)
						errcount++
					} else {
						for _, lc := range locs {
							pending = append(pending, pendingInit{typ: typ, loc: lc, val: eval})
						}
					}
				}
				// log.Print(" ", token.Val)
			} else {
				token.PushBack()
			}
		}
		if typ := token.Peek(); typ == "$CODE" || typ == "$ENDPROG" {
			break
		}
	}
	// log.Println("declared vars: ", totdecl)
	if totdecl > 0 {
		em.GenClaim(totdecl)
	}
	for _, p := range pending {
		em.GenStore(p.typ, parent, p.loc, p.val)
	}
	token.PushBack()
	/*
		for k, v := range varcoll {
			log.Println("declared variables", k, v)
		}
	*/
}

func isStartExpression() bool {
	typ := token.Peek()
	return typ == "$NAME" || typ == "$INT" || typ == "$REAL" || typ == "$CHAR" || typ == "$CHARRAY" || typ == "$TRUE" || typ == "$FALSE" || typ == "$LEFTPAR" || typ == "$MINUS" || typ == "$NOT"
}

/* value | variable | - expr | par_expr
 */
func literal() (string, string) {
	typ := "$INT"
	val := em.EMPTY
	// log.Print("literal ", token.Peek())
	switch token.Next() {
	case "$NAME":
		attr := varcoll[parent+":"+token.Val]
		typ = attr.typ
		if attr == (nameattr{}) {
			l, c := token.GetLineCol()
			log.Printf("DAP.p %v:%v -- Error, variable %v:%v is not defined", l, c, parent, token.Val)
			errcount++
		} else if attr.val == em.EMPTY {
			em.GenCopy(parent, attr.loc)
		} else {
			val = attr.val
		}
		// log.Printf("name=%v type=%v %v val=%v %v", token.Val, attr.typ, typ, attr.val, val)

	case "$INT":
		typ = "$INT"
		val = token.Val
		// em.GenConst( val )
		// log.Print("value=", token.Val)

	case "$REAL":
		typ = "$REAL"
		val = token.Val
		// em.GenReal( val )
		// log.Print("value=", token.Val)

	case "$CHAR":
		typ = "$CHAR"
		if v, err := strconv.Unquote(token.Val + string(token.Val[0])); err != nil {
			log.Print(err)
			val = token.Val[1:]
		} else {
			val = v
			// log.Print(v)
		}
		// log.Print("string=", token.Val)
		// em.GenConst( token.Val )

	case "$CHARRAY":
		typ = "$CHARRAY"
		// log.Print("string=", token.Val)
		em.GenConst(typ, token.Val)
		val = em.EMPTY

	case "$TRUE":
		typ = "$BOOL"
		val = "true"

	case "$FALSE":
		typ = "$BOOL"
		val = "false"
		// em.GenConst( token.Val )
		// log.Print("boolean=", token.Val)

	case "$LEFTPAR":
		typ, val = expression()
		expect("$RIGHTPAR", "missing )")

	case "$MINUS":
		typ, val = literal()
		if typ == "$INT" {
			nval, err := strconv.Atoi(val)
			if err != nil {
				val = em.EMPTY
				em.GenOpCmd("$NEG")
			} else {
				val = strconv.Itoa(-nval)
			}
		} else if typ == "$REAL" {
			nval, err := strconv.ParseFloat(val, 64)
			if err != nil {
				val = em.EMPTY
				em.GenOpCmd("$RNEG")
			} else {
				val = strconv.FormatFloat(-nval, 'E', -1, 64)
			}
		} else {
			l, c := token.GetLineCol()
			log.Printf("DAP.p  %v:%v -- Negation on a non numeric value", l, c)
			errcount++
			typ = "$INT" // by default to avoid further errors
		}

	case "$NOT":
		typ, val = literal()
		if typ != "$BOOL" {
			l, c := token.GetLineCol()
			log.Printf("DAP.p %v:%v -- NOT operation on a non boolean value ", l, c, typ)
			errcount++
		} else if val == em.TRUE {
			val = em.FALSE
		} else if val == em.FALSE {
			val = em.TRUE
		} else {
			em.GenOpCmd("$NOT")
		}
		typ = "$BOOL"
	}
	// log.Print("exit literal ", token.Peek())
	return typ, val
}

/*
literal {addop literal}*

	:D relexpr {mulop relexpr}* !!!
*/
func mulexpr() (string, string) {
	// log.Print("mulexpr")
	atyp := ""
	aval := em.EMPTY
	btyp := ""
	bval := em.EMPTY
	typ := ""
	for {
		btyp, bval = literal()
		//log.Print("from relexpr ", atyp, btyp)
		if atyp == "" {
			atyp = btyp
			aval = bval
		} else if atyp == "$INT" && btyp == "$INT" && typ != "$RDIV" {
			anval, erra := strconv.Atoi(aval)
			bnval, errb := strconv.Atoi(bval)
			if erra != nil || errb != nil {
				em.GenOp2Cmd(typ, atyp, aval, btyp, bval) // */%
				aval = em.EMPTY
			} else {
				switch typ {
				case "$MULT":
					aval = strconv.Itoa(anval * bnval)
				case "$DIV":
					if bnval != 0 {
						aval = strconv.Itoa(anval / bnval)
					} else {
						l, c := token.GetLineCol()
						log.Printf("DAP %v:%v -- Division by zero %v", l, c, typ)
						errcount++
						aval = em.EMPTY
					}
				case "$MOD":
					if bnval != 0 {
						aval = strconv.Itoa(anval % bnval)
					} else {
						l, c := token.GetLineCol()
						log.Printf("DAP %v:%v -- Division by zero %v", l, c, typ)
						errcount++
						aval = em.EMPTY
					}
				default:
					l, c := token.GetLineCol()
					log.Printf("DAP %v:%v -- Illegal operands on numbers %v", l, c, typ)
					errcount++
				}
			}
		} else if atyp == "$REAL" || btyp == "$REAL" || (atyp == "$INT" && btyp == "$INT" && typ == "$RDIV") {
			if atyp == "$INT" {
				em.GenPush(atyp, aval)
				if bval == em.EMPTY && aval == em.EMPTY {
					em.GenOpCmd("$SWAP")
				}
				em.GenOpCmd("$ITOR")
				if bval == em.EMPTY {
					em.GenOpCmd("$SWAP")
				}
				// em.GenOpCmd("$SWAP")
				atyp = "$REAL"
				aval = em.EMPTY
			}
			if btyp == "$INT" {
				em.GenPush(btyp, bval)
				em.GenOpCmd("$ITOR")
				btyp = "$REAL"
				bval = em.EMPTY
			}
			if typ == "$MULT" {
				typ = "$RMULT"
			} else if typ != "$RDIV" {
				l, c := token.GetLineCol()
				log.Printf("DAP %v:%v -- Illegal operation for real type %v", l, c, typ)
				errcount++
			}
			anval, erra := strconv.ParseFloat(aval, 64)
			bnval, errb := strconv.ParseFloat(bval, 64)
			if erra != nil || errb != nil {
				em.GenOp2Cmd(typ, atyp, aval, btyp, bval) // */%
				aval = em.EMPTY
			} else {
				switch typ {
				case "$RMULT":
					aval = strconv.FormatFloat(anval*bnval, 'E', -1, 64)
				case "$RDIV":
					if bnval != 0 {
						aval = strconv.FormatFloat(anval/bnval, 'E', -1, 64)
					} else {
						l, c := token.GetLineCol()
						log.Printf("DAP %v:%v -- Division by zero %v", l, c, typ)
						errcount++
						aval = em.EMPTY
					}
				default:
					l, c := token.GetLineCol()
					log.Printf("DAP %v:%v -- Illegal operands on numbers %v", l, c, typ)
					errcount++
				}
			}
		} else if atyp != btyp {
			l, c := token.GetLineCol()
			log.Printf("DAP.p %v:%v -- Mismatch mul-expression %v vs. %v", l, c, atyp, btyp)
			errcount++
		} else if atyp == "$CHARRAY" || atyp == "CHAR" { // what operations for string and characters?
			l, c := token.GetLineCol()
			log.Printf("DAP.p %v:%v -- Illegal operation %v on %v", l, c, typ, btyp)
			errcount++
		}
		if typ = token.Peek(); typ != "$MULT" && typ != "$DIV" && typ != "$MOD" && typ != "$RMULT" && typ != "$RDIV" {
			break
		}
		token.Next()
	}
	// log.Print("exit mulexpr ", typ)
	return atyp, aval
}

/*
addexpr {mulop addexpr}*

	:D mulexpr {addop mulexpr}* !!!
*/
func addexpr() (string, string) {
	// log.Print("add expression")
	atyp := ""
	aval := em.EMPTY
	btyp := ""
	bval := em.EMPTY
	typ := ""
	for {
		btyp, bval = mulexpr()
		// log.Print("from addexpr ", atyp, btyp)
		if atyp == "" {
			atyp = btyp
			aval = bval
		} else if (atyp == "$INT" || atyp == "$CHAR") && (btyp == "$INT" || btyp == "$CHAR") {
			if typ != "$PLUS" && typ != "$MINUS" && typ != "$SLEFT" && typ != "$SRITE" {
				l, c := token.GetLineCol()
				log.Printf("DAP.p %v:%v -- Illegal operators on numbers %v", l, c, typ)
				errcount++
			} else if aval == em.EMPTY && bval == em.EMPTY {
				em.GenOp2Cmd(typ, atyp, aval, btyp, bval) // +/-
				aval = em.EMPTY
			} else {
				// log.Printf("DAP.p -- doing int ops %v", typ)
				anval, erra := strconv.Atoi(aval)
				bnval, errb := strconv.Atoi(bval)
				if erra != nil || errb != nil {
					em.GenOp2Cmd(typ, atyp, aval, btyp, bval) // +/-
					aval = em.EMPTY
				} else if typ == "$SLEFT" {
					aval = strconv.Itoa(anval << bnval)
				} else if typ == "$SRITE" {
					aval = strconv.Itoa(anval >> bnval)
				} else if typ == "$BAND" {
					aval = strconv.Itoa(anval & bnval)
				} else if typ == "$BXOR" {
					aval = strconv.Itoa(anval ^ bnval)
				} else if typ == "$BOR" {
					aval = strconv.Itoa(anval | bnval)
				} else if typ == "$PLUS" {
					aval = strconv.Itoa(anval + bnval)
				} else { // $MINUS
					aval = strconv.Itoa(anval - bnval)
				}
			}
		} else if atyp == "$REAL" || btyp == "$REAL" {
			if atyp == "$INT" {
				em.GenPush(atyp, aval)
				if bval == em.EMPTY && aval == em.EMPTY { // reversed
					em.GenOpCmd("$SWAP")
				}
				em.GenOpCmd("$ITOR")
				if bval == em.EMPTY {
					em.GenOpCmd("$SWAP")
				}
				// em.GenOpCmd("$SWAP")
				atyp = "$REAL"
				aval = em.EMPTY
			} else if btyp == "$INT" {
				em.GenPush(btyp, bval)
				em.GenOpCmd("$ITOR")
				btyp = "$REAL"
				bval = em.EMPTY
			}
			if typ != "$PLUS" && typ != "$MINUS" {
				l, c := token.GetLineCol()
				log.Printf("DAP.p %v:%v -- Illegal operators on numbers %v", l, c, typ)
				errcount++
			} else if aval == em.EMPTY && bval == em.EMPTY {
				if typ == "$PLUS" {
					typ = "$RPLUS"
				} else {
					typ = "$RMINUS"
				}
				// log.Printf("DAP.p -- doing real ops %v", typ)
				em.GenOp2Cmd(typ, atyp, aval, btyp, bval) // +/-
				aval = em.EMPTY
			} else {
				anval, erra := strconv.ParseFloat(aval, 64)
				bnval, errb := strconv.ParseFloat(bval, 64)
				if erra != nil || errb != nil {
					if typ == "$MINUS" {
						typ = "$RMINUS"
					} else {
						typ = "$RPLUS"
					}
					em.GenOp2Cmd(typ, atyp, aval, btyp, bval) // +/-
					aval = em.EMPTY
				} else if typ == "$PLUS" {
					aval = strconv.FormatFloat(anval+bnval, 'E', -1, 64)
				} else { // $MINUS
					aval = strconv.FormatFloat(anval-bnval, 'E', -1, 64)
				}
			}
		} else if atyp != btyp {
			l, c := token.GetLineCol()
			log.Printf("DAP.p %v:%v -- Mismatch add-expression %v vs. %v", l, c, atyp, btyp)
			errcount++
		} else if typ == "$CHARRAY" {
			if typ != "$PLUS" { // dunno what to generate yet
				l, c := token.GetLineCol()
				log.Printf("DAP.p %v:%v -- illegal operation %v on %v", l, c, typ, atyp)
				errcount++
			} else {
				// genCode?
				l, c := token.GetLineCol()
				log.Printf("DAP.p %v:%v -- Unimplemented string operation %v", l, c, typ)
				errcount++
			}
		}
		if typ = token.Peek(); typ != "$PLUS" && typ != "$MINUS" && typ != "$SLEFT" && typ != "$SRITE" {
			break
		}
		token.Next()
	}
	// log.Print("exit add expression")
	return atyp, aval
}

/* relexpr: literal [relexpr literal]
 */
func relexpr() (string, string) {
	// log.Print("comparison")
	atyp, aval := addexpr()
	//log.Print("from lit1 ", atyp, aval)
	if typ := token.Peek(); typ == "$LEQ" || typ == "$GT" || typ == "$LT" || typ == "$GEQ" || typ == "$EQ" || typ == "$NEQ" {
		token.Next()
		btyp, bval := addexpr()
		// log.Print("from lit ", atyp, btyp)
		if atyp != btyp { // made as $INT or $REAL
			if atyp == "$INT" && btyp == "$REAL" { // between INT and REAL is allowed
				em.GenPush(atyp, aval)
				if bval == em.EMPTY && aval == em.EMPTY { // both in stack already
					em.GenOpCmd("$SWAP")
				}
				em.GenOpCmd("$ITOR")
				if bval == em.EMPTY {
					em.GenOpCmd("$SWAP")
				}
				atyp = "$REAL"
				aval = em.EMPTY
			} else if atyp == "$REAL" && btyp == "$INT" {
				em.GenPush(btyp, bval)
				em.GenOpCmd("$ITOR")
				btyp = "$REAL"
				bval = em.EMPTY
			} else {
				l, c := token.GetLineCol()
				log.Printf("DAP.p %v:%v -- Mismatch cmp-expression %v vs. %v", l, c, atyp, btyp)
				errcount++
			}
		}
		if atyp == "$REAL" { // change to real operations
			typ = em.Rmap[typ]
		}
		val := em.EMPTY
		if aval != em.EMPTY && bval != em.EMPTY {
			if atyp == "$REAL" {
				anval, erra := strconv.ParseFloat(aval, 64)
				bnval, errb := strconv.ParseFloat(bval, 64)
				cmpval := false
				if erra == nil && errb == nil {
					switch typ {
					case "$RLE":
						cmpval = anval <= bnval
					case "$RLT":
						cmpval = anval < bnval
					case "$RGE":
						cmpval = anval >= bnval
					case "$RGT":
						cmpval = anval > bnval
					case "$REQ":
						cmpval = anval == bnval
					case "$RNEQ":
						cmpval = anval != bnval
					}
					val = strconv.FormatBool(cmpval)
				}
			} else {
				anval, erra := strconv.Atoi(aval)
				bnval, errb := strconv.Atoi(bval)
				cmpval := false
				if erra == nil && errb == nil {
					switch typ {
					case "$LE":
						cmpval = anval <= bnval
					case "$LT":
						cmpval = anval < bnval
					case "$GE":
						cmpval = anval >= bnval
					case "$GT":
						cmpval = anval > bnval
					case "$EQ":
						cmpval = anval == bnval
					case "$NEQ":
						cmpval = anval != bnval
					}
					val = strconv.FormatBool(cmpval)
				}
			}
		} else {
			em.GenOp2Cmd(typ, atyp, aval, btyp, bval)
		}
		// log.Print("exit rel comparison")
		return "$BOOL", val
	}
	// log.Print("exit comparison")
	return atyp, aval
}

/*
boolean expression

	relexpr op relexpr (only, not multiple)
*/
func boolexpr() (string, string) {
	// log.Print("boolean expression")
	atyp, aval := relexpr()
	// log.Print("in boolexpr ", token.Peek(), atyp)
	for {
		typ := token.Peek()
		if typ != "$AND" && typ != "$OR" && typ != "$BAND" && typ != "$BOR" && typ != "$BXOR" {
			break
		}
		token.Next()

		if atyp == "$BOOL" && (typ == "$AND" || typ == "$OR") {
			// -- short-circuit evaluation: operand kanan BELUM di-parse di titik ini --
			// emit percabangan berdasarkan operand kiri (A) DULU, baru parsing operand kanan (B).
			// Efeknya: instruksi buat B ada di bytecode SETELAH lompatan, jadi kalau di-skip
			// saat runtime, instruksi B (misal ada "mod 0") gak akan pernah dieksekusi.
			lskip := em.GenLabel()
			lend := em.GenLabel()
			if typ == "$AND" {
				em.GenCond(lskip, aval) // A FALSE -> lompat ke lskip, gak usah evaluasi B
			} else { // $OR
				em.GenCondTrue(lskip, aval) // A TRUE -> lompat ke lskip, gak usah evaluasi B
			}
			btyp, bval := relexpr()
			if btyp != "$BOOL" {
				l, c := token.GetLineCol()
				log.Printf("DAP.p %v:%v -- Illegal operands on boolean %v", l, c, typ)
				errcount++
			}
			em.GenPush(btyp, bval) // kalau B konstanta, belum ke-push, push manual di sini
			em.GenGoto(lend)
			em.GenLoc(lskip)
			if typ == "$AND" {
				em.GenPush("$BOOL", em.FALSE)
			} else {
				em.GenPush("$BOOL", em.TRUE)
			}
			em.GenLoc(lend)
			atyp = "$BOOL"
			aval = em.EMPTY
			continue
		}

		btyp, bval := relexpr()
		// log.Print("from boolexpr ", atyp, btyp)
		// log.Print("bool expr main [", typ, "] ", atyp, " ", aval, " ", btyp, " ", bval)
		if atyp == "$INT" && btyp == "$INT" {
			if typ != "$BAND" && typ != "$BOR" && typ != "$BXOR" {
				l, c := token.GetLineCol()
				log.Printf("DAP.p %v:%v -- Illegal bit operators on numbers %v", l, c, typ)
				errcount++
			} else if aval == em.EMPTY && bval == em.EMPTY {
				em.GenOp2Cmd(typ, atyp, aval, btyp, bval) // +/-
				aval = em.EMPTY
			} else {
				// log.Printf("DAP.p -- doing bit ops %v", typ)
				anval, erra := strconv.Atoi(aval)
				bnval, errb := strconv.Atoi(bval)
				if erra != nil || errb != nil {
					em.GenOp2Cmd(typ, atyp, aval, btyp, bval) // bit oper
					aval = em.EMPTY
				} else if typ == "$BAND" {
					aval = strconv.Itoa(anval & bnval)
				} else if typ == "$BXOR" {
					aval = strconv.Itoa(anval ^ bnval)
				} else { // if typ == "$BOR"
					aval = strconv.Itoa(anval | bnval)
				}
			}
		} else if atyp == "$BOOL" && btyp == "$BOOL" {
			// log.Print("bool expr [", typ, "] ", atyp, " ", aval, " ", btyp, " ", bval)
			if typ != "$AND" && typ != "$OR" {
				l, c := token.GetLineCol()
				log.Printf("DAP.p %v:%v -- Illegal operands on boolean %v", l, c, typ)
				errcount++
			} else if aval != em.EMPTY && bval != em.EMPTY {
				if aval == bval {
				} else if typ == "$OR" {
					aval = em.TRUE
				} else {
					aval = em.FALSE
				}
			} else {
				em.GenOp2Cmd(typ, atyp, aval, btyp, bval)
				aval = em.EMPTY
			}
		}
		// log.Print("bool expr mid [", typ, "] ", atyp, " ", aval, " ", btyp, " ", bval)
	}
	// log.Print("exit boolean expression")
	return atyp, aval
}

func expression() (string, string) {
	// log.Print("expression")
	return boolexpr()
}

/* expr {, expr}*
 */
func expression_list(gen func(string, string)) int {
	// log.Print("expression list")
	count := 0
	for {
		etyp, eval := expression()
		count++
		gen(etyp, eval)
		if typ := token.Peek(); typ != "$COMMA" {
			// log.Print("expression list ", typ)
			break
		}
		token.Next()
	}
	// log.Print("exit expression list")
	return count
}

/* variable <- expr
 */
func assignment(lvl int) (string, int) {
	token.Next()
	// log.Print("assignment ", token.Val)
	assgline := token.GetLine() // this assignment line number
	attr := varcoll[parent+":"+token.Val]
	vtyp := attr.typ
	if attr.loc <= 0 {
		l, c := token.GetLineCol()
		log.Printf("DAP.p %v:%v -- Not a variable %v", l, c, token.Val)
		errcount++
	}
	// keep token, then...
	expect("$ASSG", "<- expected")
	etyp, eval := expression()
	if vtyp == "$REAL" && etyp == "$INT" {
		em.GenPush(etyp, eval)
		em.GenOpCmd("$ITOR")
		etyp = "$REAL"
		eval = em.EMPTY
	} else if vtyp == "$INT" && etyp == "$REAL" {
		em.GenPush(etyp, eval)
		em.GenOpCmd("$RTOI")
		etyp = "$INT"
		eval = em.EMPTY
	} else if vtyp != etyp {
		l, c := token.GetLineCol()
		log.Printf("DAP.p %v:%v -- Type mismatch in assignment", l, c)
		errcount++
	}
	em.GenStore(etyp, attr.parent, attr.loc, eval)
	em.CollectAssg(attr.parent, assgline, attr.loc)
	return attr.parent, attr.loc
}

/* input variable_list
 */
func input_stmt(lvl int) {
	//log.Print("input stmt")
	token.Next()
	has_lpar := token.Peek() == "$LEFTPAR"
	if has_lpar {
		token.Next()
	}
	inpline := token.GetLine() // this input statement line number
	namelist := variable_list(em.GenInp)
	for _, v := range namelist {
		attr := varcoll[parent+":"+v]
		if attr.loc <= 0 {
			l, c := token.GetLineCol()
			log.Printf("DAP.p %v:%v -- Not a variable %v", l, c, v)
			errcount++
		} else {
			em.CollectAssg(parent, inpline, attr.loc)
		}
	}
	if has_lpar {
		expect("$RIGHTPAR", "right parentheses is missing")
	}
}

/* output expression_list
 */
func output_stmt(lvl int) {
	//log.Print("output stmt")
	token.Next()
	has_lpar := token.Peek() == "$LEFTPAR"
	if has_lpar {
		token.Next()
	}
	expression_list(em.GenOut)
	if has_lpar {
		expect("$RIGHTPAR", "right parentheses is missing")
	}
}

/* for assignment to/downto expression do code_list endfor
 */
func for_stmt(lvl int) {
	// log.Print("for stmt")
	expect("$FOR", "for expected")
	l1 := em.GenLabel()
	l2 := em.GenLabel()
	parent, loc := assignment(lvl)
	incr := token.Peek() != "$DOWNTO"
	token.Next()
	totyp, toval := addexpr()
	if totyp != "$INT" {
		l, c := token.GetLineCol()
		log.Printf("DAP.p %v:%v -- Forloop iterator must be integer", l, c)
		errcount++
	}
	em.GenPush(totyp, toval)
	em.GenLoc(l1)
	em.GenDup()
	em.GenCopy(parent, loc)
	if incr { // target value op current value
		em.GenOpCmd("$GEQ")
	} else {
		em.GenOpCmd("$LEQ")
	}
	em.GenCond(l2, em.EMPTY)
	if skip("$DO") {
		em.GenLine(token.GetLineCol())
	}
	code_block(lvl + 1)
	// if skip("ENDFOR") { // optional
	// 	em.GenLine(token.GetLineCol())
	// }
	expectEnd("$ENDFOR", "$FOR", "endfor expected")
	em.GenLine(token.GetLineCol())
	// log.Print("for stmt inc/dec here ", incr)
	em.GenCopy(parent, loc)
	if incr {
		em.GenConst(totyp, "1")
	} else {
		em.GenConst(totyp, "-1")
	}
	em.GenOpCmd("$PLUS")
	em.GenStore(totyp, parent, loc, em.EMPTY)
	em.GenGoto(l1)
	em.GenLoc(l2)
	em.GenPop()
}

/* while bool_expr do code_list endwhile
 */
func while_stmt(lvl int) {
	//log.Print("while stmt")
	expect("$WHILE", "while expected")
	l1 := em.GenLabel()
	l2 := em.GenLabel()
	em.GenLoc(l1)
	em.GenLine(token.GetLineCol()) //!
	etyp, eval := boolexpr()
	if etyp != "$BOOL" {
		l, c := token.GetLineCol()
		log.Printf("DAP.p %v:%v -- Non boolean condition", l, c)
		errcount++
	} else if eval == em.TRUE {
		l, c := token.GetLineCol()
		log.Printf("DAP.p %v:%v -- Infinite loop!", l, c)
		errcount++
	} else if eval == em.FALSE {
		l, c := token.GetLineCol()
		log.Printf("DAP.p %v:%v -- Loop never entered!", l, c)
		errcount++
	}
	em.GenCond(l2, eval)
	expect("$DO", "do expected") // or skip("$WHILE")
	code_block(lvl + 1)
	// if skip("$ENDWHILE") { // optional
	// 	em.GenLine(token.GetLineCol()) //!
	// }
	expectEnd("$ENDWHILE", "$WHILE", "endwhile expected")
	em.GenLine(token.GetLineCol()) //!
	em.GenGoto(l1)
	em.GenLoc(l2)
}

/* repeat code_list until bool_expr
 */
func repeat_stmt(lvl int) {
	// log.Print("repeat-until stmt")
	expect("$REPEAT", "repeat expected")
	l1 := em.GenLabel()
	em.GenLoc(l1)
	em.GenLine(token.GetLineCol()) //!
	code_block(lvl + 1)
	expect("$UNTIL", "until expected")
	em.GenLine(token.GetLineCol())
	etyp, eval := boolexpr()
	if etyp != "$BOOL" {
		l, c := token.GetLineCol()
		log.Printf("DAP.p %v:%v -- Non boolean condition", l, c)
		errcount++
	} else if eval == em.TRUE {
		l, c := token.GetLineCol()
		log.Printf("DAP.p %v:%v -- Useless repeat-until", l, c)
		errcount++
	} else if eval == em.FALSE {
		l, c := token.GetLineCol()
		log.Printf("DAP.p  %v:%v -- Infinite repeat-until loop", l, c)
		errcount++
	}
	em.GenCond(l1, eval)
}

/*
if bool_expr then code_list

	{elif bool_expr then code_list}* [else code_list] endif
*/
func if_stmt(lvl int) {
	// log.Print("if-then-else stmt")
	expect("$IF", "if expected")
	etyp, eval := boolexpr()
	if etyp != "$BOOL" {
		l, c := token.GetLineCol()
		log.Printf("DAP.p %v:%v -- Non boolean condition", l, c)
		errcount++
	} else if eval == em.TRUE {
		l, c := token.GetLineCol()
		log.Printf("DAP.p %v:%v -- Then part is always executed", l, c)
		errcount++
	} else if eval == em.FALSE {
		l, c := token.GetLineCol()
		log.Printf("DAP.p %v:%v -- Never entering the then part", l, c)
		errcount++
	}
	lels := em.GenLabel()
	lfin := em.GenLabel()
	em.GenCond(lels, eval)
	expect("$THEN", "then expected") // or skip("$THEN")
	code_block(lvl + 1)
	typ := token.Next()
	for typ == "$ELIF" || typ == "$ELSE" {
		em.GenGoto(lfin)
		em.GenLoc(lels)
		// validasi posisi/indentasi kata kunci ("elif" atau "else") ITU SENDIRI dulu,
		// SEBELUM ngintip token berikutnya -- soalnya ngintip (Peek) mindahin "posisi baca
		// sekarang", jadi kalau syncStartBlock dipanggil SESUDAH ngintip, dia bakal salah
		// baca posisi (ngecek posisi token setelahnya, bukan posisi "elif"/"else" itu sendiri).
		syncStartBlock(lvl)
		// "else if" (2 kata, nempel sebaris) dianggap alias "elif" -- tapi kalau "if"-nya
		// ada di baris BARU setelah "else", itu artinya if bersarang biasa, bukan alias.
		isElseIf := typ == "$ELSE" && token.Peek() == "$IF" && !token.First
		if isElseIf {
			token.Next() // konsumsi "if" yang nempel di belakang "else"
		}
		if typ == "$ELIF" || isElseIf {
			lels = em.GenLabel()
			em.GenLine(token.GetLineCol()) //!
			etyp, eval = boolexpr()
			if etyp != "$BOOL" {
				l, c := token.GetLineCol()
				log.Printf("DAP.p %v:%v -- Non boolean condition", l, c)
				errcount++
			} else if eval == em.TRUE {
				l, c := token.GetLineCol()
				log.Printf("DAP.p %v:%v -- Else if part is always executed", l, c)
				errcount++
			} else if eval == em.FALSE {
				l, c := token.GetLineCol()
				log.Printf("DAP.p %v:%v -- This else if part is never entered", l, c)
				errcount++
			}
			em.GenCond(lels, eval)
			expect("$THEN", "then of elif expected") // or skip("$THEN")
			code_block(lvl + 1)
			typ = token.Next()
			continue
		}
		// else beneran (bukan alias elseif)
		em.GenLine(token.GetLineCol()) //!
		code_block(lvl + 1)
		expectEnd("$ENDIF", "$IF", "endif expected")
		typ = "$DONE"
		break
	}
	if typ != "$DONE" {
		em.GenLoc(lels)
		if typ == "$ENDIF" {
			// cocok
		} else if strings.EqualFold(token.Val, "end") && token.Peek() == "$IF" && !token.First {
			token.Next() // consumsi "if"
		} else {
			l, c := token.GetLineCol()
			log.Printf("DAP.p %v:%v -- endif expected", l, c)
			errcount++
		}
	}
	em.GenLoc(lfin)
	em.GenLine(token.GetLineCol()) //!
	// log.Print("endif", token.Typ, token.Val, token.Cno)
}

/*
case expr of {expr : code_list}* [otherwise code_list] endcase

	labels either "expr : ...", "expr ) ...", or "expr :) ..."
*/
func case_stmt(lvl int) {
	// log.Print("case stmt")
	expect("$CASE", "case expected")
	ctyp, cval := expression()
	if cval != em.EMPTY {
		l, c := token.GetLineCol()
		log.Printf("DAP.p %v:%v -- Useless switch/case, has constant expression", l, c)
		// errcount++
		em.GenConst(ctyp, cval)
	}
	lfin := em.GenLabel()
	expect("$OF", "of expected")
	for isStartExpression() && !(token.Peek() == "$NAME" && strings.EqualFold(token.Val, "end")) {
		em.GenLine(token.GetLineCol()) //!
		em.GenDup()                    // make a copy of case expression, vs. label expression
		ltyp, lval := expression()
		if ltyp != ctyp {
			l, c := token.GetLineCol()
			log.Printf("DAP.p %v.%v -- mismatch case label type", l, c)
			errcount++
		} else if lval == em.EMPTY {
			l, c := token.GetLineCol()
			log.Printf("DAP.p %v:%v -- Should have constant value as label", l, c)
			// errcount++
		}
		lnex := em.GenLabel()
		em.GenCase(lnex, ltyp, lval)

		skip("$COLON")
		skip("$RIGHTPAR")
		code_block(lvl + 1)
		em.GenGoto(lfin)
		em.GenLoc(lnex)
	}
	if token.Typ == "$DEFAULT" { // default block
		em.GenLine(token.GetLineCol()) //!
		skip("$DEFAULT")
		skip("$COLON")
		skip("$RIGHTPAR")
		code_block(lvl + 1)
	}
	em.GenLoc(lfin)
	em.GenPop()
	// if skip("$ENDCASE") { // optional
	// 	em.GenLine(token.GetLineCol()) //!
	// }
	expectEnd("$ENDCASE", "$CASE", "endcase expected")
	em.GenLine(token.GetLineCol()) //!
}

var tab = 0

func syncStartBlock(lvl int) {
	// log.Print("starting new block ", token.Typ)
	if !token.First {
		for !token.First {
			token.Next()
		}
		token.PushBack()
		l, c := token.GetLineCol()
		if errcount == 0 {
			log.Printf("DAP.p %v:%v -- Must be at the beginning of line", l, c)
			errcount++
		}
	} else if tab == 0 && lvl == 1 {
		tab = token.Cno
	} else if lvl*tab != token.Cno {
		l, c := token.GetLineCol()
		log.Printf("DAP.p %v:%v -- Inconsistent block indentation %v vs. %v", l, c, lvl*tab, token.Cno)
		errcount++
	}
}

func checkBlockLevel(sp, lvl int) {
	if tab == 0 && lvl == 1 {
		tab = sp
	} else if lvl*tab != sp {
		l, c := token.GetLineCol()
		log.Printf("DAP.p %v:%v -- Inconsistent block indentation %v vs. %v", l, c, lvl*tab, sp)
		errcount++
	}
}

// tab decreases, at least (can be more) by one tab
func upLevel(sp, lvl int) bool {
	return (lvl-1)*tab >= sp
}

/* code CR { code CR }*
 */
func code_block(lvl int) {
	// log.Print("block of codes ", lvl)
	for {
		typ := token.Peek()
		/*
						if !token.First {
							for !token.First {
								token.Next()
							}
							token.PushBack()
							log.Print("must be at the beginning of line")
			                errcount++
						} else {
							checkBlockLevel(token.Cno, lvl)
						}
		*/
		// log.Print("indent", token.Cno, lvl, token.Typ, token.Val)
		if typ == "$ENDPROG" || upLevel(token.Cno, lvl) {
			break
		}
		syncStartBlock(lvl)
		switch typ {
		case "$NAME":
			em.GenLine(token.GetLineCol())
			assignment(lvl)

		case "$INPUT":
			em.GenLine(token.GetLineCol())
			input_stmt(lvl)

		case "$OUTPUT":
			em.GenLine(token.GetLineCol())
			output_stmt(lvl)

		case "$FOR":
			for_stmt(lvl)

		case "$WHILE":
			while_stmt(lvl)

		case "$REPEAT":
			repeat_stmt(lvl)

		case "$IF":
			em.GenLine(token.GetLineCol())
			if_stmt(lvl)

		case "$CASE":
			em.GenLine(token.GetLineCol())
			case_stmt(lvl)

		default: // remove the impending token
			token.Next()
		}
	}
}

/* block of codes
 */
func algorithm() {
	// log.Print("algorithm")
	code_block(1)
}

func programBlock() {
	sync("$PROGRAM", "$DICT", "$ENDPROG")
	if token.Next() != "$PROGRAM" {
		l, c := token.GetLineCol()
		log.Printf("DAP.p %v:%v -- Missing program header, found %v/%v", l, c, token.Val, token.Typ)
		errcount++
	}
	em.GenLine(token.GetLineCol()) // trace starts from program line
	if token.Next() != "$NAME" {
		l, c := token.GetLineCol()
		log.Printf("DAP.p %v:%v -- Missing program name, found %v/%v", l, c, token.Val, token.Typ)
		errcount++
	} else {
		// log.Println("Parsing", token.Val)
	}
	// sync("$DICT", "$CODE", "$ENDPROG")
	if typ := token.Next(); typ != "$DICT" && typ != "$GLOBAL" && typ != "$LOCAL" {
		l, c := token.GetLineCol()
		log.Printf("DAP.p %v:%v -- Missing data section, found %v/%v", l, c, token.Val, token.Typ)
		errcount++
	}
	declaration()
	sync("$CODE", "$NAME", "$ENDPROG")
	if token.Next() != "$CODE" {
		l, c := token.GetLineCol()
		log.Printf("DAP.p %v:%v -- Missing code section, found %v/%v", l, c, token.Val, token.Typ)
		errcount++
	}
	em.GenLine(token.GetLineCol()) //!
	algorithm()
	em.GenLine(token.GetLineCol())
	em.GenExit()
	sync("$ENDPROG", "$ENDPROG", "$ENDPROG")
	if token.Next() != "$ENDPROG" {
		l, c := token.GetLineCol()
		log.Printf("DAP.p %v:%v -- Missing end program, found %v/%v", l, c, token.Val, token.Typ)
		errcount++
	} else {
		// log.Println("Program ends")
	}
}

func endParse() {
	/*
		for token.IsAvail() {
			log.Printf("%v:%v: %v %v\n", token.Lno, token.Cno, token.Typ, token.Val)
		}
	*/
	if err := token.Err(); err != nil {
		log.Fatalf("DAP.p *** Fatal error %v", err)
	}
	toterr := token.ErrCount() + errcount
	if toterr > 0 {
		log.Fatalf("*** DAP compilation error count %v", toterr)
	} else {
		log.Printf("*** DAP compilation successful")
	}
}

func Compile(t *scanner.Token) {
	log.Printf("*** DAP compiling")
	token = t
	programBlock()
	endParse()
}

func ProcessSymbols() {
	for vname, vattr := range varcoll {
		em.CollectVariable(vattr.parent, strings.TrimPrefix(vname, vattr.parent)[1:], vattr.typ, vattr.val, vattr.loc)
	}
}

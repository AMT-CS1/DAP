# -*- coding: utf-8 -*-
import subprocess, os, sys, json

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
SUITE = os.path.join(ROOT, "test_suite")
EXE = os.path.join(ROOT, "dap.exe" if sys.platform == "win32" or os.name == "nt" else "dap")

os.makedirs(SUITE, exist_ok=True)

TESTS = []  # list of dicts

def T(id, category, code, mode, expected, stdin="", forbid=None, note=""):
    TESTS.append(dict(id=id, category=category, code=code.strip("\n") + "\n",
                       mode=mode, expected=expected, stdin=stdin, forbid=forbid, note=note))

# ================= KATEGORI A =================
T("A1","A","""
program A1
kamus
    const PI : real = 3.14
algoritma
    output(PI)
endprogram
""", "A", "3.1400000")

T("A2","A","""
program A2
kamus
    const PI = 3.14
algoritma
    output(PI)
endprogram
""", "B", "type expected")

T("A3","A","""
program A3
kamus
    const N : integer = 5
algoritma
    output(N)
endprogram
""", "A", "5")

T("A4","A","""
program A4
kamus
    const X : real = 3
algoritma
    output(X)
endprogram
""", "A", "3.0000000")

T("A5","A","""
program A5
kamus
    const OK : boolean = 5
algoritma
    output(OK)
endprogram
""", "B", "declared as")

T("A6","A","""
program A6
kamus
    constant PI : real = 3.14
algoritma
    output(PI)
endprogram
""", "A", "3.1400000")

# ================= KATEGORI B =================
T("B1","B","""
program B1
kamus
    x : integer
algoritma
    x <- 7 div 2
    output(x)
endprogram
""", "A", "3")

T("B2","B","""
program B2
kamus
    y : real
algoritma
    y <- 7 / 2
    output(y)
endprogram
""", "A", "3.5000000")

T("B3","B","""
program B3
kamus
    b : real
    y : real
algoritma
    b <- 7.5
    y <- b div 2
    output(y)
endprogram
""", "B", "Illegal operation for real type")

T("B4","B","""
program B4
kamus
    a : integer
    y : real
algoritma
    a <- 7
    y <- a div 2.5
    output(y)
endprogram
""", "B", "Illegal operation for real type")

T("B5","B","""
program B5
kamus
    y : real
algoritma
    y <- 7.5 / 2
    output(y)
endprogram
""", "A", "3.7500000")

T("B6","B","""
program B6
kamus
    y : real
algoritma
    y <- 7 / 2.5
    output(y)
endprogram
""", "A", "2.8000000")

T("B7","B","""
program B7
kamus
    y : real
algoritma
    y <- 7.5 / 2.5
    output(y)
endprogram
""", "A", "3.0000000")

T("B8","B","""
program B8
kamus
    x : integer
algoritma
    x <- 7 mod 2
    output(x)
endprogram
""", "A", "1")

T("B9","B","""
program B9
kamus
    y : real
algoritma
    y <- 7 mod 2.0
    output(y)
endprogram
""", "B", "Illegal operation for real type")

# ================= KATEGORI C =================
T("C1","C","""
program C1
kamus
    x : integer
algoritma
    x <- 3 / 4
    output(x)
endprogram
""", "A", "0")

T("C2","C","""
program C2
kamus
    y : real
algoritma
    y <- 3 div 4
    output(y)
endprogram
""", "A", "0.0000000")

# ================= KATEGORI D =================
T("D1","D","""
program D1
kamus
    n : integer
    hasil : boolean
algoritma
    input(n)
    hasil <- (n != 0) and (10 div n == 1)
    output(hasil)
endprogram
""", "A", "false", stdin="0")

T("D2","D","""
program D2
kamus
    n : integer
    hasil : boolean
algoritma
    input(n)
    hasil <- (n == 0) or (10 div n == 1)
    output(hasil)
endprogram
""", "A", "true", stdin="0")

T("D3","D","""
program D3
kamus
    n : integer
    hasil : boolean
algoritma
    input(n)
    hasil <- (n != 0) and (10 div n == 1)
    output(hasil)
endprogram
""", "A", "false", stdin="5")

# ================= KATEGORI E =================
T("E1","E","""
program E1
kamus
    d1, d2, d3, d4 : integer
    hasil : boolean
algoritma
    input(d1, d2, d3, d4)
    hasil <- (d1<d2) and (d2<d3) and (d3<d4)
    output(hasil)
endprogram
""", "A", "true", stdin="1 2 3 4")

T("E2","E","""
program E2
kamus
    d1, d2, d3, d4 : integer
    hasil : boolean
algoritma
    input(d1, d2, d3, d4)
    hasil <- (d1<d2) and (d2<d3) and (d3<d4)
    output(hasil)
endprogram
""", "A", "false", stdin="5 2 3 4")

T("E3","E","""
program E3
kamus
    n : integer
    hasil : boolean
algoritma
    input(n)
    hasil <- (n==1) or (n==2) or (n==3) or (n==4) or (n==5)
    output(hasil)
endprogram
""", "A", "true", stdin="4")

T("E4","E","""
program E4
kamus
    a, b, c, hasil : integer
algoritma
    a <- 15
    b <- 9
    c <- 6
    hasil <- a & b & c
    output(hasil)
endprogram
""", "A", "0")

# ================= KATEGORI F =================
T("F1","F","""
program F1
kamus
    x : integer
algoritma
    x <- 0
    if x>0 then
        output(1)
    else if x<0 then
        output(-1)
    else
        output(0)
    endif
endprogram
""", "A", "0")

T("F2","F","""
program F2
kamus
    x : integer
algoritma
    x <- 5
    if x>0 then
        output(1)
    else if x<0 then
        output(-1)
    else
        output(0)
    endif
endprogram
""", "A", "1")

T("F3","F","""
program F3
kamus
    x : integer
algoritma
    x <- -5
    if x>0 then
        output(1)
    else if x<0 then
        output(-1)
    else
        output(0)
    endif
endprogram
""", "A", "-1")

T("F4a","F","""
program F4a
kamus
    x : integer
algoritma
    x <- 0
    if x>0 then
        output(1)
    else if x<0 then
        output(-1)
    else
        output(0)
    endif
endprogram
""", "A", "0", note="elseif alias (flat), x=0 -> masuk else terakhir")

T("F4b","F","""
program F4b
kamus
    x : integer
algoritma
    x <- 0
    if x>0 then
        output(1)
    else
        if x<0 then
            output(-1)
        endif
    endif
endprogram
""", "A", "", note="if baru baris baru dalam else = nested if BENERAN (bukan alias), butuh endif sendiri; x=0 tidak match cabang manapun -> tidak ada output sama sekali (beda dgn F4a)")

# ================= KATEGORI G =================
T("G1","G","""
program G1
kamus
    i : integer
algoritma
    for i <- 1 to 3 do
        output(i)
endprogram
""", "B", "endfor expected")

T("G2","G","""
program G2
kamus
    i : integer
algoritma
    i <- 0
    while i < 3 do
        output(i)
        i <- i + 1
endprogram
""", "B", "endwhile expected")

T("G3","G","""
program G3
kamus
    x : integer
algoritma
    x <- 1
    if x > 0 then
        output(x)
endprogram
""", "B", "endif expected")

T("G4","G","""
program G4
kamus
    x : integer
algoritma
    x <- 1
    case x of
        1 : output(100)
        2 : output(200)
endprogram
""", "B", "endcase expected")

T("G5","G","""
program G5
kamus
    i : integer
algoritma
    for i <- 1 to 3 do
        output(i)
    endfor
endprogram
""", "A", "1\n2\n3")

# ================= KATEGORI H =================
T("H1","H","""
program H1
kamus
algoritma
    output(42)
endprogram
""", "A", "42")

# ================= KATEGORI I =================
T("I1","I","""
program I1
kamus
    const pi : real = 3.14
    r : real
    volume : real
algoritma
    input(r)
    volume <- (4.0/3.0) * pi * r * r * r
    output(volume)
endprogram
""", "NOLEAK", "DAP --", stdin="2")

# ================= KATEGORI J =================
T("J1","J","""
program J1
kamus
    n, r : integer
algoritma
    input(n)
    r <- 0
    while r*r <= n do
        r<-r+1
    endwhile
    output(r-1)
endprogram
""", "A", "4", stdin="16")

T("J2","J","""
program J2
kamus
    x : integer
algoritma
    x <- -5
    output(x)
endprogram
""", "A", "-5")

T("J3","J","""
program J3
kamus
    x : integer
algoritma
    x <- -5+3
    output(x)
endprogram
""", "A", "-2")

T("J4","J","""
program J4
kamus
    a, b, x : integer
algoritma
    a <- 10
    b <- 3
    x <- a-b
    output(x)
endprogram
""", "A", "7")

# ================= KATEGORI K =================
T("K1","K","""
Program K1
Kamus
    x : integer
Algoritma
    x <- 5
    output(x)
EndProgram
""", "A", "5")

T("K2","K","""
PROGRAM K2
KAMUS
    X : INTEGER
ALGORITMA
    X <- 5
    OUTPUT(X)
ENDPROGRAM
""", "A", "5")

T("K3","K","""
program K3
kamus
    a, b : integer
    hasil : boolean
algoritma
    a <- 5
    b <- 10
    if a > 0 AND b > 0 then
        hasil <- true
    else
        hasil <- false
    endif
    output(hasil)
endprogram
""", "A", "true")

T("K4","K","""
program K4
kamus
    x : integer
algoritma
    INPUT(x)
    OUTPUT(x)
endprogram
""", "A", "7", stdin="7")

T("K5","K","""
program K5
kamus
    For : integer
algoritma
    For <- 5
    output(For)
endprogram
""", "FAIL", "", note="'For' cocok keyword 'for' case-insensitive -> gagal jadi nama variabel")

T("K6","K","""
program K6
kamus
    jumlahTotal : integer
algoritma
    jumlahTotal <- 100
    output(jumlahTotal)
endprogram
""", "A", "100")

T("K7","K","""
program Program
kamus
    x : integer
algoritma
    x <- 5
    output(x)
endprogram
""", "FAIL", "", note="Nama program bernama 'Program' bentrok dengan keyword program case-insensitive")

T("K8","K","""
Program HitungLuas
kamus
    sisi, luas : integer
algoritma
    sisi <- 5
    luas <- sisi * sisi
    output(luas)
endprogram
""", "A", "25", note="PascalCase header dan nama program")

T("K9","K","""
program K9
kamus
    i : integer
    x : integer
algoritma
    x <- 10
    For i <- 1 To 2 Do
        Output(i)
    EndFor
    If x > 5 Then
        Output(x)
    EndIf
endprogram
""", "A", "1\n2\n10", note="Control flow keywords dalam PascalCase")

T("K10","K","""
program K10
kamus
    a, b : integer
    c : boolean
algoritma
    a <- 10 DIV 3
    b <- 10 Mod 3
    c <- Not False And (a > 0 Or b > 0)
    output(a)
    output(b)
    output(c)
endprogram
""", "A", "3\n1\ntrue", note="Operator aritmatika & logika case-insensitive (UPPERCASE & PascalCase)")

T("K11","K","""
program K11
kamus
    i : Integer
    r : REAL
    b : Boolean
    c : Character
algoritma
    i <- 10
    r <- 2.5
    b <- True
    c <- 'X'
    output(i)
    output(r)
    output(b)
    output(c)
endprogram
""", "A", "10\n2.5000000\ntrue\nX", note="Tipe data primitif case-insensitive (PascalCase & UPPERCASE)")

T("K12","K","""
program K12
kamus
    b1, b2, b3, b4 : boolean
algoritma
    b1 <- True
    b2 <- FALSE
    b3 <- Benar
    b4 <- SALAH
    output(b1)
    output(b2)
    output(b3)
    output(b4)
endprogram
""", "A", "true\nfalse\ntrue\nfalse", note="Literal boolean bahasa Inggris dan Indonesia case-insensitive")

T("K13","K","""
program K13
kamus
    x : integer
algoritma
    x <- 99
    Tulis(x)
endprogram
""", "A", "99", note="Kata kunci I/O bahasa Indonesia case-insensitive (Tulis)")

T("K14","K","""
program K14
kamus
    Program : integer
algoritma
    Program <- 10
    output(Program)
endprogram
""", "FAIL", "", note="Variabel bernama 'Program' bentrok dengan keyword program case-insensitive")

T("K15","K","""
program K15
kamus
    nilaiUjian : integer
algoritma
    NilaiUjian <- 85
    output(NilaiUjian)
endprogram
""", "FAIL", "", note="Nama variabel case-sensitive (NilaiUjian != nilaiUjian)")

T("K16","K","""
program K16
kamus
    const MAX_VAL : integer = 100
algoritma
    output(MAX_VAL)
endprogram
""", "A", "100", note="Identifier konstanta UPPERCASE")

# ================= KATEGORI L =================
T("L1","L","""
program L1
kamus
    i : integer
algoritma
    for i <- 1 to 3 do
        output(i)
    end for
endprogram
""", "A", "1\n2\n3")

T("L2","L","""
program L2
kamus
    i : integer
algoritma
    i <- 0
    while i < 3 do
        output(i)
        i <- i + 1
    end while
endprogram
""", "A", "0\n1\n2")

T("L3","L","""
program L3
kamus
    x : integer
algoritma
    x <- 5
    if x > 0 then
        output(x)
    end if
endprogram
""", "A", "5")

T("L4","L","""
program L4
kamus
    x : integer
algoritma
    x <- -5
    if x > 0 then
        output(1)
    else
        output(-1)
    end if
endprogram
""", "A", "-1")

T("L5","L","""
program L5
kamus
    x : integer
algoritma
    x <- 2
    case x of
    1 :
        output(100)
    2 :
        output(200)
    end case
endprogram
""", "A", "200")

T("L6","L","""
program L6
kamus
    x : integer
algoritma
    x <- -5
    if x > 0 then
        output(1)
    else
        if x < 0 then
            output(-1)
        end
        if
    endif
endprogram
""", "FAIL", "", note="'end' baris sendiri lalu 'if' baris baru -> BUKAN alias (butuh nempel sebaris) -> harus tetap gagal compile")

# ================= KATEGORI M =================
T("M1","M","""
program M1
kamus
    kode : character
algoritma
    kode <- 'A'
    if kode == 'A' then
        output(1)
    else
        output(0)
    endif
endprogram
""", "A", "1", forbid="Inconsistence keywords")

# ================= KATEGORI N =================
T("N1","N","""
program N1
kamus
    total : integer <- 9000
algoritma
    output(total)
endprogram
""", "A", "9000")

T("N2","N","""
program N2
kamus
    a, b, c : integer <- 5
algoritma
    output(a)
    output(b)
    output(c)
endprogram
""", "A", "5\n5\n5")

T("N3","N","""
program N3
kamus
    x : real <- 5
algoritma
    output(x)
endprogram
""", "A", "5.0000000")

T("N4","N","""
program N4
kamus
    x : integer <- 7.9
algoritma
    output(x)
endprogram
""", "A", "7")

T("N5","N","""
program N5
kamus
    flag : boolean <- true
algoritma
    output(flag)
endprogram
""", "A", "true")

T("N6","N","""
program N6
kamus
    x : integer
algoritma
    x <- 42
    output(x)
endprogram
""", "A", "42")

# ================= KATEGORI O =================
T("O1","O","""
program O1
kamus
    x, i : integer
algoritma
    i <- 0
    repeat
        output(x)
        i <- i+1
    until i >= 3
endprogram
""", "B", "Illegal access to uninitialized variable")

T("O2","O","""
program O2
kamus
    x : integer
algoritma
    x <- 5 div 0
    output(x)
endprogram
""", "B", "Division by zero")

T("O3","O","""
program O3
kamus
    x : integer
algoritma
    x <- 7 mod 0
    output(x)
endprogram
""", "B", "Division by zero")

T("O4","O","""
program O4
kamus
    a, b : integer
algoritma
    a <- 5
    b <- 0
    output(a div b)
endprogram
""", "B", "Illegal division by zero")

T("O5","O","""
program O5
kamus
    a, b : integer
algoritma
    a <- 7
    b <- 0
    output(a mod b)
endprogram
""", "B", "Illegal modulo division by zero")

T("O6","O","""
program O6
kamus
    x : real
algoritma
    x <- 5.0 / 0.0
    output(x)
endprogram
""", "B", "Division by zero")

T("O7","O","""
program O7
kamus
    a, b, c, d : integer
algoritma
    a <- 17
    b <- 5
    c <- a div b
    d <- a mod b
    output(c)
    output(d)
endprogram
""", "A", "3\n2")

# ================= KATEGORI P =================
T("P1","P","""
program P1
kamus
    n : integer
    n : real
algoritma
    output(n)
endprogram
""", "B", "already declared")

T("P2","P","""
program P2
kamus
    const N : integer = 5
    N : integer
algoritma
    output(N)
endprogram
""", "B", "already declared")

T("P3","P","""
program P3
kamus
    a : integer
    b : real
algoritma
    a <- 1
    b <- 2.5
    output(a)
    output(b)
endprogram
""", "A", "1\n2.5000000")

# ================= KATEGORI Q =================
T("Q1","Q","""
program Q1
kamus
    x : integer
algoritma
    x <- 2
    case x of
    1 :
        output(100)
    2 :
        output(200)
    3 :
        output(300)
    endcase
endprogram
""", "A", "200", note="label harus SEJAJAR dgn 'case' (col sama), body 1 level lbh dalam -- kalau label & body sejajar, parser salah scan (lihat catatan laporan)")

T("Q2","Q","""
program Q2
kamus
    x : integer
algoritma
    x <- 9
    case x of
    1 :
        output(100)
    2 :
        output(200)
    otherwise :
        output(999)
    endcase
endprogram
""", "A", "999")

T("Q3","Q","""
program Q3
kamus
    k : character
algoritma
    k <- 'b'
    case k of
    'a' :
        output(1)
    'b' :
        output(2)
    'c' :
        output(3)
    endcase
endprogram
""", "A", "2")

# ================= KATEGORI R =================
T("R1","R","""
program R1
kamus
    clo1, clo2, clo3 : integer
algoritma
    clo1 <- 2
    clo2 <- 2
    clo3 <- 2
    while clo1>=0 && clo2>=0 && clo3>=0 do
        output(clo1)
        clo1 <- clo1-1
        clo2 <- clo2-1
        clo3 <- clo3-1
    endwhile
endprogram
""", "A", "2\n1\n0")

T("R2","R","""
program R2
kamus
    a, b, c, d : integer
algoritma
    a <- 1
    b <- 2
    c <- 3
    d <- 4
    while (a != b && (a != c) && (a != d)) do
        output(a)
        a <- a+1
    endwhile
endprogram
""", "A", "1")

# ================= INTERAKSI ANTAR FITUR =================
T("INT1","INTERAKSI","""
Program INT1
Kamus
    const BATAS : integer = 3
    a, b, c : integer
    hasil : boolean
Algoritma
    a <- 1
    b <- 2
    c <- 3
    hasil <- (a < b) AND (b < c) AND (c == BATAS)
    OUTPUT(hasil)
EndProgram
""", "A", "true", note="const + keyword case-insensitive + AND berantai")

T("INT2","INTERAKSI","""
program INT2
kamus
    x : integer <- 10
    x : integer
algoritma
    output(x)
endprogram
""", "B", "already declared", note="inisialisasi + duplikat tetap terdeteksi")

T("INT3","INTERAKSI","""
program INT3
kamus
    nilai : integer
    grade : character
algoritma
    nilai <- 75
    if nilai >= 90 THEN
        grade <- 'A'
    else if nilai >= 80 AND nilai < 90 then
        grade <- 'B'
    else if nilai >= 70 and nilai < 80 then
        grade <- 'C'
    else
        grade <- 'D'
    endif
    output(grade)
endprogram
""", "A", "C", note="else-if alias + AND berantai + keyword case-insensitive dalam if besar")

T("INT4","INTERAKSI","""
program INT4
kamus
    i, j : integer
algoritma
    i <- 0
    while i < 2 do
        for j <- 1 to 3 do
            output(j)
        i <- i+1
    endwhile
endprogram
""", "B", "endfor expected", note="for tanpa endfor di dalam while yg py endwhile; cek lokasi error benar (baris 'i <- i+1', bukan baris endwhile)")

T("INT5","INTERAKSI","""
program INT5
kamus
    const KKM : real = 75.0
    nilai : real
    status : character
algoritma
    input(nilai)
    if nilai >= 90.0 then
        status <- 'A'
    else if nilai >= 80.0 then
        status <- 'B'
    else if nilai >= KKM then
        status <- 'C'
    else
        status <- 'D'
    endif
    output(status)
endprogram
""", "A", "C", stdin="77", note="skenario nyata kalkulator nilai")

T("INT6","INTERAKSI","""
program INT6
kamus
    x, y, z : integer
algoritma
    x <- 1
    y <- 2
    z <- 3
    output(x, y)
    output(z)
endprogram
""", "A", "1 2\n3", note="multi-arg output sebaris dipisah spasi, newline di akhir output statement")

# ================= KATEGORI OUT (MULTI-ARG OUTPUT VALID) =================
T("OUT1","OUT","""
program OUT1
kamus
    a, b, c : integer
algoritma
    a <- 1
    b <- 2
    c <- 3
    output(a, b, c)
endprogram
""", "A", "1 2 3", note="multi-arg integer dipisah spasi sebaris")

T("OUT2","OUT","""
program OUT2
kamus
    n : integer
    r : real
    b : boolean
algoritma
    n <- 100
    r <- 3.14
    b <- true
    output("Total:", n, r, b)
endprogram
""", "A", "Total: 100 3.1400000 true", note="multi-arg tipe campuran (string, int, real, bool)")

T("OUT3","OUT","""
program OUT3
kamus
    a, b, c, d, e : integer
algoritma
    a <- 1
    b <- 2
    c <- 3
    d <- 4
    e <- 5
    output(a, b)
    output(c, d)
    output(e)
endprogram
""", "A", "1 2\n3 4\n5", note="beberapa pemanggilan output berurutan")

T("OUT4","OUT","""
program OUT4
kamus
    x, y : integer
algoritma
    x <- 10
    y <- 20
    tulis(x, y)
    tulis(x + y)
endprogram
""", "A", "10 20\n30", note="alias bahasa Indonesia tulis()")

T("OUT5","OUT","""
program OUT5
kamus
    x, y : integer
algoritma
    x <- 10
    y <- 5
    output(x + y, x * y, x - y)
endprogram
""", "A", "15 50 5", note="ekspresi aritmatika kompleks dalam output")

T("OUT6","OUT","""
program OUT6
kamus
    s1, s2 : string
algoritma
    s1 <- "Halo"
    s2 <- "Dunia"
    output(s1, s2)
endprogram
""", "A", "Halo Dunia", note="multi-arg string")

# ================= KATEGORI OUT_ERR (KASUS ERROR & BRUTEFORCE SINTAKS) =================
T("OUT_ERR1","OUT_ERR","""
program OUT_ERR1
kamus
    x : integer
algoritma
    x <- 1
    output(x, )
endprogram
""", "FAIL", "", note="trailing comma tanpa ekspresi")

T("OUT_ERR2","OUT_ERR","""
program OUT_ERR2
kamus
    x, y : integer
algoritma
    x <- 1
    y <- 2
    output(x, y
endprogram
""", "B", "right parentheses is missing", note="kurung tutup hilang")

T("OUT_ERR3","OUT_ERR","""
program OUT_ERR3
kamus
    x, y : integer
algoritma
    x <- 1
    y <- 2
    output x, y)
endprogram
""", "FAIL", "", note="kurung buka hilang tapi ada kurung tutup")

T("OUT_ERR4","OUT_ERR","""
program OUT_ERR4
kamus
    x, y : integer
algoritma
    x <- 1
    y <- 2
    output(x; y)
endprogram
""", "FAIL", "", note="pemisah titik-koma bukan koma")

T("OUT_ERR5","OUT_ERR","""
program OUT_ERR5
kamus
    x, y : integer
algoritma
    x <- 1
    y <- 2
    output(x: y)
endprogram
""", "FAIL", "", note="karakter titik dua dalam argumen output")

T("OUT_ERR6","OUT_ERR","""
program OUT_ERR6
kamus
    x, y : integer
algoritma
    x <- 1
    y <- 2
    output(x <- y)
endprogram
""", "FAIL", "", note="assignment di dalam pemanggilan output")

T("OUT_ERR7","OUT_ERR","""
program OUT_ERR7
kamus
    x, y : integer
algoritma
    x <- 1
    y <- 2
    output(x,, y)
endprogram
""", "FAIL", "", note="koma ganda antar argumen")

T("OUT_ERR8","OUT_ERR","""
program OUT_ERR8
kamus
    x : integer
algoritma
    output()
endprogram
""", "FAIL", "", note="argumen kosong")

T("OUT_ERR9","OUT_ERR","""
program OUT_ERR9
kamus
    x, y : integer
algoritma
    x <- 1
    y <- 2
    output[x, y]
endprogram
""", "FAIL", "", note="kurung siku bukan kurung bulat")

T("OUT_ERR10","OUT_ERR","""
program OUT_ERR10
kamus
    x : integer
algoritma
    output(x, varGaib)
endprogram
""", "B", "is not defined", note="variabel tidak terdefinisi dalam output")

T("OUT_ERR11","OUT_ERR","""
program OUT_ERR11
kamus
    x : integer
algoritma
    x <- 5
    output(x + )
endprogram
""", "FAIL", "", note="operator menggantung tanpa operan kanan")

T("OUT_ERR12","OUT_ERR","""
program OUT_ERR12
kamus
    x, y : integer
algoritma
    x <- 5
    y <- 2
    output(x * * y)
endprogram
""", "FAIL", "", note="operator ganda tidak valid")

T("OUT_ERR13","OUT_ERR","""
program OUT_ERR13
kamus
    r : real
algoritma
    r <- 3.5
    output(r div 2)
endprogram
""", "B", "Illegal operation for real type", note="operasi div pada real")

T("OUT_ERR14","OUT_ERR","""
program OUT_ERR14
kamus
    a, b : integer
algoritma
    a <- 10
    b <- 0
    output(a div b)
endprogram
""", "B", "Illegal division by zero", note="runtime error pembagian nol dalam output")

T("OUT_ERR15","OUT_ERR","""
program OUT_ERR15
kamus
    a, b : integer
algoritma
    a <- 10
    b <- 0
    output(a mod b)
endprogram
""", "B", "Illegal modulo division by zero", note="runtime error modulo nol dalam output")

T("OUT_ERR16","OUT_ERR","""
program OUT_ERR16
kamus
    x, y : integer
algoritma
    x <- 10
    output(x, y)
endprogram
""", "B", "Illegal access to uninitialized variable", note="runtime error akses variabel belum diinisialisasi")

T("OUT_ERR17","OUT_ERR","""
program OUT_ERR17
kamus
    x : integer
algoritma
    output("halo)
endprogram
""", "FAIL", "", note="string literal tanpa kutip penutup")

T("OUT_ERR18","OUT_ERR","""
program OUT_ERR18
kamus
    x : integer
algoritma
    x <- 5
    output(x and 5)
endprogram
""", "FAIL", "", note="operator boolean dengan operan integer")

# ================= KATEGORI ERR_PROG (STRUKTUR PROGRAM) =================
T("ERR_PROG1","ERR_PROG","""
kamus
    x : integer
algoritma
    x <- 1
endprogram
""", "B", "Missing program header", note="header program hilang")

T("ERR_PROG2","ERR_PROG","""
program
kamus
    x : integer
algoritma
    x <- 1
endprogram
""", "B", "Missing program name", note="nama program hilang")

T("ERR_PROG3","ERR_PROG","""
program P
algoritma
    x <- 1
endprogram
""", "B", "Missing data section", note="bagian kamus/data hilang")

T("ERR_PROG4","ERR_PROG","""
program P
kamus
    x : integer
    x <- 1
endprogram
""", "B", "Missing code section", note="bagian algoritma hilang")

# ================= KATEGORI ERR_DECL (DEKLARASI VARIABEL & KONSTANTA) =================
T("ERR_DECL1","ERR_DECL","""
program ERR_DECL1
kamus
    constant N : integer = 5
algoritma
    N <- 10
endprogram
""", "B", "Not a variable N", note="assignment ulang pada konstanta")

T("ERR_DECL2","ERR_DECL","""
program ERR_DECL2
kamus
    x : float
algoritma
    x <- 1.5
endprogram
""", "FAIL", "", note="tipe data tidak valid / deprecated (float)")

T("ERR_DECL3","ERR_DECL","""
program ERR_DECL3
kamus
    x : integer
    x : real
algoritma
    x <- 1
endprogram
""", "B", "already declared", note="deklarasi variabel duplikat tipe berbeda")

# ================= KATEGORI ERR_INPUT (PERINTAH INPUT) =================
T("ERR_INPUT1","ERR_INPUT","""
program ERR_INPUT1
kamus
    x : integer
algoritma
    input(x + 1)
endprogram
""", "FAIL", "", note="ekspresi matematika sebagai target input")

T("ERR_INPUT2","ERR_INPUT","""
program ERR_INPUT2
kamus
    x : integer
algoritma
    input(x
endprogram
""", "B", "right parentheses is missing", note="kurung tutup input hilang")

T("ERR_INPUT3","ERR_INPUT","""
program ERR_INPUT3
kamus
    x : integer
algoritma
    input(varGaib)
endprogram
""", "B", "Not a variable varGaib", note="target input variabel tidak terdaftar")

# ================= KATEGORI ERR_ASSG (PENUGASAN / ASSIGNMENT) =================
T("ERR_ASSG1","ERR_ASSG","""
program ERR_ASSG1
kamus
    x : integer
algoritma
    x <- "teks"
endprogram
""", "B", "Type mismatch in assignment", note="assign string ke integer")

T("ERR_ASSG2","ERR_ASSG","""
program ERR_ASSG2
kamus
    x : integer
algoritma
    x <- true
endprogram
""", "B", "Type mismatch in assignment", note="assign boolean ke integer")

T("ERR_ASSG3","ERR_ASSG","""
program ERR_ASSG3
kamus
    b : boolean
algoritma
    b <- 10
endprogram
""", "B", "Type mismatch in assignment", note="assign integer ke boolean")

T("ERR_ASSG4","ERR_ASSG","""
program ERR_ASSG4
kamus
    x, y : integer
algoritma
    x <- 1
    y := 2
endprogram
""", "B", "Don't mix assignment styles", note="mencampur gaya assignment <- dan :=")

# ================= KATEGORI ERR_IF (PERCABANGAN IF) =================
T("ERR_IF1","ERR_IF","""
program ERR_IF1
kamus
    x : integer
algoritma
    x <- 10
    if x then
        output(x)
    endif
endprogram
""", "B", "Non boolean condition", note="kondisi if berupa integer")

T("ERR_IF2","ERR_IF","""
program ERR_IF2
kamus
    x : integer
algoritma
    x <- 10
    if x > 0
        output(x)
    endif
endprogram
""", "B", "then expected", note="kata kunci then hilang")

T("ERR_IF3","ERR_IF","""
program ERR_IF3
kamus
    x : integer
algoritma
    x <- 10
    if x > 0 then
        output(x)
endprogram
""", "B", "endif expected", note="kata kunci penutup endif hilang")

T("ERR_IF4","ERR_IF","""
program ERR_IF4
kamus
    x : integer
algoritma
    x <- 10
    if x > 0 then
        output(x)
    else
        output(0)
    else
        output(-1)
    endif
endprogram
""", "B", "endif expected", note="klausa else ganda")

# ================= KATEGORI ERR_WHILE (PERULANGAN WHILE) =================
T("ERR_WHILE1","ERR_WHILE","""
program ERR_WHILE1
kamus
    x : integer
algoritma
    x <- 10
    while x do
        x <- x - 1
    endwhile
endprogram
""", "B", "Non boolean condition", note="kondisi while berupa integer")

T("ERR_WHILE2","ERR_WHILE","""
program ERR_WHILE2
kamus
    x : integer
algoritma
    while true do
        output(1)
    endwhile
endprogram
""", "B", "Infinite loop!", note="kondisi while literal true (infinite loop terdeteksi)")

T("ERR_WHILE3","ERR_WHILE","""
program ERR_WHILE3
kamus
    x : integer
algoritma
    while false do
        output(1)
    endwhile
endprogram
""", "B", "Loop never entered!", note="kondisi while literal false (dead code)")

T("ERR_WHILE4","ERR_WHILE","""
program ERR_WHILE4
kamus
    x : integer
algoritma
    x <- 5
    while x > 0
        x <- x - 1
    endwhile
endprogram
""", "B", "do expected", note="kata kunci do hilang pada while")

T("ERR_WHILE5","ERR_WHILE","""
program ERR_WHILE5
kamus
    x : integer
algoritma
    x <- 5
    while x > 0 do
        x <- x - 1
endprogram
""", "B", "endwhile expected", note="kata kunci penutup endwhile hilang")

# ================= KATEGORI ERR_FOR (PERULANGAN FOR) =================
T("ERR_FOR1","ERR_FOR","""
program ERR_FOR1
kamus
    r : real
algoritma
    for r <- 1.0 to 5.0 do
        output(r)
    endfor
endprogram
""", "B", "Forloop iterator must be integer", note="iterator for bertipe real")

T("ERR_FOR2","ERR_FOR","""
program ERR_FOR2
kamus
    i : integer
algoritma
    for i <- 1 do
        output(i)
    endfor
endprogram
""", "FAIL", "", note="kata kunci to hilang pada for")

T("ERR_FOR3","ERR_FOR","""
program ERR_FOR3
kamus
    i : integer
algoritma
    for i <- 1 to 5 do
        output(i)
endprogram
""", "B", "endfor expected", note="kata kunci penutup endfor hilang")

# ================= KATEGORI ERR_REPEAT (PERULANGAN REPEAT-UNTIL) =================
T("ERR_REPEAT1","ERR_REPEAT","""
program ERR_REPEAT1
kamus
    x : integer
algoritma
    x <- 0
    repeat
        x <- x + 1
    until x
endprogram
""", "B", "Non boolean condition", note="kondisi until berupa integer")

T("ERR_REPEAT2","ERR_REPEAT","""
program ERR_REPEAT2
kamus
    x : integer
algoritma
    repeat
        output(1)
    until true
endprogram
""", "B", "Useless repeat-until", note="kondisi until literal true")

T("ERR_REPEAT3","ERR_REPEAT","""
program ERR_REPEAT3
kamus
    x : integer
algoritma
    repeat
        output(1)
    until false
endprogram
""", "B", "Infinite repeat-until loop", note="kondisi until literal false (infinite loop)")

T("ERR_REPEAT4","ERR_REPEAT","""
program ERR_REPEAT4
kamus
    x : integer
algoritma
    repeat
        output(1)
endprogram
""", "B", "until expected", note="kata kunci until hilang")

# ================= KATEGORI ERR_CASE (PERCABANGAN CASE/SWITCH) =================
T("ERR_CASE1","ERR_CASE","""
program ERR_CASE1
kamus
    x : integer
algoritma
    case 5 of
    1 :
        output(1)
    endcase
endprogram
""", "B", "Useless switch/case, has constant expression", note="ekspresi case berupa konstanta literal")

T("ERR_CASE2","ERR_CASE","""
program ERR_CASE2
kamus
    x : integer
algoritma
    x <- 2
    case x of
    'a' :
        output(1)
    endcase
endprogram
""", "B", "mismatch case label type", note="tipe label case tidak cocok dengan pemilih")

T("ERR_CASE3","ERR_CASE","""
program ERR_CASE3
kamus
    x : integer
algoritma
    x <- 2
    case x of
    1 :
        output(1)
endprogram
""", "B", "endcase expected", note="kata kunci penutup endcase hilang")

# ================= KATEGORI ERR_STYLE (GAYA PENULISAN & INDENTASI) =================
T("ERR_STYLE1","ERR_STYLE","""
program ERR_STYLE1
// comment slash
{ comment curly }
kamus
    x : integer
algoritma
    x <- 1
endprogram
""", "B", "Don't mix", note="mencampur gaya komentar // dan {}")

T("ERR_STYLE2","ERR_STYLE","""
program ERR_STYLE2
kamus
    x : integer
algoritma
    x <- 1
      output(x)
endprogram
""", "B", "Inconsistent block indentation", note="indentasi blok tidak konsisten")

# =====================================================================
# runner
# =====================================================================

def run_one(t):
    fname = os.path.join(SUITE, f"test_{t['id']}.dap")
    with open(fname, "w", encoding="utf-8", newline="\n") as f:
        f.write(t["code"])
    try:
        proc = subprocess.run(
            [EXE, "-run", fname],
            input=t["stdin"], capture_output=True, text=True, timeout=10
        )
        timed_out = False
        out, err, rc = proc.stdout, proc.stderr, proc.returncode
    except subprocess.TimeoutExpired as e:
        timed_out = True
        out = (e.stdout or "") if isinstance(e.stdout, str) else (e.stdout or b"").decode(errors="replace")
        err = (e.stderr or "") if isinstance(e.stderr, str) else (e.stderr or b"").decode(errors="replace")
        rc = None

    result = dict(t)
    result["stdout"] = out
    result["stderr"] = err
    result["rc"] = rc
    result["timed_out"] = timed_out

    if timed_out:
        result["pass"] = False
        result["reason"] = "TIMEOUT (10s)"
        return result

    ok = True
    reasons = []

    mode = t["mode"]
    if mode == "A":
        got = out.strip("\n")
        exp = t["expected"]
        if got != exp:
            ok = False
            reasons.append(f"stdout mismatch: got={got!r} expected={exp!r}")
    elif mode == "B":
        if t["expected"] not in err:
            ok = False
            reasons.append(f"stderr missing substring {t['expected']!r}")
    elif mode == "FAIL":
        if rc == 0:
            ok = False
            reasons.append("expected compile failure (nonzero exit) but succeeded")
    elif mode == "NOLEAK":
        combined = out + err
        if t["expected"] in combined:
            ok = False
            reasons.append(f"forbidden substring {t['expected']!r} found in output")

    if t.get("forbid") and t["forbid"] in err:
        ok = False
        reasons.append(f"forbidden substring {t['forbid']!r} found in stderr")

    result["pass"] = ok
    result["reason"] = "; ".join(reasons) if reasons else "OK"
    return result

def main():
    results = []
    for t in TESTS:
        r = run_one(t)
        results.append(r)
        status = "PASS" if r["pass"] else "FAIL"
        print(f"[{status}] {r['id']} ({r['category']})" + (f" -- {r['reason']}" if not r['pass'] else ""))

    with open(os.path.join(SUITE, "results.json"), "w", encoding="utf-8") as f:
        json.dump(results, f, indent=2, ensure_ascii=False)

    total = len(results)
    passed = sum(1 for r in results if r["pass"])
    print(f"\nTOTAL: {passed}/{total} PASS")

if __name__ == "__main__":
    main()

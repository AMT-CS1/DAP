# -*- coding: utf-8 -*-
import subprocess, os, sys, json

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
SUITE = os.path.join(ROOT, "test_suite")
EXE = os.path.join(ROOT, "dap_test.exe")

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

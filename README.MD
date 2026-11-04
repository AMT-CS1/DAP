# DAP

Compiler & emulator untuk bahasa pseudocode **DAP** (Deskripsi Algoritma Pseudocode), dipakai untuk pengajaran algoritma dasar di mata kuliah **AMT-CS1** (Algoritma & Pemrograman) Telkom University. Ditulis dalam Go, dipakai sebagai autograder di LMS (Moodle + CodeRunner).

---

## Tentang

DAP adalah bahasa pseudocode dengan sintaks ala-Indonesia (`program`, `kamus`, `algoritma`, `endprogram`, `<-` untuk assignment, dsb) yang dirancang supaya mahasiswa bisa fokus belajar logika algoritma tanpa terbebani sintaks bahasa pemrograman "beneran". Compiler ini:

- Melakukan parsing pseudocode DAP dan langsung meng-generate bytecode (tidak ada AST terpisah)
- Menjalankan bytecode tersebut lewat VM stack-based buatan sendiri
- Bisa dipakai lewat CLI (`-run file.dap`) atau diintegrasikan ke CodeRunner/Moodle sebagai autograder

---

## Mulai Cepat

```bash
# Clone & masuk folder
git clone <url-repo-ini>
cd DAP

# Build
go build -o dap .

# Jalankan contoh program
./dap -run alpro/contoh.dap
```

Contoh program DAP paling sederhana:

```
program HelloWorld
kamus
    nama : string
algoritma
    input(nama)
    output("Halo, ", nama)
endprogram
```

Sintaks lengkap bahasa DAP (tipe data, operator, struktur kontrol, aturan indentasi, dsb) ada di [`DOKUMENTASI_BAHASA_DAP.md`](./DOKUMENTASI_BAHASA_DAP.md).

---

## Cara Build

### Build biasa (buat testing lokal)
```bash
go build -o dap .
```

### Build untuk deployment ke server Linux (LMS/CodeRunner)
Compiler ini di-deploy sebagai binary statis di server Jobe/CodeRunner. Cross-compile dari OS apapun (Windows/Mac/Linux):

```bash
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o dap-linux-amd64 .
```

Kalau server butuh transfer lewat teks (bukan upload file binary langsung), encode ke base64:
```bash
base64 -w 0 dap-linux-amd64 > dap-linux-amd64.b64.txt
```

> **Catatan**: jangan commit hasil build (`.exe`, `dap-linux-amd64`, `*.b64.txt`) ke git — itu artefak hasil compile, bukan source. Sudah masuk `.gitignore`.

---

## Struktur Project

```
.
├── dap.go              # Entry point, CLI flags (-run, -o, -console, -animate)
├── scanner/            # Lexer/tokenizer
├── parser/             # Recursive descent parser + code generation
├── emulator/           # Stack-based VM (assemble & execute bytecode)
├── ui/                 # Web UI untuk mode animasi (opsional, butuh gopherjs)
├── alpro/              # Contoh-contoh program .dap
└── test_suite/         # Test suite otomatis (lihat bagian Testing)
```

Penjelasan lebih dalam soal cara kerja tiap bagian (pipeline scanner→parser→emulator, sistem indentasi, cara short-circuit bekerja di level bytecode, dsb) ada di [`ARCHITECTURE.md`](./ARCHITECTURE.md).

---

## Testing

Ada test suite otomatis yang generate puluhan file `.dap` kasus uji, jalanin lewat compiler, dan cocokkan hasilnya:

```bash
go build -o dap_test.exe .    # nama binary HARUS persis ini, dicari otomatis oleh script
python test_suite/run_tests.py
```

Jalankan dari root folder project. Hasil detail (stdout/stderr per test) tersimpan di `test_suite/results.json`.

---

## Dokumentasi Lengkap

| Dokumen | Isi |
|---|---|
| [`DOKUMENTASI_BAHASA_DAP.md`](./DOKUMENTASI_BAHASA_DAP.md) | Sintaks lengkap bahasa DAP + penjelasan mekanisme (indentasi, short-circuit, dll) beserta contoh |
| [`ARCHITECTURE.md`](./ARCHITECTURE.md) | Arsitektur compiler: pipeline, struktur data kunci, cara kerja internal |
| [`CHANGELOG.md`](./CHANGELOG.md) | Riwayat perubahan/revisi dari versi asli |

---

## Keterbatasan yang Diketahui

Beberapa hal yang belum didukung atau perlu diwaspadai saat menulis soal:

- Array (`array [1..n] of tipe`) belum berfungsi penuh
- Fungsi/prosedur dan rekursi belum didukung
- `^` adalah operator XOR bitwise, **bukan** pangkat — hati-hati salah pakai
- Label `case` harus sejajar kolom dengan `case` itu sendiri (bukan diindentasi lebih dalam)

Daftar lengkap ada di bagian akhir [`DOKUMENTASI_BAHASA_DAP.md`](./DOKUMENTASI_BAHASA_DAP.md#9-keterbatasan--bug-yang-diketahui-belum-diperbaiki-per-dokumen-ini-ditulis).

---

## Catatan Pengembangan

- Perubahan pada `parser/parser.go`, `scanner/scanner.go`, atau `emulator/emulator.go` sebaiknya selalu diuji lewat `test_suite/run_tests.py` sebelum di-deploy ke LMS, mengingat sifatnya yang saling terkait erat (perubahan kecil di satu tempat bisa berdampak tak terduga di tempat lain).
- Binary hasil build (`.exe`, `dap-linux-amd64`, `*.b64.txt`) tidak masuk version control — build ulang sesuai kebutuhan dari source.
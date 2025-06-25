## 📄 Deskripsi
Sistem backend Payroll Management yang dibangun menggunakan Golang (Gin Framework) dan PostgreSQL, memungkinkan perusahaan untuk:

- Mencatat kehadiran karyawan
- Mencatat lembur dan reimbursement
- Menjalankan proses payroll
- Menghasilkan slip gaji
- Melihat ringkasan total payroll perusahaan
- Mencatat setiap request yang dilakukan untuk kebutuhan audit

## 🏗️ Teknologi
- Golang (Gin)
- PostgreSQL
- GORM
- JWT Authentication
- Unit & Integration Test (Testify, httptest)
- API JSON

## 📁 Struktur Folder
```
.
├── cmd/                # Entry point aplikasi
├── internal/           # Kode utama
│   ├── handlers/       # Controller HTTP
│   ├── models/         # Struct dan definisi tabel DB
│   ├── services/       # Business logic
│   ├── middlewares/    # JWT dan middleware lainnya
├── pkg/                # Helper (DB, JWT, hash)
├── config/             # Konfigurasi
├── tests/              # Unit & integration tests
├── docs/               # Dokumentasi API
├── go.mod
└── README.md
```

## ⚙️ Instalasi
```
git clone https://github.com/nambelaas/payroll-system-go.git
cd payroll-system-go

go mod tidy
```

## 🚀 Menjalankan Aplikasi
```
go run cmd/main.go
```

## 🧪 Testing
```
go test ./tests/... -v
```

## 📦 Api Endpoint
| Endpoint |  Method  | Akses | Deskripsi|
|:-----|:--------:|:------:|:-----:|
| `/login`   | **Post** | Public |Login(JWT)|
| `/employee/attendance/checkin` |  **Post**  | Employee | Check in absen |
| `/employee/attendance/checkout` |  **Post**  | Employee | Check out absen |
| `/employee/reimbursement` |  **Post**  | Employee | Ajukan reimbursement |
| `/employee/payslip/:payslip_period_id` |  **Get**  | Employee | Lihat slip gaji |
| `/admin/payroll/period` |  **Post**  | Admin | Buat periode payroll |
| `/admin/payroll/run` |  **Post**  | Admin | Jalankan payroll |
| `/admin/payslip/summary/:payroll_period_id` |  **Get**  | Admin | Ringkasan slip semua karyawan |

## 👤 Role & Auth
- Admin: Bisa mengelola payroll period dan menjalankan payroll.
- Employee: Hanya bisa mencatat data sendiri.

_Gunakan JWT token dalam `Authorization: Bearer <token>` untuk akses._

## 📚 Dokumentasi Tambahan
Lihat `docs/api.md` untuk spesifikasi endpoint lebih detail

---
## 🛡️ Lisensi
MIT © 2025 - nambelaas

---
## 📬 Kontak
Jika ada pertanyaan:
📧 salmaniseif@gmail.com
🔗 GitHub: nambelaas
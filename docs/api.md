## 📘 Dokumentasi Endpoint Payroll API

### ✍️ Format Umum

- **Base URL**: `http://localhost:8080`
- **Auth**: Gunakan JWT token di header
  ```
  Authorization: Bearer <token>
  ```

---

## 🔐 Auth

### 🔸 Login

**POST** `/login`

#### Request Body:

```json
{
  "username": "admin1",
  "password": "password123"
}
```

#### Response (200 OK):

```json
{
  "token": "JWT_TOKEN_HERE"
}
```

---

## 👩‍💼 Employee Endpoint

### ✅ Check in

**POST** `/employee/attendance/checkin`

- Check in hanya bisa dilakukan 1x per hari
- Tidak bisa dilakukan pada weekend

#### Header:

```
Authorization: Bearer <employee-token>
```

#### Response (201 Created):

```json
{
  "message": "Check in success",
  "time": "09:00"
}
```

### ❎ Check out
**POST** `/employee/attendance/checkout`

- Check in hanya bisa dilakukan 1x per hari

#### Header:

```
Authorization: Bearer <employee-token>
```

#### Response (201 Created):

```json
{
  "message": "Check out successfully",
  "time": "17:00"
}
```

---

### ⏱️ Lembur

**POST** `/employee/overtime`

#### Request Body:

```json
{
  "date": "2025-06-20",
  "hours": 2
}
```

- Maksimum 3 jam per hari
- Hanya bisa diajukan setelah jam kerja selesai

#### Response (201 Created):

```json
{
  "message": "Overtime submitted"
}
```

---

### 💰 Reimbursement

**POST** `/employee/reimbursement`

#### Request Body:

```json
{
  "amount": 150000,
  "description": "Transport Grab"
}
```

#### Response (201 Created):

```json
{
  "message": "Reimbursement submitted"
}
```

---

### 📄 Generate Slip Gaji

**GET** `/employee/payslip/:payroll_period_id`

#### Response (200 OK):

```json
{
  "prorated_salary": 6000000,
  "attendance_hours": 160,
  "overtime_hours": 4,
  "overtime_pay": 480000,
  "reimbursements": [
    {
      "amount": 150000,
      "description": "Transport"
    }
  ],
  "total_take_home": 6630000
}
```

---

## 🛠️ Admin Endpoint

### 📆 Tambah Periode Gaji

**POST** `/admin/payroll/period`

#### Request Body:

```json
{
  "start_date": "2025-06-01",
  "end_date": "2025-06-30"
}
```

#### Response:

```json
{
  "message": "Payroll period created",
  "id": 1
}
```

---

### ⚙️ Jalankan Payroll

**POST** `/admin/payroll/run`

#### Request Body:

```json
{
  "payroll_period_id": 1
}
```

#### Response:

```json
{
  "message": "Payroll processed successfully"
}
```

---

### 📊 Rekap Payslip Semua Karyawan

**GET** `/admin/payslip/summary/:payroll_period_id`

#### Response:

```json
{
  "employees": [
    {
      "employee_id": 1,
      "prorated_salary": 6000000,
      "overtime_pay": 240000,
      "reimbursement_sum": 100000,
      "total_take_home": 6340000
    },
    {
      "employee_id": 2,
      "prorated_salary": 5500000,
      "overtime_pay": 0,
      "reimbursement_sum": 0,
      "total_take_home": 5500000
    }
  ],
  "total_take_home_all": 11840000
}
```

---

## ⚠️ Error Format

Semua error response memiliki format seperti berikut:

```json
{
  "error": "Deskripsi error di sini"
}
```

---

## 📦 Catatan Tambahan

- Semua tanggal dalam format `YYYY-MM-DD`
- Token JWT berbeda antara admin dan employee (berdasarkan role)
- Overtime dihitung 2x dari rate per jam
- Take-home pay = prorated salary + overtime + reimbursement


## 📖 Dokumentasi Database

## 🏷️ List Table
```
.
├── users
├── employees
├── payroll_periods
├── attendances
├── overtimes
├── reimbursements
├── payslips
├── audit_logs
...
```
---
### 👤 Users
|Field|Primary|Foreign|Tipe|Default|Description|
|:---|:----:|:--:|:----|:---|:----|
|ID|✅||bigint|increment|Id user|
|Username|❎||text|not null|Username|
|Password|❎||text|not null|Hashed password|
|Role|❎||text|not null|Peran user|

### 👥 Employees
|Field|Primary|Foreign|Tipe|Default|Description|
|:---|:----:|:---:|:----|:---|:----|
|ID|✅||bigint|increment|Id karyawan|
|UserId|❎|Users[ID]|bigint|not null|Id user|
|Name|❎||text|not null|Nama karyawan|
|Salary|❎||numeric|not null|Gaji karyawan|

### 👥 Payroll Periods
|Field|Primary|Foreign|Tipe|Default|Description|
|:---|:----:|:---:|:----|:---|:----|
|ID|✅||bigint|increment|Id payroll|
|Start Date|❎||timestamp|not null|Waktu awal payroll|
|End Date|❎||timestamp|not null|Waktu terakhir payroll|
|Is Closed|❎||boolean|false|Flag payroll sudah dijalankan|
|Created At|❎||timestamp|not null|Waktu record terbuat|
|Updated At|❎||timestamp|null|Waktu record terupdate|
|Created By|❎||text|not null|User yang menambahkan record|
|Updated By|❎||text|null|User yang mengubah record|

### 🗓️ Attendances
|Field|Primary|Foreign|Tipe|Default|Description|
|:---|:----:|:---:|:----|:---|:----|
|ID|✅||bigint|increment|Id attendance|
|EmployeeId|❎|Employees[ID]|bigint|not null|Id karyawan|
|Date|❎||timestamp|not null|Tanggal attendance|
|Check In|❎||timestamp|not null|Waktu check in|
|Check Out|❎||timestamp|null|Waktu check out|
|Created At|❎||timestamp|not null|Waktu record terbuat|
|Updated At|❎||timestamp|null|Waktu record terupdate|
|Created By|❎||text|not null|User yang menambahkan record|
|Updated By|❎||text|null|User yang mengubah record|

### ⏳ Overtimes
|Field|Primary|Foreign|Tipe|Default|Description|
|:---|:----:|:---:|:----|:---|:----|
|ID|✅||bigint|increment|Id overtimes|
|EmployeeId|❎|Employees[ID]|bigint|not null|Id karyawan|
|Date|❎||timestamp|not null|Tanggal overtime|
|Start Time|❎||timestamp|not null|Jam awal overtime|
|End Time|❎||timestamp|null|Jam akhir overtime|
|Hours|❎||numeric|null|Total jam overtime|
|Created At|❎||timestamp|not null|Waktu record terbuat|
|Updated At|❎||timestamp|null|Waktu record terupdate|
|Created By|❎||text|not null|User yang menambahkan record|
|Updated By|❎||text|null|User yang mengubah record|

### 💸 Reimbursements
|Field|Primary|Foreign|Tipe|Default|Description|
|:---|:----:|:---:|:----|:---|:----|
|ID|✅||bigint|increment|Id reimbursements|
|EmployeeId|❎|Employees[ID]|bigint|not null|Id karyawan|
|Description|❎||text|not null|Deskripsi reimburse|
|Amount|❎||numeric|not null|Nominal yang direimburse|
|Status|❎||text|pending|Status reimburse|
|Date|❎||timestamp|not null|Tanggal reimburse|
|Created At|❎||timestamp|not null|Waktu record terbuat|
|Updated At|❎||timestamp|null|Waktu record terupdate|
|Created By|❎||text|not null|User yang menambahkan record|
|Updated By|❎||text|null|User yang mengubah record|

### 💰 Payslip
|Field|Primary|Foreign|Tipe|Default|Description|
|:---|:----:|:---:|:----|:---|:----|
|ID|✅||bigint|increment|Id payslips|
|EmployeeId|❎|Employees[ID]|bigint|not null|Id karyawan|
|Payroll Period Id|❎|Payroll Periods[ID]|bigint|not null|Id karyawan|
|Attendance Hours|❎||numeric|not null|Jumlah jam attendance|
|Prorated Salary|❎||numeric|not null|Gaji sesuai jam kerja|
|Overtime Hours|❎||numeric|not null|Total jam overtime|
|Overtime Pay|❎||numeric|not null|Total upah overtime|
|Reimbursement Sum|❎||numeric|not null|Total reimbursement|
|Total Take Home|❎||numeric|not null|Total upah digabung dengan reimburse dan overtime|
|Created At|❎||timestamp|not null|Waktu record terbuat|
|Updated At|❎||timestamp|null|Waktu record terupdate|
|Created By|❎||text|not null|User yang menambahkan record|
|Updated By|❎||text|null|User yang mengubah record|

### ⚙️ Audit Logs
|Field|Primary|Foreign|Tipe|Default|Description|
|:---|:----:|:---:|:----|:---|:----|
|ID|✅||bigint|increment|Id payslips|
|Request Id|❎||text|not null|Id unik untuk setiap request|
|User Id|❎|Users[ID]|bigint|not null|Id user|
|Endpoint|❎||text|not null|Endpoint yang digunakan|
|Status|❎||bigint|not null|Status dari response|
|Result|❎||text|not null|Hasil response|
|IpAddress|❎||text|not null|Lokasi ip|
|Created At|❎||timestamp|not null|Waktu record terbuat|
|Updated At|❎||timestamp|null|Waktu record terupdate|
# 🤖 Sistem Kendali Motor Servo Berbasis Golang (Gin Framework) & PostgreSQL

Dokumentasi lengkap dan mendalam untuk proyek **Sistem Pengendali Motor Servo** menggunakan arsitektur REST API dengan bahasa pemrograman **Go (Golang)**, basis data relasional **PostgreSQL**, dan perangkat keras IoT berbasis mikrokontroler (ESP32 / ESP8266 / Arduino) untuk menggerakkan **Motor Servo (SG90 / MG995 / MG996R)**.

---

## 📑 Daftar Isi

1. [Gambaran Umum Proyek](#-gambaran-umum-proyek)
2. [Arsitektur Sistem & Alur Komunikasi](#-arsitektur-sistem--alur-komunikasi)
3. [Struktur Direktori Proyek](#-struktur-direktori-proyek)
4. [Penjelasan Rinci Kode File per File](#-penjelasan-rinci-kode-file-per-file)
   - [1. `main.go`](#1-maingo)
   - [2. `routers/servo_router.go`](#2-routersservo_routergo)
   - [3. `controllers/servo_controller.go`](#3-controllersservo_controllergo)
   - [4. `repositories/servo_repo.go`](#4-repositoriesservo_repogo)
   - [5. `entities/servo_entity.go`](#5-entitiesservo_entitygo)
   - [6. `database/database.go`](#6-databasedatabasego)
   - [7. `database/sql_migrations/1_initiate.sql`](#7-databasesql_migrations1_initiatesql)
   - [8. `go.mod` dan Dependencies](#8-gomod-dan-dependencies)
5. [Konfigurasi Basis Data PostgreSQL & Migrasi](#-konfigurasi-basis-data-postgresql--migrasi)
6. [Dokumentasi Lengkap REST API](#-dokumentasi-lengkap-rest-api)
7. [Integrasi Perangkat Keras (Hardware Motor Servo & Mikrokontroler)](#-integrasi-perangkat-keras-hardware-motor-servo--mikrokontroler)
   - [Kebutuhan Komponen](#kebutuhan-komponen)
   - [Skema Rangkaian (Pinout & Wiring)](#skema-rangkaian-pinout--wiring)
   - [Kode Contoh Mikrokontroler (ESP32 / Arduino C++)](#kode-contoh-mikrokontroler-esp32--arduino-c)
8. [Temuan Bug & Penyesuaian Kode Penting](#-temuan-bug--penyesuaian-kode-penting)
9. [Panduan Menjalankan Proyek Langkah demi Langkah](#-panduan-menjalankan-proyek-langkah-demi-langkah)

---

## 🌟 Gambaran Umum Proyek

Proyek ini merupakan implementasi sistem **Internet of Things (IoT)** yang menghubungkan dunia digital (aplikasi web/mobile) dengan perangkat fisik di dunia nyata melalui perantara **RESTful API Backend**.

Tujuan utama sistem ini adalah:
1. Menyediakan antarmuka komunikasi (API) berbasis protokol HTTP untuk membaca status dan memerintahkan pergerakan sudut motor servo.
2. Menyimpan data status/kondisi motor servo secara persisten di dalam basis data relasional **PostgreSQL**.
3. Menjalankan mekanisme migrasi skema tabel otomatis menggunakan library `sql-migrate` berbasis file `.sql` yang di-embed langsung ke dalam binary aplikasi Golang.
4. Menyediakan endpoint yang dapat diakses secara berkala (*polling*) atau secara waktu nyata (*event-driven*) oleh mikrokontroler (seperti **ESP32**, **NodeMCU ESP8266**, atau **Arduino**) yang terhubung ke jaringan lokal/WiFi.

---

## 🏗 Arsitektur Sistem & Alur Komunikasi

Pola arsitektur yang digunakan memisahkan tanggung jawab kode (*Separation of Concerns*) menjadi beberapa lapisan (*Layered Architecture*):

```
+-------------------------------------------------------------+
|              Klien (Web, Mobile App, Postman)               |
+-------------------------------------------------------------+
                              |
                     HTTP Request (REST API)
                              v
+-------------------------------------------------------------+
|                 Golang Web Server (Gin Framework)           |
|                                                             |
|   [Routers] --------> [Controllers] --------> [Repositories]|
|  (Rute URL)          (Logika HTTP)          (Query ke DB)   |
+-------------------------------------------------------------+
          |                                       |
     Read / Write                            Read / Write
          |                                       |
          v                                       v
+-----------------------+              +----------------------+
|  Basis Data PostgreSQL |              | Mikrokontroler IoT   |
|   (Tabel: servo)      |              |   (ESP32 / ESP8266)  |
+-----------------------+              +----------------------+
                                                  |
                                             Sinyal PWM
                                                  v
                                       +----------------------+
                                       |      Motor Servo     |
                                       | (Rotasi 0° s/d 180°) |
                                       +----------------------+
```

### Penjelasan Alur Kerja:

1. **Inisialisasi Data**:
   - Administrator atau sistem memanggil endpoint `POST /servo/init-proj`.
   - Data awal dengan `id = 1` dan status sudut `0` dibuat di basis data PostgreSQL.
2. **Pengubahan Sudut / Status Servo**:
   - Pengguna mengirimkan perintah rotasi melalui antarmuka (misalnya mengubah posisi servo menjadi 90 derajat atau mode buka/tutup).
   - Pengguna menembak endpoint `PUT /servo/update/90`.
   - Controller membaca parameter `90`, memanggil Repository untuk mengeksekusi perintah SQL `UPDATE` ke PostgreSQL.
   - Status baru tersimpan aman di tabel basis data.
3. **Eksekusi Fisik oleh Mikrokontroler**:
   - Mikrokontroler (ESP32) yang terhubung ke WiFi secara periodik melakukan HTTP GET request ke endpoint `GET /servo/status`.
   - Server membalas dengan format JSON berisi nilai status servo saat ini.
   - ESP32 mengekstrak nilai sudut tersebut dan mengirimkan sinyal gelombang PWM (*Pulse Width Modulation*) ke kabel sinyal Motor Servo.
   - Motor servo berputar secara fisik menyesuaikan sudut yang ditentukan.

---

## 📁 Struktur Direktori Proyek

```text
testbab7/
│
├── controllers/
│   └── servo_controller.go        # Menangani HTTP request & response
├── database/
│   ├── sql_migrations/
│   │   └── 1_initiate.sql         # Skrip DDL SQL pembuatan tabel servo
│   └── database.go                # Logika eksekusi migrasi database otomatis
├── entities/
│   └── servo_entity.go            # Definisi Struct model data servo
├── repositories/
│   └── servo_repo.go              # Lapisan akses data (Query SQL ke PostgreSQL)
├── routers/
│   └── servo_router.go            # Registrasi routing Gin Engine & mapping URL
├── go.mod                         # Konfigurasi dependensi Go Modules
├── go.sum                         # Hash checksum dependensi Go
├── main.go                        # Titik masuk utama (entry point) aplikasi
└── README.md                      # Dokumentasi komprehensif proyek
```

---

## 🔍 Penjelasan Rinci Kode File per File

### 1. `main.go`

File `main.go` merupakan titik masuk (*entry point*) utama dari seluruh siklus hidup aplikasi server Golang.

```go
package testbab7

import (
	"database/sql"
	"fmt"
	"log"
	"mcs_bab_7/database"
	"mcs_bab_7/routers"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = ""
	dbName   = ""
)

var (
	DB  *sql.DB
	err error
)

func main() {
	var PORT = ":50001"
	psqlInfo := fmt.Sprintf(
		`host=%s port=%d user=%s password=%s dbname=%s sslmode=disable`,
		host, port, user, password, dbName,
	)

	DB, err = sql.Open("posgres", psqlInfo)

	if err != nil {
		log.Fatal("gagal membuka DB", psqlInfo)
	}

	database.DBMigrations(DB)

	defer DB.Close()

	routers.StartServer().Run(PORT)
	fmt.Println("DB berhasil tersambung")
}
```

#### Rincian Fungsi & Komponen:
- **Import Packages**:
  - `database/sql`: Package standar Golang untuk menyediakan antarmuka generik basis data SQL.
  - `_ "github.com/lib/pq"`: Driver database PostgreSQL murni untuk Golang. Penggunaan blank identifier `_` bertujuan untuk menjalankan fungsi `init()` dari package tersebut sehingga driver terdaftar ke `database/sql`.
  - `mcs_bab_7/database`: Modul internal untuk menangani migrasi skema tabel.
  - `mcs_bab_7/routers`: Modul internal untuk konfigurasi rute HTTP Gin.
- **Konstanta Koneksi Database**:
  - `host`, `port`, `user`, `password`, `dbName`: Parameter kredensial koneksi ke server PostgreSQL lokal/remote.
- **`psqlInfo`**:
  - String koneksi (*Data Source Name* / DSN) dengan format standar PostgreSQL: `host=... port=... user=... password=... dbname=... sslmode=disable`. Opsi `sslmode=disable` digunakan agar koneksi lokal tidak mewajibkan enkripsi sertifikat SSL.
- **`sql.Open(...)`**:
  - Membuka pool koneksi ke basis data.
- **`database.DBMigrations(DB)`**:
  - Menjalankan fungsi migrasi otomatis sebelum server HTTP dinyalakan, memastikan tabel yang dibutuhkan sudah terbuat.
- **`defer DB.Close()`**:
  - Menutup koneksi database ketika fungsi `main()` selesai dieksekusi agar tidak terjadi kebocoran memori (*resource leak*).
- **`routers.StartServer().Run(PORT)`**:
  - Menginisialisasi Gin Engine dan mengikat (*bind*) aplikasi ke port jaringan `:50001`.

---

### 2. `routers/servo_router.go`

File ini bertanggung jawab mengatur lalu lintas HTTP dan memetakan alamat URL (endpoint) beserta metode HTTP-nya ke *controller handler* yang sesuai.

```go
package routers

import (
	"mcs_bab_7/controllers"

	"github.com/gin-gonic/gin"
)

func StartServer() *gin.Engine {
	routers := gin.Default()
	routers.POST("/servo/init-proj", controllers.InitProj)
	routers.GET("/servo/status", controllers.GetStatus)
	routers.PUT("/servo/update/:srv_status", controllers.UpdateStatus)
	return routers
}
```

#### Rincian Rute:
- `gin.Default()`: Membuat instans Gin Engine yang sudah dilengkapi dengan middleware bawaan: `Logger` (mencatat log request di terminal) dan `Recovery` (mencegah server crash jika terjadi panic 500).
- `POST /servo/init-proj`: Rute untuk inisialisasi baris data awal servo.
- `GET /servo/status`: Rute untuk mengambil seluruh status servo yang tersimpan di database.
- `PUT /servo/update/:srv_status`: Rute dinamis dengan parameter URL `:srv_status` untuk memperbarui nilai status sudut servo.

---

### 3. `controllers/servo_controller.go`

Controller berfungsi sebagai jembatan antara protokol HTTP dengan logika aplikasi. Controller menerima request dari klien, memvalidasi input/parameter, memanggil repository, dan mengembalikan respon berformat JSON ke klien.

```go
package controllers

import (
	"mcs_bab_7/database"
	"mcs_bab_7/entities"
	"mcs_bab_7/repositories"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func InitProj(c *gin.Context) {
	err := repositories.InitProj(database.DbConnection)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{})
}

func GetStatus(c *gin.Context) {
	var result gin.H
	status, err := repositories.GetStatus(database.DbConnection)

	if err != nil {
		result = gin.H{
			"result": err.Error(),
		}
	} else {
		result = gin.H{
			"result": status,
		}
	}
	c.JSON(http.StatusOK, result)
}

func UpdateStatus(c *gin.Context) {
	var status entities.Servo
	srv_status, _ := strconv.Atoi(c.Param("srv_status"))
	status.StatusServo = srv_status
	err := repositories.UpdateStatus(database.DbConnection, status)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"StatusServo": status.StatusServo})
}
```

#### Rincian Fungsi:
1. `InitProj(c *gin.Context)`:
   - Menghubungi `repositories.InitProj` dengan meneruskan pointer koneksi `database.DbConnection`.
   - Mengirim status HTTP 500 jika terjadi kesalahan query, atau HTTP 200 OK jika berhasil.
2. `GetStatus(c *gin.Context)`:
   - Memanggil `repositories.GetStatus` untuk mengambil daftar servo.
   - Membungkus hasil dalam map `gin.H{"result": status}` dan mengembalikan HTTP 200 JSON.
3. `UpdateStatus(c *gin.Context)`:
   - Membaca parameter path dari URL menggunakan `c.Param("srv_status")`.
   - Mengonversi string tersebut ke tipe integer dengan `strconv.Atoi(...)`.
   - Mengisi struct `entities.Servo` lalu memanggil `repositories.UpdateStatus`.
   - Mengembalikan HTTP 200 dengan nilai `StatusServo` terbaru.

---

### 4. `repositories/servo_repo.go`

Repository adalah lapisan data (*Data Access Layer*) yang langsung mengeksekusi query SQL terhadap database PostgreSQL.

```go
package repositories

import (
	"database/sql"
	"mcs_bab_7/entities"
)

func InitProj(db *sql.DB) (err error) {
	sql := "masukkan ke Servo (id, srv_status) values(1, 0)"
	_, err = db.Query(sql)
	return err
}

func GetStatus(db *sql.DB) (result []entities.Servo, err error) {
	sql := "SELECT * FROM servo"
	rows, err := db.Query(sql)

	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var data entities.Servo
		err = rows.Scan(&data.Id, &data.StatusServo)

		if err != nil {
			return
		}
		result = append(result, data)
	}
	return
}

func UpdateStatus(db *sql.DB, status entities.Servo) (err error) {
	sql := "UPDATE servo SET StatusServo = $1 WHERE id = 1"
	_, err = db.Exec(sql, status.StatusServo)
	return
}
```

#### Rincian Fungsi:
1. `InitProj`:
   - Bertujuan memasukkan data awal row pertama ke tabel.
2. `GetStatus`:
   - Melakukan query `SELECT * FROM servo`.
   - Melakukan iterasi terhadap setiap baris data (`rows.Next()`).
   - Melakukan mapping nilai kolom database ke properti struct menggunakan `rows.Scan(&data.Id, &data.StatusServo)`.
   - Mengumpulkan data ke dalam slice `[]entities.Servo`.
3. `UpdateStatus`:
   - Menggunakan `db.Exec(...)` dengan parameter parameterized query `$1` untuk memperbarui status servo dengan `id = 1`.
   - Penggunaan `$1` mencegah bahaya celah keamanan **SQL Injection**.

---

### 5. `entities/servo_entity.go`

File ini memodelkan struktur entitas data servo dalam bentuk Go Struct.

```go
package entities

type Servo struct {
	Id          int `json:"id"`
	StatusServo int `json:"srv_status"`
}
```

#### Rincian Struct:
- `Id`: ID unik penanda motor servo (misal ID 1 untuk servo utama). Tag `json:"id"` menentukan nama kunci saat diserialisasi ke JSON.
- `StatusServo`: Nilai status atau derajat rotasi motor servo (misalnya 0 s/d 180 derajat). Tag `json:"srv_status"` memastikan field ini muncul sebagai `"srv_status"` di dalam respons JSON API.

---

### 6. `database/database.go`

Menangani proses migrasi basis data otomatis menggunakan fitur `embed` bawaan Go 1.16+ dan library pihak ketiga `sql-migrate`.

```go
package database

import (
	"database/sql"
	"embed"
	"fmt"

	migrate "github.com/rubenv/sql-migrate"
)

//go:embed sql_migrations/*.sql
var dbMigrations embed.FS
var DbConnection *sql.DB

func DBMigrations(dbParam *sql.DB) {
	migrations := &migrate.EmbedFileSystemMigrationSource{
		FileSystem: dbMigrations,
		Root:       "sql_migrations",
	}

	n, errs := migrate.Exec(dbParam, "postgres", migrations, migrate.Up)

	if errs != nil {
		panic(errs)
	}

	DbConnection = dbParam
	fmt.Println("Migrasi berhasil", n, migrations)
}
```

#### Rincian Fungsi:
- `//go:embed sql_migrations/*.sql`: Direktif compiler Golang untuk menyematkan seluruh file `.sql` yang ada di folder `sql_migrations` ke dalam variabel `dbMigrations` bertipe `embed.FS`. Keuntungannya: file migrasi menyatu ke dalam file binary eksekutabel sehingga tidak membutuhkan file SQL eksternal saat dideploy.
- `DbConnection`: Variabel global untuk menyimpan referensi koneksi `*sql.DB` yang dapat diakses oleh package lain (seperti controller).
- `migrate.Exec(...)`: Mengeksekusi migrasi ke atas (`migrate.Up`). Jika tabel belum ada, tabel akan dibuat otomatis ke dalam PostgreSQL.

---

### 7. `database/sql_migrations/1_initiate.sql`

File definisi skema SQL untuk pembuatan tabel awal database:

```sql
-- +migrate Up
-- +Migrate StatementBegin

CREATE TABLE servo (
    id INTERGER PRIMARY KEY,
    status_servo INTERGER
)

-- +Migrate StatementEnd
```

#### Penjelasan:
- Tag `-- +migrate Up`: Menandakan instruksi SQL yang dijalankan saat proses migrasi naik (*apply changes*).
- Tag `-- +Migrate StatementBegin` dan `-- +Migrate StatementEnd`: Pembatas blok statement SQL.
- Statement `CREATE TABLE servo`: Membuat tabel bernama `servo` dengan 2 kolom: `id` (kunci utama) dan `status_servo` (nilai status sudut perputaran).

---

### 8. `go.mod` dan Dependencies

File manifest modul Go (`mcs_bab_7`) yang mendaftarkan dependensi pihak ketiga utama:
- `github.com/gin-gonic/gin`: Web framework berkecepatan tinggi dengan routing fleksibel.
- `github.com/lib/pq`: Driver resmi PostgreSQL untuk `database/sql`.
- `github.com/rubenv/sql-migrate`: Alat bantu migrasi skema database berbasis SQL berversi.

---

## 🗄 Konfigurasi Basis Data PostgreSQL & Migrasi

### 1. Instalasi PostgreSQL
Pastikan layanan PostgreSQL sudah terpasang dan aktif di komputer Anda:
- **Windows**: Menggunakan PostgreSQL Installer resmi (EnterpriseDB) dan pgAdmin 4.
- **Linux (Ubuntu/Debian)**: `sudo apt install postgresql postgresql-contrib`
- **Docker**:
  ```bash
  docker run --name pg-servo -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=rahasia -e POSTGRES_DB=servo_db -p 5432:5432 -d postgres:15
  ```

### 2. Membuat Basis Data Baru
Buka terminal / psql CLI / pgAdmin dan buat database baru:
```sql
CREATE DATABASE servo_db;
```

### 3. Mengatur Konfigurasi di `main.go`
Buka `main.go` dan sesuaikan konstanta berikut dengan kredensial PostgreSQL lokal Anda:
```go
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "password_postgres_anda"   // Masukkan password akun postgres Anda
	dbName   = "servo_db"                // Nama database yang dibuat
)
```

---

## 📡 Dokumentasi Lengkap REST API

Server berjalan pada basis URL: `http://localhost:50001` (atau ganti `localhost` dengan alamat IP LAN server Anda saat dihubungkan ke ESP32).

---

### 1. Inisialisasi Proyek (Buat Baris Data Awal)

Menyisipkan entri servo default ke dalam basis data dengan nilai awal sudut 0.

- **URL**: `/servo/init-proj`
- **Method**: `POST`
- **Headers**: `Content-Type: application/json`
- **Body**: *None*
- **Response**:
  - **Status 200 OK**:
    ```json
    {}
    ```
  - **Status 500 Internal Server Error**:
    ```json
    {
      "error": "pq: duplicate key value violates unique constraint \"servo_pkey\""
    }
    ```
- **Contoh cURL**:
  ```bash
  curl -X POST http://localhost:50001/servo/init-proj
  ```

---

### 2. Membaca Status Servo (Status Polling)

Mengambil daftar seluruh status motor servo yang tersimpan di dalam database. Endpoint ini yang dipanggil secara berulang (*polling*) oleh mikrokontroler.

- **URL**: `/servo/status`
- **Method**: `GET`
- **Headers**: `None`
- **Body**: *None*
- **Response**:
  - **Status 200 OK**:
    ```json
    {
      "result": [
        {
          "id": 1,
          "srv_status": 90
        }
      ]
    }
    ```
- **Contoh cURL**:
  ```bash
  curl -X GET http://localhost:50001/servo/status
  ```

---

### 3. Memperbarui Status / Sudut Servo

Mengubah nilai sudut rotasi atau status servo berdasarkan parameter URL yang dikirim.

- **URL**: `/servo/update/:srv_status`
- **Method**: `PUT`
- **URL Parameter**:
  - `srv_status` *(integer)*: Nilai sudut perputaran baru (misalnya `0`, `45`, `90`, `180`).
- **Body**: *None*
- **Response**:
  - **Status 200 OK**:
    ```json
    {
      "StatusServo": 90
    }
    ```
- **Contoh cURL**:
  ```bash
  # Menggerakkan servo ke sudut 90 derajat
  curl -X PUT http://localhost:50001/servo/update/90

  # Menggerakkan servo ke sudut 180 derajat
  curl -X PUT http://localhost:50001/servo/update/180

  # Mengembalikan servo ke sudut 0 derajat
  curl -X PUT http://localhost:50001/servo/update/0
  ```

---

## ⚙️ Integrasi Perangkat Keras (Hardware Motor Servo & Mikrokontroler)

### Kebutuhan Komponen:
1. **Mikrokontroler**: ESP32 Dev Module atau NodeMCU ESP8266 (memiliki modul WiFi bawaan).
2. **Motor Servo**: Micro Servo SG90 (tegangan kerja 4.8V - 6V, rentang putaran 0° - 180°) atau MG995/MG996R untuk beban torsi lebih besar.
3. **Kabel Jumper**: Male-to-Male atau Male-to-Female sesuai pin header board.
4. **Catu Daya / Power Supply**: Disarankan sumber daya eksternal 5V (1A - 2A) jika menggunakan servo beban berat, atau pin VIN/5V mikrokontroler untuk servo SG90 kecil.
5. **Kapasitor (Opsional)**: Kapasitor elektrolit 100uF - 470uF antara VCC dan GND servo untuk meredam lonjakan arus (*voltage spike/brownout*).

---

### Skema Rangkaian (Pinout & Wiring)

Motor servo standar memiliki 3 buah kabel warna:

| Kabel Motor Servo | Fungsi | Pin pada ESP32 / Mikrokontroler |
| :--- | :--- | :--- |
| **Cokelat / Hitam** | Ground (GND) | Pin `GND` |
| **Merah** | Daya Positif (VCC / 5V) | Pin `5V` / `VIN` (atau VCC Catu Daya Eksternal) |
| **Oranye / Kuning** | Sinyal PWM (Control Signal) | Pin GPIO (Contoh: `GPIO 13` atau `GPIO 18`) |

> [!IMPORTANT]
> Jika menggunakan catu daya eksternal untuk motor servo, **pastikan ground (GND) catu daya eksternal dihubungkan bersama (Common Ground) ke pin GND mikrokontroler**. Kegagalan menghubungkan ground bersama akan menyebabkan sinyal PWM tidak terbaca dengan benar oleh servo.

---

### Kode Contoh Mikrokontroler (ESP32 / Arduino C++)

Berikut adalah kode lengkap untuk diunggah (*flash*) ke mikrokontroler **ESP32** menggunakan Arduino IDE. Program ini akan menghubungkan ESP32 ke WiFi, mengambil data JSON dari endpoint `GET /servo/status`, mengekstrak nilai sudut, dan memutar motor servo secara otomatis.

> **Library yang dibutuhkan di Arduino IDE:**
> 1. `ESP32Servo` (by Kevin Harrington) - install via Library Manager.
> 2. `ArduinoJson` (by Benoit Blanchon versi 6 atau 7) - install via Library Manager.

```cpp
#include <WiFi.h>
#include <HTTPClient.h>
#include <ArduinoJson.h>
#include <ESP32Servo.h>

// 1. Konfigurasi Jaringan WiFi
const char* ssid     = "NAMA_WIFI_ANDA";
const char* password = "PASSWORD_WIFI_ANDA";

// 2. Alamat API Golang (Gunakan IP LAN komputer server Anda)
// Jangan gunakan 'localhost' karena localhost bagi ESP32 adalah dirinya sendiri!
const char* serverUrl = "http://192.168.1.100:50001/servo/status";

// 3. Konfigurasi Servo
Servo myServo;
const int servoPin = 13;   // Pin GPIO ESP32 yang terhubung ke kabel oranye/kuning servo
int lastAngle = -1;        // Variabel penyimpan sudut sebelumnya

void setup() {
  Serial.begin(115200);
  delay(1000);

  // Alokasi timer PWM untuk ESP32
  ESP32PWM::allocateTimer(0);
  ESP32PWM::allocateTimer(1);
  ESP32PWM::allocateTimer(2);
  ESP32PWM::allocateTimer(3);
  
  myServo.setPeriodHertz(50);             // Standar 50Hz untuk motor servo
  myServo.attach(servoPin, 500, 2400);    // Pulsa minimum 500us dan maksimum 2400us

  // Koneksi ke WiFi
  Serial.print("Menghubungkan ke WiFi: ");
  Serial.println(ssid);
  WiFi.begin(ssid, password);

  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.print(".");
  }

  Serial.println("\nWiFi Berhasil Terhubung!");
  Serial.print("IP ESP32: ");
  Serial.println(WiFi.localIP());
}

void loop() {
  // Pastikan koneksi WiFi masih aktif
  if (WiFi.status() == WL_CONNECTED) {
    HTTPClient http;

    Serial.print("Menghubungi Server: ");
    Serial.println(serverUrl);

    http.begin(serverUrl);
    int httpResponseCode = http.GET();

    if (httpResponseCode == HTTP_CODE_OK) { // HTTP 200
      String payload = http.getString();
      Serial.println("Respon dari Server: " + payload);

      // Parsing JSON menggunakan ArduinoJson
      // Format respon: {"result": [{"id": 1, "srv_status": 90}]}
      DynamicJsonDocument doc(1024);
      DeserializationError error = deserializeJson(doc, payload);

      if (!error) {
        // Ambil elemen pertama dari array 'result'
        int targetAngle = doc["result"][0]["srv_status"];

        // Batasi sudut antara 0 hingga 180 derajat demi keamanan motor servo
        targetAngle = constrain(targetAngle, 0, 180);

        // Hanya gerakkan servo jika sudut mengalami perubahan
        if (targetAngle != lastAngle) {
          Serial.printf("Menggerakkan Servo dari %d° ke %d°\n", lastAngle, targetAngle);
          myServo.write(targetAngle);
          lastAngle = targetAngle;
        } else {
          Serial.println("Sudut servo tidak berubah.");
        }
      } else {
        Serial.print("Gagal parsing JSON: ");
        Serial.println(error.c_str());
      }
    } else {
      Serial.print("HTTP GET Gagal, kode error: ");
      Serial.println(httpResponseCode);
    }

    http.end(); // Bebaskan resource HTTP
  } else {
    Serial.println("Koneksi WiFi Terputus!");
  }

  // Interval polling ke server setiap 2 detik (2000 ms)
  delay(2000);
}
```

---

## ⚠️ Temuan Bug & Penyesuaian Kode Penting

Sebelum menjalankan proyek, terdapat beberapa kesalahan pengetikan (*typo*) dan ketidakkonsistenan sintaks pada kode sumber bawaan yang **wajib diperbaiki** agar program dapat dikompilasi dan berjalan tanpa error:

### 1. Koreksi pada `main.go`
- **Nama Package**:
  - *Saat ini*: `package testbab7`
  - *Perbaikan*: Harus diubah menjadi `package main` agar Golang mengenali file ini sebagai binary eksekutabel yang memiliki fungsi `main()`.
- **Nama Driver Database**:
  - *Saat ini*: `sql.Open("posgres", psqlInfo)`
  - *Perbaikan*: Driver PostgreSQL yang valid bernama `"postgres"` (kurang huruf **t**). Ubah menjadi:
    ```go
    DB, err = sql.Open("postgres", psqlInfo)
    ```
- **Kredensial Database**:
  - Variabel `password` dan `dbName` masih bernilai string kosong `""`. Wajib diisi sesuai konfigurasi database PostgreSQL lokal Anda.

### 2. Koreksi pada `database/sql_migrations/1_initiate.sql`
- **Typo Tipe Data**:
  - *Saat ini*: `id INTERGER PRIMARY KEY`, `status_servo INTERGER`
  - *Perbaikan*: PostgreSQL tidak mengenal tipe `INTERGER` (kelebihan huruf **R**). Tipe data yang benar adalah `INTEGER` atau `INT`:
    ```sql
    -- +migrate Up
    -- +Migrate StatementBegin

    CREATE TABLE IF NOT EXISTS servo (
        id INTEGER PRIMARY KEY,
        status_servo INTEGER
    );

    -- +Migrate StatementEnd
    ```

### 3. Koreksi pada `repositories/servo_repo.go`
- **Sintaks Bahasa Indonesia pada SQL Query**:
  - *Saat ini*: `sql := "masukkan ke Servo (id, srv_status) values(1, 0)"`
  - *Perbaikan*: Perintah SQL harus menggunakan sintaks standar SQL dalam bahasa Inggris `INSERT INTO`:
    ```go
    sql := "INSERT INTO servo (id, status_servo) VALUES (1, 0) ON CONFLICT (id) DO NOTHING"
    ```
- **Ketidakkonsistenan Nama Kolom**:
  - Di file migrasi nama kolom adalah `status_servo`.
  - Di fungsi `UpdateStatus`: `UPDATE servo SET StatusServo = $1 WHERE id = 1`.
  - *Perbaikan*: Samakan nama kolom menjadi `status_servo`:
    ```go
    sql := "UPDATE servo SET status_servo = $1 WHERE id = 1"
    ```

---

## 🚀 Panduan Menjalankan Proyek Langkah demi Langkah

Ikuti panduan berikut untuk menjalankan seluruh sistem dari awal hingga siap dihubungkan ke motor servo:

### Langkah 1: Persiapan Lingkungan (Prerequisites)
Pastikan telah menginstal:
- **Go**: versi 1.20 ke atas (`go version`)
- **PostgreSQL**: versi 13 ke atas
- **Arduino IDE**: jika menguji perangkat keras mikrokontroler

### Langkah 2: Buat Database di PostgreSQL
Buka terminal / CMD / PowerShell dan jalankan:
```bash
psql -U postgres
```
Lalu eksekusi perintah SQL:
```sql
CREATE DATABASE servo_db;
\q
```

### Langkah 3: Konfigurasi File Sumber
1. Buka `main.go`, ubah `package testbab7` menjadi `package main`.
2. Perbaiki driver `"posgres"` menjadi `"postgres"`.
3. Masukkan password dan nama database:
   ```go
   password = "password_anda"
   dbName   = "servo_db"
   ```
4. Perbaiki file migrasi `1_initiate.sql` dan repository `servo_repo.go` sesuai bagian [Temuan Bug & Penyesuaian Kode](#-temuan-bug--penyesuaian-kode-penting).

### Langkah 4: Sinkronisasi Dependensi
Jalankan perintah berikut di root folder proyek:
```bash
go mod tidy
```

### Langkah 5: Jalankan Server Golang
Mulai server backend dengan perintah:
```bash
go run main.go
```
Jika berhasil, terminal akan menampilkan log migrasi berhasil dan server Gin mulai mendengarkan (*listening*) pada port `50001`:
```text
Migrasi berhasil 1 ...
[GIN-debug] POST   /servo/init-proj          --> mcs_bab_7/controllers.InitProj (3 handlers)
[GIN-debug] GET    /servo/status             --> mcs_bab_7/controllers.GetStatus (3 handlers)
[GIN-debug] PUT    /servo/update/:srv_status --> mcs_bab_7/controllers.UpdateStatus (3 handlers)
[GIN-debug] Listening and serving HTTP on :50001
```

### Langkah 6: Uji Coba Endpoint Menggunakan cURL / Postman

1. **Inisialisasi Data Default**:
   ```bash
   curl -X POST http://localhost:50001/servo/init-proj
   ```
2. **Cek Status Servo Awal**:
   ```bash
   curl -X GET http://localhost:50001/servo/status
   ```
   *Output*:
   ```json
   {"result":[{"id":1,"srv_status":0}]}
   ```
3. **Ubah Posisi Servo ke Sudut 90 Derajat**:
   ```bash
   curl -X PUT http://localhost:50001/servo/update/90
   ```
   *Output*:
   ```json
   {"StatusServo":90}
   ```
4. **Verifikasi Ulang Status**:
   ```bash
   curl -X GET http://localhost:50001/servo/status
   ```
   *Output*:
   ```json
   {"result":[{"id":1,"srv_status":90}]}
   ```

### Langkah 7: Pengujian Terpadu dengan Motor Servo Fisik
1. Cek alamat IP LAN komputer server (jalankan `ipconfig` di Windows atau `ifconfig` di Linux/macOS, contoh: `192.168.1.50`).
2. Masukkan IP tersebut ke variabel `serverUrl` pada kode Arduino ESP32:
   ```cpp
   const char* serverUrl = "http://192.168.1.50:50001/servo/status";
   ```
3. Unggah (*upload*) kode ke ESP32.
4. Buka **Serial Monitor** pada baudrate `115200`.
5. Ubah sudut servo melalui endpoint `PUT /servo/update/180`.
6. Amati motor servo akan berputar secara otomatis mengikuti instruksi yang tersimpan di PostgreSQL!

---

## 📜 Lisensi & Penutup

Proyek ini dibuat sebagai modul pembelajaran sistem terdistribusi, komunikasi mikrokontroler berbasis IoT (*Internet of Things*), serta integrasi backend modern dengan basis data relasional.


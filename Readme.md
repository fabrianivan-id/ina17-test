# Concert Ticket Booking System

## Stack
- Frontend: React.js + Vite + TailwindCSS
- Backend: Golang (Gin framework)
- Database: MySQL + Redis

## Setup
### Development

#### Backend
```
cd booking-service
go run main.go
```

#### Frontend
```
cd frontend
npm install
npm run dev
```


# 1. SQL & Index Optimization

## Struktur Database

```sql
-- Tabel users (user-service)
CREATE TABLE users (
    id CHAR(36) PRIMARY KEY,         -- UUID
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabel concerts (booking-service)
CREATE TABLE concerts (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    date DATETIME NOT NULL,
    venue VARCHAR(255) NOT NULL,
    capacity INT NOT NULL
);

-- Tabel bookings (booking-service)
CREATE TABLE bookings (
    id CHAR(36) PRIMARY KEY,         -- UUID
    user_id CHAR(36) NOT NULL,
    concert_id INT NOT NULL,
    quantity INT NOT NULL CHECK (quantity > 0),
    status ENUM('pending', 'confirmed', 'cancelled') DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (concert_id) REFERENCES concerts(id)
);

-- Tabel payments (payment-service)
CREATE TABLE payments (
    id CHAR(36) PRIMARY KEY,         -- UUID
    booking_id CHAR(36) NOT NULL UNIQUE,
    amount DECIMAL(10,2) NOT NULL,
    status ENUM('pending', 'success', 'failed') DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (booking_id) REFERENCES bookings(id)
);
```

## Strategi Optimasi Index

- **Composite Index pada Bookings**
    ```sql
    CREATE INDEX idx_bookings_user_concert ON bookings(user_id, concert_id);
    ```
    - Optimasi query pencarian booking berdasarkan user dan konser
    - Mempercepat JOIN antara bookings dan concerts

- **Partial Index untuk Status Pending**
    ```sql
    CREATE INDEX idx_payments_pending ON payments(status) WHERE status = 'pending';
    ```
    - Optimasi query pembayaran tertunda
    - Mengurangi ukuran indeks

- **Covering Index untuk Reports**
    ```sql
    CREATE INDEX idx_concert_date_venue ON concerts(date, venue) INCLUDE (title, capacity);
    ```
    - Memungkinkan query laporan mengambil data langsung dari indeks

- **Hash Index untuk UUID**
    ```sql
    CREATE INDEX idx_users_id_hash ON users USING HASH (id);
    ```
    - Optimasi lookup berdasarkan UUID (equality checks)

---

# 2. High-Load Optimization (100.000 hit/detik)

## Database Connection Pooling

- Konfigurasi pool koneksi (Golang: `SetMaxOpenConns`, `SetMaxIdleConns`)
- Cegah overload DB dengan membatasi koneksi simultan

## Read Replicas

- Arahkan 90% traffic read ke replica database
- Gunakan load balancer untuk distribusi query

## Redis Caching

```go
// Contoh caching dengan Redis
func GetConcertDetails(concertID int) (Concert, error) {
    cacheKey := fmt.Sprintf("concert:%d", concertID)
    if cached, err := redis.Get(cacheKey); err == nil {
        return cached, nil
    }
    // Query database jika tidak ada di cache
    // ...
    redis.SetEx(cacheKey, concert, 5*time.Minute)
}
```
- Cache response API yang sering diakses (TTL 5-10 menit)
- Cache hasil query kompleks

## Rate Limiting

```go
// Middleware rate limiting di Golang
func RateLimiter(limit int) gin.HandlerFunc {
    store := redis.NewRateStore()
    return tollbooth_gin.LimitHandler(
        tollbooth.NewLimiter(limit, nil),
    )
}
```
- Batasi 100 request/detik per IP
- Gunakan algoritma Token Bucket

## Asynchronous Processing

- Gunakan message queue (RabbitMQ/Kafka) untuk:
    - Proses pembayaran
    - Update inventory
    - Kirim email konfirmasi
- Langsung kembalikan response "processing" ke client

## Database Sharding

- Partisi data berdasarkan wilayah geografis (`shard_key = user_location`)
- Distribusi beban ke cluster database

---

# 3. API Security (POST /api/users)

## Input Validation

```go
type UserRequest struct {
    Name     string `json:"name" binding:"required,min=3,max=100"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
}
```
- Validasi format email dan kompleksitas password
- Reject input dengan karakter berbahaya (`;`, `DROP`, dll)

## Parameterized Queries

```go
db.Exec("INSERT INTO users VALUES (?, ?, ?, ?)", id, name, email, hashedPassword)
```
- Cegah SQL injection

## Password Hashing

```go
hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
```
- Gunakan bcrypt dengan cost factor 12-14

## Rate Limiting

- Batasi 5 request/menit per IP untuk endpoint registrasi
- Cegah brute-force attack

## HTTPS Enforcement

```go
router.Use(tlsMiddleware)
```
- Redirect HTTP ke HTTPS
- Gunakan HSTS header

## Content Security Policy

- Header `Content-Type: application/json` wajib
- Reject request dengan content type tidak valid

---

# 4. CI/CD Workflow

## Diagram

*(Tambahkan diagram di sini jika ada)*

## Tahapan Detail

### Commit

- Developer push ke branch `feature/*`
- Commit message follow convention (`feat:`, `fix:`, dll)

### CI Pipeline (GitHub Actions)

```yaml
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4
      - name: Set up Go
        uses: actions/setup-go@v4
      - name: Run tests
        run: go test ./...
```

### Code Quality Gate

- Unit test coverage > 80%
- Static analysis (SonarQube)
- Security scan (Trivy for containers)

### Artifact Build

- Build Docker image untuk tiap microservice
- Push image ke AWS ECR dengan tag commit-sha

### Staging Deployment

- Deploy ke Kubernetes staging namespace
- Auto-run integration tests (Postman/Selenium)

### Approval Workflow

- Manual approval via Slack notification
- QA team verifikasi fitur

### Production Deployment

- Blue-green deployment menggunakan service mesh (Istio)
- Canary release (10% traffic awal)
- Auto-rollback jika error rate > 5%

### Monitoring

- Log aggregation (ELK Stack)
- Real-time monitoring (Prometheus/Grafana)
- Alerting untuk error rate > 1%

---

# 5. Code Review & Debugging (Golang)

## Masalah dan Perbaikan

### SQL Injection

**Masalah:** id langsung dimasukkan ke query tanpa sanitasi  
**Perbaikan:** Gunakan prepared statement

```go
func FindUserByID(id string) (*User, error) {
    row := db.QueryRow("SELECT * FROM users WHERE id = ?", id)
    // ...
}
```

### Error Handling Tidak Spesifik

**Masalah:** Error database dikembalikan sebagai "Internal server error"  
**Perbaikan:** Pisahkan error not found

```go
user, err := db.FindUserByID(id)
if errors.Is(err, sql.ErrNoRows) {
    c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
    return
} else if err != nil {
    log.Printf("Database error: %v", err)
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Database failure"})
    return
}
```

### Konversi ID Tanpa Validasi

**Masalah:** id dari path langsung digunakan tanpa validasi format  
**Perbaikan:** Validasi UUID

```go
if !isValidUUID(id) {
    c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
    return
}
```

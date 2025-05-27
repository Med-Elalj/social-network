# Zone01Oujda Forum Server

This project sets up a local HTTPS-enabled Go web server with environment configuration and TLS certificate generation handled by a Python script.

---

## 🔧 Requirements

- Python 3.7+
- Go 1.24+
- `pip` (Python package manager)

---

## 📦 Setup Instructions

### 1. Install Python dependencies

```bash
pip install cryptography
```

### 2. Generate TLS certificates and .env file

```bash
python3 init.py
```

### 3. Build and run the Go server

```bash
go get
go run .
```
#### The server will start on:
```
https://localhost:8080
```

### REACT NEEDS SEPERATE TERMAL TO RUN OR RUN IT IN BACK GROUND
```React
npx create-react-app react-front
cd react-front
npm start

```
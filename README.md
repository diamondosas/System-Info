<div align="center">
  <img src="logo.png" alt="System Info Dashboard Logo" width="120"/>

  <h1>System Information Dashboard</h1>
  <p>A sleek, robust, and modern system specification dashboard built with Go and vanilla web technologies.</p>

  <p>
    <img src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go" />
    <img src="https://img.shields.io/badge/HTML5-E34F26?style=for-the-badge&logo=html5&logoColor=white" alt="HTML5" />
    <img src="https://img.shields.io/badge/CSS3-1572B6?style=for-the-badge&logo=css3&logoColor=white" alt="CSS3" />
    <img src="https://img.shields.io/badge/JavaScript-F7DF1E?style=for-the-badge&logo=javascript&logoColor=black" alt="JavaScript" />
    <img src="https://img.shields.io/badge/Wails-cc0000?style=for-the-badge&logo=wails&logoColor=white" alt="Wails" />
    <img src="https://img.shields.io/badge/License-MIT-green.svg?style=for-the-badge" alt="License: MIT" />
  </p>
</div>

---

## 🔧 Features

- **Comprehensive Metrics:** Collects and displays real-time data for CPU, Memory (including Swap), Storage, GPU, Network, OS, Battery, Sensors, and Processes.
- **Unified Architecture:** The Go backend directly serves the web dashboard files and the REST API—no external proxy needed!
- **Sleek Web Dashboard:** A modular vanilla JavaScript frontend (`Frontend/`) with an elegant dark theme.
- **Desktop Application:** Includes a Wails-powered desktop interface (`Backend/frontend/`) for native OS integration.
- **Cross-Platform Support:** Works efficiently across Windows and Linux environments.

---

## 🚀 Quick Start

### 1. Web Dashboard Mode

The simplest way to use the application is via the built-in HTTP server.

```bash
# Clone the repository
git clone https://github.com/DiamondOsas/System-Info.git
cd System-Info/Backend

# Build the CLI/Server backend
go build -o system-info-server

# Run the server
./system-info-server
```

> **Note:** The server will start on port `9999`. Navigate to `http://localhost:9999` in your web browser to view the live dashboard.

### 2. Desktop Mode (Wails)

If you prefer a native desktop application, you can build the Wails app.

```bash
# Prerequisites: Ensure you have Wails installed (https://wails.io/docs/gettingstarted/installation)
cd Backend

# Build the desktop app
wails build
```

The executable will be located in the `Backend/build/bin/` directory.

---

## 📊 Example API Output

The Go backend exposes a clean JSON HTTP API at `GET /api/specs`. Here is a sample:

```json
{
  "os": {
    "name": "linux 24.04",
    "version": "6.8.0",
    "architecture": "x86_64",
    "uptime": "0 days, 2 hours, 14 minutes",
    "boot_time": "2026-04-16 07:43:32"
  },
  "cpu": {
    "model": "Intel Core i7",
    "cores": 8,
    "threads": 16,
    "frequency": 3200,
    "usage_percent": 15.2
  },
  "memory": {
    "total": 16384,
    "available": 8500,
    "used": 7884,
    "swap_total": 4096,
    "swap_free": 3096,
    "swap_used": 1000,
    "swap_usage_percent": 24.4
  }
}
```

---

## 📂 Project Structure

```text
/
├── Backend/              ← Go backend source code
│   ├── frontend/         ← Wails desktop app frontend (JS/CSS)
│   ├── *.go              ← Modular Go files (cpu.go, memory.go, etc.)
├── Frontend/             ← Web Dashboard (HTML, CSS, JS modules)
│   ├── index.html
│   ├── styles.css
│   ├── js/               ← Modular JavaScript
├── logo.png              ← Project Logo
└── README.md             ← Documentation
```

---

## 🤝 Contribute

Contributions are highly welcome! To contribute:

1. Fork the repository.
2. Create a new branch (`git checkout -b feature/amazing-feature`).
3. Commit your changes (`git commit -m 'Add amazing feature'`).
4. Push to the branch (`git push origin feature/amazing-feature`).
5. Open a Pull Request.

---

## 📜 License & Contact

**License:** MIT License

**Author:** Diamond
**GitHub:** [DiamondOsas](https://github.com/DiamondOsas)

<div align="center">
  <i>Built with ❤️ for performance enthusiasts and system administrators.</i>
</div>

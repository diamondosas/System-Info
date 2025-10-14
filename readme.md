<div align="center">
  <img src="https://icons.iconarchive.com/icons/osullivanluke/orb-os-x/128/System-Info-icon.png" alt="System Info Dashboard Logo" width="100"/>
</div>

# System-Info

A simple system specification dashboard: **backend in Go** gathers hardware/OS details, and a **frontend visualizes them**.

---

## 🔧 What It Does

- Collects system metrics: CPU, memory, storage, GPU, network, OS, battery, processes, and sensors.
- Exposes data through a clean **JSON HTTP API** (e.g., `GET /api/specs`).
- Displays the specs in a clean **web dashboard** (vanilla JS + HTML/CSS).

---

## 📂 Project Structure

/
├── Backend/      ← Go code collecting system data
├── Frontend/     ← HTML, CSS, JS dashboard
├── logo.png      ← Your centered logo image
├── theme.png     ← Example theme image
├── api_post.txt
├── prompt.txt
├── README.md     ← This file
└── .gitignore


---

## 🚀 Quick Start

1.  **Build / Run** the Go backend (on your OS).
2.  **Start** the frontend (serve HTML/JS files).
3.  **Navigate** to the dashboard in your browser.
4.  Dashboard fetches data from the backend and displays live metrics.

> **Tip:** Adjust the host/port or API path if needed in your frontend JavaScript code.

---

## 📊 Example Output / Metrics

Here is a sample of the data structure the API exposes:

```json
{
"cpu": {
  "model": "Intel Core i7",
  "cores": 8,
  "usage_percent": 15.2
},
"memory": {
  "total": "16 GB",
  "available": "8.5 GB"
},
"storage": [
  {
    "device": "/dev/sda1",
    "used_gb": 450
  }
],
"os": {
  "platform": "Linux",
  "version": "5.15"
}
}
```
🔮 What’s Next / Ideas
We can make this project even better!

Better styling / theming (dark mode, transitions).

Auto-refresh or websocket streaming for real-time updates.

More sensors (temperatures, fan speeds).

Add filtering or search in processes.

Containerize the application (Docker).

Add user authentication or remote access.

🤝 Contribute
Fork the repository, make your changes, and send a Pull Request (PR)! If you add new features, please remember to update this README.

📜 License & Contact
License: MIT License (or your choice)

Author: Diamond
GitHub: DiamondOsas

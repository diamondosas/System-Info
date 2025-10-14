<!-- Inline SVG logo placeholder — replace with your own if you have one -->
<svg width="100" height="100" viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
  <circle cx="50" cy="50" r="48" stroke="#333" stroke-width="4" fill="none" />
  <text x="50%" y="54%" text-anchor="middle" fill="#333" font-size="36" font-family="Arial" dy=".3em">SI</text>
</svg>

# System-Info

A simple system specification dashboard: backend in Go gathers hardware/OS details, and a frontend visualizes them. :contentReference[oaicite:0]{index=0}

---

## 🔧 What It Does

- Collects system metrics: CPU, memory, storage, GPU, network, OS, battery, processes, sensors. :contentReference[oaicite:1]{index=1}  
- Exposes data through a JSON HTTP API (e.g. `GET /api/specs`) :contentReference[oaicite:2]{index=2}  
- Displays the specs in a clean web dashboard (vanilla JS + HTML/CSS) :contentReference[oaicite:3]{index=3}  

---

## 📂 Project Structure

/
├── Backend/ ← Go code collecting system data
├── Frontend/ ← HTML, CSS, JS dashboard
├── logo.png
├── theme.png
├── api_post.txt
├── prompt.txt
├── README.md ← this file
└── .gitignore

yaml
Copy code

---

## 🚀 Quick Start

1. Build / run the Go backend (on your OS)  
2. Start frontend (serve HTML/JS)  
3. Navigate to dashboard in browser  
4. Dashboard fetches data from backend and displays live metrics  

Adjust host/port or API path if needed in frontend JS.

---

## 📊 Example Output / Metrics

You should show screenshots or sample JSON with things like:

```json
{
  "cpu": { … },
  "memory": { … },
  "storage": [ … ],
  "gpu": { … },
  "network": [ … ],
  "os": { … },
  "battery": { … },
  "processes": [ … ]
}
And in the dashboard: charts, tables, dynamic updates, etc.

🔮 What’s Next / Ideas
Better styling / theming (dark mode, transitions)

Auto-refresh or websocket streaming

More sensors (temperatures, fan speeds)

Add filtering or search in processes

Containerize (Docker)

Add user authentication or remote access

🤝 Contribute
Fork, tweak, send PRs. If you add features, update this README.

📜 License & Contact
MIT License (or your choice)
Author: Diamond
GitHub: DiamondOsas

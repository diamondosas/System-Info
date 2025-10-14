<!-- Inline SVG logo placeholder — replace with your own if you have one -->
<div align="center">
  ![My App Logo]([https://www.google.com/url?sa=i&url=https%3A%2F%2Fwww.iconarchive.com%2Fshow%2Forb-os-x-icons-by-osullivanluke%2FSystem-Info-icon.html&psig=AOvVaw1N-J1I8da2B47g5wiZhe0h&ust=1760547855043000&source=images&cd=vfe&opi=89978449&ved=0CBUQjRxqFwoTCJj92IKWpJADFQAAAAAdAAAAABAE)
</div>


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
```

And the rest of your text:

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

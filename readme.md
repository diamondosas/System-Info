  You are to build a full-stack application that collects ALL possible system information from a computer and displays it in a clean web dashboard.

  ==============================
  📌 General Requirements
  ==============================
  - The application consists of:
    1. **Backend in Go** → Runs as an .exe (Windows optimized), collects full system specs, and exposes them via an HTTP API (JSON). For gui use wails(Vanilla Js )
    2. **Frontend ** → A dashboard that fetches data from the Go backend and displays all details nicely.
  - Leave comments in the code where I can adjust styles, Dont use external apis only the vue and font awsome font files , and add/remove features.
  - Use **JSON** as the communication format between backend and frontend.
  - Keep code simple

  USe theme.png for both fornt end and backend

  ==============================
  📌FrontEnd Requirements
  ==============================
  Use pure html css adnd vanilla js python flask as the server to recicve data as json form the go application that then displays it live to the user on the web

  The website theme shold be in theme



  ==============================
  📌 Backend (Go)
  ==============================
  When the app is opened it should show it actively collecting all info of the appp and then telling the user to close it and return back t the website the user presses close and the app still runs in the background to collect realtime info like ram cpu info and alll other things 
  - Use `gopsutil` and Windows-specific APIs (WMI) to collect all available information:
    - **CPU**: model, cores, threads, frequency, usage %, cache sizes
    - **Memory (RAM)**: total, available, used, swap memory
    - **Storage**: each disk/partition, model, type (HDD/SSD), capacity, free space
    - **GPU**: vendor, model, VRAM (via WMI on Windows, fallback to OpenCL/DirectX/GL if needed)
    - **Network**: all interfaces, IP addresses, MAC, current bandwidth usage
    - **OS**: name, version, build, architecture, uptime
    - **Battery**: status, percentage, charging/discharging
    - **Sensors (if available)**: CPU/GPU temperature, fan speed
    - **Processes**: running processes with name, PID, memory usage
  - Expose this data through an HTTP API at `http://localhost:9999/api/specs`.
  - Add a **health check endpoint** at `/api/health`.
  Also add other data that the app can collect 
  - JSON output should be structured clearly:
    ```json
    {
      "cpu": {...},
      "memory": {...},
      "storage": [...],
      "gpu": {...},
      "network": [...],
      "os": {...},
      "battery": {...},
      "sensors": {...},
      "processes": [...]
    }
  For the gui use pure vanilla(JS) with normal css

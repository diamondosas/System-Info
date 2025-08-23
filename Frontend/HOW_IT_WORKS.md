# System Info Dashboard - How It Works

## 1. Overview

This document breaks down the frontend of the System Information Dashboard. The frontend is a single-page web application built with HTML, CSS, and vanilla JavaScript. It is served by a simple Python Flask server that also acts as a proxy to the main Go backend, which is responsible for collecting the system data.

---

## 2. Architecture

The frontend operates on a three-tier model. The JavaScript running in the browser communicates with the Python Flask server, which in turn communicates with the Go data collection backend.

```
┌─────────────────┐   HTTP   ┌──────────────────┐   HTTP   ┌─────────────────┐
│   Web Browser   │◀───JSON──▶│  Python Server   │◀───JSON──▶│   Go Backend    │
│ (HTML/CSS/JS)   │          │ (Flask - Port 5000)│          │  (Port 9999)    │
└─────────────────┘          └──────────────────┘          └─────────────────┘
```

*   **Browser**: Renders the dashboard and requests data every 2 seconds.
*   **Python Server**: Serves the static files (HTML, CSS, JS) and forwards data requests to the Go backend.
*   **Go Backend**: The primary data source that collects system information.

---

## 3. File Breakdown

*   `index.html`: The skeleton of the application. It defines the structure and all the elements that will display the data.
*   `styles.css`: The visual layer. It contains all the styling rules for the dark blue theme, layout, responsiveness, and animations.
*   `script.js`: The brain of the application. It handles fetching data, processing it, and updating the UI in real-time.
*   `server.py`: The web server. It serves the application to the browser and acts as a bridge to the Go backend.

---

## 4. How Each File Works

### `index.html` (The Structure)

The HTML file is built using semantic tags (`<header>`, `<main>`, `<section>`, `<footer>`).

*   **Header**: Contains the logo and the main title.
*   **Main Content**: The `.container` holds all the information cards.
*   **Cards**: Each piece of information (OS, CPU, Memory, etc.) is organized into a `<section>` with the class `.card`. Each card has a unique `id` (e.g., `id="cpu-info"`) so the JavaScript can easily target it.
*   **Data Placeholders**: Inside the cards, `<span>` elements with unique `id`s (e.g., `id="cpu-model"`) act as placeholders. The JavaScript finds these `id`s and injects the data into them.
*   **File Links**: It links to the `styles.css` for styling, the `script.js` for logic, and external libraries like Google Fonts and Font Awesome for typography and icons.

### `styles.css` (The Look)

The stylesheet is responsible for the entire visual presentation.

*   **Theme Definition**: At the top, a comment block defines the core color palette for the "Dark Blue Theme". These colors are used throughout the file.
*   **Global Styles**: The `body` selector sets the global background color, font, and text color.
*   **Card Styling**: The `.card` class is the most important selector. It defines the background, border, shadow, and hover effects for all the information panels.
*   **Layout**: The `.row` class uses a `display: grid` to create a responsive two-column layout that automatically stacks to a single column on smaller screens.
*   **Component Styling**: There are specific styles for the progress bars (`.progress-bar`), tables (`.table-container`), and other components to ensure they match the theme.
*   **Responsive Design**: The `@media` query at the end of the file contains styles specifically for screens narrower than 768px, ensuring the dashboard is usable on mobile devices.

### `script.js` (The Logic)

This script brings the dashboard to life. It executes when the HTML document is fully loaded (`DOMContentLoaded`).

#### Configuration & State Management

```javascript
const CONFIG = {
    API_ENDPOINT: '/api/data',
    REFRESH_INTERVAL: 2000, // 2 seconds
    EXECUTABLE_PATH: './system-info-collector.exe' // Path to the backend executable
};

const state = {
    isConnected: false,
    refreshInterval: null,
    connectionRetryInterval: null
};
```

*   `CONFIG`: Contains configuration constants for the application
  * `API_ENDPOINT`: The URL endpoint to fetch system data from
  * `REFRESH_INTERVAL`: How often to refresh the data in milliseconds (2 seconds)
  * `EXECUTABLE_PATH`: Path to download the backend executable
*   `state`: Tracks the application state
  * `isConnected`: Boolean indicating if we're successfully connected to the backend
  * `refreshInterval`: Reference to the interval that refreshes data every 2 seconds
  * `connectionRetryInterval`: Reference to the interval that retries connection every 1 second

#### Core Functions

```javascript
function updateText(id, value) {
    const element = document.getElementById(id);
    if (element) {
        element.textContent = value;
    }
}
```

*   `updateText(id, value)`: Finds an HTML element by its ID and updates its text content. Used for simple text updates like CPU model or OS name.

```javascript
function updateProgressBar(id, percentage) {
    const element = document.getElementById(id);
    if (element) {
        element.style.width = `${percentage}%`;
        const textElement = document.getElementById(`${id}-text`);
        if(textElement) textElement.textContent = `${percentage}%`;
    }
}
```

*   `updateProgressBar(id, percentage)`: Updates the width of a progress bar element and its associated text to show a percentage value. Used for CPU and Memory usage.

```javascript
function updateConnectionStatus(status, message) {
    elements.statusText.textContent = message;
    
    // Remove all status classes
    elements.statusIcon.classList.remove('connected', 'disconnected', 'connecting');
    
    // Add appropriate status class
    elements.statusIcon.classList.add(status);
    
    // Update global connection state
    state.isConnected = (status === 'connected');
}
```

*   `updateConnectionStatus(status, message)`: Updates the connection status indicator in the UI. It changes the icon color and text based on the connection status (connecting, connected, disconnected).

```javascript
function formatBytes(bytes, decimals = 2) {
    if (bytes === 0) return '0 Bytes';
    
    const k = 1024;
    const dm = decimals < 0 ? 0 : decimals;
    const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
    
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    
    return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
}
```

*   `formatBytes(bytes, decimals = 2)`: Converts a byte value into a human-readable format (KB, MB, GB, etc.). For example, 1024 bytes becomes "1 KB".

#### Data Fetching Functions

```javascript
async function fetchSystemInfo() {
    try {
        const response = await fetch(CONFIG.API_ENDPOINT);
        
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        const data = await response.json();
        return data;
    } catch (error) {
        console.error('Error fetching system info:', error);
        throw error;
    }
}
```

*   `fetchSystemInfo()`: An asynchronous function that fetches system data from the backend API endpoint using the browser's `fetch` API. It handles HTTP errors and returns parsed JSON data.

```javascript
async function fetchData() {
    try {
        const data = await fetchSystemInfo();
        updateAllInfo(data);
    } catch (error) {
        updateConnectionStatus('disconnected', `Error: ${error.message}`);
        console.error('Error fetching system info:', error);
        // Start connection retries if not already running
        if (!state.connectionRetryInterval) {
            startConnectionRetries();
        }
    }
}
```

*   `fetchData()`: Fetches data from the backend and updates the UI. If there's an error, it updates the connection status and starts connection retries if not already running.

```javascript
function startConnectionRetries() {
    // Clear any existing retry interval
    if (state.connectionRetryInterval) {
        clearInterval(state.connectionRetryInterval);
    }
    
    // Set up new retry interval (every 1 second)
    state.connectionRetryInterval = setInterval(async () => {
        if (!state.isConnected) {
            try {
                const data = await fetchSystemInfo();
                updateAllInfo(data);
                // If successful, clear the retry interval
                clearInterval(state.connectionRetryInterval);
                state.connectionRetryInterval = null;
            } catch (error) {
                // Continue retrying if still disconnected
                updateConnectionStatus('disconnected', `Error: ${error.message} (Retrying...)`);
                console.error('Connection retry failed:', error);
            }
        }
    }, 1000);
}
```

*   `startConnectionRetries()`: Implements continuous background attempts to connect to the backend every second when not connected. It creates an interval that tries to fetch data every second until successful.

```javascript
function startAutoRefresh() {
    // Clear any existing interval
    if (state.refreshInterval) {
        clearInterval(state.refreshInterval);
    }
    
    // Set up new interval
    state.refreshInterval = setInterval(fetchData, CONFIG.REFRESH_INTERVAL);
}
```

*   `startAutoRefresh()`: Sets up an interval to automatically refresh the data every 2 seconds (as defined in CONFIG.REFRESH_INTERVAL).

#### UI Update Functions

```javascript
function updateAllInfo(data) {
    // Update connection status to connected
    updateConnectionStatus('connected', 'Connected to backend');
    
    // OS Information
    updateText('os-name', data.os?.name || 'N/A');
    updateText('os-version', data.os?.version || 'N/A');
    updateText('os-arch', data.os?.architecture || 'N/A');
    updateText('os-uptime', data.os?.uptime || 'N/A');

    // CPU Information
    updateText('cpu-model', data.cpu?.model || 'N/A');
    updateText('cpu-cores', data.cpu?.cores || 'N/A');
    updateText('cpu-threads', data.cpu?.threads || 'N/A');
    updateText('cpu-frequency', data.cpu?.frequency ? `${data.cpu.frequency} MHz` : 'N/A');
    // Round CPU usage to whole number
    const cpuUsage = Math.round(data.cpu?.usage_percent || 0);
    updateProgressBar('cpu-usage', cpuUsage);

    // Memory Information
    updateText('mem-total', data.memory?.total ? formatBytes(data.memory.total * 1024 * 1024) : 'N/A');
    updateText('mem-available', data.memory?.available ? formatBytes(data.memory.available * 1024 * 1024) : 'N/A');
    updateText('mem-used', data.memory?.used ? formatBytes(data.memory.used * 1024 * 1024) : 'N/A');
    const memoryUsage = data.memory?.total && data.memory?.used ? 
        Math.round((data.memory.used / data.memory.total) * 100) : 0;
    updateProgressBar('mem-usage', memoryUsage);

    // GPU Information
    updateText('gpu-vendor', data.gpu?.vendor || 'N/A');
    updateText('gpu-model', data.gpu?.model || 'N/A');
    updateText('gpu-vram', data.gpu?.vram ? formatBytes(data.gpu.vram * 1024 * 1024) : 'N/A');

    // Storage Information
    const storageContainer = document.getElementById('storage-devices');
    if (storageContainer) {
        storageContainer.innerHTML = '';
        if (data.storage && Array.isArray(data.storage)) {
            data.storage.forEach(device => {
                const deviceDiv = document.createElement('div');
                deviceDiv.className = 'device';
                deviceDiv.innerHTML = `
                    <h4>Disk: ${device.device || 'N/A'}</h4>
                    <p><strong>Model:</strong> ${device.model || 'N/A'}</p>
                    <p><strong>Type:</strong> ${device.type || 'N/A'}</p>
                    <p><strong>Capacity:</strong> ${device.capacity ? formatBytes(device.capacity * 1024 * 1024) : 'N/A'}</p>
                    <p><strong>Free Space:</strong> ${device.free ? formatBytes(device.free * 1024 * 1024) : 'N/A'}</p>
                `;
                storageContainer.appendChild(deviceDiv);
            });
        } else {
            storageContainer.innerHTML = '<p>No storage devices found</p>';
        }
    }

    // Network Information
    const networkContainer = document.getElementById('network-interfaces');
    if (networkContainer) {
        networkContainer.innerHTML = '';
        if (data.network && Array.isArray(data.network)) {
            // Filter to only show Ethernet and Wi-Fi interfaces
            const filteredNetworks = data.network.filter(iface => 
                iface.interface && 
                (iface.interface.toLowerCase().includes('ethernet') || 
                 iface.interface.toLowerCase().includes('wi-fi') ||
                 iface.interface.toLowerCase().includes('wifi'))
            );
            
            filteredNetworks.forEach(iface => {
                const ifaceDiv = document.createElement('div');
                ifaceDiv.className = 'interface';
                ifaceDiv.innerHTML = `
                    <h4>Interface: ${iface.interface || 'N/A'}</h4>
                    <p><strong>IP Address:</strong> ${iface.ip_address || 'N/A'}</p>
                    <p><strong>MAC Address:</strong> ${iface.mac_address || 'N/A'}</p>
                    <p><strong>Bandwidth:</strong> ${iface.bandwidth_down ? `${formatBytes(iface.bandwidth_down)}/s` : 'N/A'}</p>
                `;
                networkContainer.appendChild(ifaceDiv);
            });
            
            // Show message if no Ethernet or Wi-Fi interfaces found
            if (filteredNetworks.length === 0) {
                networkContainer.innerHTML = '<p>No Ethernet or Wi-Fi interfaces found</p>';
            }
        } else {
            networkContainer.innerHTML = '<p>No network interfaces found</p>';
        }
    }

    // Battery Information
    updateText('battery-status', data.battery?.status || 'N/A');
    updateText('battery-percentage', data.battery?.percentage ? `${data.battery.percentage}%` : 'N/A');

    // Sensors Information
    updateText('cpu-temp', data.sensors?.cpu_temp ? `${data.sensors.cpu_temp}°C` : 'N/A');
    updateText('gpu-temp', data.sensors?.gpu_temp ? `${data.sensors.gpu_temp}°C` : 'N/A');

    // Processes Information
    const processesBody = document.getElementById('processes-body');
    if (processesBody) {
        processesBody.innerHTML = '';
        if (data.processes && Array.isArray(data.processes)) {
            // Sort processes by memory usage (descending)
            const sortedProcesses = [...data.processes].sort((a, b) => 
                (b.memory_mb || 0) - (a.memory_mb || 0)
            );
            
            // Show top 20 processes
            const topProcesses = sortedProcesses.slice(0, 20);
            
            topProcesses.forEach(proc => {
                const row = document.createElement('tr');
                row.innerHTML = `
                    <td>${proc.pid || 'N/A'}</td>
                    <td>${proc.name || 'N/A'}</td>
                    <td>${proc.memory_mb ? formatBytes(proc.memory_mb * 1024 * 1024) : 'N/A'}</td>
                `;
                processesBody.appendChild(row);
            });
        } else {
            const row = document.createElement('tr');
            row.innerHTML = '<td colspan="3">No processes found</td>';
            processesBody.appendChild(row);
        }
    }
}
```

*   `updateAllInfo(data)`: The main rendering function that takes the data object from the backend and populates the entire dashboard. It calls `updateText` and `updateProgressBar` for each piece of information, and dynamically creates HTML for multi-item sections like Storage and Network. It also implements:
  * Filtering network interfaces to only show Ethernet and Wi-Fi
  * Rounding CPU usage to whole numbers

#### Event Handlers

```javascript
function handleDownloadClick() {
    // Create a temporary link element
    const link = document.createElement('a');
    link.href = CONFIG.EXECUTABLE_PATH;
    link.download = 'system-info-collector.exe';
    link.style.display = 'none';
    
    // Add to DOM, click, and remove
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
}
```

*   `handleDownloadClick()`: Handles the download button click by creating a temporary link element that triggers the download of the backend executable.

#### Initialization

```javascript
function init() {
    // Set initial connection status
    updateConnectionStatus('connecting', 'Connecting to backend...');
    
    // Initialize event listeners
    initializeEventListeners();
    
    // Start data fetching
    fetchData();
    startAutoRefresh();
    
    // Start connection retries in case of disconnection
    startConnectionRetries();
}

// Start the application when the DOM is loaded
document.addEventListener('DOMContentLoaded', init);
```

*   `init()`: Initializes the application when the DOM is fully loaded. It sets the initial connection status, initializes event listeners, starts the initial data fetch, starts the auto-refresh interval, and starts connection retries.
*   `document.addEventListener('DOMContentLoaded', init)`: Registers the init function to run when the DOM is fully loaded.

### `server.py` (The Bridge)

This is a lightweight web server using the Flask framework.

*   **Flask App**: It initializes a Flask application.
*   **API Route (`/api/data`)**: This is the most important part. When the frontend's JavaScript requests this URL, the Flask server makes its own HTTP request to the Go backend (`http://localhost:9999/api/specs`). It then forwards the JSON response from the Go backend directly to the frontend. This makes it a **proxy**.
*   **Frontend Route (`/`)**: This route serves the `index.html` file, allowing you to view the dashboard in your browser.
*   **Static Files**: The server is configured to automatically serve other files like `styles.css` and `script.js` from the same directory.
*   **Execution**: When you run `python server.py`, it starts the server on `http://localhost:5000`.

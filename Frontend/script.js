// Configuration
const CONFIG = {
    API_ENDPOINT: '/api/data',
    REFRESH_INTERVAL: 2000, // 2 seconds
    EXECUTABLE_PATH: './info-collector.exe' // Path to the backend executable
};

// DOM Elements
const elements = {
    downloadBtn: document.getElementById('download-btn'),
    connectionStatus: document.getElementById('connection-status'),
    statusIcon: document.querySelector('.status-icon'),
    statusText: document.getElementById('status-text')
};

// State Management
const state = {
    isConnected: false,
    refreshInterval: null,
    connectionRetryInterval: null
};

// Utility Functions

/**
 * Updates text content of an element by ID
 * @param {string} id - Element ID
 * @param {string} value - Text content to set
 */
function updateText(id, value) {
    const element = document.getElementById(id);
    if (element) {
        element.textContent = value;
    }
}

/**
 * Updates progress bar width and text
 * @param {string} id - Progress bar ID
 * @param {number} percentage - Percentage value (0-100)
 */
function updateProgressBar(id, percentage) {
    const element = document.getElementById(id);
    if (element) {
        element.style.width = `${percentage}%`;
        const textElement = document.getElementById(`${id}-text`);
        if(textElement) textElement.textContent = `${percentage}%`;
    }
}

/**
 * Updates the connection status indicator
 * @param {string} status - Status type: 'connected', 'disconnected', 'connecting'
 * @param {string} message - Status message to display
 */
function updateConnectionStatus(status, message) {
    elements.statusText.textContent = message;
    
    // Remove all status classes
    elements.statusIcon.classList.remove('connected', 'disconnected', 'connecting');
    
    // Add appropriate status class
    elements.statusIcon.classList.add(status);
    
    // Update global connection state
    state.isConnected = (status === 'connected');
}

/**
 * Formats bytes into human-readable format
 * @param {number} bytes - Number of bytes
 * @param {number} decimals - Decimal places
 * @returns {string} Formatted string
 */
function formatBytes(bytes, decimals = 2) {
    if (bytes === 0) return '0 Bytes';
    
    const k = 1024;
    const dm = decimals < 0 ? 0 : decimals;
    const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
    
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    
    return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
}

// UI Update Functions

/**
 * Updates all system information on the page
 * @param {Object} data - System information data
 */
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

// Data Fetching Functions

/**
 * Fetches system information from the backend
 * @returns {Promise<Object>} System information data
 */
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

/**
 * Fetches data and updates the UI
 */
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

/**
 * Attempts to connect to the backend continuously every second when disconnected
 */
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

/**
 * Stops the connection retry attempts
 */
function stopConnectionRetries() {
    if (state.connectionRetryInterval) {
        clearInterval(state.connectionRetryInterval);
        state.connectionRetryInterval = null;
    }
}

/**
 * Starts the auto-refresh interval
 */
function startAutoRefresh() {
    // Clear any existing interval
    if (state.refreshInterval) {
        clearInterval(state.refreshInterval);
    }
    
    // Set up new interval
    state.refreshInterval = setInterval(fetchData, CONFIG.REFRESH_INTERVAL);
}

/**
 * Stops the auto-refresh interval
 */
function stopAutoRefresh() {
    if (state.refreshInterval) {
        clearInterval(state.refreshInterval);
        state.refreshInterval = null;
    }
}

// Event Handlers

/**
 * Handles the download button click
 */
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

/**
 * Initializes event listeners
 */
function initializeEventListeners() {
    // Download button
    if (elements.downloadBtn) {
        elements.downloadBtn.addEventListener('click', handleDownloadClick);
    }
    
    // Handle page visibility changes to optimize resource usage
    document.addEventListener('visibilitychange', () => {
        if (document.hidden) {
            stopAutoRefresh();
        } else {
            startAutoRefresh();
        }
    });
}

// Initialization

/**
 * Initializes the application
 */
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

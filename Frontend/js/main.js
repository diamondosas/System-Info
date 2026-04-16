import { fetchSystemInfo, CONFIG } from './api.js';
import { updateAllInfo, updateConnectionStatus } from './ui.js';

let state = {
    isConnected: false,
    refreshInterval: null,
    connectionRetryInterval: null
};

async function fetchData() {
    try {
        const data = await fetchSystemInfo();
        state.isConnected = true;
        updateAllInfo(data);
    } catch (error) {
        state.isConnected = false;
        updateConnectionStatus('disconnected', `Error: ${error.message}`);
        if (!state.connectionRetryInterval) {
            startConnectionRetries();
        }
    }
}

function startConnectionRetries() {
    if (state.connectionRetryInterval) {
        clearInterval(state.connectionRetryInterval);
    }

    state.connectionRetryInterval = setInterval(async () => {
        if (!state.isConnected) {
            try {
                const data = await fetchSystemInfo();
                state.isConnected = true;
                updateAllInfo(data);
                clearInterval(state.connectionRetryInterval);
                state.connectionRetryInterval = null;
            } catch (error) {
                updateConnectionStatus('disconnected', `Error: ${error.message} (Retrying...)`);
            }
        }
    }, 1000);
}

function startAutoRefresh() {
    if (state.refreshInterval) {
        clearInterval(state.refreshInterval);
    }
    state.refreshInterval = setInterval(fetchData, CONFIG.REFRESH_INTERVAL);
}

function stopAutoRefresh() {
    if (state.refreshInterval) {
        clearInterval(state.refreshInterval);
        state.refreshInterval = null;
    }
}

function init() {
    updateConnectionStatus('connecting', 'Connecting to backend...');

    document.addEventListener('visibilitychange', () => {
        if (document.hidden) {
            stopAutoRefresh();
        } else {
            startAutoRefresh();
        }
    });

    fetchData();
    startAutoRefresh();
    startConnectionRetries();
}

document.addEventListener('DOMContentLoaded', init);

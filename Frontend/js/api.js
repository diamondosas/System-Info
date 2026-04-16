export const CONFIG = {
    API_ENDPOINT: '/api/specs',
    REFRESH_INTERVAL: 2000,
    EXECUTABLE_PATH: './info-collector.exe'
};

export async function fetchSystemInfo() {
    const response = await fetch(CONFIG.API_ENDPOINT);
    if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
    }
    return response.json();
}

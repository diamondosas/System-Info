import './main.css';
import logo from './assets/images/logo.png';
import { GetSpecs } from '../wailsjs/go/main/App';

function formatBytes(bytes, decimals = 2) {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const dm = decimals < 0 ? 0 : decimals;
    const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
}

document.querySelector('#app').innerHTML = `
    <div class="container" style="--wails-draggable:drag">
        <div class="header">
            <img id="logo" class="logo" src="${logo}" alt="App Logo" style="--wails-draggable:no-drag">
            <h1 style="--wails-draggable:no-drag">System Info Dashboard</h1>
        </div>

        <div class="content" style="--wails-draggable:no-drag">
            <div id="loading-state" class="status-container">
                <div class="status-indicator loading">
                    <div class="spinner"></div>
                    <span>Collecting system information...</span>
                </div>
            </div>
            
            <div id="data-state" style="display: none; width: 100%;">
                <div class="grid">
                    <div class="card">
                        <h3>OS</h3>
                        <div class="value" id="val-os">Loading...</div>
                        <div class="sub" id="sub-os"></div>
                    </div>
                    <div class="card">
                        <h3>CPU</h3>
                        <div class="value" id="val-cpu">Loading...</div>
                        <div class="sub" id="sub-cpu"></div>
                    </div>
                    <div class="card">
                        <h3>Memory</h3>
                        <div class="value" id="val-mem">Loading...</div>
                        <div class="sub" id="sub-mem"></div>
                    </div>
                </div>
                <div style="margin-top: 15px; font-size: 0.9em; text-align: center; color: #8b949e;">
                    Running Web Dashboard Server at <a href="http://localhost:9999" target="_blank" style="color: #388bfd; text-decoration: none;">http://localhost:9999</a>
                </div>
            </div>
            
            <div id="error-state" class="status-container" style="display: none; color: #f44336;">
                <div style="padding: 15px; background: rgba(244, 67, 54, 0.1); border-radius: 8px;">
                    Error collecting data. See console.
                </div>
            </div>
        </div>

        <div class="footer">
            <button id="close-btn" class="close-btn" style="--wails-draggable:no-drag">Close</button>
        </div>
    </div>
`;

const loadingState = document.getElementById('loading-state');
const dataState = document.getElementById('data-state');
const errorState = document.getElementById('error-state');

async function init() {
    document.getElementById('close-btn').addEventListener('click', () => {
        window.runtime.WindowHide();
    });

    try {
        const specs = await GetSpecs();

        document.getElementById('val-os').textContent = specs.os.name;
        document.getElementById('sub-os').textContent = specs.os.architecture;
        
        document.getElementById('val-cpu').textContent = specs.cpu.model.substring(0, 20) + (specs.cpu.model.length > 20 ? '...' : '');
        document.getElementById('sub-cpu').textContent = \`\${specs.cpu.cores} Cores, \${specs.cpu.threads} Threads\`;
        
        document.getElementById('val-mem').textContent = \`\${formatBytes(specs.memory.used * 1024 * 1024)} / \${formatBytes(specs.memory.total * 1024 * 1024)}\`;
        const memPercent = Math.round((specs.memory.used / specs.memory.total) * 100) || 0;
        document.getElementById('sub-mem').textContent = \`\${memPercent}% Used\`;
        
        loadingState.style.display = 'none';
        dataState.style.display = 'block';
    } catch (error) {
        console.error("Error collecting specs:", error);
        loadingState.style.display = 'none';
        errorState.style.display = 'block';
    }
}

document.addEventListener('DOMContentLoaded', init);

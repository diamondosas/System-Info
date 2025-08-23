// Import CSS files
import './main.css';

// Import logo with error handling
import logo from './assets/images/logo.png';
import { Greet } from '../wailsjs/go/main/App';
// Remove the incorrect import - WindowStartDrag doesn't exist in Wails v2
// import { WindowStartDrag } from '../wailsjs/runtime/runtime.js';

// Create the collection interface HTML structure
try {
    document.querySelector('#app').innerHTML = `
        <div class="container" id="container" style="--wails-draggable:drag">
            <div class="header">
                <img id="logo" class="logo" src="${logo}" alt="App Logo">
                <h1>Collect System Information</h1>
            </div>
            
            <div class="content">
                <div class="status-container">
                    <div class="status-indicator" id="status-indicator">
                        <div class="spinner"></div>
                        <span id="status-text">Initializing system information collection...</span>
                    </div>
                </div>
                
                <div class="progress-container">
                    <div class="progress-bar">
                        <div class="progress-fill" id="progress-fill"></div>
                    </div>
                    <div class="progress-text" id="progress-text">0%</div>
                    <div class="completion-message" id="completion-message" style="display: none;">
                        Click Close and Return to Webpage to View System Info
                    </div>
                </div>
            </div>
            
            <div class="footer">
                <button id="close-btn" class="close-btn" style="--wails-draggable:no-drag">Close</button>
            </div>
        </div>
    `;
} catch (error) {
    console.error('Error creating HTML structure:', error);
    document.querySelector('#app').innerHTML = '<div style="color: white; padding: 20px;">Error loading application interface</div>';
}

// DOM Elements
const statusIndicator = document.getElementById('status-indicator');
const statusText = document.getElementById('status-text');
const progressFill = document.getElementById('progress-fill');
const progressText = document.getElementById('progress-text');
const closeBtn = document.getElementById('close-btn');
const completionMessage = document.getElementById('completion-message');

// Debug: Check if elements were found
console.log('DOM Elements:', {
    statusIndicator,
    statusText,
    progressFill,
    progressText,
    closeBtn,
    completionMessage
});

// Update status indicator
function updateStatus(text, isLoading = false) {
    statusText.textContent = text;
    if (isLoading) {
        statusIndicator.classList.add('loading');
    } else {
        statusIndicator.classList.remove('loading');
    }
}

// Update progress bar
function updateProgress(percentage) {
    progressFill.style.width = `${percentage}%`;
    progressText.textContent = `${Math.round(percentage)}%`;
}

// Add visual indication to close button without animation
function highlightCloseButton() {
    closeBtn.style.backgroundColor = '#4a86e8';
    closeBtn.style.boxShadow = '0 0 10px #4a86e8';
}

// Simulate collection process
async function simulateCollection() {
    // Steps in the collection process
    const steps = [
        { text: "Initializing system information collection...", progress: 0 },
        { text: "Collecting CPU information...", progress: 10 },
        { text: "Collecting memory information...", progress: 20 },
        { text: "Collecting storage information...", progress: 30 },
        { text: "Collecting GPU information...", progress: 40 },
        { text: "Collecting network information...", progress: 50 },
        { text: "Collecting OS information...", progress: 60 },
        { text: "Collecting battery information...", progress: 70 },
        { text: "Collecting sensor information...", progress: 80 },
        { text: "Collecting process information...", progress: 90 },
        { text: "Finalizing data collection...", progress: 100 }
    ];
    
    // Process each step with a delay to show progress
    for (const step of steps) {
        updateStatus(step.text, true);
        updateProgress(step.progress);
        // Wait for 300ms to simulate work
        await new Promise(resolve => setTimeout(resolve, 300));
    }
    
    // Show completion message
    updateStatus("System information collection complete!");
    
    // Hide spinner
    const spinner = document.querySelector('.spinner');
    if (spinner) {
        spinner.style.display = 'none';
    }
    
    // Show completion message
    completionMessage.style.display = 'block';
    
    // Highlight close button
    highlightCloseButton();
}

// Initialize the application
async function init() {
    try {
        // Window dragging is now handled by CSS --wails-draggable:drag attribute
        // No JavaScript event listeners needed for dragging
        
        // Set up close button with Wails runtime function
        closeBtn.addEventListener('click', () => {
            // Use Wails runtime function to hide the window instead of quitting
            window.runtime.WindowHide();
        });
        
        // Start simulation
        await simulateCollection();
        
        // Show initial greeting
        try {
            const greeting = await Greet("User");
            console.log(greeting);
        } catch (err) {
            console.error("Error with greeting:", err);
        }
    } catch (error) {
        console.error("Error initializing application:", error);
        // Show error message in the UI
        const appElement = document.querySelector('#app');
        if (appElement) {
            appElement.innerHTML = '<div style="color: white; padding: 20px;">Error initializing application: ' + error.message + '</div>';
        }
    }
}

// Initialize when the DOM is loaded
document.addEventListener('DOMContentLoaded', init);
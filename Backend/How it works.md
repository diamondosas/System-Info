# System Information Collector - How It Works

This document provides a detailed overview of the System Information Collector application architecture and functionality, with specific explanations of the JavaScript implementation.

## Overview

The System Information Collector is a Go application built with the Wails framework that collects detailed system information from Windows machines and displays it in a web-based dashboard. The application consists of two main components:

1. **Backend (Go)** - Collects system information using various system APIs
2. **Frontend (HTML/CSS/JavaScript)** - Displays the information collection process in a user-friendly interface within a Wails desktop window

## Architecture

```
┌─────────────────────────────────────┐
│         Wails Frontend (Web)        │
│  ┌──────────────────────────────┐   │
│  │     Collection Interface     │   │
│  │  (HTML/CSS/JavaScript)       │   │
│  └──────────────────────────────┘   │
└─────────────────┬───────────────────┘
                  │ HTTP API Calls
                  ▼
┌─────────────────────────────────────┐
│           Backend (Go)              │
│  ┌──────────────────────────────┐   │
│  │      Wails Application       │   │
│  │    (main.go + app.go)        │   │
│  └──────────────────────────────┘   │
│  ┌──────────────────────────────┐   │
│  │    System Information        │   │
│  │       Collector              │   │
│  │ (gopsutil + WMI Queries)     │   │
│  └──────────────────────────────┘   │
│  ┌──────────────────────────────┐   │
│  │     Embedded HTTP Server     │   │
│  │   (Provides REST API)        │   │
│  └──────────────────────────────┘   │
└─────────────────────────────────────┘
```

## Backend Components

### Main Application Structure

- **main.go**: Entry point of the application. Initializes the Wails application with configuration settings.
- **app.go**: Contains the core application logic including:
  - System information collection functions
  - HTTP server implementation
  - Data structures for representing system information

### System Information Collection

The backend collects system information using two primary methods:

1. **gopsutil Library**: A cross-platform library for retrieving system information
   - CPU information (model, cores, usage)
   - Memory information (total, available, used)
   - Disk usage statistics
   - Network interface details
   - Running processes

2. **Windows Management Instrumentation (WMI) Queries**: Direct Windows API calls for detailed system information
   - Operating system details (name, version, uptime)
   - GPU information (vendor, model, VRAM)
   - Storage devices (disk drives, partitions, logical disks)
   - Battery status
   - System sensors (temperature)

### HTTP Server

The backend includes an embedded HTTP server that exposes two endpoints:

1. `GET /api/specs` - Returns all collected system information in JSON format
2. `GET /api/health` - Returns a simple health check response

The server runs on port 9999 and is started when the application launches.

## Frontend Components (Wails Collection Interface)

### Technology Stack

- **Vanilla JavaScript**: No frameworks, just plain JavaScript for DOM manipulation
- **CSS3**: For styling and responsive design

### Key Files

- **main.js**: Main JavaScript file that:
  - Creates the collection interface UI dynamically
  - Simulates the collection process with progress updates
  - Handles user interactions (close button)
- **main.css**: Styling for the collection interface
- **style.css**: Global styling for the application

### Detailed JavaScript Code Explanation

#### DOM Initialization and Structure

```javascript
document.querySelector('#app').innerHTML = `
    <div class="container">
        <div class="header">
            <img id="logo" class="logo" src="${logo}" alt="App Logo">
            <h1>System Information Collector</h1>
        </div>
        <!-- Rest of the HTML structure -->
        <div class="footer">
            <button id="close-btn" class="close-btn">Close</button>
        </div>
    </div>
`;
```

This code creates the entire collection interface UI structure dynamically when the application loads. It includes:
- A header with logo and title
- A status indicator showing collection progress
- A progress bar to visualize the collection process
- A footer with a close button

Key changes from the previous version:
1. Removed the "✕" close button from the header
2. Added a new "Close" button in the footer with text instead of just an icon
3. Added a completion message that appears after collection is finished

#### DOM Element Selection

```javascript
const statusIndicator = document.getElementById('status-indicator');
const statusText = document.getElementById('status-text');
const progressFill = document.getElementById('progress-fill');
const progressText = document.getElementById('progress-text');
const closeBtn = document.getElementById('close-btn');
const completionMessage = document.getElementById('completion-message');
```

These lines select DOM elements by their IDs and store them in variables for later use. This is more efficient than querying the DOM repeatedly, as it only needs to be done once when the page loads.

New element added:
- `completionMessage`: References the completion message div that appears after collection

#### UI Update Functions

```javascript
function updateStatus(text, isLoading = false) {
    statusText.textContent = text;
    if (isLoading) {
        statusIndicator.classList.add('loading');
    } else {
        statusIndicator.classList.remove('loading');
    }
}
```

This function updates the status indicator at the top of the interface:
- `statusText.textContent = text;` - Changes the displayed text to show the current collection step
- `statusIndicator.classList.add('loading');` - Adds a loading class to show the spinner animation
- `statusIndicator.classList.remove('loading');` - Removes the loading class to hide the spinner animation

```javascript
function updateProgress(percentage) {
    progressFill.style.width = `${percentage}%`;
    progressText.textContent = `${Math.round(percentage)}%`;
}
```

This function updates the progress bar:
- `progressFill.style.width = `${percentage}%`;` - Sets the width of the progress fill element based on the percentage
- `progressText.textContent = `${Math.round(percentage)}%`;` - Updates the progress text to show the current percentage

#### New Highlight Function

```javascript
function highlightCloseButton() {
    closeBtn.style.backgroundColor = '#4a86e8';
    closeBtn.style.boxShadow = '0 0 10px #4a86e8';
}
```

This new function adds a visual indication to the close button without animation:
- `closeBtn.style.backgroundColor = '#4a86e8';` - Sets the button background to a solid blue color
- `closeBtn.style.boxShadow = '0 0 10px #4a86e8';` - Adds a blue glow effect around the button

This creates a clear visual cue that draws the user's attention to the close button after collection is complete, without the distracting color changes.

#### Collection Process Simulation

```javascript
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
```

This asynchronous function simulates the collection process with several new features:
1. Defines steps with text descriptions and progress percentages
2. Iterates through each step with a delay to simulate work
3. Updates the status text and progress bar for each step
4. Shows a completion message when finished
5. **New**: Hides the loading spinner after completion
6. **New**: Shows the completion message ("Return to the webpage to view your system information")
7. **New**: Highlights the close button with a glow effect to draw attention to it

#### Event Handling

```javascript
closeBtn.addEventListener('click', () => {
    window.runtime.WindowHide();
});
```

This line sets up an event listener for the close button:
- When clicked, it hides the window using the `window.runtime.WindowHide()` method instead of closing the entire application
- This allows the application to continue running in the background while still providing the API data

#### Application Initialization

```javascript
async function init() {
    // Make window draggable by clicking anywhere on the container
    const container = document.getElementById('container');
    container.addEventListener('mousedown', (e) => {
        // Check if the click is not on the close button
        if (e.target !== closeBtn) {
            WindowDrag();
        }
    });
    
    // Set up close button with Wails runtime function
    closeBtn.addEventListener('click', () => {
        // Use Wails runtime function to close the window
        window.runtime.Quit();
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
}

// Initialize when the DOM is loaded
document.addEventListener('DOMContentLoaded', init);
```

The initialization function has been updated with two important new features:

1. **Window Dragging Functionality**:
   - Window dragging is now handled by CSS attributes (`--wails-draggable:drag`) instead of JavaScript event listeners
   - The container div has `style="--wails-draggable:drag"` which makes the entire window draggable
   - The close button has `style="--wails-draggable:no-drag"` which prevents it from interfering with the dragging functionality

2. **Enhanced Close Button**:
   - `closeBtn.addEventListener('click', () => { window.runtime.Quit(); });` - Uses the Wails runtime `Quit()` function instead of `window.close()` for proper application termination

The `document.addEventListener('DOMContentLoaded', init);` line ensures that the initialization function is called only after the HTML document has been completely loaded and parsed.

### CSS Styling Explanation

The CSS has been updated to implement the requested visual changes:

#### Acrylic Effect and Rounded Corners

```css
html, body {
    background-color: rgba(27, 38, 54, 0.7); /* Translucent background for acrylic effect */
    border-radius: 15px; /* Rounded corners */
    backdrop-filter: blur(10px); /* Blur effect for acrylic look */
}
```

These properties create the acrylic effect:
- `background-color: rgba(27, 38, 54, 0.7);` - Makes the background translucent (70% opacity)
- `border-radius: 15px;` - Rounds the corners of the window
- `backdrop-filter: blur(10px);` - Applies a blur effect to elements behind the window

#### New Close Button Styling

```css
.close-btn {
    background-color: #4a86e8;
    color: white;
    border: none;
    padding: 12px 30px;
    border-radius: 25px;
    cursor: pointer;
    font-size: 1rem;
    font-weight: 500;
    transition: background-color 0.3s ease;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.close-btn:hover {
    transform: translateY(-2px);
    box-shadow: 0 6px 8px rgba(0, 0, 0, 0.15);
}
```

The close button CSS has been completely redesigned:
- Changed from a small circular icon button to a larger pill-shaped text button
- Added padding to make it wider and taller
- Changed border-radius to create a pill shape
- Added box-shadow for depth
- Added hover effects to lift the button and increase shadow

#### Completion Message

```css
.completion-message {
    margin-top: 20px;
    font-size: 1.1rem;
    color: #a0c4ff;
    font-weight: bold;
}
```

This new CSS class styles the completion message:
- Adds margin to separate it from the progress bar
- Uses a larger font size for visibility
- Uses a light blue color to match the theme
- Makes the text bold for emphasis

### Data Flow

1. When the application starts:
   - The UI is initialized with the new acrylic-styled interface
   - JavaScript simulates the collection process with progress updates
   - The backend collects actual system information in the background
2. During collection:
   - The spinner animation indicates ongoing work
   - Progress text and bar show completion percentage
   - Status text describes the current collection step
3. After collection completes:
   - The spinner disappears
   - A completion message appears
   - The close button begins animating with color changes
   - The user is prompted to close the window and return to the web dashboard

## Development Workflow

### Backend Development

1. Go code is written in `.go` files
2. Dependencies are managed with Go modules (go.mod/go.sum)
3. The application can be built and run directly with Go tools

### Frontend Development

1. The frontend is a standard web application using modern JavaScript
2. Wails automatically generates JavaScript bindings for Go functions (in wailsjs/)
3. CSS is used for styling with a responsive design

### Integration

Wails integrates the frontend and backend by:
1. Embedding the frontend build output in the Go binary
2. Providing a development mode where the frontend runs separately
3. Generating JavaScript bindings for calling Go functions from the frontend

## Deployment

The application can be distributed as a single executable file that:
1. Starts the backend Go application
2. Embeds the frontend assets
3. Runs an HTTP server for API access
4. Displays the frontend in a frameless desktop window with acrylic styling

Users can access the dashboard through a web browser by navigating to `http://localhost:9999` while the application is running.
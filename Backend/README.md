# System Information Collector

## About

This application collects detailed system information from Windows machines and displays it in a web-based dashboard.

## Features

- Collects comprehensive system information including:
  - CPU details (model, cores, threads, frequency, usage)
  - Memory (RAM) information
  - Storage details (disks, partitions, capacity, free space)
  - GPU information
  - Network interfaces and bandwidth usage
  - OS details (name, version, architecture, uptime)
  - Battery status
  - Sensor data (temperature)
  - Running processes
  
- Provides a frameless collection interface that shows the progress of information gathering
- Exposes all collected data through a REST API
- Accessible through a web browser at http://localhost:9999

## Live Development

To run in live development mode, run `wails dev` in the project directory. This will run a Vite development
server that will provide very fast hot reload of your frontend changes. If you want to develop in a browser
and have access to your Go methods, there is also a dev server that runs on http://localhost:34115. Connect
to this in your browser, and you can call your Go code from devtools.

You can also use the `dev.bat` script to start the development server.

## Building

To build a redistributable, production mode package, use `wails build`.

You can also use the `build.bat` script to build the application.

## Usage

1. Run the application
2. The frameless collection interface will appear, showing the progress of information gathering
3. Close the collection interface window
4. Open a web browser and navigate to http://localhost:9999 to view the dashboard with all collected system information

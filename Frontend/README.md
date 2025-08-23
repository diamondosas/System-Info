# System Information Dashboard - Frontend

This is the frontend component of the System Information Dashboard. It's a web-based interface that displays system information collected by the Go backend service.

## Prerequisites

- Python 3.6 or higher
- Required Python packages:
  - Flask
  - Flask-CORS
  - Requests

Install the required packages with:
```bash
pip install flask flask-cors requests
```

## Running the Server

To run the frontend server, execute:
```bash
python server.py
```

By default, the server will:
- Run on all network interfaces (0.0.0.0) on port 5000
- Proxy requests to the Go backend at http://localhost:9999/api/specs

## Network Configuration

When accessing the dashboard from another machine:

1. **Find the IP address** of the machine running the Flask server:
   - On Windows: `ipconfig`
   - On macOS/Linux: `ifconfig` or `ip addr`

2. **Set the Go backend URL** if it's running on a different machine:
   ```bash
   # On Windows (Command Prompt)
   set GO_BACKEND_URL=http://192.168.1.100:9999/api/specs && python server.py
   
   # On Windows (PowerShell)
   $env:GO_BACKEND_URL="http://192.168.1.100:9999/api/specs"; python server.py
   
   # On macOS/Linux
   GO_BACKEND_URL=http://192.168.1.100:9999/api/specs python server.py
   ```
   Replace `192.168.1.100` with the actual IP address of the machine running the Go backend.

3. **Access the dashboard** from any device on the same network by navigating to:
   ```
   http://YOUR_SERVER_IP:5000
   ```
   Replace `YOUR_SERVER_IP` with the actual IP address of the machine running the Flask server.

## Troubleshooting

If you can't access the dashboard from another machine:

1. Check that the server is running and accessible:
   ```bash
   # Test if the server is responding
   curl http://YOUR_SERVER_IP:5000/health
   ```

2. Ensure that your firewall allows connections on port 5000.

3. Verify that the Go backend is running and accessible from the Flask server machine:
   ```bash
   # Test if the Go backend is responding
   curl http://localhost:9999/api/specs
   ```

4. Check the Flask server console for any error messages.
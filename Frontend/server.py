# This Python script runs a simple Flask web server.
# Its purpose is to serve the frontend (index.html, css, js) and
# to act as a proxy for the Go backend API.

from flask import Flask, jsonify, send_from_directory, request
import requests # You may need to install this: pip install requests
import json

# --- Configuration ---
# You can change these settings if needed.
PORT = 5000
GO_BACKEND_URL = "http://localhost:9999/api/specs"

app = Flask(__name__, static_folder='.', static_url_path='')

# --- API Route ---
# This route fetches data from the Go backend and sends it to the frontend.
@app.route('/api/data')
def get_system_data():
    try:
        # Make a request to the Go backend
        response = requests.get(GO_BACKEND_URL, timeout=5) # 5 second timeout
        response.raise_for_status()  # Raise an exception for bad status codes (4xx or 5xx)
        
        # Return the JSON data from the Go backend
        return jsonify(response.json())

    except requests.exceptions.RequestException as e:
        # If the Go backend is not available or there's an error, 
        # return an error message.
        print(f"Error connecting to Go backend: {e}")
        error_message = {
            "error": "Could not connect to the Go backend.",
            "details": str(e)
        }
        return jsonify(error_message), 503 # Service Unavailable

# --- Frontend Route ---
# This route serves the main index.html file.
@app.route('/')
def serve_index():
    return send_from_directory('.', 'index.html')

# --- Health Check Route ---
@app.route('/health')
def health_check():
    return jsonify({"status": "ok"}), 200

if __name__ == '__main__':
    print(f"Frontend server running at http://localhost:{PORT}")
    print(f"Proxying API requests to: {GO_BACKEND_URL}")
    # To run this server, execute `python server.py` in your terminal.
    app.run(port=PORT, debug=True)

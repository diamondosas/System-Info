import { updateText, updateProgressBar, formatBytes } from './utils.js';

export function updateConnectionStatus(status, message) {
    const statusText = document.getElementById('status-text');
    const statusIcon = document.querySelector('.status-icon');
    if(statusText) statusText.textContent = message;

    if(statusIcon) {
        statusIcon.classList.remove('connected', 'disconnected', 'connecting');
        statusIcon.classList.add(status);
    }
}

export function updateAllInfo(data) {
    updateConnectionStatus('connected', 'Connected to backend');

    updateText('os-name', data.os?.name || 'N/A');
    updateText('os-version', data.os?.version || 'N/A');
    updateText('os-arch', data.os?.architecture || 'N/A');
    updateText('os-uptime', data.os?.uptime || 'N/A');
    updateText('os-boottime', data.os?.boot_time || 'N/A');

    updateText('cpu-model', data.cpu?.model || 'N/A');
    updateText('cpu-cores', data.cpu?.cores || 'N/A');
    updateText('cpu-threads', data.cpu?.threads || 'N/A');
    updateText('cpu-frequency', data.cpu?.frequency ? `${data.cpu.frequency} MHz` : 'N/A');
    const cpuUsage = Math.round(data.cpu?.usage_percent || 0);
    updateProgressBar('cpu-usage', cpuUsage);

    updateText('mem-total', data.memory?.total ? formatBytes(data.memory.total * 1024 * 1024) : 'N/A');
    updateText('mem-available', data.memory?.available ? formatBytes(data.memory.available * 1024 * 1024) : 'N/A');
    updateText('mem-used', data.memory?.used ? formatBytes(data.memory.used * 1024 * 1024) : 'N/A');

    updateText('swap-total', data.memory?.swap_total ? formatBytes(data.memory.swap_total * 1024 * 1024) : 'N/A');
    updateText('swap-used', data.memory?.swap_used ? formatBytes(data.memory.swap_used * 1024 * 1024) : 'N/A');

    const memoryUsage = data.memory?.total && data.memory?.used ?
        Math.round((data.memory.used / data.memory.total) * 100) : 0;
    updateProgressBar('mem-usage', memoryUsage);

    const swapUsage = data.memory?.swap_total && data.memory?.swap_used ?
        Math.round((data.memory.swap_used / data.memory.swap_total) * 100) : 0;
    updateProgressBar('swap-usage', swapUsage);

    updateText('gpu-vendor', data.gpu?.vendor || 'N/A');
    updateText('gpu-model', data.gpu?.model || 'N/A');
    updateText('gpu-vram', data.gpu?.vram ? formatBytes(data.gpu.vram * 1024 * 1024) : 'N/A');

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

    const networkContainer = document.getElementById('network-interfaces');
    if (networkContainer) {
        networkContainer.innerHTML = '';
        if (data.network && Array.isArray(data.network)) {
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

            if (filteredNetworks.length === 0) {
                networkContainer.innerHTML = '<p>No Ethernet or Wi-Fi interfaces found</p>';
            }
        } else {
            networkContainer.innerHTML = '<p>No network interfaces found</p>';
        }
    }

    updateText('battery-status', data.battery?.status || 'N/A');
    updateText('battery-percentage', data.battery?.percentage ? `${data.battery.percentage}%` : 'N/A');

    updateText('cpu-temp', data.sensors?.cpu_temp ? `${data.sensors.cpu_temp}°C` : 'N/A');
    updateText('gpu-temp', data.sensors?.gpu_temp ? `${data.sensors.gpu_temp}°C` : 'N/A');

    const processesBody = document.getElementById('processes-body');
    if (processesBody) {
        processesBody.innerHTML = '';
        if (data.processes && Array.isArray(data.processes)) {
            const sortedProcesses = [...data.processes].sort((a, b) =>
                (b.memory_mb || 0) - (a.memory_mb || 0)
            );

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

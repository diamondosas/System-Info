export function updateText(id, value) {
    const element = document.getElementById(id);
    if (element) {
        element.textContent = value;
    }
}

export function updateProgressBar(id, percentage) {
    const element = document.getElementById(id);
    if (element) {
        element.style.width = `${percentage}%`;
        const textElement = document.getElementById(`${id}-text`);
        if(textElement) textElement.textContent = `${percentage}%`;
    }
}

export function formatBytes(bytes, decimals = 2) {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const dm = decimals < 0 ? 0 : decimals;
    const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
}

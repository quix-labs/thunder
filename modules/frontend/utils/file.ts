export function downloadJSON(obj: any, filename: string = 'data.json'): void {
    const jsonStr = JSON.stringify(obj, null, 2);
    const blob = new Blob([jsonStr], { type: 'application/json' });
    const file = new File([blob], filename, { type: 'application/json' });
    const url = URL.createObjectURL(file);

    const a = document.createElement('a');
    a.href = url;
    a.download = filename;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
}
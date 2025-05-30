export function getMinDateTime() {
    const now = new Date();
    now.setHours(now.getHours() + 24);
    // Format as "YYYY-MM-DDTHH:MM" for input[type="datetime-local"]
    return now.toISOString().slice(0, 16);
}

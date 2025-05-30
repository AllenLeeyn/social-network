export function formatDate(dateString, options = {}) {
  
  if (!dateString) return '';
  const date = new Date(dateString);

  // Default options: e.g., "May 30, 2025, 5:22 PM"
  const defaultOptions = {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  };

  return date.toLocaleString(undefined, { ...defaultOptions, ...options });
}

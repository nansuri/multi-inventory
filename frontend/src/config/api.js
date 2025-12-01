// API Configuration
// Uses VITE_API_BASE_URL from environment variables, or relative path for production
const apiBase = import.meta.env.VITE_API_BASE_URL || '';

export { apiBase };

// config.js
// Export the API base URL from Vite's environment variables
// Fallback to localhost if not defined
export const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8080/api/v1";

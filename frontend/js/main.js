const API_URL = "/api";

async function fetchAPI(endpoint, options = {}) {
  const token = localStorage.getItem("token");
  const headers = { "Content-Type": "application/json" };
  if (token) headers["Authorization"] = `Bearer ${token}`;

  const response = await fetch(`${API_URL}${endpoint}`, {
    ...options,
    headers,
  });
  if (!response.ok) throw new Error("API Error");
  return await response.json();
}

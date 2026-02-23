// fetchUsers.tsx
import { API_BASE_URL } from './../config/config';

export interface User {
  id: string;
  email: string;
  username: string;
  domain: string;
  password: string;
}

/**
 * Fetch users from the backend. By default, fetch 4000 at a time.
 */
export async function fetchUsers(
  email?: string,
  domain?: string,
  limit: number = 4000
): Promise<User[]> {
  const params = new URLSearchParams();

  if (email) {
    params.append('email', email);
  }
  if (domain) {
    params.append('domain', domain);
  }

  // Add limit param
  params.append('limit', limit.toString());

  const response = await fetch(`${API_BASE_URL}?${params.toString()}`);
  if (!response.ok) {
    throw new Error('Failed to fetch users');
  }

  const data = await response.json();
  return data;
}

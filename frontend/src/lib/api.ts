const API_BASE_URL = (
  import.meta.env.PUBLIC_API_BASE_URL ?? 'http://localhost:8080'
).replace(/\/$/, '');

export interface HoleScore {
  hole_number: number;
  par: number;
  strokes: number;
}

export interface User {
  id: number;
  username: string;
}

// Data the browser sends when creating a round.
export interface RoundInput {
  // Must match an existing user; the backend resolves the user's database ID.
  username: string;
  course_name: string;
  date: string; // RFC 3339 timestamp including timezone.
  holes: HoleScore[];
}

// Data the backend returns when retrieving a round.
export interface Round extends RoundInput {
  id: number;
}

// Shared response handling for all API calls.
async function readResponse<T>(response: Response): Promise<T> {
  if (!response.ok) {
    const body = await response.json().catch(() => null);

    throw new Error(
      typeof body?.error === 'string'
        ? body.error
        : `Request failed (${response.status})`
    );
  }

  return response.json() as Promise<T>;
}

export async function checkAPI(): Promise<string> {
  const response = await fetch(`${API_BASE_URL}/healthz`);
  const data = await readResponse<{ message: string }>(response);
  return data.message;
}

export async function saveRound(
  round: RoundInput
): Promise<{ id: number; message: string }> {
  const response = await fetch(`${API_BASE_URL}/api/v1/rounds`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(round),
  });

  return readResponse<{ id: number; message: string }>(response);
}

export async function getRound(id: number): Promise<Round> {
  const response = await fetch(`${API_BASE_URL}/api/v1/rounds/${id}`);
  return readResponse<Round>(response);
}

export async function getUserRounds(username: string): Promise<Round[]> {
  const query = new URLSearchParams({ username });
  const response = await fetch(`${API_BASE_URL}/api/v1/rounds?${query}`);
  return readResponse<Round[]>(response);
}

export async function getAllUsers(): Promise<User[]> {
  const response = await fetch(`${API_BASE_URL}/api/v1/users`);
  return readResponse<User[]>(response);
}

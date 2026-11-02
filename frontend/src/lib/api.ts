const API_BASE_URL = (
  import.meta.env.PUBLIC_API_BASE_URL ?? 'http://localhost:8080'
).replace(/\/$/, '');

export interface HoleScore {
  hole_number: number;
  par: number;
  strokes: number;
}

// Data the browser sends when creating a round.
export interface RoundInput {
  course_name: string;
  date: string;
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
): Promise<{ message: string }> {
  const response = await fetch(`${API_BASE_URL}/api/v1/rounds`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(round),
  });

  return readResponse<{ message: string }>(response);
}

export async function getRound(id: number): Promise<Round> {
  const response = await fetch(`${API_BASE_URL}/api/v1/rounds/${id}`);
  return readResponse<Round>(response);
}
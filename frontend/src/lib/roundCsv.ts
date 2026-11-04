import type { HoleScore, RoundInput } from './api';

// Expected format: header, Par row, then one player row.
export function parseRoundCsv(text: string): RoundInput {
  const lines = text.trim().split('\n');
  if (lines.length !== 3) {
    throw new Error('Choose a scorecard CSV containing a header, a Par row and one player row.');
  }

  const headers = lines[0].split(',').map(value => value.trim());
  const parValues = lines[1].split(',').map(value => value.trim());
  const playerValues = lines[2].split(',').map(value => value.trim());

  if (parValues[0] !== 'Par') {
    throw new Error('The second CSV row must contain the par values.');
  }

  const usernameColumn = headers.indexOf('PlayerName');
  const username = playerValues[usernameColumn];
  const courseColumn = headers.indexOf('CourseName');
  const dateColumn = headers.indexOf('StartDate');
  const courseName = playerValues[courseColumn];
  const startDate = playerValues[dateColumn] ?? '';
  // CSV: 2026-10-10 1542+0200 → API: 2026-10-10T15:42:00+02:00
  const [day, time] = startDate.split(' ');
  if (!day || !/^\d{4}[+-]\d{4}$/.test(time ?? '')) {
    throw new Error('StartDate must use YYYY-MM-DD HHmm+HHmm, for example 2026-10-10 1542+0200.');
  }
  const hours = time.slice(0, 2);
  const minutes = time.slice(2, 4);
  const offset = time.slice(4, 7) + ':' + time.slice(7, 9);
  const date = `${day}T${hours}:${minutes}:00${offset}`;

  const holes: HoleScore[] = [];

  for (let column = 0; column < headers.length; column++) {
    const columnName = headers[column];
    if (!/^Hole\d+$/.test(columnName)) continue;

    // For example, Hole1 becomes hole number 1.
    const holeNumber = Number(columnName.replace('Hole', ''));
    const par = Number(parValues[column]);
    const strokes = Number(playerValues[column]);

    holes.push({
      hole_number: holeNumber,
      par,
      strokes,
    });
  }

  return {
    username,
    course_name: courseName,
    date,
    holes,
  };
}

import type { RoundInput } from './api';

// Expected format: header, Par row, then one player row.
export function parseRoundCsv(text: string): RoundInput {
  const rows = text.trim().split(/\r?\n/).map(line => line.split(',').map(value => value.trim()));
  const [header, par, player] = rows;
  if (rows.length !== 3 || par?.[0] !== 'Par') {
    throw new Error('Choose a scorecard CSV containing a header, a Par row and one player row.');
  }
  const column = (name: string) => header.indexOf(name);

  return {
    course_name: player[column('CourseName')],
    date: (player[column('StartDate')] ?? '').slice(0, 10),
    holes: header.flatMap((name, index) => /^Hole\d+$/.test(name) ? [{
      hole_number: Number(name.slice(4)),
      par: Number(par[index]),
      strokes: Number(player[index]),
    }] : []),
  };
}

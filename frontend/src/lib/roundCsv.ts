import type { RoundInput } from './api';

// Parse CSV records, including quoted commas, escaped quotes and line breaks.
function records(text: string): string[][] {
  const rows: string[][] = [];
  let row: string[] = [];
  let field = '';
  let quoted = false;
  let closed = false;
  const finishField = () => { row.push(field); field = ''; closed = false; };
  const finishRow = () => {
    finishField();
    if (row.some(value => value.trim() !== '')) rows.push(row);
    row = [];
  };
  text = text.replace(/^\uFEFF/, '');
  for (let i = 0; i < text.length; i++) {
    const char = text[i];
    if (quoted) {
      if (char === '"') {
        if (text[i + 1] === '"') { field += '"'; i++; }
        else { quoted = false; closed = true; }
      } else field += char;
    } else if (char === ',') finishField();
    else if (char === '\n' || char === '\r') {
      if (char === '\r' && text[i + 1] === '\n') i++;
      finishRow();
    } else if (char === '"') {
      if (field !== '' || closed) throw new Error('Invalid quote in CSV field.');
      quoted = true;
    } else {
      if (closed) throw new Error('Unexpected text after a quoted CSV field.');
      field += char;
    }
  }
  if (quoted) throw new Error('Unclosed quoted CSV field.');
  finishRow();
  return rows;
}

export function parseRoundCsv(text: string): RoundInput[] {
  const [header, ...rows] = records(text);
  if (!header || rows.length === 0) throw new Error('CSV must contain a header and at least one hole.');
  const columns = header.map(value => value.trim());
  if (new Set(columns).size !== columns.length) throw new Error('CSV contains duplicate column names.');
  if (columns.includes('PlayerName')) return parseScorecard(columns, rows);
  for (const name of ['course_name', 'date', 'hole_number', 'par', 'strokes']) {
    if (!columns.includes(name)) throw new Error(`Missing CSV column: ${name}.`);
  }
  const rounds = new Map<string, RoundInput>();
  rows.forEach((row, index) => {
    const fail = (message: string): never => { throw new Error(`CSV record ${index + 2}: ${message}`); };
    if (row.length !== columns.length) fail('column count does not match the header.');
    const value = (name: string) => (row[columns.indexOf(name)] ?? '').trim();
    const course = value('course_name');
    const date = value('date');
    if (!course) fail('course_name is required.');
    if (!/^\d{4}-\d{2}-\d{2}$/.test(date) || date.startsWith('0000') ||
        !Number.isFinite(Date.parse(date)) || new Date(date).toISOString().slice(0, 10) !== date) {
      fail('date must be a valid date in YYYY-MM-DD format.');
    }
    const positiveInteger = (name: string) => {
      const raw = value(name);
      const number = Number(raw);
      if (!/^\d+$/.test(raw) || !Number.isSafeInteger(number) || number <= 0 || number > 2147483647) {
        fail(`${name} must be a positive integer no greater than 2147483647.`);
      }
      return number;
    };
    const hole = {
      hole_number: positiveInteger('hole_number'),
      par: positiveInteger('par'),
      strokes: positiveInteger('strokes'),
    };
    const key = JSON.stringify([value('round_key'), course, date]);
    let round = rounds.get(key);
    if (!round) {
      round = { course_name: course, date, holes: [] };
      rounds.set(key, round);
    }
    if (round.holes.some(existing => existing.hole_number === hole.hole_number)) {
      fail(`duplicate hole ${hole.hole_number}. Use round_key to separate rounds on the same course and date.`);
    }
    round.holes.push(hole);
  });
  return [...rounds.values()];
}

// Scorecard exports have one Par row and one player row per round.
function parseScorecard(columns: string[], rows: string[][]): RoundInput[] {
  for (const name of ['CourseName', 'LayoutName', 'StartDate']) {
    if (!columns.includes(name)) throw new Error(`Missing CSV column: ${name}.`);
  }
  const holes = columns.filter(name => /^Hole[1-9]\d*$/.test(name));
  if (holes.length === 0) throw new Error('CSV must contain Hole1, Hole2, … columns.');
  const groups = new Map<string, { par?: string[]; players: string[][] }>();
  const value = (row: string[], name: string) => (row[columns.indexOf(name)] ?? '').trim();
  for (const [index, row] of rows.entries()) {
    if (row.length !== columns.length) throw new Error(`CSV record ${index + 2}: column count does not match the header.`);
    const player = value(row, 'PlayerName');
    if (!player) throw new Error(`CSV record ${index + 2}: PlayerName is required.`);
    const key = JSON.stringify(['CourseName', 'LayoutName', 'StartDate', 'EndDate'].map(name => value(row, name)));
    let group = groups.get(key);
    if (!group) { group = { players: [] }; groups.set(key, group); }
    if (player === 'Par') {
      if (group.par) throw new Error('Multiple Par rows found for the same scorecard.');
      group.par = row;
    } else group.players.push(row);
  }
  const rounds: RoundInput[] = [];
  for (const group of groups.values()) {
    if (!group.par) throw new Error('Each scorecard needs a Par row.');
    if (group.players.length !== 1) {
      throw new Error('Each scorecard must contain exactly one player row. Keep the Par row and the player you want to import.');
    }
    const player = group.players[0];
    const start = value(player, 'StartDate');
    if (!/^\d{4}-\d{2}-\d{2}(?:$|[ T])/.test(start)) throw new Error('StartDate must begin with YYYY-MM-DD.');
    // Reuse the normal round validation. CSV summary totals are not authoritative.
    const encode = (field: string) => '"' + field.replaceAll('"', '""') + '"';
    const normalized = ['course_name,date,hole_number,par,strokes', ...holes.map(name => [
      value(player, 'CourseName'), start.slice(0, 10), name.slice(4),
      value(group.par!, name), value(player, name),
    ].map(encode).join(','))].join('\n');
    rounds.push(...parseRoundCsv(normalized));
  }
  return rounds;
}

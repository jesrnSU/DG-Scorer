import { test } from 'node:test';
import assert from 'node:assert/strict';
import { readFile } from 'node:fs/promises';
import { parseRoundCsv } from './roundCsv.ts';
const header = 'course_name,date,hole_number,par,strokes\n';

test('example groups holes into rounds and parses quoted commas', async () => {
  const csv = await readFile(new URL('../../public/examples/rounds.csv', import.meta.url), 'utf8');
  const rounds = parseRoundCsv(csv);
  assert.equal(rounds.length, 2);
  assert.equal(rounds[0].holes.length, 2);
  assert.equal(rounds[1].course_name, 'Practice, North');
  assert.equal('round_key' in rounds[0], false);
});
test('handles BOM, CRLF, escaped quotes, embedded newlines and blank rows', () => {
  const rounds = parseRoundCsv('\uFEFF' + header.replace('\n', '\r\n') +
    '"Course ""North"",\nSweden",2026-11-04,1,3,4\r\n\r\n');
  assert.equal(rounds[0].course_name, 'Course "North",\nSweden');
});
test('round_key separates rounds at the same course on the same day', () => {
  assert.equal(parseRoundCsv('round_key,' + header +
    'a,Practice,2026-11-04,1,3,4\nb,Practice,2026-11-04,1,3,2').length, 2);
});
for (const [name, csv, message] of [
  ['empty file', '', /header/],
  ['missing column', 'date\n2026-11-04', /Missing CSV column/],
  ['duplicate columns', header.trim() + ',par\nPractice,2026-11-04,1,3,4,3', /duplicate column/],
  ['wrong column count', header + 'Practice,2026-11-04,1,3', /column count/],
  ['invalid date', header + 'Practice,2026-02-30,1,3,4', /valid date/],
  ['empty course', header + ',2026-11-04,1,3,4', /course_name/],
  ['zero score', header + 'Practice,2026-11-04,1,3,0', /strokes/],
  ['fractional score', header + 'Practice,2026-11-04,1,3,2.5', /strokes/],
  ['integer overflow', header + 'Practice,2026-11-04,1,3,2147483648', /strokes/],
  ['duplicate hole', header + 'Practice,2026-11-04,1,3,4\nPractice,2026-11-04,1,3,3', /duplicate hole/],
  ['unclosed quote', header + '"Practice,2026-11-04,1,3,4', /Unclosed/],
  ['text after quote', header + '"Practice"x,2026-11-04,1,3,4', /Unexpected text/],
]) test(`rejects ${name}`, () => assert.throws(() => parseRoundCsv(csv), message));

test('imports supplied scorecard using Par and player rows', async () => {
  const csv = await readFile(new URL('../../public/examples/scorecard.csv', import.meta.url), 'utf8');
  const [round] = parseRoundCsv(csv);
  assert.equal(round.course_name, 'Visättra');
  assert.equal(round.date, '2026-10-10');
  assert.equal(round.holes.length, 18);
  assert.deepEqual(round.holes[0], { hole_number: 1, par: 3, strokes: 4 });
  assert.equal(round.holes.reduce((sum, hole) => sum + hole.strokes, 0), 57);
  assert.equal(round.holes.reduce((sum, hole) => sum + hole.par, 0), 60);
  const lines = csv.trim().split('\n');
  assert.throws(() => parseRoundCsv([lines[0], lines[2]].join('\n')), /Par row/);
  assert.throws(() => parseRoundCsv(csv.trim() + '\n' + lines[2]), /exactly one player/);
  assert.throws(() => parseRoundCsv(csv.replace('4,3,2,4,3,5', '0,3,2,4,3,5')), /strokes/);
});

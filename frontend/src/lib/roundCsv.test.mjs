import { test } from 'node:test';
import assert from 'node:assert/strict';
import { readFile } from 'node:fs/promises';
import { parseRoundCsv } from './roundCsv.ts';

test('converts the supplied scorecard into the round expected by the API', async () => {
  const csv = await readFile(new URL('../../public/examples/scorecard.csv', import.meta.url), 'utf8');
  const round = parseRoundCsv(csv);
  assert.equal(round.course_name, 'Visättra');
  assert.equal(round.username, 'JesperNi');
  assert.equal(round.date, '2026-10-10T15:42:00+02:00');
  const laterRound = parseRoundCsv(csv.replaceAll('1542+0200', '1742+0200'));
  assert.notEqual(round.date, laterRound.date);
  assert.equal(round.holes.length, 18);
  assert.deepEqual(round.holes[0], { hole_number: 1, par: 3, strokes: 4 });
  assert.equal(round.holes.reduce((sum, hole) => sum + hole.strokes, 0), 57);
  assert.equal(round.holes.reduce((sum, hole) => sum + hole.par, 0), 60);
});

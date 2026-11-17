package com.jesper.dgscorekeeper.ui

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.size
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.semantics.contentDescription
import androidx.compose.ui.semantics.semantics
import androidx.compose.ui.unit.dp
import androidx.wear.compose.material.Button
import androidx.wear.compose.material.MaterialTheme
import androidx.wear.compose.material.Text
import com.jesper.dgscorekeeper.R
import com.jesper.dgscorekeeper.round.HoleInfo
import com.jesper.dgscorekeeper.round.MIN_SCORE

@Composable
fun HoleScreen(
    holeNumber: Int,
    scores: List<Int>,
    hole: HoleInfo,
    onIncrease: (Int) -> Unit,
    onDecrease: (Int) -> Unit,
    onNext: () -> Unit,
    onPrevious: () -> Unit,
) {
    ScreenLayout(
        title = stringResource(R.string.compact_hole_details, holeNumber, hole.lengthMeters, hole.par),
        onNext = onNext,
        onPrevious = onPrevious,
        scrollable = false,
        compact = true,
    ) {
        Column(
            modifier = Modifier.weight(1f).fillMaxWidth(),
            verticalArrangement = Arrangement.SpaceEvenly,
        ) {
            scores.chunked(2).forEachIndexed { row, players ->
                Row(
                    modifier = Modifier.fillMaxWidth().weight(1f),
                    horizontalArrangement = Arrangement.spacedBy(4.dp),
                    verticalAlignment = Alignment.CenterVertically,
                ) {
                    players.forEachIndexed { column, score ->
                        val player = row * 2 + column
                        PlayerScore(
                            player = player,
                            score = score,
                            onIncrease = { onIncrease(player) },
                            onDecrease = { onDecrease(player) },
                            modifier = Modifier.weight(1f),
                        )
                    }
                    if (players.size == 1 && scores.size > 1) Spacer(Modifier.weight(1f))
                }
            }
        }
    }
}

@Composable
private fun PlayerScore(
    player: Int,
    score: Int,
    onIncrease: () -> Unit,
    onDecrease: () -> Unit,
    modifier: Modifier = Modifier,
) {
    val decreaseLabel = stringResource(R.string.decrease_player_score, player + 1)
    val increaseLabel = stringResource(R.string.increase_player_score, player + 1)
    Column(
        modifier = modifier,
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.spacedBy(2.dp, Alignment.CenterVertically),
    ) {
        Text(stringResource(R.string.player_number, player + 1), style = MaterialTheme.typography.caption2)
        Row(
            modifier = Modifier.fillMaxWidth(),
            horizontalArrangement = Arrangement.SpaceEvenly,
            verticalAlignment = Alignment.CenterVertically,
        ) {
            Button(
                modifier = Modifier.size(28.dp).semantics { contentDescription = decreaseLabel },
                enabled = score > MIN_SCORE,
                onClick = onDecrease,
            ) { Text(stringResource(R.string.decrease_score)) }
            Text(score.toString(), style = MaterialTheme.typography.title3, maxLines = 1)
            Button(
                modifier = Modifier.size(28.dp).semantics { contentDescription = increaseLabel },
                onClick = onIncrease,
            ) { Text(stringResource(R.string.increase_score)) }
        }
    }
}

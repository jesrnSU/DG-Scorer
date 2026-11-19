package com.jesper.dgscorekeeper.ui

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Row
import androidx.compose.runtime.Composable
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.unit.dp
import androidx.wear.compose.material.Button
import androidx.wear.compose.material.Chip
import androidx.wear.compose.material.Text
import com.jesper.dgscorekeeper.R
import com.jesper.dgscorekeeper.round.RoundUiState

@Composable
fun RoundSetupScreen(
    state: RoundUiState,
    onPlayerCount: (Int) -> Unit,
    onStart: () -> Unit,
    onPrevious: () -> Unit,
) {
    ScreenLayout(title = stringResource(R.string.round_setup), onPrevious = onPrevious) {
        Text(state.course.name)
        Text(stringResource(R.string.sample_course))
        Text(stringResource(R.string.players, state.playerCount))
        Row(horizontalArrangement = Arrangement.spacedBy(16.dp)) {
            Button(enabled = state.playerCount > 1, onClick = { onPlayerCount(state.playerCount - 1) }) {
                Text(stringResource(R.string.decrease_score))
            }
            Button(enabled = state.playerCount < 4, onClick = { onPlayerCount(state.playerCount + 1) }) {
                Text(stringResource(R.string.increase_score))
            }
        }
        Chip(onClick = onStart, label = { Text(stringResource(R.string.start_round)) })
    }
}

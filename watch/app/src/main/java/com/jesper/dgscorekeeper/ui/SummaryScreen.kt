package com.jesper.dgscorekeeper.ui

import androidx.compose.runtime.Composable
import androidx.compose.ui.res.stringResource
import androidx.wear.compose.material.MaterialTheme
import androidx.wear.compose.material.Text
import com.jesper.dgscorekeeper.R

@Composable
fun SummaryScreen(
    totals: List<Int>,
    onNext: () -> Unit,
    onPrevious: () -> Unit,
) {
    ScreenLayout(
        title = stringResource(R.string.total_score),
        onNext = onNext,
        onPrevious = onPrevious,
    ) {
        totals.forEachIndexed { player, total ->
            Text(
                text = stringResource(R.string.player_total, player + 1, total),
                style = MaterialTheme.typography.title2,
                color = MaterialTheme.colors.onBackground,
            )
        }
    }
}

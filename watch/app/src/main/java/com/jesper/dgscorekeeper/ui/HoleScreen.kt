package com.jesper.dgscorekeeper.ui

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Row
import androidx.compose.runtime.Composable
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.unit.dp
import androidx.wear.compose.material.Button
import androidx.wear.compose.material.MaterialTheme
import androidx.wear.compose.material.Text
import com.jesper.dgscorekeeper.R
import com.jesper.dgscorekeeper.round.MIN_SCORE

@Composable
fun HoleScreen(
    holeNumber: Int,
    score: Int,
    onIncrease: () -> Unit,
    onDecrease: () -> Unit,
    onNext: () -> Unit,
    onPrevious: () -> Unit,
) {
    ScreenLayout(
        title = stringResource(R.string.hole_number, holeNumber),
        onNext = onNext,
        onPrevious = onPrevious,
    ) {
        Text(
            text = score.toString(),
            style = MaterialTheme.typography.display1,
            color = MaterialTheme.colors.onBackground,
        )
        Row(
            horizontalArrangement = Arrangement.spacedBy(16.dp),
        ) {
            Button(
                enabled = score > MIN_SCORE,
                onClick = onDecrease,
            ) {
                Text(
                    text = stringResource(R.string.decrease_score),
                )
            }
            Button(onClick = onIncrease) {
                Text(
                    text = stringResource(R.string.increase_score),
                )
            }
        }
    }
}

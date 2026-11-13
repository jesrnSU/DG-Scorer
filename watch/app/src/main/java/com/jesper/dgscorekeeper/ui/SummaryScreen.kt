package com.jesper.dgscorekeeper.ui

import androidx.compose.runtime.Composable
import androidx.compose.ui.res.stringResource
import androidx.wear.compose.material.MaterialTheme
import androidx.wear.compose.material.Text
import com.jesper.dgscorekeeper.R

@Composable
fun SummaryScreen(totalScore: Int, onNext: () -> Unit, onPrevious: () -> Unit) {
    ScreenLayout(
        title = stringResource(R.string.total_score),
        onNext = onNext,
        onPrevious = onPrevious,
    ) {
        Text(
            text = totalScore.toString(),
            style = MaterialTheme.typography.display1,
            color = MaterialTheme.colors.onBackground,
        )
    }
}

package com.jesper.dgscorekeeper.ui

import androidx.compose.runtime.Composable
import com.jesper.dgscorekeeper.round.RoundScreen
import com.jesper.dgscorekeeper.round.RoundUiState

@Composable
fun ScorekeeperApp(
    state: RoundUiState,
    onStart: () -> Unit,
    onIncrease: () -> Unit,
    onDecrease: () -> Unit,
    onNext: () -> Unit,
    onPrevious: () -> Unit,
) {
    when (val screen = state.screen) {
        RoundScreen.Menu -> MenuScreen(onStart = onStart)
        is RoundScreen.Hole -> HoleScreen(
            holeNumber = screen.index + 1,
            score = state.scores[screen.index],
            onIncrease = onIncrease,
            onDecrease = onDecrease,
            onNext = onNext,
            onPrevious = onPrevious,
        )
        RoundScreen.Summary -> SummaryScreen(
            totalScore = state.totalScore,
            onNext = onNext,
            onPrevious = onPrevious,
        )
    }
}

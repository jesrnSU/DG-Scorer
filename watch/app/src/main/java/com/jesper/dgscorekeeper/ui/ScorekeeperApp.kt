package com.jesper.dgscorekeeper.ui

import androidx.compose.runtime.Composable
import com.jesper.dgscorekeeper.round.RoundScreen
import com.jesper.dgscorekeeper.round.RoundUiState

@Composable
fun ScorekeeperApp(
    state: RoundUiState,
    onStart: () -> Unit,
    onOpenSetup: () -> Unit,
    onPlayerCount: (Int) -> Unit,
    onIncrease: (Int) -> Unit,
    onDecrease: (Int) -> Unit,
    onNext: () -> Unit,
    onPrevious: () -> Unit,
) {
    when (val screen = state.screen) {
        RoundScreen.Menu -> {
            MenuScreen(onStart = onOpenSetup)
        }

        RoundScreen.Setup -> RoundSetupScreen(state, onPlayerCount, onStart, onPrevious)

        is RoundScreen.Hole -> {
            HoleScreen(
                holeNumber = screen.index + 1,
                scores = state.playerScores.take(state.playerCount).map { it[screen.index] },
                hole = state.course.holes[screen.index],
                onIncrease = onIncrease,
                onDecrease = onDecrease,
                onNext = onNext,
                onPrevious = onPrevious,
            )
        }
        RoundScreen.Summary -> {
            SummaryScreen(
                totals = state.totals,
                onNext = onNext,
                onPrevious = onPrevious,
            )
        }
    }
}

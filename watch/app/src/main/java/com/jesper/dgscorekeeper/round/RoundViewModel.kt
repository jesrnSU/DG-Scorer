package com.jesper.dgscorekeeper.round

import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.setValue
import androidx.lifecycle.SavedStateHandle
import androidx.lifecycle.ViewModel

class RoundViewModel(private val savedStateHandle: SavedStateHandle) : ViewModel() {
    var uiState by mutableStateOf(restoreState())
        private set

    fun startRound() {
        if (uiState.screen == RoundScreen.Menu) {
            update(uiState.copy(screen = RoundScreen.Hole(0)))
        }
    }

    fun increaseScore() = changeScore(1)

    fun decreaseScore() = changeScore(-1)

    fun next() {
        when (val screen = uiState.screen) {
            RoundScreen.Menu -> Unit
            is RoundScreen.Hole -> update(uiState.copy(
                screen = if (screen.index == HOLE_COUNT - 1) RoundScreen.Summary
                    else RoundScreen.Hole(screen.index + 1),
            ))
            RoundScreen.Summary -> update(RoundUiState())
        }
    }

    fun previous() {
        when (val screen = uiState.screen) {
            RoundScreen.Menu -> Unit
            is RoundScreen.Hole -> update(uiState.copy(
                screen = RoundScreen.Hole((screen.index - 1).coerceAtLeast(0)),
            ))
            RoundScreen.Summary -> update(uiState.copy(screen = RoundScreen.Hole(HOLE_COUNT - 1)))
        }
    }

    private fun changeScore(amount: Int) {
        val hole = uiState.screen as? RoundScreen.Hole ?: return
        update(uiState.copy(scores = uiState.scores.mapIndexed { index, score ->
            if (index == hole.index) (score + amount).coerceAtLeast(MIN_SCORE) else score
        }))
    }

    private fun update(state: RoundUiState) {
        uiState = state
        // Store Bundle-compatible values so Android can restore a destroyed activity/process.
        savedStateHandle["scores"] = state.scores.toIntArray()
        savedStateHandle["page"] = when (val screen = state.screen) {
            RoundScreen.Menu -> -1
            is RoundScreen.Hole -> screen.index
            RoundScreen.Summary -> HOLE_COUNT
        }
    }

    private fun restoreState(): RoundUiState {
        val scores = savedStateHandle.get<IntArray>("scores")?.toList()
            ?.takeIf { it.size == HOLE_COUNT && it.all { score -> score >= MIN_SCORE } }
            ?: List(HOLE_COUNT) { DEFAULT_SCORE }
        val screen = when (val page = savedStateHandle.get<Int>("page")) {
            in 0 until HOLE_COUNT -> RoundScreen.Hole(requireNotNull(page))
            HOLE_COUNT -> RoundScreen.Summary
            else -> RoundScreen.Menu
        }
        return RoundUiState(screen = screen, scores = scores)
    }
}

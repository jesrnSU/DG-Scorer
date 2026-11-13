package com.jesper.dgscorekeeper.round

import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.setValue
import androidx.lifecycle.SavedStateHandle
import androidx.lifecycle.ViewModel

class RoundViewModel(
    private val savedStateHandle: SavedStateHandle,
) : ViewModel() {
    var uiState by mutableStateOf(restoreState())
        private set

    fun startRound() {
        if (uiState.screen != RoundScreen.Menu) {
            return
        }

        val nextState = uiState.copy(screen = RoundScreen.Hole(0))
        update(nextState)
    }

    fun increaseScore() {
        changeScore(1)
    }

    fun decreaseScore() {
        changeScore(-1)
    }

    fun next() {
        val nextScreen = when (val screen = uiState.screen) {
            RoundScreen.Menu -> {
                return
            }

            is RoundScreen.Hole -> {
                if (screen.index == HOLE_COUNT - 1) {
                    RoundScreen.Summary
                } else {
                    RoundScreen.Hole(screen.index + 1)
                }
            }

            RoundScreen.Summary -> {
                update(RoundUiState())
                return
            }
        }

        val nextState = uiState.copy(screen = nextScreen)
        update(nextState)
    }

    fun previous() {
        val previousScreen = when (val screen = uiState.screen) {
            RoundScreen.Menu -> {
                return
            }

            is RoundScreen.Hole -> {
                if (screen.index == 0) {
                    RoundScreen.Menu
                } else {
                    RoundScreen.Hole(screen.index - 1)
                }
            }

            RoundScreen.Summary -> {
                RoundScreen.Hole(HOLE_COUNT - 1)
            }
        }

        val nextState = uiState.copy(screen = previousScreen)
        update(nextState)
    }

    private fun changeScore(amount: Int) {
        val screen = uiState.screen
        if (screen !is RoundScreen.Hole) {
            return
        }

        val currentScore = uiState.scores[screen.index]
        val newScore = (currentScore + amount).coerceAtLeast(MIN_SCORE)
        val updatedScores = uiState.scores.toMutableList()
        updatedScores[screen.index] = newScore

        val nextState = uiState.copy(scores = updatedScores.toList())
        update(nextState)
    }

    private fun update(state: RoundUiState) {
        uiState = state

        // Store Bundle-compatible values so Android can restore the round.
        val savedPage = when (val screen = state.screen) {
            RoundScreen.Menu -> -1
            is RoundScreen.Hole -> screen.index
            RoundScreen.Summary -> HOLE_COUNT
        }

        savedStateHandle["scores"] = state.scores.toIntArray()
        savedStateHandle["page"] = savedPage
    }

    private fun restoreState(): RoundUiState {
        val scores = restoreScores()
        val screen = restoreScreen()

        return RoundUiState(
            screen = screen,
            scores = scores,
        )
    }

    private fun restoreScores(): List<Int> {
        val savedScores = savedStateHandle.get<IntArray>("scores")
        if (savedScores == null || savedScores.size != HOLE_COUNT) {
            return RoundUiState().scores
        }

        for (score in savedScores) {
            if (score < MIN_SCORE) {
                return RoundUiState().scores
            }
        }

        return savedScores.toList()
    }

    private fun restoreScreen(): RoundScreen {
        val savedPage = savedStateHandle.get<Int>("page")
        if (savedPage == null) {
            return RoundScreen.Menu
        }

        return when (savedPage) {
            in 0 until HOLE_COUNT -> RoundScreen.Hole(savedPage)
            HOLE_COUNT -> RoundScreen.Summary
            else -> RoundScreen.Menu
        }
    }
}

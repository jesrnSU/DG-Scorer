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

    fun openSetup() {
        if (uiState.screen == RoundScreen.Menu) update(uiState.copy(screen = RoundScreen.Setup))
    }

    fun setPlayerCount(count: Int) {
        if (uiState.screen == RoundScreen.Setup) {
            update(uiState.copy(playerCount = count.coerceIn(1, 4)))
        }
    }

    fun startRound() {
        if (uiState.screen == RoundScreen.Setup) update(uiState.copy(screen = RoundScreen.Hole(0)))
    }

    fun increaseScore(player: Int = 0) = changeScore(player, 1)
    fun decreaseScore(player: Int = 0) = changeScore(player, -1)

    fun next() {
        val nextScreen = when (val screen = uiState.screen) {
            RoundScreen.Menu, RoundScreen.Setup -> {
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
            RoundScreen.Setup -> RoundScreen.Menu
            RoundScreen.Menu -> {
                return
            }

            is RoundScreen.Hole -> {
                if (screen.index == 0) {
                    RoundScreen.Setup
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

    private fun changeScore(player: Int, amount: Int) {
        val screen = uiState.screen
        if (screen !is RoundScreen.Hole || player !in 0 until uiState.playerCount) return
        val scores = uiState.playerScores.mapIndexed { index, holes ->
            if (index != player) holes else holes.mapIndexed { hole, score ->
                if (hole == screen.index) (score + amount).coerceAtLeast(MIN_SCORE) else score
            }
        }
        update(uiState.copy(playerScores = scores))
    }

    private fun update(state: RoundUiState) {
        uiState = state

        // Store Bundle-compatible values so Android can restore the round.
        val savedPage = when (val screen = state.screen) {
            RoundScreen.Setup -> -2
            RoundScreen.Menu -> -1
            is RoundScreen.Hole -> screen.index
            RoundScreen.Summary -> HOLE_COUNT
        }

        savedStateHandle["scores"] = state.playerScores.flatten().toIntArray()
        savedStateHandle["playerCount"] = state.playerCount
        savedStateHandle["page"] = savedPage
    }

    private fun restoreState(): RoundUiState {
        val scores = restoreScores()
        val screen = restoreScreen()

        return RoundUiState(
            screen = screen,
            playerScores = scores,
            playerCount = (savedStateHandle.get<Int>("playerCount") ?: 1).coerceIn(1, 4),
        )
    }

    private fun restoreScores(): List<List<Int>> {
        val saved = savedStateHandle.get<IntArray>("scores") ?: return RoundUiState().playerScores
        if (saved.any { it < MIN_SCORE }) return RoundUiState().playerScores
        return when (saved.size) {
            HOLE_COUNT * 4 -> saved.toList().chunked(HOLE_COUNT)
            HOLE_COUNT -> listOf(saved.toList()) + RoundUiState().playerScores.drop(1)
            else -> RoundUiState().playerScores
        }
    }

    private fun restoreScreen(): RoundScreen {
        val savedPage = savedStateHandle.get<Int>("page")
        if (savedPage == null) {
            return RoundScreen.Menu
        }

        return when (savedPage) {
            in 0 until HOLE_COUNT -> RoundScreen.Hole(savedPage)
            -2 -> RoundScreen.Setup
            HOLE_COUNT -> RoundScreen.Summary
            else -> RoundScreen.Menu
        }
    }
}

package com.jesper.dgscorekeeper.round

import androidx.lifecycle.SavedStateHandle
import org.junit.Assert.assertEquals
import org.junit.Test

class RoundViewModelTest {
    @Test
    fun `new round starts on hole one with eighteen default scores`() {
        val model = RoundViewModel(SavedStateHandle())
        assertEquals(RoundScreen.Menu, model.uiState.screen)
        model.openSetup()
        model.startRound()
        assertEquals(RoundScreen.Hole(0), model.uiState.screen)
        val expectedScores = List(18) {
            3
        }
        assertEquals(expectedScores, model.uiState.scores)
        assertEquals(54, model.uiState.totalScore)
    }

    @Test
    fun `scores belong to each hole and cannot fall below one`() {
        val model = RoundViewModel(SavedStateHandle())
        model.openSetup()
        model.startRound()
        repeat(5) {
            model.decreaseScore()
        }
        model.next()
        model.increaseScore()
        model.previous()
        assertEquals(1, model.uiState.scores[0])
        assertEquals(4, model.uiState.scores[1])
        assertEquals(53, model.uiState.totalScore)
    }

    @Test
    fun `summary allows returning to hole eighteen and advancing resets the round`() {
        val model = RoundViewModel(SavedStateHandle())
        model.openSetup()
        model.startRound()
        repeat(17) {
            model.next()
        }
        assertEquals(RoundScreen.Hole(17), model.uiState.screen)
        model.increaseScore()
        model.next()
        assertEquals(RoundScreen.Summary, model.uiState.screen)
        assertEquals(55, model.uiState.totalScore)
        model.previous()
        assertEquals(RoundScreen.Hole(17), model.uiState.screen)
        assertEquals(4, model.uiState.scores[17])
        model.next()
        model.next()
        assertEquals(RoundUiState(), model.uiState)
        model.openSetup()
        model.startRound()
        assertEquals(RoundScreen.Hole(0), model.uiState.screen)
        val expectedScores = List(18) {
            3
        }
        assertEquals(expectedScores, model.uiState.scores)
    }

    @Test
    fun `previous from hole one returns to setup and preserves scores`() {
        val handle = SavedStateHandle()
        val model = RoundViewModel(handle)
        model.openSetup()
        model.startRound()
        model.increaseScore()

        model.previous()

        assertEquals(RoundScreen.Setup, model.uiState.screen)
        assertEquals(4, model.uiState.scores[0])

        val restoredModel = RoundViewModel(handle)
        assertEquals(RoundScreen.Setup, restoredModel.uiState.screen)
        restoredModel.openSetup()
        restoredModel.startRound()
        assertEquals(RoundScreen.Hole(0), restoredModel.uiState.screen)
        assertEquals(4, restoredModel.uiState.scores[0])
    }

    @Test
    fun `menu ignores score changes and swipes`() {
        val model = RoundViewModel(SavedStateHandle())
        model.increaseScore()
        model.decreaseScore()
        model.next()
        model.previous()
        assertEquals(RoundUiState(), model.uiState)
    }

    @Test
    fun `saved state restores hole scores and summary and records reset`() {
        val handle = SavedStateHandle()
        val model = RoundViewModel(handle)
        model.openSetup()
        model.startRound()
        model.increaseScore()
        model.next()
        assertEquals(model.uiState, RoundViewModel(handle).uiState)
        repeat(17) {
            model.next()
        }
        assertEquals(RoundScreen.Summary, RoundViewModel(handle).uiState.screen)
        model.next()
        assertEquals(RoundUiState(), RoundViewModel(handle).uiState)
    }
    @Test
    fun `setup bounds players and preserves independent scores`() {
        val handle = SavedStateHandle()
        val model = RoundViewModel(handle)
        model.startRound()
        assertEquals(RoundScreen.Menu, model.uiState.screen)
        model.openSetup()
        model.setPlayerCount(9)
        assertEquals(4, model.uiState.playerCount)
        model.setPlayerCount(0)
        assertEquals(1, model.uiState.playerCount)
        model.setPlayerCount(4)
        model.next()
        assertEquals(RoundScreen.Setup, model.uiState.screen)
        model.startRound()
        model.increaseScore(3)
        model.decreaseScore(1)
        model.increaseScore(4)
        assertEquals(listOf(3, 2, 3, 4), model.uiState.playerScores.map { it[0] })
        model.next()
        model.increaseScore(0)
        assertEquals(model.uiState, RoundViewModel(handle).uiState)
        assertEquals(listOf(55, 53, 54, 55), model.uiState.totals)
    }
}

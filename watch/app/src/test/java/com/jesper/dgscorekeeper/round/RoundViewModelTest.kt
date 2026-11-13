package com.jesper.dgscorekeeper.round

import androidx.lifecycle.SavedStateHandle
import org.junit.Assert.assertEquals
import org.junit.Test

class RoundViewModelTest {
    @Test
    fun `new round starts on hole one with eighteen default scores`() {
        val model = RoundViewModel(SavedStateHandle())
        assertEquals(RoundScreen.Menu, model.uiState.screen)
        model.startRound()
        assertEquals(RoundScreen.Hole(0), model.uiState.screen)
        assertEquals(List(18) { 3 }, model.uiState.scores)
        assertEquals(54, model.uiState.totalScore)
    }

    @Test
    fun `scores belong to each hole and cannot fall below one`() {
        val model = RoundViewModel(SavedStateHandle())
        model.startRound()
        model.previous()
        assertEquals(RoundScreen.Hole(0), model.uiState.screen)
        repeat(5) { model.decreaseScore() }
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
        model.startRound()
        repeat(17) { model.next() }
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
        model.startRound()
        assertEquals(RoundScreen.Hole(0), model.uiState.screen)
        assertEquals(List(18) { 3 }, model.uiState.scores)
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
        model.startRound()
        model.increaseScore()
        model.next()
        assertEquals(model.uiState, RoundViewModel(handle).uiState)
        repeat(17) { model.next() }
        assertEquals(RoundScreen.Summary, RoundViewModel(handle).uiState.screen)
        model.next()
        assertEquals(RoundUiState(), RoundViewModel(handle).uiState)
    }
}

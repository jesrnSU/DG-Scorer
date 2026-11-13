package com.jesper.dgscorekeeper.round

const val HOLE_COUNT = 18
const val DEFAULT_SCORE = 3
const val MIN_SCORE = 1

sealed interface RoundScreen {
    data object Menu : RoundScreen
    data class Hole(val index: Int) : RoundScreen
    data object Summary : RoundScreen
}

data class RoundUiState(
    val screen: RoundScreen = RoundScreen.Menu,
    val scores: List<Int> = List(HOLE_COUNT) { DEFAULT_SCORE },
) {
    val totalScore: Int get() = scores.sum()
}

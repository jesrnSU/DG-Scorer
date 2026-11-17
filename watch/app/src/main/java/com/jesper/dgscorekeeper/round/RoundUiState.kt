package com.jesper.dgscorekeeper.round

const val HOLE_COUNT = 18
const val DEFAULT_SCORE = 3
const val MIN_SCORE = 1

sealed interface RoundScreen {
    data object Setup : RoundScreen

    data object Menu : RoundScreen

    data class Hole(
        val index: Int,
    ) : RoundScreen

    data object Summary : RoundScreen
}

// Replace this fixture with the nearest course returned by a location-backed repository.
data class HoleInfo(val lengthMeters: Int, val par: Int)
data class CourseInfo(val name: String, val holes: List<HoleInfo>)
val PLACEHOLDER_COURSE = CourseInfo(
    name = "Sample Disc Golf Course",
    holes = List(HOLE_COUNT) { HoleInfo(lengthMeters = 90, par = 3) },
)

data class RoundUiState(
    val screen: RoundScreen = RoundScreen.Menu,
    val playerCount: Int = 1,
    val course: CourseInfo = PLACEHOLDER_COURSE,
    val playerScores: List<List<Int>> = List(4) { List(HOLE_COUNT) { DEFAULT_SCORE } },
) {
    val scores: List<Int> get() = playerScores.first()
    val totalScore: Int get() = scores.sum()
    val totals: List<Int> get() = playerScores.take(playerCount).map { it.sum() }
}

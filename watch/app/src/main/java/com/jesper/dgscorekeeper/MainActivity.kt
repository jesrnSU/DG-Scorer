package com.jesper.dgscorekeeper

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.viewModels
import androidx.wear.compose.material.MaterialTheme
import com.jesper.dgscorekeeper.round.RoundViewModel
import com.jesper.dgscorekeeper.ui.ScorekeeperApp

class MainActivity : ComponentActivity() {
    private val roundViewModel: RoundViewModel by viewModels()

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContent {
            MaterialTheme {
                ScorekeeperApp(
                    state = roundViewModel.uiState,
                    onStart = roundViewModel::startRound,
                    onIncrease = roundViewModel::increaseScore,
                    onDecrease = roundViewModel::decreaseScore,
                    onNext = roundViewModel::next,
                    onPrevious = roundViewModel::previous,
                )
            }
        }
    }
}

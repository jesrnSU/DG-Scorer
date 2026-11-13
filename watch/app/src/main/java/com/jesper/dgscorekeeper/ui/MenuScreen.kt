package com.jesper.dgscorekeeper.ui

import androidx.compose.runtime.Composable
import androidx.compose.ui.res.stringResource
import androidx.wear.compose.material.Button
import androidx.wear.compose.material.Text
import com.jesper.dgscorekeeper.R

@Composable
fun MenuScreen(
    onStart: () -> Unit,
) {
    ScreenLayout(
        title = stringResource(R.string.app_name),
    ) {
        Button(onClick = onStart) {
            Text(
                text = stringResource(R.string.start),
            )
        }
    }
}

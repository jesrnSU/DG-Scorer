package com.jesper.dgscorekeeper.ui

import androidx.compose.foundation.background
import androidx.compose.foundation.gestures.detectHorizontalDragGestures
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.ColumnScope
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.padding
import androidx.compose.runtime.Composable
import androidx.compose.runtime.rememberUpdatedState
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.input.pointer.pointerInput
import androidx.compose.ui.platform.LocalDensity
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import androidx.wear.compose.material.MaterialTheme
import androidx.wear.compose.material.Text

@Composable
internal fun ScreenLayout(
    title: String,
    onNext: (() -> Unit)? = null,
    onPrevious: (() -> Unit)? = null,
    content: @Composable ColumnScope.() -> Unit,
) {
    val swipeThreshold = with(LocalDensity.current) { 40.dp.toPx() }
    val currentNext = rememberUpdatedState(onNext)
    val currentPrevious = rememberUpdatedState(onPrevious)
    val gesturesEnabled = onNext != null || onPrevious != null
    Column(
        modifier = Modifier
            .fillMaxSize()
            .background(MaterialTheme.colors.background)
            .pointerInput(gesturesEnabled, swipeThreshold) {
                if (gesturesEnabled) {
                    var dragDistance = 0f
                    detectHorizontalDragGestures(
                        onDragStart = { dragDistance = 0f },
                        onDragCancel = { dragDistance = 0f },
                        onDragEnd = {
                            when {
                                dragDistance <= -swipeThreshold -> currentNext.value?.invoke()
                                dragDistance >= swipeThreshold -> currentPrevious.value?.invoke()
                            }
                            dragDistance = 0f
                        },
                        onHorizontalDrag = { change, amount ->
                            change.consume()
                            dragDistance += amount
                        },
                    )
                }
            }
            .padding(horizontal = 24.dp, vertical = 16.dp),
        verticalArrangement = Arrangement.spacedBy(8.dp, Alignment.CenterVertically),
        horizontalAlignment = Alignment.CenterHorizontally,
    ) {
        Text(
            text = title,
            style = MaterialTheme.typography.title3,
            color = MaterialTheme.colors.onBackground,
            textAlign = TextAlign.Center,
        )
        content()
    }
}

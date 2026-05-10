package com.aarcsx.krisho.core.designsystem.theme

import androidx.compose.foundation.isSystemInDarkTheme
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.lightColorScheme
import androidx.compose.material3.darkColorScheme
import androidx.compose.runtime.Composable
import androidx.compose.ui.graphics.Color

private val LightColorScheme = lightColorScheme(
    primary = ForestGreen,
    onPrimary = SurfaceWhite,
    secondary = SoftSage,
    onSecondary = Neutral10,
    background = WarmBackground,
    onBackground = Neutral10,
    surface = SurfaceWhite,
    onSurface = Neutral10,
    error = AlertRed,
    onError = SurfaceWhite
)

@Composable
fun KrishoTheme(
    darkTheme: Boolean = isSystemInDarkTheme(),
    content: @Composable () -> Unit
) {
    // Agriculture apps often stay light/warm to feel organic, 
    // but we can provide a muted dark mode if needed.
    val colorScheme = if (darkTheme) {
        // Simple dark scheme for now
        darkColorScheme(
            primary = SoftSage,
            onPrimary = Neutral10,
            secondary = ForestGreen,
            background = Color(0xFF121212),
            surface = Color(0xFF1E1E1E)
        )
    } else {
        LightColorScheme
    }

    MaterialTheme(
        colorScheme = colorScheme,
        typography = Typography,
        content = content
    )
}

package com.aarcsx.krisho.features.home

import androidx.compose.runtime.Immutable

@Immutable
data class WeatherInfo(
    val temperature: String,
    val condition: String,
    val humidity: String,
    val wind: String
)

@Immutable
data class DiseaseAlert(
    val title: String,
    val description: String
)

@Immutable
data class RecentScan(
    val id: String,
    val cropName: String,
    val status: String,
    val date: String,
    val imageUrl: String
)

@Immutable
data class Recommendation(
    val id: String,
    val title: String,
    val description: String,
    val imageUrl: String
)

@Immutable
data class HomeUiState(
    val userName: String = "Farmer",
    val location: String = "Detecting...",
    val weather: WeatherInfo? = null,
    val alerts: List<DiseaseAlert> = emptyList(),
    val recentScans: List<RecentScan> = emptyList(),
    val recommendations: List<Recommendation> = emptyList(),
    val isLoading: Boolean = false,
    val error: String? = null
)

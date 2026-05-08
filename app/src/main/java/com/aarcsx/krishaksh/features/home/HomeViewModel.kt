package com.aarcsx.krishaksh.features.home

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.delay
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import javax.inject.Inject

@HiltViewModel
class HomeViewModel @Inject constructor() : ViewModel() {

    private val _uiState = MutableStateFlow(HomeUiState())
    val uiState: StateFlow<HomeUiState> = _uiState.asStateFlow()

    init {
        loadHomeData()
    }

    fun loadHomeData() {
        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isLoading = true)
            // Mock delay
            delay(500)
            _uiState.value = HomeUiState(
                userName = "Ram Singh",
                location = "Bhiwani, Haryana",
                weather = WeatherInfo(
                    temperature = "28°C",
                    condition = "Sunny Day",
                    humidity = "60%",
                    wind = "5 km/h"
                ),
                alerts = listOf(
                    DiseaseAlert(
                        title = "Nearby Alerts",
                        description = "Yellow Rust alert in neighboring farms. Inspect your wheat crop today."
                    )
                ),
                recentScans = listOf(
                    RecentScan("1", "Wheat Scan", "Healthy", "2 hours ago", "https://images.unsplash.com/photo-1574323347407-f5e1ad6d020b?q=80&w=2187&auto=format&fit=crop"),
                    RecentScan("2", "Mustard Scan", "Aphid Warning", "Yesterday", "https://images.unsplash.com/photo-1599424423956-6f81014ccbe0?q=80&w=2169&auto=format&fit=crop")
                ),
                recommendations = listOf(
                    Recommendation("1", "Eco-Green Fertilizer", "Boost your winter yield by up to 20% with organic nutrients.", "https://images.unsplash.com/photo-1628352081506-83c43123ed6d?q=80&w=2196&auto=format&fit=crop")
                ),
                isLoading = false
            )
        }
    }
}

package com.aarcsx.krishaksh.features.scan

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
class ScanViewModel @Inject constructor() : ViewModel() {

    private val _uiState = MutableStateFlow(ScanUiState())
    val uiState: StateFlow<ScanUiState> = _uiState.asStateFlow()

    fun analyzeImage() {
        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isAnalyzing = true)
            // Mock AI Analysis delay
            delay(3000)
            _uiState.value = _uiState.value.copy(
                isAnalyzing = false,
                result = ScanResult(
                    diseaseName = "Late Blight",
                    confidence = "94%",
                    symptoms = "Small, dark, water-soaked spots on lower leaves.",
                    prevention = "Use certified disease-free seeds and rotate crops.",
                    treatment = "Apply recommended fungicides early in the season.",
                    recommendationTitle = "Copper Fungicide",
                    recommendationDesc = "Effective against Late Blight in Potatoes."
                )
            )
        }
    }

    fun reset() {
        _uiState.value = ScanUiState()
    }
}

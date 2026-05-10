package com.aarcsx.krishaksh.features.scan

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.aarcsx.krishaksh.core.common.ApiResult
import com.aarcsx.krishaksh.core.repository.ScanRepository
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import javax.inject.Inject

@HiltViewModel
class ScanViewModel @Inject constructor(
    private val scanRepository: ScanRepository,
    private val userRepository: com.aarcsx.krishaksh.core.repository.UserRepository
) : ViewModel() {

    private val _uiState = MutableStateFlow(ScanUiState())
    val uiState: StateFlow<ScanUiState> = _uiState.asStateFlow()

    private val _userCrops = MutableStateFlow<List<String>>(emptyList())
    val userCrops: StateFlow<List<String>> = _userCrops.asStateFlow()

    init {
        loadUserCrops()
    }

    private fun loadUserCrops() {
        viewModelScope.launch {
            userRepository.getProfile().collect { res ->
                if (res is ApiResult.Success) {
                    _userCrops.value = res.data.crops
                }
            }
        }
    }

    fun analyzeImage(imageBytes: ByteArray = ByteArray(0), cropName: String = "Wheat") {
        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isAnalyzing = true, error = null)
            when (val res = scanRepository.runScanLifecycle(cropName = cropName, imageBytes = imageBytes)) {
                is ApiResult.Success -> {
                    val s = res.data
                    _uiState.value = _uiState.value.copy(
                        isAnalyzing = false,
                        result = ScanResult(
                            diseaseName = s.diseaseName,
                            confidence = s.confidence,
                            symptoms = s.symptoms,
                            prevention = s.prevention,
                            treatment = s.treatment,
                            recommendationTitle = s.recommendationTitle,
                            recommendationDesc = s.recommendationDesc
                        )
                    )
                }
                is ApiResult.Error -> {
                    _uiState.value = _uiState.value.copy(isAnalyzing = false, error = res.message)
                }
                ApiResult.Loading -> Unit
            }
        }
    }

    fun reset() {
        _uiState.value = ScanUiState()
    }
}

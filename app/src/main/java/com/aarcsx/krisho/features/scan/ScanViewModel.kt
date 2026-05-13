package com.aarcsx.krisho.features.scan

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.aarcsx.krisho.core.common.ApiResult
import com.aarcsx.krisho.core.repository.AppRepository
import com.aarcsx.krisho.core.repository.ScanRepository
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import javax.inject.Inject

data class CropItem(
    val name: String,
    val emoji: String
)

@HiltViewModel
class ScanViewModel @Inject constructor(
    private val scanRepository: ScanRepository,
    private val appRepository: AppRepository
) : ViewModel() {

    private val _uiState = MutableStateFlow(ScanUiState())
    val uiState: StateFlow<ScanUiState> = _uiState.asStateFlow()

    private val _userCrops = MutableStateFlow<List<CropItem>>(emptyList())
    val userCrops: StateFlow<List<CropItem>> = _userCrops.asStateFlow()

    init {
        loadAvailableCrops()
    }

    private fun loadAvailableCrops() {
        viewModelScope.launch {
            when (val res = appRepository.getCrops()) {
                is ApiResult.Success -> {
                    _userCrops.value = res.data.map { CropItem(it.name, it.emoji) }
                }
                is ApiResult.Error -> {
                    // keep default crop list in UI fallback
                }
                ApiResult.Loading -> Unit
            }
        }
    }

    fun analyzeImage(imageBytes: ByteArray = ByteArray(0), cropName: String = "Wheat") {
        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isAnalyzing = true, error = null, infoMessage = null)
            when (val res = scanRepository.runScanLifecycle(cropName = cropName, imageBytes = imageBytes)) {
                is ApiResult.Success -> {
                    val s = res.data
                    if (s.remoteScanId.isNullOrBlank()) {
                        _uiState.value = _uiState.value.copy(
                            isAnalyzing = false,
                            infoMessage = "Image captured and queued offline. It will sync automatically when the network is available.",
                            error = null
                        )
                    } else {
                        _uiState.value = _uiState.value.copy(
                            isAnalyzing = false,
                            result = ScanResult(
                                imageUrl = s.imageUrl,
                                diseaseName = s.diseaseName,
                                confidence = s.confidence,
                                symptoms = s.symptoms,
                                prevention = s.prevention,
                                treatment = s.treatment,
                                recommendationTitle = s.recommendationTitle,
                                recommendationDesc = s.recommendationDesc
                            ),
                            infoMessage = null
                        )
                    }
                }
                is ApiResult.Error -> {
                    _uiState.value = _uiState.value.copy(isAnalyzing = false, error = res.message, infoMessage = null)
                }
                ApiResult.Loading -> Unit
            }
        }
    }

    fun reset() {
        _uiState.value = ScanUiState()
    }
}

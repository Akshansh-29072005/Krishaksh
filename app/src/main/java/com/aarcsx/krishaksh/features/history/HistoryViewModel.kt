package com.aarcsx.krishaksh.features.history

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.aarcsx.krishaksh.core.local.room.entity.ScanEntity
import com.aarcsx.krishaksh.core.repository.ScanRepository
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.*
import javax.inject.Inject

data class HistoryUiState(
    val scans: List<ScanEntity> = emptyList(),
    val isLoading: Boolean = false
)

@HiltViewModel
class HistoryViewModel @Inject constructor(
    private val scanRepository: ScanRepository
) : ViewModel() {

    private val _uiState = MutableStateFlow(HistoryUiState())
    val uiState: StateFlow<HistoryUiState> = _uiState.asStateFlow()

    init {
        loadScans()
    }

    private fun loadScans() {
        scanRepository.getRecentScans()
            .onStart { _uiState.update { it.copy(isLoading = true) } }
            .onEach { list ->
                _uiState.update { it.copy(scans = list, isLoading = false) }
            }
            .launchIn(viewModelScope)
    }
}
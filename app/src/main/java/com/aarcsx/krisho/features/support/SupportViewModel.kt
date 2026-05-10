package com.aarcsx.krisho.features.support

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.aarcsx.krisho.core.common.voice.VoiceRecorderManager
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.Job
import kotlinx.coroutines.delay
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.update
import kotlinx.coroutines.launch
import java.io.File
import javax.inject.Inject

data class SupportUiState(
    val message: String = "",
    val isRecording: Boolean = false,
    val recordingDuration: Long = 0L,
    val audioFile: File? = null,
    val isSending: Boolean = false,
    val error: String? = null,
    val success: Boolean = false,
    val isCallbackRequested: Boolean = false
)

@HiltViewModel
class SupportViewModel @Inject constructor(
    private val voiceRecorderManager: VoiceRecorderManager
) : ViewModel() {

    private val _uiState = MutableStateFlow(SupportUiState())
    val uiState: StateFlow<SupportUiState> = _uiState.asStateFlow()

    private var recordingJob: Job? = null

    fun onMessageChange(message: String) {
        _uiState.update { it.copy(message = message) }
    }

    fun startRecording() {
        try {
            val file = voiceRecorderManager.startRecording()
            _uiState.update { it.copy(isRecording = true, audioFile = file, recordingDuration = 0L) }
            
            recordingJob?.cancel()
            recordingJob = viewModelScope.launch {
                while (true) {
                    delay(1000)
                    _uiState.update { it.copy(recordingDuration = it.recordingDuration + 1) }
                }
            }
        } catch (e: Exception) {
            _uiState.update { it.copy(error = "Failed to start recording") }
        }
    }

    fun stopRecording() {
        voiceRecorderManager.stopRecording()
        recordingJob?.cancel()
        _uiState.update { it.copy(isRecording = false) }
    }

    fun sendSupportRequest() {
        viewModelScope.launch {
            _uiState.update { it.copy(isSending = true) }
            // TODO: Integrate with real Support API
            delay(2000)
            _uiState.update { it.copy(isSending = false, success = true, message = "", audioFile = null) }
        }
    }

    fun requestCallback(phoneNumber: String) {
        viewModelScope.launch {
            _uiState.update { it.copy(isSending = true) }
            // TODO: API call to request callback
            delay(1500)
            _uiState.update { it.copy(isSending = false, isCallbackRequested = true) }
        }
    }

    fun clearError() {
        _uiState.update { it.copy(error = null) }
    }

    fun clearSuccess() {
        _uiState.update { it.copy(success = false) }
    }

    override fun onCleared() {
        super.onCleared()
        if (_uiState.value.isRecording) {
            voiceRecorderManager.stopRecording()
        }
    }
}
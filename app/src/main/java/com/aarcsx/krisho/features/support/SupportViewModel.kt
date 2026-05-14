package com.aarcsx.krisho.features.support

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.aarcsx.krisho.core.common.voice.VoiceRecorderManager
import com.aarcsx.krisho.core.network.api.UserApiService
import com.aarcsx.krisho.core.network.dto.CreateTicketDto
import com.aarcsx.krisho.core.network.dto.SupportTicketDto
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.Job
import kotlinx.coroutines.delay
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.update
import kotlinx.coroutines.launch
import okhttp3.MediaType.Companion.toMediaType
import okhttp3.MultipartBody
import okhttp3.RequestBody.Companion.asRequestBody
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
    val isCallbackRequested: Boolean = false,
    val tickets: List<SupportTicketDto> = emptyList(),
    val selectedTicketId: String? = null
)

@HiltViewModel
class SupportViewModel @Inject constructor(
    private val voiceRecorderManager: VoiceRecorderManager,
    private val apiService: UserApiService
) : ViewModel() {

    private val _uiState = MutableStateFlow(SupportUiState())
    val uiState: StateFlow<SupportUiState> = _uiState.asStateFlow()

    private var recordingJob: Job? = null

    init {
        loadTickets()
    }

    private fun loadTickets() {
        viewModelScope.launch {
            try {
                val response = apiService.getTickets()
                if (response.isSuccessful) {
                    val tickets = response.body()?.data ?: emptyList()
                    _uiState.update { it.copy(tickets = tickets) }
                }
            } catch (e: Exception) {
                _uiState.update { it.copy(error = "Failed to load tickets: ${e.message}") }
            }
        }
    }

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
            try {
                _uiState.update { it.copy(isSending = true, error = null) }
                
                val currentState = _uiState.value
                val request = CreateTicketDto(
                    title = "Support Request",
                    description = currentState.message,
                    priority = "high",
                    request_callback = false
                )
                
                val response = apiService.createTicket(request)
                
                if (response.isSuccessful) {
                    // Upload voice if available
                    response.body()?.data?.let { ticket ->
                        if (currentState.audioFile != null) {
                            uploadVoiceAttachment(ticket.id, currentState.audioFile)
                        }
                    }
                    
                    _uiState.update { 
                        it.copy(
                            isSending = false, 
                            success = true, 
                            message = "", 
                            audioFile = null,
                            recordingDuration = 0L
                        ) 
                    }
                    
                    // Reload tickets
                    delay(1000)
                    loadTickets()
                } else {
                    _uiState.update { 
                        it.copy(
                            isSending = false, 
                            error = "Failed to create support request"
                        ) 
                    }
                }
            } catch (e: Exception) {
                _uiState.update { 
                    it.copy(
                        isSending = false, 
                        error = "Error: ${e.message}"
                    ) 
                }
            }
        }
    }

    private suspend fun uploadVoiceAttachment(ticketId: String, audioFile: File) {
        try {
            val mediaType = "audio/mpeg".toMediaType()
            val requestBody = audioFile.asRequestBody(mediaType)
            val filePart = MultipartBody.Part.createFormData("voice", audioFile.name, requestBody)
            apiService.uploadVoice(ticketId, filePart)
        } catch (e: Exception) {
            _uiState.update { it.copy(error = "Failed to upload voice: ${e.message}") }
        }
    }

    fun requestCallback(ticketId: String) {
        viewModelScope.launch {
            try {
                _uiState.update { it.copy(isSending = true, error = null) }
                
                val response = apiService.requestCallback(ticketId)
                
                if (response.isSuccessful) {
                    _uiState.update { 
                        it.copy(
                            isSending = false, 
                            isCallbackRequested = true
                        ) 
                    }
                    loadTickets()
                } else {
                    _uiState.update { 
                        it.copy(
                            isSending = false, 
                            error = "Failed to request callback"
                        ) 
                    }
                }
            } catch (e: Exception) {
                _uiState.update { 
                    it.copy(
                        isSending = false, 
                        error = "Error: ${e.message}"
                    ) 
                }
            }
        }
    }

    fun clearError() {
        _uiState.update { it.copy(error = null) }
    }

    fun clearSuccess() {
        _uiState.update { it.copy(success = false) }
    }

    fun selectTicket(ticketId: String) {
        _uiState.update { it.copy(selectedTicketId = ticketId) }
    }

    override fun onCleared() {
        super.onCleared()
        if (_uiState.value.isRecording) {
            voiceRecorderManager.stopRecording()
        }
    }
}
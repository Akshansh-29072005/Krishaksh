package com.aarcsx.krisho.features.profile

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.aarcsx.krisho.core.common.ApiResult
import com.aarcsx.krisho.core.repository.UserRepository
import com.aarcsx.krisho.core.network.dto.UserMeDto
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.*
import kotlinx.coroutines.launch
import javax.inject.Inject

data class ProfileUiState(
    val profile: UserMeDto? = null,
    val isLoading: Boolean = false,
    val error: String? = null,
    val showLanguageDialog: Boolean = false,
    val showCropsDialog: Boolean = false,
    val selectedLanguage: String = "English",
    val isUpdating: Boolean = false
)

@HiltViewModel
class ProfileViewModel @Inject constructor(
    private val userRepository: UserRepository
) : ViewModel() {

    private val _uiState = MutableStateFlow(ProfileUiState())
    val uiState: StateFlow<ProfileUiState> = _uiState.asStateFlow()

    init {
        loadProfile()
        viewModelScope.launch {
            userRepository.languageSetting.collect { lang ->
                _uiState.update { it.copy(selectedLanguage = if (lang == "hi") "Hindi" else "English") }
            }
        }
    }

    fun loadProfile() {
        viewModelScope.launch {
            userRepository.getProfile().collect { result ->
                when (result) {
                    is ApiResult.Loading -> _uiState.update { it.copy(isLoading = true) }
                    is ApiResult.Success -> _uiState.update { it.copy(profile = result.data, isLoading = false, error = null) }
                    is ApiResult.Error -> _uiState.update { it.copy(isLoading = false, error = result.message) }
                }
            }
        }
    }

    fun toggleLanguageDialog(show: Boolean) {
        _uiState.update { it.copy(showLanguageDialog = show) }
    }

    fun toggleCropsDialog(show: Boolean) {
        _uiState.update { it.copy(showCropsDialog = show) }
    }

    fun updateLanguage(langCode: String) {
        viewModelScope.launch {
            userRepository.saveLanguage(langCode)
            toggleLanguageDialog(false)
        }
    }

    fun updateCrops(crops: List<String>) {
        viewModelScope.launch {
            _uiState.update { it.copy(isUpdating = true) }
            val result = userRepository.updateCrops(crops)
            if (result is ApiResult.Success) {
                _uiState.update { it.copy(profile = result.data, isUpdating = false) }
            } else if (result is ApiResult.Error) {
                _uiState.update { it.copy(isUpdating = false, error = result.message) }
            }
            toggleCropsDialog(false)
        }
    }

    fun logout() {
        viewModelScope.launch {
            userRepository.logout()
        }
    }
}
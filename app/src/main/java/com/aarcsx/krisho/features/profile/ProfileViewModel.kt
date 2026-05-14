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

import com.aarcsx.krisho.core.network.api.WeatherApiService
import com.aarcsx.krisho.core.util.LocationProvider

data class ProfileUiState(
    val profile: UserMeDto? = null,
    val location: String = "Detecting...",
    val isLoading: Boolean = false,
    val error: String? = null,
    val showLanguageDialog: Boolean = false,
    val showPhoneDialog: Boolean = false,
    val phoneInput: String = "",
    val selectedLanguage: String = "English",
    val isUpdating: Boolean = false
)

@HiltViewModel
class ProfileViewModel @Inject constructor(
    private val userRepository: UserRepository,
    private val locationProvider: LocationProvider,
    private val weatherApiService: WeatherApiService
) : ViewModel() {

    private val _uiState = MutableStateFlow(ProfileUiState())
    val uiState: StateFlow<ProfileUiState> = _uiState.asStateFlow()

    init {
        loadProfile()
        loadLocation()
        viewModelScope.launch {
            userRepository.languageSetting.collect { lang ->
                _uiState.update { it.copy(selectedLanguage = if (lang == "hi") "Hindi" else "English") }
            }
        }
    }

    private fun loadLocation() {
        viewModelScope.launch {
            try {
                val location = locationProvider.getCurrentLocation()
                if (location != null) {
                    val weatherRes = weatherApiService.getCurrentWeather(location.latitude, location.longitude)
                    if (weatherRes.isSuccessful) {
                        val weatherData = weatherRes.body()?.data
                        _uiState.update { it.copy(location = weatherData?.location_name ?: "Unknown Location") }
                    }
                } else {
                    _uiState.update { it.copy(location = "GPS Disabled") }
                }
            } catch (e: Exception) {
                _uiState.update { it.copy(location = "Location Offline") }
            }
        }
    }

    fun loadProfile() {
        viewModelScope.launch {
            userRepository.getProfile().collect { result ->
                when (result) {
                    is ApiResult.Loading -> _uiState.update { it.copy(isLoading = true) }
                    is ApiResult.Success -> _uiState.update { it.copy(profile = result.data, phoneInput = result.data.phone ?: "", isLoading = false, error = null) }
                    is ApiResult.Error -> _uiState.update { it.copy(isLoading = false, error = result.message) }
                }
            }
        }
    }

    fun toggleLanguageDialog(show: Boolean) {
        _uiState.update { it.copy(showLanguageDialog = show) }
    }

    fun togglePhoneDialog(show: Boolean) {
        _uiState.update { it.copy(showPhoneDialog = show) }
    }

    fun updatePhoneInput(phone: String) {
        _uiState.update { it.copy(phoneInput = phone) }
    }

    fun updatePhone(phone: String) {
        viewModelScope.launch {
            _uiState.update { it.copy(isUpdating = true, error = null) }
            when (val result = userRepository.updatePhone(phone)) {
                is ApiResult.Success -> _uiState.update { it.copy(profile = result.data, phoneInput = result.data.phone ?: "", isUpdating = false, error = null) }
                is ApiResult.Error -> _uiState.update { it.copy(isUpdating = false, error = result.message) }
                else -> _uiState.update { it.copy(isUpdating = false, error = "Unable to update phone number") }
            }
            togglePhoneDialog(false)
        }
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

    fun logout(onSuccess: () -> Unit = {}) {
        viewModelScope.launch {
            userRepository.logout()
            onSuccess()
        }
    }
}
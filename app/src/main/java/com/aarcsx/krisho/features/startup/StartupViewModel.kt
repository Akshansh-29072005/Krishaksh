package com.aarcsx.krisho.features.startup

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.aarcsx.krisho.BuildConfig
import com.aarcsx.krisho.core.common.ApiResult
import com.aarcsx.krisho.core.network.dto.AppConfigDto
import com.aarcsx.krisho.core.repository.AppRepository
import com.aarcsx.krisho.core.repository.AuthRepository
import com.aarcsx.krisho.navigation.Screen
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.launch
import javax.inject.Inject

data class StartupUiState(
    val isLoading: Boolean = false,
    val appConfig: AppConfigDto? = null,
    val authToken: String? = null,
    val errorMessage: String? = null
) {
    val startDestination: String = if (authToken.isNullOrBlank()) Screen.Auth.route else Screen.Home.route
    val isForceUpdateRequired: Boolean = appConfig != null && BuildConfig.VERSION_CODE < appConfig.minimum_version_code
}

@HiltViewModel
class AppStartupViewModel @Inject constructor(
    private val authRepository: AuthRepository,
    private val appRepository: AppRepository
) : ViewModel() {

    private val _uiState = MutableStateFlow(StartupUiState(isLoading = true))
    val uiState: StateFlow<StartupUiState> = _uiState.asStateFlow()

    init {
        refreshStartupState()
    }

    fun refreshStartupState() {
        viewModelScope.launch {
            _uiState.value = StartupUiState(isLoading = true)
            val authToken = authRepository.jwtToken.first()
            when (val result = appRepository.getAppConfig()) {
                is ApiResult.Success -> _uiState.value = StartupUiState(
                    isLoading = false,
                    appConfig = result.data,
                    authToken = authToken
                )
                is ApiResult.Error -> _uiState.value = StartupUiState(
                    isLoading = false,
                    authToken = authToken,
                    errorMessage = result.message
                )
                ApiResult.Loading -> _uiState.value = StartupUiState(isLoading = true)
            }
        }
    }
}

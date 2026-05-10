package com.aarcsx.krishaksh.features.auth

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.aarcsx.krishaksh.core.common.ApiResult
import com.aarcsx.krishaksh.core.repository.AuthRepository
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.update
import kotlinx.coroutines.launch
import javax.inject.Inject

data class AuthUiState(
    val isLoading: Boolean = false,
    val error: String? = null,
    val isSuccess: Boolean = false
)

@HiltViewModel
class AuthViewModel @Inject constructor(
    private val authRepository: AuthRepository
) : ViewModel() {

    private val _uiState = MutableStateFlow(AuthUiState())
    val uiState: StateFlow<AuthUiState> = _uiState.asStateFlow()

    fun onGoogleSignInResult(idToken: String?) {
        if (idToken == null) {
            _uiState.update { it.copy(error = "Google Sign In failed") }
            return
        }

        viewModelScope.launch {
            _uiState.update { it.copy(isLoading = true, error = null) }
            val result = authRepository.loginWithGoogle(idToken)
            when (result) {
                is ApiResult.Success -> _uiState.update { it.copy(isLoading = false, isSuccess = true) }
                is ApiResult.Error -> _uiState.update { it.copy(isLoading = false, error = result.message) }
                ApiResult.Loading -> _uiState.update { it.copy(isLoading = true) }
            }
        }
    }

    fun clearError() {
        _uiState.update { it.copy(error = null) }
    }
}
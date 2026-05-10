package com.aarcsx.krishaksh.features.home

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.aarcsx.krishaksh.core.common.ApiResult
import com.aarcsx.krishaksh.core.repository.ScanRepository
import com.aarcsx.krishaksh.core.repository.UserRepository
import com.aarcsx.krishaksh.core.network.api.UserApiService
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.collectLatest
import kotlinx.coroutines.flow.combine
import kotlinx.coroutines.flow.catch
import kotlinx.coroutines.flow.onStart
import kotlinx.coroutines.flow.update
import kotlinx.coroutines.launch
import kotlinx.coroutines.flow.flow
import javax.inject.Inject

@HiltViewModel
class HomeViewModel @Inject constructor(
    private val userRepository: UserRepository,
    private val scanRepository: ScanRepository,
    private val userApiService: UserApiService
) : ViewModel() {

    private val _uiState = MutableStateFlow(HomeUiState())
    val uiState: StateFlow<HomeUiState> = _uiState.asStateFlow()

    init {
        loadHomeData()
    }

    fun loadHomeData() {
        viewModelScope.launch {
            combine(
                userRepository.getProfile().onStart { emit(ApiResult.Loading) },
                scanRepository.getRecentScans(),
                flow {
                    try {
                        emit(userApiService.featuredAds())
                    } catch (e: Exception) {
                        emit(null)
                    }
                }
            ) { profileRes, scans, adsRes ->
                val profile = (profileRes as? ApiResult.Success)?.data
                val ads = adsRes?.body()?.data ?: emptyList()

                HomeUiState(
                    userName = profile?.name ?: "Farmer",
                    location = profile?.email ?: "India",
                    weather = WeatherInfo(
                        temperature = "30°C",
                        condition = "Clear Sky",
                        humidity = "45%",
                        wind = "10 km/h"
                    ),
                    alerts = emptyList(),
                    recentScans = scans.take(5).map { entity ->
                        RecentScan(
                            id = entity.id.toString(),
                            cropName = entity.cropName,
                            status = entity.diseaseName,
                            date = "Recently",
                            imageUrl = entity.imageUrl
                        )
                    },
                    recommendations = ads.map { ad ->
                        Recommendation(
                            id = ad.id,
                            title = ad.title,
                            description = "Ad",
                            imageUrl = ad.image_url
                        )
                    },
                    isLoading = (profileRes is ApiResult.Loading)
                )
            }.catch { e ->
                e.printStackTrace()
                _uiState.update { it.copy(isLoading = false, error = "Connection Error") }
            }.collectLatest { state ->
                _uiState.value = state
            }
        }
    }

    override fun onCleared() {
        super.onCleared()
        // Cancel any pending work if needed
    }
}

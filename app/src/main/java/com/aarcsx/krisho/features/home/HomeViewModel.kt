package com.aarcsx.krisho.features.home

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.aarcsx.krisho.core.common.ApiResult
import com.aarcsx.krisho.core.repository.ScanRepository
import com.aarcsx.krisho.core.repository.UserRepository
import com.aarcsx.krisho.core.network.api.UserApiService
import com.aarcsx.krisho.core.network.api.WeatherApiService
import com.aarcsx.krisho.core.util.LocationProvider
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
    private val userApiService: UserApiService,
    private val weatherApiService: WeatherApiService,
    private val locationProvider: LocationProvider
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
                },
                flow {
                    try {
                        val location = locationProvider.getCurrentLocation()
                        val lat = location?.latitude ?: 28.6139
                        val lon = location?.longitude ?: 77.2090
                        emit(weatherApiService.getCurrentWeather(lat, lon))
                    } catch (e: Exception) {
                        emit(null)
                    }
                }
            ) { profileRes, scans, adsRes, weatherRes ->
                val profile = (profileRes as? ApiResult.Success)?.data
                val ads = adsRes?.body()?.data ?: emptyList()
                val weatherData = weatherRes?.body()?.data

                HomeUiState(
                    userName = profile?.name ?: "Farmer",
                    location = weatherData?.location_name ?: "India",
                    weather = weatherData?.let {
                        WeatherInfo(
                            temperature = it.temperature,
                            condition = it.condition,
                            humidity = it.humidity,
                            wind = it.wind_speed
                        )
                    },
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

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
            _uiState.update { it.copy(isLoading = true, error = null) }
            
            // 1. Refresh scans (background)
            scanRepository.refreshScans()

            // 2. Fetch Profile & Scans
            launch {
                combine(
                    userRepository.getProfile(),
                    scanRepository.getRecentScans()
                ) { profileRes, scans ->
                    val profile = (profileRes as? ApiResult.Success)?.data
                    _uiState.update { currentState ->
                        currentState.copy(
                            userName = profile?.name ?: currentState.userName,
                            recentScans = scans.take(5).map { entity ->
                                RecentScan(
                                    id = entity.id.toString(),
                                    cropName = entity.cropName,
                                    status = entity.diseaseName,
                                    date = "Recently",
                                    imageUrl = entity.imageUrl
                                )
                            },
                            isLoading = profileRes is ApiResult.Loading
                        )
                    }
                }.catch { it.printStackTrace() }.collectLatest {}
            }

            // 3. Fetch Ads
            launch {
                try {
                    val adsRes = userApiService.featuredAds()
                    val ads = adsRes.body()?.data ?: emptyList()
                    _uiState.update { currentState ->
                        currentState.copy(
                            recommendations = ads.map { ad ->
                                Recommendation(
                                    id = ad.id,
                                    title = ad.title,
                                    description = "Ad",
                                    imageUrl = ad.image_url
                                )
                            }
                        )
                    }
                } catch (e: Exception) {
                    e.printStackTrace()
                }
            }

            // 4. Fetch Location & Weather
            launch {
                try {
                    val location = locationProvider.getCurrentLocation()
                    if (location != null) {
                        val weatherRes = weatherApiService.getCurrentWeather(location.latitude, location.longitude)
                        if (weatherRes.isSuccessful) {
                            val weatherData = weatherRes.body()?.data
                            _uiState.update { currentState ->
                                currentState.copy(
                                    location = weatherData?.location_name ?: "Unknown Location",
                                    weather = weatherData?.let {
                                        WeatherInfo(
                                            temperature = it.temperature,
                                            condition = it.condition,
                                            humidity = it.humidity,
                                            wind = it.wind_speed
                                        )
                                    }
                                )
                            }
                        }
                    } else {
                        _uiState.update { it.copy(location = "GPS Disabled") }
                    }
                } catch (e: Exception) {
                    e.printStackTrace()
                    _uiState.update { it.copy(location = "Location Error") }
                }
            }
        }
    }

    fun reloadLocationWeather() {
        viewModelScope.launch {
            _uiState.update { it.copy(location = "Detecting...") }
            try {
                val location = locationProvider.getCurrentLocation()
                if (location != null) {
                    val weatherRes = weatherApiService.getCurrentWeather(location.latitude, location.longitude)
                    if (weatherRes.isSuccessful) {
                        val weatherData = weatherRes.body()?.data
                        _uiState.update { currentState ->
                            currentState.copy(
                                location = weatherData?.location_name ?: "Unknown Location",
                                weather = weatherData?.let {
                                    WeatherInfo(
                                        temperature = it.temperature,
                                        condition = it.condition,
                                        humidity = it.humidity,
                                        wind = it.wind_speed
                                    )
                                }
                            )
                        }
                    } else {
                        _uiState.update { it.copy(location = "Weather Unavailable") }
                    }
                } else {
                    _uiState.update { it.copy(location = "Location Unavailable") }
                }
            } catch (e: Exception) {
                e.printStackTrace()
                _uiState.update { it.copy(location = "Location Error") }
            }
        }
    }

    override fun onCleared() {
        super.onCleared()
    }
}

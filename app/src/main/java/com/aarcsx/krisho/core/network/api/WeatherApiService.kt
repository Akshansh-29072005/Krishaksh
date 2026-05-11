package com.aarcsx.krisho.core.network.api

import com.aarcsx.krisho.core.network.dto.ApiEnvelope
import retrofit2.Response
import retrofit2.http.GET
import retrofit2.http.Query

interface WeatherApiService {
    @GET("weather")
    suspend fun getCurrentWeather(
        @Query("lat") lat: Double,
        @Query("lon") lon: Double
    ): Response<ApiEnvelope<WeatherResponseDto>>
}

data class WeatherResponseDto(
    val temperature: String,
    val condition: String,
    val humidity: String,
    val wind_speed: String,
    val location_name: String
)

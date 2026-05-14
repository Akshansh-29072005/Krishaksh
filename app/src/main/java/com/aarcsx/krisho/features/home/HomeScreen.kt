package com.aarcsx.krisho.features.home

import android.Manifest
import android.os.Build
import androidx.activity.compose.rememberLauncherForActivityResult
import androidx.activity.result.contract.ActivityResultContracts
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.verticalScroll
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.aarcsx.krisho.R
import androidx.hilt.navigation.compose.hiltViewModel
import com.aarcsx.krisho.core.designsystem.components.*

@Composable
fun HomeScreen(
    viewModel: HomeViewModel = hiltViewModel(),
    onScanClick: () -> Unit = {}
) {
    val uiState by viewModel.uiState.collectAsState()
    
    // Permission Launcher for Notifications
    val notificationPermissionLauncher = rememberLauncherForActivityResult(
        ActivityResultContracts.RequestPermission()
    ) { _ -> }

    // Launcher for Location and Camera
    val generalPermissionLauncher = rememberLauncherForActivityResult(
        ActivityResultContracts.RequestMultiplePermissions()
    ) { permissions ->
        // When user responds to the permission dialog, reload location & weather
        val locationGranted = permissions[Manifest.permission.ACCESS_FINE_LOCATION] == true ||
                permissions[Manifest.permission.ACCESS_COARSE_LOCATION] == true
        if (locationGranted) {
            viewModel.reloadLocationWeather()
        }
    }

    // Request vital permissions on launch
    LaunchedEffect(Unit) {
        generalPermissionLauncher.launch(
            arrayOf(
                Manifest.permission.CAMERA,
                Manifest.permission.ACCESS_FINE_LOCATION,
                Manifest.permission.ACCESS_COARSE_LOCATION
            )
        )
    }

    Scaffold(
        containerColor = Color(0xFFF5F7F2)
    ) { padding ->
        Box(modifier = Modifier.fillMaxSize().padding(padding)) {
            when {
                uiState.isLoading -> {
                    CircularProgressIndicator(
                        modifier = Modifier.align(Alignment.Center),
                        color = ForestGreen
                    )
                }
                uiState.error != null -> {
                    Column(
                        modifier = Modifier.fillMaxSize(),
                        verticalArrangement = Arrangement.Center,
                        horizontalAlignment = Alignment.CenterHorizontally
                    ) {
                        Text(text = uiState.error!!)
                        Button(onClick = { viewModel.loadHomeData() }) {
                            Text("Retry")
                        }
                    }
                }
                else -> {
                    HomeContent(
                        uiState = uiState,
                        onScanClick = onScanClick,
                        onNotificationClick = {
                            if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.TIRAMISU) {
                                notificationPermissionLauncher.launch(Manifest.permission.POST_NOTIFICATIONS)
                            }
                        }
                    )
                }
            }
        }
    }
}

@Composable
private fun HomeContent(
    uiState: HomeUiState,
    onScanClick: () -> Unit,
    onNotificationClick: () -> Unit
) {
    Column(
        modifier = Modifier
            .fillMaxSize()
            .verticalScroll(rememberScrollState())
            .padding(bottom = 32.dp),
        verticalArrangement = Arrangement.spacedBy(8.dp)
    ) {
        // App Title
        Text(
            text = "Krisho",
            modifier = Modifier
                .fillMaxWidth()
                .padding(top = 8.dp),
            textAlign = androidx.compose.ui.text.style.TextAlign.Center,
            style = MaterialTheme.typography.displaySmall,
            fontWeight = FontWeight.Black,
            color = ForestGreen,
            letterSpacing = (-1).sp
        )

        // 1. Header Section
        AgriHeader(
            userName = uiState.userName,
            location = uiState.location,
            onNotificationClick = onNotificationClick
        )

        // 2. Hero Scan Section
        AgriHeroSection(onScanClick = onScanClick)

        // 3. Weather + Alert Section
        AgriWeatherAlertRow(
            temp = uiState.weather?.temperature ?: "--°C",
            condition = uiState.weather?.condition ?: "Checking...",
            windSpeed = uiState.weather?.wind ?: "-- km/h",
            weatherType = when (uiState.weather?.condition?.lowercase()) {
                "sunny", "clear" -> WeatherType.SUNNY
                "rain", "showers", "drizzle" -> WeatherType.RAINY
                "cloudy", "overcast", "mist", "haze" -> WeatherType.CLOUDY
                "storm", "thunderstorm" -> WeatherType.STORMY
                else -> WeatherType.SUNNY
            },
            alertTitle = uiState.alerts.firstOrNull()?.title ?: stringResource(R.string.safe_zone),
            alertDesc = uiState.alerts.firstOrNull()?.description ?: stringResource(R.string.no_disease_outbreaks, uiState.location)
        )

        // 4. Featured Partners Section
        if (uiState.recommendations.isNotEmpty()) {
            AgriAdCarousel(partners = uiState.recommendations.map { AdPartner(it.id, it.title, it.imageUrl) })
        } else {
            // Default Fallback Partners
            val partners = listOf(
                AdPartner("1", "GreenEarth Agri", "https://images.unsplash.com/photo-1599424423956-6f81014ccbe0?q=80&w=2169&auto=format&fit=crop"),
                AdPartner("2", "EcoFarm Solutions", "https://images.unsplash.com/photo-1574323347407-f5e1ad6d020b?q=80&w=2187&auto=format&fit=crop"),
                AdPartner("3", "BioCure Genetics", "https://images.unsplash.com/photo-1628352081506-83c43123ed6d?q=80&w=2196&auto=format&fit=crop"),
                AdPartner("4", "AgriCorp Seeds", "https://images.unsplash.com/photo-1585314062340-f1a5a7c9328d?q=80&w=1000&auto=format&fit=crop")
            )
            AgriAdCarousel(partners = partners)
        }

        // 5. Recent Scans
        if (uiState.recentScans.isNotEmpty()) {
            AgriRecentScansRow(
                scans = uiState.recentScans.map {
                    RecentScanData(
                        id = it.id,
                        cropName = it.cropName,
                        status = it.status,
                        time = it.date,
                        imageUrl = it.imageUrl,
                        isHealthy = it.status.contains("Healthy", ignoreCase = true)
                    )
                },
                onScanClick = { /* TODO */ }
            )
        }

        // 6. Disease Prevention Recommendation
        uiState.recommendations.firstOrNull()?.let { rec ->
            AgriRecommendationCard(
                title = rec.title,
                description = rec.description,
                onClick = { /* TODO */ }
            )
        } ?: AgriRecommendationCard(
            title = stringResource(R.string.stay_proactive),
            description = stringResource(R.string.scan_regularly_prompt),
            onClick = { /* TODO */ }
        )
    }
}

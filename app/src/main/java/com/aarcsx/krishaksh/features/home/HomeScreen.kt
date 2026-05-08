package com.aarcsx.krishaksh.features.home

import android.Manifest
import android.os.Build
import androidx.activity.compose.rememberLauncherForActivityResult
import androidx.activity.result.contract.ActivityResultContracts
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.hilt.navigation.compose.hiltViewModel
import com.aarcsx.krishaksh.core.designsystem.components.*

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
    ) { _ -> }

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
    LazyColumn(
        modifier = Modifier.fillMaxSize(),
        verticalArrangement = Arrangement.spacedBy(8.dp),
        contentPadding = PaddingValues(bottom = 32.dp)
    ) {
        // App Title
        item {
            Text(
                text = "Krishaksh",
                modifier = Modifier.padding(start = 24.dp, top = 16.dp),
                style = MaterialTheme.typography.displaySmall,
                fontWeight = FontWeight.Black,
                color = ForestGreen,
                letterSpacing = (-1).sp
            )
        }

        // 1. Header Section
        item {
            AgriHeader(
                userName = uiState.userName,
                location = uiState.location,
                onNotificationClick = onNotificationClick
            )
        }

        // 2. Hero Scan Section
        item {
            AgriHeroSection(onScanClick = onScanClick)
        }

        // 3. Weather + Alert Section
        item {
            AgriWeatherAlertRow(
                temp = uiState.weather?.temperature ?: "28°C",
                condition = uiState.weather?.condition ?: "Sunny",
                windSpeed = "12 km/h",
                weatherType = WeatherType.SUNNY,
                alertTitle = uiState.alerts.firstOrNull()?.title ?: "No Alerts",
                alertDesc = uiState.alerts.firstOrNull()?.description ?: "Your area is safe."
            )
        }

        // 4. Recent Scans Section
        item {
            val scanItems = uiState.recentScans.map {
                RecentScanData(
                    id = it.id,
                    cropName = it.cropName,
                    status = it.status,
                    time = it.date,
                    imageUrl = it.imageUrl,
                    isHealthy = it.status == "Healthy"
                )
            }
            AgriRecentScansRow(
                scans = scanItems,
                onScanClick = { /* TODO: Navigate to detail */ }
            )
        }

        // 5. Recommendation Section
        item {
            AgriRecommendationCard(
                title = "Rice Blast Prevention",
                description = "Humidity is high. Use recommended fungicides to protect your crop.",
                onClick = { /* TODO */ }
            )
        }
    }
}

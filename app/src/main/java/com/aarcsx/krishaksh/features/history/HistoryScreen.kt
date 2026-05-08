package com.aarcsx.krishaksh.features.history

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import com.aarcsx.krishaksh.core.designsystem.components.*
import com.aarcsx.krishaksh.features.home.*

@Composable
fun HistoryScreen(onBackClick: () -> Unit = {}) {
    Scaffold(
        containerColor = Color(0xFFF5F7F2)
    ) { paddingValues ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(paddingValues)
        ) {
            AgriSubHeader(
                title = "Scan History",
                onBackClick = onBackClick
            )
            
            LazyColumn(
                modifier = Modifier.fillMaxSize(),
                contentPadding = PaddingValues(horizontal = 24.dp, vertical = 8.dp),
                verticalArrangement = Arrangement.spacedBy(16.dp)
            ) {
                items(10) {
                    AgriRecentScanCard(
                        cropName = "Crop $it",
                        status = if (it % 2 == 0) "Healthy" else "Diseased",
                        isHealthy = it % 2 == 0,
                        time = "May 08, 2026",
                        imageUrl = "",
                        onClick = { /* TODO */ }
                    )
                }
                
                // Add extra spacer at the bottom for better scroll feel
                item { Spacer(modifier = Modifier.height(80.dp)) }
            }
        }
    }
}

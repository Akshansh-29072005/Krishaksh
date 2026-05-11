package com.aarcsx.krisho.features.scan

import androidx.compose.foundation.*
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.grid.*
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.layout.ContentScale
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.aarcsx.krisho.R
import com.aarcsx.krisho.core.designsystem.theme.ForestGreen
import com.aarcsx.krisho.core.designsystem.components.AgriSubHeader

@Composable
fun CropSelectionScreen(
    onCropSelected: (String) -> Unit,
    onBackClick: () -> Unit,
    viewModel: ScanViewModel = androidx.hilt.navigation.compose.hiltViewModel()
) {
    val userCrops by viewModel.userCrops.collectAsState()
    var selectedCrop by remember { mutableStateOf("") }
    val defaultCrops = remember { listOf("Rice", "Brinjal", "Maize", "Tomato", "Potato", "Wheat") }
    val crops = userCrops?.ifEmpty { defaultCrops } ?: defaultCrops

    Scaffold(
        containerColor = Color(0xFFF5F7F2)
    ) { padding ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(padding)
        ) {
            AgriSubHeader(
                title = "Select Crop",
                onBackClick = onBackClick
            )

            Column(
                modifier = Modifier
                    .fillMaxSize()
                    .padding(horizontal = 24.dp)
            ) {
                Text(
                    "Which crop are you scanning today?",
                    style = MaterialTheme.typography.titleLarge,
                    fontWeight = FontWeight.Bold,
                    color = ForestGreen,
                    modifier = Modifier.padding(top = 24.dp, bottom = 8.dp)
                )

                Text(
                    "Select a crop to proceed with diagnosis",
                    style = MaterialTheme.typography.bodyMedium,
                    color = Color.Gray,
                    modifier = Modifier.padding(bottom = 24.dp)
                )

                LazyVerticalGrid(
                    columns = GridCells.Fixed(3),
                    horizontalArrangement = Arrangement.spacedBy(12.dp),
                    verticalArrangement = Arrangement.spacedBy(16.dp),
                    modifier = Modifier.weight(1f)
                ) {
                    items(crops) { crop ->
                        CropCard(
                            name = crop,
                            isSelected = selectedCrop == crop,
                            onClick = { selectedCrop = crop }
                        )
                    }
                }

                Button(
                    onClick = { onCropSelected(selectedCrop) },
                    enabled = selectedCrop.isNotEmpty(),
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(vertical = 24.dp)
                        .height(56.dp),
                    shape = RoundedCornerShape(16.dp),
                    colors = ButtonDefaults.buttonColors(containerColor = ForestGreen)
                ) {
                    Text("Start Scanning", fontWeight = FontWeight.Bold, fontSize = 18.sp)
                }
            }
        }
    }
}

@Composable
fun CropCard(
    name: String,
    isSelected: Boolean,
    onClick: () -> Unit
) {
    val emoji = when (name) {
        "Rice" -> "🌾"
        "Brinjal" -> "🍆"
        "Maize" -> "🌽"
        "Tomato" -> "🍅"
        "Potato" -> "🥔"
        "Wheat" -> "🌾"
        else -> "🌿"
    }

    Card(
        modifier = Modifier
            .fillMaxWidth()
            .clickable { onClick() },
        shape = RoundedCornerShape(24.dp),
        colors = CardDefaults.cardColors(
            containerColor = if (isSelected) ForestGreen.copy(alpha = 0.1f) else Color.White
        ),
        border = if (isSelected) BorderStroke(2.dp, ForestGreen) else null,
        elevation = CardDefaults.cardElevation(defaultElevation = 0.dp)
    ) {
        Column(
            modifier = Modifier
                .padding(16.dp)
                .fillMaxWidth(),
            horizontalAlignment = Alignment.CenterHorizontally,
            verticalArrangement = Arrangement.Center
        ) {
            Box(
                modifier = Modifier
                    .size(64.dp)
                    .background(Color(0xFFF8F9FA), CircleShape)
                    .border(1.dp, Color(0xFFEEEEEE), CircleShape),
                contentAlignment = Alignment.Center
            ) {
                Text(text = emoji, fontSize = 32.sp)
            }
            Spacer(modifier = Modifier.height(12.dp))
            Text(
                text = name,
                fontWeight = FontWeight.SemiBold,
                fontSize = 15.sp,
                color = if (isSelected) ForestGreen else Color(0xFF424242)
            )
        }
    }
}

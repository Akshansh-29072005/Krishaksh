package com.aarcsx.krishaksh.features.scan

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Chat
import androidx.compose.material.icons.filled.Support
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import com.aarcsx.krishaksh.core.designsystem.components.*
import com.aarcsx.krishaksh.core.designsystem.theme.ForestGreen
import androidx.compose.ui.graphics.SolidColor

import coil.compose.AsyncImage
import androidx.compose.ui.layout.ContentScale

@Composable
fun ResultScreen(
    result: ScanResult,
    onBackClick: () -> Unit
) {
    Scaffold(
        containerColor = Color(0xFFF5F7F2)
    ) { padding ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(padding)
        ) {
            AgriSubHeader(
                title = "Analysis Result",
                onBackClick = onBackClick
            )

            LazyColumn(
                modifier = Modifier.fillMaxSize(),
                contentPadding = PaddingValues(horizontal = 24.dp, vertical = 8.dp)
            ) {
                item {
                    Card(
                        modifier = Modifier.fillMaxWidth(),
                        shape = RoundedCornerShape(28.dp),
                        colors = CardDefaults.cardColors(containerColor = Color.White)
                    ) {
                        Column {
                            AsyncImage(
                                model = "file:///home/akshansh/.gemini/antigravity/brain/ebfaf9b6-f281-4ce2-a0dc-749993e3e660/potato_leaf_blight_1778254271532.png",
                                contentDescription = null,
                                modifier = Modifier
                                    .fillMaxWidth()
                                    .height(220.dp),
                                contentScale = ContentScale.Crop
                            )
                            Column(modifier = Modifier.padding(24.dp)) {
                                Text(
                                    text = result.diseaseName,
                                    style = MaterialTheme.typography.headlineMedium,
                                    color = ForestGreen,
                                    fontWeight = FontWeight.Bold
                                )
                                
                                val confidenceValue = result.confidence.replace("%", "").toDoubleOrNull() ?: 0.0
                                
                                Row(verticalAlignment = Alignment.CenterVertically) {
                                    Text(
                                        text = "Confidence: ${result.confidence}",
                                        style = MaterialTheme.typography.bodyMedium,
                                        color = if (confidenceValue < 60.0) Color(0xFFC62828) else ForestGreen.copy(alpha = 0.7f),
                                        fontWeight = if (confidenceValue < 60.0) FontWeight.Bold else FontWeight.Normal
                                    )
                                    if (confidenceValue < 60.0) {
                                        Spacer(modifier = Modifier.width(8.dp))
                                        Icon(
                                            imageVector = Icons.Default.Warning,
                                            contentDescription = "Low Confidence",
                                            tint = Color(0xFFC62828),
                                            modifier = Modifier.size(16.dp)
                                        )
                                    }
                                }

                                if (confidenceValue < 40.0) {
                                    Spacer(modifier = Modifier.height(8.dp))
                                    Surface(
                                        color = Color(0xFFFFEBEE),
                                        shape = RoundedCornerShape(8.dp)
                                    ) {
                                        Text(
                                            text = "Uncertain prediction. Please retake photo with better lighting.",
                                            modifier = Modifier.padding(8.dp),
                                            style = MaterialTheme.typography.labelSmall,
                                            color = Color(0xFFC62828)
                                        )
                                    }
                                }
                            }
                        }
                    }
                }

                item { Spacer(modifier = Modifier.height(24.dp)) }

                item {
                    ResultSection("Symptoms", result.symptoms)
                }

                item { Spacer(modifier = Modifier.height(24.dp)) }

                item {
                    ResultSection("Prevention", result.prevention)
                }

                item { Spacer(modifier = Modifier.height(24.dp)) }

                item {
                    ResultSection("Treatment", result.treatment)
                }

                item { Spacer(modifier = Modifier.height(32.dp)) }

                item {
                    Text(
                        "Recommended Product",
                        style = MaterialTheme.typography.titleMedium,
                        fontWeight = FontWeight.Bold,
                        color = ForestGreen
                    )
                    Spacer(modifier = Modifier.height(12.dp))
                    AgriRecommendationCard(
                        title = result.recommendationTitle,
                        description = result.recommendationDesc,
                        onClick = { /* TODO */ }
                    )
                }

                item { Spacer(modifier = Modifier.height(32.dp)) }

                item {
                    Row(
                        modifier = Modifier.fillMaxWidth(),
                        horizontalArrangement = Arrangement.spacedBy(16.dp)
                    ) {
                        Button(
                            onClick = { /* TODO */ },
                            modifier = Modifier.weight(1f).height(56.dp),
                            colors = ButtonDefaults.buttonColors(containerColor = ForestGreen),
                            shape = RoundedCornerShape(16.dp)
                        ) {
                            Icon(Icons.Default.Chat, contentDescription = null)
                            Spacer(modifier = Modifier.width(8.dp))
                            Text("Expert")
                        }
                        OutlinedButton(
                            onClick = { /* TODO */ },
                            modifier = Modifier.weight(1f).height(56.dp),
                            shape = RoundedCornerShape(16.dp),
                            colors = ButtonDefaults.outlinedButtonColors(contentColor = ForestGreen),
                            border = ButtonDefaults.outlinedButtonBorder.copy(brush = SolidColor(ForestGreen))
                        ) {
                            Icon(Icons.Default.Support, contentDescription = null)
                            Spacer(modifier = Modifier.width(8.dp))
                            Text("Support")
                        }
                    }
                }
                
                item { Spacer(modifier = Modifier.height(40.dp)) }
            }
        }
    }
}

@Composable
private fun ResultSection(title: String, content: String) {
    Column {
        Text(
            text = title,
            style = MaterialTheme.typography.titleMedium,
            fontWeight = FontWeight.Bold,
            color = ForestGreen
        )
        Spacer(modifier = Modifier.height(8.dp))
        Text(
            text = content,
            style = MaterialTheme.typography.bodyMedium,
            color = Color.Gray
        )
    }
}

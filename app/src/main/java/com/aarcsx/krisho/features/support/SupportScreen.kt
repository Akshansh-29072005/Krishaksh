package com.aarcsx.krisho.features.support

import android.Manifest
import androidx.activity.compose.rememberLauncherForActivityResult
import androidx.activity.result.contract.ActivityResultContracts
import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.SolidColor
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.hilt.navigation.compose.hiltViewModel
import com.aarcsx.krisho.core.designsystem.components.*

@Composable
fun SupportScreen(
    onBackClick: () -> Unit = {},
    onGoToProfile: () -> Unit = {},
    viewModel: SupportViewModel = hiltViewModel()
) {
    val uiState by viewModel.uiState.collectAsState()

    val permissionLauncher = rememberLauncherForActivityResult(
        ActivityResultContracts.RequestPermission()
    ) { isGranted ->
        if (isGranted) {
            viewModel.startRecording()
        }
    }

    Scaffold(
        containerColor = Color(0xFFF5F7F2)
    ) { paddingValues ->
        LazyColumn(
            modifier = Modifier
                .fillMaxSize()
                .padding(paddingValues)
        ) {
            item {
                AgriSubHeader(
                    title = "Help & Support",
                    onBackClick = onBackClick
                )
            }

            item {
                Column(
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(horizontal = 24.dp)
                        .padding(top = 24.dp),
                    verticalArrangement = Arrangement.spacedBy(24.dp)
                ) {
                    Card(
                        modifier = Modifier.fillMaxWidth(),
                        shape = RoundedCornerShape(28.dp),
                        colors = CardDefaults.cardColors(containerColor = Color.White)
                    ) {
                        Column(modifier = Modifier.padding(24.dp)) {
                            Text(
                                "How can we help you today?",
                                style = MaterialTheme.typography.titleMedium,
                                fontWeight = FontWeight.Bold,
                                color = ForestGreen
                            )
                            Spacer(modifier = Modifier.height(16.dp))
                            OutlinedTextField(
                                value = uiState.message,
                                onValueChange = { viewModel.onMessageChange(it) },
                                placeholder = { Text("Describe your issue...") },
                                modifier = Modifier
                                    .fillMaxWidth()
                                    .height(140.dp),
                                shape = RoundedCornerShape(16.dp),
                                colors = OutlinedTextFieldDefaults.colors(
                                    focusedBorderColor = ForestGreen,
                                    cursorColor = ForestGreen,
                                    focusedTextColor = Color.Black,
                                    unfocusedTextColor = Color.Black.copy(alpha = 0.8f)
                                )
                            )

                            if (uiState.audioFile != null) {
                                Spacer(modifier = Modifier.height(8.dp))
                                Row(verticalAlignment = Alignment.CenterVertically) {
                                    Icon(Icons.Default.Mic, contentDescription = null, tint = ForestGreen)
                                    Spacer(modifier = Modifier.width(8.dp))
                                    Text(
                                        "Voice message recorded (${uiState.recordingDuration}s)",
                                        style = MaterialTheme.typography.bodySmall,
                                        color = ForestGreen
                                    )
                                }
                            }

                            Spacer(modifier = Modifier.height(24.dp))
                            Row(horizontalArrangement = Arrangement.spacedBy(16.dp)) {
                                Button(
                                    onClick = { viewModel.sendSupportRequest() },
                                    modifier = Modifier
                                        .weight(1f)
                                        .height(56.dp),
                                    shape = RoundedCornerShape(16.dp),
                                    enabled = (uiState.message.isNotBlank() || uiState.audioFile != null) && !uiState.isSending,
                                    colors = ButtonDefaults.buttonColors(containerColor = ForestGreen)
                                ) {
                                    if (uiState.isSending) {
                                        CircularProgressIndicator(
                                            modifier = Modifier.size(24.dp),
                                            color = Color.White,
                                            strokeWidth = 2.dp
                                        )
                                    } else {
                                        Icon(Icons.Default.Send, contentDescription = null)
                                        Spacer(modifier = Modifier.width(8.dp))
                                        Text("Send")
                                    }
                                }

                                val micColor = if (uiState.isRecording) Color.Red else ForestGreen
                                OutlinedButton(
                                    onClick = {
                                        if (uiState.isRecording) {
                                            viewModel.stopRecording()
                                        } else {
                                            permissionLauncher.launch(Manifest.permission.RECORD_AUDIO)
                                        }
                                    },
                                    modifier = Modifier
                                        .weight(1f)
                                        .height(56.dp),
                                    shape = RoundedCornerShape(16.dp),
                                    colors = ButtonDefaults.outlinedButtonColors(contentColor = micColor),
                                    border = ButtonDefaults.outlinedButtonBorder.copy(
                                        brush = SolidColor(micColor)
                                    )
                                ) {
                                    Icon(
                                        if (uiState.isRecording) Icons.Default.Stop else Icons.Default.Mic,
                                        contentDescription = null
                                    )
                                    Spacer(modifier = Modifier.width(8.dp))
                                    Text(if (uiState.isRecording) "Stop" else "Voice")
                                }
                            }
                        }
                    }

                    if (uiState.error != null) {
                        Card(
                            modifier = Modifier.fillMaxWidth(),
                            shape = RoundedCornerShape(16.dp),
                            colors = CardDefaults.cardColors(containerColor = Color(0xFFFFEBEE))
                        ) {
                            Row(
                                modifier = Modifier.padding(16.dp),
                                verticalAlignment = Alignment.CenterVertically
                            ) {
                                Icon(Icons.Default.Error, contentDescription = null, tint = Color(0xFFD32F2F))
                                Spacer(modifier = Modifier.width(12.dp))
                                Text(
                                    uiState.error!!,
                                    style = MaterialTheme.typography.bodySmall,
                                    color = Color(0xFFD32F2F)
                                )
                            }
                        }
                    }

                    if (uiState.success) {
                        Card(
                            modifier = Modifier.fillMaxWidth(),
                            shape = RoundedCornerShape(16.dp),
                            colors = CardDefaults.cardColors(containerColor = ForestGreen.copy(alpha = 0.1f))
                        ) {
                            Row(
                                modifier = Modifier.padding(16.dp),
                                verticalAlignment = Alignment.CenterVertically
                            ) {
                                Icon(Icons.Default.CheckCircle, contentDescription = null, tint = ForestGreen)
                                Spacer(modifier = Modifier.width(12.dp))
                                Text(
                                    "Support request sent successfully!",
                                    style = MaterialTheme.typography.bodySmall,
                                    color = ForestGreen,
                                    fontWeight = FontWeight.Medium
                                )
                            }
                        }
                    }
                }
            }

            if (uiState.tickets.isNotEmpty()) {
                item {
                    Spacer(modifier = Modifier.height(24.dp))
                    Text(
                        "Your Support Tickets",
                        style = MaterialTheme.typography.titleMedium,
                        fontWeight = FontWeight.Bold,
                        color = ForestGreen,
                        modifier = Modifier.padding(horizontal = 24.dp)
                    )
                    Spacer(modifier = Modifier.height(12.dp))
                }

                items(uiState.tickets) { ticket ->
                    TicketCard(
                        ticket = ticket,
                        canRequestCallback = uiState.userPhone?.isNotBlank() == true,
                        onCallbackClick = { viewModel.requestCallback(ticket.id) },
                        modifier = Modifier.padding(horizontal = 24.dp, vertical = 8.dp)
                    )
                }

                item {
                    Spacer(modifier = Modifier.height(32.dp))
                }
            }
        }
    }

    // Success dialog
    if (uiState.success) {
        LaunchedEffect(uiState.success) {
            viewModel.clearSuccess()
        }
    }

    if (uiState.showPhoneRequiredDialog) {
        AlertDialog(
            onDismissRequest = { viewModel.clearPhoneRequiredDialog() },
            title = { Text("Phone number required") },
            text = {
                Text("You must add a phone number in your profile before requesting a callback.")
            },
            confirmButton = {
                TextButton(onClick = {
                    viewModel.clearPhoneRequiredDialog()
                    onGoToProfile()
                }) {
                    Text("Update Phone")
                }
            },
            dismissButton = {
                TextButton(onClick = { viewModel.clearPhoneRequiredDialog() }) {
                    Text("Cancel")
                }
            }
        )
    }
}

@Composable
fun TicketCard(
    ticket: com.aarcsx.krisho.core.network.dto.SupportTicketDto,
    canRequestCallback: Boolean,
    onCallbackClick: () -> Unit,
    modifier: Modifier = Modifier
) {
    Card(
        modifier = modifier.fillMaxWidth(),
        shape = RoundedCornerShape(16.dp),
        colors = CardDefaults.cardColors(containerColor = Color.White),
        elevation = CardDefaults.cardElevation(defaultElevation = 4.dp)
    ) {
        Column(modifier = Modifier.padding(16.dp)) {
            // Header: Title and Status
            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.SpaceBetween,
                verticalAlignment = Alignment.CenterVertically
            ) {
                Text(
                    ticket.title,
                    style = MaterialTheme.typography.titleSmall,
                    fontWeight = FontWeight.Bold,
                    modifier = Modifier.weight(1f)
                )
                
                val statusColor = when (ticket.status) {
                    "open" -> Color(0xFFFFA500)
                    "resolved" -> Color(0xFF4CAF50)
                    "closed" -> Color(0xFF9E9E9E)
                    else -> Color(0xFF2196F3)
                }
                
                Surface(
                    shape = RoundedCornerShape(8.dp),
                    color = statusColor.copy(alpha = 0.2f)
                ) {
                    Text(
                        ticket.status.uppercase(),
                        style = MaterialTheme.typography.labelSmall,
                        color = statusColor,
                        modifier = Modifier.padding(horizontal = 8.dp, vertical = 4.dp),
                        fontWeight = FontWeight.SemiBold
                    )
                }
            }

            Spacer(modifier = Modifier.height(8.dp))

            // Description
            Text(
                ticket.description,
                style = MaterialTheme.typography.bodySmall,
                color = Color.Gray,
                maxLines = 2,
                modifier = Modifier.fillMaxWidth()
            )

            Spacer(modifier = Modifier.height(12.dp))

            // Callback Status
            if (ticket.callback_requested && canRequestCallback) {
                Row(
                    modifier = Modifier
                        .fillMaxWidth()
                        .background(
                            color = ForestGreen.copy(alpha = 0.1f),
                            shape = RoundedCornerShape(8.dp)
                        )
                        .padding(8.dp),
                    verticalAlignment = Alignment.CenterVertically,
                    horizontalArrangement = Arrangement.spacedBy(8.dp)
                ) {
                    Icon(
                        Icons.Default.Phone,
                        contentDescription = null,
                        tint = ForestGreen,
                        modifier = Modifier.size(16.dp)
                    )
                    Text(
                        "Callback: ${ticket.callback_status}",
                        style = MaterialTheme.typography.labelSmall,
                        color = ForestGreen,
                        fontWeight = FontWeight.Medium
                    )
                }
            } else {
                Button(
                    onClick = onCallbackClick,
                    modifier = Modifier
                        .fillMaxWidth()
                        .height(40.dp),
                    shape = RoundedCornerShape(8.dp),
                    colors = ButtonDefaults.buttonColors(
                        containerColor = ForestGreen.copy(alpha = 0.2f),
                        contentColor = ForestGreen
                    )
                ) {
                    Icon(Icons.Default.Call, contentDescription = null, modifier = Modifier.size(16.dp))
                    Spacer(modifier = Modifier.width(6.dp))
                    Text("Request Callback", style = MaterialTheme.typography.labelSmall)
                }
            }
        }
    }
}

package com.aarcsx.krishaksh.features.support

import android.Manifest
import androidx.activity.compose.rememberLauncherForActivityResult
import androidx.activity.result.contract.ActivityResultContracts
import androidx.compose.foundation.layout.*
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
import com.aarcsx.krishaksh.core.designsystem.components.*

@Composable
fun SupportScreen(
    onBackClick: () -> Unit = {},
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
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(paddingValues)
        ) {
            AgriSubHeader(
                title = "Help & Support",
                onBackClick = onBackClick
            )

            Column(
                modifier = Modifier
                    .fillMaxSize()
                    .padding(horizontal = 24.dp),
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
                                cursorColor = ForestGreen
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

                if (uiState.isCallbackRequested) {
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
                                "Callback requested! We will call you soon.",
                                style = MaterialTheme.typography.bodyMedium,
                                color = ForestGreen,
                                fontWeight = FontWeight.Medium
                            )
                        }
                    }
                } else {
                    var showPhoneDialog by remember { mutableStateOf(false) }
                    
                    Button(
                        onClick = { showPhoneDialog = true },
                        modifier = Modifier
                            .fillMaxWidth()
                            .height(64.dp),
                        shape = RoundedCornerShape(16.dp),
                        colors = ButtonDefaults.buttonColors(
                            containerColor = ForestGreen.copy(alpha = 0.1f),
                            contentColor = ForestGreen
                        )
                    ) {
                        Icon(Icons.Default.Call, contentDescription = null)
                        Spacer(modifier = Modifier.width(12.dp))
                        Text("Request Callback", fontWeight = FontWeight.SemiBold)
                    }

                    if (showPhoneDialog) {
                        var phone by remember { mutableStateOf("") }
                        AlertDialog(
                            onDismissRequest = { showPhoneDialog = false },
                            title = { Text("Enter Phone Number") },
                            text = {
                                OutlinedTextField(
                                    value = phone,
                                    onValueChange = { if (it.all { char -> char.isDigit() }) phone = it },
                                    label = { Text("Phone Number") },
                                    singleLine = true,
                                    modifier = Modifier.fillMaxWidth(),
                                    colors = OutlinedTextFieldDefaults.colors(
                                        focusedBorderColor = ForestGreen,
                                        focusedLabelColor = ForestGreen
                                    )
                                )
                            },
                            confirmButton = {
                                TextButton(
                                    onClick = {
                                        if (phone.length >= 10) {
                                            viewModel.requestCallback(phone)
                                            showPhoneDialog = false
                                        }
                                    }
                                ) {
                                    Text("Request", color = ForestGreen)
                                }
                            },
                            dismissButton = {
                                TextButton(onClick = { showPhoneDialog = false }) {
                                    Text("Cancel")
                                }
                            }
                        )
                    }
                }
            }
        }
    }

    // Success dialog
    if (uiState.success) {
        AlertDialog(
            onDismissRequest = { viewModel.clearSuccess() },
            confirmButton = {
                TextButton(onClick = { viewModel.clearSuccess() }) {
                    Text("OK", color = ForestGreen)
                }
            },
            title = { Text("Success") },
            text = { Text("Your support request has been sent. We'll get back to you soon.") }
        )
    }

    // Error message
    if (uiState.error != null) {
        LaunchedEffect(uiState.error) {
            // Show snackbar or toast
            viewModel.clearError()
        }
    }
}

package com.aarcsx.krishaksh.features.scan

import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.ArrowBack
import androidx.compose.material.icons.filled.Camera
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import com.aarcsx.krishaksh.core.designsystem.theme.ForestGreen

import android.content.Context
import androidx.camera.core.CameraSelector
import androidx.camera.core.Preview
import androidx.camera.lifecycle.ProcessCameraProvider
import androidx.camera.view.PreviewView
import androidx.compose.ui.viewinterop.AndroidView
import androidx.core.content.ContextCompat
import androidx.lifecycle.compose.LocalLifecycleOwner
import androidx.compose.ui.platform.LocalContext

@Composable
fun ScanScreen(
    onBackClick: () -> Unit,
    onCaptured: () -> Unit,
    viewModel: ScanViewModel = androidx.hilt.navigation.compose.hiltViewModel()
) {
    val context = LocalContext.current
    val lifecycleOwner = LocalLifecycleOwner.current

    var showCropSelection by remember { mutableStateOf(false) }
    var selectedCrop by remember { mutableStateOf("") }
    val userCrops by viewModel.userCrops.collectAsState()
    val defaultCrops = listOf("Wheat", "Rice", "Tomato", "Potato", "Cotton", "Mustard")
    val crops = userCrops.ifEmpty { defaultCrops }

    Box(modifier = Modifier.fillMaxSize().background(Color.Black)) {
        // Camera Preview Live Feed
        AndroidView(
            factory = { ctx ->
                val previewView = PreviewView(ctx)
                val cameraProviderFuture = ProcessCameraProvider.getInstance(ctx)
                cameraProviderFuture.addListener({
                    val cameraProvider = cameraProviderFuture.get()
                    val preview = Preview.Builder().build().also {
                        it.setSurfaceProvider(previewView.surfaceProvider)
                    }
                    val cameraSelector = CameraSelector.DEFAULT_BACK_CAMERA

                    try {
                        cameraProvider.unbindAll()
                        cameraProvider.bindToLifecycle(
                            lifecycleOwner, cameraSelector, preview
                        )
                    } catch(exc: Exception) {
                        // Log or handle camera binding failure
                    }
                }, ContextCompat.getMainExecutor(ctx))
                previewView
            },
            modifier = Modifier.fillMaxSize()
        )

        // Overlays
        IconButton(
            onClick = onBackClick,
            modifier = Modifier.padding(16.dp).align(Alignment.TopStart)
        ) {
            Icon(Icons.Default.ArrowBack, contentDescription = "Back", tint = Color.White)
        }

        Column(
            modifier = Modifier
                .align(Alignment.BottomCenter)
                .padding(bottom = 48.dp),
            horizontalAlignment = Alignment.CenterHorizontally
        ) {
            Text(
                "Point camera at the crop leaf",
                color = Color.White,
                modifier = Modifier.padding(bottom = 24.dp)
            )

            Surface(
                onClick = {
                    // Trigger crop selection dialog before analyzing
                    showCropSelection = true
                },
                shape = CircleShape,
                color = Color.White,
                modifier = Modifier.size(80.dp)
            ) {
                Box(contentAlignment = Alignment.Center) {
                    Icon(
                        Icons.Default.Camera,
                        contentDescription = "Capture",
                        tint = ForestGreen,
                        modifier = Modifier.size(32.dp)
                    )
                }
            }
        }

        // Crop Selection Dialog
        if (showCropSelection) {
            AlertDialog(
                onDismissRequest = { showCropSelection = false },
                title = { Text("Select Crop Category") },
                text = {
                    Column {
                        Text("Which crop did you just scan?")
                        Spacer(modifier = Modifier.height(16.dp))
                        crops.forEach { crop ->
                            Row(
                                modifier = Modifier
                                    .fillMaxWidth()
                                    .clickable { selectedCrop = crop }
                                    .padding(vertical = 8.dp),
                                verticalAlignment = Alignment.CenterVertically
                            ) {
                                RadioButton(
                                    selected = selectedCrop == crop,
                                    onClick = { selectedCrop = crop },
                                    colors = RadioButtonDefaults.colors(selectedColor = ForestGreen)
                                )
                                Spacer(modifier = Modifier.width(8.dp))
                                Text(text = crop)
                            }
                        }
                    }
                },
                confirmButton = {
                    Button(
                        onClick = {
                            showCropSelection = false
                            onCaptured() // Proceed to LoadingScreen (API placeholder)
                        },
                        colors = ButtonDefaults.buttonColors(containerColor = ForestGreen),
                        enabled = selectedCrop.isNotEmpty()
                    ) {
                        Text("Upload & Analyze")
                    }
                },
                dismissButton = {
                    TextButton(onClick = { showCropSelection = false }) {
                        Text("Cancel", color = Color.Gray)
                    }
                }
            )
        }
    }
}

@Composable
fun LoadingScreen() {
    Box(modifier = Modifier.fillMaxSize().background(ForestGreen), contentAlignment = Alignment.Center) {
        Column(horizontalAlignment = Alignment.CenterHorizontally) {
            CircularProgressIndicator(color = Color.White)
            Spacer(modifier = Modifier.height(24.dp))
            Text(
                "AI is analyzing your crop...",
                color = Color.White,
                style = MaterialTheme.typography.titleMedium
            )
        }
    }
}

package com.aarcsx.krisho.features.scan

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
import androidx.compose.ui.zIndex
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import com.aarcsx.krisho.core.designsystem.theme.ForestGreen

import android.content.Context
import android.graphics.Bitmap
import android.graphics.BitmapFactory
import android.net.Uri
import androidx.compose.ui.graphics.asImageBitmap
import androidx.camera.core.*
import androidx.camera.lifecycle.ProcessCameraProvider
import androidx.camera.view.PreviewView
import androidx.compose.ui.viewinterop.AndroidView
import androidx.core.content.ContextCompat
import androidx.lifecycle.compose.LocalLifecycleOwner
import androidx.compose.ui.platform.LocalContext
import java.io.ByteArrayOutputStream
import java.io.File
import java.util.concurrent.ExecutorService
import java.util.concurrent.Executors

@Composable
fun ScanScreen(
    cropName: String,
    onBackClick: () -> Unit,
    onCaptured: (ScanResult) -> Unit,
    viewModel: ScanViewModel = androidx.hilt.navigation.compose.hiltViewModel()
) {
    val context = LocalContext.current
    val lifecycleOwner = LocalLifecycleOwner.current

    val uiState by viewModel.uiState.collectAsState()
    var capturedImage by remember { mutableStateOf<Bitmap?>(null) }
    var imageCapture: ImageCapture? by remember { mutableStateOf(null) }

    val cameraExecutor = remember { Executors.newSingleThreadExecutor() }

    LaunchedEffect(uiState.result) {
        uiState.result?.let {
            onCaptured(it)
            viewModel.reset()
        }
    }

    Box(modifier = Modifier.fillMaxSize()) {
        if (uiState.isAnalyzing) {
            Box(
                modifier = Modifier
                    .fillMaxSize()
                    .background(ForestGreen)
                    .zIndex(10f),
                contentAlignment = Alignment.Center
            ) {
                LoadingScreen()
            }
        } else {
            if (capturedImage == null) {
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

                            imageCapture = ImageCapture.Builder()
                                .setCaptureMode(ImageCapture.CAPTURE_MODE_MINIMIZE_LATENCY)
                                .build()

                            val cameraSelector = CameraSelector.DEFAULT_BACK_CAMERA

                            try {
                                cameraProvider.unbindAll()
                                cameraProvider.bindToLifecycle(
                                    lifecycleOwner, cameraSelector, preview, imageCapture
                                )
                            } catch (exc: Exception) {
                                // Handle camera binding failure
                            }
                        }, ContextCompat.getMainExecutor(ctx))
                        previewView
                    },
                    modifier = Modifier.fillMaxSize()
                )
            } else {
                // Captured Image Overlay
                androidx.compose.foundation.Image(
                    bitmap = capturedImage!!.asImageBitmap(),
                    contentDescription = "Captured Image",
                    modifier = Modifier.fillMaxSize(),
                    contentScale = androidx.compose.ui.layout.ContentScale.Crop
                )
            }
        }

        // Overlays
        if (!uiState.isAnalyzing) {
            IconButton(
                onClick = onBackClick,
                modifier = Modifier.padding(16.dp).align(Alignment.TopStart)
            ) {
                Icon(Icons.Default.ArrowBack, contentDescription = "Back", tint = Color.White)
            }
        }

        if (capturedImage == null && !uiState.isAnalyzing) {
            Column(
                modifier = Modifier
                    .align(Alignment.BottomCenter)
                    .padding(bottom = 48.dp),
                horizontalAlignment = Alignment.CenterHorizontally
            ) {
                Text(
                    "Point camera at the $cropName leaf",
                    color = Color.White,
                    modifier = Modifier.padding(bottom = 24.dp)
                )

                Surface(
                    onClick = {
                        val capture = imageCapture ?: return@Surface
                        capture.takePicture(
                            ContextCompat.getMainExecutor(context),
                            object : ImageCapture.OnImageCapturedCallback() {
                                override fun onCaptureSuccess(image: ImageProxy) {
                                    val buffer = image.planes[0].buffer
                                    val bytes = ByteArray(buffer.remaining())
                                    buffer.get(bytes)
                                    capturedImage =
                                        BitmapFactory.decodeByteArray(bytes, 0, bytes.size)
                                    image.close()
                                }

                                override fun onError(exc: ImageCaptureException) {
                                    // Handle error
                                }
                            }
                        )
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
        }
        
        // Captured image URI is handled in ScanViewModel, but we show confirmation buttons here
        if (!uiState.isAnalyzing && capturedImage != null) {
            // Error overlay
            uiState.error?.let { error ->
                Box(
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(horizontal = 32.dp)
                        .padding(bottom = 120.dp)
                        .background(Color.Red.copy(alpha = 0.8f), RoundedCornerShape(8.dp))
                        .padding(8.dp)
                        .align(Alignment.BottomCenter)
                ) {
                    Text(error, color = Color.White, style = MaterialTheme.typography.bodySmall)
                }
            }

            // Confirmation Buttons
            Row(
                modifier = Modifier
                    .align(Alignment.BottomCenter)
                    .padding(bottom = 48.dp)
                    .fillMaxWidth()
                    .padding(horizontal = 32.dp),
                horizontalArrangement = Arrangement.spacedBy(16.dp)
            ) {
                Button(
                    onClick = { capturedImage = null },
                    modifier = Modifier.weight(1f).height(56.dp),
                    colors = ButtonDefaults.buttonColors(containerColor = Color.DarkGray),
                    shape = RoundedCornerShape(16.dp)
                ) {
                    Text("Cancel")
                }
                Button(
                    onClick = {
                        capturedImage?.let { bitmap ->
                            val stream = ByteArrayOutputStream()
                            bitmap.compress(Bitmap.CompressFormat.JPEG, 90, stream)
                            viewModel.analyzeImage(stream.toByteArray(), cropName)
                        }
                    },
                    modifier = Modifier.weight(1f).height(56.dp),
                    colors = ButtonDefaults.buttonColors(containerColor = ForestGreen),
                    shape = RoundedCornerShape(16.dp)
                ) {
                    Text("Analyze")
                }
            }
        }
    }
}

@Composable
fun LoadingScreen() {
    Box(
        modifier = Modifier
            .fillMaxSize()
            .background(ForestGreen),
        contentAlignment = Alignment.Center
    ) {
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

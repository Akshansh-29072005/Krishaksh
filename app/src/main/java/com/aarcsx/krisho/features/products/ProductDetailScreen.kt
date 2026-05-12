package com.aarcsx.krisho.features.products

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.verticalScroll
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.layout.ContentScale
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import coil.compose.AsyncImage
import com.aarcsx.krisho.core.designsystem.components.*
import com.aarcsx.krisho.core.models.Product

import androidx.hilt.navigation.compose.hiltViewModel
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.runtime.LaunchedEffect
import com.aarcsx.krisho.core.local.room.entity.ProductEntity

import android.app.Activity

@Composable
fun ProductDetailScreen(
    productId: String,
    onBackClick: () -> Unit,
    viewModel: ProductDetailViewModel = hiltViewModel()
) {
    val uiState by viewModel.uiState.collectAsState()
    val context = androidx.compose.ui.platform.LocalContext.current
    val activity = context as? Activity

    LaunchedEffect(productId) {
        viewModel.loadProduct(productId)
    }

    if (uiState.orderInitiated) {
        AlertDialog(
            onDismissRequest = { /* Handle dismiss */ },
            confirmButton = {
                Button(onClick = onBackClick) {
                    Text("OK")
                }
            },
            title = { Text("Payment Successful!") },
            text = { Text("Your payment was successful (ID: ${uiState.paymentSuccessId}). Your order for ${uiState.product?.name} has been placed.") },
            containerColor = WarmWhite,
            titleContentColor = ForestGreen,
            textContentColor = Color.DarkGray
        )
    }

    if (uiState.paymentError != null) {
        AlertDialog(
            onDismissRequest = { viewModel.clearPaymentError() },
            confirmButton = {
                Button(onClick = { viewModel.clearPaymentError() }) {
                    Text("Retry")
                }
            },
            title = { Text("Payment Failed") },
            text = { Text(uiState.paymentError!!) },
            containerColor = WarmWhite,
            titleContentColor = Color.Red,
            textContentColor = Color.DarkGray
        )
    }

    Scaffold(
        containerColor = WarmWhite,
        bottomBar = {
            if (uiState.product != null) {
                Surface(
                    modifier = Modifier.fillMaxWidth(),
                    shadowElevation = 8.dp,
                    color = Color.White
                ) {
                    Button(
                        onClick = { activity?.let { viewModel.onBuyNow(it) } },
                        modifier = Modifier
                            .fillMaxWidth()
                            .padding(24.dp)
                            .height(56.dp),
                        shape = RoundedCornerShape(16.dp),
                        colors = ButtonDefaults.buttonColors(containerColor = ForestGreen)
                    ) {
                        Text("Buy Now", fontWeight = FontWeight.Bold, fontSize = 18.sp)
                    }
                }
            }
        }
    ) { paddingValues ->
        Box(modifier = Modifier.fillMaxSize().padding(paddingValues)) {
            when {
                uiState.isLoading -> {
                    CircularProgressIndicator(
                        modifier = Modifier.align(Alignment.Center),
                        color = ForestGreen
                    )
                }
                uiState.error != null -> {
                    Text(
                        text = uiState.error!!,
                        modifier = Modifier.align(Alignment.Center),
                        color = Color.Red
                    )
                }
                uiState.product != null -> {
                    val product = uiState.product!!
                    Column(
                        modifier = Modifier
                            .fillMaxSize()
                            .verticalScroll(rememberScrollState())
                    ) {
                        AgriSubHeader(
                            title = "Product Details",
                            onBackClick = onBackClick
                        )

                        AsyncImage(
                            model = product.imageUrl,
                            contentDescription = null,
                            modifier = Modifier
                                .fillMaxWidth()
                                .height(300.dp)
                                .padding(horizontal = 24.dp)
                                .clip(RoundedCornerShape(28.dp)),
                            contentScale = ContentScale.Crop
                        )

                        Column(
                            modifier = Modifier.padding(24.dp)
                        ) {
                            Row(
                                modifier = Modifier.fillMaxWidth(),
                                horizontalArrangement = Arrangement.SpaceBetween,
                                verticalAlignment = Alignment.Top
                            ) {
                                Column(modifier = Modifier.weight(1f)) {
                                    Text(
                                        text = product.name,
                                        style = MaterialTheme.typography.headlineMedium,
                                        fontWeight = FontWeight.Bold,
                                        color = ForestGreen
                                    )
                                    Text(
                                        text = "by ${product.companyName}",
                                        style = MaterialTheme.typography.bodyMedium,
                                        color = Color.Gray
                                    )
                                }
                                Text(
                                    text = "₹${product.price}",
                                    style = MaterialTheme.typography.headlineSmall,
                                    fontWeight = FontWeight.ExtraBold,
                                    color = ForestGreen
                                )
                            }

                            Spacer(Modifier.height(32.dp))

                            Text(
                                text = "Description",
                                style = MaterialTheme.typography.titleLarge,
                                fontWeight = FontWeight.Bold,
                                color = ForestGreen
                            )
                            Spacer(Modifier.height(8.dp))
                            Text(
                                text = product.description,
                                style = MaterialTheme.typography.bodyLarge,
                                lineHeight = 24.sp,
                                color = Color.DarkGray
                            )

                            Spacer(Modifier.height(24.dp))

                            if (product.usageInstructions.isNotEmpty()) {
                                Card(
                                    modifier = Modifier.fillMaxWidth(),
                                    shape = RoundedCornerShape(20.dp),
                                    colors = CardDefaults.cardColors(containerColor = Color(0xFFE8F5E9))
                                ) {
                                    Column(Modifier.padding(20.dp)) {
                                        Text(
                                            text = "Unit / Instructions",
                                            fontWeight = FontWeight.Bold,
                                            color = ForestGreen
                                        )
                                        Spacer(Modifier.height(8.dp))
                                        Text(
                                            text = product.usageInstructions,
                                            style = MaterialTheme.typography.bodyMedium,
                                            color = ForestGreenLight
                                        )
                                    }
                                }
                            }
                            
                            Spacer(Modifier.height(32.dp))
                        }
                    }
                }
            }
        }
    }
}

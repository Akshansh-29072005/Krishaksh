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

@Composable
fun ProductDetailScreen(
    productId: String,
    onBackClick: () -> Unit
) {
    // Mock product lookup - in a real app this would come from a database/API via ViewModel
    val product = getMockProduct(productId)

    Scaffold(
        containerColor = WarmWhite,
        bottomBar = {
            Surface(
                modifier = Modifier.fillMaxWidth(),
                shadowElevation = 8.dp,
                color = Color.White
            ) {
                Button(
                    onClick = { /* Buy Now */ },
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
    ) { paddingValues ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(paddingValues)
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
                            text = "by ${product.company}",
                            style = MaterialTheme.typography.bodyMedium,
                            color = Color.Gray
                        )
                    }
                    Text(
                        text = product.price,
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

                Card(
                    modifier = Modifier.fillMaxWidth(),
                    shape = RoundedCornerShape(20.dp),
                    colors = CardDefaults.cardColors(containerColor = Color(0xFFE8F5E9))
                ) {
                    Column(Modifier.padding(20.dp)) {
                        Text(
                            text = "Why use this?",
                            fontWeight = FontWeight.Bold,
                            color = ForestGreen
                        )
                        Spacer(Modifier.height(8.dp))
                        Text(
                            text = product.usage,
                            style = MaterialTheme.typography.bodyMedium,
                            color = ForestGreenLight
                        )
                    }
                }
                
                Spacer(Modifier.height(32.dp))
            }
        }
    }
}

private fun getMockProduct(id: String): Product {
    return Product(
        id = id,
        name = "Organic Growth Booster",
        price = "₹499.00",
        description = "A premium organic fertilizer derived from seaweed extract and composted organic matter. It provides essential micronutrients and promotes healthy root development.",
        usage = "Mix 20ml per liter of water. Apply once every 15 days during the growing season for best results.",
        company = "GreenEarth Agri Solutions",
        imageUrl = "https://images.unsplash.com/photo-1585314062340-f1a5a7c9328d?q=80&w=1000&auto=format&fit=crop"
    )
}

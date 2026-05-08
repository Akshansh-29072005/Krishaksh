package com.aarcsx.krishaksh.features.products

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.grid.GridCells
import androidx.compose.foundation.lazy.grid.LazyVerticalGrid
import androidx.compose.foundation.lazy.grid.items
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.unit.dp
import com.aarcsx.krishaksh.core.designsystem.components.*
import com.aarcsx.krishaksh.core.models.Product

@Composable
fun ProductsScreen(
    onBackClick: () -> Unit = {},
    onProductClick: (String) -> Unit = {}
) {
    var searchQuery by remember { mutableStateOf("") }
    
    // Mock products list
    val allProducts = remember {
        listOf(
            Product("1", "Organic Growth Booster", "₹499", "Premium organic fertilizer", "Mix 20ml/L", "GreenEarth", "https://images.unsplash.com/photo-1585314062340-f1a5a7c9328d?q=80&w=1000&auto=format&fit=crop"),
            Product("2", "Neem Oil Spray", "₹299", "Natural pest control", "Spray on leaves", "EcoFarm", "https://images.unsplash.com/photo-1615485290382-441e4d0c9cb5?q=80&w=1000&auto=format&fit=crop"),
            Product("3", "NPK 19:19:19", "₹150", "Balanced water soluble fertilizer", "Soil application", "AgriCorp", "https://images.unsplash.com/photo-1592982537447-6f2a6a0c3c1b?q=80&w=1000&auto=format&fit=crop"),
            Product("4", "Bio-Fungicide", "₹350", "Prevent fungal infections", "Seed treatment", "BioCure", "https://images.unsplash.com/photo-1558449028-b53a39d100fc?q=80&w=1000&auto=format&fit=crop"),
            Product("5", "Seed Starter Mix", "₹199", "Premium soil for seeds", "Nursery use", "HortiPlus", "https://images.unsplash.com/photo-1590080875515-8a3d8d77d12a?q=80&w=1000&auto=format&fit=crop"),
            Product("6", "Drip Irrigation Kit", "₹1200", "Water saving solution", "Field setup", "IrrigateX", "https://images.unsplash.com/photo-1590682680695-43b964a3ae17?q=80&w=1000&auto=format&fit=crop")
        )
    }

    val filteredProducts = remember(searchQuery) {
        if (searchQuery.isEmpty()) allProducts
        else allProducts.filter { it.name.contains(searchQuery, ignoreCase = true) }
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
                title = "Agri-Marketplace",
                onBackClick = onBackClick
            )
            
            AgriSearchBar(
                value = searchQuery,
                onValueChange = { searchQuery = it },
                modifier = Modifier.padding(bottom = 16.dp)
            )
            
            LazyVerticalGrid(
                columns = GridCells.Fixed(2),
                modifier = Modifier.fillMaxSize(),
                contentPadding = PaddingValues(24.dp),
                horizontalArrangement = Arrangement.spacedBy(16.dp),
                verticalArrangement = Arrangement.spacedBy(16.dp)
            ) {
                items(filteredProducts) { product ->
                    AgriProductCard(
                        name = product.name,
                        price = product.price,
                        company = product.company,
                        imageUrl = product.imageUrl,
                        onClick = { onProductClick(product.id) }
                    )
                }
                
                item { Spacer(modifier = Modifier.height(80.dp)) }
            }
        }
    }
}

package com.aarcsx.krishaksh.core.designsystem.components

import androidx.compose.foundation.BorderStroke
import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyRow
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.ArrowBack
import androidx.compose.material.icons.filled.*
import androidx.compose.material.icons.outlined.CloudQueue
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.draw.shadow
import androidx.compose.ui.graphics.Brush
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.layout.ContentScale
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import coil.compose.AsyncImage

// Design Tokens
val ForestGreen = Color(0xFF1B5E20)
val ForestGreenLight = Color(0xFF2E7D32)
val WarmWhite = Color(0xFFF5F7F2)
val SoftGray = Color(0xFFE0E0E0)
val WarningRed = Color(0xFFC62828)
val InfoBlue = Color(0xFF1565C0)

val AgriCornerRadius = 28.dp
val AgriPadding = 24.dp

@Composable
fun AgriSubHeader(
    title: String,
    onBackClick: () -> Unit,
    modifier: Modifier = Modifier,
    trailingAction: (@Composable () -> Unit)? = null
) {
    Row(
        modifier = modifier
            .fillMaxWidth()
            .padding(horizontal = 24.dp, vertical = 16.dp),
        verticalAlignment = Alignment.CenterVertically
    ) {
        IconButton(
            onClick = onBackClick,
            modifier = Modifier
                .size(48.dp)
                .background(Color.White, CircleShape)
        ) {
            Icon(
                Icons.AutoMirrored.Filled.ArrowBack,
                contentDescription = "Back",
                tint = ForestGreen
            )
        }
        Spacer(modifier = Modifier.width(16.dp))
        Text(
            text = title,
            style = MaterialTheme.typography.headlineSmall,
            fontWeight = FontWeight.Bold,
            color = ForestGreen,
            modifier = Modifier.weight(1f)
        )
        if (trailingAction != null) {
            trailingAction()
        }
    }
}

@Composable
fun AgriHeader(
    userName: String,
    location: String,
    onNotificationClick: () -> Unit
) {
    Row(
        modifier = Modifier
            .fillMaxWidth()
            .padding(horizontal = AgriPadding, vertical = 20.dp),
        horizontalArrangement = Arrangement.SpaceBetween,
        verticalAlignment = Alignment.CenterVertically
    ) {
        Column {
            Text(
                text = "Namaste, $userName",
                style = MaterialTheme.typography.headlineSmall,
                fontWeight = FontWeight.Bold,
                color = ForestGreen
            )
            Row(verticalAlignment = Alignment.CenterVertically) {
                Icon(
                    imageVector = Icons.Default.LocationOn,
                    contentDescription = null,
                    tint = ForestGreenLight,
                    modifier = Modifier.size(16.dp)
                )
                Spacer(Modifier.width(4.dp))
                Text(
                    text = location,
                    style = MaterialTheme.typography.bodyMedium,
                    color = Color.Gray
                )
            }
        }
        IconButton(
            onClick = onNotificationClick,
            modifier = Modifier
                .background(Color.White, CircleShape)
                .shadow(2.dp, CircleShape)
        ) {
            Icon(Icons.Default.Notifications, contentDescription = "Notifications", tint = ForestGreen)
        }
    }
}

@Composable
fun AgriHeroSection(
    onScanClick: () -> Unit
) {
    Card(
        modifier = Modifier
            .fillMaxWidth()
            .padding(horizontal = AgriPadding)
            .shadow(4.dp, RoundedCornerShape(AgriCornerRadius)),
        shape = RoundedCornerShape(AgriCornerRadius),
        colors = CardDefaults.cardColors(containerColor = Color.White)
    ) {
        Column(
            modifier = Modifier
                .fillMaxWidth()
                .background(
                    Brush.verticalGradient(
                        colors = listOf(Color.White, Color(0xFFE8F5E9))
                    )
                )
                .padding(AgriPadding),
            horizontalAlignment = Alignment.CenterHorizontally
        ) {
            Box(
                modifier = Modifier
                    .size(80.dp)
                    .background(Color(0xFFF1F8E9), CircleShape),
                contentAlignment = Alignment.Center
            ) {
                Icon(
                    imageVector = Icons.Default.PhotoCamera,
                    contentDescription = null,
                    modifier = Modifier.size(40.dp),
                    tint = ForestGreen
                )
            }
            Spacer(Modifier.height(16.dp))
            Text(
                text = "Scan Your Crop",
                style = MaterialTheme.typography.headlineSmall,
                fontWeight = FontWeight.Bold,
                color = ForestGreen
            )
            Spacer(Modifier.height(8.dp))
            Text(
                text = "AI-powered disease detection for healthier crops",
                style = MaterialTheme.typography.bodyMedium,
                color = Color.Gray,
                textAlign = androidx.compose.ui.text.style.TextAlign.Center,
                modifier = Modifier.padding(horizontal = 16.dp)
            )
            Spacer(Modifier.height(24.dp))
            Button(
                onClick = onScanClick,
                modifier = Modifier
                    .fillMaxWidth()
                    .height(56.dp),
                shape = RoundedCornerShape(16.dp),
                colors = ButtonDefaults.buttonColors(containerColor = ForestGreen)
            ) {
                Text("Start Scanning", fontWeight = FontWeight.Bold, fontSize = 16.sp)
            }
        }
    }
}


enum class WeatherType {
    SUNNY, RAINY, CLOUDY, STORMY
}

@Composable
fun AgriWeatherAlertRow(
    temp: String,
    condition: String,
    windSpeed: String,
    weatherType: WeatherType,
    alertTitle: String,
    alertDesc: String
) {
    val (bgColor, icon, iconColor) = when (weatherType) {
        WeatherType.SUNNY -> Triple(Color(0xFFFFF9C4), Icons.Default.WbSunny, Color(0xFFFBC02D))
        WeatherType.RAINY -> Triple(Color(0xFFE1F5FE), Icons.Default.CloudQueue, Color(0xFF0288D1)) // Using Cloud for rain placeholder
        WeatherType.CLOUDY -> Triple(Color(0xFFF5F5F5), Icons.Default.Cloud, Color(0xFF757575))
        WeatherType.STORMY -> Triple(Color(0xFFEDE7F6), Icons.Default.Thunderstorm, Color(0xFF512DA8))
    }

    Row(
        modifier = Modifier
            .fillMaxWidth()
            .padding(horizontal = AgriPadding, vertical = 16.dp),
        horizontalArrangement = Arrangement.spacedBy(16.dp)
    ) {
        // Weather Card
        Card(
            modifier = Modifier.weight(1f),
            shape = RoundedCornerShape(AgriCornerRadius),
            colors = CardDefaults.cardColors(containerColor = bgColor)
        ) {
            Column(Modifier.padding(16.dp)) {
                Icon(icon, contentDescription = null, tint = iconColor)
                Spacer(Modifier.height(8.dp))
                Text(text = temp, fontWeight = FontWeight.Bold, fontSize = 20.sp, color = ForestGreen)
                Text(text = condition, style = MaterialTheme.typography.bodySmall, color = ForestGreen.copy(0.6f))
                Spacer(Modifier.height(4.dp))
                Row(verticalAlignment = Alignment.CenterVertically) {
                    Icon(
                        imageVector = Icons.Default.Air,
                        contentDescription = null,
                        modifier = Modifier.size(12.dp),
                        tint = ForestGreen.copy(0.6f)
                    )
                    Spacer(Modifier.width(4.dp))
                    Text(text = windSpeed, fontSize = 10.sp, color = ForestGreen.copy(0.6f))
                }
            }
        }
        // Alert Card
        Card(
            modifier = Modifier.weight(1f),
            shape = RoundedCornerShape(AgriCornerRadius),
            colors = CardDefaults.cardColors(containerColor = Color(0xFFFFF3E0))
        ) {
            Column(Modifier.padding(16.dp)) {
                Icon(Icons.Default.Warning, contentDescription = null, tint = WarningRed)
                Spacer(Modifier.height(8.dp))
                Text(text = alertTitle, fontWeight = FontWeight.Bold, fontSize = 14.sp, color = WarningRed)
                Text(text = alertDesc, style = MaterialTheme.typography.bodySmall, color = Color.DarkGray)
            }
        }
    }
}

@Composable
fun AgriSearchBar(
    value: String,
    onValueChange: (String) -> Unit,
    placeholder: String = "Search products...",
    modifier: Modifier = Modifier
) {
    TextField(
        value = value,
        onValueChange = onValueChange,
        modifier = modifier
            .fillMaxWidth()
            .padding(horizontal = AgriPadding)
            .shadow(2.dp, RoundedCornerShape(16.dp)),
        placeholder = { Text(placeholder, color = Color.Gray) },
        leadingIcon = { Icon(Icons.Default.Search, contentDescription = null, tint = ForestGreen) },
        colors = TextFieldDefaults.colors(
            focusedContainerColor = Color.White,
            unfocusedContainerColor = Color.White,
            focusedIndicatorColor = Color.Transparent,
            unfocusedIndicatorColor = Color.Transparent,
            cursorColor = ForestGreen
        ),
        shape = RoundedCornerShape(16.dp),
        singleLine = true
    )
}

@Composable
fun AgriProductCard(
    name: String,
    price: String,
    company: String,
    imageUrl: String,
    onClick: () -> Unit
) {
    Card(
        modifier = Modifier
            .fillMaxWidth()
            .clickable { onClick() },
        shape = RoundedCornerShape(20.dp),
        colors = CardDefaults.cardColors(containerColor = Color.White),
        elevation = CardDefaults.cardElevation(defaultElevation = 2.dp)
    ) {
        Column {
            AsyncImage(
                model = imageUrl,
                contentDescription = null,
                modifier = Modifier
                    .fillMaxWidth()
                    .height(140.dp)
                    .clip(RoundedCornerShape(topStart = 20.dp, topEnd = 20.dp)),
                contentScale = ContentScale.Crop
            )
            Column(Modifier.padding(12.dp)) {
                Text(
                    text = name,
                    fontWeight = FontWeight.Bold,
                    color = ForestGreen,
                    maxLines = 1,
                    overflow = androidx.compose.ui.text.style.TextOverflow.Ellipsis
                )
                Text(text = company, fontSize = 12.sp, color = Color.Gray)
                Spacer(Modifier.height(8.dp))
                Text(
                    text = price,
                    fontWeight = FontWeight.ExtraBold,
                    color = ForestGreen,
                    fontSize = 16.sp
                )
            }
        }
    }
}

@Composable
fun AgriRecentScanCard(
    cropName: String,
    status: String,
    time: String,
    imageUrl: String,
    isHealthy: Boolean,
    onClick: () -> Unit
) {
    Card(
        modifier = Modifier
            .fillMaxWidth()
            .clickable { onClick() },
        shape = RoundedCornerShape(20.dp),
        colors = CardDefaults.cardColors(containerColor = Color.White)
    ) {
        Row(verticalAlignment = Alignment.CenterVertically) {
            AsyncImage(
                model = imageUrl,
                contentDescription = null,
                modifier = Modifier
                    .size(80.dp)
                    .clip(RoundedCornerShape(20.dp)),
                contentScale = ContentScale.Crop
            )
            Column(Modifier.padding(16.dp)) {
                Text(text = cropName, fontWeight = FontWeight.Bold, color = ForestGreen)
                Text(
                    text = status, 
                    fontSize = 12.sp, 
                    color = if (isHealthy) ForestGreenLight else WarningRed
                )
                Text(text = time, fontSize = 10.sp, color = Color.Gray)
            }
        }
    }
}

@Composable
fun AgriRecentScansRow(
    scans: List<RecentScanData>,
    onScanClick: (String) -> Unit
) {
    Column {
        PaddingValues(horizontal = AgriPadding).let {
            Text(
                text = "Recent Scans",
                style = MaterialTheme.typography.titleLarge,
                fontWeight = FontWeight.Bold,
                color = ForestGreen,
                modifier = Modifier.padding(horizontal = AgriPadding)
            )
        }
        Spacer(Modifier.height(12.dp))
        LazyRow(
            contentPadding = PaddingValues(horizontal = AgriPadding),
            horizontalArrangement = Arrangement.spacedBy(16.dp)
        ) {
            items(scans) { scan ->
                Card(
                    modifier = Modifier
                        .width(160.dp)
                        .clickable { onScanClick(scan.id) },
                    shape = RoundedCornerShape(20.dp),
                    colors = CardDefaults.cardColors(containerColor = Color.White)
                ) {
                    Column {
                        AsyncImage(
                            model = scan.imageUrl,
                            contentDescription = null,
                            modifier = Modifier
                                .fillMaxWidth()
                                .height(100.dp)
                                .clip(RoundedCornerShape(topStart = 20.dp, topEnd = 20.dp)),
                            contentScale = ContentScale.Crop
                        )
                        Column(Modifier.padding(12.dp)) {
                            Text(text = scan.cropName, fontWeight = FontWeight.Bold, color = ForestGreen)
                            Text(text = scan.status, fontSize = 12.sp, color = if (scan.isHealthy) ForestGreenLight else WarningRed)
                            Text(text = scan.time, fontSize = 10.sp, color = Color.Gray)
                        }
                    }
                }
            }
        }
    }
}

@Composable
fun AgriRecommendationCard(
    title: String,
    description: String,
    onClick: () -> Unit
) {
    Card(
        modifier = Modifier
            .fillMaxWidth()
            .padding(AgriPadding)
            .clickable { onClick() },
        shape = RoundedCornerShape(AgriCornerRadius),
        colors = CardDefaults.cardColors(containerColor = ForestGreen)
    ) {
        Row(
            modifier = Modifier.padding(20.dp),
            verticalAlignment = Alignment.CenterVertically
        ) {
            Column(modifier = Modifier.weight(1f)) {
                Text(text = title, fontWeight = FontWeight.Bold, color = Color.White, fontSize = 16.sp)
                Spacer(Modifier.height(4.dp))
                Text(text = description, style = MaterialTheme.typography.bodySmall, color = Color(0xFFE8F5E9))
            }
            Icon(Icons.Default.ChevronRight, contentDescription = null, tint = Color.White)
        }
    }
}

data class RecentScanData(
    val id: String,
    val cropName: String,
    val status: String,
    val time: String,
    val imageUrl: String,
    val isHealthy: Boolean
)

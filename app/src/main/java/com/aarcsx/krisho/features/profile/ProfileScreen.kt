package com.aarcsx.krisho.features.profile

import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.*
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.ArrowForwardIos
import androidx.compose.material.icons.filled.AccountCircle
import androidx.compose.material.icons.filled.Check
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.material.icons.filled.Logout
import androidx.hilt.navigation.compose.hiltViewModel
import com.aarcsx.krisho.core.designsystem.components.*
import com.aarcsx.krisho.core.designsystem.theme.ForestGreen

@Composable
fun ProfileScreen(
    onBackClick: () -> Unit = {},
    onNavigateToHelp: () -> Unit = {},
    onNavigateToPrivacy: () -> Unit = {},
    onNavigateToTerms: () -> Unit = {},
    onLogoutSuccess: () -> Unit = {},
    viewModel: ProfileViewModel = hiltViewModel()
) {
    val uiState by viewModel.uiState.collectAsState()

    Scaffold(
        containerColor = Color(0xFFF5F7F2)
    ) { paddingValues ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(paddingValues)
        ) {
            AgriSubHeader(
                title = "My Profile",
                onBackClick = onBackClick,
                trailingAction = {
                    TextButton(onClick = { viewModel.logout(onLogoutSuccess) }) {
                        Text(
                            "Logout",
                            color = MaterialTheme.colorScheme.error,
                            fontWeight = FontWeight.Bold
                        )
                    }
                }
            )

            Column(
                modifier = Modifier
                    .fillMaxSize()
                    .padding(horizontal = 24.dp),
                verticalArrangement = Arrangement.spacedBy(16.dp)
            ) {
                // Profile Header Card
                Card(
                    modifier = Modifier.fillMaxWidth(),
                    shape = RoundedCornerShape(28.dp),
                    colors = CardDefaults.cardColors(containerColor = Color.White)
                ) {
                    Row(
                        modifier = Modifier
                            .padding(24.dp)
                            .fillMaxWidth(),
                        verticalAlignment = Alignment.CenterVertically
                    ) {
                        Box(
                            modifier = Modifier
                                .size(64.dp)
                                .background(ForestGreen.copy(alpha = 0.1f), CircleShape),
                            contentAlignment = Alignment.Center
                        ) {
                            Icon(
                                Icons.Default.AccountCircle,
                                contentDescription = null,
                                modifier = Modifier.size(40.dp),
                                tint = ForestGreen
                            )
                        }
                        Spacer(modifier = Modifier.width(16.dp))
                        Column {
                            Text(
                                uiState.profile?.name ?: "Loading...",
                                style = MaterialTheme.typography.titleLarge,
                                fontWeight = FontWeight.Bold,
                                color = ForestGreen
                            )
                            Text(
                                "Farmer • ${uiState.location}",
                                style = MaterialTheme.typography.bodyMedium,
                                color = Color.Gray
                            )
                        }
                    }
                }

                Spacer(modifier = Modifier.height(8.dp))

                // Settings List
                ProfileSettingItem(
                    title = "Language",
                    value = uiState.selectedLanguage,
                    onClick = { viewModel.toggleLanguageDialog(true) }
                )
                ProfileSettingItem(
                    title = "Crops Grown",
                    value = uiState.profile?.crops?.joinToString(", ") ?: "Select Crops",
                    onClick = { viewModel.toggleCropsDialog(true) }
                )
                ProfileSettingItem("Help Center", onClick = onNavigateToHelp)
                ProfileSettingItem("Privacy Policy", onClick = onNavigateToPrivacy)
                ProfileSettingItem("Terms of Service", onClick = onNavigateToTerms)

                Spacer(modifier = Modifier.height(16.dp))

                Column(
                    modifier = Modifier.fillMaxWidth(),
                    horizontalAlignment = Alignment.CenterHorizontally
                ) {
                    Text(
                        text = "App Version",
                        style = MaterialTheme.typography.labelMedium,
                        color = Color.Gray
                    )
                    Text(
                        text = "1.0.2",
                        style = MaterialTheme.typography.titleMedium,
                        fontWeight = FontWeight.Bold,
                        color = ForestGreen
                    )
                }

                Spacer(modifier = Modifier.height(16.dp))

                Button(
                    onClick = { viewModel.logout(onLogoutSuccess) },
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(horizontal = 24.dp)
                        .height(56.dp),
                    shape = RoundedCornerShape(16.dp),
                    colors = ButtonDefaults.buttonColors(containerColor = MaterialTheme.colorScheme.error)
                ) {
                    Icon(Icons.Default.Logout, contentDescription = null)
                    Spacer(modifier = Modifier.width(8.dp))
                    Text("Logout", fontWeight = FontWeight.Bold)
                }
            }
        }
    }

    if (uiState.showLanguageDialog) {
        LanguageSelectionDialog(
            currentLang = if (uiState.selectedLanguage == "Hindi") "hi" else "en",
            onDismiss = { viewModel.toggleLanguageDialog(false) },
            onLanguageSelected = { viewModel.updateLanguage(it) }
        )
    }

    if (uiState.showCropsDialog) {
        CropsSelectionDialog(
            selectedCrops = uiState.profile?.crops ?: emptyList(),
            onDismiss = { viewModel.toggleCropsDialog(false) },
            onCropsSelected = { viewModel.updateCrops(it) }
        )
    }
}

@Composable
fun LanguageSelectionDialog(
    currentLang: String,
    onDismiss: () -> Unit,
    onLanguageSelected: (String) -> Unit
) {
    AlertDialog(
        onDismissRequest = onDismiss,
        title = { Text("Select Language") },
        text = {
            Column {
                LanguageItem("English", "en", currentLang == "en") { onLanguageSelected("en") }
                LanguageItem("Hindi", "hi", currentLang == "hi") { onLanguageSelected("hi") }
            }
        },
        confirmButton = {
            TextButton(onClick = onDismiss) { Text("Cancel", color = ForestGreen) }
        }
    )
}

@Composable
fun LanguageItem(label: String, code: String, isSelected: Boolean, onClick: () -> Unit) {
    Row(
        modifier = Modifier
            .fillMaxWidth()
            .clickable { onClick() }
            .padding(vertical = 12.dp),
        verticalAlignment = Alignment.CenterVertically
    ) {
        RadioButton(selected = isSelected, onClick = onClick, colors = RadioButtonDefaults.colors(selectedColor = ForestGreen))
        Spacer(modifier = Modifier.width(12.dp))
        Text(label, style = MaterialTheme.typography.bodyLarge)
    }
}

@Composable
fun CropsSelectionDialog(
    selectedCrops: List<String>,
    onDismiss: () -> Unit,
    onCropsSelected: (List<String>) -> Unit
) {
    val allCrops = listOf("Wheat", "Rice", "Tomato", "Potato", "Cotton", "Mustard", "Sugarcane", "Maize")
    val currentSelection = remember { mutableStateListOf<String>().apply { addAll(selectedCrops) } }

    AlertDialog(
        onDismissRequest = onDismiss,
        title = { Text("Crops Grown") },
        text = {
            LazyColumn(modifier = Modifier.heightIn(max = 400.dp)) {
                items(allCrops) { crop ->
                    val isChecked = currentSelection.contains(crop)
                    Row(
                        modifier = Modifier
                            .fillMaxWidth()
                            .clickable {
                                if (isChecked) currentSelection.remove(crop) else currentSelection.add(crop)
                            }
                            .padding(vertical = 8.dp),
                        verticalAlignment = Alignment.CenterVertically
                    ) {
                        Checkbox(
                            checked = isChecked,
                            onCheckedChange = {
                                if (it) currentSelection.add(crop) else currentSelection.remove(crop)
                            },
                            colors = CheckboxDefaults.colors(checkedColor = ForestGreen)
                        )
                        Spacer(modifier = Modifier.width(8.dp))
                        Text(crop)
                    }
                }
            }
        },
        confirmButton = {
            Button(
                onClick = { onCropsSelected(currentSelection.toList()) },
                colors = ButtonDefaults.buttonColors(containerColor = ForestGreen)
            ) {
                Text("Save")
            }
        },
        dismissButton = {
            TextButton(onClick = onDismiss) { Text("Cancel", color = ForestGreen) }
        }
    )
}

@Composable
fun ProfileSettingItem(
    title: String,
    value: String = "",
    onClick: () -> Unit = {}
) {
    Card(
        modifier = Modifier
            .fillMaxWidth()
            .clickable { onClick() },
        shape = RoundedCornerShape(16.dp),
        colors = CardDefaults.cardColors(containerColor = Color.White)
    ) {
        ListItem(
            colors = ListItemDefaults.colors(containerColor = Color.Transparent),
            headlineContent = { 
                Text(title, fontWeight = FontWeight.SemiBold, color = ForestGreen) 
            },
            supportingContent = if (value.isNotEmpty()) {
                { Text(value, color = Color.Gray) }
            } else null,
            trailingContent = {
                Icon(
                    Icons.AutoMirrored.Filled.ArrowForwardIos,
                    contentDescription = null,
                    modifier = Modifier.size(16.dp),
                    tint = Color.Gray
                )
            }
        )
    }
}

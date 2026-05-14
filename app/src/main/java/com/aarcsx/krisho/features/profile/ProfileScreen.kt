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
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.input.KeyboardType
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
                            val displayName = uiState.profile?.name ?: uiState.profile?.email ?: "Loading..."
                            Text(
                                displayName,
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
                    title = "Phone number",
                    value = uiState.profile?.phone ?: "Add phone number",
                    onClick = { viewModel.togglePhoneDialog(true) }
                )
                ProfileSettingItem("Help Center", onClick = onNavigateToHelp)
                ProfileSettingItem("Privacy Policy", onClick = onNavigateToPrivacy)
                ProfileSettingItem("Terms of Service", onClick = onNavigateToTerms)

                Spacer(modifier = Modifier.height(24.dp))

                Column(
                    modifier = Modifier.fillMaxWidth(),
                    horizontalAlignment = Alignment.CenterHorizontally
                ) {
                    Row(
                        verticalAlignment = Alignment.CenterVertically,
                        horizontalArrangement = Arrangement.spacedBy(4.dp)
                    ) {
                        Text(
                            text = "App Version",
                            style = MaterialTheme.typography.labelMedium,
                            color = Color.Gray
                        )
                        Text(
                            text = "1.0.2",
                            style = MaterialTheme.typography.labelMedium,
                            fontWeight = FontWeight.Bold,
                            color = ForestGreen
                        )
                    }
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

    if (uiState.showPhoneDialog) {
        AlertDialog(
            onDismissRequest = { viewModel.togglePhoneDialog(false) },
            title = { Text("Update Phone Number") },
            text = {
                Column {
                    Text(
                        "Enter a phone number so our team can call you back.",
                        style = MaterialTheme.typography.bodyMedium,
                        color = Color.Gray
                    )
                    Spacer(modifier = Modifier.height(16.dp))
                    OutlinedTextField(
                        value = uiState.phoneInput,
                        onValueChange = { viewModel.updatePhoneInput(it) },
                        placeholder = { Text("Enter phone number") },
                        singleLine = true,
                        keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Number),
                        modifier = Modifier.fillMaxWidth(),
                        colors = OutlinedTextFieldDefaults.colors(
                            focusedBorderColor = ForestGreen,
                            cursorColor = ForestGreen
                        )
                    )
                }
            },
            confirmButton = {
                TextButton(onClick = { viewModel.updatePhone(uiState.phoneInput) }) {
                    Text("Save", color = ForestGreen)
                }
            },
            dismissButton = {
                TextButton(onClick = { viewModel.togglePhoneDialog(false) }) {
                    Text("Cancel", color = ForestGreen)
                }
            }
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

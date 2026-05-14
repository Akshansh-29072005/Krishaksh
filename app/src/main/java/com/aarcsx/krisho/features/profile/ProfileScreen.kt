package com.aarcsx.krisho.features.profile

import android.app.Activity
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
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.input.KeyboardType
import androidx.compose.ui.unit.dp
import androidx.compose.material.icons.filled.Logout
import com.aarcsx.krisho.R
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
                title = stringResource(R.string.my_profile),
                onBackClick = onBackClick,
                trailingAction = {
                    TextButton(onClick = { viewModel.logout(onLogoutSuccess) }) {
                        Text(
                            stringResource(R.string.logout),
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
                                stringResource(R.string.farmer_location, uiState.location),
                                style = MaterialTheme.typography.bodyMedium,
                                color = Color.Gray
                            )
                        }
                    }
                }

                Spacer(modifier = Modifier.height(8.dp))

                // Settings List
                ProfileSettingItem(
                    title = stringResource(R.string.language),
                    value = uiState.selectedLanguage,
                    onClick = { viewModel.toggleLanguageDialog(true) }
                )
                ProfileSettingItem(
                    title = stringResource(R.string.phone_number),
                    value = uiState.profile?.phone ?: stringResource(R.string.add_phone_number),
                    onClick = { viewModel.togglePhoneDialog(true) }
                )
                ProfileSettingItem(stringResource(R.string.help_center), onClick = onNavigateToHelp)
                ProfileSettingItem(stringResource(R.string.privacy_policy), onClick = onNavigateToPrivacy)
                ProfileSettingItem(stringResource(R.string.terms_of_service), onClick = onNavigateToTerms)

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
                            text = stringResource(R.string.app_version),
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

    val activity = LocalContext.current as? Activity

    if (uiState.showLanguageDialog) {
        LanguageSelectionDialog(
            currentLang = if (uiState.selectedLanguage == "Hindi") "hi" else "en",
            onDismiss = { viewModel.toggleLanguageDialog(false) },
            onLanguageSelected = {
                viewModel.updateLanguage(it) {
                    activity?.recreate()
                }
            }
        )
    }

    if (uiState.showPhoneDialog) {
        AlertDialog(
            onDismissRequest = { viewModel.togglePhoneDialog(false) },
            title = { Text(stringResource(R.string.update_phone_number)) },
            text = {
                Column {
                    Text(
                        stringResource(R.string.enter_phone_prompt),
                        style = MaterialTheme.typography.bodyMedium,
                        color = Color.Gray
                    )
                    Spacer(modifier = Modifier.height(16.dp))
                    OutlinedTextField(
                        value = uiState.phoneInput,
                        onValueChange = { viewModel.updatePhoneInput(it) },
                        placeholder = { Text(stringResource(R.string.enter_phone_number)) },
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
                    Text(stringResource(R.string.save), color = ForestGreen)
                }
            },
            dismissButton = {
                TextButton(onClick = { viewModel.togglePhoneDialog(false) }) {
                    Text(stringResource(R.string.cancel), color = ForestGreen)
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
        title = { Text(stringResource(R.string.select_language)) },
        text = {
            Column {
                LanguageItem(stringResource(R.string.english), "en", currentLang == "en") { onLanguageSelected("en") }
                LanguageItem(stringResource(R.string.hindi), "hi", currentLang == "hi") { onLanguageSelected("hi") }
            }
        },
        confirmButton = {
            TextButton(onClick = onDismiss) { Text(stringResource(R.string.cancel), color = ForestGreen) }
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

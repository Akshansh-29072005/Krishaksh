package com.aarcsx.krishaksh.features.legal

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.verticalScroll
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.ArrowBack
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.unit.dp
import com.aarcsx.krishaksh.R

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun LegalScreen(
    title: String,
    content: String,
    onBackClick: () -> Unit
) {
    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text(title) },
                navigationIcon = {
                    IconButton(onClick = onBackClick) {
                        Icon(Icons.Default.ArrowBack, contentDescription = "Back")
                    }
                }
            )
        }
    ) { padding ->
        Column(
            modifier = Modifier
                .padding(padding)
                .fillMaxSize()
                .verticalScroll(rememberScrollState())
                .padding(16.dp)
        ) {
            Text(
                text = content,
                style = MaterialTheme.typography.bodyLarge
            )
        }
    }
}

@Composable
fun HelpCenterScreen(onBackClick: () -> Unit) {
    LegalScreen(
        title = "Help Center",
        content = "Welcome to Krishaksh Help Center.\n\n" +
                "1. How to scan a crop?\n" +
                "Go to the Scan tab, select your crop, and click a clear photo of the affected area.\n\n" +
                "2. What crops are supported?\n" +
                "We currently support Wheat, Mustard, and Cotton. More crops are coming soon.\n\n" +
                "3. How to contact support?\n" +
                "You can use the Support tab to raise a ticket or send a voice message.",
        onBackClick = onBackClick
    )
}

@Composable
fun PrivacyPolicyScreen(onBackClick: () -> Unit) {
    LegalScreen(
        title = "Privacy Policy",
        content = "Your privacy is important to us. Krishaksh collects minimum data required for crop diagnosis and history tracking.\n\n" +
                "1. Data Collection: We collect crop images, approximate location for disease alerts, and profile info.\n\n" +
                "2. Data Usage: Images are used to improve our AI accuracy.\n\n" +
                "3. Third Parties: We do not sell your personal data to third parties.",
        onBackClick = onBackClick
    )
}

@Composable
fun TermsOfServiceScreen(onBackClick: () -> Unit) {
    LegalScreen(
        title = "Terms of Service",
        content = "By using Krishaksh, you agree to these terms:\n\n" +
                "1. Use for Information: Our AI diagnosis is for information purposes. Consult an agriculture expert for critical decisions.\n\n" +
                "2. Fair Use: Do not misuse our API or upload harmful content.\n\n" +
                "3. Updates: We may update these terms occasionally.",
        onBackClick = onBackClick
    )
}

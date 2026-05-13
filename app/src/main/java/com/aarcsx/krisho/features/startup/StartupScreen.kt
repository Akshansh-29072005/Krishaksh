package com.aarcsx.krisho.features.startup

import android.content.Intent
import android.net.Uri
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Refresh
import androidx.compose.material3.Button
import androidx.compose.material3.ButtonDefaults
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.material3.Icon
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Surface
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.hilt.navigation.compose.hiltViewModel
import com.aarcsx.krisho.BuildConfig
import com.aarcsx.krisho.core.designsystem.theme.ForestGreen
import com.aarcsx.krisho.navigation.KrishoNavGraph

@Composable
fun AppStartupScreen(viewModel: AppStartupViewModel = hiltViewModel()) {
    val uiState by viewModel.uiState.collectAsState()

    val errorMessage = uiState.errorMessage
    KrishoThemeSurface {
        when {
            uiState.isLoading -> StartupLoadingScreen()
            uiState.isForceUpdateRequired -> UpdateRequiredScreen(uiState.appConfig!!)
            errorMessage != null -> StartupErrorScreen(errorMessage, viewModel::refreshStartupState)
            else -> KrishoNavGraph(startDestination = uiState.startDestination)
        }
    }
}

@Composable
private fun KrishoThemeSurface(content: @Composable () -> Unit) {
    Surface(
        modifier = Modifier.fillMaxSize(),
        color = MaterialTheme.colorScheme.background
    ) {
        content()
    }
}

@Composable
private fun StartupLoadingScreen() {
    Box(
        modifier = Modifier.fillMaxSize(),
        contentAlignment = Alignment.Center
    ) {
        CircularProgressIndicator()
    }
}

@Composable
private fun StartupErrorScreen(errorMessage: String, onRetry: () -> Unit) {
    val context = LocalContext.current
    Column(
        modifier = Modifier
            .fillMaxSize()
            .padding(24.dp),
        verticalArrangement = Arrangement.Center,
        horizontalAlignment = Alignment.CenterHorizontally
    ) {
        Text(
            text = "Oops, we couldn’t verify app version",
            style = MaterialTheme.typography.titleLarge,
            fontWeight = FontWeight.Bold
        )
        Text(
            text = errorMessage,
            style = MaterialTheme.typography.bodyMedium,
            modifier = Modifier.padding(vertical = 12.dp)
        )
        Button(
            onClick = onRetry,
            modifier = Modifier.padding(top = 8.dp),
            colors = ButtonDefaults.buttonColors(containerColor = ForestGreen)
        ) {
            Icon(imageVector = Icons.Default.Refresh, contentDescription = "Retry", modifier = Modifier.size(18.dp))
            Text(text = "Retry", modifier = Modifier.padding(start = 8.dp), fontWeight = FontWeight.Bold)
        }
        Button(
            onClick = {
                val intent = Intent(Intent.ACTION_VIEW, Uri.parse("https://play.google.com/store/apps/details?id=com.aarcsx.krisho"))
                context.startActivity(intent)
            },
            modifier = Modifier.padding(top = 16.dp),
            colors = ButtonDefaults.buttonColors(containerColor = ForestGreen)
        ) {
            Text("Get Latest App", fontWeight = FontWeight.Bold)
        }
    }
}

@Composable
fun UpdateRequiredScreen(config: com.aarcsx.krisho.core.network.dto.AppConfigDto) {
    val context = LocalContext.current
    Column(
        modifier = Modifier
            .fillMaxSize()
            .padding(24.dp),
        verticalArrangement = Arrangement.Center,
        horizontalAlignment = Alignment.CenterHorizontally
    ) {
        Text(
            text = "Update Required",
            style = MaterialTheme.typography.titleLarge,
            fontWeight = FontWeight.Bold
        )
        Text(
            text = config.message ?: "A newer version of Krisho is required for security patches and bug fixes.",
            style = MaterialTheme.typography.bodyMedium,
            modifier = Modifier.padding(vertical = 12.dp)
        )
        Text(
            text = "Current version: ${BuildConfig.VERSION_NAME} (${BuildConfig.VERSION_CODE})",
            style = MaterialTheme.typography.bodyMedium
        )
        Text(
            text = "Required version: ${config.latest_version_name ?: config.minimum_version_code.toString()}",
            style = MaterialTheme.typography.bodyMedium,
            modifier = Modifier.padding(vertical = 8.dp)
        )
        Button(
            onClick = {
                config.update_url?.takeIf { it.isNotBlank() }?.let { url ->
                    val intent = Intent(Intent.ACTION_VIEW, Uri.parse(url))
                    context.startActivity(intent)
                }
            },
            enabled = !config.update_url.isNullOrBlank(),
            colors = ButtonDefaults.buttonColors(containerColor = ForestGreen)
        ) {
            Text("Update App", fontWeight = FontWeight.Bold)
        }
    }
}

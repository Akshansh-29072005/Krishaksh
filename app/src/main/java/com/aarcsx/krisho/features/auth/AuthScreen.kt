package com.aarcsx.krisho.features.auth

import androidx.compose.foundation.Image
import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.hilt.navigation.compose.hiltViewModel
import androidx.activity.compose.rememberLauncherForActivityResult
import androidx.activity.result.contract.ActivityResultContracts
import com.aarcsx.krisho.R
import com.aarcsx.krisho.core.designsystem.components.ForestGreen
import com.aarcsx.krisho.core.auth.GoogleSignInManager
import com.google.android.gms.auth.api.signin.GoogleSignIn
import androidx.compose.ui.platform.LocalContext

@Composable
fun AuthScreen(
    onAuthSuccess: () -> Unit,
    viewModel: AuthViewModel = hiltViewModel()
) {
    val uiState by viewModel.uiState.collectAsState()
    val context = LocalContext.current
    
    val launcher = rememberLauncherForActivityResult(
        contract = ActivityResultContracts.StartActivityForResult()
    ) { result ->
        android.util.Log.d("AuthScreen", "Result received: ${result.resultCode}")
        
        // If result is 0 (Canceled), the user backed out or play services is broken
        if (result.resultCode == 0) {
            android.util.Log.e("AuthScreen", "Sign-in was canceled or failed at the activity level")
        }

        val task = GoogleSignIn.getSignedInAccountFromIntent(result.data)
        try {
            val account = task.getResult(com.google.android.gms.common.api.ApiException::class.java)
            android.util.Log.d("AuthScreen", "Token received successfully")
            viewModel.onGoogleSignInResult(account?.idToken)
        } catch (e: com.google.android.gms.common.api.ApiException) {
            // Log the error code and check context
            android.util.Log.e("AuthScreen", "Google Sign In failed: ${e.statusCode}")
            
            if (e.statusCode == 10) {
                android.util.Log.e("AuthScreen", "Error 10: Potential misconfiguration in Google Cloud Console.")
                android.util.Log.e("AuthScreen", "Check: Package: com.aarcsx.krisho, SHA1: 7D:D6:86:90:A8:8F:BE:63:EE:5C:72:26:7F:8C:A6:03:D0:31:AB:27")
            }
            
            viewModel.onGoogleSignInResult(null)
        } catch (e: Exception) {
            android.util.Log.e("AuthScreen", "General error: ${e.message}", e)
            viewModel.onGoogleSignInResult(null)
        }
    }

    LaunchedEffect(uiState.isSuccess) {
        if (uiState.isSuccess) {
            onAuthSuccess()
        }
    }

    Column(
        modifier = Modifier
            .fillMaxSize()
            .background(Color.White)
            .padding(24.dp),
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.Center
    ) {
        // App Logo
        Image(
            painter = painterResource(id = R.drawable.ic_launcher_foreground), // Placeholder
            contentDescription = "App Logo",
            modifier = Modifier.size(120.dp)
        )
        
        Text(
            text = "Krisho",
            style = MaterialTheme.typography.headlineLarge,
            fontWeight = FontWeight.ExtraBold,
            color = ForestGreen,
            letterSpacing = 2.sp
        )
        
        Text(
            text = "Empowering Farmers with AI",
            style = MaterialTheme.typography.bodyMedium,
            color = Color.Gray,
            modifier = Modifier.padding(top = 8.dp)
        )

        Spacer(modifier = Modifier.height(64.dp))

        if (uiState.isLoading) {
            CircularProgressIndicator(color = ForestGreen)
        } else {
            // Google Sign In Button
            Button(
                onClick = { 
                    launcher.launch(viewModel.googleSignInManager.client.signInIntent)
                },
                modifier = Modifier
                    .fillMaxWidth()
                    .height(56.dp),
                shape = RoundedCornerShape(16.dp),
                colors = ButtonDefaults.buttonColors(containerColor = Color(0xFFF2F2F2), contentColor = Color.Black),
                elevation = ButtonDefaults.buttonElevation(defaultElevation = 0.dp)
            ) {
                Row(verticalAlignment = Alignment.CenterVertically) {
                    // Google Icon would go here
                    Icon(
                        painter = painterResource(id = R.drawable.ic_launcher_foreground), // Placeholder
                        contentDescription = null,
                        modifier = Modifier.size(24.dp),
                        tint = Color.Unspecified
                    )
                    Spacer(modifier = Modifier.width(12.dp))
                    Text(
                        "Sign in with Google",
                        style = MaterialTheme.typography.titleMedium,
                        fontWeight = FontWeight.Bold
                    )
                }
            }
        }

        if (uiState.error != null) {
            Spacer(modifier = Modifier.height(16.dp))
            Text(
                text = uiState.error!!,
                color = MaterialTheme.colorScheme.error,
                style = MaterialTheme.typography.bodySmall
            )
        }

        Spacer(modifier = Modifier.height(32.dp))
        
        Text(
            text = "By signing in, you agree to our Terms and Privacy Policy",
            style = MaterialTheme.typography.labelSmall,
            color = Color.LightGray,
            modifier = Modifier.padding(horizontal = 32.dp),
            textAlign = androidx.compose.ui.text.style.TextAlign.Center
        )
    }
}
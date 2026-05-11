package com.aarcsx.krisho

import android.content.Context
import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Surface
import androidx.compose.ui.Modifier
import androidx.lifecycle.lifecycleScope
import com.aarcsx.krisho.core.common.LocaleManager
import com.aarcsx.krisho.core.designsystem.theme.KrishoTheme
import com.aarcsx.krisho.core.local.datastore.PreferencesManager
import com.aarcsx.krisho.navigation.KrishoNavGraph
import com.aarcsx.krisho.navigation.Screen
import dagger.hilt.android.AndroidEntryPoint
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.launch
import javax.inject.Inject

@AndroidEntryPoint
class MainActivity : ComponentActivity() {

    @Inject
    lateinit var preferencesManager: PreferencesManager

    override fun attachBaseContext(newBase: Context) {
        // Apply saved language on startup
        super.attachBaseContext(LocaleManager.getLocaleContextWrapper(newBase, (newBase.applicationContext as KrishoApp).preferencesManager))
    }

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        
        lifecycleScope.launch {
            preferencesManager.jwtToken.collect { token ->
                val startDestination = if (token.isNullOrBlank()) Screen.Auth.route else Screen.Home.route
                
                setContent {
                    KrishoTheme {
                        Surface(
                            modifier = Modifier.fillMaxSize(),
                            color = MaterialTheme.colorScheme.background
                        ) {
                            KrishoNavGraph(startDestination = startDestination)
                        }
                    }
                }
            }
        }
    }
}

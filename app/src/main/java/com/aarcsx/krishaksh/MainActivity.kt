package com.aarcsx.krishaksh

import android.content.Context
import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Surface
import androidx.compose.ui.Modifier
import com.aarcsx.krishaksh.core.common.LocaleManager
import com.aarcsx.krishaksh.core.designsystem.theme.KrishakshTheme
import com.aarcsx.krishaksh.core.local.datastore.PreferencesManager
import com.aarcsx.krishaksh.navigation.KrishakshNavGraph
import dagger.hilt.android.AndroidEntryPoint
import javax.inject.Inject

@AndroidEntryPoint
class MainActivity : ComponentActivity() {

    @Inject
    lateinit var preferencesManager: PreferencesManager

    override fun attachBaseContext(newBase: Context) {
        // Apply saved language on startup
        super.attachBaseContext(LocaleManager.getLocaleContextWrapper(newBase, (newBase.applicationContext as KrishakshApp).preferencesManager))
    }

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContent {
            KrishakshTheme {
                Surface(
                    modifier = Modifier.fillMaxSize(),
                    color = MaterialTheme.colorScheme.background
                ) {
                    KrishakshNavGraph()
                }
            }
        }
    }
}

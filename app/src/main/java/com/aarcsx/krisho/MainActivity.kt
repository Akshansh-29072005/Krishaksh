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

import androidx.activity.enableEdgeToEdge

import com.aarcsx.krisho.core.util.PaymentManager
import com.razorpay.PaymentResultListener

@AndroidEntryPoint
class MainActivity : ComponentActivity(), PaymentResultListener {

    @Inject
    lateinit var preferencesManager: PreferencesManager

    @Inject
    lateinit var paymentManager: PaymentManager

    override fun onPaymentSuccess(razorpayPaymentId: String?) {
        lifecycleScope.launch {
            paymentManager.onPaymentSuccess(razorpayPaymentId)
        }
    }

    override fun onPaymentError(code: Int, description: String?) {
        lifecycleScope.launch {
            paymentManager.onPaymentError(code, description)
        }
    }

    override fun attachBaseContext(newBase: Context) {
        // Apply saved language on startup
        super.attachBaseContext(LocaleManager.getLocaleContextWrapper(newBase, (newBase.applicationContext as KrishoApp).preferencesManager))
    }

    override fun onCreate(savedInstanceState: Bundle?) {
        enableEdgeToEdge()
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

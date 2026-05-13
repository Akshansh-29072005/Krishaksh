package com.aarcsx.krisho

import android.content.Context
import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.lifecycle.lifecycleScope
import com.aarcsx.krisho.core.common.LocaleManager
import com.aarcsx.krisho.core.designsystem.theme.KrishoTheme
import com.aarcsx.krisho.core.local.datastore.PreferencesManager
import com.aarcsx.krisho.core.util.PaymentManager
import com.aarcsx.krisho.features.startup.AppStartupScreen
import dagger.hilt.android.AndroidEntryPoint
import kotlinx.coroutines.launch
import javax.inject.Inject

@AndroidEntryPoint
class MainActivity : ComponentActivity(), com.razorpay.PaymentResultListener {

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

        setContent {
            KrishoTheme {
                AppStartupScreen()
            }
        }
    }
}

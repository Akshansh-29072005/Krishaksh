package com.aarcsx.krisho.core.util

import android.app.Activity
import com.aarcsx.krisho.BuildConfig
import com.razorpay.Checkout
import org.json.JSONObject
import timber.log.Timber
import javax.inject.Inject
import javax.inject.Singleton
import kotlinx.coroutines.flow.MutableSharedFlow
import kotlinx.coroutines.flow.asSharedFlow

sealed class PaymentResult {
    data class Success(val paymentId: String?) : PaymentResult()
    data class Error(val code: Int, val description: String?) : PaymentResult()
}

@Singleton
class PaymentManager @Inject constructor() {
    
    private val _paymentResult = MutableSharedFlow<PaymentResult>()
    val paymentResult = _paymentResult.asSharedFlow()

    fun startPayment(
        activity: Activity,
        amount: Double,
        productName: String,
        customerEmail: String,
        customerContact: String
    ) {
        val checkout = Checkout()
        checkout.setKeyID(BuildConfig.RAZORPAY_KEY_ID)
        
        try {
            val options = JSONObject()
            options.put("name", "Krisho")
            options.put("description", productName)
            options.put("image", "https://api.krisho.aarcsx.com/static/logos/krisho_round.png")
            options.put("theme.color", "#1F4D36")
            options.put("currency", "INR")
            options.put("amount", (amount * 100).toInt()) 
            
            // Allow Razorpay to handle all methods automatically based on Dashboard
            // Removing the restrictive 'method' block helps show UPI in Test Mode
            
            val retryObj = JSONObject()
            retryObj.put("enabled", true)
            retryObj.put("max_count", 4)
            options.put("retry", retryObj)

            val prefill = JSONObject()
            prefill.put("email", customerEmail)
            prefill.put("contact", customerContact)
            // Do not force "upi" as prefill method if it's not showing, let user choose
            options.put("prefill", prefill)

            checkout.open(activity, options)
        } catch (e: Exception) {
            Timber.e(e, "Error in starting Razorpay Checkout")
        }
    }

    suspend fun onPaymentSuccess(razorpayPaymentId: String?) {
        _paymentResult.emit(PaymentResult.Success(razorpayPaymentId))
    }

    suspend fun onPaymentError(code: Int, description: String?) {
        val userFriendlyMessage = when {
            description?.contains("cancelled", ignoreCase = true) == true -> "Payment was cancelled. Please try again if you want to complete the purchase."
            code == Checkout.NETWORK_ERROR -> "Network issue detected. Please check your internet connection."
            code == Checkout.INVALID_OPTIONS -> "There was a configuration error. Please contact support."
            else -> "Payment could not be completed. Please try a different payment method."
        }
        _paymentResult.emit(PaymentResult.Error(code, userFriendlyMessage))
    }
}

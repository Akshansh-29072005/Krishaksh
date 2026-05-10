package com.aarcsx.krisho.core.payment

import android.app.Activity
import com.razorpay.Checkout
import com.razorpay.PaymentResultListener
import org.json.JSONObject
import timber.log.Timber

/**
 * Manager for handling Razorpay payments with sandbox/prod support.
 */
class PaymentManager(private val activity: Activity) : PaymentResultListener {

    private var onPaymentSuccess: ((String) -> Unit)? = null
    private var onPaymentError: ((Int, String) -> Unit)? = null

    fun startPayment(
        amountInPaise: Int,
        orderId: String,
        email: String,
        contact: String,
        onSuccess: (String) -> Unit,
        onError: (Int, String) -> Unit
    ) {
        this.onPaymentSuccess = onSuccess
        this.onPaymentError = onError

        val checkout = Checkout()
        // Key should be injected from BuildConfig or safe storage
        // checkout.setKeyID("rzp_test_xxxx") 

        try {
            val options = JSONObject().apply {
                put("name", "Krisho")
                put("description", "Agri-Store Purchase")
                put("image", "https://krisho.com/logo.png")
                put("order_id", orderId)
                put("theme.color", "#1B5E20")
                put("currency", "INR")
                put("amount", amountInPaise)
                put("prefill.email", email)
                put("prefill.contact", contact)
            }
            checkout.open(activity, options)
        } catch (e: Exception) {
            Timber.e(e, "Error starting Razorpay checkout")
            onError(-1, e.localizedMessage ?: "Init error")
        }
    }

    override fun onPaymentSuccess(razorpayPaymentId: String?) {
        Timber.d("Payment Success: $razorpayPaymentId")
        onPaymentSuccess?.invoke(razorpayPaymentId ?: "")
    }

    override fun onPaymentError(code: Int, response: String?) {
        Timber.e("Payment Error $code: $response")
        onPaymentError?.invoke(code, response ?: "Payment cancelled or failed")
    }
}
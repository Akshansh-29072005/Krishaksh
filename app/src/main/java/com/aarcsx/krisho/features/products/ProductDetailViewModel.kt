package com.aarcsx.krisho.features.products

import android.app.Activity
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.aarcsx.krisho.core.local.room.entity.ProductEntity
import com.aarcsx.krisho.core.repository.ProductRepository
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.update
import kotlinx.coroutines.launch
import javax.inject.Inject

import com.aarcsx.krisho.core.util.PaymentManager
import com.aarcsx.krisho.core.util.PaymentResult
import kotlinx.coroutines.flow.collectLatest

import com.aarcsx.krisho.core.repository.UserRepository
import com.aarcsx.krisho.core.network.dto.UserMeDto
import com.aarcsx.krisho.core.common.ApiResult

data class ProductDetailUiState(
    val product: ProductEntity? = null,
    val userProfile: UserMeDto? = null,
    val isLoading: Boolean = false,
    val error: String? = null,
    val orderInitiated: Boolean = false,
    val paymentSuccessId: String? = null,
    val paymentError: String? = null
)

@HiltViewModel
class ProductDetailViewModel @Inject constructor(
    private val repository: ProductRepository,
    private val userRepository: UserRepository,
    private val paymentManager: PaymentManager
) : ViewModel() {

    private val _uiState = MutableStateFlow(ProductDetailUiState())
    val uiState: StateFlow<ProductDetailUiState> = _uiState.asStateFlow()

    init {
        viewModelScope.launch {
            paymentManager.paymentResult.collectLatest { result ->
                when (result) {
                    is PaymentResult.Success -> {
                        _uiState.update { it.copy(
                            orderInitiated = true, 
                            paymentSuccessId = result.paymentId,
                            isLoading = false
                        ) }
                    }
                    is PaymentResult.Error -> {
                        _uiState.update { it.copy(
                            paymentError = result.description ?: "Payment Failed",
                            isLoading = false
                        ) }
                    }
                }
            }
        }
    }

    fun loadProduct(productId: String) {
        viewModelScope.launch {
            _uiState.update { it.copy(isLoading = true, error = null) }
            
            // Load user profile for payment prefill
            userRepository.getProfile().collect { result ->
                if (result is ApiResult.Success) {
                    _uiState.update { it.copy(userProfile = result.data) }
                }
            }

            val details = repository.getProductDetails(productId)
            if (details != null) {
                _uiState.update { it.copy(product = details, isLoading = false) }
            } else {
                _uiState.update { it.copy(isLoading = false, error = "Product not found") }
            }
        }
    }

    fun onBuyNow(activity: Activity) {
        val product = _uiState.value.product ?: return
        val user = _uiState.value.userProfile
        
        _uiState.update { it.copy(isLoading = true, paymentError = null) }
        paymentManager.startPayment(
            activity = activity,
            amount = product.price,
            productName = product.name,
            customerEmail = user?.email ?: "farmer@krisho.com",
            customerContact = user?.phone ?: "9999999999"
        )
    }

    fun clearPaymentError() {
        _uiState.update { it.copy(paymentError = null) }
    }
}

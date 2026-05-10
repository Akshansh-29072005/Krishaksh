package com.aarcsx.krishaksh.core.data.sync

import com.aarcsx.krishaksh.core.repository.ProductRepository
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class SyncCoordinator @Inject constructor(
    private val productRepository: ProductRepository
) {
    suspend fun syncAll() {
        productRepository.syncProducts()
    }
}

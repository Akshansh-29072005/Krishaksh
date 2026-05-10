package com.aarcsx.krisho.core.repository

import com.aarcsx.krisho.core.common.ApiResult
import com.aarcsx.krisho.core.data.remote.ProductRemoteDataSource
import com.aarcsx.krisho.core.local.room.dao.ProductDao
import com.aarcsx.krisho.core.local.room.entity.ProductEntity
import kotlinx.coroutines.flow.Flow
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class ProductRepository @Inject constructor(
    private val productDao: ProductDao,
    private val remote: ProductRemoteDataSource
) {
    fun getAllProducts(): Flow<List<ProductEntity>> = productDao.getAllProducts()

    suspend fun syncProducts(): ApiResult<Unit> {
        return when (val res = remote.getProducts()) {
            is ApiResult.Success -> {
                val entities = res.data.data.orEmpty().map {
                    ProductEntity(
                        id = it.id,
                        name = it.name,
                        price = it.price ?: 0.0,
                        imageUrl = it.image_url.orEmpty(),
                        category = it.crop_type ?: "General",
                        description = it.description.orEmpty(),
                        companyName = "Krisho Partner",
                        usageInstructions = it.unit ?: "",
                        lastUpdated = System.currentTimeMillis()
                    )
                }
                productDao.clearProducts()
                productDao.insertProducts(entities)
                ApiResult.Success(Unit)
            }
            is ApiResult.Error -> res
            ApiResult.Loading -> ApiResult.Loading
        }
    }

    suspend fun getProductDetails(id: String): ProductEntity? {
        return productDao.getProductById(id) ?: run {
            when (val remoteRes = remote.getProduct(id)) {
                is ApiResult.Success -> remoteRes.data.data?.let {
                    ProductEntity(
                        id = it.id,
                        name = it.name,
                        price = it.price ?: 0.0,
                        imageUrl = it.image_url.orEmpty(),
                        category = it.crop_type ?: "General",
                        description = it.description.orEmpty(),
                        companyName = "Krisho Partner",
                        usageInstructions = it.unit ?: "",
                        lastUpdated = System.currentTimeMillis()
                    )
                }
                else -> null
            }
        }
    }
}

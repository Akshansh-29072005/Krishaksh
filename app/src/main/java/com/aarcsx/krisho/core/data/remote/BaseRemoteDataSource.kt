package com.aarcsx.krisho.core.data.remote

import com.aarcsx.krisho.core.common.ApiResult
import com.aarcsx.krisho.core.network.util.ApiErrorParser
import retrofit2.Response

abstract class BaseRemoteDataSource {
    protected suspend fun <T> call(api: suspend () -> Response<T>): ApiResult<T> {
        return try {
            val res = api()
            if (res.isSuccessful && res.body() != null) ApiResult.Success(res.body()!!)
            else ApiResult.Error(res.code(), ApiErrorParser.parse(res))
        } catch (e: Exception) {
            e.printStackTrace()
            ApiResult.Error(message = e.message ?: "Network error", throwable = e)
        }
    }
}

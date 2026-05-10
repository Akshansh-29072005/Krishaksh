package com.aarcsx.krishaksh.core.data.remote

import com.aarcsx.krishaksh.core.common.ApiResult
import com.aarcsx.krishaksh.core.network.util.ApiErrorParser
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

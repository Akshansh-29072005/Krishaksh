package com.aarcsx.krishaksh.core.common

data class PageRequest(val page: Int = 1, val pageSize: Int = 20)
data class PageResult<T>(val items: List<T>, val nextPage: Int? = null)

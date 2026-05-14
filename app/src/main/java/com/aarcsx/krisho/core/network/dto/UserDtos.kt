package com.aarcsx.krisho.core.network.dto

data class UserMeDto(
    val id: String,
    val role: String,
    val email: String? = null,
    val name: String? = null,
    val phone: String? = null,
    val crops: List<String> = emptyList()
)

data class SupportTicketDto(
    val id: String,
    val title: String,
    val description: String,
    val status: String,
    val priority: String,
    val callback_requested: Boolean,
    val callback_status: String,
    val created_at: String,
    val updated_at: String,
    val resolved: Boolean = false,
    val resolution_response: String? = null,
    val resolution_audio_url: String? = null
)

data class SupportMessageDto(
    val id: String,
    val sender_role: String,
    val body: String,
    val created_at: String
)

data class SupportThreadDto(
    val ticket: SupportTicketDto,
    val messages: List<SupportMessageDto>
)

data class CreateTicketDto(
    val title: String,
    val description: String,
    val priority: String = "medium",
    val request_callback: Boolean = false
)

data class SendMessageDto(
    val body: String
)

data class AdCampaignDto(
    val id: String,
    val title: String,
    val image_url: String,
    val target_url: String? = null
)

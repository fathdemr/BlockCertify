package main

// @title       BlockCertify API
// @version     1.0
// @description REST API for BlockCertify. Önce /exapi/auth/login ile token alın, ardından sağ üstteki Authorize butonuna "Bearer {token}" girin.

// @contact.name BlockCertify Team

// @host     api.blockcertify.app
// @BasePath /
// @schemes  https http

// @securityDefinitions.apikey BearerAuth
// @in          header
// @name        Authorization
// @description Bearer {token} formatında girin

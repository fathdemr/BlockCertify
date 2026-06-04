package main

import (
	"BlockCertify/internal/config"
	"BlockCertify/internal/logger"
	"BlockCertify/internal/routes"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	logger.Init()

	current, _ := os.Getwd()
	fmt.Println(current)

	if err := config.InitConfigFile("./internal"); err != nil {
		panic(err)
	}
	if err := config.InitDB(); err != nil {
		panic(err)
	}
	rediPingStatus := config.RedisClient.Ping(context.Background())
	if rediPingStatus.Err() != nil {
		panic(rediPingStatus.Err())
	}

	app := gin.New()

	app.ForwardedByClientIP = true
	// Cors alayına hak veriyor. sonra kaldıracağız.
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = []string{
		"Origin",           // İstek yapılan kaynağı (domain, port) belirtir.
		"Authorization",    // Kimlik doğrulama bilgileri taşır (Bearer Token, Basic Auth vb.).
		"Content-Type",     // İstek veya yanıt içeriğinin türünü belirtir (application/json, text/html vb.).
		"Bilet",            // Özel bir kimlik doğrulama veya yetkilendirme başlığı olabilir.
		"ApiKey",           // API erişimi için kullanılan anahtar.
		"ApiSecret",        // API erişimi için gizli anahtar.
		"X-Forwarded-For",  // İstek yapan istemcinin gerçek IP adresini taşır (proxy arkasındaysa).
		"X-Real-Ip",        // Genellikle istemcinin gerçek IP adresini belirlemek için kullanılır.
		"User-Agent",       // İstek yapan cihazın veya tarayıcının bilgisini taşır (örn. Chrome, Postman).
		"Referer",          // Kullanıcının hangi sayfadan geldiğini gösterir.
		"Accept-Language",  // İstemcinin tercih ettiği dil ayarlarını içerir.
		"Accept-Encoding",  // Sunucunun hangi sıkıştırma formatlarını (gzip, deflate vb.) desteklediğini gösterir.
		"Cache-Control",    // Önbellekleme politikasını belirtir.
		"Connection",       // Bağlantının nasıl yönetileceğini belirler (keep-alive, close vb.).
		"DNT",              // "Do Not Track" talebi, kullanıcının izlenmek istemediğini belirtir.
		"X-Requested-With", // İsteğin AJAX olup olmadığını belirlemek için kullanılır.
		"Sec-Fetch-Site",   // İsteğin hangi siteden yapıldığını gösterir (same-origin, cross-site vb.).
		"Sec-Fetch-Mode",   // İsteğin türünü belirtir (cors, no-cors vb.).
		"Sec-Fetch-Dest",   // Kaynağın hangi amaçla yüklendiğini gösterir (document, script vb.).
		"X-Device-Id",      // İstemcinin cihaz ID’sini belirtmek için özel bir başlık.
		"X-Device-Model",   // İstemcinin cihaz modelini belirtmek için özel bir başlık.
		"X-OS-Version",     // İşletim sistemi sürümünü belirten özel bir başlık.
		"X-Client-Version", // Mobil uygulama istemcisinin sürümünü belirten başlık.
		"X-Platform",       // İstemcinin hangi platformdan geldiğini belirtir (iOS, Android, Web).
		"X-Timezone",       // İstemcinin bulunduğu zaman dilimini belirtir.
		"X-Session-Id",     // Kullanıcının oturum bilgisini taşır.
		"X-App-Id",         // Mobil veya web uygulamasının kimliğini belirtir.
	}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"}
	app.Use(cors.New(corsConfig))

	port := 8080

	//Public API
	//TODO middleware ekle
	exapi := app.Group("/exapi")
	routes.UserRoutes(exapi)
	routes.UniversityRoutes(exapi)
	routes.PingRoutes(exapi)
	routes.FacultyRoutes(exapi)
	routes.DepartmentRoutes(exapi)

	//Private API
	//TODO middleware ekle
	api := app.Group("/api")
	routes.DiplomaRoutes(api)
	routes.WalletRoutes(api)

	//r.Static("/public", "./public")
	//Start server
	log.Printf("Server running on port %d", port)
	if err := app.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

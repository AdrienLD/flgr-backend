package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var websiteAccess string = "https://songs.flgr.fr"

func init() {
	if err := godotenv.Load(); err != nil {
		panic("Erreur lors du chargement du fichier .env")
	}
}

func main() {
	router := gin.Default()

	// Configurer CORS
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", websiteAccess)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, command, method")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	// Routes
	router.POST("/api/gettoken", getToken)
	router.GET("/api/getplayerstate", getPlayerState)
	router.GET("/api/testtoken", testToken)
	router.POST("/api/research", research)
	router.POST("/api/tracks", getTracks)
	router.GET("/api/playpause", playPause)
	router.POST("/api/newplaylist", getPlaylist)
	/*
		router.GET("/api/getlyricsId", getLyricsById)
	*/

	// Démarrer le serveur
	router.Run(":4000")
}

func getToken(c *gin.Context) {
	var requestBody struct {
		Action string `json:"action"`
		Code   string `json:"code"`
	}

	CLIENT_ID := os.Getenv("CLIENT_ID")
	CLIENT_SECRET := os.Getenv("CLIENT_SECRET")

	fmt.Println("CLIENT_ID:", CLIENT_ID)
	fmt.Println("CLIENT_SECRET:", CLIENT_SECRET)

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides"})
		return
	}
	var params = url.Values{}
	if requestBody.Action == "gettoken" {
		params.Set("grant_type", "authorization_code")
		params.Set("code", requestBody.Code)
		params.Set("redirect_uri", websiteAccess+"/Auth")
	} else {
		refreshToken, err := c.Cookie("refresh_token")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token manquant"})
			return
		}
		params.Set("grant_type", "refresh_token")
		params.Set("refresh_token", refreshToken)
	}

	req, _ := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(params.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(CLIENT_ID, CLIENT_SECRET)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la requête"})
		return
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&data)

	if accessToken, ok := data["access_token"].(string); ok {
		c.SetCookie("token", accessToken, 3600, "/", websiteAccess, false, true)
	}

	if refreshToken, ok := data["refresh_token"].(string); ok {
		c.SetCookie("refresh_token", refreshToken, 3600*24*30, "/", websiteAccess, false, true)
	}

	c.JSON(resp.StatusCode, data)
}

func getPlayerState(c *gin.Context) {
	token, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token manquant"})
		return
	}

	req, _ := http.NewRequest("GET", "https://api.spotify.com/v1/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la requête"})
		return
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&data)

	c.JSON(resp.StatusCode, data)
}

func testToken(c *gin.Context) {
	token, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token manquant"})
		return
	}

	req, _ := http.NewRequest("GET", "https://api.spotify.com/v1/me/player", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la requête"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 204 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"status": 401}})
		return
	}

	var data map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&data)

	c.JSON(resp.StatusCode, data)
}

func research(c *gin.Context) {
	var requestBody struct {
		Titre string `json:"titre"`
		Type  string `json:"type"`
	}

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides"})
		return
	}

	token, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token manquant"})
		return
	}

	url := fmt.Sprintf("https://api.spotify.com/v1/search?q=%s&type=%s", url.QueryEscape(requestBody.Titre), requestBody.Type)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la requête"})
		return
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&data)

	c.JSON(resp.StatusCode, data)
}

func getTracks(c *gin.Context) {
	var requestBody struct {
		PlaylistId string `json:"playlistId"`
		Offset     int    `json:"offset"`
	}

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides"})
		return
	}

	token, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token manquant"})
		return
	}

	url := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks?offset=%d", requestBody.PlaylistId, requestBody.Offset)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la requête"})
		return
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&data)

	c.JSON(resp.StatusCode, data)
}

func playPause(c *gin.Context) {
	commande := c.GetHeader("command")
	method := c.GetHeader("method")

	token, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token manquant"})
		return
	}

	url := fmt.Sprintf("https://api.spotify.com/v1/me/player/%s", commande)
	req, _ := http.NewRequest(method, url, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la requête"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 204 {
		var data map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&data)
		c.JSON(resp.StatusCode, data)
	} else {
		c.JSON(resp.StatusCode, gin.H{"message": "Réponse sans contenu"})
	}
}

func getPlaylist(c *gin.Context) {
	var requestBody struct {
		PlaylistId string `json:"playlistId"`
	}

	// Récupérer les données du corps de la requête
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides"})
		return
	}

	// Récupérer le token depuis les cookies
	token, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token manquant"})
		return
	}

	// Construire la requête vers l'API Spotify
	url := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s", requestBody.PlaylistId)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	// Envoyer la requête
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la requête vers Spotify"})
		return
	}
	defer resp.Body.Close()

	// Lire la réponse
	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la lecture de la réponse"})
		return
	}

	// Retourner les données au client
	c.JSON(resp.StatusCode, data)
}

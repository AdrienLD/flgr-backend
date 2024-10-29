package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/AdrienLD/flgr-backend/songs"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var websiteAccess string = "https://songs.flgr.fr"

var httpClient = &http.Client{}

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
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	})

	// Routes
	router.POST("/api/gettoken", getToken)
	router.GET("/api/getplayerstate", getPlayerState)
	router.GET("/api/testtoken", testToken)
	router.POST("/api/research", research)
	router.GET("/api/playpause", playPause)
	router.POST("/api/newplaylist", getPlaylist)
	router.POST("/api/replaceplaylist", replaceplaylist)
	router.POST("/api/nextmusic", nextMusic)

	router.Run(":4000")
}

func replaceplaylist(c *gin.Context) {
	var requestData struct {
		PlaylistId []string `json:"playlistId"`
	}

	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides"})
		return
	}

	token, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token manquant"})
		return
	}

	var results []map[string]interface{}

	for _, playlistId := range requestData.PlaylistId {
		parts := strings.Split(playlistId, " £ ")
		genre, playlistName := parts[0], parts[1]

		var id string
		var playlistInfo map[string]interface{}

		if genre == "UserPlaylist" {
			id = playlistName
			fmt.Printf("https://api.spotify.com/v1/playlists/%s\n", playlistName)
			url := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s", playlistName)
			playlistInfo = getAPICall(url, token)
		} else {
			research := songs.ReformatPlaylistResearch(genre, playlistName)
			url := fmt.Sprintf("https://api.spotify.com/v1/search?q=%s&type=%s", url.QueryEscape(research), "playlist")
			searchResult := getAPICall(url, token)
			if searchResult != nil {
				playlistInfo = songs.GetFirstPlaylistID(searchResult)
				id, _ = playlistInfo["id"].(string)
			}
		}

		tracks := recuperateAllMusics(id, token)
		for _, track := range tracks {
			if trackMap, ok := track.(map[string]interface{}); ok {
				trackMap["playlist"] = playlistInfo
			}
		}
		results = append(results, map[string]interface{}{
			"tracks": tracks,
		})

	}

	songs.EqualizePlaylists(results)
	mixedPlaylist := songs.MixAllTracks(results)

	c.JSON(http.StatusOK, gin.H{
		"message": "Playlists traitées",
		"results": mixedPlaylist,
	})
}

func getAPICall(url, token string) map[string]interface{} {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("Erreur lors de la requête:", err)
		return nil
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		fmt.Println("Erreur de décodage JSON:", err)
		return nil
	}

	return data
}

func postAPICall(url, token string) {
	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	_, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("Erreur lors de la requête:", err)
	}
}

func putAPICall(url, token string) {
	req, _ := http.NewRequest("PUT", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	_, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("Erreur lors de la requête:", err)
	}
}

func recuperateAllMusics(playlistId string, token string) []interface{} {
	var allTracks []interface{}
	offset := 0
	for {
		response := getTracks(playlistId, offset, token)
		if items, ok := response["items"].([]interface{}); ok {
			for _, item := range items {
				if trackInfo, ok := item.(map[string]interface{}); ok {
					if track, ok := trackInfo["track"]; ok {
						allTracks = append(allTracks, track)
					}
				}
			}
		}
		if next, ok := response["next"].(string); ok && next != "" {
			offset += 100
		} else {
			break
		}
	}
	return allTracks
}

func getTracks(playlistId string, offset int, token string) map[string]interface{} {
	url := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks?offset=%d", playlistId, offset)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("getTracksV2 - HTTP request error:", err)
		return nil
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		fmt.Println("getTracksV2 - JSON decode error:", err)
		return nil
	}

	return data
}

func getToken(c *gin.Context) {
	var requestBody struct {
		Action string `json:"action"`
		Code   string `json:"code"`
	}

	CLIENT_ID := os.Getenv("CLIENT_ID")
	CLIENT_SECRET := os.Getenv("CLIENT_SECRET")

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

	resp, err := httpClient.Do(req)
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

	resp, err := httpClient.Do(req)
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

	resp, err := httpClient.Do(req)
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

	resp, err := httpClient.Do(req)
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

	resp, err := httpClient.Do(req)
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

	resp, err := httpClient.Do(req)
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

func nextMusic(c *gin.Context) {
	var requestBody struct {
		MusicId    string `json:"MusicId"`
		PositionMs int    `json:"PositionMs"`
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

	url := fmt.Sprintf("https://api.spotify.com/v1/me/player/queue?uri=spotify%%3Atrack%%3A%s", requestBody.MusicId)
	postAPICall(url, token)

	url = "https://api.spotify.com/v1/me/player/queue"
	queue := getAPICall(url, token)

	url = "https://api.spotify.com/v1/me/player"
	response := getAPICall(url, token)
	device := response["device"].(map[string]interface{})
	volume := device["volume_percent"].(float64)

	url = "https://api.spotify.com/v1/me/player/volume?volume_percent=0"
	putAPICall(url, token)

	url = "https://api.spotify.com/v1/me/player/next"

	if items, ok := queue["queue"]; ok {
		for number, item := range items.([]interface{}) {
			if itemMap, ok := item.(map[string]interface{}); ok && itemMap["id"] == requestBody.MusicId {
				for i := 0; i < number; i++ {
					postAPICall(url, token)
				}
				break
			}
		}
	}

	postAPICall(url, token)

	url = fmt.Sprintf("https://api.spotify.com/v1/me/player/seek?position_ms=%d", requestBody.PositionMs)
	putAPICall(url, token)

	url = fmt.Sprintf("https://api.spotify.com/v1/me/player/volume?volume_percent=%.0f", volume)
	putAPICall(url, token)

	c.JSON(http.StatusOK, gin.H{"message": "Musique suivante"})
}

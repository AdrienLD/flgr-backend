package songs

import (
	"fmt"
	"math/rand/v2"
)

func AreaOfCircle(radius float64) float64 {
	return 3.14 * radius * radius
}

func ExtractAllPlaylist(playlists []string) []string {
	for i := range playlists {
		fmt.Println(playlists[i])
	}
	return playlists
}

func ReformatPlaylistResearch(genre string, playlistName string) string {
	switch genre {
	case "Années":
		return "Années " + playlistName
	case "Genres":
		return playlistName + " Classics"
	case "Artistes":
		return "this is " + playlistName
	case "Télévision":
		return "musiques de " + playlistName
	case "Français":
		return "gen " + playlistName + "Français"
	default:
		return "this is " + playlistName
	}
}

func EqualizePlaylists(playlists []map[string]interface{}) {
	// Trouver la taille de la plus petite playlist
	minSize := findMinPlaylist(playlists)

	// Réduire chaque playlist à la taille de la plus petite
	for _, playlist := range playlists {
		tracks, ok := playlist["tracks"].([]interface{})
		if !ok {
			fmt.Println("Erreur: format de playlist incorrect")
			continue
		}

		// Si la playlist est plus grande que minSize, réduire sa taille
		if len(tracks) > minSize {
			playlist["tracks"] = randomSample(tracks, minSize)
		}
	}

}

func randomSample(tracks []interface{}, n int) []interface{} {
	selected := make([]interface{}, n)
	perm := rand.Perm(len(tracks)) // Génère une permutation aléatoire

	for i := 0; i < n; i++ {
		selected[i] = tracks[perm[i]]
	}

	return selected
}

func findMinPlaylist(playlists []map[string]interface{}) int {
	minSize := -1
	for _, playlist := range playlists {
		tracks, ok := playlist["tracks"].([]interface{})
		if !ok {
			fmt.Println("Erreur: format de playlist incorrect")
			continue
		}
		if minSize == -1 || len(tracks) < minSize {
			minSize = len(tracks)
		}
	}
	return minSize
}

func MixAllTracks(playlists []map[string]interface{}) []interface{} {
	var allTracks []interface{}

	for _, playlist := range playlists {
		tracks, ok := playlist["tracks"].([]interface{})
		if !ok {
			fmt.Println("Erreur: format de playlist incorrect")
			continue
		}
		allTracks = append(allTracks, tracks...)
	}
	rand.Shuffle(len(allTracks), func(i, j int) {
		allTracks[i], allTracks[j] = allTracks[j], allTracks[i]
	})

	return allTracks
}

func GetFirstPlaylistID(data map[string]interface{}) map[string]interface{} {
	if playlists, ok := data["playlists"].(map[string]interface{}); ok {
		if items, ok := playlists["items"].([]interface{}); ok {
			for _, item := range items {
				if item != nil {
					if firstItem, ok := item.(map[string]interface{}); ok {
						return firstItem
					}
				}
			}
		}
	}
	return nil
}

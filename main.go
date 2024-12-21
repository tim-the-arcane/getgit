package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// GitHubContent repräsentiert die Struktur der GitHub API Antwort für Inhalte
type GitHubContent struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Type        string `json:"type"` // "file" oder "dir"
	DownloadURL string `json:"download_url"`
}

// downloadFile lädt eine Datei von der gegebenen URL herunter und speichert sie unter outputPath
func downloadFile(fileURL, outputPath string) error {
	// HTTP GET Anfrage
	resp, err := http.Get(fileURL)
	if err != nil {
		return fmt.Errorf("Fehler beim Herunterladen der Datei: %v", err)
	}
	defer resp.Body.Close()

	// Überprüfen, ob der HTTP Status OK ist
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Fehler: HTTP Status %s", resp.Status)
	}

	// Lokales Verzeichnis erstellen, falls nicht vorhanden
	err = os.MkdirAll(filepath.Dir(outputPath), os.ModePerm)
	if err != nil {
		return fmt.Errorf("Fehler beim Erstellen des Verzeichnisses: %v", err)
	}

	// Datei erstellen
	out, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("Fehler beim Erstellen der Datei: %v", err)
	}
	defer out.Close()

	// Inhalt kopieren
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("Fehler beim Schreiben der Datei: %v", err)
	}

	return nil
}

// downloadFolder lädt den Inhalt eines GitHub-Ordners rekursiv herunter
func downloadFolder(owner, repo, branch, path, localPath string) error {
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s?ref=%s", owner, repo, path, branch)

	// HTTP GET Anfrage zur GitHub API
	resp, err := http.Get(apiURL)
	if err != nil {
		return fmt.Errorf("Fehler beim Abrufen der API: %v", err)
	}
	defer resp.Body.Close()

	// Überprüfen, ob der HTTP Status OK ist
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Fehler: HTTP Status %s bei URL %s", resp.Status, apiURL)
	}

	// JSON-Antwort dekodieren
	var contents []GitHubContent
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&contents); err != nil {
		return fmt.Errorf("Fehler beim Dekodieren der JSON-Antwort: %v", err)
	}

	// Durchlaufe alle Inhalte und handle entsprechend
	for _, item := range contents {
		// Anpassung: Vermeide das Beibehalten der Repository-Verzeichnisstruktur
		localItemPath := filepath.Join(localPath, item.Name)
		if item.Type == "file" {
			fmt.Printf("Herunterladen der Datei: %s\n", item.DownloadURL)
			err := downloadFile(item.DownloadURL, localItemPath)
			if err != nil {
				return err
			}
			fmt.Printf("Datei erfolgreich heruntergeladen: %s\n", localItemPath)
		} else if item.Type == "dir" {
			fmt.Printf("Erstelle Verzeichnis: %s\n", localItemPath)
			err := os.MkdirAll(localItemPath, os.ModePerm)
			if err != nil {
				return fmt.Errorf("Fehler beim Erstellen des Verzeichnisses: %v", err)
			}
			// Rekursiver Aufruf für Unterordner
			err = downloadFolder(owner, repo, branch, item.Path, localItemPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// parseGitHubURL analysiert die GitHub-URL und gibt die Bestandteile zurück
func parseGitHubURL(githubURL string) (owner, repo, branch, path string, isFile bool, folderName string, err error) {
	parsedURL, err := url.Parse(githubURL)
	if err != nil {
		return
	}

	segments := strings.Split(parsedURL.Path, "/")
	if len(segments) < 4 {
		err = fmt.Errorf("Ungültige GitHub-URL")
		return
	}

	owner = segments[1]
	repo = segments[2]

	// Prüfen, ob es sich um eine Datei oder einen Ordner handelt
	if segments[3] == "blob" {
		isFile = true
		if len(segments) < 5 {
			err = fmt.Errorf("Ungültige GitHub-Datei-URL")
			return
		}
		branch = segments[4]
		path = strings.Join(segments[5:], "/")
		folderName = filepath.Base(path)
	} else if segments[3] == "tree" {
		isFile = false
		if len(segments) < 5 {
			// Standardbranch, wenn nicht angegeben
			branch = "main"
			path = strings.Join(segments[4:], "/")
		} else {
			branch = segments[4]
			path = strings.Join(segments[5:], "/")
		}
		folderName = filepath.Base(path)
	} else {
		err = fmt.Errorf("URL muss '/blob/' oder '/tree/' enthalten")
		return
	}

	return
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <GitHub-URL>")
		return
	}

	githubURL := os.Args[1]
	owner, repo, branch, path, isFile, folderName, err := parseGitHubURL(githubURL)
	if err != nil {
		fmt.Printf("Fehler beim Parsen der URL: %v\n", err)
		return
	}

	if isFile {
		// Datei herunterladen
		rawURL := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/%s", owner, repo, branch, path)
		fileName := filepath.Base(path)
		if fileName == "" {
			fileName = "downloaded_file"
		}
		fmt.Printf("Herunterladen der Datei: %s\n", rawURL)
		err := downloadFile(rawURL, fileName)
		if err != nil {
			fmt.Printf("Fehler: %v\n", err)
			return
		}
		fmt.Printf("Datei erfolgreich heruntergeladen: %s\n", fileName)
	} else {
		// Ordner herunterladen
		// Anpassung: Erstelle einen Ordner im aktuellen Arbeitsverzeichnis
		localPath := filepath.Join(".", folderName)
		fmt.Printf("Erstelle den Ordner: %s\n", localPath)
		err := os.MkdirAll(localPath, os.ModePerm)
		if err != nil {
			fmt.Printf("Fehler beim Erstellen des Ordners: %v\n", err)
			return
		}
		fmt.Printf("Herunterladen des Ordners: %s\n", path)
		err = downloadFolder(owner, repo, branch, path, localPath)
		if err != nil {
			fmt.Printf("Fehler: %v\n", err)
			return
		}
		fmt.Printf("Ordner erfolgreich heruntergeladen: %s\n", localPath)
	}
}

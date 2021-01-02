package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
)

// GetAssetsRoutes get all routes of path /assets
func GetAssetsRoutes(h *Handlers) []Route {
	routes := []Route{
		{
			"add asset",
			"POST",
			"/assets",
			h.addAsset,
		},
	}
	return routes
}

type asset struct {
	Name string `json:"name"`
}

func (h *Handlers) addAsset(w http.ResponseWriter, r *http.Request) {
	dir, _ := os.Getwd()

	// expects client formData key asset: formData.append("asset", file)
	file, fileheader, err := r.FormFile("asset")

	if err != nil {
		respondWithHTTP(w, http.StatusBadRequest)
		fmt.Fprintln(w, err)
		return
	}
	defer file.Close()

	filename := fileheader.Filename
	fp := path.Join(dir, "assets/images/", filename)

	out, err := os.Create(fp)
	if err != nil {
		respondWithHTTP(w, http.StatusBadRequest)
		fmt.Fprintf(w, "Unable to create the file for writing.")
		return
	}
	defer out.Close()

	// write the content from POST to the file
	_, err = io.Copy(out, file)
	if err != nil {
		respondWithHTTP(w, http.StatusBadRequest)
		fmt.Fprintln(w, err)
	}
	respondWithHTTP(w, http.StatusOK)

}

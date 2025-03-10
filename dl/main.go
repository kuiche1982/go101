package main

import (
    "io"
    "net/http"
    "os"
)

func main() {

    fileUrl := "https://grab-ucm-test.s3-ap-southeast-1.amazonaws.com/raw_food-cms-api/staging/blob/v20190301.0904.13000.food-cms-api"

    if err := DownloadFile("avatar.jpg", fileUrl); err != nil {
        panic(err)
    }
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) error {

    // Get the data
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    // Create the file
    out, err := os.Create(filepath)
    if err != nil {
        return err
    }
    defer out.Close()

    // Write the body to file
    _, err = io.Copy(out, resp.Body)
    return err
}

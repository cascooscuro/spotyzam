# spotyzam
Update a Spotify playlist based on the library of your shazams that can be downloaded as a csv file from https://www.shazam.com/myshazam

# Downloads
Releases for linux, windows and macOS in the release page https://github.com/cascooscuro/spotyzam/releases

 # Instructions:
  1. Register an application at: https://developer.spotify.com/my-applications/
   - Use 'http://localhost:8080/callback' as the redirect URI
  2. In your cmd/shell, set the SPOTIFY_ID environment variable to the client ID you got in step 1.
     - set SPOTIFY_ID=your-spotify-client-id  (from windows cmd)
     - export SPOTIFY_ID=your-spotify-client-id (linux/osx shell)
  3. Set the SPOTIFY_SECRET environment variable to the client secret from step 1.
     - set SPOTIFY_SECRET=your-spotify-client-secret  (from windows cmd)
     - export SPOTIFY_SECRET=your-spotify-client-secret (linux/osx shell)
  4. In Spotify, Create a playlist where you want your tracks to be added 
      - get the ID of that playlist by using the share button. 
      - The ID is the string that goes after https://open.spotify.com/playlist/ 
      - e.g.  in 'https://open.spotify.com/playlist/37i9dQZF1DX6YTj07PjLwE' the id is 37i9dQZF1DX6YTj07PjLwE  
  5. The csv file expects the format used by shazam: Index,TagTime,Title,Artist,URL,TrackKey. First line is skipped.
  6. Execute: spotyzam playlist-id csv-file.csv
      - e.g. C:\Users\user\Document>spotyzam.exe 37i9dQZF1DX6YTj07PjLwE shazamlibrary.csv


# Spotify Application settings

<img src="https://user-images.githubusercontent.com/5746813/153570984-66cce31b-d7a1-435e-895f-bffc768e3c38.png" width="400" height="400">

# Building/Compiling
1. Install golang  https://go.dev/doc/install
2. Use git to clone the github repo (git clone https://github.com/cascooscuro/spotyzam.git) or just download the code from https://github.com/cascooscuro/spotyzam/archive/refs/heads/main.zip and unzip it
3. in the directory where you have spotyzam.go do "go mod init spotyzam" and after that "go mod tidy"
4. Finally, to compile it: "go build spotyzam.go"
5. You should have now a binary "spotyzam" in your folder. Youn can run it by doing "./spotyzam" or "spotyzam.exe", depending on your OS
6. follow the instructions  to create the spotify application and playlist and set the environment variable to run it

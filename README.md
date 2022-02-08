# spotyzam
Update a Spotify playlist based on the library of your shazams, that can be downloaded as a csv file in https://www.shazam.com/myshazam

 Instructions:
  1. Register an application at: https://developer.spotify.com/my-applications/
   - Use 'http://localhost:8080/callback' as the redirect URI
  2. Set the SPOTIFY_ID environment variable to the client ID you got in step 1.
     - set SPOTIFY_ID=your-spotify-client-id  (from windows cmd)
     - export SPOTIFY_ID=your-spotify-client-id (linux shell)
  3. Set the SPOTIFY_SECRET environment variable to the client secret from step 1.
     - set SPOTIFY_SECRET=your-spotify-client-secret  (from windows cmd)
     - export SPOTIFY_SECRET=your-spotify-client-secret (linux shell)
  4. In Spotify, Create a playlist where you want your tracks to be added 
      - get the ID of that playlist by using the share button. 
      - The ID is the string that goes after https://open.spotify.com/playlist/ 
      - e.g.  in 'https://open.spotify.com/playlist/37i9dQZF1DX6YTj07PjLwE' the id is 37i9dQZF1DX6YTj07PjLwE  
  5. The csv file expects the format used by shazam: Index,TagTime,Title,Artist,URL,TrackKey. First line is skipped.
  6. Execute: spotyzam playlist-id csv-file.csv
      - e.g. C:\Users\user\Document>spotyzam.exe 6R123456JTOr75ntrtymmk shazamlibrary.csv


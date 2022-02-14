//
//  1. Register an application at: https://developer.spotify.com/my-applications/
//       - Use "http://localhost:8080/callback" as the redirect URI
//  2. Set the SPOTIFY_ID environment variable to the client ID you got in step 1.
//       set SPOTIFY_ID=your-spotify-client-id  (from windows cmd)
//               export SPOTIFY_ID=your-spotify-client-id (linux shell)
//  3. Set the SPOTIFY_SECRET environment variable to the client secret from step 1.
//       set SPOTIFY_SECRET=your-spotify-client-secret  (from windows cmd)
//               export SPOTIFY_SECRET=your-spotify-client-secret (linux shell)

package main

import (
        "context"
        "fmt"
        "github.com/zmb3/spotify/v2/auth"
        "log"
        "net/http"
        "os"
        "encoding/csv"
        "bufio"
        "io"
        "github.com/zmb3/spotify/v2"
        "strings"
        "regexp"
)

// redirectURI is the OAuth redirect URI for the application.
// You must register an application at Spotify's developer portal
// and enter this value.
const redirectURI = "http://localhost:8080/callback"

var (
        auth  = spotifyauth.New(spotifyauth.WithRedirectURL(redirectURI), spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate, spotifyauth.ScopePlaylistModifyPrivate, spotifyauth.ScopePlaylistModifyPublic ))
        ch    = make(chan *spotify.Client)
        state = "abc123"
)

func main() {


        helptext:= "ERROR: invalid option use 'spotyzam playlist-id csv-file.csv' \n" +
        " Instructions: \n" + 
        "1. Register an application at: https://developer.spotify.com/my-applications/ \n" + 
        "   - Use 'http://localhost:8080/callback' as the redirect URI \n" +
        "2. Set the SPOTIFY_ID environment variable to the client ID you got in step 1. \n" +
        "set SPOTIFY_ID=your-spotify-client-id  (from windows cmd) \n" +
        "export SPOTIFY_ID=your-spotify-client-id (linux shell) \n" +
        "3. Set the SPOTIFY_SECRET environment variable to the client secret from step 1. \n" +
        "set SPOTIFY_SECRET=your-spotify-client-secret  (from windows cmd) \n" +
        "export SPOTIFY_SECRET=your-spotify-client-secret (linux shell) \n" +
        "4. Create a playlist where you want your tracks to be added \n" + 
        "    get the ID of that playlist by using the share button. \n" + 
        "    The ID is the string that goes after https://open.spotify.com/playlist/ \n" + 
        "    e.g.  in 'https://open.spotify.com/playlist/37i9dQZF1DX6YTj07PjLwE' the id is 37i9dQZF1DX6YTj07PjLwE  \n" + 
        "5. The csv file expects the format used by shazam: Index,TagTime,Title,Artist,URL,TrackKey. We skip the first line.  \n"

        if len(os.Args) < 3 {
                fmt.Println(helptext)
                os.Exit(0)   
        }

        // first start an HTTP server
        http.HandleFunc("/callback", completeAuth)
        http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
                log.Println("Got request for:", r.URL.String())
        })

        go func() {
                err := http.ListenAndServe(":8080", nil)
                if err != nil {
                        log.Fatal(err)
                }
        }()

        url := auth.AuthURL(state)
        fmt.Println("")
        fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

        // wait for auth to complete
        client := <-ch

        // use the client to make calls that require authorization
        user, err := client.CurrentUser(context.Background())
        if err != nil {
                log.Fatal(err)
        }
        fmt.Println("You are logged in as:", user.ID)

        //read args from cmd-line
        PlaylistId := os.Args[1]
        Filename := os.Args[2]

        //open csv file and read rows
        f, err := os.Open(Filename)
        if err != nil {
                panic(err)
        }
        defer f.Close()
        records, err := readSample(f)
        if err != nil {
        panic(err)
        }

        //get details of playlist
        results, err := client.GetPlaylist(context.Background(), spotify.ID(PlaylistId))

        if err != nil {
                log.Fatal(err)
        }
   
        fmt.Println("\n")
        fmt.Println("Playlist Name:", results.Name)
        fmt.Println("Playlist Description:", results.Description)
        fmt.Println("Playlist TOTAL Tracks:", results.Tracks.Total)
        fmt.Println("\n")

        results_tracks, err := client.GetPlaylistTracks(context.Background(), spotify.ID(PlaylistId))

        if err != nil {
                log.Fatal(err)
        }

        type track struct {
            albumname string
            trackname string
            trackid spotify.ID
        }

        type notfoundtrack struct {
            csvtext string
            artist string
            track string
        }



        var tracks_in_playlist []track
        var not_found_tracks []notfoundtrack

        if results_tracks.Tracks != nil && len(results_tracks.Tracks) >0 {

                for _, item := range results_tracks.Tracks  {
                        tracks_in_playlist = append(tracks_in_playlist,track{item.Track.Album.Name,item.Track.Name, item.Track.ID} )
                }
        
                for page := 1; ; page++ {
                        err = client.NextPage(context.Background(), results_tracks)
                        if err == spotify.ErrNoMorePages {
                                break
                        }
                        if err != nil {
                                log.Fatal(err)
                        }
                        for _, item := range results_tracks.Tracks  {
                                tracks_in_playlist = append(tracks_in_playlist,track{item.Track.Album.Name,item.Track.Name, item.Track.ID} )
                        }
                }
        }

        var tracklist []spotify.ID

        fmt.Println("")
        fmt.Println("There are: ", len(records), " tracks in csv file")
        fmt.Println("----Searching them via API ----" )
    
        
        for _, item := range records  {

                //item[2]=track
                //item[3]=artist
                var doweappend int = 1
                
                //remove chars from artist
                //remove everythingafter  Feat.
                reg1 := regexp.MustCompile(`Feat\..*$`)
                res1 := reg1.ReplaceAllString(item[3], "")
                temp2 := strings.Replace(res1, "&", " ", -1)
                temp3 := strings.Replace(temp2, ",", " ", -1)
                fin_artist := temp3

                //remove everything inside () or []
                reg2 := regexp.MustCompile(`[\(\[].*?[\)\]]`)
                res2 := reg2.ReplaceAllString(item[2], "")
                //remove bad words thant contains ***
                reg3 := regexp.MustCompile(`[^ ]*\*\*\*[^ ]*`)
                res3 := reg3.ReplaceAllString(res2, "")
                //remove bad words thant contains **
                reg4 := regexp.MustCompile(`[^ ]*\*\*[^ ]*`)
                res4 := reg4.ReplaceAllString(res3, "")
                fin_track := res4
                                
                //q:= "track:" + fin_track + " " + "artist:" +  fin_artist
                //searching without field filters gives better results
                q := fin_track + " " + fin_artist
                                
                fmt.Println("CSV:",item)
                fmt.Println("QUERY:",q)
                
                results, err := client.Search(context.Background(), q, spotify.SearchTypeTrack)

                if err != nil {
                    log.Fatal(err)
                }

                // handle search results
                if results.Tracks != nil && len(results.Tracks.Tracks) >0 {
                
                        for i := range tracks_in_playlist {
                                
                                if ((tracks_in_playlist[i].albumname == results.Tracks.Tracks[0].Album.Name) && (tracks_in_playlist[i].trackname == results.Tracks.Tracks[0].Name )) {
                                // track already in playlist!
                                doweappend = 0
                                fmt.Println("ALREADY in playlist")
                                fmt.Println("------------")
                                break
                                } 
                        }
                                        
                        if doweappend == 1 {                    
                                fmt.Println("FOUND NEW TRACK ::","Artist:",results.Tracks.Tracks[0].Artists[0].Name, "----",  "Album:",results.Tracks.Tracks[0].Album.Name, "----", "Track:",results.Tracks.Tracks[0].Name,"----", "ID:", results.Tracks.Tracks[0].ID )
                                tracklist = append(tracklist,  spotify.ID(results.Tracks.Tracks[0].ID)  )
                                fmt.Println("------------")
                        }
                }
                //track not found
                if results.Tracks != nil && len(results.Tracks.Tracks) == 0 {
                        not_found_tracks = append(not_found_tracks,notfoundtrack{strings.Join(item, ","),fin_artist, fin_track} )
                        fmt.Println("NOT FOUND")
                        fmt.Println("------------")
                }
        }

        if len(tracklist) == 0 {
                fmt.Println("\n")
                fmt.Println("----Nothing to do. There aren't new tracks to add to the playlist." )
                fmt.Println("\n")
                fmt.Println("There are: ", len(not_found_tracks), " tracks that were not found:")
                for _, item := range not_found_tracks  {
                        fmt.Println("NOT FOUND CSV TEXT ::", item.csvtext )
                        fmt.Println("NOT FOUND ::", "Artist:",item.artist, "----", "Track:",item.track )
        }
                os.Exit(0)

        } else {
                fmt.Println("")
                fmt.Println("There are: ", len(tracklist), " tracks to update in the playlist")
                fmt.Println("----Let's add tracks via API ----" )
        }

        // Split the slice into batches of 100 items.
        batch := 100

        fmt.Print("Progress: ")
        for i := 0; i < len(tracklist); i += batch {
                j := i + batch
                if j > len(tracklist) {
                        j = len(tracklist)
                }

                _, err := client.AddTracksToPlaylist(context.Background(), spotify.ID(PlaylistId), tracklist[i:j]...)
                fmt.Print( len(tracklist[i:j]) , ",")
                
                if err != nil {
                log.Fatal(err)
                }
        }
                
        fmt.Println("")
        fmt.Println("----Tracks Updated!! ----" )
        results_playlist_update, err := client.GetPlaylist(context.Background(), spotify.ID(PlaylistId))

        if err != nil {
                log.Fatal(err)
        }
        fmt.Println("\n")
        fmt.Println("Playlist TOTAL Tracks after the update:", results_playlist_update.Tracks.Total)
        fmt.Println("\n")

        fmt.Println("There are: ", len(not_found_tracks), " tracks that were not found:")
        for _, item := range not_found_tracks  {
                fmt.Println("NOT FOUND CSV TEXT ::", item.csvtext )
                fmt.Println("NOT FOUND ::", "Artist:",item.artist, "----", "Track:",item.track )
        }
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
        tok, err := auth.Token(r.Context(), state, r)
        if err != nil {
                http.Error(w, "Couldn't get token", http.StatusForbidden)
                log.Fatal(err)
        }
        if st := r.FormValue("state"); st != state {
                http.NotFound(w, r)
                log.Fatalf("State mismatch: %s != %s\n", st, state)
        }

        // use the token to get an authenticated client
        client := spotify.New(auth.Client(r.Context(), tok),spotify.WithRetry(true))
        fmt.Fprintf(w, "Login Completed!")
        ch <- client
}

func readSample(rs io.ReadSeeker) ([][]string, error) {
    // Skip first row (line)
    row1, err := bufio.NewReader(rs).ReadSlice('\n')
    if err != nil {
        return nil, err
    }
    _, err = rs.Seek(int64(len(row1)), io.SeekStart)
    if err != nil {
        return nil, err
    }

    // Read remaining rows
    r := csv.NewReader(rs)
    r.LazyQuotes = true
    rows, err := r.ReadAll()
    if err != nil {
        return nil, err
    }
    return rows, nil
}

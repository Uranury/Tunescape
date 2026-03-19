### Overall Scopes:
For GET:
- user-read-private 
- user-read-email 
- user-top-read 
- playlist-read-private 
- user-library-read 
- user-follow-read

For PUT:
- user-library-modify 
- user-follow-modify 
- playlist-modify-public
- playlist-modify-private
- ugc-image-upload

For POST:
- playlist-modify-public 
- playlist-modify-private

For DELETE:
- user-library-modify 
- user-follow-modify 
- playlist-modify-public 
- playlist-modify-private

## GET methods
### Main objects in responses:
- `TrackObject` - informetion about track: name, artist, album, duration etc.
  Get Track - `GET /v1/tracks/{id}`

  Response also contains SimplifiedArtistObject

  <details>
  <summary>JSON Object</summary>
  
  ```bash
  {
    "album": {
        "album_type": "album",
        "artists": [ 
            {
                "external_urls": {
                    "spotify": "https://open.spotify.com/artist/7dGJo4pcD2V6oG8kP0tJRR"
                },
                "href": "https://api.spotify.com/v1/artists/7dGJo4pcD2V6oG8kP0tJRR",
                "id": "7dGJo4pcD2V6oG8kP0tJRR",
                "name": "Eminem",
                "type": "artist",
                "uri": "spotify:artist:7dGJo4pcD2V6oG8kP0tJRR"
            }
        ],
        "external_urls": {
            "spotify": "https://open.spotify.com/album/1kTlYbs28MXw7hwO0NLYif"
        },
        "href": "https://api.spotify.com/v1/albums/1kTlYbs28MXw7hwO0NLYif",
        "id": "1kTlYbs28MXw7hwO0NLYif",
        "images": [
            {
                "url": "https://i.scdn.co/image/ab67616d0000b2731bec21e57fff76db49e15a70",
                "width": 640,
                "height": 640
            },
            {
                "url": "https://i.scdn.co/image/ab67616d00001e021bec21e57fff76db49e15a70",
                "width": 300,
                "height": 300
            },
            {
                "url": "https://i.scdn.co/image/ab67616d000048511bec21e57fff76db49e15a70",
                "width": 64,
                "height": 64
            }
        ],
        "is_playable": true,
        "name": "Encore (Deluxe Version)",
        "release_date": "2004-11-12",
        "release_date_precision": "day",
        "total_tracks": 23,
        "type": "album",
        "uri": "spotify:album:1kTlYbs28MXw7hwO0NLYif"
    },
    "artists": [ 
        {
            "external_urls": {
                "spotify": "https://open.spotify.com/artist/7dGJo4pcD2V6oG8kP0tJRR"
            },
            "href": "https://api.spotify.com/v1/artists/7dGJo4pcD2V6oG8kP0tJRR",
            "id": "7dGJo4pcD2V6oG8kP0tJRR",
            "name": "Eminem",
            "type": "artist",
            "uri": "spotify:artist:7dGJo4pcD2V6oG8kP0tJRR"
        }
    ],
    "disc_number": 1,
    "duration_ms": 250760,
    "explicit": true,
    "external_ids": {
        "isrc": "USIR10400813"
    },
    "external_urls": {
        "spotify": "https://open.spotify.com/track/561jH07mF1jHuk7KlaeF0s"
    },
    "href": "https://api.spotify.com/v1/tracks/561jH07mF1jHuk7KlaeF0s",
    "id": "561jH07mF1jHuk7KlaeF0s",
    "is_local": false,
    "is_playable": true,
    "name": "Mockingbird",
    "track_number": 16,
    "type": "track",
    "uri": "spotify:track:561jH07mF1jHuk7KlaeF0s"
  }
  
  ```
  </details>


- `ArtistObject` - information about artist: name, genre, followers, images.
  Get Artist - `GET /v1/artists/{id}`

  Response also contains ImageObject
  
  <details>
  <summary>JSON Object</summary>

  ```bash
  {
    "artists": { 
        "href": "https://api.spotify.com/v1/search?offset=0&limit=5&query=eminem&type=artist",
        "limit": 5,
        "next": "https://api.spotify.com/v1/search?offset=5&limit=5&query=eminem&type=artist",
        "offset": 0,
        "previous": null,
        "total": 877,
        "items": [
            {
                "external_urls": {
                    "spotify": "https://open.spotify.com/artist/7dGJo4pcD2V6oG8kP0tJRR"
                },
                "href": "https://api.spotify.com/v1/artists/7dGJo4pcD2V6oG8kP0tJRR",
                "id": "7dGJo4pcD2V6oG8kP0tJRR",
                "images": [
                    {
                        "url": "https://i.scdn.co/image/ab6761610000e5eba00b11c129b27a88fc72f36b",
                        "height": 640,
                        "width": 640
                    },
                    {
                        "url": "https://i.scdn.co/image/ab67616100005174a00b11c129b27a88fc72f36b",
                        "height": 320,
                        "width": 320
                    },
                    {
                        "url": "https://i.scdn.co/image/ab6761610000f178a00b11c129b27a88fc72f36b",
                        "height": 160,
                        "width": 160
                    }
                ],
                "name": "Eminem",
                "type": "artist",
                "uri": "spotify:artist:7dGJo4pcD2V6oG8kP0tJRR"
            },
            {
                "external_urls": {
                    "spotify": "https://open.spotify.com/artist/15UsOTVnJzReFVN1VCnxy4"
                },
                "href": "https://api.spotify.com/v1/artists/15UsOTVnJzReFVN1VCnxy4",
                "id": "15UsOTVnJzReFVN1VCnxy4",
                "images": [
                    {
                        "url": "https://i.scdn.co/image/ab6761610000e5ebf0c20db5ef6c6fbe5135d2e4",
                        "height": 640,
                        "width": 640
                    },
                    {
                        "url": "https://i.scdn.co/image/ab67616100005174f0c20db5ef6c6fbe5135d2e4",
                        "height": 320,
                        "width": 320
                    },
                    {
                        "url": "https://i.scdn.co/image/ab6761610000f178f0c20db5ef6c6fbe5135d2e4",
                        "height": 160,
                        "width": 160
                    }
                ],
                "name": "XXXTENTACION",
                "type": "artist",
                "uri": "spotify:artist:15UsOTVnJzReFVN1VCnxy4"
            },
            {
                "external_urls": {
                    "spotify": "https://open.spotify.com/artist/1Xyo4u8uXC1ZmMpatF05PJ"
                },
                "href": "https://api.spotify.com/v1/artists/1Xyo4u8uXC1ZmMpatF05PJ",
                "id": "1Xyo4u8uXC1ZmMpatF05PJ",
                "images": [
                    {
                        "url": "https://i.scdn.co/image/ab6761610000e5eb9e528993a2820267b97f6aae",
                        "height": 640,
                        "width": 640
                    },
                    {
                        "url": "https://i.scdn.co/image/ab676161000051749e528993a2820267b97f6aae",
                        "height": 320,
                        "width": 320
                    },
                    {
                        "url": "https://i.scdn.co/image/ab6761610000f1789e528993a2820267b97f6aae",
                        "height": 160,
                        "width": 160
                    }
                ],
                "name": "The Weeknd",
                "type": "artist",
                "uri": "spotify:artist:1Xyo4u8uXC1ZmMpatF05PJ"
            },
            {
                "external_urls": {
                    "spotify": "https://open.spotify.com/artist/5TZ0uM0InXX4LtI8hTsXlv"
                },
                "href": "https://api.spotify.com/v1/artists/5TZ0uM0InXX4LtI8hTsXlv",
                "id": "5TZ0uM0InXX4LtI8hTsXlv",
                "images": [
                    {
                        "url": "https://i.scdn.co/image/ab6761610000e5ebb1836e93e000b222ccb76f47",
                        "height": 640,
                        "width": 640
                    },
                    {
                        "url": "https://i.scdn.co/image/ab67616100005174b1836e93e000b222ccb76f47",
                        "height": 320,
                        "width": 320
                    },
                    {
                        "url": "https://i.scdn.co/image/ab6761610000f178b1836e93e000b222ccb76f47",
                        "height": 160,
                        "width": 160
                    }
                ],
                "name": "Bakr",
                "type": "artist",
                "uri": "spotify:artist:5TZ0uM0InXX4LtI8hTsXlv"
            },
            {
                "external_urls": {
                    "spotify": "https://open.spotify.com/artist/4oUcWvCNxqNZv4l7BXlE0y"
                },
                "href": "https://api.spotify.com/v1/artists/4oUcWvCNxqNZv4l7BXlE0y",
                "id": "4oUcWvCNxqNZv4l7BXlE0y",
                "images": [
                    {
                        "url": "https://i.scdn.co/image/ab6761610000e5eb3af4f3c1280d72345bb66c57",
                        "height": 640,
                        "width": 640
                    },
                    {
                        "url": "https://i.scdn.co/image/ab676161000051743af4f3c1280d72345bb66c57",
                        "height": 320,
                        "width": 320
                    },
                    {
                        "url": "https://i.scdn.co/image/ab6761610000f1783af4f3c1280d72345bb66c57",
                        "height": 160,
                        "width": 160
                    }
                ],
                "name": "Mount Eminest",
                "type": "artist",
                "uri": "spotify:artist:4oUcWvCNxqNZv4l7BXlE0y"
            }
        ]
    }
  }

  ```

  </details>

- `AlbumObject` - information about album: name, release date, total_tracks, artists, images.
  Get Album - `GET /v1/albums/{id}`

  Response also contains CopyrightObject, SimplifiedTrackObject, SimplifiedArtistObject, ImageObject
 
  <details>
  <summary>JSON Object</summary>
  
  ```bash
  {
      "album_type": "single",
      "total_tracks": 3,
      "external_urls": {
          "spotify": "https://open.spotify.com/album/46eOftG6eRzvFa1OKPzmMJ"
      },
      "href": "https://api.spotify.com/v1/albums/46eOftG6eRzvFa1OKPzmMJ",
      "id": "46eOftG6eRzvFa1OKPzmMJ",
      "images": [
          {
              "url": "https://i.scdn.co/image/ab67616d0000b27347c680e39a6cf1223355394a",
              "height": 640,
              "width": 640
          },
          {
              "url": "https://i.scdn.co/image/ab67616d00001e0247c680e39a6cf1223355394a",
              "height": 300,
              "width": 300
          },
          {
              "url": "https://i.scdn.co/image/ab67616d0000485147c680e39a6cf1223355394a",
              "height": 64,
              "width": 64
          }
      ],
      "name": "Live at Ford Field",
      "release_date": "2025-11-27",
      "release_date_precision": "day",
      "type": "album",
      "uri": "spotify:album:46eOftG6eRzvFa1OKPzmMJ",
      "artists": [ 
          {
              "external_urls": {
                  "spotify": "https://open.spotify.com/artist/4FZ3j1oH43e7cukCALsCwf"
              },
              "href": "https://api.spotify.com/v1/artists/4FZ3j1oH43e7cukCALsCwf",
              "id": "4FZ3j1oH43e7cukCALsCwf",
              "name": "Jack White",
              "type": "artist",
              "uri": "spotify:artist:4FZ3j1oH43e7cukCALsCwf"
          },
          {
              "external_urls": {
                  "spotify": "https://open.spotify.com/artist/7dGJo4pcD2V6oG8kP0tJRR"
              },
              "href": "https://api.spotify.com/v1/artists/7dGJo4pcD2V6oG8kP0tJRR",
              "id": "7dGJo4pcD2V6oG8kP0tJRR",
              "name": "Eminem",
              "type": "artist",
              "uri": "spotify:artist:7dGJo4pcD2V6oG8kP0tJRR"
          }
      ],
      "tracks": {
          "href": "https://api.spotify.com/v1/albums/46eOftG6eRzvFa1OKPzmMJ/tracks?offset=0&limit=50",
          "limit": 50,
          "next": null,
          "offset": 0,
          "previous": null,
          "total": 3,
          "items": [ //`SimplifiedTrackObject`
              {
                  "artists": [ 
                      {
                          "external_urls": {
                              "spotify": "https://open.spotify.com/artist/4FZ3j1oH43e7cukCALsCwf"
                          },
                          "href": "https://api.spotify.com/v1/artists/4FZ3j1oH43e7cukCALsCwf",
                          "id": "4FZ3j1oH43e7cukCALsCwf",
                          "name": "Jack White",
                          "type": "artist",
                          "uri": "spotify:artist:4FZ3j1oH43e7cukCALsCwf"
                      }
                  ],
                  "disc_number": 1,
                  "duration_ms": 132725,
                  "explicit": false,
                  "external_urls": {
                      "spotify": "https://open.spotify.com/track/6JD6pSfTa6jAumxb0lc9CJ"
                  },
                  "href": "https://api.spotify.com/v1/tracks/6JD6pSfTa6jAumxb0lc9CJ",
                  "id": "6JD6pSfTa6jAumxb0lc9CJ",
                  "name": "That's How I'm Feeling - Live",
                  "track_number": 1,
                  "type": "track",
                  "uri": "spotify:track:6JD6pSfTa6jAumxb0lc9CJ",
                  "is_local": false
              },
              {
                  "artists": [ 
                      {
                          "external_urls": {
                              "spotify": "https://open.spotify.com/artist/4FZ3j1oH43e7cukCALsCwf"
                          },
                          "href": "https://api.spotify.com/v1/artists/4FZ3j1oH43e7cukCALsCwf",
                          "id": "4FZ3j1oH43e7cukCALsCwf",
                          "name": "Jack White",
                          "type": "artist",
                          "uri": "spotify:artist:4FZ3j1oH43e7cukCALsCwf"
                      },
                      {
                          "external_urls": {
                              "spotify": "https://open.spotify.com/artist/7dGJo4pcD2V6oG8kP0tJRR"
                          },
                          "href": "https://api.spotify.com/v1/artists/7dGJo4pcD2V6oG8kP0tJRR",
                          "id": "7dGJo4pcD2V6oG8kP0tJRR",
                          "name": "Eminem",
                          "type": "artist",
                          "uri": "spotify:artist:7dGJo4pcD2V6oG8kP0tJRR"
                      }
                  ],
                  "disc_number": 1,
                  "duration_ms": 138197,
                  "explicit": false,
                  "external_urls": {
                      "spotify": "https://open.spotify.com/track/4XWhHCDrr1r1SPgh9GeQKu"
                  },
                  "href": "https://api.spotify.com/v1/tracks/4XWhHCDrr1r1SPgh9GeQKu",
                  "id": "4XWhHCDrr1r1SPgh9GeQKu",
                  "name": "Hello Operator / 'Till I Collapse - Live",
                  "track_number": 2,
                  "type": "track",
                  "uri": "spotify:track:4XWhHCDrr1r1SPgh9GeQKu",
                  "is_local": false
              },
              {
                  "artists": [ 
                      {
                          "external_urls": {
                              "spotify": "https://open.spotify.com/artist/4FZ3j1oH43e7cukCALsCwf"
                          },
                          "href": "https://api.spotify.com/v1/artists/4FZ3j1oH43e7cukCALsCwf",
                          "id": "4FZ3j1oH43e7cukCALsCwf",
                          "name": "Jack White",
                          "type": "artist",
                          "uri": "spotify:artist:4FZ3j1oH43e7cukCALsCwf"
                      }
                  ],
                  "disc_number": 1,
                  "duration_ms": 174773,
                  "explicit": false,
                  "external_urls": {
                      "spotify": "https://open.spotify.com/track/78WxnmK9sOrRB5XOnISlFn"
                  },
                  "href": "https://api.spotify.com/v1/tracks/78WxnmK9sOrRB5XOnISlFn",
                  "id": "78WxnmK9sOrRB5XOnISlFn",
                  "name": "Seven Nation Army - Live",
                  "track_number": 3,
                  "type": "track",
                  "uri": "spotify:track:78WxnmK9sOrRB5XOnISlFn",
                  "is_local": false
              }
          ]
      },
      "copyrights": [
          {
              "text": "© 2025 Third Face LLC & Marshall B. Mathers III",
              "type": "C"
          },
          {
              "text": "℗ 2025 Third Face LLC & Marshall B. Mathers III",
              "type": "P"
          }
      ],
      "external_ids": {
          "upc": "00199957165518"
      },
      "genres": []
  }
  ```

  </details>

## Another objects


### `SimplifiedArtistObject`
<details>
  <summary>Sample</summary>

  ```bash
  "artists": [
    {
    "external_urls": {
      "spotify": "string"
      },
    "href": "string",
    "id": "string",
    "name": "string",
    "type": "artist",
    "uri": "string"
    }
  ]
  ```
</details>

<details>
  <summary>Occur</summary>

  - TrackObject → artists[]
  - TrackObject.album.artists[]
  - AlbumObject.artists[]
  - AlbumObject.tracks.items[].artists[]
  - Search for Item → albums.items[].artists[]
  - Get Album Tracks → items[].artists[]
  - Get User's Saved Albums → items[].album.artists[]
  - Get Artist's Albums → items[].artists[]
  - Get Playlist → items.items[].item.artists[] и items.items[].item.album.artists[]
  - Get Playlist Items → items[].item.artists[] и items[].item.album.artists[]
  - Get User's Saved Tracks → items[].track.artists[] и items[].track.album.artists[]
  - Get User's Top Items → items[].artists[] и items[].album.artists[]

</details>


### `SimplifiedTrackObject`
<details>
  <summary>Sample</summary>

  ```bash
  "items": [
  {
    "artists": [
      {
        "external_urls": {
          "spotify": "string"
        },
        "href": "string",
        "id": "string",
        "name": "string",
        "type": "artist",
        "uri": "string"
      }
    ],
    "available_markets": [
      "string"
    ],
    "disc_number": 0,
    "duration_ms": 0,
    "explicit": false,
    "external_urls": {
      "spotify": "string"
    },
    "href": "string",
    "id": "string",
    "is_playable": false,
    "linked_from": {
      "external_urls": {
        "spotify": "string"
      },
      "href": "string",
      "id": "string",
      "type": "string",
      "uri": "string"
    },
    "restrictions": {
      "reason": "string"
    },
    "name": "string",
    "preview_url": "string",
    "track_number": 0,
    "type": "string",
    "uri": "string",
    "is_local": false
  }
  ]
  ```
</details>

<details>
  <summary>Occur</summary>

  - Get Album.tracks.items[] → track inside album
  - Get Album Tracks → items[]
  - Get Playlist.items.items[].item
  - Get Playlist Items.items[].item
  - Get User's Saved Tracks.items[].track
  - Get User's Top Items.items[]
  
</details>

### `SimplifiedAlbumObject`
<details>
  <summary>Sample</summary>

  ```bash
  "items": [
  {
    "album_type": "compilation",
    "total_tracks": 9,
    "available_markets": [
      "CA",
      "BR",
      "IT"
    ],
    "external_urls": {
      "spotify": "string"
    },
    "href": "string",
    "id": "2up3OPMp9Tb4dAKM2erWXQ",
    "images": [
      {
        "url": "https://i.scdn.co/image/ab67616d00001e02ff9ca10b55ce82ae553c8228",
        "height": 300,
        "width": 300
      }
    ],
    "name": "string",
    "release_date": "1981-12",
    "release_date_precision": "year",
    "restrictions": {
      "reason": "market"
    },
    "type": "album",
    "uri": "spotify:album:2up3OPMp9Tb4dAKM2erWXQ",
    "artists": [
      {
        "external_urls": {
          "spotify": "string"
        },
        "href": "string",
        "id": "string",
        "name": "string",
        "type": "artist",
        "uri": "string"
      }
    ],
    "album_group": "compilation"
  }
  ]
  ```
</details>

<details>
  <summary>Occur</summary>

  - Search for Item → albums.items[]
  - Get User's Saved Albums.items[].album
  - Get Artist's Albums.items[]
  - Get Playlist.items.items[].item.album
  - Get Playlist Items.items[].item.album
  - Get User's Saved Tracks.items[].track.album
  - Get User's Top Items.items[].album
  - Get Track's TrackObject.album
  
</details>

### `ImageObject`
<details>
  <summary>Sample</summary>

  ```bash
  "images": [
    {
      "url": "https://i.scdn.co/image/ab67616d00001e02ff9ca10b55ce82ae553c8228",
      "height": 300,
      "width": 300
    }
  ]
  ```
</details>

<details>
  <summary>Occur</summary>

  - ArtistObject.images[]
  - AlbumObject.images[]
  - Search for Item.albums.items[].images[]
  - Get User's Saved Albums.items[].album.images[]
  - Get Artist's Albums.items[].images[]
  - Get Playlist.images[]
  - Get Playlist.items.items[].item.album.images[]
  - Get Playlist Items.items[].item.album.images[]
  - Get User's Saved Tracks.items[].track.album.images[]
  - Get Current User's Profile.images[]
  - Get Followed Artists.artists.items[].images[]
  - Get User's Top Items.items[].album.images[]
  
</details>

<!-- ExternalUrlsObject
<details>
  <summary>Sample</summary>

  ```bash
  "external_urls": {
  "spotify": "string"
  }
  ```
</details> -->

### `SimplifiedPlaylistObject`
<details>
  <summary>Sample</summary>

  ```bash
  [
  {
    "collaborative": false,
    "description": "string",
    "external_urls": {
      "spotify": "string"
    },
    "href": "string",
    "id": "string",
    "images": [
      {
        "url": "https://i.scdn.co/image/ab67616d00001e02ff9ca10b55ce82ae553c8228",
        "height": 300,
        "width": 300
      }
    ],
    "name": "string",
    "owner": {
      "external_urls": {
        "spotify": "string"
      },
      "href": "string",
      "id": "string",
      "type": "user",
      "uri": "string",
      "display_name": "string"
    },
    "public": false,
    "snapshot_id": "string",
    "items": {
      "href": "string",
      "total": 0
    },
    "tracks": {
      "href": "string",
      "total": 0
    },
    "type": "string",
    "uri": "string"
  }
    ]
  ```
</details>

<details>
  <summary>Occur</summary>

  - Get Get Current User's Playlists
  
</details>

### `SimplifiedUserObject`
<details>
  <summary>Sample</summary>

  ```bash
  "owner": {
    "external_urls": {
      "spotify": "string"
    },
    "href": "string",
    "id": "string",
    "type": "user",
    "uri": "string",
    "display_name": "string"
  }
  ```
</details>

<details>
  <summary>Occur</summary>

  - Get Playlist.owner
  - Get Playlist.items.items[].added_by
  - Get Current User's Playlists.items[].owner
  - Get Playlist Items.items[].added_by
  - Get Current User's Profile
  
</details>

### `PlaylistTrackObject`
<details>
  <summary>Sample</summary>

  ```bash
  "items": [
  {
    "added_at": "string",
    "added_by": {
      "external_urls": {
        "spotify": "string"
      },
      "href": "string",
      "id": "string",
      "type": "user",
      "uri": "string"
    },
    "is_local": false,
    "item": {
      "album": {
        "album_type": "compilation",
        "total_tracks": 9,
        "available_markets": [
          "CA",
          "BR",
          "IT"
        ],
        "external_urls": {
          "spotify": "string"
        },
        "href": "string",
        "id": "2up3OPMp9Tb4dAKM2erWXQ",
        "images": [
          {
            "url": "https://i.scdn.co/image/ab67616d00001e02ff9ca10b55ce82ae553c8228",
            "height": 300,
            "width": 300
          }
        ],
        "name": "string",
        "release_date": "1981-12",
        "release_date_precision": "year",
        "restrictions": {
          "reason": "market"
        },
        "type": "album",
        "uri": "spotify:album:2up3OPMp9Tb4dAKM2erWXQ",
        "artists": [
          {
            "external_urls": {
              "spotify": "string"
            },
            "href": "string",
            "id": "string",
            "name": "string",
            "type": "artist",
            "uri": "string"
          }
        ]
      },
      "artists": [
        {
          "external_urls": {
            "spotify": "string"
          },
          "href": "string",
          "id": "string",
          "name": "string",
          "type": "artist",
          "uri": "string"
        }
      ],
      "available_markets": [
        "string"
      ],
      "disc_number": 0,
      "duration_ms": 0,
      "explicit": false,
      "external_ids": {
        "isrc": "string",
        "ean": "string",
        "upc": "string"
      },
      "external_urls": {
        "spotify": "string"
      },
      "href": "string",
      "id": "string",
      "is_playable": false,
      "linked_from": {},
      "restrictions": {
        "reason": "string"
      },
      "name": "string",
      "popularity": 0,
      "preview_url": "string",
      "track_number": 0,
      "type": "track",
      "uri": "string",
      "is_local": false
    },
    "track": {
      "album": {
        "album_type": "compilation",
        "total_tracks": 9,
        "available_markets": [
          "CA",
          "BR",
          "IT"
        ],
        "external_urls": {
          "spotify": "string"
        },
        "href": "string",
        "id": "2up3OPMp9Tb4dAKM2erWXQ",
        "images": [
          {
            "url": "https://i.scdn.co/image/ab67616d00001e02ff9ca10b55ce82ae553c8228",
            "height": 300,
            "width": 300
          }
        ],
        "name": "string",
        "release_date": "1981-12",
        "release_date_precision": "year",
        "restrictions": {
          "reason": "market"
        },
        "type": "album",
        "uri": "spotify:album:2up3OPMp9Tb4dAKM2erWXQ",
        "artists": [
          {
            "external_urls": {
              "spotify": "string"
            },
            "href": "string",
            "id": "string",
            "name": "string",
            "type": "artist",
            "uri": "string"
          }
        ]
      },
      "artists": [
        {
          "external_urls": {
            "spotify": "string"
          },
          "href": "string",
          "id": "string",
          "name": "string",
          "type": "artist",
          "uri": "string"
        }
      ],
      "available_markets": [
        "string"
      ],
      "disc_number": 0,
      "duration_ms": 0,
      "explicit": false,
      "external_ids": {
        "isrc": "string",
        "ean": "string",
        "upc": "string"
      },
      "external_urls": {
        "spotify": "string"
      },
      "href": "string",
      "id": "string",
      "is_playable": false,
      "linked_from": {},
      "restrictions": {
        "reason": "string"
      },
      "name": "string",
      "popularity": 0,
      "preview_url": "string",
      "track_number": 0,
      "type": "track",
      "uri": "string",
      "is_local": false
    }
  }
  ]
  ```
</details>

<details>
  <summary>Occur</summary>

  - Get Playlist.items.items[]
  - Get Playlist Items.items[]
  
</details>

### `SavedAlbumObject`
<details>
  <summary>Sample</summary>

  ```bash
  [
  {
    "added_at": "string",
    "album": {
      "album_type": "compilation",
      "total_tracks": 9,
      "available_markets": [
        "CA",
        "BR",
        "IT"
      ],
      "external_urls": {
        "spotify": "string"
      },
      "href": "string",
      "id": "2up3OPMp9Tb4dAKM2erWXQ",
      "images": [
        {
          "url": "https://i.scdn.co/image/ab67616d00001e02ff9ca10b55ce82ae553c8228",
          "height": 300,
          "width": 300
        }
      ],
      "name": "string",
      "release_date": "1981-12",
      "release_date_precision": "year",
      "restrictions": {
        "reason": "market"
      },
      "type": "album",
      "uri": "spotify:album:2up3OPMp9Tb4dAKM2erWXQ",
      "artists": [
        {
          "external_urls": {
            "spotify": "string"
          },
          "href": "string",
          "id": "string",
          "name": "string",
          "type": "artist",
          "uri": "string"
        }
      ],
      "tracks": {
        "href": "https://api.spotify.com/v1/me/shows?offset=0&limit=20",
        "limit": 20,
        "next": "https://api.spotify.com/v1/me/shows?offset=1&limit=1",
        "offset": 0,
        "previous": "https://api.spotify.com/v1/me/shows?offset=1&limit=1",
        "total": 4,
        "items": [
          {
            "artists": [
              {
                "external_urls": {
                  "spotify": "string"
                },
                "href": "string",
                "id": "string",
                "name": "string",
                "type": "artist",
                "uri": "string"
              }
            ],
            "available_markets": [
              "string"
            ],
            "disc_number": 0,
            "duration_ms": 0,
            "explicit": false,
            "external_urls": {
              "spotify": "string"
            },
            "href": "string",
            "id": "string",
            "is_playable": false,
            "linked_from": {
              "external_urls": {
                "spotify": "string"
              },
              "href": "string",
              "id": "string",
              "type": "string",
              "uri": "string"
            },
            "restrictions": {
              "reason": "string"
            },
            "name": "string",
            "preview_url": "string",
            "track_number": 0,
            "type": "string",
            "uri": "string",
            "is_local": false
          }
        ]
      },
      "copyrights": [
        {
          "text": "string",
          "type": "string"
        }
      ],
      "external_ids": {
        "isrc": "string",
        "ean": "string",
        "upc": "string"
      },
      "genres": [],
      "label": "string",
      "popularity": 0
    }
  }
  ]
  ```
</details>

<details>
  <summary>Occur</summary>

  - Get User's Saved Albums.items[]
  
</details>


### `SavedTrackObject`
<details>
  <summary>Sample</summary>

  ```bash
  [
  {
    "added_at": "string",
    "track": {
      "album": {
        "album_type": "compilation",
        "total_tracks": 9,
        "available_markets": [
          "CA",
          "BR",
          "IT"
        ],
        "external_urls": {
          "spotify": "string"
        },
        "href": "string",
        "id": "2up3OPMp9Tb4dAKM2erWXQ",
        "images": [
          {
            "url": "https://i.scdn.co/image/ab67616d00001e02ff9ca10b55ce82ae553c8228",
            "height": 300,
            "width": 300
          }
        ],
        "name": "string",
        "release_date": "1981-12",
        "release_date_precision": "year",
        "restrictions": {
          "reason": "market"
        },
        "type": "album",
        "uri": "spotify:album:2up3OPMp9Tb4dAKM2erWXQ",
        "artists": [
          {
            "external_urls": {
              "spotify": "string"
            },
            "href": "string",
            "id": "string",
            "name": "string",
            "type": "artist",
            "uri": "string"
          }
        ]
      },
      "artists": [
        {
          "external_urls": {
            "spotify": "string"
          },
          "href": "string",
          "id": "string",
          "name": "string",
          "type": "artist",
          "uri": "string"
        }
      ],
      "available_markets": [
        "string"
      ],
      "disc_number": 0,
      "duration_ms": 0,
      "explicit": false,
      "external_ids": {
        "isrc": "string",
        "ean": "string",
        "upc": "string"
      },
      "external_urls": {
        "spotify": "string"
      },
      "href": "string",
      "id": "string",
      "is_playable": false,
      "linked_from": {},
      "restrictions": {
        "reason": "string"
      },
      "name": "string",
      "popularity": 0,
      "preview_url": "string",
      "track_number": 0,
      "type": "track",
      "uri": "string",
      "is_local": false
    }
  }
  ]
  ```
</details>

<details>
  <summary>Occur</summary>

  - Get User's Saved Tracks.items[]
  
</details>

## Endpoints
### Search for Item
**Get Spotify catalog information about albums, artists, playlists, tracks, shows, episodes or audiobooks that match a keyword string**

Endpoint - https://api.spotify.com/v1/search

- q  - search query. The available filters are album, artist, track, year, upc, tag:hipster, tag:new, isrc, and genre.**
- type - comma-separated list of item types to search across.  Allowed values: "album", "artist", "playlist", "track", "show", "episode", "audiobook"
- limit, offset

Responses: Track/Artist/Album'Object (depend on query)

<details>
<summary>JSON Object</summary>

https://api.spotify.com/v1/search?q=the%20eminem%20show&type=album&limit=1



```bash
{
    "albums": {
        "href": "https://api.spotify.com/v1/search?offset=0&limit=1&query=the%20eminem%20show&type=album",
        "limit": 1,
        "next": "https://api.spotify.com/v1/search?offset=1&limit=1&query=the%20eminem%20show&type=album",
        "offset": 0,
        "previous": null,
        "total": 804,
        "items": [
            {
                "album_type": "album",
                "total_tracks": 20,
                "external_urls": {
                    "spotify": "https://open.spotify.com/album/2cWBwpqMsDJC1ZUwz813lo"
                },
                "href": "https://api.spotify.com/v1/albums/2cWBwpqMsDJC1ZUwz813lo",
                "id": "2cWBwpqMsDJC1ZUwz813lo",
                "images": [
                    {
                        "height": 640,
                        "url": "https://i.scdn.co/image/ab67616d0000b2736ca5c90113b30c3c43ffb8f4",
                        "width": 640
                    },
                    {
                        "height": 300,
                        "url": "https://i.scdn.co/image/ab67616d00001e026ca5c90113b30c3c43ffb8f4",
                        "width": 300
                    },
                    {
                        "height": 64,
                        "url": "https://i.scdn.co/image/ab67616d000048516ca5c90113b30c3c43ffb8f4",
                        "width": 64
                    }
                ],
                "name": "The Eminem Show",
                "release_date": "2002-05-26",
                "release_date_precision": "day",
                "type": "album",
                "uri": "spotify:album:2cWBwpqMsDJC1ZUwz813lo",
                "artists": [ 
                    {
                        "external_urls": {
                            "spotify": "https://open.spotify.com/artist/7dGJo4pcD2V6oG8kP0tJRR"
                        },
                        "href": "https://api.spotify.com/v1/artists/7dGJo4pcD2V6oG8kP0tJRR",
                        "id": "7dGJo4pcD2V6oG8kP0tJRR",
                        "name": "Eminem",
                        "type": "artist",
                        "uri": "spotify:artist:7dGJo4pcD2V6oG8kP0tJRR"
                    }
                ]
            }
        ]
    }
}

```

</details>


### Get Album Tracks
Endpoint - https://api.spotify.com/v1/albums/{id}/tracks

Responses: `SimplifiedTrackObject`
<details>
<summary>JSON Object</summary>

https://api.spotify.com/v1/albums/46eOftG6eRzvFa1OKPzmMJ/tracks?limit=2

```bash
{
    "href": "https://api.spotify.com/v1/albums/46eOftG6eRzvFa1OKPzmMJ/tracks?offset=0&limit=2",
    "items": [ //`SimplifiedTrackObject`
        {
            "artists": [ 
                {
                    "external_urls": {
                        "spotify": "https://open.spotify.com/artist/4FZ3j1oH43e7cukCALsCwf"
                    },
                    "href": "https://api.spotify.com/v1/artists/4FZ3j1oH43e7cukCALsCwf",
                    "id": "4FZ3j1oH43e7cukCALsCwf",
                    "name": "Jack White",
                    "type": "artist",
                    "uri": "spotify:artist:4FZ3j1oH43e7cukCALsCwf"
                }
            ],
            "disc_number": 1,
            "duration_ms": 132725,
            "explicit": false,
            "external_urls": {
                "spotify": "https://open.spotify.com/track/6JD6pSfTa6jAumxb0lc9CJ"
            },
            "href": "https://api.spotify.com/v1/tracks/6JD6pSfTa6jAumxb0lc9CJ",
            "id": "6JD6pSfTa6jAumxb0lc9CJ",
            "name": "That's How I'm Feeling - Live",
            "track_number": 1,
            "type": "track",
            "uri": "spotify:track:6JD6pSfTa6jAumxb0lc9CJ",
            "is_local": false
        },
        {
            "artists": [
                {
                    "external_urls": {
                        "spotify": "https://open.spotify.com/artist/4FZ3j1oH43e7cukCALsCwf"
                    },
                    "href": "https://api.spotify.com/v1/artists/4FZ3j1oH43e7cukCALsCwf",
                    "id": "4FZ3j1oH43e7cukCALsCwf",
                    "name": "Jack White",
                    "type": "artist",
                    "uri": "spotify:artist:4FZ3j1oH43e7cukCALsCwf"
                },
                {
                    "external_urls": {
                        "spotify": "https://open.spotify.com/artist/7dGJo4pcD2V6oG8kP0tJRR"
                    },
                    "href": "https://api.spotify.com/v1/artists/7dGJo4pcD2V6oG8kP0tJRR",
                    "id": "7dGJo4pcD2V6oG8kP0tJRR",
                    "name": "Eminem",
                    "type": "artist",
                    "uri": "spotify:artist:7dGJo4pcD2V6oG8kP0tJRR"
                }
            ],
            "disc_number": 1,
            "duration_ms": 138197,
            "explicit": false,
            "external_urls": {
                "spotify": "https://open.spotify.com/track/4XWhHCDrr1r1SPgh9GeQKu"
            },
            "href": "https://api.spotify.com/v1/tracks/4XWhHCDrr1r1SPgh9GeQKu",
            "id": "4XWhHCDrr1r1SPgh9GeQKu",
            "name": "Hello Operator / 'Till I Collapse - Live",
            "track_number": 2,
            "type": "track",
            "uri": "spotify:track:4XWhHCDrr1r1SPgh9GeQKu",
            "is_local": false
        }
    ],
    "limit": 2,
    "next": "https://api.spotify.com/v1/albums/46eOftG6eRzvFa1OKPzmMJ/tracks?offset=2&limit=2",
    "offset": 0,
    "previous": null,
    "total": 3
}
```

</details>

### Get User's Saved Albums
**Get a list of the albums saved in the current Spotify user's 'Your Music' library.**

Authorization scopes - user-library-read

Endpoint - https://api.spotify.com/v1/me/albums

Responses: `SavedAlbumObject`

<details>
<summary>JSON Object</summary>

  https://api.spotify.com/v1/me/albums

  ```bash
  {
    "href": "https://api.spotify.com/v1/me/albums?offset=0&limit=20",
    "items": [
        {
            "added_at": "2026-03-19T09:20:10Z",
            "album": {
                "album_type": "single",
                "total_tracks": 3,
                "external_urls": {
                    "spotify": "https://open.spotify.com/album/46eOftG6eRzvFa1OKPzmMJ"
                },
                "href": "https://api.spotify.com/v1/albums/46eOftG6eRzvFa1OKPzmMJ",
                "id": "46eOftG6eRzvFa1OKPzmMJ",
                "images": [
                    {
                        "url": "https://i.scdn.co/image/ab67616d0000b27347c680e39a6cf1223355394a",
                        "height": 640,
                        "width": 640
                    },
                    {
                        "url": "https://i.scdn.co/image/ab67616d00001e0247c680e39a6cf1223355394a",
                        "height": 300,
                        "width": 300
                    },
                    {
                        "url": "https://i.scdn.co/image/ab67616d0000485147c680e39a6cf1223355394a",
                        "height": 64,
                        "width": 64
                    }
                ],
                "name": "Live at Ford Field",
                "release_date": "2025-11-27",
                "release_date_precision": "day",
                "type": "album",
                "uri": "spotify:album:46eOftG6eRzvFa1OKPzmMJ",
                "artists": [ 
                    {
                        "external_urls": {
                            "spotify": "https://open.spotify.com/artist/4FZ3j1oH43e7cukCALsCwf"
                        },
                        "href": "https://api.spotify.com/v1/artists/4FZ3j1oH43e7cukCALsCwf",
                        "id": "4FZ3j1oH43e7cukCALsCwf",
                        "name": "Jack White",
                        "type": "artist",
                        "uri": "spotify:artist:4FZ3j1oH43e7cukCALsCwf"
                    },
                    {
                        "external_urls": {
                            "spotify": "https://open.spotify.com/artist/7dGJo4pcD2V6oG8kP0tJRR"
                        },
                        "href": "https://api.spotify.com/v1/artists/7dGJo4pcD2V6oG8kP0tJRR",
                        "id": "7dGJo4pcD2V6oG8kP0tJRR",
                        "name": "Eminem",
                        "type": "artist",
                        "uri": "spotify:artist:7dGJo4pcD2V6oG8kP0tJRR"
                    }
                ],
                "tracks": {
                    "href": "https://api.spotify.com/v1/albums/46eOftG6eRzvFa1OKPzmMJ/tracks?offset=0&limit=50",
                    "limit": 50,
                    "next": null,
                    "offset": 0,
                    "previous": null,
                    "total": 3,
                    "items": [
                        {
                            "artists": [ 
                                {
                                    "external_urls": {
                                        "spotify": "https://open.spotify.com/artist/4FZ3j1oH43e7cukCALsCwf"
                                    },
                                    "href": "https://api.spotify.com/v1/artists/4FZ3j1oH43e7cukCALsCwf",
                                    "id": "4FZ3j1oH43e7cukCALsCwf",
                                    "name": "Jack White",
                                    "type": "artist",
                                    "uri": "spotify:artist:4FZ3j1oH43e7cukCALsCwf"
                                }
                            ],
                            "disc_number": 1,
                            "duration_ms": 132725,
                            "explicit": false,
                            "external_urls": {
                                "spotify": "https://open.spotify.com/track/6JD6pSfTa6jAumxb0lc9CJ"
                            },
                            "href": "https://api.spotify.com/v1/tracks/6JD6pSfTa6jAumxb0lc9CJ",
                            "id": "6JD6pSfTa6jAumxb0lc9CJ",
                            "name": "That's How I'm Feeling - Live",
                            "track_number": 1,
                            "type": "track",
                            "uri": "spotify:track:6JD6pSfTa6jAumxb0lc9CJ",
                            "is_local": false
                        },
                        {
                            "artists": [ 
                                {
                                    "external_urls": {
                                        "spotify": "https://open.spotify.com/artist/4FZ3j1oH43e7cukCALsCwf"
                                    },
                                    "href": "https://api.spotify.com/v1/artists/4FZ3j1oH43e7cukCALsCwf",
                                    "id": "4FZ3j1oH43e7cukCALsCwf",
                                    "name": "Jack White",
                                    "type": "artist",
                                    "uri": "spotify:artist:4FZ3j1oH43e7cukCALsCwf"
                                },
                                {
                                    "external_urls": {
                                        "spotify": "https://open.spotify.com/artist/7dGJo4pcD2V6oG8kP0tJRR"
                                    },
                                    "href": "https://api.spotify.com/v1/artists/7dGJo4pcD2V6oG8kP0tJRR",
                                    "id": "7dGJo4pcD2V6oG8kP0tJRR",
                                    "name": "Eminem",
                                    "type": "artist",
                                    "uri": "spotify:artist:7dGJo4pcD2V6oG8kP0tJRR"
                                }
                            ],
                            "disc_number": 1,
                            "duration_ms": 138197,
                            "explicit": false,
                            "external_urls": {
                                "spotify": "https://open.spotify.com/track/4XWhHCDrr1r1SPgh9GeQKu"
                            },
                            "href": "https://api.spotify.com/v1/tracks/4XWhHCDrr1r1SPgh9GeQKu",
                            "id": "4XWhHCDrr1r1SPgh9GeQKu",
                            "name": "Hello Operator / 'Till I Collapse - Live",
                            "track_number": 2,
                            "type": "track",
                            "uri": "spotify:track:4XWhHCDrr1r1SPgh9GeQKu",
                            "is_local": false
                        },
                        {
                            "artists": [ 
                                {
                                    "external_urls": {
                                        "spotify": "https://open.spotify.com/artist/4FZ3j1oH43e7cukCALsCwf"
                                    },
                                    "href": "https://api.spotify.com/v1/artists/4FZ3j1oH43e7cukCALsCwf",
                                    "id": "4FZ3j1oH43e7cukCALsCwf",
                                    "name": "Jack White",
                                    "type": "artist",
                                    "uri": "spotify:artist:4FZ3j1oH43e7cukCALsCwf"
                                }
                            ],
                            "disc_number": 1,
                            "duration_ms": 174773,
                            "explicit": false,
                            "external_urls": {
                                "spotify": "https://open.spotify.com/track/78WxnmK9sOrRB5XOnISlFn"
                            },
                            "href": "https://api.spotify.com/v1/tracks/78WxnmK9sOrRB5XOnISlFn",
                            "id": "78WxnmK9sOrRB5XOnISlFn",
                            "name": "Seven Nation Army - Live",
                            "track_number": 3,
                            "type": "track",
                            "uri": "spotify:track:78WxnmK9sOrRB5XOnISlFn",
                            "is_local": false
                        }
                    ]
                },
                "copyrights": [
                    {
                        "text": "© 2025 Third Face LLC & Marshall B. Mathers III",
                        "type": "C"
                    },
                    {
                        "text": "℗ 2025 Third Face LLC & Marshall B. Mathers III",
                        "type": "P"
                    }
                ],
                "external_ids": {
                    "upc": "00199957165518"
                },
                "genres": []
            }
        }
    ],
    "limit": 20,
    "next": null,
    "offset": 0,
    "previous": null,
    "total": 1
  }
  ```

</details>

### Get Artist's Albums
Endpoint - https://api.spotify.com/v1/artists/{id}/albums

- id - The Spotify ID of the artist
- include_groups - keywords that will be used to filter, by default all album types will be returned. Valid values: album, single appears_on, compilation. For example: include_groups=album,single.

Responses: `SimplifiedAlbumObject`

<details>
<summary>JSON Object</summary>

  https://api.spotify.com/v1/artists/7dGJo4pcD2V6oG8kP0tJRR/albums?include_groups=single,appears_on&limit=2

  ```bash
  {
    "href": "https://api.spotify.com/v1/artists/7dGJo4pcD2V6oG8kP0tJRR/albums?offset=0&limit=2&include_groups=single,appears_on",
    "limit": 2,
    "next": "https://api.spotify.com/v1/artists/7dGJo4pcD2V6oG8kP0tJRR/albums?offset=2&limit=2&include_groups=single,appears_on",
    "offset": 0,
    "previous": null,
    "total": 218,
    "items": [
        {
            "album_type": "single",
            "total_tracks": 3,
            "external_urls": {
                "spotify": "https://open.spotify.com/album/46eOftG6eRzvFa1OKPzmMJ"
            },
            "href": "https://api.spotify.com/v1/albums/46eOftG6eRzvFa1OKPzmMJ",
            "id": "46eOftG6eRzvFa1OKPzmMJ",
            "images": [
                {
                    "url": "https://i.scdn.co/image/ab67616d0000b27347c680e39a6cf1223355394a",
                    "height": 640,
                    "width": 640
                },
                {
                    "url": "https://i.scdn.co/image/ab67616d00001e0247c680e39a6cf1223355394a",
                    "height": 300,
                    "width": 300
                },
                {
                    "url": "https://i.scdn.co/image/ab67616d0000485147c680e39a6cf1223355394a",
                    "height": 64,
                    "width": 64
                }
            ],
            "name": "Live at Ford Field",
            "release_date": "2025-11-27",
            "release_date_precision": "day",
            "type": "album",
            "uri": "spotify:album:46eOftG6eRzvFa1OKPzmMJ",
            "artists": [ 
                {
                    "external_urls": {
                        "spotify": "https://open.spotify.com/artist/4FZ3j1oH43e7cukCALsCwf"
                    },
                    "href": "https://api.spotify.com/v1/artists/4FZ3j1oH43e7cukCALsCwf",
                    "id": "4FZ3j1oH43e7cukCALsCwf",
                    "name": "Jack White",
                    "type": "artist",
                    "uri": "spotify:artist:4FZ3j1oH43e7cukCALsCwf"
                },
                {
                    "external_urls": {
                        "spotify": "https://open.spotify.com/artist/7dGJo4pcD2V6oG8kP0tJRR"
                    },
                    "href": "https://api.spotify.com/v1/artists/7dGJo4pcD2V6oG8kP0tJRR",
                    "id": "7dGJo4pcD2V6oG8kP0tJRR",
                    "name": "Eminem",
                    "type": "artist",
                    "uri": "spotify:artist:7dGJo4pcD2V6oG8kP0tJRR"
                }
            ]
        },
        {
            "album_type": "single",
            "total_tracks": 1,
            "external_urls": {
                "spotify": "https://open.spotify.com/album/30oInr77U6wVk6mmbFJJea"
            },
            "href": "https://api.spotify.com/v1/albums/30oInr77U6wVk6mmbFJJea",
            "id": "30oInr77U6wVk6mmbFJJea",
            "images": [
                {
                    "url": "https://i.scdn.co/image/ab67616d0000b273991d90a963704373daff6883",
                    "height": 640,
                    "width": 640
                },
                {
                    "url": "https://i.scdn.co/image/ab67616d00001e02991d90a963704373daff6883",
                    "height": 300,
                    "width": 300
                },
                {
                    "url": "https://i.scdn.co/image/ab67616d00004851991d90a963704373daff6883",
                    "height": 64,
                    "width": 64
                }
            ],
            "name": "Animals (Pt. I) [with Eminem]",
            "release_date": "2025-07-07",
            "release_date_precision": "day",
            "type": "album",
            "uri": "spotify:album:30oInr77U6wVk6mmbFJJea",
            "artists": [
                {
                    "external_urls": {
                        "spotify": "https://open.spotify.com/artist/6U3ybJ9UHNKEdsH7ktGBZ7"
                    },
                    "href": "https://api.spotify.com/v1/artists/6U3ybJ9UHNKEdsH7ktGBZ7",
                    "id": "6U3ybJ9UHNKEdsH7ktGBZ7",
                    "name": "JID",
                    "type": "artist",
                    "uri": "spotify:artist:6U3ybJ9UHNKEdsH7ktGBZ7"
                },
                {
                    "external_urls": {
                        "spotify": "https://open.spotify.com/artist/7dGJo4pcD2V6oG8kP0tJRR"
                    },
                    "href": "https://api.spotify.com/v1/artists/7dGJo4pcD2V6oG8kP0tJRR",
                    "id": "7dGJo4pcD2V6oG8kP0tJRR",
                    "name": "Eminem",
                    "type": "artist",
                    "uri": "spotify:artist:7dGJo4pcD2V6oG8kP0tJRR"
                }
            ]
        }
    ]
  }
```

</details>

### Get Playlist
**Get a playlist owned by a Spotify user**

Endpoint - https://api.spotify.com/v1/playlists/{playlist_id}

- fields - Filters for the query: a comma-separated list of the fields to return. If omitted, all fields are returned. For example, to get just the playlist''s description and URI: fields=description,uri. A dot separator can be used to specify non-reoccurring fields, while parentheses can be used to specify reoccurring fields within objects. For example, to get just the added date and user ID of the adder: fields=tracks.items(added_at,added_by.id). Use multiple parentheses to drill down into nested objects, for example: fields=tracks.items(track(name,href,album(name,href))). Fields can be excluded by prefixing them with an exclamation mark, for example: fields=tracks.items(track(name,href,album(!name,href)))

Responses: `ImageObject`, `SimplifiedUserObject`, `PlaylistTrackObject`

<details>
<summary>JSON Object</summary>

  https://api.spotify.com/v1/artists/7dGJo4pcD2V6oG8kP0tJRR/albums?include_groups=single,appears_on&limit=2

  ```bash
  {
    "collaborative": false,
    "description": "testttt",
    "external_urls": {
        "spotify": "https://open.spotify.com/playlist/3WqViNalocp0d3lNtuLbfv"
    },
    "followers": {
        "href": null,
        "total": 0
    },
    "href": "https://api.spotify.com/v1/playlists/3WqViNalocp0d3lNtuLbfv",
    "id": "3WqViNalocp0d3lNtuLbfv",
    "images": [
        {
            "height": 640,
            "url": "https://mosaic.scdn.co/640/ab67616d00001e0201c0cd5da820e6128655854fab67616d00001e023d98a0ae7c78a3a9babaf8afab67616d00001e02506c4cc93e5a6234164125e1ab67616d00001e02e6220934dd6c0d6110c45a1e",
            "width": 640
        },
        {
            "height": 300,
            "url": "https://mosaic.scdn.co/300/ab67616d00001e0201c0cd5da820e6128655854fab67616d00001e023d98a0ae7c78a3a9babaf8afab67616d00001e02506c4cc93e5a6234164125e1ab67616d00001e02e6220934dd6c0d6110c45a1e",
            "width": 300
        },
        {
            "height": 60,
            "url": "https://mosaic.scdn.co/60/ab67616d00001e0201c0cd5da820e6128655854fab67616d00001e023d98a0ae7c78a3a9babaf8afab67616d00001e02506c4cc93e5a6234164125e1ab67616d00001e02e6220934dd6c0d6110c45a1e",
            "width": 60
        }
    ],
    "name": "LLLLLL",
    "owner": {
        "display_name": "Amirkhan",
        "external_urls": {
            "spotify": "https://open.spotify.com/user/31yjfrssw7ftaqyydxwf2s5if3hy"
        },
        "href": "https://api.spotify.com/v1/users/31yjfrssw7ftaqyydxwf2s5if3hy",
        "id": "31yjfrssw7ftaqyydxwf2s5if3hy",
        "type": "user",
        "uri": "spotify:user:31yjfrssw7ftaqyydxwf2s5if3hy"
    },
    "primary_color": null,
    "public": true,
    "snapshot_id": "AAAACi6jviBU2UweokTcpCxyD8CIApxj",
    "items": {
        "href": "https://api.spotify.com/v1/playlists/3WqViNalocp0d3lNtuLbfv/items?offset=0&limit=100",
        "items": [
            {
                "added_at": "2026-03-19T01:41:39Z",
                "added_by": {
                    "external_urls": {
                        "spotify": "https://open.spotify.com/user/31yjfrssw7ftaqyydxwf2s5if3hy"
                    },
                    "href": "https://api.spotify.com/v1/users/31yjfrssw7ftaqyydxwf2s5if3hy",
                    "id": "31yjfrssw7ftaqyydxwf2s5if3hy",
                    "type": "user",
                    "uri": "spotify:user:31yjfrssw7ftaqyydxwf2s5if3hy"
                },
                "is_local": false,
                "primary_color": null,
                "item": {
                    "is_playable": true,
                    "explicit": false,
                    "type": "track",
                    "episode": false,
                    "track": true,
                    "album": {
                        "is_playable": true,
                        "type": "album",
                        "album_type": "single",
                        "href": "https://api.spotify.com/v1/albums/2Z1gnUf3nbn6DtwZSUIH54",
                        "id": "2Z1gnUf3nbn6DtwZSUIH54",
                        "images": [
                            {
                                "url": "https://i.scdn.co/image/ab67616d0000b27301c0cd5da820e6128655854f",
                                "width": 640,
                                "height": 640
                            },
                            {
                                "url": "https://i.scdn.co/image/ab67616d00001e0201c0cd5da820e6128655854f",
                                "width": 300,
                                "height": 300
                            },
                            {
                                "url": "https://i.scdn.co/image/ab67616d0000485101c0cd5da820e6128655854f",
                                "width": 64,
                                "height": 64
                            }
                        ],
                        "name": "MILLION DOLLAR BABY",
                        "release_date": "2024-04-26",
                        "release_date_precision": "day",
                        "uri": "spotify:album:2Z1gnUf3nbn6DtwZSUIH54",
                        "artists": [ 
                            {
                                "external_urls": {
                                    "spotify": "https://open.spotify.com/artist/1WaFQSHVGZQJTbf0BdxdNo"
                                },
                                "href": "https://api.spotify.com/v1/artists/1WaFQSHVGZQJTbf0BdxdNo",
                                "id": "1WaFQSHVGZQJTbf0BdxdNo",
                                "name": "Tommy Richman",
                                "type": "artist",
                                "uri": "spotify:artist:1WaFQSHVGZQJTbf0BdxdNo"
                            }
                        ],
                        "external_urls": {
                            "spotify": "https://open.spotify.com/album/2Z1gnUf3nbn6DtwZSUIH54"
                        },
                        "total_tracks": 2
                    },
                    "artists": [ 
                        {
                            "external_urls": {
                                "spotify": "https://open.spotify.com/artist/1WaFQSHVGZQJTbf0BdxdNo"
                            },
                            "href": "https://api.spotify.com/v1/artists/1WaFQSHVGZQJTbf0BdxdNo",
                            "id": "1WaFQSHVGZQJTbf0BdxdNo",
                            "name": "Tommy Richman",
                            "type": "artist",
                            "uri": "spotify:artist:1WaFQSHVGZQJTbf0BdxdNo"
                        }
                    ],
                    "disc_number": 1,
                    "track_number": 1,
                    "duration_ms": 155151,
                    "external_ids": {
                        "isrc": "QM24S2402528"
                    },
                    "external_urls": {
                        "spotify": "https://open.spotify.com/track/5AJ9hqTS2wcFQCELCFRO7A"
                    },
                    "href": "https://api.spotify.com/v1/tracks/5AJ9hqTS2wcFQCELCFRO7A",
                    "id": "5AJ9hqTS2wcFQCELCFRO7A",
                    "name": "MILLION DOLLAR BABY",
                    "uri": "spotify:track:5AJ9hqTS2wcFQCELCFRO7A",
                    "is_local": false
                },
                "video_thumbnail": {
                    "url": null
                }
            },
            {
                "added_at": "2026-03-19T01:41:46Z",
                "added_by": {
                    "external_urls": {
                        "spotify": "https://open.spotify.com/user/31yjfrssw7ftaqyydxwf2s5if3hy"
                    },
                    "href": "https://api.spotify.com/v1/users/31yjfrssw7ftaqyydxwf2s5if3hy",
                    "id": "31yjfrssw7ftaqyydxwf2s5if3hy",
                    "type": "user",
                    "uri": "spotify:user:31yjfrssw7ftaqyydxwf2s5if3hy"
                },
                "is_local": false,
                "primary_color": null,
                "item": {
                    "is_playable": true,
                    "explicit": true,
                    "type": "track",
                    "episode": false,
                    "track": true,
                    "album": {
                        "is_playable": true,
                        "type": "album",
                        "album_type": "album",
                        "href": "https://api.spotify.com/v1/albums/7MZzYkbHL9Tk3O6WeD4Z0Z",
                        "id": "7MZzYkbHL9Tk3O6WeD4Z0Z",
                        "images": [
                            {
                                "url": "https://i.scdn.co/image/ab67616d0000b273506c4cc93e5a6234164125e1",
                                "width": 640,
                                "height": 640
                            },
                            {
                                "url": "https://i.scdn.co/image/ab67616d00001e02506c4cc93e5a6234164125e1",
                                "width": 300,
                                "height": 300
                            },
                            {
                                "url": "https://i.scdn.co/image/ab67616d00004851506c4cc93e5a6234164125e1",
                                "width": 64,
                                "height": 64
                            }
                        ],
                        "name": "Relapse: Refill",
                        "release_date": "2009-05-15",
                        "release_date_precision": "day",
                        "uri": "spotify:album:7MZzYkbHL9Tk3O6WeD4Z0Z",
                        "artists": [ 
                            {
                                "external_urls": {
                                    "spotify": "https://open.spotify.com/artist/7dGJo4pcD2V6oG8kP0tJRR"
                                },
                                "href": "https://api.spotify.com/v1/artists/7dGJo4pcD2V6oG8kP0tJRR",
                                "id": "7dGJo4pcD2V6oG8kP0tJRR",
                                "name": "Eminem",
                                "type": "artist",
                                "uri": "spotify:artist:7dGJo4pcD2V6oG8kP0tJRR"
                            }
                        ],
                        "external_urls": {
                            "spotify": "https://open.spotify.com/album/7MZzYkbHL9Tk3O6WeD4Z0Z"
                        },
                        "total_tracks": 29
                    },
                    "artists": [ 
                        {
                            "external_urls": {
                                "spotify": "https://open.spotify.com/artist/7dGJo4pcD2V6oG8kP0tJRR"
                            },
                            "href": "https://api.spotify.com/v1/artists/7dGJo4pcD2V6oG8kP0tJRR",
                            "id": "7dGJo4pcD2V6oG8kP0tJRR",
                            "name": "Eminem",
                            "type": "artist",
                            "uri": "spotify:artist:7dGJo4pcD2V6oG8kP0tJRR"
                        }
                    ],
                    "disc_number": 1,
                    "track_number": 17,
                    "duration_ms": 392826,
                    "external_ids": {
                        "isrc": "USUM70964089"
                    },
                    "external_urls": {
                        "spotify": "https://open.spotify.com/track/1HR2CTi0ytRJIcik1QKdOa"
                    },
                    "href": "https://api.spotify.com/v1/tracks/1HR2CTi0ytRJIcik1QKdOa",
                    "id": "1HR2CTi0ytRJIcik1QKdOa",
                    "name": "Beautiful",
                    "uri": "spotify:track:1HR2CTi0ytRJIcik1QKdOa",
                    "is_local": false
                },
                "video_thumbnail": {
                    "url": null
                }
            },
            {
                "added_at": "2026-03-19T01:41:58Z",
                "added_by": {
                    "external_urls": {
                        "spotify": "https://open.spotify.com/user/31yjfrssw7ftaqyydxwf2s5if3hy"
                    },
                    "href": "https://api.spotify.com/v1/users/31yjfrssw7ftaqyydxwf2s5if3hy",
                    "id": "31yjfrssw7ftaqyydxwf2s5if3hy",
                    "type": "user",
                    "uri": "spotify:user:31yjfrssw7ftaqyydxwf2s5if3hy"
                },
                "is_local": false,
                "primary_color": null,
                "item": {
                    "is_playable": true,
                    "explicit": false,
                    "type": "track",
                    "episode": false,
                    "track": true,
                    "album": {
                        "is_playable": true,
                        "type": "album",
                        "album_type": "single",
                        "href": "https://api.spotify.com/v1/albums/5V729UqvhwNOcMejx0m55I",
                        "id": "5V729UqvhwNOcMejx0m55I",
                        "images": [
                            {
                                "url": "https://i.scdn.co/image/ab67616d0000b2733d98a0ae7c78a3a9babaf8af",
                                "width": 640,
                                "height": 640
                            },
                            {
                                "url": "https://i.scdn.co/image/ab67616d00001e023d98a0ae7c78a3a9babaf8af",
                                "width": 300,
                                "height": 300
                            },
                            {
                                "url": "https://i.scdn.co/image/ab67616d000048513d98a0ae7c78a3a9babaf8af",
                                "width": 64,
                                "height": 64
                            }
                        ],
                        "name": "NewJeans 'Super Shy'",
                        "release_date": "2023-07-07",
                        "release_date_precision": "day",
                        "uri": "spotify:album:5V729UqvhwNOcMejx0m55I",
                        "artists": [ 
                            {
                                "external_urls": {
                                    "spotify": "https://open.spotify.com/artist/6HvZYsbFfjnjFrWF950C9d"
                                },
                                "href": "https://api.spotify.com/v1/artists/6HvZYsbFfjnjFrWF950C9d",
                                "id": "6HvZYsbFfjnjFrWF950C9d",
                                "name": "NewJeans",
                                "type": "artist",
                                "uri": "spotify:artist:6HvZYsbFfjnjFrWF950C9d"
                            }
                        ],
                        "external_urls": {
                            "spotify": "https://open.spotify.com/album/5V729UqvhwNOcMejx0m55I"
                        },
                        "total_tracks": 2
                    },
                    "artists": [ 
                        {
                            "external_urls": {
                                "spotify": "https://open.spotify.com/artist/6HvZYsbFfjnjFrWF950C9d"
                            },
                            "href": "https://api.spotify.com/v1/artists/6HvZYsbFfjnjFrWF950C9d",
                            "id": "6HvZYsbFfjnjFrWF950C9d",
                            "name": "NewJeans",
                            "type": "artist",
                            "uri": "spotify:artist:6HvZYsbFfjnjFrWF950C9d"
                        }
                    ],
                    "disc_number": 1,
                    "track_number": 1,
                    "duration_ms": 108986,
                    "external_ids": {
                        "isrc": "USA2P2330067"
                    },
                    "external_urls": {
                        "spotify": "https://open.spotify.com/track/6rdkCkjk6D12xRpdMXy0I2"
                    },
                    "href": "https://api.spotify.com/v1/tracks/6rdkCkjk6D12xRpdMXy0I2",
                    "id": "6rdkCkjk6D12xRpdMXy0I2",
                    "name": "New Jeans",
                    "uri": "spotify:track:6rdkCkjk6D12xRpdMXy0I2",
                    "is_local": false
                },
                "video_thumbnail": {
                    "url": null
                }
            },
            {
                "added_at": "2026-03-19T01:42:11Z",
                "added_by": {
                    "external_urls": {
                        "spotify": "https://open.spotify.com/user/31yjfrssw7ftaqyydxwf2s5if3hy"
                    },
                    "href": "https://api.spotify.com/v1/users/31yjfrssw7ftaqyydxwf2s5if3hy",
                    "id": "31yjfrssw7ftaqyydxwf2s5if3hy",
                    "type": "user",
                    "uri": "spotify:user:31yjfrssw7ftaqyydxwf2s5if3hy"
                },
                "is_local": false,
                "primary_color": null,
                "item": {
                    "is_playable": true,
                    "explicit": true,
                    "type": "track",
                    "episode": false,
                    "track": true,
                    "album": {
                        "is_playable": true,
                        "type": "album",
                        "album_type": "single",
                        "href": "https://api.spotify.com/v1/albums/7baqnLVVcQUr5yUhakW9KX",
                        "id": "7baqnLVVcQUr5yUhakW9KX",
                        "images": [
                            {
                                "url": "https://i.scdn.co/image/ab67616d0000b273e6220934dd6c0d6110c45a1e",
                                "width": 640,
                                "height": 640
                            },
                            {
                                "url": "https://i.scdn.co/image/ab67616d00001e02e6220934dd6c0d6110c45a1e",
                                "width": 300,
                                "height": 300
                            },
                            {
                                "url": "https://i.scdn.co/image/ab67616d00004851e6220934dd6c0d6110c45a1e",
                                "width": 64,
                                "height": 64
                            }
                        ],
                        "name": "Lost In Euphoria",
                        "release_date": "2025-03-21",
                        "release_date_precision": "day",
                        "uri": "spotify:album:7baqnLVVcQUr5yUhakW9KX",
                        "artists": [ 
                            {
                                "external_urls": {
                                    "spotify": "https://open.spotify.com/artist/7LVC96BEVGugTAp38AajV6"
                                },
                                "href": "https://api.spotify.com/v1/artists/7LVC96BEVGugTAp38AajV6",
                                "id": "7LVC96BEVGugTAp38AajV6",
                                "name": "Lithe",
                                "type": "artist",
                                "uri": "spotify:artist:7LVC96BEVGugTAp38AajV6"
                            }
                        ],
                        "external_urls": {
                            "spotify": "https://open.spotify.com/album/7baqnLVVcQUr5yUhakW9KX"
                        },
                        "total_tracks": 7
                    },
                    "artists": [ 
                        {
                            "external_urls": {
                                "spotify": "https://open.spotify.com/artist/7LVC96BEVGugTAp38AajV6"
                            },
                            "href": "https://api.spotify.com/v1/artists/7LVC96BEVGugTAp38AajV6",
                            "id": "7LVC96BEVGugTAp38AajV6",
                            "name": "Lithe",
                            "type": "artist",
                            "uri": "spotify:artist:7LVC96BEVGugTAp38AajV6"
                        }
                    ],
                    "disc_number": 1,
                    "track_number": 1,
                    "duration_ms": 161311,
                    "external_ids": {
                        "isrc": "US38Y2524393"
                    },
                    "external_urls": {
                        "spotify": "https://open.spotify.com/track/2f7NXiO2Uyffl4Pp2AArRI"
                    },
                    "href": "https://api.spotify.com/v1/tracks/2f7NXiO2Uyffl4Pp2AArRI",
                    "id": "2f7NXiO2Uyffl4Pp2AArRI",
                    "name": "444",
                    "uri": "spotify:track:2f7NXiO2Uyffl4Pp2AArRI",
                    "is_local": false
                },
                "video_thumbnail": {
                    "url": null
                }
            },
            {
                "added_at": "2026-03-19T01:42:17Z",
                "added_by": {
                    "external_urls": {
                        "spotify": "https://open.spotify.com/user/31yjfrssw7ftaqyydxwf2s5if3hy"
                    },
                    "href": "https://api.spotify.com/v1/users/31yjfrssw7ftaqyydxwf2s5if3hy",
                    "id": "31yjfrssw7ftaqyydxwf2s5if3hy",
                    "type": "user",
                    "uri": "spotify:user:31yjfrssw7ftaqyydxwf2s5if3hy"
                },
                "is_local": false,
                "primary_color": null,
                "item": {
                    "is_playable": true,
                    "explicit": true,
                    "type": "track",
                    "episode": false,
                    "track": true,
                    "album": {
                        "is_playable": true,
                        "type": "album",
                        "album_type": "album",
                        "href": "https://api.spotify.com/v1/albums/7txGsnDSqVMoRl6RQ9XyZP",
                        "id": "7txGsnDSqVMoRl6RQ9XyZP",
                        "images": [
                            {
                                "url": "https://i.scdn.co/image/ab67616d0000b273c4fee55d7b51479627c31f89",
                                "width": 640,
                                "height": 640
                            },
                            {
                                "url": "https://i.scdn.co/image/ab67616d00001e02c4fee55d7b51479627c31f89",
                                "width": 300,
                                "height": 300
                            },
                            {
                                "url": "https://i.scdn.co/image/ab67616d00004851c4fee55d7b51479627c31f89",
                                "width": 64,
                                "height": 64
                            }
                        ],
                        "name": "HEROES & VILLAINS",
                        "release_date": "2022-12-02",
                        "release_date_precision": "day",
                        "uri": "spotify:album:7txGsnDSqVMoRl6RQ9XyZP",
                        "artists": [ 
                            {
                                "external_urls": {
                                    "spotify": "https://open.spotify.com/artist/0iEtIxbK0KxaSlF7G42ZOp"
                                },
                                "href": "https://api.spotify.com/v1/artists/0iEtIxbK0KxaSlF7G42ZOp",
                                "id": "0iEtIxbK0KxaSlF7G42ZOp",
                                "name": "Metro Boomin",
                                "type": "artist",
                                "uri": "spotify:artist:0iEtIxbK0KxaSlF7G42ZOp"
                            }
                        ],
                        "external_urls": {
                            "spotify": "https://open.spotify.com/album/7txGsnDSqVMoRl6RQ9XyZP"
                        },
                        "total_tracks": 15
                    },
                    "artists": [
                        {
                            "external_urls": {
                                "spotify": "https://open.spotify.com/artist/0iEtIxbK0KxaSlF7G42ZOp"
                            },
                            "href": "https://api.spotify.com/v1/artists/0iEtIxbK0KxaSlF7G42ZOp",
                            "id": "0iEtIxbK0KxaSlF7G42ZOp",
                            "name": "Metro Boomin",
                            "type": "artist",
                            "uri": "spotify:artist:0iEtIxbK0KxaSlF7G42ZOp"
                        },
                        {
                            "external_urls": {
                                "spotify": "https://open.spotify.com/artist/0Y5tJX1MQlPlqiwlOH1tJY"
                            },
                            "href": "https://api.spotify.com/v1/artists/0Y5tJX1MQlPlqiwlOH1tJY",
                            "id": "0Y5tJX1MQlPlqiwlOH1tJY",
                            "name": "Travis Scott",
                            "type": "artist",
                            "uri": "spotify:artist:0Y5tJX1MQlPlqiwlOH1tJY"
                        },
                        {
                            "external_urls": {
                                "spotify": "https://open.spotify.com/artist/50co4Is1HCEo8bhOyUWKpn"
                            },
                            "href": "https://api.spotify.com/v1/artists/50co4Is1HCEo8bhOyUWKpn",
                            "id": "50co4Is1HCEo8bhOyUWKpn",
                            "name": "Young Thug",
                            "type": "artist",
                            "uri": "spotify:artist:50co4Is1HCEo8bhOyUWKpn"
                        }
                    ],
                    "disc_number": 1,
                    "track_number": 6,
                    "duration_ms": 194786,
                    "external_ids": {
                        "isrc": "USUG12208787"
                    },
                    "external_urls": {
                        "spotify": "https://open.spotify.com/track/5wG3HvLhF6Y5KTGlK0IW3J"
                    },
                    "href": "https://api.spotify.com/v1/tracks/5wG3HvLhF6Y5KTGlK0IW3J",
                    "id": "5wG3HvLhF6Y5KTGlK0IW3J",
                    "name": "Trance (with Travis Scott & Young Thug)",
                    "uri": "spotify:track:5wG3HvLhF6Y5KTGlK0IW3J",
                    "is_local": false
                },
                "video_thumbnail": {
                    "url": null
                }
            },
            {
                "added_at": "2026-03-19T01:44:48Z",
                "added_by": {
                    "external_urls": {
                        "spotify": "https://open.spotify.com/user/31yjfrssw7ftaqyydxwf2s5if3hy"
                    },
                    "href": "https://api.spotify.com/v1/users/31yjfrssw7ftaqyydxwf2s5if3hy",
                    "id": "31yjfrssw7ftaqyydxwf2s5if3hy",
                    "type": "user",
                    "uri": "spotify:user:31yjfrssw7ftaqyydxwf2s5if3hy"
                },
                "is_local": false,
                "primary_color": null,
                "item": {
                    "is_playable": true,
                    "explicit": true,
                    "type": "track",
                    "episode": false,
                    "track": true,
                    "album": {
                        "is_playable": true,
                        "type": "album",
                        "album_type": "album",
                        "href": "https://api.spotify.com/v1/albums/2cWBwpqMsDJC1ZUwz813lo",
                        "id": "2cWBwpqMsDJC1ZUwz813lo",
                        "images": [
                            {
                                "url": "https://i.scdn.co/image/ab67616d0000b2736ca5c90113b30c3c43ffb8f4",
                                "width": 640,
                                "height": 640
                            },
                            {
                                "url": "https://i.scdn.co/image/ab67616d00001e026ca5c90113b30c3c43ffb8f4",
                                "width": 300,
                                "height": 300
                            },
                            {
                                "url": "https://i.scdn.co/image/ab67616d000048516ca5c90113b30c3c43ffb8f4",
                                "width": 64,
                                "height": 64
                            }
                        ],
                        "name": "The Eminem Show",
                        "release_date": "2002-05-26",
                        "release_date_precision": "day",
                        "uri": "spotify:album:2cWBwpqMsDJC1ZUwz813lo",
                        "artists": [
                            {
                                "external_urls": {
                                    "spotify": "https://open.spotify.com/artist/7dGJo4pcD2V6oG8kP0tJRR"
                                },
                                "href": "https://api.spotify.com/v1/artists/7dGJo4pcD2V6oG8kP0tJRR",
                                "id": "7dGJo4pcD2V6oG8kP0tJRR",
                                "name": "Eminem",
                                "type": "artist",
                                "uri": "spotify:artist:7dGJo4pcD2V6oG8kP0tJRR"
                            }
                        ],
                        "external_urls": {
                            "spotify": "https://open.spotify.com/album/2cWBwpqMsDJC1ZUwz813lo"
                        },
                        "total_tracks": 20
                    },
                    "artists": [
                        {
                            "external_urls": {
                                "spotify": "https://open.spotify.com/artist/7dGJo4pcD2V6oG8kP0tJRR"
                            },
                            "href": "https://api.spotify.com/v1/artists/7dGJo4pcD2V6oG8kP0tJRR",
                            "id": "7dGJo4pcD2V6oG8kP0tJRR",
                            "name": "Eminem",
                            "type": "artist",
                            "uri": "spotify:artist:7dGJo4pcD2V6oG8kP0tJRR"
                        }
                    ],
                    "disc_number": 1,
                    "track_number": 10,
                    "duration_ms": 290320,
                    "external_ids": {
                        "isrc": "USIR10211038"
                    },
                    "external_urls": {
                        "spotify": "https://open.spotify.com/track/7lQ8MOhq6IN2w8EYcFNSUk"
                    },
                    "href": "https://api.spotify.com/v1/tracks/7lQ8MOhq6IN2w8EYcFNSUk",
                    "id": "7lQ8MOhq6IN2w8EYcFNSUk",
                    "name": "Without Me",
                    "uri": "spotify:track:7lQ8MOhq6IN2w8EYcFNSUk",
                    "is_local": false
                },
                "video_thumbnail": {
                    "url": null
                }
            }
        ],
        "limit": 100,
        "next": null,
        "offset": 0,
        "previous": null,
        "total": 6
    },
    "type": "playlist",
    "uri": "spotify:playlist:3WqViNalocp0d3lNtuLbfv"
  }
  ```

</details>


### Get Current User's Playlists
**Get a list of the playlists owned or followed by the current Spotify user**

Authorization scopes - playlist-read-private

Endpoint - https://api.spotify.com/v1/me/playlists

Responses: `SimplifiedPlaylistObject`

<details>
<summary>JSON Object</summary>

  https://api.spotify.com/v1/me/playlists?limit=1

  ```bash
  {
    "href": "https://api.spotify.com/v1/me/playlists?offset=0&limit=1",
    "limit": 1,
    "next": "https://api.spotify.com/v1/me/playlists?offset=1&limit=1",
    "offset": 0,
    "previous": null,
    "total": 4,
    "items": [
        {
            "collaborative": false,
            "description": "orororororro",
            "external_urls": {
                "spotify": "https://open.spotify.com/playlist/4EDuU8w7IGq1YrY3Aq6oKy"
            },
            "href": "https://api.spotify.com/v1/playlists/4EDuU8w7IGq1YrY3Aq6oKy",
            "id": "4EDuU8w7IGq1YrY3Aq6oKy",
            "images": null,
            "name": "star platinum",
            "owner": {
                "display_name": "Amirkhan",
                "external_urls": {
                    "spotify": "https://open.spotify.com/user/31yjfrssw7ftaqyydxwf2s5if3hy"
                },
                "href": "https://api.spotify.com/v1/users/31yjfrssw7ftaqyydxwf2s5if3hy",
                "id": "31yjfrssw7ftaqyydxwf2s5if3hy",
                "type": "user",
                "uri": "spotify:user:31yjfrssw7ftaqyydxwf2s5if3hy"
            },
            "primary_color": null,
            "public": true,
            "snapshot_id": "AAAAAmG5GJ/Kg5cQ7xpoVOXe0pEOuUtH",
            "items": {
                "href": "https://api.spotify.com/v1/playlists/4EDuU8w7IGq1YrY3Aq6oKy/items",
                "total": 0
            },
            "type": "playlist",
            "uri": "spotify:playlist:4EDuU8w7IGq1YrY3Aq6oKy"
        }
   ]
  }
  ```

</details>

### Get Playlist Items
**Get full details of the items of a playlist owned by a Spotify user**

Authorization scopes - playlist-read-private

Endpoint - https://api.spotify.com/v1/playlists/{playlist_id}/items

- fields - Filters for the query: a comma-separated list of the fields to return. If omitted, all fields are returned. For example, to get just the playlist''s description and URI: fields=description,uri. A dot separator can be used to specify non-reoccurring fields, while parentheses can be used to specify reoccurring fields within objects. For example, to get just the added date and user ID of the adder: fields=tracks.items(added_at,added_by.id). Use multiple parentheses to drill down into nested objects, for example: fields=tracks.items(track(name,href,album(name,href))). Fields can be excluded by prefixing them with an exclamation mark, for example: fields=tracks.items(track(name,href,album(!name,href)))

Responses: `PlaylistTrackObject`

<details>
<summary>JSON Object</summary>

  https://api.spotify.com/v1/playlists/3WqViNalocp0d3lNtuLbfv/items?limit=1

  ```bash
  {
    "href": "https://api.spotify.com/v1/playlists/3WqViNalocp0d3lNtuLbfv/items?offset=0&limit=1",
    "items": [
        {
            "added_at": "2026-03-19T01:41:39Z",
            "added_by": {
                "external_urls": {
                    "spotify": "https://open.spotify.com/user/31yjfrssw7ftaqyydxwf2s5if3hy"
                },
                "href": "https://api.spotify.com/v1/users/31yjfrssw7ftaqyydxwf2s5if3hy",
                "id": "31yjfrssw7ftaqyydxwf2s5if3hy",
                "type": "user",
                "uri": "spotify:user:31yjfrssw7ftaqyydxwf2s5if3hy"
            },
            "is_local": false,
            "primary_color": null,
            "item": {
                "is_playable": true,
                "explicit": false,
                "type": "track",
                "episode": false,
                "track": true,
                "album": {
                    "is_playable": true,
                    "type": "album",
                    "album_type": "single",
                    "href": "https://api.spotify.com/v1/albums/2Z1gnUf3nbn6DtwZSUIH54",
                    "id": "2Z1gnUf3nbn6DtwZSUIH54",
                    "images": [
                        {
                            "height": 640,
                            "url": "https://i.scdn.co/image/ab67616d0000b27301c0cd5da820e6128655854f",
                            "width": 640
                        },
                        {
                            "height": 300,
                            "url": "https://i.scdn.co/image/ab67616d00001e0201c0cd5da820e6128655854f",
                            "width": 300
                        },
                        {
                            "height": 64,
                            "url": "https://i.scdn.co/image/ab67616d0000485101c0cd5da820e6128655854f",
                            "width": 64
                        }
                    ],
                    "name": "MILLION DOLLAR BABY",
                    "release_date": "2024-04-26",
                    "release_date_precision": "day",
                    "uri": "spotify:album:2Z1gnUf3nbn6DtwZSUIH54",
                    "artists": [
                        {
                            "external_urls": {
                                "spotify": "https://open.spotify.com/artist/1WaFQSHVGZQJTbf0BdxdNo"
                            },
                            "href": "https://api.spotify.com/v1/artists/1WaFQSHVGZQJTbf0BdxdNo",
                            "id": "1WaFQSHVGZQJTbf0BdxdNo",
                            "name": "Tommy Richman",
                            "type": "artist",
                            "uri": "spotify:artist:1WaFQSHVGZQJTbf0BdxdNo"
                        }
                    ],
                    "external_urls": {
                        "spotify": "https://open.spotify.com/album/2Z1gnUf3nbn6DtwZSUIH54"
                    },
                    "total_tracks": 2
                },
                "artists": [
                    {
                        "external_urls": {
                            "spotify": "https://open.spotify.com/artist/1WaFQSHVGZQJTbf0BdxdNo"
                        },
                        "href": "https://api.spotify.com/v1/artists/1WaFQSHVGZQJTbf0BdxdNo",
                        "id": "1WaFQSHVGZQJTbf0BdxdNo",
                        "name": "Tommy Richman",
                        "type": "artist",
                        "uri": "spotify:artist:1WaFQSHVGZQJTbf0BdxdNo"
                    }
                ],
                "disc_number": 1,
                "track_number": 1,
                "duration_ms": 155151,
                "external_ids": {
                    "isrc": "QM24S2402528"
                },
                "external_urls": {
                    "spotify": "https://open.spotify.com/track/5AJ9hqTS2wcFQCELCFRO7A"
                },
                "href": "https://api.spotify.com/v1/tracks/5AJ9hqTS2wcFQCELCFRO7A",
                "id": "5AJ9hqTS2wcFQCELCFRO7A",
                "name": "MILLION DOLLAR BABY",
                "uri": "spotify:track:5AJ9hqTS2wcFQCELCFRO7A",
                "is_local": false
            },
            "video_thumbnail": {
                "url": null
            }
        }
    ],
    "limit": 1,
    "next": "https://api.spotify.com/v1/playlists/3WqViNalocp0d3lNtuLbfv/items?offset=1&limit=1",
    "offset": 0,
    "previous": null,
    "total": 6
  }
  ```

</details>


### Get User's Saved Tracks
Authorization scopes - user-library-read

Endpoint https://api.spotify.com/v1/me/tracks

Responses: `SavedTrackObject`

<details>
<summary>JSON Object</summary>

  https://api.spotify.com/v1/me/tracks?limit=2

  ```bash
  {
    "href": "https://api.spotify.com/v1/me/tracks?offset=0&limit=2",
    "items": [
        {
            "added_at": "2026-03-19T01:40:37Z",
            "track": {
                "album": {
                    "album_type": "album",
                    "artists": [
                        {
                            "external_urls": {
                                "spotify": "https://open.spotify.com/artist/7dGJo4pcD2V6oG8kP0tJRR"
                            },
                            "href": "https://api.spotify.com/v1/artists/7dGJo4pcD2V6oG8kP0tJRR",
                            "id": "7dGJo4pcD2V6oG8kP0tJRR",
                            "name": "Eminem",
                            "type": "artist",
                            "uri": "spotify:artist:7dGJo4pcD2V6oG8kP0tJRR"
                        }
                    ],
                    "external_urls": {
                        "spotify": "https://open.spotify.com/album/7MZzYkbHL9Tk3O6WeD4Z0Z"
                    },
                    "href": "https://api.spotify.com/v1/albums/7MZzYkbHL9Tk3O6WeD4Z0Z",
                    "id": "7MZzYkbHL9Tk3O6WeD4Z0Z",
                    "images": [
                        {
                            "height": 640,
                            "width": 640,
                            "url": "https://i.scdn.co/image/ab67616d0000b273506c4cc93e5a6234164125e1"
                        },
                        {
                            "height": 300,
                            "width": 300,
                            "url": "https://i.scdn.co/image/ab67616d00001e02506c4cc93e5a6234164125e1"
                        },
                        {
                            "height": 64,
                            "width": 64,
                            "url": "https://i.scdn.co/image/ab67616d00004851506c4cc93e5a6234164125e1"
                        }
                    ],
                    "is_playable": true,
                    "name": "Relapse: Refill",
                    "release_date": "2009-05-15",
                    "release_date_precision": "day",
                    "total_tracks": 29,
                    "type": "album",
                    "uri": "spotify:album:7MZzYkbHL9Tk3O6WeD4Z0Z"
                },
                "artists": [
                    {
                        "external_urls": {
                            "spotify": "https://open.spotify.com/artist/7dGJo4pcD2V6oG8kP0tJRR"
                        },
                        "href": "https://api.spotify.com/v1/artists/7dGJo4pcD2V6oG8kP0tJRR",
                        "id": "7dGJo4pcD2V6oG8kP0tJRR",
                        "name": "Eminem",
                        "type": "artist",
                        "uri": "spotify:artist:7dGJo4pcD2V6oG8kP0tJRR"
                    }
                ],
                "disc_number": 1,
                "duration_ms": 392826,
                "explicit": true,
                "external_ids": {
                    "isrc": "USUM70964089"
                },
                "external_urls": {
                    "spotify": "https://open.spotify.com/track/1HR2CTi0ytRJIcik1QKdOa"
                },
                "href": "https://api.spotify.com/v1/tracks/1HR2CTi0ytRJIcik1QKdOa",
                "id": "1HR2CTi0ytRJIcik1QKdOa",
                "is_local": false,
                "is_playable": true,
                "name": "Beautiful",
                "track_number": 17,
                "type": "track",
                "uri": "spotify:track:1HR2CTi0ytRJIcik1QKdOa"
            }
        },
        {
            "added_at": "2025-07-08T17:45:07Z",
            "track": {
                "album": {
                    "album_type": "album",
                    "artists": [
                        {
                            "external_urls": {
                                "spotify": "https://open.spotify.com/artist/3hjgDpzMEj8wWDo8vXqywg"
                            },
                            "href": "https://api.spotify.com/v1/artists/3hjgDpzMEj8wWDo8vXqywg",
                            "id": "3hjgDpzMEj8wWDo8vXqywg",
                            "name": "Pasosh",
                            "type": "artist",
                            "uri": "spotify:artist:3hjgDpzMEj8wWDo8vXqywg"
                        }
                    ],
                    "external_urls": {
                        "spotify": "https://open.spotify.com/album/0Xa0Qjo2BoqnKZ8k1ZCZyq"
                    },
                    "href": "https://api.spotify.com/v1/albums/0Xa0Qjo2BoqnKZ8k1ZCZyq",
                    "id": "0Xa0Qjo2BoqnKZ8k1ZCZyq",
                    "images": [
                        {
                            "height": 640,
                            "width": 640,
                            "url": "https://i.scdn.co/image/ab67616d0000b2731fc6cbdd9ddfdaf7f4adeb72"
                        },
                        {
                            "height": 300,
                            "width": 300,
                            "url": "https://i.scdn.co/image/ab67616d00001e021fc6cbdd9ddfdaf7f4adeb72"
                        },
                        {
                            "height": 64,
                            "width": 64,
                            "url": "https://i.scdn.co/image/ab67616d000048511fc6cbdd9ddfdaf7f4adeb72"
                        }
                    ],
                    "is_playable": true,
                    "name": "Пасош",
                    "release_date": "2014-03-20",
                    "release_date_precision": "day",
                    "total_tracks": 14,
                    "type": "album",
                    "uri": "spotify:album:0Xa0Qjo2BoqnKZ8k1ZCZyq"
                },
                "artists": [
                    {
                        "external_urls": {
                            "spotify": "https://open.spotify.com/artist/3hjgDpzMEj8wWDo8vXqywg"
                        },
                        "href": "https://api.spotify.com/v1/artists/3hjgDpzMEj8wWDo8vXqywg",
                        "id": "3hjgDpzMEj8wWDo8vXqywg",
                        "name": "Pasosh",
                        "type": "artist",
                        "uri": "spotify:artist:3hjgDpzMEj8wWDo8vXqywg"
                    }
                ],
                "disc_number": 1,
                "duration_ms": 205217,
                "explicit": false,
                "external_ids": {
                    "isrc": "QMFME1998996"
                },
                "external_urls": {
                    "spotify": "https://open.spotify.com/track/1Uy7H0ewneFhYQp3F2CAab"
                },
                "href": "https://api.spotify.com/v1/tracks/1Uy7H0ewneFhYQp3F2CAab",
                "id": "1Uy7H0ewneFhYQp3F2CAab",
                "is_local": false,
                "is_playable": true,
                "name": "я очень устал",
                "track_number": 8,
                "type": "track",
                "uri": "spotify:track:1Uy7H0ewneFhYQp3F2CAab"
            }
        }
    ],
    "limit": 2,
    "next": "https://api.spotify.com/v1/me/tracks?offset=2&limit=2",
    "offset": 0,
    "previous": null,
    "total": 4
  } 
  ```

</details>


### Get Current User's Profile
**Get metadatas about user**

Authorization scopes - user-read-private user-read-email

Endpoint - https://api.spotify.com/v1/me

Responses: `ImageObject`, `SimplifiedUserObject`

<details>
<summary>JSON Object</summary>

  https://api.spotify.com/v1/me/tracks?limit=2

  ```bash
  {
    "country": "KZ",
    "display_name": "Amirkhan",
    "email": "ospanovamirkhan5@gmail.com",
    "explicit_content": {
        "filter_enabled": false,
        "filter_locked": false
    },
    "external_urls": {
        "spotify": "https://open.spotify.com/user/31yjfrssw7ftaqyydxwf2s5if3hy"
    },
    "followers": {
        "href": null,
        "total": 0
    },
    "href": "https://api.spotify.com/v1/users/31yjfrssw7ftaqyydxwf2s5if3hy",
    "id": "31yjfrssw7ftaqyydxwf2s5if3hy",
    "images": [],
    "product": "premium",
    "type": "user",
    "uri": "spotify:user:31yjfrssw7ftaqyydxwf2s5if3hy"
  }
  ```

</details>

### Get User's Top Items
**Get the current user's top artists or tracks based on calculated affinity**

Authorization scopes - user-top-read

Endpoint - https://api.spotify.com/v1/me/top/{type}
- type - artists or tracks
- time_range - set the time range to calculate affinities. Valid values: long_term (calculated from ~1 year of data and including all new data as it becomes available), medium_term (approximately last 6 months), short_term (approximately last 4 weeks). Default: medium_term

Responses: one of ArtistObject or TrackObject

<details>
<summary>JSON Object</summary>

  https://api.spotify.com/v1/me/top/tracks?limit=1

  ```bash
  {
    "items": [
        {
            "album": {
                "album_type": "album",
                "artists": [
                    {
                        "external_urls": {
                            "spotify": "https://open.spotify.com/artist/7dGJo4pcD2V6oG8kP0tJRR"
                        },
                        "href": "https://api.spotify.com/v1/artists/7dGJo4pcD2V6oG8kP0tJRR",
                        "id": "7dGJo4pcD2V6oG8kP0tJRR",
                        "name": "Eminem",
                        "type": "artist",
                        "uri": "spotify:artist:7dGJo4pcD2V6oG8kP0tJRR"
                    }
                ],
                "external_urls": {
                    "spotify": "https://open.spotify.com/album/2cWBwpqMsDJC1ZUwz813lo"
                },
                "href": "https://api.spotify.com/v1/albums/2cWBwpqMsDJC1ZUwz813lo",
                "id": "2cWBwpqMsDJC1ZUwz813lo",
                "images": [
                    {
                        "height": 640,
                        "url": "https://i.scdn.co/image/ab67616d0000b2736ca5c90113b30c3c43ffb8f4",
                        "width": 640
                    },
                    {
                        "height": 300,
                        "url": "https://i.scdn.co/image/ab67616d00001e026ca5c90113b30c3c43ffb8f4",
                        "width": 300
                    },
                    {
                        "height": 64,
                        "url": "https://i.scdn.co/image/ab67616d000048516ca5c90113b30c3c43ffb8f4",
                        "width": 64
                    }
                ],
                "is_playable": true,
                "name": "The Eminem Show",
                "release_date": "2002-05-26",
                "release_date_precision": "day",
                "total_tracks": 20,
                "type": "album",
                "uri": "spotify:album:2cWBwpqMsDJC1ZUwz813lo"
            },
            "artists": [
                {
                    "external_urls": {
                        "spotify": "https://open.spotify.com/artist/7dGJo4pcD2V6oG8kP0tJRR"
                    },
                    "href": "https://api.spotify.com/v1/artists/7dGJo4pcD2V6oG8kP0tJRR",
                    "id": "7dGJo4pcD2V6oG8kP0tJRR",
                    "name": "Eminem",
                    "type": "artist",
                    "uri": "spotify:artist:7dGJo4pcD2V6oG8kP0tJRR"
                }
            ],
            "disc_number": 1,
            "duration_ms": 290320,
            "explicit": true,
            "external_ids": {
                "isrc": "USIR10211038"
            },
            "external_urls": {
                "spotify": "https://open.spotify.com/track/7lQ8MOhq6IN2w8EYcFNSUk"
            },
            "href": "https://api.spotify.com/v1/tracks/7lQ8MOhq6IN2w8EYcFNSUk",
            "id": "7lQ8MOhq6IN2w8EYcFNSUk",
            "is_local": false,
            "is_playable": true,
            "name": "Without Me",
            "track_number": 10,
            "type": "track",
            "uri": "spotify:track:7lQ8MOhq6IN2w8EYcFNSUk"
        }
    ],
    "total": 22,
    "limit": 1,
    "offset": 0,
    "href": "https://api.spotify.com/v1/me/top/tracks?limit=1",
    "next": "https://api.spotify.com/v1/me/top/tracks?offset=1&limit=1",
    "previous": null
  }
  ```

</details>

### Get Followed Artists
**Get the current user's followed artists**

Authorization scopes - user-follow-read

Endpoint - https://api.spotify.com/v1/me/following

Responses: ArtistObject

<details>
<summary>JSON Object</summary>

  https://api.spotify.com/v1/me/following?type=artist&limit=1

  ```bash
  {
    "artists": {
        "href": "https://api.spotify.com/v1/me/following?type=artist&limit=1",
        "limit": 1,
        "next": "https://api.spotify.com/v1/me/following?type=artist&limit=1&after=7dGJo4pcD2V6oG8kP0tJRR",
        "cursors": {
            "after": "7dGJo4pcD2V6oG8kP0tJRR"
        },
        "total": 3,
        "items": [
            {
                "external_urls": {
                    "spotify": "https://open.spotify.com/artist/7dGJo4pcD2V6oG8kP0tJRR"
                },
                "href": "https://api.spotify.com/v1/artists/7dGJo4pcD2V6oG8kP0tJRR",
                "id": "7dGJo4pcD2V6oG8kP0tJRR",
                "images": [
                    {
                        "url": "https://i.scdn.co/image/ab6761610000e5eba00b11c129b27a88fc72f36b",
                        "height": 640,
                        "width": 640
                    },
                    {
                        "url": "https://i.scdn.co/image/ab67616100005174a00b11c129b27a88fc72f36b",
                        "height": 320,
                        "width": 320
                    },
                    {
                        "url": "https://i.scdn.co/image/ab6761610000f178a00b11c129b27a88fc72f36b",
                        "height": 160,
                        "width": 160
                    }
                ],
                "name": "Eminem",
                "type": "artist",
                "uri": "spotify:artist:7dGJo4pcD2V6oG8kP0tJRR"
            }
        ]
    }
  }
  ```
</details>

### Check User's Saved Items
**Check if one or more items are already saved in the current user's library. Accepts Spotify URIs for tracks, albums, episodes, shows, audiobooks, artists, users, and playlists**

Authorization scopes - user-library-read user-follow-read playlist-read-private

Endpoint
https://api.spotify.com/v1/me/library/contains

- uris - a comma-separated list of Spotify URIs. Maximum: 40 URIs.
Supported URI types:
spotify:track:{id}
spotify:album:{id}
spotify:episode:{id}
spotify:show:{id}
spotify:audiobook:{id}
spotify:artist:{id}
spotify:user:{id}
spotify:playlist:{id}
Example: https://api.spotify.com/v1/me/library/contains?uris=spotify%3Atrack%3A7a3LWj5xSFhFRYmztS8wgK%2Cspotify%3Aalbum%3A4aawyAB9vmqN3uQ7FjRGTy%2Cspotify%3Aartist%3A2takcwOaAZWiXQijPHIx7B

#### PUT methods
### Save Items to Library
**Add one or more items to the current user's library. Accepts Spotify URIs for tracks, albums, episodes, shows, audiobooks, users, and playlists**

Authorization scopes - user-library-modify user-follow-modify playlist-modify-public

Endpoint - https://api.spotify.com/v1/me/library

- uris - A comma-separated list of Spotify URIs. Maximum: 40 URIs
Supported URI types:
spotify:track:{id}
spotify:album:{id}
spotify:episode:{id}
spotify:show:{id}
spotify:audiobook:{id}
spotify:user:{id}
spotify:playlist:{id}
Example: https://api.spotify.com/v1/me/library?uris=spotify%3Atrack%3A7a3LWj5xSFhFRYmztS8wgK%2Cspotify%3Aalbum%3A4aawyAB9vmqN3uQ7FjRGTy

### Change Playlist Details
**Change a playlist's name and public/private state. (The user must, of course, own the playlist.)**

Authorization scopes - playlist-modify-public playlist-modify-private

Endpoint - https://api.spotify.com/v1/playlists/{playlist_id}
Request body
{
    "name": "Updated Playlist Name",
    "description": "Updated playlist description",
    "public": false
}

### Update Playlist Items
**Either reorder or replace items in a playlist depending on the request's parameters. To reorder items, include range_start, insert_before, range_length and snapshot_id in the request's body. To replace items, include uris as either a query parameter or in the request's body. Replacing items in a playlist will overwrite its existing items. This operation can be used for replacing or clearing items in a playlist**

Authorization scopes - playlist-modify-public playlist-modify-private

Endpoint - https://api.spotify.com/v1/playlists/{playlist_id}/items
Request body
{
    "range_start": 1,
    "insert_before": 3,
    "range_length": 2
}

### Add Custom Playlist Cover Image
**Replace the image used to represent a specific playlist**

Authorization scopes - ugc-image-upload playlist-modify-public playlist-modify-private

Endpoint - https://api.spotify.com/v1/playlists/{playlist_id}/images

Body image/jpeg
Base64 encoded JPEG image data, maximum payload size is 256 KB.

Example: "/9j/2wCEABoZGSccJz4lJT5CLy8vQkc9Ozs9R0dHR0dHR0dHR0dHR0dHR0dHR0dHR0dHR0dHR0dHR0dHR0dHR0dHR0dHR0cBHCcnMyYzPSYmPUc9Mj1HR0dEREdHR0dHR0dHR0dHR0dHR0dHR0dHR0dHR0dHR0dHR0dHR0dHR0dHR0dHR0dHR//dAAQAAf/uAA5BZG9iZQBkwAAAAAH/wAARCAABAAEDACIAAREBAhEB/8QASwABAQAAAAAAAAAAAAAAAAAAAAYBAQAAAAAAAAAAAAAAAAAAAAAQAQAAAAAAAAAAAAAAAAAAAAARAQAAAAAAAAAAAAAAAAAAAAD/2gAMAwAAARECEQA/AJgAH//Z"


#### POST methods
### Add Items to Playlist
**Add one or more items to a user's playlist**

Authorization scopes - playlist-modify-public playlist-modify-private


Endpoint - https://api.spotify.com/v1/playlists/{playlist_id}/items
Request body
{
    "uris": [
        "string"
    ],
    "position": 0
}
- position - is the position to insert the items, a zero-based index (to place to the first place -> position = 0)
- uris - a comma-separated list of Spotify URIs to add, can be track or episode URIs. For example:
uris=spotify:track:4iV5W9uYEdYUVa79Axb7Rh, spotify:track:1301WleyT98MSxVHPZCA6M, spotify:episode:512ojhOuo1ktJprKbVcKyQ. A maximum of 100 items can be added in one request

### Create Playlist
**Create a playlist for the current Spotify user. (The playlist will be empty until you add tracks.) Each user is generally limited to a maximum of 11000 playlists**

Authorization scopes - playlist-modify-public playlist-modify-private

Endpoint - https://api.spotify.com/v1/me/playlists
Request body
{
    "name": "New Playlist",
    "description": "New playlist description",
    "public": false
}


#### DELETE methods
### Remove Items from Library
**Remove one or more items from the current user's library. Accepts Spotify URIs for tracks, albums, episodes, shows, audiobooks, users, and playlists**

Authorization scopes - user-library-modify user-follow-modify playlist-modify-public

Endpoint - https://api.spotify.com/v1/me/library
- uris - a comma-separated list of Spotify URIs. Maximum: 40 URIs.
Supported URI types:
spotify:track:{id}
spotify:album:{id}
spotify:episode:{id}
spotify:show:{id}
spotify:audiobook:{id}
spotify:user:{id}
spotify:playlist:{id}
Example: https://api.spotify.com/v1/me/library?uris=spotify%3Atrack%3A7a3LWj5xSFhFRYmztS8wgK%2Cspotify%3Aalbum%3A4aawyAB9vmqN3uQ7FjRGTy

### Remove Playlist Items
**Remove one or more items from a user's playlist**

Authorization scopes - playlist-modify-public playlist-modify-private

Endpoint - https://api.spotify.com/v1/playlists/{playlist_id}/items
Request body
{
    "items": [
        {
            "uri": "string"
        }
    ],
    "snapshot_id": "string" 
}
- snapshot_id is the playlist's snapshot ID against which you want to make the changes. The API will validate that the specified items exist and in the specified positions and make the changes, even if more recent changes have been made to the playlist

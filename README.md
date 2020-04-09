# music-share

Small service to share music between different providers (spotify and yandex music as example). You can pass link to a track/album/artist and get links for this entity in other music services (see example below)


Example API request:
```
localhost:4000/link?url=https://music.yandex.ru/artist/312625
```
Response:
```
{
  "spotify": "https://open.spotify.com/artist/7kxOVclB0zQamtBR0syCrg",
  "yandex": "https://music.yandex.ru/artist/312625"
}
```

You need to provide [Spotify credentials](https://developer.spotify.com/documentation/general/guides/authorization-guide/) (see *Client Credentials Flow*)

.env should look like this:
```
SPOTIFY_CREDENTIALS=xxx
```

where `xxx` is your credentials
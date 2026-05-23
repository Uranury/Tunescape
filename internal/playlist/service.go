package playlist

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gitlab.com/Uranury/tunescape/internal/snapshot"
	"gitlab.com/Uranury/tunescape/internal/spotify"
)

type Service interface {
	CreateFromLatestSnapshot(ctx context.Context, userID uuid.UUID) (*Response, error)
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]Playlist, error)
}

type service struct {
	repo        Repository
	snapshotSvc snapshot.Service
	spotifySvc  spotify.Service
}

func NewService(repo Repository, snapshotSvc snapshot.Service, spotifySvc spotify.Service) Service {
	return &service{repo: repo, snapshotSvc: snapshotSvc, spotifySvc: spotifySvc}
}

func (s *service) CreateFromLatestSnapshot(ctx context.Context, userID uuid.UUID) (*Response, error) {
	snap, err := s.snapshotSvc.GetLatestSnapshot(ctx, userID)
	if err != nil {
		return nil, err
	}

	const maxPlaylistTracks = 10

	tracks := snap.Tracks
	if len(tracks) > maxPlaylistTracks {
		tracks = tracks[:maxPlaylistTracks]
	}

	trackURIs := make([]string, 0, len(tracks))
	for _, t := range tracks {
		if t.SpotifyID == "" {
			continue
		}
		trackURIs = append(trackURIs, "spotify:track:"+t.SpotifyID)
	}

	if len(trackURIs) == 0 {
		return nil, fmt.Errorf("no tracks found for playlist creation")
	}

	name := "Tunescape Top Tracks · " + snap.CreatedAt.Format("Jan 2, 2006")

	result, err := s.spotifySvc.CreatePlaylist(ctx, userID, name, trackURIs)
	if err != nil {
		return nil, fmt.Errorf("create spotify playlist: %w", err)
	}

	embedURL := "https://open.spotify.com/embed/playlist/" + result.ID + "?utm_source=generator&theme=0"
	p := &Playlist{
		SpotifyPlaylistID: result.ID,
		Name:              name,
		ExternalURL:       result.ExternalURL,
		EmbedURL:          embedURL,
	}
	if err := s.repo.Insert(ctx, userID, p); err != nil {
		return nil, fmt.Errorf("persist playlist: %w", err)
	}

	return &Response{
		PlaylistID:  result.ID,
		Name:        name,
		ExternalURL: result.ExternalURL,
		EmbedURL:    embedURL,
	}, nil
}

func (s *service) ListByUserID(ctx context.Context, userID uuid.UUID) ([]Playlist, error) {
	return s.repo.ListByUserID(ctx, userID)
}

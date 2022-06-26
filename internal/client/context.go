package client

import "context"

type (
	songKey   struct{}
	statusKey struct{}
)

func ContextWithSong(ctx context.Context, song Song) context.Context {
	return context.WithValue(ctx, songKey{}, song)
}

func ContextHasSong(ctx context.Context) bool {
	return ctx.Value(songKey{}) != nil
}

func SongFromContext(ctx context.Context) Song {
	return ctx.Value(songKey{}).(Song)
}

func ContextWithStatus(ctx context.Context, status Status) context.Context {
	return context.WithValue(ctx, statusKey{}, status)
}

func ContextHasStatus(ctx context.Context) bool {
	return ctx.Value(songKey{}) != nil
}

func StatusFromContext(ctx context.Context) Status {
	return ctx.Value(statusKey{}).(Status)
}

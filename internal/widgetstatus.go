package internal

import "github.com/gdamore/tcell/v2"

// [playing]     Logic/Russ - *Therapy Music* - Vinyl Days (2022)    [rzscxu]
//                 repeat random single consume crossfade update

// var _ Widget = &WidgetStatus{}

type WidgetStatus struct {
	widgetCommon
}

func NewWidgetStatus(events chan<- tcell.Event) *WidgetStatus {
	return nil
}

type WidgetStatusData struct {
	Song   SongType
	Status StatusType
	// Album string
	// Artist string
	// Title string
	// Date string
	// Filename string
	// State string
	// Repeat string
	// Random string
	// Single string
	// Consume string
}

package common

import (
	dcc "github.com/OpenDiablo2/dcc/pkg"
)

type DCCSequenceProvider struct {
	Sequences []*dcc.Direction
}

func (d *DCCSequenceProvider) SequenceCount() int {
	return len(d.Sequences)
}

func (d *DCCSequenceProvider) FrameCount(sequenceId int) int {
	if sequenceId < 0 || sequenceId >= len(d.Sequences) {
		return 0
	}

	return len(d.Sequences[sequenceId].Frames())
}

func (d *DCCSequenceProvider) FrameWidth(sequenceId, frameId int) int {
	if sequenceId < 0 || sequenceId >= len(d.Sequences) {
		return 0
	}

	if frameId < 0 || frameId >= len(d.Sequences[sequenceId].Frames()) {
		return 0
	}

	return d.Sequences[sequenceId].Frames()[frameId].Width
}

func (d *DCCSequenceProvider) FrameHeight(sequenceId, frameId int) int {
	if sequenceId < 0 || sequenceId >= len(d.Sequences) {
		return 0
	}

	if frameId < 0 || frameId >= len(d.Sequences[sequenceId].Frames()) {
		return 0
	}

	return d.Sequences[sequenceId].Frame(frameId).Height
}

func (d *DCCSequenceProvider) GetColorIndexAt(sequenceId, frameId, x, y int) uint8 {
	if sequenceId < 0 || sequenceId >= len(d.Sequences) {
		return 0
	}

	if frameId < 0 || frameId >= len(d.Sequences[sequenceId].Frames()) {
		return 0
	}

	return d.Sequences[sequenceId].Frame(frameId).ColorIndexAt(x, y)
}

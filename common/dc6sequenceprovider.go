package common

import dc6 "github.com/OpenDiablo2/dc6/pkg"

type DC6SequenceProvider struct {
	Sequences []*dc6.Direction
}

func (d *DC6SequenceProvider) SequenceCount() int {
	return len(d.Sequences)
}

func (d *DC6SequenceProvider) FrameCount(sequenceId int) int {
	if sequenceId < 0 || sequenceId >= len(d.Sequences) {
		return 0
	}

	return len(d.Sequences[sequenceId].Frames)
}

func (d *DC6SequenceProvider) FrameWidth(sequenceId, frameId int) int {
	if sequenceId < 0 || sequenceId >= len(d.Sequences) {
		return 0
	}

	if frameId < 0 || frameId >= len(d.Sequences[sequenceId].Frames) {
		return 0
	}

	return int(d.Sequences[sequenceId].Frames[frameId].Width)
}

func (d *DC6SequenceProvider) FrameHeight(sequenceId, frameId int) int {
	if sequenceId < 0 || sequenceId >= len(d.Sequences) {
		return 0
	}

	if frameId < 0 || frameId >= len(d.Sequences[sequenceId].Frames) {
		return 0
	}

	return int(d.Sequences[sequenceId].Frames[frameId].Height)
}

func (d *DC6SequenceProvider) GetColorIndexAt(sequenceId, frameId, x, y int) uint8 {
	if sequenceId < 0 || sequenceId >= len(d.Sequences) {
		return 0
	}

	if frameId < 0 || frameId >= len(d.Sequences[sequenceId].Frames) {
		return 0
	}

	return d.Sequences[sequenceId].Frames[frameId].ColorIndexAt(x, y)
}

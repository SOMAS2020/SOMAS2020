package foraging

// Copy returns a deep copy of the DeerPopulationModel.
func (dp DeerPopulationModel) Copy() DeerPopulationModel {
	ret := dp
	return ret
}

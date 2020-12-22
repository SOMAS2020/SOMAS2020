package disasters

// Copy returns a deep copy of Environment.
func (e Environment) Copy() Environment {
	ret := e
	return ret
}

// Copy returns a deep copy of DisasterReport.
func (d DisasterReport) Copy() DisasterReport {
	ret := d
	return ret
}

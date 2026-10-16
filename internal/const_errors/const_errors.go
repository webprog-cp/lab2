package const_errors

type ConstError string

func (e ConstError) Error() string {
	return string(e)
}

const (
	NoEntityByID        = ConstError("no entity with this ID found")
	InvalidReference    = ConstError("referenced entity does not exist")
	EntityHasDependents = ConstError("entity has dependent records")
	ShelterAtCapacity   = ConstError("shelter is at full capacity")
	CapacityBelowUsage  = ConstError("capacity is below current occupancy")
)

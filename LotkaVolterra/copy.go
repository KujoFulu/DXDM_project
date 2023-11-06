// Copy takes a pointer of Ecosystem object, and returns a new Ecosystem object with the same attributes.
func Copy(ecosystem Ecosystem) *Ecosystem {
	// initialize a new Ecosystem object
	newEcosystem := &Ecosystem{}

	// copy the species slice
	newEcosystem.species = ecosystem.species

	// copy the interaction matrix
	newEcosystem.interaction = ecosystem.interaction

	// copy the deathGrowth matrix
	newEcosystem.deathGrowth = ecosystem.deathGrowth

	return newEcosystem
}
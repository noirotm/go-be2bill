package be2bill

type Amount interface {
	Immediate() bool
	Options() Options
}

type SingleAmount int

func (SingleAmount) Immediate() bool {
	return true
}

func (SingleAmount) Options() Options {
	return nil
}

type FragmentedAmount Options

func (FragmentedAmount) Immediate() bool {
	return false
}

func (p FragmentedAmount) Options() Options {
	return Options(p)
}

package hostenv

func UpateSpecsRev(specsRev string) error {
	upateRev(specsRev, "specs")
	return nil
}

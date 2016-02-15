package outputs

func ShipStdout(event []byte) (err error) {
	Log.Debug(string(event))
	return
}

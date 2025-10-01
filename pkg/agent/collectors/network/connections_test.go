package network

/*
var isSSPresent bool

const ssValidInput = "udp  UNCONN     0      0                                          0.0.0.0:5353                                            0.0.0.0:* "
const ssInvalidInput = "UDP ESTAB 0.0.0.0/24 :5353"

func TestIntegrationArchOpenPorts_Grab(t *testing.T) {

	if _, err := exec.LookPath(systemManagerCommand); err == nil {
		isSSPresent = true
	}

	if !isSSPresent {
		t.Skip("Runnable only on ss based package system")
	}
	if utils.IsCI() {
		t.Skip("[Test skipped] Can't fetch connections in container")
	}

	d := systemManagerPortsGrabber{}
	portsList, err := d.Grab()
	assert.NoError(t, err)
	assert.NotNil(t, portsList)
	assert.NotEmpty(t, portsList)

	raw, err := json.Marshal(portsList)
	assert.NoError(t, err)

	assert.NotNil(t, raw)
	assert.True(t, len(raw) > 100)
}
*/

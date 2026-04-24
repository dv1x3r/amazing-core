package gsf

type Connection struct {
	remoteIP string
	platform Platform
}

func (c *Connection) RemoteIP() string {
	return c.remoteIP
}

func (c *Connection) Platform() Platform {
	return c.platform
}

func (c *Connection) SetPlatform(platform Platform) {
	c.platform = platform
}

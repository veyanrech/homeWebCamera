package camera

type macOSCamera struct {
	DevicesNames   []string
	PicturesFolder string
}

func NewMacOSCamera(dn []string, folder string) Camera {
	return &macOSCamera{
		DevicesNames:   dn,
		PicturesFolder: folder,
	}
}

func (c *macOSCamera) TakePicture() error {
	return nil
}

package metadata

import "testing"

func TestGetVNDBImages(t *testing.T) {
	images, err := GetVNDBImages("v20424")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Covers: %d", len(images.Covers))

	for _, image := range images.Covers {
		t.Logf(
			"Cover: %dx%d %s",
			image.Width,
			image.Height,
			image.URL,
		)
	}

	t.Logf("Backgrounds: %d", len(images.Backgrounds))

	for _, image := range images.Backgrounds {
		t.Logf(
			"Background: %dx%d %s",
			image.Width,
			image.Height,
			image.URL,
		)
	}
}

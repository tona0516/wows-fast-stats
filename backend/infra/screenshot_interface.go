package infra

type ScreenshotInterface interface {
	Save(path string, base64Data string) error
}

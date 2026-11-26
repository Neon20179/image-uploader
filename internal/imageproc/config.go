package imageproc

import "fmt"

type Config struct {
	MaxWidth       int
	MaxHeight      int
	MaxSizeBytes   int
	InitialQuality int
	MinQuality     int
	QualityStep    int
	EncodeMethod   int
}

func (c Config) Validate() error {
	if c.MaxWidth <= 0 || c.MaxHeight <= 0 || c.MaxSizeBytes <= 0 {
		return fmt.Errorf("image dimensions and size limit must be positive")
	}
	if c.MinQuality < 0 || c.InitialQuality > 100 || c.MinQuality > c.InitialQuality || c.QualityStep <= 0 {
		return fmt.Errorf("quality must satisfy 0 <= min <= initial <= 100 and step > 0")
	}
	if c.EncodeMethod < 0 || c.EncodeMethod > 6 {
		return fmt.Errorf("WebP encode method must be between 0 and 6")
	}
	return nil
}

package pkg

import (
	"fmt"

	"github.com/kirklin/go-swd"
)

type Detector struct {
	swd *swd.SWD
}

func NewDetector() (*Detector, error) {
	detector, err := swd.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create sensitive word detector: %w", err)
	}

	customWords := map[string]swd.Category{
		"涉黄":       swd.Pornography,
		"涉政":       swd.Political,
		"赌博词汇":   swd.Gambling,
		"毒品词汇":   swd.Drugs,
		"脏话词汇":   swd.Profanity,
		"歧视词汇":   swd.Discrimination,
		"诈骗词汇":   swd.Scam,
		"自定义词汇": swd.Custom,
	}
	if err := detector.AddWords(customWords); err != nil {
		return nil, fmt.Errorf("failed to add custom words to sensitive word detector: %w", err)
	}

	return &Detector{swd: detector}, nil
}

func (d *Detector) ContainsSensitiveWords(text string) bool {
	return d.swd.Detect(text)
}

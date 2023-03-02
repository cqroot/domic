package configs

import "github.com/cqroot/gmdots/internal/dot"

func RangeDotConfigs(handleFunc func(dotName string, dotConfig dot.DotConfig)) error {
	dots, err := dot.Dots()
	if err != nil {
		return err
	}

	defDotConfigs, err := DefaultDotConfigs()
	if err != nil {
		return err
	}

	for _, dotName := range dots {
		dotConfig, ok := defDotConfigs[dotName]
		if !ok {
			continue
		}

		if dotConfig.Dest == "-" {
			continue
		}

		handleFunc(dotName, dotConfig)
	}

	return nil
}

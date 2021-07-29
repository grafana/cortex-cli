package commands

import (
	"encoding/json"
	"os"

	"github.com/grafana-tools/sdk"
	"github.com/grafana/cortex-tools/pkg/analyse"
	"gopkg.in/alecthomas/kingpin.v2"
)

type DashboardAnalyseCommand struct {
	DashFilesList []string
	outputFile    string
}

func (cmd *DashboardAnalyseCommand) run(k *kingpin.ParseContext) error {
	output := &analyse.MetricsInGrafana{}

	for _, file := range cmd.DashFilesList {
		var board sdk.Board
		buf, err := loadFile(file)
		if err != nil {
			return err
		}
		if err = json.Unmarshal(buf, &board); err != nil {
			return (err)
		}
		parseMetricsInBoard(output, board)
	}

	err := writeOut(output, cmd.outputFile)
	if err != nil {
		return err
	}
	return nil
}

func loadFile(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileinfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	filesize := fileinfo.Size()
	buffer := make([]byte, filesize)

	_, err = file.Read(buffer)
	if err != nil {
		return nil, err
	}

	return buffer, nil
}

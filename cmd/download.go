package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"tapo/cloud/tapo"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var downloadCmd = &cobra.Command{
	Use: "download",
	Run: runDownload,
}

func init() {
	downloadCmd.Flags().Int("since", 1, "Number of days to download")
	downloadCmd.Flags().Int("days", 1, "Number of days to download")

	rootCmd.AddCommand(downloadCmd)
}

func runDownload(cmd *cobra.Command, args []string) {

	token := viper.GetString("TOKEN")
	termID := viper.GetString("TERM_ID")

	since, _ := cmd.Flags().GetInt("since")
	days, _ := cmd.Flags().GetInt("days")

	pwd, _ := os.Getwd()
	ffmpeg, _ := exec.LookPath(".\\ffmpeg")

	client := tapo.NewTapoCareClient("https://euw1-app-tapo-care.i.tplinknbu.com", token, termID)

	startDate := time.Now().Add(time.Hour * -24 * time.Duration(max(since, days)-1))
	endDate := startDate.Add(time.Hour * 24 * time.Duration(days))

	if resp, err := client.GetVideosDevices(); err == nil {

		for _, device := range resp.DeviceList {

			for page := 0; ; page++ {

				req := tapo.ListActivitiesByDateRequest{
					DeviceId: device.DeviceId,
					Page:     page,
					PageSize: 100,

					StartTime: startDate.Format("2006-01-02 00:00:00"),
					EndTime:   endDate.Format("2006-01-02 00:00:00"),

					Source:           "1",
					EventTypeFilters: []tapo.EventType{},
				}

				resp, err := client.ListActivitiesByDate(&req)

				if err != nil {
					log.Fatalln(err)
				}

				for _, activity := range resp.Listing {

					basePath := filepath.Join(pwd, device.Alias)
					filePath := getFilePath(basePath, activity.Event.EventLocalTime, ".mp4")

					if _, err := os.Stat(filePath); os.IsExist(err) {
						continue
					}

					if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
						log.Println(err)
						continue
					}

					url := fmt.Sprintf("%s&token=%s", activity.Event.Data.Video.StreamUrl, token)

					cmd := exec.Command(ffmpeg, "-i", url, filePath)
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr

					if err := cmd.Run(); err != nil {
						log.Println(err)
						continue
					}
				}

				if resp.Total <= (resp.Page+1)*resp.PageSize {
					break
				}
			}
		}
	}
}

package cmd

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"tapo/cloud/tapo"
	"time"

	"github.com/connesc/cipherio"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var downloadOldCmd = &cobra.Command{
	Use: "download-old",
	Run: runDownloadOld,
}

var httpClient = &http.Client{
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}

func init() {
	downloadOldCmd.Flags().Int("since", 1, "Number of days to download")
	downloadOldCmd.Flags().Int("days", 1, "Number of days to download")

	rootCmd.AddCommand(downloadOldCmd)
}

func runDownloadOld(cmd *cobra.Command, args []string) {

	token := viper.GetString("TOKEN")
	termID := viper.GetString("TERM_ID")

	since, _ := cmd.Flags().GetInt("since")
	days, _ := cmd.Flags().GetInt("days")

	pwd, _ := os.Getwd()

	startDate := time.Now().Add(time.Hour * -24 * time.Duration(max(since, days)-1))
	endDate := startDate.Add(time.Hour * 24 * time.Duration(since))

	client := tapo.NewTapoCareClient("https://euw1-app-tapo-care.i.tplinknbu.com", token, termID)

	if resp, err := client.GetVideosDevices(); err != nil {

		for _, device := range resp.DeviceList {

			for page := 0; ; page++ {

				resp, err := client.GetVideosList(&tapo.GetVideosListRequest{
					DeviceId:  device.DeviceId,
					StartTime: startDate.Format("2006-01-02 00:00:00"),
					EndTime:   endDate.Format("2006-01-02 00:00:00"),
					Order:     "desc",
					Page:      page,
					PageSize:  100,
				})

				if err != nil {
					log.Fatalln(err)
				}

				for _, index := range resp.Index {

					for _, video := range index.Video {

						basePath := filepath.Join(pwd, device.Alias)
						filePath := getFilePath(basePath, index.EventLocalTime, ".ts")

						if _, err := os.Stat(filePath); os.IsExist(err) {
							continue
						}

						if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
							log.Println(err)
							continue
						}

						file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)

						if err != nil {
							log.Println(err)
							continue
						}

						defer file.Close()

						_, err = readVideo(video, file)

						if err != nil {
							log.Println(err)
							continue
						}
					}
				}

				if resp.Total <= (resp.Page+1)*resp.PageSize {
					break
				}
			}
		}
	}
}

func getFilePath(basePath string, eventLocalTime string, ext string) string {

	folderName, _, _ := strings.Cut(eventLocalTime, " ")
	fileName := strings.ReplaceAll(eventLocalTime, ":", "-") + ext

	return filepath.Join(basePath, folderName, fileName)
}

func readVideo(video tapo.IndexVideo, writer io.Writer) (int64, error) {

	resp, err := httpClient.Get(video.Uri)

	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	reader, err := getReader(video, resp.Body)

	if err != nil {
		return 0, err
	}

	return io.Copy(writer, reader)
}

func getReader(video tapo.IndexVideo, reader io.Reader) (io.Reader, error) {

	if video.EncryptionMethod == tapo.ENCRYPTION_METHOD_AES_128_CBC {
		return getBlockReader(video.DecryptionInfo, reader)
	}

	return reader, nil
}

func getBlockReader(decryptionInfo tapo.DecryptionInfo, reader io.Reader) (io.Reader, error) {

	blockMode, err := getBlockMode(decryptionInfo, reader)

	if err != nil {
		return nil, err
	}

	return cipherio.NewBlockReader(reader, blockMode), nil
}

func getBlockMode(decryptionInfo tapo.DecryptionInfo, reader io.Reader) (cipher.BlockMode, error) {

	block, err := getCipher(decryptionInfo)

	if err != nil {
		return nil, err
	}

	iv, err := getIV(decryptionInfo, block, reader)

	if err != nil {
		return nil, err
	}

	return cipher.NewCBCDecrypter(block, iv), nil
}

func getCipher(decryptionInfo tapo.DecryptionInfo) (cipher.Block, error) {

	key, err := getKey(decryptionInfo)

	if err != nil {
		return nil, err
	}

	return aes.NewCipher(key)
}

func getKey(decryptionInfo tapo.DecryptionInfo) ([]byte, error) {

	if decryptionInfo.Key != "" {
		return base64.StdEncoding.DecodeString(decryptionInfo.Key)
	}

	resp, err := httpClient.Get(decryptionInfo.KeyUri)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func getIV(decryptionInfo tapo.DecryptionInfo, block cipher.Block, reader io.Reader) ([]byte, error) {

	iv := make([]byte, block.BlockSize())
	n, err := readIV(decryptionInfo, iv, reader)

	if err == nil && n != len(iv) {
		err = errors.New("IV is too short")
	}

	return iv, err
}

func readIV(decryptionInfo tapo.DecryptionInfo, iv []byte, reader io.Reader) (int, error) {

	if decryptionInfo.IV != "" {
		return base64.StdEncoding.Decode(iv, []byte(decryptionInfo.IV))
	}

	return reader.Read(iv)
}

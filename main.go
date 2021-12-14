package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/schollz/progressbar/v3"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"time"
)

const version = "1.0 GO"
const paperEndpoint = "https://papermc.io/api/v2"
const spigetEndpoint = "https://api.spiget.org/v2"
const runOs = runtime.GOOS

func main() {
	log.Println("Using Koper version", version)

	rawArgs := os.Args
	args := os.Args[1:]
	if len(args) == 0 {
		log.Println("Not enough arguments. Usage:", rawArgs[0], "<command> [command args]. Use --help for details")
		return
	}
	switch strings.ToLower(args[0]) {
	case "--help", "-h", "help", "h":
		log.Println(`List of all commands:
		install <plugin_name> [path_to_server_folder] - installs <plugin> to your server's plugins folder. You have to run this in your server folder'
		setup_server <brand (for example paper, spigot, tuinity etc.)> <minecraft_version> <server_name (name of the folder where server must be)]>`)
		break
	case "install":
		argsLen := len(args)
		var pluginsFolder string
		if argsLen == 3 {
			if _, err := os.Stat(args[2]); os.IsNotExist(err){
				log.Fatalln("Invalid server folder path. Your path: ", args[2])
				return
			}
			pluginsFolder = args[2] + string(os.PathSeparator) + "plugins"
			test := os.Mkdir(pluginsFolder, os.ModeDir)
			if test != nil {
				if !(os.IsExist(test) || os.IsNotExist(test)) {
					log.Fatalln("Cannot create plugins folder", test)
				}
			}
		} else if argsLen == 2 {
			if _, err := os.Stat("plugins"); os.IsNotExist(err) {
				log.Fatalln("Running not in server folder or using vanilla software")
				return
			}
			pluginsFolder = "plugins"
		} else {
			log.Fatalln("Not enough arguments. Usage:", rawArgs[0], "<command> [command args]. Use --help for details")
			return
		}
		DownloadPlugin(pluginsFolder, args[1])
		break
	case "setup_server":
		if len(args) != 4 {
			log.Fatalln("Not enough arguments. Usage:", rawArgs[0], "<command> [command args]. Use --help for details")
			return
		}
		fmt.Print("Do you agree with Mojang EULA (https://account.mojang.com/documents/minecraft_eula)? yes/no: ")
		var text string
		_, _ = fmt.Scanln(&text)
		text = strings.ToLower(strings.Replace(text, "\n", "", -1))
		text = strings.Replace(text, " ", "", -1)
		if strings.Compare("no", text) == 0 {
			fmt.Println("Well there is no way ¯\\_(ツ)_/¯")
			return
		} else if strings.Compare("yes", text) != 0 {
			fmt.Println("Unknown option. Please, rerun program and answer correctly.")
			return
		}
		name := args[3]
		log.Println("Installing", args[1], args[2], "to", args[3])
		if runOs == "windows" {
			log.Println("Detected platform: windows")
		} else {
			log.Println("Detected platform: unix-like")
		}
		downloadServer(args[1], args[2], name)
		log.Println("Created startup script")
		log.Println("Created plugins folder")
		os.Mkdir(name + string(os.PathSeparator) + "plugins", os.ModeAppend)
		now := time.Now().Local()
		zone, _ := now.Zone()
		_ = ioutil.WriteFile(name + string(os.PathSeparator) + "eula.txt", []byte(fmt.Sprintf(`#By changing the setting below to TRUE you are indicating your agreement to our EULA (https://account.mojang.com/documents/minecraft_eula).
#You also agree that tacos are tasty, and the best food in the world.
#%s %s %d %d:%d:%d %s %d
eula=true
`, now.Weekday().String(), now.Month().String(), now.Day(), now.Hour(), now.Minute(), now.Second(), zone, now.Year())), 0644)
		log.Println("Created EULA.txt")
		log.Println("Complete setting up")
		break
	}
}

func downloadServer(brand string, version string, name string) {
	if _, err := os.Stat("/" + name); errors.Is(err, os.ErrNotExist) {
		_ = os.Mkdir(name, 0777)
	}
	var fileName string
	switch strings.ToLower(brand) {
	case "paper":
		resp, err := http.Get(paperEndpoint + "/projects/paper/versions/" + version)
		if err != nil {
			log.Fatalln("Unable to fetch latest Paper build on version", version)
			return
		}

		body, err := ioutil.ReadAll(resp.Body)
		var versionInfo map[string]interface{}
		err = json.Unmarshal(body, &versionInfo)
		buildString := fmt.Sprint(versionInfo["builds"])
		buildString = strings.Replace(buildString, "[", "", -1)
		buildString = strings.Replace(buildString, "]", "", -1)
		builds := strings.Split(buildString, " ")
		if len(builds) == 0{
			log.Fatalln("No builds")
			return
		}
		build := builds[len(builds) - 1]
		fileName = "paper-" + version + "-" + build + ".jar"
		DownloadFile(paperEndpoint+ "/projects/paper/versions/" + version +"/builds/" + build +"/downloads/" + fileName, name + string(os.PathSeparator) + fileName)
		break
	case "airplane":
		var airplanePath string
		if version == "1.16" {
			airplanePath = "Airplane-1.16"
		} else if version == "1.17"{
			airplanePath = "Airplane-1.17"
		} else {
			log.Fatalln("Airplane supports only 1.16.5 and 1.17(latest release, like 1.17.1)")
		}
		fileName = "launcher-airplane.jar"
		DownloadFile("https://ci.tivy.ca/job/" + airplanePath + "/lastSuccessfulBuild/artifact/launcher-airplane.jar", name + string(os.PathSeparator) + fileName)
		break
	case "tuinity", "tunity":
		var tuinityPath string
		if version == "1.12.2" {
			tuinityPath = "Tuinity"
		} else if version == "1.17"{
			tuinityPath = "Tuinity-1.17"
		} else {
			log.Fatalln("Tuinity supports only 1.12.2 and 1.17")
		}
		fileName = "tuinity-paperclip.jar"
		DownloadFile("https://ci.codemc.io/job/Spottedleaf/job/" + tuinityPath + "/lastSuccessfulBuild/artifact/tuinity-paperclip.jar", name + string(os.PathSeparator) + fileName)
		break
	case "craftbukkit", "bukkit":
		log.Fatalln("Because getbukkit are greedy they require $20 per month for their api we don't support auto download")
	case "spigot":
		buildToolsPath := "BuildTools.jar"
		os.MkdirAll(name, os.ModeAppend)
		DownloadFile("https://hub.spigotmc.org/jenkins/job/BuildTools/lastSuccessfulBuild/artifact/target/BuildTools.jar", name + string(os.PathSeparator) + buildToolsPath)
		log.Println("Downloading and building latest spigot build for Minecraft version", version, "Please wait, this could take a while. Usually takes ~10 minutes")
		fileName = "spigot-" + version + ".jar"
		cmd := exec.Command("java", "-jar", buildToolsPath, "--rev", version)
		cmd.Dir = name
		stderr, _ := cmd.StderrPipe()
		err := cmd.Start()
		if err != nil {
			return
		}

		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			m := scanner.Text()
			log.Println(m)
		}
		err = cmd.Wait()
		if err != nil {
			log.Fatalln("Failed executing of build")
			return
		}
		break
	default:
		log.Fatalln("Unknown server brand", brand, `Available server brands:
Paper (Recommended for all)
Spigot (Recommenced for small testings)
Tuinity (Recommended for servers on 1.12.2 or 1.17 with 200+ online)
Airplane (Recommended for servers  1.16.5 and 1.17.1 with 200+ online)`)
	}
	os.MkdirAll(name + "/plugins", os.ModeAppend)
	if runOs == "windows" {
		_ = ioutil.WriteFile(name + "/start.bat", []byte(fmt.Sprintf(`@ECHO OFF
:loop
java -Xms8G -Xmx8G -jar %s nogui
TIMEOUT /T 5
goto loop`, fileName)), 0777)

	} else {
		_ = ioutil.WriteFile(name + "/start.sh", []byte(fmt.Sprintf(`#!/bin/sh
            BINDIR=$(dirname "$(readlink -fn "$0")")
            cd "\$BINDIR"
            while true
            do
                java -Xmx2G -DIReallyKnowWhatIAmDoingISwear -jar %s nogui
                echo "Use Ctrl + C for stop server restarting"
                echo "Restart in:"
                for i in 5 4 3 2 1
                do
                    echo "$i..."
                    sleep 1
                done
                echo "RESTART"
            done`, fileName)), 0777)
	}
}

func DownloadFile(url string, dest string) {
	file := path.Base(url)
	start := time.Now()
	err := os.RemoveAll(dest)
	if errors.Is(err, os.ErrExist) {
		log.Println("Found old version, deleting")
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln("Server respond with error", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.45 Safari/537.36")
	resp, _ := http.DefaultClient.Do(req)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalln("Unable to close response body")
		}
	}(resp.Body)
	if resp == nil {
		log.Fatalln("Unable to access", url)
		return
	}

	f, _ := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY, 0644)
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatalln("Unable to close file")
		}
	}(f)

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"Downloading " + file,
	)

	_, err = io.Copy(io.MultiWriter(f, bar), resp.Body)
	if err != nil {
		log.Fatalln("Unable to copy file,", err)
		return
	}

	elapsed := time.Since(start)
	log.Printf("Download completed in %s", elapsed)
}

func DownloadFileW(url string, dest string, fallback string) {
	_ = path.Base(url)
	start := time.Now()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln("Server respond with error", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.45 Safari/537.36")
	resp, _ := http.DefaultClient.Do(req)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalln("Unable to close response body")
		}
	}(resp.Body)
	if resp == nil {
		log.Fatalln("Unable to access", url)
		return
	}
	var realDest string

	var fileName string
	contentDispose := resp.Header.Get("Content-Disposition")
	if contentDispose == ""{
		contentType :=  resp.Header.Get("Content-Type")
		if !strings.Contains(contentType, "java"){
			log.Fatalln("Current plugin has only external source and file from external source is just a webpage. To install plugin go manually search for yours resource. Content-Type:", contentType)
			return
		}
		fileName = fallback
	} else {
		fileName = contentDispose
	}
	realDest = dest + string(os.PathSeparator) + fileName
	f, _ := os.OpenFile(realDest, os.O_CREATE|os.O_WRONLY, 0644)
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatalln("Unable to close file")
		}
	}(f)

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"Downloading " + fallback,
	)

	_, err = io.Copy(io.MultiWriter(f, bar), resp.Body)
	if err != nil {
		log.Fatalln("Unable to copy file,", err)
		return
	}

	elapsed := time.Since(start)
	log.Printf("Download completed in %s", elapsed)
}

func DownloadPlugin(pathToPlugins string, plugin string)  {
	resp, err := http.Get(spigetEndpoint + "/search/resources/" + plugin +"?size=1")
	if err != nil {
		log.Fatalln("Unable to fetch resource with name", plugin)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	var resources []interface{}
	err = json.Unmarshal(body, &resources)
	if err != nil {
		log.Fatalln("Something went wring with json (Resources query)", err)
	}
	if len(resources) == 0 {
		log.Fatalln("No plugins found by query:", plugin)
		return
	}
	var resource map[string]interface{}
	resourceJson := interfaceToJson(resources[0])
	err = json.Unmarshal([]byte(resourceJson), &resource)
	if err != nil {
		log.Fatalln("Something went wring with json (Resource map)", err)
	}
	var fileData map[string]interface{}
	err = json.Unmarshal([]byte(interfaceToJson(resource["file"])), &fileData)
	if err != nil {
		log.Fatalln("Something went wring with json (Resource file info)", err)
	}
	fileType := fmt.Sprint(fileData["type"])
	folderPath := pathToPlugins
	var versionData map[string]interface{}
	err = json.Unmarshal([]byte(interfaceToJson(resource["version"])), &versionData)

	var pluginName string
	if contains(resource, "name"){
		pluginName = fmt.Sprint(resource["name"])
	} else {
		pluginName = plugin
	}
	name := pluginName + "_" + fmt.Sprint(versionData["id"]) + ".jar"
	if contains(resource, "donationLink") {
		log.Println(pluginName, "developer's donation link:", fmt.Sprint(resource["donationLink"]))
	}
	if fileType == "external" {
		DownloadFileW(fmt.Sprint(fileData["externalUrl"]), folderPath, name)
	} else {
		resourceData := strings.Split(fmt.Sprint(fileData["url"]), "/")
		resourceId := strings.Split(resourceData[1], ".")[1]
		DownloadFileW(spigetEndpoint + "/resources/" + string(resourceId) + "/download", folderPath,  name)
	}
}

func contains(actualMap map[string]interface{}, key string) bool {
	_, ok := actualMap[key]
	return ok
}

func createKeyValuePairs(m map[string]interface{}) string {
	b := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(b, "%s=\"%s\"\n", key, fmt.Sprint(value))
	}
	return b.String()
}

func interfaceToJson(m interface{}) string {
	mJson, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	return string(mJson)
}
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"gopkg.in/yaml.v2"
)

type RegexPattern struct {
	Nome  string `yaml:"nome"`
	Regex string `yaml:"regex"`
}

type Config struct {
	RegexPatterns []RegexPattern `yaml:"regex_patterns"`
}

func main() {
	// Declaração de flags de linha de comando
	helpFlag := flag.Bool("h", false, "Mostra a mensagem de ajuda")
	inputFileFlag := flag.String("l", "", "Especifica um arquivo de entrada")
	flag.Parse()

	if *helpFlag {
		showHelp()
		return
	}

	var scanner *bufio.Scanner

	// Verifica se foi fornecido um arquivo de entrada
	if *inputFileFlag != "" {
		file, err := os.Open(*inputFileFlag)
		if err != nil {
			fmt.Printf("Erro ao abrir o arquivo de entrada: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()
		scanner = bufio.NewScanner(file)
	} else {
		scanner = bufio.NewScanner(os.Stdin)
	}

	// Verifica e cria o diretório .refgex no diretório home do usuário
	homeDir, err := getHomeDir()
	if err != nil {
		fmt.Println("Erro ao obter o diretório home:", err)
		os.Exit(1)
	}

	refgexDir := filepath.Join(homeDir, ".refgex")
	err = createDirectory(refgexDir)
	if err != nil {
		fmt.Println("Erro ao criar o diretório .refgex:", err)
		os.Exit(1)
	}

	// Define o caminho do arquivo regex.yml no diretório .refgex
	regexFilePath := filepath.Join(refgexDir, "regex.yml")

	// Verifica se o arquivo regex.yml existe e o baixa se necessário
	if !fileExists(regexFilePath) {
		err := downloadRegexYAML(regexFilePath)
		if err != nil {
			fmt.Println("Erro ao baixar regex.yml:", err)
			os.Exit(1)
		}
	}

	// Lê o arquivo regex.yml com os padrões de regex
	yamlFile, err := ioutil.ReadFile(regexFilePath)
	if err != nil {
		fmt.Println("Erro ao ler o arquivo regex.yml:", err)
		os.Exit(1)
	}

	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		fmt.Println("Erro ao fazer o parsing do arquivo regex.yml:", err)
		os.Exit(1)
	}

	// Compila os padrões de regex
	regexMap := make(map[string]*regexp.Regexp)
	for _, pattern := range config.RegexPatterns {
		compiledRegex, err := regexp.Compile(pattern.Regex)
		if err != nil {
			fmt.Printf("Erro ao compilar o regex '%s': %v\n", pattern.Nome, err)
			continue
		}
		regexMap[pattern.Nome] = compiledRegex
	}

	// Lê as linhas do arquivo de entrada ou entrada padrão
	for scanner.Scan() {
		linha := scanner.Text()
		for nome, regex := range regexMap {
			if regex.MatchString(linha) {
				match := regex.FindString(linha)
				fmt.Printf("Padrão capturado: %s\n", nome)
				fmt.Printf("Texto capturado: %s\n", match)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Erro ao ler a entrada padrão:", err)
		os.Exit(1)
	}
}

func showHelp() {
	fmt.Println("Uso:")
	fmt.Println(" cat seu_arquivo.txt,js,html... |  refgex ")
	fmt.Println("Opções:")
	fmt.Println("  -h           Mostra esta mensagem de ajuda")
	fmt.Println("  -l arquivo.txt  Especifica um arquivo de entrada")
}

func getHomeDir() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	}
	return user.HomeDir, nil
}

func createDirectory(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func downloadRegexYAML(filePath string) error {
	url := "https://raw.githubusercontent.com/eikehacker1/refgex/main/regex.yml"
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Erro ao baixar o arquivo: %s", resp.Status)
	}

	// Lê o corpo da resposta HTTP e o salva no arquivo
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filePath, body, 0644)
	if err != nil {
		return err
	}

	return nil
}

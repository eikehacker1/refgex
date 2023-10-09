# refgex
This tool is capable of picking up patterns in different types of files using regex fed into the regex.yml file that is in your environment variable in Linux
regex example:
```yaml
  - name: Artifactory API Token
     regex: '(?:\s|=|:|"|^)AKC[a-zA-Z0-9]{10,}'
```
And to use it you need to have Language GO installed on your machine:
```bash
sudo apt update -y && sudo apt upgrade -y ; sudo apt install golang
```
install tool:
```bash
go install -v github.com/eikehacker1/refgex@latest 
```
move to binary 
```bash
cp ~/go/bin/refgex /usr/bin
```
use:
```bash
cat file(.js,.html.txt[...]) | refgex 
```
ou:
```bash
refgex -l file.txt 
```

### the file with regex is at:
~/.regex/regex.yaml
#### You can add more regex besides dis downloaded automatically from my github

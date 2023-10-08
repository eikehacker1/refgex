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

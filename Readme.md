# Bulk Github webhook updater

## Download and usage
1. Clone git repo ``https://github.com/TheYkk/webhook-updater.git``
2. Go to cloned folder ``cd webook-updater``
3. Download go imports ``go mod download``
4. Create personel token for add webhooks [Link](https://docs.github.com/en/free-pro-team@latest/github/authenticating-to-github/creating-a-personal-access-token)
5. Start using
```bash
export URL=https://star.theykk.net/webhook 
export SECRET=supersecret
export GITHUB_TOKEN=personeltoken
export GITHUB_USER=TheYkk

go run . 
```
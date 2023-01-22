# Meana server agent
## How to build an agent
```bash
git clone git@github.com:meana-io/meana-agent.git
go build -o ./dist/meana ./cmd/meana-agent/main.go
```

You need to run agent at least once to create initial configuration file:

```bash
cd dist
sudo ./meana
```

```
---------------------
Provide meana config
Enter server address: 
Enter UUID: 
---------------------
```

## How to create agent system.d service
```bash
sudo nano /etc/systemd/system/meana.service
```

```bash
[Unit]
Description=Meana agent

[Service]
User=root
WorkingDirectory=PATH_TO_AGENT_DIR/dist
ExecStart=PATH_TO_AGENT_DIR/dist/meana
Restart=always

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl start meana
sudo systemctl enable meana
```


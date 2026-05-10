# App-TKJ - Deployment Guide

## Prerequisites

- Proxmox VE with LXC support
- LXC 203 (App Server)
- LXC 201 (PostgreSQL Database)
- Go 1.19+ installed on LXC 203

## Database Setup (LXC 201)

1. Connect to PostgreSQL:
```bash
psql -U postgres
```

2. Create database and user:
```sql
CREATE DATABASE db_web_tkj;
CREATE USER app_tkj WITH PASSWORD 'your_secure_password';
GRANT ALL PRIVILEGES ON DATABASE db_web_tkj TO app_tkj;
```

3. Run schema:
```bash
psql -U app_tkj -d db_web_tkj -f schema.sql
```

## Application Deployment (LXC 203)

1. Create application directory:
```bash
sudo mkdir -p /opt/app-tkj
sudo chown -R $USER:$USER /opt/app-tkj
```

2. Copy application files:
```bash
cd /workspace/app-tkj
cp app-tkj /opt/app-tkj/
cp schema.sql /opt/app-tkj/
cp app-tkj.service /etc/systemd/system/
```

3. Create environment file:
```bash
cat > /opt/app-tkj/.env << EOF
PORT=8090
DATABASE_URL=postgres://app_tkj:your_secure_password@192.168.88.201:5432/db_web_tkj
SESSION_KEY=$(openssl rand -hex 32)
EOF
```

4. Enable and start service:
```bash
sudo systemctl daemon-reload
sudo systemctl enable app-tkj
sudo systemctl start app-tkj
```

5. Check status:
```bash
sudo systemctl status app-tkj
curl http://localhost:8090/health
```

## Firewall Configuration

Allow port 8090:
```bash
sudo ufw allow 8090/tcp
```

## Update Instructions

```bash
cd /opt/app-tkj
git pull  # if using git
go build -o app-tkj
sudo systemctl restart app-tkj
```

## Default Credentials

- Username: `admin`
- Password: `admin123`

**⚠️ Change default credentials immediately after first login!**

## Logs

View logs:
```bash
sudo journalctl -u app-tkj -f
```

## Health Check

Endpoint: `http://your-server:8090/health`

Response:
```json
{ "status": "ok" }
```

## Backup

Database backup:
```bash
pg_dump -h 192.168.88.201 -U app_tkj db_web_tkj > backup_$(date +%Y%m%d).sql
```

#!/usr/bin/env sh

CLICKHOUSE_DB="${CLICKHOUSE_DB:-database}";
CLICKHOUSE_USER="${CLICKHOUSE_USER:-user}";
CLICKHOUSE_PASSWORD="${CLICKHOUSE_PASSWORD:-password}";

cat <<EOT >> /etc/clickhouse-server/users.d/user.xml
<yandex>
  <!-- Docs: <https://clickhouse.tech/docs/en/operations/settings/settings_users/> -->
  <users>
    <default>
      <profile>default</profile>
      <networks>
        <ip>::/0</ip>
      </networks>
      <password>default</password>
      <quota>default</quota>
    </default>
  </users>
</yandex>
EOT
#cat /etc/clickhouse-server/users.d/user.xml;

clickhouse-client --query "CREATE DATABASE IF NOT EXISTS ${CLICKHOUSE_DB}";

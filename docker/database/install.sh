#!/bin/bash
psql -U postgres -f /opt/global_chat/sql/install.sql
psql -U marcus -d global_chat -f /opt/global_chat/sql/structure.sql
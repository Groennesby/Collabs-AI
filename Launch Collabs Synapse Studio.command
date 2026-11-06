#!/bin/zsh
cd "/Users/macm1pro/Documents/Collabs Local AI 14b"
export OLLAMA_FLASH_ATTENTION=1
export OLLAMA_KV_CACHE_TYPE=q8_0
./bin/collabs-studio

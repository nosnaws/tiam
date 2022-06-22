#!/bin/bash

arcade="../rules/battlesnake play -g wrapped -m arcade_maze -W 19 -H 21 -n shai -u http://localhost:8082 -n tiam -u http://localhost:8081 -n local -u http://localhost:8080 --hazardDamagePerTurn 100 --output game_out.txt"
wrapped="../rules/battlesnake play -g wrapped -W 11 -H 11 -n shai -u http://localhost:8082 -n tiam -u http://localhost:8081 -n local -u http://localhost:8080 -v --output game_out.txt"
standard="../rules/battlesnake play -g standard -W 11 -H 11 -n shai -u http://localhost:8082 -n tiam -u http://localhost:8081 -n local -u http://localhost:8080 -v --output game_out.txt"
gameResults=()

for i in {0..29}
do
  $(../rules/battlesnake play -g wrapped -m arcade_maze -W 19 -H 21 -n tiam -u http://localhost:8081 -n local -u http://localhost:8080 --hazardDamagePerTurn 100 --output game_out.txt -v)
  gameResults[$i]=$(cat game_out.txt | tail -1)
done

echo "${gameResults}"

for gr in ${gameResults[@]}; do
  echo "${gr}"
done

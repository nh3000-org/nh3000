#EX: tail -f log.file | nhlog   -loglang  eng  -serverip  nats://nats.newhorizons3000.org:4222  -logpattern  [ERR]  -logalias  LOGALIAS
#- serverip - NATS nats://xxxxx.yyy:port
# -logalias - make unique for each instance, become DEVICE.device in NATS
cat makelinux.sh | ./nhlog  -loglang eng -serverip nats://192.168.0.5:4222 -logpattern export -logalias loggertest
